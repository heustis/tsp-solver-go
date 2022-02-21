package circuit_test

import (
	"testing"

	"github.com/heustis/tsp-solver-go/circuit"
	"github.com/heustis/tsp-solver-go/model"
	"github.com/heustis/tsp-solver-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestNewConvexConcaveByEdge(t *testing.T) {
	assert := assert.New(t)
	circuit := circuit.NewConvexConcaveByEdge(model2d.DeduplicateVertices([]model.CircuitVertex{
		// Note: the circuit is sorted by DeduplicateVertices(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}),
		model2d.BuildPerimiter,
		false,
	)

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

func TestUpdate_ConvexConcaveByEdge(t *testing.T) {
	assert := assert.New(t)
	circuit := circuit.NewConvexConcaveByEdge(model2d.DeduplicateVertices([]model.CircuitVertex{
		// Note: the circuit is sorted by DeduplicateVertices(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}),
		model2d.BuildPerimiter,
		false,
	)

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

func TestNewConvexConcaveByEdge_WithUpdates(t *testing.T) {
	assert := assert.New(t)
	circuit := circuit.NewConvexConcaveByEdge(model2d.DeduplicateVertices([]model.CircuitVertex{
		// Note: the circuit is sorted by DeduplicateVertices(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}),
		model2d.BuildPerimiter,
		true,
	)

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

func TestUpdate_ConvexConcaveByEdge_WithUpdates(t *testing.T) {
	assert := assert.New(t)
	circuit := circuit.NewConvexConcaveByEdge(model2d.DeduplicateVertices([]model.CircuitVertex{
		// Note: the circuit is sorted by DeduplicateVertices(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}),
		model2d.BuildPerimiter,
		true,
	)

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
