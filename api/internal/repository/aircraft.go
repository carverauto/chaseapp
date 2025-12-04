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

// AircraftRepository handles aircraft data access.
type AircraftRepository struct {
	pool *pgxpool.Pool
}

// NewAircraftRepository creates a new AircraftRepository.
func NewAircraftRepository(pool *pgxpool.Pool) *AircraftRepository {
	return &AircraftRepository{pool: pool}
}

// Upsert creates or updates an aircraft by ICAO code.
func (r *AircraftRepository) Upsert(ctx context.Context, input model.UpsertAircraftInput) (*model.Aircraft, error) {
	now := time.Now()

	metadataJSON, err := json.Marshal(input.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		INSERT INTO aircraft (
			icao, callsign, registration, latitude, longitude, altitude,
			ground_speed, track, vertical_rate, aircraft_type, category,
			operator, on_ground, squawk, emergency, metadata, last_seen_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
		ON CONFLICT (icao) DO UPDATE SET
			callsign = COALESCE(EXCLUDED.callsign, aircraft.callsign),
			registration = COALESCE(EXCLUDED.registration, aircraft.registration),
			latitude = COALESCE(EXCLUDED.latitude, aircraft.latitude),
			longitude = COALESCE(EXCLUDED.longitude, aircraft.longitude),
			altitude = COALESCE(EXCLUDED.altitude, aircraft.altitude),
			ground_speed = COALESCE(EXCLUDED.ground_speed, aircraft.ground_speed),
			track = COALESCE(EXCLUDED.track, aircraft.track),
			vertical_rate = COALESCE(EXCLUDED.vertical_rate, aircraft.vertical_rate),
			aircraft_type = COALESCE(EXCLUDED.aircraft_type, aircraft.aircraft_type),
			category = COALESCE(EXCLUDED.category, aircraft.category),
			operator = COALESCE(EXCLUDED.operator, aircraft.operator),
			on_ground = EXCLUDED.on_ground,
			squawk = COALESCE(EXCLUDED.squawk, aircraft.squawk),
			emergency = EXCLUDED.emergency,
			metadata = COALESCE(EXCLUDED.metadata, aircraft.metadata),
			last_seen_at = $17,
			updated_at = NOW()
		RETURNING id, icao, callsign, registration, latitude, longitude, altitude,
			ground_speed, track, vertical_rate, aircraft_type, category, operator,
			on_ground, squawk, emergency, cluster_id, metadata,
			first_seen_at, last_seen_at, created_at, updated_at`

	var aircraft model.Aircraft
	var metadataBytes []byte

	err = r.pool.QueryRow(ctx, query,
		input.ICAO, input.Callsign, input.Registration, input.Latitude, input.Longitude,
		input.Altitude, input.GroundSpeed, input.Track, input.VerticalRate,
		input.AircraftType, input.Category, input.Operator, input.OnGround,
		input.Squawk, input.Emergency, metadataJSON, now,
	).Scan(
		&aircraft.ID, &aircraft.ICAO, &aircraft.Callsign, &aircraft.Registration,
		&aircraft.Latitude, &aircraft.Longitude, &aircraft.Altitude,
		&aircraft.GroundSpeed, &aircraft.Track, &aircraft.VerticalRate,
		&aircraft.AircraftType, &aircraft.Category, &aircraft.Operator,
		&aircraft.OnGround, &aircraft.Squawk, &aircraft.Emergency,
		&aircraft.ClusterID, &metadataBytes,
		&aircraft.FirstSeenAt, &aircraft.LastSeenAt, &aircraft.CreatedAt, &aircraft.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to upsert aircraft: %w", err)
	}

	if len(metadataBytes) > 0 {
		json.Unmarshal(metadataBytes, &aircraft.Metadata)
	}

	return &aircraft, nil
}

// GetByID retrieves an aircraft by ID.
func (r *AircraftRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Aircraft, error) {
	query := `
		SELECT id, icao, callsign, registration, latitude, longitude, altitude,
			   ground_speed, track, vertical_rate, aircraft_type, category, operator,
			   on_ground, squawk, emergency, cluster_id, metadata,
			   first_seen_at, last_seen_at, created_at, updated_at
		FROM aircraft
		WHERE id = $1`

	return r.scanAircraft(r.pool.QueryRow(ctx, query, id))
}

// GetByICAO retrieves an aircraft by ICAO code.
func (r *AircraftRepository) GetByICAO(ctx context.Context, icao string) (*model.Aircraft, error) {
	query := `
		SELECT id, icao, callsign, registration, latitude, longitude, altitude,
			   ground_speed, track, vertical_rate, aircraft_type, category, operator,
			   on_ground, squawk, emergency, cluster_id, metadata,
			   first_seen_at, last_seen_at, created_at, updated_at
		FROM aircraft
		WHERE icao = $1`

	return r.scanAircraft(r.pool.QueryRow(ctx, query, icao))
}

// List retrieves aircraft with pagination and filtering.
func (r *AircraftRepository) List(ctx context.Context, opts model.AircraftListOptions) (*model.AircraftListResult, error) {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.Limit < 1 || opts.Limit > 100 {
		opts.Limit = 50
	}

	offset := (opts.Page - 1) * opts.Limit

	baseQuery := `FROM aircraft WHERE 1=1`
	args := []interface{}{}
	argNum := 1

	if opts.Category != "" {
		baseQuery += fmt.Sprintf(" AND category = $%d", argNum)
		args = append(args, opts.Category)
		argNum++
	}
	if opts.ClusterID != "" {
		baseQuery += fmt.Sprintf(" AND cluster_id = $%d", argNum)
		args = append(args, opts.ClusterID)
		argNum++
	}
	if opts.OnGround != nil {
		baseQuery += fmt.Sprintf(" AND on_ground = $%d", argNum)
		args = append(args, *opts.OnGround)
		argNum++
	}

	// Geographic bounding box filter
	if opts.MinLat != nil && opts.MaxLat != nil && opts.MinLng != nil && opts.MaxLng != nil {
		baseQuery += fmt.Sprintf(" AND latitude BETWEEN $%d AND $%d AND longitude BETWEEN $%d AND $%d",
			argNum, argNum+1, argNum+2, argNum+3)
		args = append(args, *opts.MinLat, *opts.MaxLat, *opts.MinLng, *opts.MaxLng)
		argNum += 4
	}

	// Get total count
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count aircraft: %w", err)
	}

	// Get aircraft
	selectQuery := fmt.Sprintf(`
		SELECT id, icao, callsign, registration, latitude, longitude, altitude,
			   ground_speed, track, vertical_rate, aircraft_type, category, operator,
			   on_ground, squawk, emergency, cluster_id, metadata,
			   first_seen_at, last_seen_at, created_at, updated_at
		%s ORDER BY last_seen_at DESC LIMIT $%d OFFSET $%d`,
		baseQuery, argNum, argNum+1)

	args = append(args, opts.Limit, offset)

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list aircraft: %w", err)
	}
	defer rows.Close()

	var aircraft []model.Aircraft
	for rows.Next() {
		a, err := r.scanAircraftRow(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan aircraft: %w", err)
		}
		aircraft = append(aircraft, *a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate aircraft: %w", err)
	}

	totalPages := (total + opts.Limit - 1) / opts.Limit

	return &model.AircraftListResult{
		Aircraft:   aircraft,
		Total:      total,
		Page:       opts.Page,
		Limit:      opts.Limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateCluster assigns aircraft to a cluster.
func (r *AircraftRepository) UpdateCluster(ctx context.Context, id uuid.UUID, clusterID string) error {
	query := `UPDATE aircraft SET cluster_id = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, clusterID)
	return err
}

// Delete removes an aircraft.
func (r *AircraftRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM aircraft WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete aircraft: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteStale removes aircraft not seen since the given time.
func (r *AircraftRepository) DeleteStale(ctx context.Context, since time.Time) (int64, error) {
	query := `DELETE FROM aircraft WHERE last_seen_at < $1`
	result, err := r.pool.Exec(ctx, query, since)
	if err != nil {
		return 0, fmt.Errorf("failed to delete stale aircraft: %w", err)
	}
	return result.RowsAffected(), nil
}

// AddHistory records a position in aircraft history.
func (r *AircraftRepository) AddHistory(ctx context.Context, aircraftID uuid.UUID, lat, lng float64, altitude, speed, track *int) error {
	query := `
		INSERT INTO aircraft_history (aircraft_id, latitude, longitude, altitude, ground_speed, track)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, query, aircraftID, lat, lng, altitude, speed, track)
	return err
}

// GetHistory retrieves position history for an aircraft.
func (r *AircraftRepository) GetHistory(ctx context.Context, aircraftID uuid.UUID, limit int) ([]model.AircraftHistory, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	query := `
		SELECT id, aircraft_id, latitude, longitude, altitude, ground_speed, track, recorded_at
		FROM aircraft_history
		WHERE aircraft_id = $1
		ORDER BY recorded_at DESC
		LIMIT $2`

	rows, err := r.pool.Query(ctx, query, aircraftID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get aircraft history: %w", err)
	}
	defer rows.Close()

	var history []model.AircraftHistory
	for rows.Next() {
		var h model.AircraftHistory
		err := rows.Scan(&h.ID, &h.AircraftID, &h.Latitude, &h.Longitude,
			&h.Altitude, &h.GroundSpeed, &h.Track, &h.RecordedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan aircraft history: %w", err)
		}
		history = append(history, h)
	}

	return history, rows.Err()
}

// DeleteOldHistory removes history older than the given time.
func (r *AircraftRepository) DeleteOldHistory(ctx context.Context, before time.Time) (int64, error) {
	query := `DELETE FROM aircraft_history WHERE recorded_at < $1`
	result, err := r.pool.Exec(ctx, query, before)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old history: %w", err)
	}
	return result.RowsAffected(), nil
}

// Helper to scan a single aircraft from a row.
func (r *AircraftRepository) scanAircraft(row pgx.Row) (*model.Aircraft, error) {
	var aircraft model.Aircraft
	var metadataBytes []byte

	err := row.Scan(
		&aircraft.ID, &aircraft.ICAO, &aircraft.Callsign, &aircraft.Registration,
		&aircraft.Latitude, &aircraft.Longitude, &aircraft.Altitude,
		&aircraft.GroundSpeed, &aircraft.Track, &aircraft.VerticalRate,
		&aircraft.AircraftType, &aircraft.Category, &aircraft.Operator,
		&aircraft.OnGround, &aircraft.Squawk, &aircraft.Emergency,
		&aircraft.ClusterID, &metadataBytes,
		&aircraft.FirstSeenAt, &aircraft.LastSeenAt, &aircraft.CreatedAt, &aircraft.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get aircraft: %w", err)
	}

	if len(metadataBytes) > 0 {
		json.Unmarshal(metadataBytes, &aircraft.Metadata)
	}

	return &aircraft, nil
}

// Helper to scan aircraft from rows.
func (r *AircraftRepository) scanAircraftRow(rows pgx.Rows) (*model.Aircraft, error) {
	var aircraft model.Aircraft
	var metadataBytes []byte

	err := rows.Scan(
		&aircraft.ID, &aircraft.ICAO, &aircraft.Callsign, &aircraft.Registration,
		&aircraft.Latitude, &aircraft.Longitude, &aircraft.Altitude,
		&aircraft.GroundSpeed, &aircraft.Track, &aircraft.VerticalRate,
		&aircraft.AircraftType, &aircraft.Category, &aircraft.Operator,
		&aircraft.OnGround, &aircraft.Squawk, &aircraft.Emergency,
		&aircraft.ClusterID, &metadataBytes,
		&aircraft.FirstSeenAt, &aircraft.LastSeenAt, &aircraft.CreatedAt, &aircraft.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if len(metadataBytes) > 0 {
		json.Unmarshal(metadataBytes, &aircraft.Metadata)
	}

	return &aircraft, nil
}
