package circuit_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerimeter_ConvexConcaveByEdge(t *testing.T) {
	assert := assert.New(t)
	circuit := circuit.NewConvexConcaveByEdge([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	},
		model2d.DeduplicateVertices,
		&model2d.PerimeterBuilder2D{},
		false,
	).(*circuit.ConvexConcaveByEdge)

	circuit.Prepare()
	circuit.BuildPerimiter()

	assert.Len(circuit.Vertices, 8)

	assert.Len(circuit.GetAttachedVertices(), 6)
	assert.Equal(model2d.NewVertex2D(-15, -15), circuit.GetAttachedVertices()[0])
	assert.Equal(model2d.NewVertex2D(15, -15), circuit.GetAttachedVertices()[1])
	assert.Equal(model2d.NewVertex2D(8, 5), circuit.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(9, 6), circuit.GetAttachedVertices()[3])
	assert.Equal(model2d.NewVertex2D(3, 13), circuit.GetAttachedVertices()[4])
	assert.Equal(model2d.NewVertex2D(-7, 6), circuit.GetAttachedVertices()[5])

	assert.InDelta(96.50213879006101, circuit.GetLength(), model.Threshold)

	assert.Len(circuit.GetUnattachedVertices(), 2)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[3]])
}

func TestPrepare_ConvexConcaveByEdge(t *testing.T) {
	assert := assert.New(t)
	circuit := circuit.NewConvexConcaveByEdge([]model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(-15-model.Threshold/3.0, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(-15+model.Threshold/3.0, -15-model.Threshold/3.0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(7, 6),
		model2d.NewVertex2D(-7, 6),
	},
		model2d.DeduplicateVertices,
		&model2d.PerimeterBuilder2D{},
		false,
	).(*circuit.ConvexConcaveByEdge)

	circuit.Prepare()

	assert.NotNil(circuit.Vertices)
	assert.Len(circuit.Vertices, 7)
	assert.ElementsMatch(circuit.Vertices, []model.CircuitVertex{
		model2d.NewVertex2D(-15+model.Threshold/3.0, -15-model.Threshold/3.0),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(7, 6),
		model2d.NewVertex2D(15, -15),
	})

	assert.NotNil(circuit.GetUnattachedVertices())
	assert.Len(circuit.GetUnattachedVertices(), 0)

	assert.NotNil(circuit.GetAttachedVertices())
	assert.Len(circuit.GetAttachedVertices(), 0)

	assert.Equal(0.0, circuit.GetLength())
}

func TestUpdate_ConvexConcaveByEdge(t *testing.T) {
	assert := assert.New(t)
	circuit := circuit.NewConvexConcaveByEdge([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	},
		model2d.DeduplicateVertices,
		&model2d.PerimeterBuilder2D{},
		false,
	).(*circuit.ConvexConcaveByEdge)

	circuit.Prepare()
	circuit.BuildPerimiter()

	attached := circuit.GetAttachedVertices()
	assert.Len(circuit.Vertices, 8)
	assert.Len(attached, 6)
	assert.Len(circuit.GetUnattachedVertices(), 2)

	circuit.Update(circuit.FindNextVertexAndEdge())

	assert.Len(circuit.Vertices, 8)
	attached = circuit.GetAttachedVertices()
	assert.Len(attached, 7)
	assert.Len(circuit.GetUnattachedVertices(), 1)

	assert.Equal(circuit.Vertices[0], attached[0])
	assert.Equal(circuit.Vertices[7], attached[1])
	assert.Equal(circuit.Vertices[3], attached[2])
	assert.Equal(circuit.Vertices[5], attached[3])
	assert.Equal(circuit.Vertices[6], attached[4])
	assert.Equal(circuit.Vertices[4], attached[5])
	assert.Equal(circuit.Vertices[1], attached[6])

	assert.Equal(model2d.NewVertex2D(3, 0), attached[2])
	assert.Equal(model2d.NewVertex2D(8, 5), attached[3])
	assert.Equal(model2d.NewVertex2D(9, 6), attached[4])

	circuit.Update(circuit.FindNextVertexAndEdge())

	assert.Len(circuit.Vertices, 8)
	attached = circuit.GetAttachedVertices()
	assert.Len(attached, 8)
	assert.Len(circuit.GetUnattachedVertices(), 0)

	assert.Equal(circuit.Vertices[0], attached[0])
	assert.Equal(circuit.Vertices[7], attached[1])
	assert.Equal(circuit.Vertices[2], attached[2])
	assert.Equal(circuit.Vertices[3], attached[3])
	assert.Equal(circuit.Vertices[5], attached[4])
	assert.Equal(circuit.Vertices[6], attached[5])
	assert.Equal(circuit.Vertices[4], attached[6])
	assert.Equal(circuit.Vertices[1], attached[7])

	assert.Equal(model2d.NewVertex2D(0, 0), attached[2])
	assert.Equal(model2d.NewVertex2D(3, 0), attached[3])
	assert.Equal(model2d.NewVertex2D(8, 5), attached[4])
	assert.Equal(model2d.NewVertex2D(9, 6), attached[5])

	circuit.Update(circuit.FindNextVertexAndEdge())

	assert.Len(circuit.Vertices, 8)
	assert.Len(circuit.GetAttachedVertices(), 8)
	assert.Len(circuit.GetUnattachedVertices(), 0)
}

func TestBuildPerimeter_ConvexConcave_ByEdge_WithUpdates(t *testing.T) {
	assert := assert.New(t)
	circuit := circuit.NewConvexConcaveByEdge([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	},
		model2d.DeduplicateVertices,
		&model2d.PerimeterBuilder2D{},
		true,
	).(*circuit.ConvexConcaveByEdge)

	circuit.Prepare()
	circuit.BuildPerimiter()

	assert.Len(circuit.Vertices, 8)

	assert.Len(circuit.GetAttachedVertices(), 6)
	assert.Equal(model2d.NewVertex2D(-15, -15), circuit.GetAttachedVertices()[0])
	assert.Equal(model2d.NewVertex2D(15, -15), circuit.GetAttachedVertices()[1])
	assert.Equal(model2d.NewVertex2D(8, 5), circuit.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(9, 6), circuit.GetAttachedVertices()[3])
	assert.Equal(model2d.NewVertex2D(3, 13), circuit.GetAttachedVertices()[4])
	assert.Equal(model2d.NewVertex2D(-7, 6), circuit.GetAttachedVertices()[5])

	assert.InDelta(96.50213879006101, circuit.GetLength(), model.Threshold)

	assert.Len(circuit.GetUnattachedVertices(), 2)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[3]])
}

func TestPrepare_ConvexConcave_ByEdge_WithUpdates(t *testing.T) {
	assert := assert.New(t)
	circuit := circuit.NewConvexConcaveByEdge([]model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(-15-model.Threshold/3.0, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(-15+model.Threshold/3.0, -15-model.Threshold/3.0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(7, 6),
		model2d.NewVertex2D(-7, 6),
	},
		model2d.DeduplicateVertices,
		&model2d.PerimeterBuilder2D{},
		true,
	).(*circuit.ConvexConcaveByEdge)

	circuit.Prepare()

	assert.NotNil(circuit.Vertices)
	assert.Len(circuit.Vertices, 7)
	assert.ElementsMatch(circuit.Vertices, []model.CircuitVertex{
		model2d.NewVertex2D(-15+model.Threshold/3.0, -15-model.Threshold/3.0),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(7, 6),
		model2d.NewVertex2D(15, -15),
	})

	assert.NotNil(circuit.GetUnattachedVertices())
	assert.Len(circuit.GetUnattachedVertices(), 0)

	assert.NotNil(circuit.GetAttachedVertices())
	assert.Len(circuit.GetAttachedVertices(), 0)

	assert.Equal(0.0, circuit.GetLength())
}

func TestUpdate_ConvexConcave_ByEdge_WithUpdates(t *testing.T) {
	assert := assert.New(t)
	circuit := circuit.NewConvexConcaveByEdge([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	},
		model2d.DeduplicateVertices,
		&model2d.PerimeterBuilder2D{},
		true,
	).(*circuit.ConvexConcaveByEdge)

	circuit.Prepare()
	circuit.BuildPerimiter()

	attached := circuit.GetAttachedVertices()
	assert.Len(circuit.Vertices, 8)
	assert.Len(attached, 6)
	assert.Len(circuit.GetUnattachedVertices(), 2)

	circuit.Update(circuit.FindNextVertexAndEdge())

	assert.Len(circuit.Vertices, 8)
	attached = circuit.GetAttachedVertices()
	assert.Len(attached, 7)
	assert.Len(circuit.GetUnattachedVertices(), 1)

	assert.Equal(circuit.Vertices[0], attached[0])
	assert.Equal(circuit.Vertices[7], attached[1])
	assert.Equal(circuit.Vertices[3], attached[2])
	assert.Equal(circuit.Vertices[5], attached[3])
	assert.Equal(circuit.Vertices[6], attached[4])
	assert.Equal(circuit.Vertices[4], attached[5])
	assert.Equal(circuit.Vertices[1], attached[6])

	assert.Equal(model2d.NewVertex2D(3, 0), attached[2])
	assert.Equal(model2d.NewVertex2D(8, 5), attached[3])
	assert.Equal(model2d.NewVertex2D(9, 6), attached[4])

	circuit.Update(circuit.FindNextVertexAndEdge())

	assert.Len(circuit.Vertices, 8)
	attached = circuit.GetAttachedVertices()
	assert.Len(attached, 8)
	assert.Len(circuit.GetUnattachedVertices(), 0)

	assert.Equal(circuit.Vertices[0], attached[0])
	assert.Equal(circuit.Vertices[7], attached[1])
	assert.Equal(circuit.Vertices[2], attached[2])
	assert.Equal(circuit.Vertices[3], attached[3])
	assert.Equal(circuit.Vertices[5], attached[4])
	assert.Equal(circuit.Vertices[6], attached[5])
	assert.Equal(circuit.Vertices[4], attached[6])
	assert.Equal(circuit.Vertices[1], attached[7])

	assert.Equal(model2d.NewVertex2D(0, 0), attached[2])
	assert.Equal(model2d.NewVertex2D(3, 0), attached[3])
	assert.Equal(model2d.NewVertex2D(8, 5), attached[4])
	assert.Equal(model2d.NewVertex2D(9, 6), attached[5])

	circuit.Update(circuit.FindNextVertexAndEdge())

	assert.Len(circuit.Vertices, 8)
	assert.Len(circuit.GetAttachedVertices(), 8)
	assert.Len(circuit.GetUnattachedVertices(), 0)
}
