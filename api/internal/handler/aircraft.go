package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"chaseapp.tv/api/internal/model"
	"chaseapp.tv/api/internal/repository"
	"chaseapp.tv/api/pkg/dbscan"
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
	var input model.ClusterAircraftInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.EpsilonMeters <= 0 {
		Error(w, http.StatusBadRequest, "eps_meters must be greater than zero")
		return
	}
	if input.MinPoints <= 0 {
		input.MinPoints = 3
	}
	if len(input.Points) == 0 {
		JSON(w, http.StatusOK, model.ClusterResponse{Clusters: []model.ClusterResult{}})
		return
	}

	points := make([]dbscan.Point, 0, len(input.Points))
	for _, p := range input.Points {
		// Basic coordinate validation
		if p.Latitude < -90 || p.Latitude > 90 || p.Longitude < -180 || p.Longitude > 180 {
			Error(w, http.StatusBadRequest, "Invalid coordinates in points")
			return
		}

		points = append(points, dbscan.Point{
			ID:  p.ID,
			Lat: p.Latitude,
			Lng: p.Longitude,
			Metadata: map[string]any{
				"point": p,
			},
		})
	}

	clusters := dbscan.ClusterPoints(points, input.EpsilonMeters, input.MinPoints)
	results := make([]model.ClusterResult, 0, len(clusters))

	for _, c := range clusters {
		if len(c.Points) == 0 {
			continue
		}

		var (
			sumLat float64
			sumLng float64
		)

		result := model.ClusterResult{
			ID:     c.ID,
			Points: make([]model.ClusterPoint, 0, len(c.Points)),
			Size:   len(c.Points),
		}

		for _, pt := range c.Points {
			cp, ok := pt.Metadata["point"].(model.ClusterPoint)
			if !ok {
				continue
			}
			result.Points = append(result.Points, cp)
			sumLat += cp.Latitude
			sumLng += cp.Longitude

			if cp.Category == model.AircraftCategoryMedia {
				result.MediaPresent = true
			}
		}

		count := float64(len(result.Points))
		if count > 0 {
			result.CentroidLat = sumLat / count
			result.CentroidLng = sumLng / count
		}

		results = append(results, result)
	}

	h.logger.Info("aircraft clustered",
		slog.Int("clusters", len(results)),
		slog.Int("points", len(input.Points)),
	)

	JSON(w, http.StatusOK, model.ClusterResponse{Clusters: results})
}
