package model

// ClusterAircraftInput represents a clustering request payload.
type ClusterAircraftInput struct {
	Points        []ClusterPoint `json:"points"`
	EpsilonMeters float64        `json:"eps_meters"` // Neighborhood radius in meters
	MinPoints     int            `json:"min_points"` // Minimum points to form a cluster
}

// ClusterPoint represents an aircraft point to be clustered.
type ClusterPoint struct {
	ID        string           `json:"id,omitempty"`
	ICAO      string           `json:"icao,omitempty"`
	Callsign  string           `json:"callsign,omitempty"`
	Latitude  float64          `json:"latitude"`
	Longitude float64          `json:"longitude"`
	Altitude  *int             `json:"altitude,omitempty"`
	Category  AircraftCategory `json:"category,omitempty"`
	OnGround  bool             `json:"on_ground"`
	Metadata  map[string]any   `json:"metadata,omitempty"`
}

// ClusterResult represents a single cluster output.
type ClusterResult struct {
	ID           string         `json:"id"`
	Points       []ClusterPoint `json:"points"`
	CentroidLat  float64        `json:"centroid_lat"`
	CentroidLng  float64        `json:"centroid_lng"`
	MediaPresent bool           `json:"media_present"`
	Size         int            `json:"size"`
}

// ClusterResponse represents the clustering response payload.
type ClusterResponse struct {
	Clusters []ClusterResult `json:"clusters"`
}
