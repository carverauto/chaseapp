package worker

import (
	"context"
	"log/slog"
	"time"

	"chaseapp.tv/api/internal/model"
	"chaseapp.tv/api/internal/realtime"
	"chaseapp.tv/api/internal/search"
)

// IndexerWorker listens for chase events and updates Typesense.
type IndexerWorker struct {
	subscriber *realtime.Subscriber
	search     *search.Client
	logger     *slog.Logger
}

// NewIndexerWorker creates a new indexer worker.
func NewIndexerWorker(subscriber *realtime.Subscriber, searchClient *search.Client, logger *slog.Logger) *IndexerWorker {
	return &IndexerWorker{
		subscriber: subscriber,
		search:     searchClient,
		logger:     logger,
	}
}

// Start subscribes to chase events and blocks until context cancellation.
func (w *IndexerWorker) Start(ctx context.Context) error {
	if w.subscriber == nil || w.search == nil {
		return nil
	}

	err := w.subscriber.SubscribeChases(func(event string, chase *model.Chase) {
		if chase == nil {
			return
		}
		switch event {
		case realtime.SubjectChaseCreated, realtime.SubjectChaseUpdated, realtime.SubjectChaseLive:
			if err := w.search.UpsertChase(context.Background(), chase); err != nil {
				w.logger.Warn("indexer upsert failed", slog.Any("error", err))
			}
		case realtime.SubjectChaseDeleted, realtime.SubjectChaseEnded:
			if err := w.search.DeleteChase(context.Background(), chase.ID); err != nil {
				w.logger.Warn("indexer delete failed", slog.Any("error", err))
			}
		}
	})
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

// StatsWorker placeholder for stats aggregation.
func StatsWorker(ctx context.Context, logger *slog.Logger) {
	RunInterval(ctx, 10*time.Minute, func(ctx context.Context) {
		logger.Info("stats worker tick")
	})
}
