package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()
	req := tc.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "chaseapp",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "chaseapp",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	t.Cleanup(func() { container.Terminate(context.Background()) })

	host, err := container.Host(ctx)
	require.NoError(t, err)
	mapped, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	connString := fmt.Sprintf("postgres://chaseapp:password@%s:%s/chaseapp?sslmode=disable", host, mapped.Port())
	config, err := pgxpool.ParseConfig(connString)
	require.NoError(t, err)
	pool, err := pgxpool.NewWithConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(func() { pool.Close() })

	schema := `
CREATE TABLE chases (
	id UUID PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT,
	chase_type TEXT,
	location JSONB,
	city TEXT,
	state TEXT,
	country TEXT,
	live BOOLEAN DEFAULT false,
	started_at TIMESTAMP,
	ended_at TIMESTAMP,
	thumbnail_url TEXT,
	streams JSONB,
	view_count INT DEFAULT 0,
	share_count INT DEFAULT 0,
	source TEXT,
	source_url TEXT,
	created_by UUID,
	metadata JSONB,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	deleted_at TIMESTAMP
);`

	_, err = pool.Exec(ctx, schema)
	require.NoError(t, err)

	return pool
}

func TestCountChases(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewChaseRepository(pool)
	ctx := context.Background()

	now := time.Now()
	rows := []struct {
		title string
		live  bool
	}{
		{"Chase 1", true},
		{"Chase 2", false},
		{"Chase 3", false},
	}
	for _, row := range rows {
		_, err := pool.Exec(ctx, `INSERT INTO chases (id, title, chase_type, live, created_at, updated_at) VALUES ($1, $2, 'chase', $3, $4, $4)`,
			uuid.New(), row.title, row.live, now)
		require.NoError(t, err)
	}

	total, live, err := repo.CountChases(ctx)
	require.NoError(t, err)
	require.Equal(t, 3, total)
	require.Equal(t, 1, live)
}
