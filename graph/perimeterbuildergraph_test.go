package graph_test

import (
	"testing"

	"github.com/heustis/tsp-solver-go/graph"
	"github.com/heustis/tsp-solver-go/model"
	"github.com/stretchr/testify/assert"
)

func TestNewGraphCircuit_ShouldBuildPerimiter(t *testing.T) {
	assert := assert.New(t)

	seed := int64(2)
	gen := graph.GraphGenerator{
		NumVertices: 20,
		MaxEdges:    6,
		MinEdges:    4,
		Seed:        &seed,
	}

	g := gen.Create()
	defer g.Delete()

	edges, unattached := graph.BuildPerimiter(graph.ToCircuitVertexArray(g.GetVertices()))

	assert.Len(edges, 7)
	assert.Len(unattached, 13)

	length := 0.0
	for _, e := range edges {
		length += e.GetLength()
	}
	assert.InDelta(50980.6999004202, length, model.Threshold)
}

func TestNewGraphCircuit_ShouldBuildPerimiter2(t *testing.T) {
	assert := assert.New(t)

	seed := int64(4)
	gen := graph.GraphGenerator{
		NumVertices: 20,
		MaxEdges:    6,
		MinEdges:    4,
		Seed:        &seed,
	}

	g := gen.Create()
	defer g.Delete()

	edges, unattached := graph.BuildPerimiter(graph.ToCircuitVertexArray(g.GetVertices()))

	assert.Len(edges, 5)
	assert.Len(unattached, 15)

	length := 0.0
	for _, e := range edges {
		length += e.GetLength()
	}
	assert.InDelta(39582.40043724108, length, model.Threshold)
}
