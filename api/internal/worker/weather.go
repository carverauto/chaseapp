package worker

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"chaseapp.tv/api/internal/external"
)

// WeatherWorker periodically polls weather alerts or external data.
type WeatherWorker struct {
	client   *external.Client
	logger   *slog.Logger
	interval time.Duration
}

// NewWeatherWorker creates a WeatherWorker.
func NewWeatherWorker(client *external.Client, logger *slog.Logger) *WeatherWorker {
	return &WeatherWorker{
		client:   client,
		logger:   logger,
		interval: 10 * time.Minute,
	}
}

// Start begins polling for weather alerts.
func (w *WeatherWorker) Start(ctx context.Context) {
	if w.client == nil {
		return
	}
	RunInterval(ctx, w.interval, func(ctx context.Context) {
		data, err := w.client.GetWeatherAlerts(ctx, "")
		if err != nil {
			w.logger.Warn("weather polling failed", slog.Any("error", err))
			return
		}
		var payload struct {
			Features []json.RawMessage `json:"features"`
		}
		_ = json.Unmarshal(data, &payload)
		w.logger.Info("weather polling complete", slog.Int("alerts", len(payload.Features)))
	})
}
