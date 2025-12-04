package handler

import (
	"net/http"
)

// GetChatToken generates a JWT token for the chat service.
// POST /api/v1/auth/chat-token
func GetChatToken(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement JWT generation for chat service
	Error(w, http.StatusNotImplemented, "Not implemented")
}
