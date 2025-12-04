package push

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"chaseapp.tv/api/internal/config"
)

const fcmEndpoint = "https://fcm.googleapis.com/v1/projects/%s/messages:send"

// FCMClient sends push notifications via FCM HTTP v1 API.
type FCMClient struct {
	client *http.Client
	cfg    config.PushConfig
}

// NewFCMClient creates a new FCM client.
func NewFCMClient(cfg config.PushConfig) (*FCMClient, error) {
	if cfg.FCMProjectID == "" || cfg.FCMKeyPath == "" {
		return nil, fmt.Errorf("fcm configuration missing")
	}

	// TODO: load service account and mint OAuth2 token (placeholder for now).

	return &FCMClient{
		client: &http.Client{Timeout: 10 * time.Second},
		cfg:    cfg,
	}, nil
}

// FCMMessage represents a simplified FCM message.
type FCMMessage struct {
	Token string
	Title string
	Body  string
	Topic string
	Data  map[string]string
}

// Send pushes a notification to a token or topic.
func (c *FCMClient) Send(ctx context.Context, msg FCMMessage) error {
	target := map[string]any{}
	if msg.Token != "" {
		target["token"] = msg.Token
	} else if msg.Topic != "" {
		target["topic"] = msg.Topic
	} else {
		return fmt.Errorf("fcm requires token or topic")
	}

	notification := map[string]string{
		"title": msg.Title,
		"body":  msg.Body,
	}

	payload := map[string]any{
		"message": map[string]any{
			"notification": notification,
			"data":         msg.Data,
		},
	}
	for k, v := range target {
		payload["message"].(map[string]any)[k] = v
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal fcm payload: %w", err)
	}

	// Placeholder bearer token. In production, mint from service account JWT.
	bearer := "REPLACE_WITH_OAUTH_TOKEN"

	url := fmt.Sprintf(fcmEndpoint, c.cfg.FCMProjectID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+bearer)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("fcm request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("fcm returned status %d", resp.StatusCode)
	}
	return nil
}
