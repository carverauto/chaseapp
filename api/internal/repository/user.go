package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"chaseapp.tv/api/internal/model"
)

// UserRepository handles user data access.
type UserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// Create creates a new user.
func (r *UserRepository) Create(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	id := uuid.New()
	now := time.Now()

	query := `
		INSERT INTO users (id, external_id, email, display_name, photo_url, provider, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	var user model.User
	err := r.pool.QueryRow(ctx, query,
		id, input.ExternalID, input.Email, input.DisplayName, input.PhotoURL, input.Provider, now, now,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.ExternalID = input.ExternalID
	user.Email = input.Email
	user.DisplayName = input.DisplayName
	user.PhotoURL = input.PhotoURL
	user.Provider = input.Provider
	user.NotificationsEnabled = true

	return &user, nil
}

// GetByID retrieves a user by ID.
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
		SELECT id, external_id, email, display_name, photo_url, provider,
			   notifications_enabled, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL`

	var user model.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.ExternalID, &user.Email, &user.DisplayName, &user.PhotoURL,
		&user.Provider, &user.NotificationsEnabled, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByExternalID retrieves a user by external OAuth ID.
func (r *UserRepository) GetByExternalID(ctx context.Context, externalID string) (*model.User, error) {
	query := `
		SELECT id, external_id, email, display_name, photo_url, provider,
			   notifications_enabled, created_at, updated_at, last_login_at
		FROM users
		WHERE external_id = $1 AND deleted_at IS NULL`

	var user model.User
	err := r.pool.QueryRow(ctx, query, externalID).Scan(
		&user.ID, &user.ExternalID, &user.Email, &user.DisplayName, &user.PhotoURL,
		&user.Provider, &user.NotificationsEnabled, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by external ID: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, external_id, email, display_name, photo_url, provider,
			   notifications_enabled, created_at, updated_at, last_login_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL`

	var user model.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.ExternalID, &user.Email, &user.DisplayName, &user.PhotoURL,
		&user.Provider, &user.NotificationsEnabled, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// Update updates a user.
func (r *UserRepository) Update(ctx context.Context, id uuid.UUID, input model.UpdateUserInput) (*model.User, error) {
	// Get current user
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.DisplayName != nil {
		user.DisplayName = *input.DisplayName
	}
	if input.PhotoURL != nil {
		user.PhotoURL = *input.PhotoURL
	}
	if input.NotificationsEnabled != nil {
		user.NotificationsEnabled = *input.NotificationsEnabled
	}

	query := `
		UPDATE users SET
			email = $2, display_name = $3, photo_url = $4, notifications_enabled = $5,
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`

	err = r.pool.QueryRow(ctx, query,
		id, user.Email, user.DisplayName, user.PhotoURL, user.NotificationsEnabled,
	).Scan(&user.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// Delete soft-deletes a user.
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateLastLogin updates the last login timestamp.
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET last_login_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// FindOrCreate finds a user by external ID or creates a new one.
func (r *UserRepository) FindOrCreate(ctx context.Context, input model.CreateUserInput) (*model.User, bool, error) {
	// Try to find existing user
	user, err := r.GetByExternalID(ctx, input.ExternalID)
	if err == nil {
		// Update last login
		r.UpdateLastLogin(ctx, user.ID)
		return user, false, nil
	}

	if !errors.Is(err, ErrNotFound) {
		return nil, false, err
	}

	// Create new user
	user, err = r.Create(ctx, input)
	if err != nil {
		return nil, false, err
	}

	return user, true, nil
}
