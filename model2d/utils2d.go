package model2d

import (
	"sort"

	"github.com/fealos/lee-tsp-go/model"
)

func DeduplicateVertices(vertices []model.CircuitVertex) []model.CircuitVertex {
	// Note: we aren't using a set for deduplication due to using the Threshold for equality checks
	uniqueVertices := make([]model.CircuitVertex, 0, len(vertices))

	// Sort by X (then Y for same X)
	sort.Slice(vertices, func(indexA int, indexB int) bool {
		vA := vertices[indexA].(*Vertex2D)
		vB := vertices[indexB].(*Vertex2D)
		return vA.X < vB.X || (vA.X <= vB.X+model.Threshold && vA.Y <= vB.Y)
	})

	// traverse the sorted listed, adding unique points to the deduplicated list
	for sourceIndex := 0; sourceIndex < len(vertices); {
		v := vertices[sourceIndex]
		uniqueVertices = append(uniqueVertices, v)

		nextIndex := sourceIndex + 1
		for ; nextIndex < len(vertices); nextIndex++ {
			if !v.Equals(vertices[nextIndex]) {
				break
			}
		}
		sourceIndex = nextIndex
	}
	return uniqueVertices
}

func findFarthestPoint(target model.CircuitVertex, points []model.CircuitVertex) model.CircuitVertex {
	var farthestPoint model.CircuitVertex
	farthestDistance := 0.0

	for _, point := range points {
		if distance := point.DistanceTo(target); distance > farthestDistance {
			farthestDistance = distance
			farthestPoint = point
		}
	}

	return farthestPoint
}
