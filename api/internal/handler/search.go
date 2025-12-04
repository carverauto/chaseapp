package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"chaseapp.tv/api/internal/search"
)

// SearchHandler handles search endpoints.
type SearchHandler struct {
	client *search.Client
	logger *slog.Logger
}

// NewSearchHandler creates a new SearchHandler.
func NewSearchHandler(client *search.Client, logger *slog.Logger) *SearchHandler {
	return &SearchHandler{client: client, logger: logger}
}

// Search performs a Typesense query across chases.
// GET /api/v1/search?q=...
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		Error(w, http.StatusBadRequest, "q is required")
		return
	}

	page := parseInt(r.URL.Query().Get("page"), 1)
	limit := parseInt(r.URL.Query().Get("limit"), 20)

	result, err := h.client.Search(r.Context(), query, page, limit)
	if err != nil {
		h.logger.Error("search failed", slog.Any("error", err))
		Error(w, http.StatusBadGateway, "Search failed")
		return
	}

	JSON(w, http.StatusOK, result)
}

func parseInt(val string, def int) int {
	if val == "" {
		return def
	}
	if n, err := strconv.Atoi(val); err == nil {
		return n
	}
	return def
}
