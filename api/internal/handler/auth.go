package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"chaseapp.tv/api/internal/auth"
	"chaseapp.tv/api/internal/middleware"
)

// AuthHandler handles authentication-related endpoints.
type AuthHandler struct {
	signer *auth.ChatTokenSigner
	logger *slog.Logger
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(signer *auth.ChatTokenSigner, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		signer: signer,
		logger: logger,
	}
}

type chatTokenRequest struct {
	Permissions []string `json:"permissions"`
}

type chatTokenResponse struct {
	Token string `json:"token"`
}

// GetChatToken generates a JWT token for the chat service.
// POST /api/v1/auth/chat-token
func (h *AuthHandler) GetChatToken(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok || user.ID == "" {
		Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req chatTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
		Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	perms := req.Permissions
	if len(perms) == 0 {
		perms = []string{"chat:read", "chat:write"}
	}

	token, err := h.signer.Sign(user.ID, user.Email, perms)
	if err != nil {
		h.logger.Error("failed to sign chat token", slog.Any("error", err))
		Error(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	JSON(w, http.StatusOK, chatTokenResponse{Token: token})
}
