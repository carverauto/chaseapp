package model

import (
	"time"

	"github.com/google/uuid"
)

// AircraftCategory represents the type/category of aircraft.
type AircraftCategory string

const (
	AircraftCategoryMedia          AircraftCategory = "media"
	AircraftCategoryLawEnforcement AircraftCategory = "law_enforcement"
	AircraftCategoryMilitary       AircraftCategory = "military"
	AircraftCategoryMedical        AircraftCategory = "medical"
	AircraftCategoryFirefighting   AircraftCategory = "firefighting"
	AircraftCategoryGeneral        AircraftCategory = "general"
)

// Aircraft represents an ADSB-tracked aircraft.
type Aircraft struct {
	ID uuid.UUID `json:"id"`

	// Aircraft identification
	ICAO         string `json:"icao"`                    // ICAO 24-bit address (hex)
	Callsign     string `json:"callsign,omitempty"`      // Flight callsign
	Registration string `json:"registration,omitempty"` // Tail number

	// Position
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
	Altitude     *int     `json:"altitude,omitempty"`      // Feet
	GroundSpeed  *int     `json:"ground_speed,omitempty"`  // Knots
	Track        *int     `json:"track,omitempty"`         // Heading 0-359
	VerticalRate *int     `json:"vertical_rate,omitempty"` // ft/min

	// Aircraft info
	AircraftType string           `json:"aircraft_type,omitempty"` // ICAO type code
	Category     AircraftCategory `json:"category,omitempty"`
	Operator     string           `json:"operator,omitempty"`

	// Status
	OnGround  bool   `json:"on_ground"`
	Squawk    string `json:"squawk,omitempty"`    // Transponder code
	Emergency string `json:"emergency,omitempty"` // Emergency status

	// Tracking
	LastSeenAt  time.Time `json:"last_seen_at"`
	FirstSeenAt time.Time `json:"first_seen_at"`

	// Clustering
	ClusterID string `json:"cluster_id,omitempty"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AircraftHistory represents a historical position record.
type AircraftHistory struct {
	ID         uuid.UUID `json:"id"`
	AircraftID uuid.UUID `json:"aircraft_id"`

	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Altitude    *int    `json:"altitude,omitempty"`
	GroundSpeed *int    `json:"ground_speed,omitempty"`
	Track       *int    `json:"track,omitempty"`

	RecordedAt time.Time `json:"recorded_at"`
}

// UpsertAircraftInput represents the input for creating or updating an aircraft.
type UpsertAircraftInput struct {
	ICAO         string  `json:"icao" validate:"required"`
	Callsign     string  `json:"callsign,omitempty"`
	Registration string  `json:"registration,omitempty"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
	Altitude     *int     `json:"altitude,omitempty"`
	GroundSpeed  *int     `json:"ground_speed,omitempty"`
	Track        *int     `json:"track,omitempty"`
	VerticalRate *int     `json:"vertical_rate,omitempty"`
	AircraftType string   `json:"aircraft_type,omitempty"`
	Category     AircraftCategory `json:"category,omitempty"`
	Operator     string   `json:"operator,omitempty"`
	OnGround     bool     `json:"on_ground"`
	Squawk       string   `json:"squawk,omitempty"`
	Emergency    string   `json:"emergency,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// AircraftListOptions represents options for listing aircraft.
type AircraftListOptions struct {
	Page      int              `json:"page"`
	Limit     int              `json:"limit"`
	Category  AircraftCategory `json:"category,omitempty"`
	ClusterID string           `json:"cluster_id,omitempty"`
	OnGround  *bool            `json:"on_ground,omitempty"`
	// Bounding box for geographic filtering
	MinLat *float64 `json:"min_lat,omitempty"`
	MaxLat *float64 `json:"max_lat,omitempty"`
	MinLng *float64 `json:"min_lng,omitempty"`
	MaxLng *float64 `json:"max_lng,omitempty"`
}

// AircraftListResult represents a paginated list of aircraft.
type AircraftListResult struct {
	Aircraft   []Aircraft `json:"aircraft"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalPages int        `json:"total_pages"`
}
