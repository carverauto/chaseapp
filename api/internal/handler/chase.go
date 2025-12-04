package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"chaseapp.tv/api/internal/middleware"
	"chaseapp.tv/api/internal/model"
	"chaseapp.tv/api/internal/realtime"
	"chaseapp.tv/api/internal/repository"
)

// ChaseHandler handles chase-related HTTP requests.
type ChaseHandler struct {
	repo      *repository.ChaseRepository
	publisher *realtime.Publisher
	logger    *slog.Logger
}

// NewChaseHandler creates a new ChaseHandler.
func NewChaseHandler(repo *repository.ChaseRepository, publisher *realtime.Publisher, logger *slog.Logger) *ChaseHandler {
	return &ChaseHandler{
		repo:      repo,
		publisher: publisher,
		logger:    logger,
	}
}

// List returns a paginated list of chases.
// GET /api/v1/chases
func (h *ChaseHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	opts := model.ChaseListOptions{
		Page:  1,
		Limit: 20,
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			opts.Page = p
		}
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			opts.Limit = l
		}
	}

	if live := r.URL.Query().Get("live"); live != "" {
		b := live == "true" || live == "1"
		opts.Live = &b
	}

	if chaseType := r.URL.Query().Get("type"); chaseType != "" {
		opts.ChaseType = model.ChaseType(chaseType)
	}

	if city := r.URL.Query().Get("city"); city != "" {
		opts.City = city
	}

	if state := r.URL.Query().Get("state"); state != "" {
		opts.State = state
	}

	result, err := h.repo.List(ctx, opts)
	if err != nil {
		h.logger.Error("failed to list chases", slog.Any("error", err))
		Error(w, http.StatusInternalServerError, "Failed to retrieve chases")
		return
	}

	JSON(w, http.StatusOK, result)
}

// Create creates a new chase.
// POST /api/v1/chases
func (h *ChaseHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input model.CreateChaseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if input.Title == "" {
		Error(w, http.StatusBadRequest, "Title is required")
		return
	}
	if input.ChaseType == "" {
		Error(w, http.StatusBadRequest, "Chase type is required")
		return
	}

	// Get user ID from context if authenticated
	var createdBy *uuid.UUID
	if userIDStr, ok := ctx.Value(middleware.UserIDKey).(string); ok && userIDStr != "" {
		if uid, err := uuid.Parse(userIDStr); err == nil {
			createdBy = &uid
		}
	}

	chase, err := h.repo.Create(ctx, input, createdBy)
	if err != nil {
		h.logger.Error("failed to create chase",
			slog.Any("error", err),
			slog.String("title", input.Title),
		)
		Error(w, http.StatusInternalServerError, "Failed to create chase")
		return
	}

	h.logger.Info("chase created",
		slog.String("id", chase.ID.String()),
		slog.String("title", chase.Title),
	)

	h.publishChaseEvent(realtime.SubjectChaseCreated, chase)
	if chase.Live {
		h.publishChaseEvent(realtime.SubjectChaseLive, chase)
	}

	JSON(w, http.StatusCreated, chase)
}

// Get retrieves a single chase by ID.
// GET /api/v1/chases/{id}
func (h *ChaseHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid chase ID")
		return
	}

	chase, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			Error(w, http.StatusNotFound, "Chase not found")
			return
		}
		h.logger.Error("failed to get chase",
			slog.Any("error", err),
			slog.String("id", id.String()),
		)
		Error(w, http.StatusInternalServerError, "Failed to retrieve chase")
		return
	}

	// Optionally increment view count
	if r.URL.Query().Get("track_view") == "true" {
		go func() {
			if err := h.repo.IncrementViewCount(ctx, id); err != nil {
				h.logger.Warn("failed to increment view count",
					slog.Any("error", err),
					slog.String("id", id.String()),
				)
			}
		}()
	}

	JSON(w, http.StatusOK, chase)
}

// Update updates an existing chase.
// PUT /api/v1/chases/{id}
func (h *ChaseHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid chase ID")
		return
	}

	var input model.UpdateChaseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	chase, wasLive, err := h.repo.Update(ctx, id, input)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			Error(w, http.StatusNotFound, "Chase not found")
			return
		}
		h.logger.Error("failed to update chase",
			slog.Any("error", err),
			slog.String("id", id.String()),
		)
		Error(w, http.StatusInternalServerError, "Failed to update chase")
		return
	}

	h.logger.Info("chase updated",
		slog.String("id", chase.ID.String()),
		slog.Bool("live", chase.Live),
	)

	h.publishChaseEvent(realtime.SubjectChaseUpdated, chase)

	if !wasLive && chase.Live {
		h.publishChaseEvent(realtime.SubjectChaseLive, chase)
	}

	if wasLive && !chase.Live {
		h.publishChaseEvent(realtime.SubjectChaseEnded, chase)
	}

	JSON(w, http.StatusOK, chase)
}

// Delete soft-deletes a chase.
// DELETE /api/v1/chases/{id}
func (h *ChaseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid chase ID")
		return
	}

	chase, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			Error(w, http.StatusNotFound, "Chase not found")
			return
		}
		h.logger.Error("failed to load chase before delete",
			slog.Any("error", err),
			slog.String("id", id.String()),
		)
		Error(w, http.StatusInternalServerError, "Failed to delete chase")
		return
	}

	if err := h.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			Error(w, http.StatusNotFound, "Chase not found")
			return
		}
		h.logger.Error("failed to delete chase",
			slog.Any("error", err),
			slog.String("id", id.String()),
		)
		Error(w, http.StatusInternalServerError, "Failed to delete chase")
		return
	}

	h.logger.Info("chase deleted", slog.String("id", id.String()))

	h.publishChaseEvent(realtime.SubjectChaseDeleted, chase)

	w.WriteHeader(http.StatusNoContent)
}

// GetBundle returns an offline data bundle of recent chases.
// GET /api/v1/chases/bundle
func (h *ChaseHandler) GetBundle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get live chases
	liveChases, err := h.repo.GetLiveChases(ctx)
	if err != nil {
		h.logger.Error("failed to get live chases", slog.Any("error", err))
		Error(w, http.StatusInternalServerError, "Failed to generate bundle")
		return
	}

	// Get recent ended chases
	recentOpts := model.ChaseListOptions{
		Page:  1,
		Limit: 50,
	}
	recentResult, err := h.repo.List(ctx, recentOpts)
	if err != nil {
		h.logger.Error("failed to get recent chases", slog.Any("error", err))
		Error(w, http.StatusInternalServerError, "Failed to generate bundle")
		return
	}

	bundle := map[string]any{
		"live":   liveChases,
		"recent": recentResult.Chases,
		"meta": map[string]any{
			"live_count":   len(liveChases),
			"recent_count": len(recentResult.Chases),
		},
	}

	JSON(w, http.StatusOK, bundle)
}

// IncrementShare increments the share count for a chase.
// POST /api/v1/chases/{id}/share
func (h *ChaseHandler) IncrementShare(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid chase ID")
		return
	}

	if err := h.repo.IncrementShareCount(ctx, id); err != nil {
		h.logger.Error("failed to increment share count",
			slog.Any("error", err),
			slog.String("id", id.String()),
		)
		Error(w, http.StatusInternalServerError, "Failed to record share")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ChaseHandler) publishChaseEvent(subject string, chase *model.Chase) {
	if h.publisher == nil || chase == nil {
		return
	}

	if err := h.publisher.PublishChase(subject, chase); err != nil {
		h.logger.Error("failed to publish chase event",
			slog.Any("error", err),
			slog.String("subject", subject),
			slog.String("chase_id", chase.ID.String()),
			slog.Time("occurred_at", time.Now().UTC()),
		)
	}
}
