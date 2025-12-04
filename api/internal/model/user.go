package model

import (
	"time"

	"github.com/google/uuid"
)

// AuthProvider represents the OAuth provider.
type AuthProvider string

const (
	AuthProviderGoogle   AuthProvider = "google"
	AuthProviderApple    AuthProvider = "apple"
	AuthProviderFacebook AuthProvider = "facebook"
	AuthProviderTwitter  AuthProvider = "twitter"
)

// User represents a user account.
type User struct {
	ID                   uuid.UUID    `json:"id"`
	ExternalID           string       `json:"external_id"` // OAuth provider user ID
	Email                string       `json:"email,omitempty"`
	DisplayName          string       `json:"display_name,omitempty"`
	PhotoURL             string       `json:"photo_url,omitempty"`
	Provider             AuthProvider `json:"provider"`
	NotificationsEnabled bool         `json:"notifications_enabled"`
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
	LastLoginAt          *time.Time   `json:"last_login_at,omitempty"`
	DeletedAt            *time.Time   `json:"-"`
}

// CreateUserInput represents the input for creating a user.
type CreateUserInput struct {
	ExternalID  string       `json:"external_id" validate:"required"`
	Email       string       `json:"email,omitempty" validate:"omitempty,email"`
	DisplayName string       `json:"display_name,omitempty"`
	PhotoURL    string       `json:"photo_url,omitempty"`
	Provider    AuthProvider `json:"provider" validate:"required,oneof=google apple facebook twitter"`
}

// UpdateUserInput represents the input for updating a user.
type UpdateUserInput struct {
	Email                *string `json:"email,omitempty" validate:"omitempty,email"`
	DisplayName          *string `json:"display_name,omitempty"`
	PhotoURL             *string `json:"photo_url,omitempty"`
	NotificationsEnabled *bool   `json:"notifications_enabled,omitempty"`
}
