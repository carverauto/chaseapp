package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"chaseapp.tv/api/internal/middleware"
	"chaseapp.tv/api/internal/model"
	"chaseapp.tv/api/internal/repository"
)

// PushHandler handles push notification related HTTP requests.
type PushHandler struct {
	tokenRepo *repository.PushTokenRepository
	userRepo  *repository.UserRepository
	logger    *slog.Logger
}

// NewPushHandler creates a new PushHandler.
func NewPushHandler(tokenRepo *repository.PushTokenRepository, userRepo *repository.UserRepository, logger *slog.Logger) *PushHandler {
	return &PushHandler{
		tokenRepo: tokenRepo,
		userRepo:  userRepo,
		logger:    logger,
	}
}

// SubscribeRequest represents a push subscription request.
type SubscribeRequest struct {
	Token      string   `json:"token"`
	Platform   string   `json:"platform"`
	DeviceID   string   `json:"device_id,omitempty"`
	DeviceName string   `json:"device_name,omitempty"`
	AppVersion string   `json:"app_version,omitempty"`
	Topics     []string `json:"topics,omitempty"`
}

// Subscribe registers a device for push notifications.
// POST /api/v1/push/subscribe
func (h *PushHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Token == "" {
		Error(w, http.StatusBadRequest, "Token is required")
		return
	}

	platform := model.Platform(req.Platform)
	if platform != model.PlatformIOS && platform != model.PlatformAndroid &&
		platform != model.PlatformWeb && platform != model.PlatformSafari {
		Error(w, http.StatusBadRequest, "Invalid platform")
		return
	}

	// Get user ID from context if authenticated
	var userID *uuid.UUID
	if userIDStr, ok := ctx.Value(middleware.UserIDKey).(string); ok && userIDStr != "" {
		if uid, err := uuid.Parse(userIDStr); err == nil {
			userID = &uid
		}
	}

	input := model.CreatePushTokenInput{
		Token:      req.Token,
		Platform:   platform,
		DeviceID:   req.DeviceID,
		DeviceName: req.DeviceName,
		AppVersion: req.AppVersion,
		Topics:     req.Topics,
	}

	token, err := h.tokenRepo.Create(ctx, userID, input)
	if err != nil {
		h.logger.Error("failed to subscribe push token",
			slog.Any("error", err),
			slog.String("platform", string(platform)),
		)
		Error(w, http.StatusInternalServerError, "Failed to register device")
		return
	}

	h.logger.Info("push token registered",
		slog.String("id", token.ID.String()),
		slog.String("platform", string(token.Platform)),
	)

	JSON(w, http.StatusCreated, map[string]any{
		"id":       token.ID,
		"platform": token.Platform,
		"topics":   token.SubscribedTopics,
	})
}

// UnsubscribeRequest represents a push unsubscription request.
type UnsubscribeRequest struct {
	Token    string   `json:"token"`
	Platform string   `json:"platform"`
	Topics   []string `json:"topics,omitempty"` // If empty, fully unsubscribe
}

// Unsubscribe removes a device from push notifications.
// POST /api/v1/push/unsubscribe
func (h *PushHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req UnsubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Token == "" {
		Error(w, http.StatusBadRequest, "Token is required")
		return
	}

	platform := model.Platform(req.Platform)

	// If topics provided, just unsubscribe from those topics
	if len(req.Topics) > 0 {
		token, err := h.tokenRepo.GetByToken(ctx, req.Token, platform)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				Error(w, http.StatusNotFound, "Token not found")
				return
			}
			h.logger.Error("failed to get token", slog.Any("error", err))
			Error(w, http.StatusInternalServerError, "Failed to unsubscribe")
			return
		}

		if err := h.tokenRepo.Unsubscribe(ctx, token.ID, req.Topics); err != nil {
			h.logger.Error("failed to unsubscribe from topics", slog.Any("error", err))
			Error(w, http.StatusInternalServerError, "Failed to unsubscribe")
			return
		}

		h.logger.Info("unsubscribed from topics",
			slog.String("token_id", token.ID.String()),
			slog.Any("topics", req.Topics),
		)

		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Full unsubscription - deactivate or delete token
	if err := h.tokenRepo.DeactivateByToken(ctx, req.Token, platform); err != nil {
		h.logger.Error("failed to deactivate token", slog.Any("error", err))
		Error(w, http.StatusInternalServerError, "Failed to unsubscribe")
		return
	}

	h.logger.Info("push token deactivated",
		slog.String("platform", string(platform)),
	)

	w.WriteHeader(http.StatusNoContent)
}

// GetSafariPushPackage generates a Safari push notification package.
// GET /api/v1/push/safari-package
// This is kept as a standalone function for now as it requires special handling.
func GetSafariPushPackage(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Safari push package generation
	// This requires P12 certificate signing and ZIP file creation
	// Will be implemented in Phase 10
	Error(w, http.StatusNotImplemented, "Safari push package not yet implemented")
}
