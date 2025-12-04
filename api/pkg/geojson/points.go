package geojson

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Point represents a 2D coordinate (X = longitude, Y = latitude).
type Point struct {
	X float64
	Y float64
}

// ExtractPoints parses GeoJSON input and returns all coordinate points it contains.
func ExtractPoints(data []byte) ([]Point, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	points, err := parseGeoJSON(raw)
	if err != nil {
		return nil, err
	}
	if len(points) == 0 {
		return nil, errors.New("no coordinates provided")
	}
	return points, nil
}

func parseGeoJSON(obj map[string]any) ([]Point, error) {
	typ, ok := obj["type"].(string)
	if !ok {
		return nil, errors.New("geojson type is required")
	}

	switch typ {
	case "FeatureCollection":
		features, ok := obj["features"].([]any)
		if !ok {
			return nil, errors.New("features must be an array")
		}
		var points []Point
		for _, f := range features {
			fm, ok := f.(map[string]any)
			if !ok {
				continue
			}
			if fm["type"] == "Feature" {
				pts, err := parseFeature(fm)
				if err != nil {
					return nil, err
				}
				points = append(points, pts...)
			}
		}
		return points, nil

	case "Feature":
		return parseFeature(obj)

	default:
		return parseGeometry(obj)
	}
}

func parseFeature(obj map[string]any) ([]Point, error) {
	geom, ok := obj["geometry"].(map[string]any)
	if !ok {
		return nil, errors.New("feature.geometry is required")
	}
	return parseGeometry(geom)
}

func parseGeometry(obj map[string]any) ([]Point, error) {
	typ, ok := obj["type"].(string)
	if !ok {
		return nil, errors.New("geometry type is required")
	}

	coords := obj["coordinates"]

	switch typ {
	case "Point":
		p, err := coordPair(coords)
		if err != nil {
			return nil, err
		}
		return []Point{p}, nil

	case "MultiPoint", "LineString":
		arr, ok := coords.([]any)
		if !ok {
			return nil, errors.New("coordinates must be an array")
		}
		var pts []Point
		for _, c := range arr {
			p, err := coordPair(c)
			if err != nil {
				return nil, err
			}
			pts = append(pts, p)
		}
		return pts, nil

	case "MultiLineString", "Polygon":
		arr, ok := coords.([]any)
		if !ok {
			return nil, errors.New("coordinates must be an array")
		}
		var pts []Point
		for _, ring := range arr {
			ringArr, ok := ring.([]any)
			if !ok {
				continue
			}
			for _, c := range ringArr {
				p, err := coordPair(c)
				if err != nil {
					return nil, err
				}
				pts = append(pts, p)
			}
		}
		return pts, nil

	case "MultiPolygon":
		arr, ok := coords.([]any)
		if !ok {
			return nil, errors.New("coordinates must be an array")
		}
		var pts []Point
		for _, poly := range arr {
			polyArr, ok := poly.([]any)
			if !ok {
				continue
			}
			for _, ring := range polyArr {
				ringArr, ok := ring.([]any)
				if !ok {
					continue
				}
				for _, c := range ringArr {
					p, err := coordPair(c)
					if err != nil {
						return nil, err
					}
					pts = append(pts, p)
				}
			}
		}
		return pts, nil
	default:
		return nil, fmt.Errorf("unsupported geometry type: %s", typ)
	}
}

func coordPair(v any) (Point, error) {
	arr, ok := v.([]any)
	if !ok || len(arr) < 2 {
		return Point{}, errors.New("coordinate must be [x, y]")
	}
	lon, ok1 := toFloat(arr[0])
	lat, ok2 := toFloat(arr[1])
	if !ok1 || !ok2 {
		return Point{}, errors.New("coordinates must be numbers")
	}
	return Point{X: lon, Y: lat}, nil
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}
