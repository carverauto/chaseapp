package push

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	"chaseapp.tv/api/internal/config"
)

const apnsURL = "https://api.push.apple.com/3/device/"

// APNsClient sends notifications to Apple Push Notification service.
type APNsClient struct {
	httpClient *http.Client
	cfg        config.PushConfig
}

// NewAPNsClient creates a new APNs client using a P8 key converted to PKCS#12.
// NOTE: This is a simplified placeholder that relies on APNsKeyPath being a PKCS#12/P8 client cert.
func NewAPNsClient(cfg config.PushConfig) (*APNsClient, error) {
	if cfg.APNsKeyPath == "" || cfg.APNsKeyID == "" || cfg.APNsTeamID == "" || cfg.APNsBundleID == "" {
		return nil, fmt.Errorf("apns configuration missing")
	}

	// Load client certificate (for sandbox/dev you may use JWT; here we expect a TLS client cert).
	cert, err := tls.LoadX509KeyPair(cfg.APNsKeyPath, cfg.APNsKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load apns key pair: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	return &APNsClient{
		httpClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: &http.Transport{TLSClientConfig: tlsConfig},
		},
		cfg: cfg,
	}, nil
}

// APNsMessage is the simplified payload for APNs.
type APNsMessage struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Topic string `json:"-"`
}

// Send pushes a notification to a specific device token.
func (c *APNsClient) Send(ctx context.Context, deviceToken string, msg APNsMessage) error {
	if deviceToken == "" {
		return fmt.Errorf("device token is required")
	}

	payload := map[string]any{
		"aps": map[string]any{
			"alert": map[string]string{
				"title": msg.Title,
				"body":  msg.Body,
			},
			"sound": "default",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apnsURL+deviceToken, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	topic := msg.Topic
	if topic == "" {
		topic = c.cfg.APNsBundleID
	}

	req.Header.Set("apns-topic", topic)
	req.Header.Set("apns-id", uuid.New().String())
	req.Header.Set("content-type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("apns request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("apns returned status %d", resp.StatusCode)
	}

	return nil
}
