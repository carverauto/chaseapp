package external

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"chaseapp.tv/api/internal/config"
)

const (
	defaultHTTPTimeout   = 15 * time.Second
	defaultShortCacheTTL = 2 * time.Minute
	defaultLongCacheTTL  = 5 * time.Minute
)

// Client fetches data from external data sources with basic caching.
type Client struct {
	httpClient *http.Client
	cfg        config.ExternalConfig
	cache      *cache
	logger     *slog.Logger
}

// Launches groups recent and upcoming launches.
type Launches struct {
	Upcoming json.RawMessage `json:"upcoming"`
	Recent   json.RawMessage `json:"recent"`
}

// NewClient creates a new external client.
func NewClient(cfg config.ExternalConfig, logger *slog.Logger) *Client {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: defaultHTTPTimeout,
		},
		cache:  newCache(),
		logger: logger,
	}
}

// GetQuakes retrieves earthquake GeoJSON from USGS.
func (c *Client) GetQuakes(ctx context.Context) (json.RawMessage, error) {
	base := strings.TrimSuffix(c.cfg.USGSBaseURL, "/")
	url := fmt.Sprintf("%s/earthquakes/feed/v1.0/summary/all_hour.geojson", base)
	return c.fetchJSON(ctx, "quakes:all_hour", url, defaultShortCacheTTL, nil)
}

// GetBoats retrieves vessel data from AISHub.
func (c *Client) GetBoats(ctx context.Context) (json.RawMessage, error) {
	if c.cfg.AISHubAPIKey == "" {
		return nil, errors.New("AISHub API key is not configured")
	}

	u, err := url.Parse(c.cfg.AISHubBaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid AISHub base URL: %w", err)
	}

	q := u.Query()
	q.Set("username", c.cfg.AISHubAPIKey)
	q.Set("format", "1")
	q.Set("output", "json")
	q.Set("compress", "0")
	u.RawQuery = q.Encode()

	return c.fetchJSON(ctx, "boats:latest", u.String(), defaultShortCacheTTL, nil)
}

// GetLaunches retrieves recent and upcoming rocket launches.
func (c *Client) GetLaunches(ctx context.Context) (*Launches, error) {
	base := strings.TrimSuffix(c.cfg.LaunchLibraryBaseURL, "/")
	upcomingURL := fmt.Sprintf("%s/launch/upcoming/?limit=20&mode=list", base)
	recentURL := fmt.Sprintf("%s/launch/previous/?limit=10&mode=list", base)

	upcoming, err := c.fetchJSON(ctx, "launches:upcoming", upcomingURL, defaultLongCacheTTL, nil)
	if err != nil {
		return nil, err
	}

	recent, err := c.fetchJSON(ctx, "launches:recent", recentURL, defaultLongCacheTTL, nil)
	if err != nil {
		return nil, err
	}

	return &Launches{
		Upcoming: upcoming,
		Recent:   recent,
	}, nil
}

// GetWeatherAlerts retrieves active weather alerts from NOAA/NWS.
func (c *Client) GetWeatherAlerts(ctx context.Context, area string) (json.RawMessage, error) {
	base := strings.TrimSuffix(c.cfg.NOAABaseURL, "/")
	u, err := url.Parse(base + "/alerts/active")
	if err != nil {
		return nil, fmt.Errorf("invalid NOAA base URL: %w", err)
	}

	if area != "" {
		q := u.Query()
		q.Set("area", strings.ToUpper(area))
		u.RawQuery = q.Encode()
	}

	headers := map[string]string{
		"Accept":     "application/geo+json",
		"User-Agent": "chaseapp-api/1.0",
	}

	cacheKey := "weather:active"
	if area != "" {
		cacheKey = cacheKey + ":" + strings.ToUpper(area)
	}

	return c.fetchJSON(ctx, cacheKey, u.String(), defaultShortCacheTTL, headers)
}

// fetchJSON performs an HTTP GET with caching and returns the decoded body as json.RawMessage.
func (c *Client) fetchJSON(ctx context.Context, cacheKey, targetURL string, ttl time.Duration, headers map[string]string) (json.RawMessage, error) {
	if data, ok := c.cache.get(cacheKey); ok {
		return data, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("upstream returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if !json.Valid(body) {
		c.logger.Warn("received invalid JSON from external service",
			slog.String("url", targetURL),
		)
		return nil, errors.New("external service returned invalid JSON")
	}

	raw := json.RawMessage(body)
	c.cache.set(cacheKey, raw, ttl)
	return raw, nil
}
