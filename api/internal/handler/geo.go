package handler

import (
	"net/http"
)

// GetBoundingRectangle calculates the minimum bounding rectangle for GeoJSON features.
// POST /api/v1/geo/bounding-rect
func GetBoundingRectangle(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement rotating calipers algorithm
	Error(w, http.StatusNotImplemented, "Not implemented")
}
