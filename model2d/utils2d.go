package model2d

import (
	"math"
	"sort"

	"github.com/fealos/lee-tsp-go/model"
)

func deduplicateVertices(vertices []*Vertex2D) []*Vertex2D {
	// Note: we aren't using a set for deduplication due to using the Threshold for equality checks
	uniqueVertices := make([]*Vertex2D, 0, len(vertices))

	// Sort by X (then Y for same X)
	sort.Slice(vertices, func(indexA int, indexB int) bool {
		vA := vertices[indexA]
		vB := vertices[indexB]
		return vA.X < vB.X || (vA.X <= vB.X+model.Threshold && vA.Y <= vB.Y)
	})

	// traverse the sorted listed, adding unquue points to the deduplicated list
	for sourceIndex := 0; sourceIndex < len(vertices); {
		v := vertices[sourceIndex]
		uniqueVertices = append(uniqueVertices, v)

		nextIndex := sourceIndex + 1
		for ; nextIndex < len(vertices); nextIndex++ {
			v2 := vertices[nextIndex]
			if math.Abs(v2.X-v.X) > model.Threshold || math.Abs(v2.Y-v.Y) > model.Threshold {
				break
			}
		}
		sourceIndex = nextIndex
	}
	return uniqueVertices
}

func findFarthestPoint(target *Vertex2D, points []*Vertex2D) *Vertex2D {
	var farthestPoint *Vertex2D
	farthestDistance := 0.0

	for _, point := range points {
		if distance := point.DistanceTo(target); distance > farthestDistance {
			farthestDistance = distance
			farthestPoint = point
		}
	}

	return farthestPoint
}
