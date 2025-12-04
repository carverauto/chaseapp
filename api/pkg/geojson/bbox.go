package geojson

import (
	"errors"
	"math"
)

// Rectangle represents a minimum bounding rectangle.
type Rectangle struct {
	Points [4]Point
	Area   float64
}

// MinimumBoundingRectangle computes the minimum-area rectangle containing all points
// using a convex hull + rotating calipers approach.
func MinimumBoundingRectangle(points []Point) (Rectangle, error) {
	if len(points) == 0 {
		return Rectangle{}, errors.New("no points provided")
	}
	if len(points) == 1 {
		p := points[0]
		return Rectangle{
			Points: [4]Point{p, p, p, p},
			Area:   0,
		}, nil
	}

	// Degenerate 2-point case: axis-aligned rectangle from endpoints.
	if len(points) == 2 {
		p1, p2 := points[0], points[1]
		minX := math.Min(p1.X, p2.X)
		maxX := math.Max(p1.X, p2.X)
		minY := math.Min(p1.Y, p2.Y)
		maxY := math.Max(p1.Y, p2.Y)
		return Rectangle{
			Points: [4]Point{
				{X: minX, Y: minY},
				{X: maxX, Y: minY},
				{X: maxX, Y: maxY},
				{X: minX, Y: maxY},
			},
			Area: (maxX - minX) * (maxY - minY),
		}, nil
	}

	hull := convexHull(points)
	if len(hull) < 3 {
		// Fallback to 2-point handling
		return MinimumBoundingRectangle(hull)
	}

	var best Rectangle
	best.Area = math.Inf(1)

	for i := 0; i < len(hull); i++ {
		p1 := hull[i]
		p2 := hull[(i+1)%len(hull)]

		// Angle to align edge with the x-axis.
		edgeAngle := math.Atan2(p2.Y-p1.Y, p2.X-p1.X)
		cosTheta := math.Cos(-edgeAngle)
		sinTheta := math.Sin(-edgeAngle)

		var minX, maxX = math.Inf(1), math.Inf(-1)
		var minY, maxY = math.Inf(1), math.Inf(-1)

		for _, p := range hull {
			rx := p.X*cosTheta - p.Y*sinTheta
			ry := p.X*sinTheta + p.Y*cosTheta

			if rx < minX {
				minX = rx
			}
			if rx > maxX {
				maxX = rx
			}
			if ry < minY {
				minY = ry
			}
			if ry > maxY {
				maxY = ry
			}
		}

		area := (maxX - minX) * (maxY - minY)
		if area < best.Area {
			// Define rectangle in rotated space.
			rect := [4]Point{
				{X: minX, Y: minY},
				{X: maxX, Y: minY},
				{X: maxX, Y: maxY},
				{X: minX, Y: maxY},
			}

			// Rotate back to original space.
			cosBack := math.Cos(edgeAngle)
			sinBack := math.Sin(edgeAngle)
			var unrotated [4]Point
			for idx, rp := range rect {
				x := rp.X*cosBack - rp.Y*sinBack
				y := rp.X*sinBack + rp.Y*cosBack
				unrotated[idx] = Point{X: x, Y: y}
			}

			best = Rectangle{
				Points: unrotated,
				Area:   area,
			}
		}
	}

	return best, nil
}

// convexHull computes the convex hull of a set of points using the monotonic chain algorithm.
func convexHull(points []Point) []Point {
	if len(points) <= 1 {
		return append([]Point{}, points...)
	}

	pts := append([]Point{}, points...)
	sortPoints(pts)

	var lower []Point
	for _, p := range pts {
		for len(lower) >= 2 && cross(lower[len(lower)-2], lower[len(lower)-1], p) <= 0 {
			lower = lower[:len(lower)-1]
		}
		lower = append(lower, p)
	}

	var upper []Point
	for i := len(pts) - 1; i >= 0; i-- {
		p := pts[i]
		for len(upper) >= 2 && cross(upper[len(upper)-2], upper[len(upper)-1], p) <= 0 {
			upper = upper[:len(upper)-1]
		}
		upper = append(upper, p)
	}

	// Concatenate lower and upper to get full hull, removing the last point of each (duplicate of first).
	return append(lower[:len(lower)-1], upper[:len(upper)-1]...)
}

func cross(o, a, b Point) float64 {
	return (a.X-o.X)*(b.Y-o.Y) - (a.Y-o.Y)*(b.X-o.X)
}
