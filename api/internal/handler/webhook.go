package handler

import (
	"net/http"
)

// SendDiscordWebhook sends a formatted message to a Discord webhook.
// POST /api/v1/webhooks/discord
func SendDiscordWebhook(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Discord webhook integration
	Error(w, http.StatusNotImplemented, "Not implemented")
}
