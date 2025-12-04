package worker

import (
	"context"
	"encoding/json"
	"log/slog"

	"chaseapp.tv/api/internal/realtime"
)

// AircraftEventWorker listens for aircraft.updated events.
type AircraftEventWorker struct {
	subscriber *realtime.Subscriber
	logger     *slog.Logger
}

// NewAircraftEventWorker creates a new worker.
func NewAircraftEventWorker(sub *realtime.Subscriber, logger *slog.Logger) *AircraftEventWorker {
	return &AircraftEventWorker{
		subscriber: sub,
		logger:     logger,
	}
}

// Start subscribes and logs aircraft updates.
func (w *AircraftEventWorker) Start(ctx context.Context) error {
	if w.subscriber == nil {
		return nil
	}

	if err := w.subscriber.SubscribeAircraftUpdated(func(payload json.RawMessage) {
		w.logger.Info("aircraft.updated event received", slog.Int("bytes", len(payload)))
	}); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
