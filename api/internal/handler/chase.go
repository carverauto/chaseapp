package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

// ListChases returns a paginated list of chases.
// GET /api/v1/chases
func ListChases(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement with repository
	JSON(w, http.StatusOK, map[string]any{
		"chases": []any{},
		"total":  0,
		"page":   1,
		"limit":  20,
	})
}

// CreateChase creates a new chase.
// POST /api/v1/chases
func CreateChase(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement with repository
	Error(w, http.StatusNotImplemented, "Not implemented")
}

// GetChase returns a single chase by ID.
// GET /api/v1/chases/{id}
func GetChase(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// TODO: Implement with repository
	_ = id
	Error(w, http.StatusNotImplemented, "Not implemented")
}

// UpdateChase updates an existing chase.
// PUT /api/v1/chases/{id}
func UpdateChase(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// TODO: Implement with repository
	_ = id
	Error(w, http.StatusNotImplemented, "Not implemented")
}

// DeleteChase deletes a chase.
// DELETE /api/v1/chases/{id}
func DeleteChase(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// TODO: Implement with repository
	_ = id
	Error(w, http.StatusNotImplemented, "Not implemented")
}

// GetChasesBundle returns an offline data bundle.
// GET /api/v1/chases/bundle
func GetChasesBundle(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement bundle generation
	Error(w, http.StatusNotImplemented, "Not implemented")
}
