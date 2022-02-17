package graph_test

import (
	"testing"

	"github.com/heustis/lee-tsp-go/graph"
	"github.com/stretchr/testify/assert"
)

func TestNewGraph_ShouldCreateAnEdgeFromEveryVertexToEveryOtherVertex(t *testing.T) {
	assert := assert.New(t)

	gen := &graph.GraphGenerator{
		MaxEdges:    6,
		MinEdges:    2,
		NumVertices: uint32(20),
	}

	g := gen.Create()
	assert.NotNil(g)
	defer g.Delete()

	for _, start := range g.GetVertices() {
		for _, destination := range g.GetVertices() {
			edge := start.EdgeTo(destination)
			assert.NotNil(edge)
			assert.Equal(start, edge.GetStart())
			assert.Equal(destination, edge.GetEnd())
		}
	}
}
