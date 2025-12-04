package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"chaseapp.tv/api/internal/config"
	"chaseapp.tv/api/internal/webhook"
)

// WebhookHandler handles outbound webhook integrations.
type WebhookHandler struct {
	discord *webhook.Client
	logger  *slog.Logger
}

// DiscordClient exposes the Discord client for internal workers.
func (h *WebhookHandler) DiscordClient() *webhook.Client {
	return h.discord
}

// NewWebhookHandler creates a webhook handler.
func NewWebhookHandler(cfg config.ExternalConfig, logger *slog.Logger) (*WebhookHandler, error) {
	if cfg.DiscordWebhook == "" {
		logger.Warn("discord webhook URL not configured")
		return &WebhookHandler{
			discord: nil,
			logger:  logger,
		}, nil
	}

	client, err := webhook.NewClient(cfg.DiscordWebhook)
	if err != nil {
		return nil, err
	}
	return &WebhookHandler{
		discord: client,
		logger:  logger,
	}, nil
}

type discordWebhookRequest struct {
	Content string                 `json:"content,omitempty"`
	Title   string                 `json:"title,omitempty"`
	Body    string                 `json:"body,omitempty"`
	URL     string                 `json:"url,omitempty"`
	Color   int                    `json:"color,omitempty"`
	Fields  []webhook.EmbedField   `json:"fields,omitempty"`
	Extras  map[string]interface{} `json:"extras,omitempty"` // optional extra fields for templating
}

// SendDiscordWebhook sends a formatted message to a Discord webhook.
// POST /api/v1/webhooks/discord
func (h *WebhookHandler) SendDiscordWebhook(w http.ResponseWriter, r *http.Request) {
	if h.discord == nil {
		Error(w, http.StatusServiceUnavailable, "Discord webhook not configured")
		return
	}

	var req discordWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	embed := webhook.Embed{
		Title:       req.Title,
		Description: req.Body,
		URL:         req.URL,
		Color:       req.Color,
		Fields:      req.Fields,
	}

	msg := webhook.Message{
		Content: req.Content,
	}
	// Only include embed if it has content
	if embed.Title != "" || embed.Description != "" || len(embed.Fields) > 0 {
		msg.Embeds = []webhook.Embed{embed}
	}

	if err := h.discord.Send(r.Context(), msg); err != nil {
		h.logger.Error("failed to send discord webhook", slog.Any("error", err))
		Error(w, http.StatusBadGateway, "Failed to send webhook")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
