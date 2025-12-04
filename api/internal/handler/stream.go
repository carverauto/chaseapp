package handler

import (
	"net/http"
)

// ExtractStreamURLs scrapes news network pages for stream URLs.
// POST /api/v1/streams/extract
func ExtractStreamURLs(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement web scraping with Colly
	Error(w, http.StatusNotImplemented, "Not implemented")
}
