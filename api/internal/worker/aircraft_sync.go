// Package worker contains background workers.
package worker

import (
	"context"
	"log/slog"
	"time"

	"chaseapp.tv/api/internal/repository"
)

// AircraftSyncWorker periodically prunes stale aircraft data.
type AircraftSyncWorker struct {
	repo     *repository.AircraftRepository
	logger   *slog.Logger
	interval time.Duration
	ttl      time.Duration

	stop    chan struct{}
	stopped chan struct{}
}

// NewAircraftSyncWorker creates a new worker with default intervals.
func NewAircraftSyncWorker(repo *repository.AircraftRepository, logger *slog.Logger) *AircraftSyncWorker {
	return &AircraftSyncWorker{
		repo:     repo,
		logger:   logger,
		interval: 5 * time.Minute,
		ttl:      15 * time.Minute,
		stop:     make(chan struct{}),
		stopped:  make(chan struct{}),
	}
}

// Start launches the worker loop.
func (w *AircraftSyncWorker) Start() {
	go w.run()
}

// Stop signals the worker to stop and waits for completion or context timeout.
func (w *AircraftSyncWorker) Stop(ctx context.Context) {
	select {
	case <-w.stopped:
		return
	case <-w.stop:
		return
	default:
		close(w.stop)
	}

	select {
	case <-w.stopped:
	case <-ctx.Done():
	}
}

func (w *AircraftSyncWorker) run() {
	defer close(w.stopped)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Run once on start
	w.prune(time.Now())

	for {
		select {
		case t := <-ticker.C:
			w.prune(t)
		case <-w.stop:
			return
		}
	}
}

func (w *AircraftSyncWorker) prune(now time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	before := now.Add(-w.ttl)
	removed, err := w.repo.DeleteStale(ctx, before)
	if err != nil {
		w.logger.Error("failed to prune stale aircraft", slog.Any("error", err))
		return
	}

	if removed > 0 {
		w.logger.Info("pruned stale aircraft", slog.Int64("removed", removed))
	}
}
