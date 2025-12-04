package scraper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"chaseapp.tv/api/internal/model"
)

const (
	defaultTimeout   = 15 * time.Second
	defaultUserAgent = "chaseapp-stream-extractor/1.0"
)

var (
	m3u8Regex = regexp.MustCompile(`https?://[^"'\\s]+\\.m3u8`)
	mp4Regex  = regexp.MustCompile(`https?://[^"'\\s]+\\.mp4`)
)

// Extractor fetches pages and extracts media stream URLs.
type Extractor struct {
	client *http.Client
}

// NewExtractor creates a new extractor with sensible defaults.
func NewExtractor() *Extractor {
	return &Extractor{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// Extract downloads the target URL and returns discovered stream URLs.
func (e *Extractor) Extract(ctx context.Context, target string) ([]model.Stream, error) {
	if target == "" {
		return nil, errors.New("target URL is required")
	}

	u, err := url.Parse(target)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid URL: %s", target)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("received status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return extractStreamsFromBody(u.Host, string(body)), nil
}

func extractStreamsFromBody(host, body string) []model.Stream {
	network := inferNetwork(host)
	seen := make(map[string]model.Stream)

	for _, match := range m3u8Regex.FindAllString(body, -1) {
		seen[match] = model.Stream{
			URL:     match,
			Type:    "m3u8",
			Network: network,
		}
	}

	for _, match := range mp4Regex.FindAllString(body, -1) {
		if _, ok := seen[match]; ok {
			continue
		}
		seen[match] = model.Stream{
			URL:     match,
			Type:    "mp4",
			Network: network,
		}
	}

	streams := make([]model.Stream, 0, len(seen))
	for _, s := range seen {
		streams = append(streams, s)
	}
	return streams
}

func inferNetwork(host string) string {
	host = strings.ToLower(host)
	switch {
	case strings.Contains(host, "nbc"):
		return "NBC LA"
	case strings.Contains(host, "abc7") || strings.Contains(host, "abc"):
		return "ABC7"
	case strings.Contains(host, "cbs"):
		return "CBS News"
	default:
		return host
	}
}
