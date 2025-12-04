// Package model defines domain models for the API.
package model

import (
	"time"

	"github.com/google/uuid"
)

// ChaseType represents the type of chase/event.
type ChaseType string

const (
	ChaseTypeChase    ChaseType = "chase"
	ChaseTypeRocket   ChaseType = "rocket"
	ChaseTypeWeather  ChaseType = "weather"
	ChaseTypeAircraft ChaseType = "aircraft"
)

// Location represents a geographic location.
type Location struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Address string  `json:"address,omitempty"`
}

// Stream represents a live stream source.
type Stream struct {
	URL     string `json:"url"`
	Network string `json:"network,omitempty"`
	Type    string `json:"type,omitempty"` // m3u8, mp4, etc.
}

// Chase represents a live event (police chase, rocket launch, etc.).
type Chase struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	ChaseType   ChaseType  `json:"chase_type"`
	Location    *Location  `json:"location,omitempty"`
	City        string     `json:"city,omitempty"`
	State       string     `json:"state,omitempty"`
	Country     string     `json:"country,omitempty"`
	Live        bool       `json:"live"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	EndedAt     *time.Time `json:"ended_at,omitempty"`

	ThumbnailURL string   `json:"thumbnail_url,omitempty"`
	Streams      []Stream `json:"streams,omitempty"`

	ViewCount  int `json:"view_count"`
	ShareCount int `json:"share_count"`

	Source    string `json:"source,omitempty"`
	SourceURL string `json:"source_url,omitempty"`

	CreatedBy *uuid.UUID             `json:"created_by,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

// CreateChaseInput represents the input for creating a chase.
type CreateChaseInput struct {
	Title        string                 `json:"title" validate:"required,min=1,max=500"`
	Description  string                 `json:"description,omitempty"`
	ChaseType    ChaseType              `json:"chase_type" validate:"required,oneof=chase rocket weather aircraft"`
	Location     *Location              `json:"location,omitempty"`
	City         string                 `json:"city,omitempty"`
	State        string                 `json:"state,omitempty"`
	Country      string                 `json:"country,omitempty"`
	Live         bool                   `json:"live"`
	ThumbnailURL string                 `json:"thumbnail_url,omitempty"`
	Streams      []Stream               `json:"streams,omitempty"`
	Source       string                 `json:"source,omitempty"`
	SourceURL    string                 `json:"source_url,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateChaseInput represents the input for updating a chase.
type UpdateChaseInput struct {
	Title        *string                `json:"title,omitempty"`
	Description  *string                `json:"description,omitempty"`
	Location     *Location              `json:"location,omitempty"`
	City         *string                `json:"city,omitempty"`
	State        *string                `json:"state,omitempty"`
	Live         *bool                  `json:"live,omitempty"`
	ThumbnailURL *string                `json:"thumbnail_url,omitempty"`
	Streams      []Stream               `json:"streams,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ChaseListOptions represents options for listing chases.
type ChaseListOptions struct {
	Page      int       `json:"page"`
	Limit     int       `json:"limit"`
	Live      *bool     `json:"live,omitempty"`
	ChaseType ChaseType `json:"chase_type,omitempty"`
	City      string    `json:"city,omitempty"`
	State     string    `json:"state,omitempty"`
}

// ChaseListResult represents a paginated list of chases.
type ChaseListResult struct {
	Chases     []Chase `json:"chases"`
	Total      int     `json:"total"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	TotalPages int     `json:"total_pages"`
}
