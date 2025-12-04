package worker

import (
	"context"
	"log/slog"
	"time"
)

// WeatherWorker periodically polls weather alerts or external data.
func WeatherWorker(ctx context.Context, logger *slog.Logger) {
	RunInterval(ctx, 10*time.Minute, func(ctx context.Context) {
		logger.Info("weather polling tick")
		// TODO: implement weather polling and storage
	})
}
