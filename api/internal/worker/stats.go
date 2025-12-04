package worker

import (
	"context"
	"log/slog"
	"time"
)

// StatsWorker runs periodic statistics aggregation.
func StatsWorker(ctx context.Context, logger *slog.Logger) {
	RunInterval(ctx, 15*time.Minute, func(ctx context.Context) {
		logger.Info("stats aggregation tick")
		// TODO: implement stats aggregation logic
	})
}
