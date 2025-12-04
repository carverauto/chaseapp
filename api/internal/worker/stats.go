package worker

import (
	"context"
	"log/slog"
	"time"

	"chaseapp.tv/api/internal/repository"
)

// StatsWorker runs periodic statistics aggregation.
type StatsWorker struct {
	repo     *repository.ChaseRepository
	logger   *slog.Logger
	interval time.Duration
}

// NewStatsWorker creates a StatsWorker.
func NewStatsWorker(repo *repository.ChaseRepository, logger *slog.Logger) *StatsWorker {
	return &StatsWorker{
		repo:     repo,
		logger:   logger,
		interval: 15 * time.Minute,
	}
}

// Start begins periodic aggregation.
func (w *StatsWorker) Start(ctx context.Context) {
	RunInterval(ctx, w.interval, func(ctx context.Context) {
		total, live, err := w.repo.CountChases(ctx)
		if err != nil {
			w.logger.Warn("stats aggregation failed", slog.Any("error", err))
			return
		}
		w.logger.Info("stats aggregation",
			slog.Int("total_chases", total),
			slog.Int("live_chases", live),
		)
	})
}
