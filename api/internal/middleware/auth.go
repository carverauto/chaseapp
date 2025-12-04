package middleware

import (
	"context"
	"net/http"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const (
	// UserIDKey is the context key for user ID.
	UserIDKey contextKey = "user_id"
	// UserEmailKey is the context key for user email.
	UserEmailKey contextKey = "user_email"
)

// User represents an authenticated user extracted from request headers.
type User struct {
	ID    string
	Email string
}

// UserFromContext extracts the authenticated user from context.
func UserFromContext(ctx context.Context) (*User, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return nil, false
	}

	email, _ := ctx.Value(UserEmailKey).(string)

	return &User{
		ID:    userID,
		Email: email,
	}, true
}

// Auth extracts user information from Kong-provided headers.
// Kong validates JWTs and passes user info via X-User-ID and X-User-Email headers.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		userEmail := r.Header.Get("X-User-Email")

		// Add user info to context if present
		ctx := r.Context()
		if userID != "" {
			ctx = context.WithValue(ctx, UserIDKey, userID)
		}
		if userEmail != "" {
			ctx = context.WithValue(ctx, UserEmailKey, userEmail)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth ensures the request has a valid user context.
// Use this for endpoints that require authentication.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := UserFromContext(r.Context()); !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
