package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"chaseapp.tv/api/internal/model"
	"chaseapp.tv/api/internal/realtime"
	"chaseapp.tv/api/internal/repository"
	"chaseapp.tv/api/pkg/scraper"
)

// StreamHandler handles stream extraction.
type StreamHandler struct {
	repo      *repository.ChaseRepository
	extractor *scraper.Extractor
	publisher *realtime.Publisher
	logger    *slog.Logger
}

// NewStreamHandler creates a new StreamHandler.
func NewStreamHandler(repo *repository.ChaseRepository, extractor *scraper.Extractor, publisher *realtime.Publisher, logger *slog.Logger) *StreamHandler {
	return &StreamHandler{
		repo:      repo,
		extractor: extractor,
		publisher: publisher,
		logger:    logger,
	}
}

type extractStreamsRequest struct {
	ChaseID string   `json:"chase_id"`
	URLs    []string `json:"urls"`
}

type extractStreamsResponse struct {
	Chase         *model.Chase   `json:"chase"`
	Extracted     []model.Stream `json:"extracted"`
	ExistingCount int            `json:"existing_count"`
}

// ExtractStreamURLs scrapes news network pages for stream URLs.
// POST /api/v1/streams/extract
func (h *StreamHandler) ExtractStreamURLs(w http.ResponseWriter, r *http.Request) {
	var req extractStreamsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.ChaseID == "" {
		Error(w, http.StatusBadRequest, "chase_id is required")
		return
	}
	if len(req.URLs) == 0 {
		Error(w, http.StatusBadRequest, "at least one URL is required")
		return
	}

	chaseID, err := uuid.Parse(req.ChaseID)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid chase_id")
		return
	}

	chase, err := h.repo.GetByID(r.Context(), chaseID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			Error(w, http.StatusNotFound, "Chase not found")
			return
		}
		h.logger.Error("failed to load chase", slog.Any("error", err))
		Error(w, http.StatusInternalServerError, "Failed to load chase")
		return
	}

	merged := mergeStreams(chase.Streams, nil)
	var extracted []model.Stream

	for _, url := range req.URLs {
		streams, err := h.extractor.Extract(r.Context(), url)
		if err != nil {
			h.logger.Warn("stream extraction failed", slog.Any("error", err), slog.String("url", url))
			continue
		}
		for _, s := range streams {
			if _, ok := merged[s.URL]; ok {
				continue
			}
			merged[s.URL] = s
			extracted = append(extracted, s)
		}
	}

	if len(extracted) == 0 {
		JSON(w, http.StatusOK, extractStreamsResponse{
			Chase:         chase,
			Extracted:     []model.Stream{},
			ExistingCount: len(chase.Streams),
		})
		return
	}

	// Persist updated streams
	var updatedStreams []model.Stream
	for _, stream := range merged {
		updatedStreams = append(updatedStreams, stream)
	}

	updatedChase, _, err := h.repo.Update(r.Context(), chaseID, model.UpdateChaseInput{
		Streams: updatedStreams,
	})
	if err != nil {
		h.logger.Error("failed to update chase streams", slog.Any("error", err))
		Error(w, http.StatusInternalServerError, "Failed to update chase streams")
		return
	}

	if h.publisher != nil {
		if err := h.publisher.PublishChase(realtime.SubjectChaseUpdated, updatedChase); err != nil {
			h.logger.Warn("failed to publish chase updated event", slog.Any("error", err))
		}
	}

	JSON(w, http.StatusOK, extractStreamsResponse{
		Chase:         updatedChase,
		Extracted:     extracted,
		ExistingCount: len(chase.Streams),
	})
}

func mergeStreams(existing []model.Stream, additional []model.Stream) map[string]model.Stream {
	result := make(map[string]model.Stream)
	for _, s := range existing {
		result[s.URL] = s
	}
	for _, s := range additional {
		result[s.URL] = s
	}
	return result
}
