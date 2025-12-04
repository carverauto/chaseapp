package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"chaseapp.tv/api/internal/model"
	"chaseapp.tv/api/internal/repository"
)

// AircraftHandler handles aircraft-related HTTP requests.
type AircraftHandler struct {
	repo   *repository.AircraftRepository
	logger *slog.Logger
}

// NewAircraftHandler creates a new AircraftHandler.
func NewAircraftHandler(repo *repository.AircraftRepository, logger *slog.Logger) *AircraftHandler {
	return &AircraftHandler{
		repo:   repo,
		logger: logger,
	}
}

// List returns a paginated list of aircraft.
// GET /api/v1/aircraft
func (h *AircraftHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	opts := model.AircraftListOptions{
		Page:  1,
		Limit: 50,
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

	if category := r.URL.Query().Get("category"); category != "" {
		opts.Category = model.AircraftCategory(category)
	}

	if clusterID := r.URL.Query().Get("cluster_id"); clusterID != "" {
		opts.ClusterID = clusterID
	}

	if onGround := r.URL.Query().Get("on_ground"); onGround != "" {
		b := onGround == "true" || onGround == "1"
		opts.OnGround = &b
	}

	// Parse bounding box
	if minLat := r.URL.Query().Get("min_lat"); minLat != "" {
		if v, err := strconv.ParseFloat(minLat, 64); err == nil {
			opts.MinLat = &v
		}
	}
	if maxLat := r.URL.Query().Get("max_lat"); maxLat != "" {
		if v, err := strconv.ParseFloat(maxLat, 64); err == nil {
			opts.MaxLat = &v
		}
	}
	if minLng := r.URL.Query().Get("min_lng"); minLng != "" {
		if v, err := strconv.ParseFloat(minLng, 64); err == nil {
			opts.MinLng = &v
		}
	}
	if maxLng := r.URL.Query().Get("max_lng"); maxLng != "" {
		if v, err := strconv.ParseFloat(maxLng, 64); err == nil {
			opts.MaxLng = &v
		}
	}

	result, err := h.repo.List(ctx, opts)
	if err != nil {
		h.logger.Error("failed to list aircraft", slog.Any("error", err))
		Error(w, http.StatusInternalServerError, "Failed to retrieve aircraft")
		return
	}

	JSON(w, http.StatusOK, result)
}

// Cluster performs DBSCAN clustering on aircraft positions.
// POST /api/v1/aircraft/cluster
func (h *AircraftHandler) Cluster(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement DBSCAN clustering
	// This will be implemented in Phase 5 with the pkg/dbscan package
	Error(w, http.StatusNotImplemented, "Clustering not yet implemented")
}
