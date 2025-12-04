// Package main is the entry point for the ChaseApp API server.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"chaseapp.tv/api/internal/config"
	"chaseapp.tv/api/internal/database"
	"chaseapp.tv/api/internal/server"
)

func main() {
	// Load .env file if present (development)
	_ = godotenv.Load()

	// Set up structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	// Initialize database connection
	ctx := context.Background()
	db, err := database.New(ctx, cfg.Database, logger)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("connected to database",
		slog.String("host", cfg.Database.Host),
		slog.String("database", cfg.Database.Database),
	)

	// Run migrations if enabled via environment
	if os.Getenv("DB_AUTO_MIGRATE") == "true" {
		migrationsPath := os.Getenv("DB_MIGRATIONS_PATH")
		if migrationsPath == "" {
			migrationsPath = "migrations"
		}
		if err := db.RunMigrations(migrationsPath); err != nil {
			logger.Error("failed to run migrations", slog.Any("error", err))
			os.Exit(1)
		}
	}

	// Create and start server
	srv := server.New(cfg, logger, db.Pool)

	// Handle graceful shutdown
	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start server in goroutine
	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	logger.Info("server started",
		slog.String("host", cfg.Server.Host),
		slog.Int("port", cfg.Server.Port),
	)

	// Wait for interrupt signal
	<-shutdownCtx.Done()
	stop()

	logger.Info("shutting down gracefully")

	// Create shutdown context with timeout
	gracefulCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(gracefulCtx); err != nil {
		logger.Error("shutdown error", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("server stopped")
}
