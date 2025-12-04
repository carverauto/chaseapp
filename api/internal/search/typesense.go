package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"chaseapp.tv/api/internal/config"
	"chaseapp.tv/api/internal/model"
)

const chaseCollection = "chases"

// Client provides minimal Typesense operations.
type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

// NewClient creates a Typesense client.
func NewClient(cfg config.SearchConfig) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("typesense API key is required")
	}
	base := strings.TrimSuffix(cfg.URL(), "/")
	return &Client{
		baseURL: base,
		apiKey:  cfg.APIKey,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// EnsureCollection creates the chase collection if it does not exist.
func (c *Client) EnsureCollection(ctx context.Context) error {
	schema := map[string]any{
		"name": chaseCollection,
		"fields": []map[string]any{
			{"name": "id", "type": "string"},
			{"name": "title", "type": "string"},
			{"name": "description", "type": "string", "optional": true},
			{"name": "chase_type", "type": "string", "facet": true},
			{"name": "city", "type": "string", "optional": true, "facet": true},
			{"name": "state", "type": "string", "optional": true, "facet": true},
			{"name": "country", "type": "string", "optional": true, "facet": true},
			{"name": "live", "type": "bool", "facet": true},
			{"name": "started_at", "type": "int64", "optional": true},
			{"name": "ended_at", "type": "int64", "optional": true},
			{"name": "created_at", "type": "int64"},
		},
		"default_sorting_field": "created_at",
	}

	// Try creating; if already exists, ignore.
	if err := c.do(ctx, http.MethodPost, "/collections", schema, nil); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return err
	}
	return nil
}

// UpsertChase indexes a chase document.
func (c *Client) UpsertChase(ctx context.Context, chase *model.Chase) error {
	doc := mapChase(chase)
	path := fmt.Sprintf("/collections/%s/documents?action=upsert", chaseCollection)
	return c.do(ctx, http.MethodPost, path, doc, nil)
}

// DeleteChase removes a chase document.
func (c *Client) DeleteChase(ctx context.Context, id uuid.UUID) error {
	path := fmt.Sprintf("/collections/%s/documents/%s", chaseCollection, id.String())
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

// Search returns matched chases.
func (c *Client) Search(ctx context.Context, query string, page, perPage int) (*SearchResult, error) {
	if perPage <= 0 || perPage > 50 {
		perPage = 20
	}
	if page <= 0 {
		page = 1
	}

	payload := map[string]any{
		"q":        query,
		"query_by": "title,description",
		"page":     page,
		"per_page": perPage,
	}

	var result SearchResult
	path := fmt.Sprintf("/collections/%s/documents/search", chaseCollection)
	if err := c.do(ctx, http.MethodPost, path, payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SearchResult represents a Typesense search response.
type SearchResult struct {
	Found int `json:"found"`
	Page  int `json:"page"`
	Hits  []struct {
		Document map[string]any `json:"document"`
	} `json:"hits"`
}

func (c *Client) do(ctx context.Context, method, path string, payload any, out any) error {
	var body *bytes.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshal payload: %w", err)
		}
		body = bytes.NewReader(data)
	} else {
		body = bytes.NewReader(nil)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-TYPESENSE-API-KEY", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("typesense request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("typesense returned status %d", resp.StatusCode)
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

func mapChase(chase *model.Chase) map[string]any {
	doc := map[string]any{
		"id":          chase.ID.String(),
		"title":       chase.Title,
		"description": chase.Description,
		"chase_type":  chase.ChaseType,
		"city":        chase.City,
		"state":       chase.State,
		"country":     chase.Country,
		"live":        chase.Live,
		"created_at":  chase.CreatedAt.Unix(),
	}
	if chase.StartedAt != nil {
		doc["started_at"] = chase.StartedAt.Unix()
	}
	if chase.EndedAt != nil {
		doc["ended_at"] = chase.EndedAt.Unix()
	}
	return doc
}
