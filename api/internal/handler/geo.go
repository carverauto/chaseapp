package handler

import (
	"io"
	"log/slog"
	"net/http"

	"chaseapp.tv/api/pkg/geojson"
)

// GeoHandler handles geospatial utilities.
type GeoHandler struct {
	logger *slog.Logger
}

// NewGeoHandler creates a new GeoHandler.
func NewGeoHandler(logger *slog.Logger) *GeoHandler {
	return &GeoHandler{logger: logger}
}

// GetBoundingRectangle calculates the minimum bounding rectangle for GeoJSON features.
// POST /api/v1/geo/bounding-rect
func (h *GeoHandler) GetBoundingRectangle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		Error(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	points, err := geojson.ExtractPoints(body)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	rect, err := geojson.MinimumBoundingRectangle(points)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	coords := rectangleToPolygon(rect)
	response := map[string]any{
		"type": "Feature",
		"geometry": map[string]any{
			"type":        "Polygon",
			"coordinates": coords,
		},
		"properties": map[string]any{
			"area": rect.Area,
		},
	}

	JSON(w, http.StatusOK, response)
}

func rectangleToPolygon(rect geojson.Rectangle) [][][]float64 {
	pts := rect.Points
	coords := make([][]float64, 0, 5)
	for _, p := range pts {
		coords = append(coords, []float64{p.X, p.Y})
	}
	// Close the polygon by repeating the first point.
	coords = append(coords, []float64{pts[0].X, pts[0].Y})
	return [][][]float64{coords}
}
