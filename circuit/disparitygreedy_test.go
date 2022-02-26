package circuit_test

import (
	"testing"

	"github.com/heustis/tsp-solver-go/circuit"
	"github.com/heustis/tsp-solver-go/model"
	"github.com/heustis/tsp-solver-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestNewDisparityGreedy(t *testing.T) {
	assert := assert.New(t)
	vertices := model2d.DeduplicateVertices([]model.CircuitVertex{
		// Note: the circuit is sorted by DeduplicateVertices(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	})
	c := circuit.NewDisparityGreedy(vertices, model2d.BuildPerimiter, false)

	assert.Len(c.GetAttachedVertices(), 5)
	assert.Equal(model2d.NewVertex2D(-15, -15), c.GetAttachedVertices()[0])
	assert.Equal(model2d.NewVertex2D(15, -15), c.GetAttachedVertices()[1])
	assert.Equal(model2d.NewVertex2D(9, 6), c.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(3, 13), c.GetAttachedVertices()[3])
	assert.Equal(model2d.NewVertex2D(-7, 6), c.GetAttachedVertices()[4])

	assert.InDelta(95.73863479511238, c.GetLength(), model.Threshold)

	assert.Len(c.GetAttachedEdges(), 5)
	assert.True(vertices[0].EdgeTo(vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(vertices[7].EdgeTo(vertices[6]).Equals(c.GetAttachedEdges()[1]))
	assert.True(vertices[6].EdgeTo(vertices[4]).Equals(c.GetAttachedEdges()[2]))
	assert.True(vertices[4].EdgeTo(vertices[1]).Equals(c.GetAttachedEdges()[3]))
	assert.True(vertices[1].EdgeTo(vertices[0]).Equals(c.GetAttachedEdges()[4]))

	assert.Len(c.GetUnattachedVertices(), 3)
	assert.True(c.GetUnattachedVertices()[vertices[2]])
	assert.True(c.GetUnattachedVertices()[vertices[3]])
	assert.True(c.GetUnattachedVertices()[vertices[5]])
}

func TestUpdate_DisparityGreedy(t *testing.T) {
	assert := assert.New(t)
	vertices := model2d.DeduplicateVertices([]model.CircuitVertex{
		// Note: the circuit is sorted by DeduplicateVertices(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	})
	c := circuit.NewDisparityGreedy(vertices, model2d.BuildPerimiter, false)

	assert.Len(c.GetAttachedVertices(), 5)
	assert.Len(c.GetAttachedEdges(), 5)
	assert.Len(c.GetUnattachedVertices(), 3)

	// Init circuit := (-15, -15), (15, -15), (9,6), (3,13), (-7,6)

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.GetAttachedVertices(), 6)
	assert.Len(c.GetAttachedEdges(), 6)
	assert.Len(c.GetUnattachedVertices(), 2)

	assert.True(vertices[0].EdgeTo(vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(vertices[7].EdgeTo(vertices[3]).Equals(c.GetAttachedEdges()[1]))
	assert.True(vertices[3].EdgeTo(vertices[6]).Equals(c.GetAttachedEdges()[2]))
	assert.True(vertices[6].EdgeTo(vertices[4]).Equals(c.GetAttachedEdges()[3]))
	assert.True(vertices[4].EdgeTo(vertices[1]).Equals(c.GetAttachedEdges()[4]))
	assert.True(vertices[1].EdgeTo(vertices[0]).Equals(c.GetAttachedEdges()[5]))

	assert.Equal(model2d.NewVertex2D(3, 0), c.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(9, 6), c.GetAttachedVertices()[3])

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.GetAttachedVertices(), 7)
	assert.Len(c.GetAttachedEdges(), 7)
	assert.Len(c.GetUnattachedVertices(), 1)

	assert.True(vertices[0].EdgeTo(vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(vertices[7].EdgeTo(vertices[3]).Equals(c.GetAttachedEdges()[1]))
	assert.True(vertices[3].EdgeTo(vertices[5]).Equals(c.GetAttachedEdges()[2]))
	assert.True(vertices[5].EdgeTo(vertices[6]).Equals(c.GetAttachedEdges()[3]))
	assert.True(vertices[6].EdgeTo(vertices[4]).Equals(c.GetAttachedEdges()[4]))
	assert.True(vertices[4].EdgeTo(vertices[1]).Equals(c.GetAttachedEdges()[5]))
	assert.True(vertices[1].EdgeTo(vertices[0]).Equals(c.GetAttachedEdges()[6]))

	assert.Equal(model2d.NewVertex2D(3, 0), c.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(8, 5), c.GetAttachedVertices()[3])
	assert.Equal(model2d.NewVertex2D(9, 6), c.GetAttachedVertices()[4])

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.GetAttachedVertices(), 8)
	assert.Len(c.GetAttachedEdges(), 8)
	assert.Len(c.GetUnattachedVertices(), 0)

	assert.True(vertices[0].EdgeTo(vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(vertices[7].EdgeTo(vertices[2]).Equals(c.GetAttachedEdges()[1]))
	assert.True(vertices[2].EdgeTo(vertices[3]).Equals(c.GetAttachedEdges()[2]))
	assert.True(vertices[3].EdgeTo(vertices[5]).Equals(c.GetAttachedEdges()[3]))
	assert.True(vertices[5].EdgeTo(vertices[6]).Equals(c.GetAttachedEdges()[4]))
	assert.True(vertices[6].EdgeTo(vertices[4]).Equals(c.GetAttachedEdges()[5]))
	assert.True(vertices[4].EdgeTo(vertices[1]).Equals(c.GetAttachedEdges()[6]))
	assert.True(vertices[1].EdgeTo(vertices[0]).Equals(c.GetAttachedEdges()[7]))

	assert.Equal(model2d.NewVertex2D(0, 0), c.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(3, 0), c.GetAttachedVertices()[3])
	assert.Equal(model2d.NewVertex2D(8, 5), c.GetAttachedVertices()[4])
	assert.Equal(model2d.NewVertex2D(9, 6), c.GetAttachedVertices()[5])

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.GetAttachedVertices(), 8)
	assert.Len(c.GetAttachedEdges(), 8)
	assert.Len(c.GetUnattachedVertices(), 0)
}

func TestUpdate_DisparityGreedyRelative(t *testing.T) {
	assert := assert.New(t)
	vertices := model2d.DeduplicateVertices([]model.CircuitVertex{
		// Note: the circuit is sorted by DeduplicateVertices(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	})
	c := circuit.NewDisparityGreedy(vertices, model2d.BuildPerimiter, true)

	assert.Len(c.GetAttachedVertices(), 5)
	assert.Len(c.GetAttachedEdges(), 5)
	assert.Len(c.GetUnattachedVertices(), 3)

	// Init circuit := (-15, -15), (15, -15), (9,6), (3,13), (-7,6)

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.GetAttachedVertices(), 6)
	assert.Len(c.GetAttachedEdges(), 6)
	assert.Len(c.GetUnattachedVertices(), 2)

	assert.True(vertices[0].EdgeTo(vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(vertices[7].EdgeTo(vertices[3]).Equals(c.GetAttachedEdges()[1]))
	assert.True(vertices[3].EdgeTo(vertices[6]).Equals(c.GetAttachedEdges()[2]))
	assert.True(vertices[6].EdgeTo(vertices[4]).Equals(c.GetAttachedEdges()[3]))
	assert.True(vertices[4].EdgeTo(vertices[1]).Equals(c.GetAttachedEdges()[4]))
	assert.True(vertices[1].EdgeTo(vertices[0]).Equals(c.GetAttachedEdges()[5]))

	assert.Equal(model2d.NewVertex2D(3, 0), c.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(9, 6), c.GetAttachedVertices()[3])

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.GetAttachedVertices(), 7)
	assert.Len(c.GetAttachedEdges(), 7)
	assert.Len(c.GetUnattachedVertices(), 1)

	assert.True(vertices[0].EdgeTo(vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(vertices[7].EdgeTo(vertices[3]).Equals(c.GetAttachedEdges()[1]))
	assert.True(vertices[3].EdgeTo(vertices[5]).Equals(c.GetAttachedEdges()[2]))
	assert.True(vertices[5].EdgeTo(vertices[6]).Equals(c.GetAttachedEdges()[3]))
	assert.True(vertices[6].EdgeTo(vertices[4]).Equals(c.GetAttachedEdges()[4]))
	assert.True(vertices[4].EdgeTo(vertices[1]).Equals(c.GetAttachedEdges()[5]))
	assert.True(vertices[1].EdgeTo(vertices[0]).Equals(c.GetAttachedEdges()[6]))

	assert.Equal(model2d.NewVertex2D(3, 0), c.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(8, 5), c.GetAttachedVertices()[3])
	assert.Equal(model2d.NewVertex2D(9, 6), c.GetAttachedVertices()[4])

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.GetAttachedVertices(), 8)
	assert.Len(c.GetAttachedEdges(), 8)
	assert.Len(c.GetUnattachedVertices(), 0)

	assert.True(vertices[0].EdgeTo(vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(vertices[7].EdgeTo(vertices[2]).Equals(c.GetAttachedEdges()[1]))
	assert.True(vertices[2].EdgeTo(vertices[3]).Equals(c.GetAttachedEdges()[2]))
	assert.True(vertices[3].EdgeTo(vertices[5]).Equals(c.GetAttachedEdges()[3]))
	assert.True(vertices[5].EdgeTo(vertices[6]).Equals(c.GetAttachedEdges()[4]))
	assert.True(vertices[6].EdgeTo(vertices[4]).Equals(c.GetAttachedEdges()[5]))
	assert.True(vertices[4].EdgeTo(vertices[1]).Equals(c.GetAttachedEdges()[6]))
	assert.True(vertices[1].EdgeTo(vertices[0]).Equals(c.GetAttachedEdges()[7]))

	assert.Equal(model2d.NewVertex2D(0, 0), c.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(3, 0), c.GetAttachedVertices()[3])
	assert.Equal(model2d.NewVertex2D(8, 5), c.GetAttachedVertices()[4])
	assert.Equal(model2d.NewVertex2D(9, 6), c.GetAttachedVertices()[5])

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.GetAttachedVertices(), 8)
	assert.Len(c.GetAttachedEdges(), 8)
	assert.Len(c.GetUnattachedVertices(), 0)
}

func TestUpdate_ShouldNotRemoveAttachedInteriorPointFromPerimeterIfNewEdgeIsCloserThanPreviousEdge_DisparityGreedy(t *testing.T) {
	assert := assert.New(t)
	vertices := model2d.DeduplicateVertices([]model.CircuitVertex{
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(4.7, 2.0),
		model2d.NewVertex2D(5.0, 2.25),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6.0, 2.5),
		model2d.NewVertex2D(10, 0),
	})
	c := circuit.NewDisparityGreedy(vertices, model2d.BuildPerimiter, false)

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 4)
	assert.Equal(vertices[0], c.GetAttachedVertices()[0])
	assert.Equal(vertices[5], c.GetAttachedVertices()[1])
	assert.Equal(vertices[4], c.GetAttachedVertices()[2])
	assert.Equal(vertices[3], c.GetAttachedVertices()[3])

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 5)
	assert.Equal(vertices[0], c.GetAttachedVertices()[0])
	assert.Equal(vertices[1], c.GetAttachedVertices()[1])
	assert.Equal(vertices[5], c.GetAttachedVertices()[2])
	assert.Equal(vertices[4], c.GetAttachedVertices()[3])
	assert.Equal(vertices[3], c.GetAttachedVertices()[4])

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Equal(vertices[0], c.GetAttachedVertices()[0])
	assert.Equal(vertices[1], c.GetAttachedVertices()[1])
	assert.Equal(vertices[2], c.GetAttachedVertices()[2])
	assert.Equal(vertices[5], c.GetAttachedVertices()[3])
	assert.Equal(vertices[4], c.GetAttachedVertices()[4])
	assert.Equal(vertices[3], c.GetAttachedVertices()[5])

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Equal(vertices[0], c.GetAttachedVertices()[0])
	assert.Equal(vertices[1], c.GetAttachedVertices()[1])
	assert.Equal(vertices[2], c.GetAttachedVertices()[2])
	assert.Equal(vertices[5], c.GetAttachedVertices()[3])
	assert.Equal(vertices[4], c.GetAttachedVertices()[4])
	assert.Equal(vertices[3], c.GetAttachedVertices()[5])

	v, _ := c.FindNextVertexAndEdge()
	assert.Nil(v)
}

func TestUpdate_DisparityGreedy_ShouldPanicIfEdgeIsNotInCircuit(t *testing.T) {
	assert := assert.New(t)
	vertices := model2d.DeduplicateVertices([]model.CircuitVertex{
		// Note: the circuit is sorted by DeduplicateVertices(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	})
	c := circuit.NewDisparityGreedy(vertices, model2d.BuildPerimiter, false)

	assert.Len(c.GetAttachedVertices(), 5)
	assert.Len(c.GetAttachedEdges(), 5)
	assert.Len(c.GetUnattachedVertices(), 3)

	assert.Panics(func() { c.Update(vertices[2], vertices[0].EdgeTo(vertices[4])) })
}
