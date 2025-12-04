// Package repository provides data access layer implementations.
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

// ErrNotFound is returned when a resource is not found.
var ErrNotFound = errors.New("resource not found")

// ChaseRepository handles chase data access.
type ChaseRepository struct {
	pool *pgxpool.Pool
}

// NewChaseRepository creates a new ChaseRepository.
func NewChaseRepository(pool *pgxpool.Pool) *ChaseRepository {
	return &ChaseRepository{pool: pool}
}

// Create creates a new chase.
func (r *ChaseRepository) Create(ctx context.Context, input model.CreateChaseInput, createdBy *uuid.UUID) (*model.Chase, error) {
	id := uuid.New()
	now := time.Now()

	locationJSON, err := json.Marshal(input.Location)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal location: %w", err)
	}

	streamsJSON, err := json.Marshal(input.Streams)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal streams: %w", err)
	}

	metadataJSON, err := json.Marshal(input.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	var startedAt *time.Time
	if input.Live {
		startedAt = &now
	}

	query := `
		INSERT INTO chases (
			id, title, description, chase_type, location, city, state, country,
			live, started_at, thumbnail_url, streams, source, source_url,
			created_by, metadata, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		) RETURNING id, created_at, updated_at`

	var chase model.Chase
	err = r.pool.QueryRow(ctx, query,
		id, input.Title, input.Description, input.ChaseType, locationJSON,
		input.City, input.State, input.Country, input.Live, startedAt,
		input.ThumbnailURL, streamsJSON, input.Source, input.SourceURL,
		createdBy, metadataJSON, now, now,
	).Scan(&chase.ID, &chase.CreatedAt, &chase.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create chase: %w", err)
	}

	// Populate the chase with input data
	chase.Title = input.Title
	chase.Description = input.Description
	chase.ChaseType = input.ChaseType
	chase.Location = input.Location
	chase.City = input.City
	chase.State = input.State
	chase.Country = input.Country
	chase.Live = input.Live
	chase.StartedAt = startedAt
	chase.ThumbnailURL = input.ThumbnailURL
	chase.Streams = input.Streams
	chase.Source = input.Source
	chase.SourceURL = input.SourceURL
	chase.CreatedBy = createdBy
	chase.Metadata = input.Metadata

	return &chase, nil
}

// GetByID retrieves a chase by ID.
func (r *ChaseRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Chase, error) {
	query := `
		SELECT id, title, description, chase_type, location, city, state, country,
			   live, started_at, ended_at, thumbnail_url, streams, view_count, share_count,
			   source, source_url, created_by, metadata, created_at, updated_at
		FROM chases
		WHERE id = $1 AND deleted_at IS NULL`

	var chase model.Chase
	var locationJSON, streamsJSON, metadataJSON []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&chase.ID, &chase.Title, &chase.Description, &chase.ChaseType,
		&locationJSON, &chase.City, &chase.State, &chase.Country,
		&chase.Live, &chase.StartedAt, &chase.EndedAt, &chase.ThumbnailURL,
		&streamsJSON, &chase.ViewCount, &chase.ShareCount,
		&chase.Source, &chase.SourceURL, &chase.CreatedBy, &metadataJSON,
		&chase.CreatedAt, &chase.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get chase: %w", err)
	}

	// Unmarshal JSON fields
	if len(locationJSON) > 0 {
		if err := json.Unmarshal(locationJSON, &chase.Location); err != nil {
			return nil, fmt.Errorf("failed to unmarshal location: %w", err)
		}
	}
	if len(streamsJSON) > 0 {
		if err := json.Unmarshal(streamsJSON, &chase.Streams); err != nil {
			return nil, fmt.Errorf("failed to unmarshal streams: %w", err)
		}
	}
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &chase.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return &chase, nil
}

// List retrieves chases with pagination and filtering.
func (r *ChaseRepository) List(ctx context.Context, opts model.ChaseListOptions) (*model.ChaseListResult, error) {
	// Set defaults
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.Limit < 1 || opts.Limit > 100 {
		opts.Limit = 20
	}

	offset := (opts.Page - 1) * opts.Limit

	// Build query with filters
	baseQuery := `FROM chases WHERE deleted_at IS NULL`
	args := []interface{}{}
	argNum := 1

	if opts.Live != nil {
		baseQuery += fmt.Sprintf(" AND live = $%d", argNum)
		args = append(args, *opts.Live)
		argNum++
	}
	if opts.ChaseType != "" {
		baseQuery += fmt.Sprintf(" AND chase_type = $%d", argNum)
		args = append(args, opts.ChaseType)
		argNum++
	}
	if opts.City != "" {
		baseQuery += fmt.Sprintf(" AND city = $%d", argNum)
		args = append(args, opts.City)
		argNum++
	}
	if opts.State != "" {
		baseQuery += fmt.Sprintf(" AND state = $%d", argNum)
		args = append(args, opts.State)
		argNum++
	}

	// Get total count
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count chases: %w", err)
	}

	// Get chases
	selectQuery := fmt.Sprintf(`
		SELECT id, title, description, chase_type, location, city, state, country,
			   live, started_at, ended_at, thumbnail_url, streams, view_count, share_count,
			   source, source_url, created_by, metadata, created_at, updated_at
		%s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		baseQuery, argNum, argNum+1)

	args = append(args, opts.Limit, offset)

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list chases: %w", err)
	}
	defer rows.Close()

	var chases []model.Chase
	for rows.Next() {
		var chase model.Chase
		var locationJSON, streamsJSON, metadataJSON []byte

		err := rows.Scan(
			&chase.ID, &chase.Title, &chase.Description, &chase.ChaseType,
			&locationJSON, &chase.City, &chase.State, &chase.Country,
			&chase.Live, &chase.StartedAt, &chase.EndedAt, &chase.ThumbnailURL,
			&streamsJSON, &chase.ViewCount, &chase.ShareCount,
			&chase.Source, &chase.SourceURL, &chase.CreatedBy, &metadataJSON,
			&chase.CreatedAt, &chase.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chase: %w", err)
		}

		// Unmarshal JSON fields
		if len(locationJSON) > 0 {
			json.Unmarshal(locationJSON, &chase.Location)
		}
		if len(streamsJSON) > 0 {
			json.Unmarshal(streamsJSON, &chase.Streams)
		}
		if len(metadataJSON) > 0 {
			json.Unmarshal(metadataJSON, &chase.Metadata)
		}

		chases = append(chases, chase)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate chases: %w", err)
	}

	totalPages := (total + opts.Limit - 1) / opts.Limit

	return &model.ChaseListResult{
		Chases:     chases,
		Total:      total,
		Page:       opts.Page,
		Limit:      opts.Limit,
		TotalPages: totalPages,
	}, nil
}

// Update updates a chase.
func (r *ChaseRepository) Update(ctx context.Context, id uuid.UUID, input model.UpdateChaseInput) (*model.Chase, error) {
	// Get current chase first
	chase, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Track if Live status changed from true to false
	wasLive := chase.Live
	var endedAt *time.Time

	// Apply updates
	if input.Title != nil {
		chase.Title = *input.Title
	}
	if input.Description != nil {
		chase.Description = *input.Description
	}
	if input.Location != nil {
		chase.Location = input.Location
	}
	if input.City != nil {
		chase.City = *input.City
	}
	if input.State != nil {
		chase.State = *input.State
	}
	if input.Live != nil {
		chase.Live = *input.Live
		// If chase ended, set EndedAt
		if wasLive && !*input.Live {
			now := time.Now()
			endedAt = &now
			chase.EndedAt = endedAt
		}
	}
	if input.ThumbnailURL != nil {
		chase.ThumbnailURL = *input.ThumbnailURL
	}
	if input.Streams != nil {
		chase.Streams = input.Streams
	}
	if input.Metadata != nil {
		chase.Metadata = input.Metadata
	}

	locationJSON, _ := json.Marshal(chase.Location)
	streamsJSON, _ := json.Marshal(chase.Streams)
	metadataJSON, _ := json.Marshal(chase.Metadata)

	query := `
		UPDATE chases SET
			title = $2, description = $3, location = $4, city = $5, state = $6,
			live = $7, ended_at = $8, thumbnail_url = $9, streams = $10,
			metadata = $11, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`

	err = r.pool.QueryRow(ctx, query,
		id, chase.Title, chase.Description, locationJSON, chase.City, chase.State,
		chase.Live, endedAt, chase.ThumbnailURL, streamsJSON, metadataJSON,
	).Scan(&chase.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update chase: %w", err)
	}

	return chase, nil
}

// Delete soft-deletes a chase.
func (r *ChaseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE chases SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete chase: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// IncrementViewCount increments the view count for a chase.
func (r *ChaseRepository) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE chases SET view_count = view_count + 1 WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// IncrementShareCount increments the share count for a chase.
func (r *ChaseRepository) IncrementShareCount(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE chases SET share_count = share_count + 1 WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// GetLiveChases retrieves all currently live chases.
func (r *ChaseRepository) GetLiveChases(ctx context.Context) ([]model.Chase, error) {
	result, err := r.List(ctx, model.ChaseListOptions{
		Live:  boolPtr(true),
		Limit: 100,
	})
	if err != nil {
		return nil, err
	}
	return result.Chases, nil
}

func boolPtr(b bool) *bool {
	return &b
}
