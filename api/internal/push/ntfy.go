package push

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"chaseapp.tv/api/internal/config"
)

// NtfyClient sends notifications to an ntfy/Gotify topic endpoint.
type NtfyClient struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewNtfyClient creates a new ntfy client.
func NewNtfyClient(cfg config.PushConfig) (*NtfyClient, error) {
	if cfg.NtfyURL == "" {
		return nil, errors.New("ntfy url not configured")
	}
	base := strings.TrimSuffix(cfg.NtfyURL, "/")
	return &NtfyClient{
		baseURL: base,
		token:   cfg.NtfyToken,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// Publish sends a simple notification to the given topic.
func (c *NtfyClient) Publish(ctx context.Context, topic, title, message string) error {
	if topic == "" {
		return errors.New("topic is required")
	}
	url := fmt.Sprintf("%s/%s", c.baseURL, topic)

	body := &bytes.Buffer{}
	body.WriteString(message)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Title", title)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("send ntfy request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy returned status %d", resp.StatusCode)
	}
	return nil
}
