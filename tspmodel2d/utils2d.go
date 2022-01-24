package tspmodel2d

import (
	"math/rand"
	"sort"
	"time"

	"github.com/fealos/lee-tsp-go/tspmodel"
)

// DeduplicateVertices returns a copy of the supplied array (or slice) with duplicates removed.
// This sorts the vertices in the array by X, then Y, in the source array, which modifies its ordering.
// If the source data has vertices at {X-Threshold, Y}, {X, Y}, and {X+Threshold, Y}, both {X-Threshold, Y} and {X+Threshold, Y} will be in the deduplicated set.
func DeduplicateVertices(vertices []tspmodel.CircuitVertex) []tspmodel.CircuitVertex {
	// Note: we aren't using a set for deduplication due to using the Threshold for equality checks
	uniqueVertices := make([]tspmodel.CircuitVertex, 0, len(vertices))

	// Sort by X (then Y for same X)
	sort.Slice(vertices, func(indexA int, indexB int) bool {
		vA := vertices[indexA].(*Vertex2D)
		vB := vertices[indexB].(*Vertex2D)
		return vA.X < vB.X || (vA.X <= vB.X+tspmodel.Threshold && vA.Y <= vB.Y)
	})

	// traverse the sorted listed, adding unique points to the deduplicated list
	for sourceIndex := 0; sourceIndex < len(vertices); {
		v := vertices[sourceIndex]
		uniqueVertices = append(uniqueVertices, v)

		// skip indixes until we encounter a vertex that is sufficiently different from the current vertex to be considered unique
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

// GenerateVertices creates a new array of 2-dimensional vertices, containing the specified number of vertices.
func GenerateVertices(size int) []tspmodel.CircuitVertex {
	var vertices []tspmodel.CircuitVertex
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		vertices = append(vertices, NewVertex2D(r.Float64()*10000, r.Float64()*10000))
	}
	return vertices
}
