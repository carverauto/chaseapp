package handler

import (
	"net/http"
)

// ListAircraft returns aircraft within a bounding box.
// GET /api/v1/aircraft?bounds=lat1,lon1,lat2,lon2
func ListAircraft(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement with repository
	JSON(w, http.StatusOK, map[string]any{
		"aircraft": []any{},
		"total":    0,
	})
}

// ClusterAircraft applies DBSCAN clustering to aircraft data.
// POST /api/v1/aircraft/cluster
func ClusterAircraft(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement DBSCAN clustering
	Error(w, http.StatusNotImplemented, "Not implemented")
}
