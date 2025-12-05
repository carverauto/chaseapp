package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBoundingRectangle(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	handler := NewGeoHandler(logger)

	body := `{
		"type": "FeatureCollection",
		"features": [
			{
				"type": "Feature",
				"geometry": {
					"type": "Polygon",
					"coordinates": [[[0,0],[2,0],[2,1],[0,1],[0,0]]]
				},
				"properties": {}
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/geo/bounding-rect", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler.GetBoundingRectangle(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)

	var payload map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))

	geom, ok := payload["geometry"].(map[string]any)
	require.True(t, ok, "geometry missing")
	require.Equal(t, "Polygon", geom["type"])

	coords, ok := geom["coordinates"].([]any)
	require.True(t, ok, "coordinates missing")
	require.Len(t, coords, 1)

	props, ok := payload["properties"].(map[string]any)
	require.True(t, ok, "properties missing")
	area, ok := props["area"].(float64)
	require.True(t, ok, "area missing")
	require.Greater(t, area, 0.0)
}
