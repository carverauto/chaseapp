package search

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"chaseapp.tv/api/internal/config"
	"chaseapp.tv/api/internal/model"
)

type stubTransport func(*http.Request) (*http.Response, error)

func (s stubTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return s(req)
}

func TestNewClientRequiresAPIKey(t *testing.T) {
	_, err := NewClient(config.SearchConfig{})
	if err == nil {
		t.Fatalf("expected error when API key missing")
	}
}

func TestUpsertChaseSendsRequest(t *testing.T) {
	t.Helper()

	var captured struct {
		path    string
		headers http.Header
		body    map[string]any
	}

	cfg := testSearchConfig(t, "http://typesense.local:8108", "secret-key")
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}
	client.baseURL = cfg.URL()
	client.http = &http.Client{Transport: stubTransport(func(r *http.Request) (*http.Response, error) {
		captured.path = r.URL.Path + "?" + r.URL.RawQuery
		captured.headers = r.Header.Clone()
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&captured.body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{}`)),
			Header:     make(http.Header),
		}, nil
	})}

	now := time.Now().UTC()
	chase := &model.Chase{
		ID:        uuid.New(),
		Title:     "Test Chase",
		ChaseType: model.ChaseTypeChase,
		Live:      true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := client.UpsertChase(context.Background(), chase); err != nil {
		t.Fatalf("upsert failed: %v", err)
	}

	if captured.path != "/collections/chases/documents?action=upsert" {
		t.Fatalf("unexpected request path: %s", captured.path)
	}
	if got := captured.headers.Get("X-TYPESENSE-API-KEY"); got != "secret-key" {
		t.Fatalf("expected API key header to be set, got %q", got)
	}
	if captured.body["id"] != chase.ID.String() || captured.body["title"] != chase.Title {
		t.Fatalf("request body missing chase fields: %+v", captured.body)
	}
}

func TestDoReturnsErrorOnNonSuccess(t *testing.T) {
	cfg := testSearchConfig(t, "http://typesense.local:8108", "secret-key")
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}
	client.baseURL = cfg.URL()
	client.http = &http.Client{Transport: stubTransport(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("boom")),
			Header:     make(http.Header),
		}, nil
	})}

	err = client.DeleteChase(context.Background(), uuid.New())
	if err == nil || !strings.Contains(err.Error(), "status 500") {
		t.Fatalf("expected error for non-2xx response, got: %v", err)
	}
}

func testSearchConfig(t *testing.T, rawURL, apiKey string) config.SearchConfig {
	t.Helper()
	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("failed to parse url: %v", err)
	}
	host, portStr, err := net.SplitHostPort(parsed.Host)
	if err != nil {
		t.Fatalf("failed to split host/port: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("invalid port: %v", err)
	}
	return config.SearchConfig{
		Host:     host,
		Port:     port,
		Protocol: parsed.Scheme,
		APIKey:   apiKey,
	}
}
