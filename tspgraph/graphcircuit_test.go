package tspgraph_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/tspgraph"
	"github.com/fealos/lee-tsp-go/tspmodel"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerimiter_Should(t *testing.T) {
	assert := assert.New(t)

	seed := int64(2)
	gen := tspgraph.GraphGenerator{
		NumVertices: 20,
		MaxEdges:    6,
		MinEdges:    4,
		Seed:        &seed,
	}

	g := gen.Create()
	defer g.Delete()

	circuit := tspgraph.NewGraphCircuit(g)
	defer circuit.Delete()
	circuit.Prepare()
	circuit.BuildPerimiter()

	perimiter := circuit.GetAttachedVertices()
	assert.Len(perimiter, 7)
	unattached := circuit.GetUnattachedVertices()
	assert.Len(unattached, 13)

	assert.InDelta(50980.6999004202, circuit.GetLength(), tspmodel.Threshold)
}

func TestBuildPerimiter_3(t *testing.T) {
	assert := assert.New(t)

	seed := int64(4)
	gen := tspgraph.GraphGenerator{
		NumVertices: 20,
		MaxEdges:    6,
		MinEdges:    4,
		Seed:        &seed,
	}

	g := gen.Create()
	defer g.Delete()

	circuit := tspgraph.NewGraphCircuit(g)
	defer circuit.Delete()
	circuit.Prepare()
	circuit.BuildPerimiter()

	perimiter := circuit.GetAttachedVertices()
	assert.Len(perimiter, 5)
	unattached := circuit.GetUnattachedVertices()
	assert.Len(unattached, 15)

	assert.InDelta(39582.40043724108, circuit.GetLength(), tspmodel.Threshold)
}

func TestPrepare_ShouldInitializeEdgesAndUnattachedVertices(t *testing.T) {
	assert := assert.New(t)

	seed := int64(2)
	numVertices := 20
	gen := tspgraph.GraphGenerator{
		NumVertices: uint32(numVertices),
		MaxEdges:    6,
		MinEdges:    4,
		Seed:        &seed,
	}

	g := gen.Create()
	defer g.Delete()

	circuit := tspgraph.NewGraphCircuit(g)

	for start := 0; start < numVertices; start++ {
		for end := 0; end < numVertices; end++ {
			assert.Nil(circuit.GetEdgeFor(g.Vertices[start], g.Vertices[end]))
		}
	}
	assert.Nil(circuit.GetUnattachedVertices())
	assert.NotNil(circuit.GetAttachedVertices())
	assert.Len(circuit.GetAttachedVertices(), 0)

	circuit.Prepare()

	assert.Equal(0.0, circuit.GetLength())
	assert.NotNil(circuit.GetUnattachedVertices())
	assert.NotNil(circuit.GetAttachedVertices())
	assert.Len(circuit.GetAttachedVertices(), 0)

	for start := 0; start < numVertices; start++ {
		assert.True(circuit.GetUnattachedVertices()[g.Vertices[start]])
		for end := 0; end < numVertices; end++ {
			edge := circuit.GetEdgeFor(g.Vertices[start], g.Vertices[end])
			assert.NotNil(edge)
			assert.Equal(g.Vertices[start], edge.GetStart())
			assert.Equal(g.Vertices[end], edge.GetEnd())
		}
	}

	// Also test Delete here, since we can compare it to Prepare.
	circuit.Delete()
	for start := 0; start < numVertices; start++ {
		for end := 0; end < numVertices; end++ {
			assert.Nil(circuit.GetEdgeFor(g.Vertices[start], g.Vertices[end]))
		}
	}
	assert.Len(circuit.GetUnattachedVertices(), 0)
	assert.NotNil(circuit.GetAttachedVertices())
	assert.Len(circuit.GetAttachedVertices(), 0)
}
