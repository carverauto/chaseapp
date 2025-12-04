package handler

import (
	"log/slog"
	"net/http"

	"chaseapp.tv/api/internal/external"
)

// ExternalHandler handles external data endpoints.
type ExternalHandler struct {
	client *external.Client
	logger *slog.Logger
}

// NewExternalHandler creates a new ExternalHandler.
func NewExternalHandler(client *external.Client, logger *slog.Logger) *ExternalHandler {
	return &ExternalHandler{
		client: client,
		logger: logger,
	}
}

// GetQuakes returns earthquake data from USGS.
// GET /api/v1/quakes
func (h *ExternalHandler) GetQuakes(w http.ResponseWriter, r *http.Request) {
	data, err := h.client.GetQuakes(r.Context())
	if err != nil {
		h.logger.Error("failed to fetch quakes", slog.Any("error", err))
		Error(w, http.StatusBadGateway, "Failed to fetch earthquake data")
		return
	}
	JSON(w, http.StatusOK, data)
}

// GetBoats returns vessel data from AISHub.
// GET /api/v1/boats
func (h *ExternalHandler) GetBoats(w http.ResponseWriter, r *http.Request) {
	data, err := h.client.GetBoats(r.Context())
	if err != nil {
		h.logger.Error("failed to fetch boats", slog.Any("error", err))
		Error(w, http.StatusBadGateway, "Failed to fetch vessel data")
		return
	}
	JSON(w, http.StatusOK, data)
}

// GetLaunches returns rocket launch data.
// GET /api/v1/launches
func (h *ExternalHandler) GetLaunches(w http.ResponseWriter, r *http.Request) {
	data, err := h.client.GetLaunches(r.Context())
	if err != nil {
		h.logger.Error("failed to fetch launches", slog.Any("error", err))
		Error(w, http.StatusBadGateway, "Failed to fetch launch data")
		return
	}
	JSON(w, http.StatusOK, data)
}

// GetWeatherAlerts returns active weather alerts from NOAA/NWS.
// GET /api/v1/weather/alerts
func (h *ExternalHandler) GetWeatherAlerts(w http.ResponseWriter, r *http.Request) {
	area := r.URL.Query().Get("area")
	data, err := h.client.GetWeatherAlerts(r.Context(), area)
	if err != nil {
		h.logger.Error("failed to fetch weather alerts", slog.Any("error", err))
		Error(w, http.StatusBadGateway, "Failed to fetch weather alerts")
		return
	}
	JSON(w, http.StatusOK, data)
}
