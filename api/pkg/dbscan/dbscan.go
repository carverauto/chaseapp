// Package dbscan provides a simple DBSCAN clustering implementation for geospatial points.
package dbscan

import (
	"math"
	"strconv"
)

// Point represents a geospatial point.
type Point struct {
	ID       string
	Lat      float64
	Lng      float64
	Metadata map[string]any
}

// Cluster represents a set of points assigned to the same cluster.
type Cluster struct {
	ID     string
	Points []Point
}

// ClusterPoints performs DBSCAN clustering.
// epsMeters defines the neighborhood radius in meters.
// minPoints defines the minimum number of points to form a dense region.
func ClusterPoints(points []Point, epsMeters float64, minPoints int) []Cluster {
	if epsMeters <= 0 || minPoints <= 0 || len(points) == 0 {
		return nil
	}

	labels := make([]int, len(points)) // -1 noise, 0 unvisited, >=1 cluster id
	clusterID := 0

	for i := range points {
		if labels[i] != 0 {
			continue // already processed
		}

		neighbors := regionQuery(points, i, epsMeters)
		if len(neighbors) < minPoints {
			labels[i] = -1 // noise
			continue
		}

		clusterID++
		expandCluster(points, labels, i, neighbors, clusterID, epsMeters, minPoints)
	}

	return buildClusters(points, labels, clusterID)
}

func expandCluster(points []Point, labels []int, pointIdx int, neighbors []int, clusterID int, epsMeters float64, minPoints int) {
	labels[pointIdx] = clusterID

	for i := 0; i < len(neighbors); i++ {
		nIdx := neighbors[i]

		if labels[nIdx] == -1 {
			labels[nIdx] = clusterID // convert noise to border
		}
		if labels[nIdx] != 0 {
			continue // already assigned
		}

		labels[nIdx] = clusterID

		nNeighbors := regionQuery(points, nIdx, epsMeters)
		if len(nNeighbors) >= minPoints {
			neighbors = append(neighbors, nNeighbors...)
		}
	}
}

func regionQuery(points []Point, idx int, epsMeters float64) []int {
	target := points[idx]
	var neighbors []int
	for i := range points {
		if i == idx {
			continue
		}
		if haversineMeters(target.Lat, target.Lng, points[i].Lat, points[i].Lng) <= epsMeters {
			neighbors = append(neighbors, i)
		}
	}
	return neighbors
}

func buildClusters(points []Point, labels []int, clusterCount int) []Cluster {
	if clusterCount == 0 {
		return nil
	}

	clusters := make([]Cluster, clusterCount)
	for i := 1; i <= clusterCount; i++ {
		clusters[i-1] = Cluster{
			ID:     generateClusterID(i),
			Points: []Point{},
		}
	}

	for idx, label := range labels {
		if label <= 0 {
			continue // noise or unassigned
		}
		clusters[label-1].Points = append(clusters[label-1].Points, points[idx])
	}

	return clusters
}

func generateClusterID(id int) string {
	return "cluster-" + strconv.Itoa(id)
}

// haversineMeters calculates the great-circle distance between two points.
func haversineMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000.0 // meters

	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRadians(lat1))*math.Cos(toRadians(lat2))*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

func toRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}
