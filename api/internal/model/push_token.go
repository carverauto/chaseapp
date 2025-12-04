package model

import (
	"time"

	"github.com/google/uuid"
)

// Platform represents the push notification platform.
type Platform string

const (
	PlatformIOS     Platform = "ios"
	PlatformAndroid Platform = "android"
	PlatformWeb     Platform = "web"
	PlatformSafari  Platform = "safari"
)

// PushToken represents a device's push notification token.
type PushToken struct {
	ID     uuid.UUID  `json:"id"`
	UserID *uuid.UUID `json:"user_id,omitempty"`

	// Token info
	Token    string   `json:"token"`
	Platform Platform `json:"platform"`

	// Device info
	DeviceID   string `json:"device_id,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
	AppVersion string `json:"app_version,omitempty"`

	// Subscription preferences
	SubscribedTopics []string `json:"subscribed_topics"`

	// Status
	IsActive   bool       `json:"is_active"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`

	// Metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreatePushTokenInput represents the input for registering a push token.
type CreatePushTokenInput struct {
	Token      string   `json:"token" validate:"required"`
	Platform   Platform `json:"platform" validate:"required,oneof=ios android web safari"`
	DeviceID   string   `json:"device_id,omitempty"`
	DeviceName string   `json:"device_name,omitempty"`
	AppVersion string   `json:"app_version,omitempty"`
	Topics     []string `json:"topics,omitempty"`
}

// UpdatePushTokenInput represents the input for updating a push token.
type UpdatePushTokenInput struct {
	DeviceName *string  `json:"device_name,omitempty"`
	AppVersion *string  `json:"app_version,omitempty"`
	Topics     []string `json:"topics,omitempty"`
	IsActive   *bool    `json:"is_active,omitempty"`
}
