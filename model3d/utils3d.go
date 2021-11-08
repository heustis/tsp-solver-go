package model3d

import (
	"github.com/fealos/lee-tsp-go/model"
)

func DeduplicateVertices3D(vertices []model.CircuitVertex) []model.CircuitVertex {
	// Note: we aren't using a set for deduplication due to using the Threshold for equality checks
	uniqueVertices := make([]model.CircuitVertex, 0, len(vertices))

	// Unlike 2D, sorting 3D vertices can result in duplicate entries in the output, due to having to sort by Y or Z first.
	// If we sort by Y first, it is possible to have a point that matches in X and Y and is significantly different in Z in between points that match in X, Y, and Z.
	// So, we are doing an O(n*n) deduplication, by checking each added point to see if it is a duplicate of the current point.
	for _, v := range vertices {
		shouldAdd := true
		for _, added := range uniqueVertices {
			if v.Equals(added) {
				shouldAdd = false
				break
			}
		}

		if shouldAdd {
			uniqueVertices = append(uniqueVertices, v)
		}
	}
	return uniqueVertices
}
