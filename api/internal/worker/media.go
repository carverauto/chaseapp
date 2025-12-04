package worker

import (
	"context"
	"log/slog"
	"time"

	"chaseapp.tv/api/internal/model"
	"chaseapp.tv/api/internal/repository"
	"chaseapp.tv/api/pkg/scraper"
)

// MediaWorker attempts to extract streams for live chases missing streams.
type MediaWorker struct {
	repo      *repository.ChaseRepository
	extractor *scraper.Extractor
	logger    *slog.Logger
	interval  time.Duration
}

// NewMediaWorker creates a MediaWorker.
func NewMediaWorker(repo *repository.ChaseRepository, extractor *scraper.Extractor, logger *slog.Logger) *MediaWorker {
	return &MediaWorker{
		repo:      repo,
		extractor: extractor,
		logger:    logger,
		interval:  20 * time.Minute,
	}
}

// Start begins periodic extraction.
func (w *MediaWorker) Start(ctx context.Context) {
	if w.repo == nil || w.extractor == nil {
		return
	}

	RunInterval(ctx, w.interval, func(ctx context.Context) {
		liveChases, err := w.repo.GetLiveChases(ctx)
		if err != nil {
			w.logger.Warn("media worker failed to load live chases", slog.Any("error", err))
			return
		}
		for _, chase := range liveChases {
			if len(chase.Streams) > 0 || chase.SourceURL == "" {
				continue
			}
			streams, err := w.extractor.Extract(ctx, chase.SourceURL)
			if err != nil {
				w.logger.Warn("media extraction failed", slog.Any("error", err), slog.String("chase_id", chase.ID.String()))
				continue
			}
			if len(streams) == 0 {
				continue
			}
			_, _, err = w.repo.Update(ctx, chase.ID, model.UpdateChaseInput{
				Streams: streams,
			})
			if err != nil {
				w.logger.Warn("failed to save extracted streams", slog.Any("error", err), slog.String("chase_id", chase.ID.String()))
				continue
			}
			w.logger.Info("extracted streams for chase", slog.String("chase_id", chase.ID.String()), slog.Int("streams", len(streams)))
		}
	})
}
