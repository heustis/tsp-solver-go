package graph_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/graph"
	"github.com/fealos/lee-tsp-go/model"
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

	circuit := graph.NewGraphCircuit(g)
	defer circuit.Delete()

	perimiter := circuit.GetAttachedVertices()
	assert.Len(perimiter, 7)
	unattached := circuit.GetUnattachedVertices()
	assert.Len(unattached, 13)

	assert.InDelta(50980.6999004202, circuit.GetLength(), model.Threshold)
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

	circuit := graph.NewGraphCircuit(g)
	defer circuit.Delete()

	perimiter := circuit.GetAttachedVertices()
	assert.Len(perimiter, 5)
	unattached := circuit.GetUnattachedVertices()
	assert.Len(unattached, 15)

	assert.InDelta(39582.40043724108, circuit.GetLength(), model.Threshold)
}
