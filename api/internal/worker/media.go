package worker

import (
	"context"
	"log/slog"
	"time"
)

// MediaWorker is a placeholder for MP4 link extraction.
func MediaWorker(ctx context.Context, logger *slog.Logger) {
	RunInterval(ctx, 20*time.Minute, func(ctx context.Context) {
		logger.Info("media extraction tick")
		// TODO: implement MP4 link extraction logic
	})
}
