package graph

import "github.com/fealos/lee-tsp-go/model"

func ToCircuitVertexArray(g []*GraphVertex) []model.CircuitVertex {
	vertices := make([]model.CircuitVertex, len(g))
	for i, v := range g {
		vertices[i] = v
	}
	return vertices
}