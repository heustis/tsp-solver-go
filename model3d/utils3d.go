package model3d

import (
	"math/rand"
	"time"

	"github.com/heustis/tsp-solver-go/model"
)

// GenerateVertices creates a new array of 3-dimensional vertices, containing the specified number of vertices.
func GenerateVertices(size int) []model.CircuitVertex {
	var vertices []model.CircuitVertex
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		vertices = append(vertices, NewVertex3D(r.Float64()*10000, r.Float64()*10000, r.Float64()*10000))
	}
	return model.DeduplicateVerticesNoSorting(vertices)
}
