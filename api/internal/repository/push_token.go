package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"chaseapp.tv/api/internal/model"
)

// PushTokenRepository handles push token data access.
type PushTokenRepository struct {
	pool *pgxpool.Pool
}

// NewPushTokenRepository creates a new PushTokenRepository.
func NewPushTokenRepository(pool *pgxpool.Pool) *PushTokenRepository {
	return &PushTokenRepository{pool: pool}
}

// Create creates a new push token.
func (r *PushTokenRepository) Create(ctx context.Context, userID *uuid.UUID, input model.CreatePushTokenInput) (*model.PushToken, error) {
	id := uuid.New()
	now := time.Now()

	query := `
		INSERT INTO push_tokens (id, user_id, token, platform, device_id, device_name, app_version, subscribed_topics, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (token, platform) DO UPDATE SET
			user_id = COALESCE(EXCLUDED.user_id, push_tokens.user_id),
			device_id = COALESCE(EXCLUDED.device_id, push_tokens.device_id),
			device_name = COALESCE(EXCLUDED.device_name, push_tokens.device_name),
			app_version = COALESCE(EXCLUDED.app_version, push_tokens.app_version),
			subscribed_topics = COALESCE(EXCLUDED.subscribed_topics, push_tokens.subscribed_topics),
			is_active = true,
			last_used_at = NOW(),
			updated_at = NOW()
		RETURNING id, user_id, token, platform, device_id, device_name, app_version,
			subscribed_topics, is_active, last_used_at, metadata, created_at, updated_at`

	var token model.PushToken
	var metadataBytes []byte

	topics := input.Topics
	if topics == nil {
		topics = []string{}
	}

	err := r.pool.QueryRow(ctx, query,
		id, userID, input.Token, input.Platform, input.DeviceID,
		input.DeviceName, input.AppVersion, topics, now, now,
	).Scan(
		&token.ID, &token.UserID, &token.Token, &token.Platform,
		&token.DeviceID, &token.DeviceName, &token.AppVersion,
		&token.SubscribedTopics, &token.IsActive, &token.LastUsedAt,
		&metadataBytes, &token.CreatedAt, &token.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create push token: %w", err)
	}

	if len(metadataBytes) > 0 {
		json.Unmarshal(metadataBytes, &token.Metadata)
	}

	return &token, nil
}

// GetByID retrieves a push token by ID.
func (r *PushTokenRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.PushToken, error) {
	query := `
		SELECT id, user_id, token, platform, device_id, device_name, app_version,
			   subscribed_topics, is_active, last_used_at, metadata, created_at, updated_at
		FROM push_tokens
		WHERE id = $1`

	return r.scanToken(r.pool.QueryRow(ctx, query, id))
}

// GetByToken retrieves a push token by the token string.
func (r *PushTokenRepository) GetByToken(ctx context.Context, tokenStr string, platform model.Platform) (*model.PushToken, error) {
	query := `
		SELECT id, user_id, token, platform, device_id, device_name, app_version,
			   subscribed_topics, is_active, last_used_at, metadata, created_at, updated_at
		FROM push_tokens
		WHERE token = $1 AND platform = $2`

	return r.scanToken(r.pool.QueryRow(ctx, query, tokenStr, platform))
}

// GetByUserID retrieves all push tokens for a user.
func (r *PushTokenRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.PushToken, error) {
	query := `
		SELECT id, user_id, token, platform, device_id, device_name, app_version,
			   subscribed_topics, is_active, last_used_at, metadata, created_at, updated_at
		FROM push_tokens
		WHERE user_id = $1 AND is_active = true
		ORDER BY last_used_at DESC NULLS LAST`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tokens: %w", err)
	}
	defer rows.Close()

	return r.scanTokens(rows)
}

// GetByTopic retrieves all active tokens subscribed to a topic.
func (r *PushTokenRepository) GetByTopic(ctx context.Context, topic string) ([]model.PushToken, error) {
	query := `
		SELECT id, user_id, token, platform, device_id, device_name, app_version,
			   subscribed_topics, is_active, last_used_at, metadata, created_at, updated_at
		FROM push_tokens
		WHERE $1 = ANY(subscribed_topics) AND is_active = true`

	rows, err := r.pool.Query(ctx, query, topic)
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens by topic: %w", err)
	}
	defer rows.Close()

	return r.scanTokens(rows)
}

// GetByPlatformAndTopic retrieves tokens for a specific platform and topic.
func (r *PushTokenRepository) GetByPlatformAndTopic(ctx context.Context, platform model.Platform, topic string) ([]model.PushToken, error) {
	query := `
		SELECT id, user_id, token, platform, device_id, device_name, app_version,
			   subscribed_topics, is_active, last_used_at, metadata, created_at, updated_at
		FROM push_tokens
		WHERE platform = $1 AND $2 = ANY(subscribed_topics) AND is_active = true`

	rows, err := r.pool.Query(ctx, query, platform, topic)
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens by platform and topic: %w", err)
	}
	defer rows.Close()

	return r.scanTokens(rows)
}

// Update updates a push token.
func (r *PushTokenRepository) Update(ctx context.Context, id uuid.UUID, input model.UpdatePushTokenInput) (*model.PushToken, error) {
	token, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.DeviceName != nil {
		token.DeviceName = *input.DeviceName
	}
	if input.AppVersion != nil {
		token.AppVersion = *input.AppVersion
	}
	if input.Topics != nil {
		token.SubscribedTopics = input.Topics
	}
	if input.IsActive != nil {
		token.IsActive = *input.IsActive
	}

	query := `
		UPDATE push_tokens SET
			device_name = $2, app_version = $3, subscribed_topics = $4,
			is_active = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	err = r.pool.QueryRow(ctx, query,
		id, token.DeviceName, token.AppVersion, token.SubscribedTopics, token.IsActive,
	).Scan(&token.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update push token: %w", err)
	}

	return token, nil
}

// Subscribe adds topics to a token's subscriptions.
func (r *PushTokenRepository) Subscribe(ctx context.Context, id uuid.UUID, topics []string) error {
	query := `
		UPDATE push_tokens SET
			subscribed_topics = array_cat(subscribed_topics, $2::text[]),
			updated_at = NOW()
		WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, topics)
	return err
}

// Unsubscribe removes topics from a token's subscriptions.
func (r *PushTokenRepository) Unsubscribe(ctx context.Context, id uuid.UUID, topics []string) error {
	query := `
		UPDATE push_tokens SET
			subscribed_topics = array(
				SELECT unnest(subscribed_topics) EXCEPT SELECT unnest($2::text[])
			),
			updated_at = NOW()
		WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, topics)
	return err
}

// Deactivate marks a token as inactive (e.g., after push failure).
func (r *PushTokenRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE push_tokens SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// DeactivateByToken marks a token as inactive by token string.
func (r *PushTokenRepository) DeactivateByToken(ctx context.Context, tokenStr string, platform model.Platform) error {
	query := `UPDATE push_tokens SET is_active = false, updated_at = NOW() WHERE token = $1 AND platform = $2`
	_, err := r.pool.Exec(ctx, query, tokenStr, platform)
	return err
}

// Delete removes a push token.
func (r *PushTokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM push_tokens WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete push token: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteByToken removes a push token by token string.
func (r *PushTokenRepository) DeleteByToken(ctx context.Context, tokenStr string, platform model.Platform) error {
	query := `DELETE FROM push_tokens WHERE token = $1 AND platform = $2`
	_, err := r.pool.Exec(ctx, query, tokenStr, platform)
	return err
}

// UpdateLastUsed updates the last_used_at timestamp.
func (r *PushTokenRepository) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE push_tokens SET last_used_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// CountByTopic returns the count of active tokens subscribed to a topic.
func (r *PushTokenRepository) CountByTopic(ctx context.Context, topic string) (int, error) {
	query := `SELECT COUNT(*) FROM push_tokens WHERE $1 = ANY(subscribed_topics) AND is_active = true`
	var count int
	err := r.pool.QueryRow(ctx, query, topic).Scan(&count)
	return count, err
}

// Helper to scan a single token from a row.
func (r *PushTokenRepository) scanToken(row pgx.Row) (*model.PushToken, error) {
	var token model.PushToken
	var metadataBytes []byte

	err := row.Scan(
		&token.ID, &token.UserID, &token.Token, &token.Platform,
		&token.DeviceID, &token.DeviceName, &token.AppVersion,
		&token.SubscribedTopics, &token.IsActive, &token.LastUsedAt,
		&metadataBytes, &token.CreatedAt, &token.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan push token: %w", err)
	}

	if len(metadataBytes) > 0 {
		json.Unmarshal(metadataBytes, &token.Metadata)
	}

	return &token, nil
}

// Helper to scan multiple tokens.
func (r *PushTokenRepository) scanTokens(rows pgx.Rows) ([]model.PushToken, error) {
	var tokens []model.PushToken
	for rows.Next() {
		var token model.PushToken
		var metadataBytes []byte

		err := rows.Scan(
			&token.ID, &token.UserID, &token.Token, &token.Platform,
			&token.DeviceID, &token.DeviceName, &token.AppVersion,
			&token.SubscribedTopics, &token.IsActive, &token.LastUsedAt,
			&metadataBytes, &token.CreatedAt, &token.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan push token: %w", err)
		}

		if len(metadataBytes) > 0 {
			json.Unmarshal(metadataBytes, &token.Metadata)
		}

		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}
