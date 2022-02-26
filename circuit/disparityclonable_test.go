package circuit_test

import (
	"testing"

	"github.com/heustis/tsp-solver-go/circuit"
	"github.com/heustis/tsp-solver-go/model"
	"github.com/heustis/tsp-solver-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestNewDisparityClonable(t *testing.T) {
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
	c := circuit.NewDisparityClonable(vertices, model2d.BuildPerimiter)

	attached := c.GetAttachedVertices()
	assert.Len(attached, 5)
	assert.Equal(model2d.NewVertex2D(-15, -15), attached[0])
	assert.Equal(model2d.NewVertex2D(15, -15), attached[1])
	assert.Equal(model2d.NewVertex2D(9, 6), attached[2])
	assert.Equal(model2d.NewVertex2D(3, 13), attached[3])
	assert.Equal(model2d.NewVertex2D(-7, 6), attached[4])

	assert.InDelta(95.73863479511238, c.GetLength(), model.Threshold)

	edges := c.GetAttachedEdges()
	assert.Len(edges, 5)
	assert.True(vertices[0].EdgeTo(vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(vertices[7].EdgeTo(vertices[6]).Equals(c.GetAttachedEdges()[1]))
	assert.True(vertices[6].EdgeTo(vertices[4]).Equals(c.GetAttachedEdges()[2]))
	assert.True(vertices[4].EdgeTo(vertices[1]).Equals(c.GetAttachedEdges()[3]))
	assert.True(vertices[1].EdgeTo(vertices[0]).Equals(c.GetAttachedEdges()[4]))

	unattached := c.GetUnattachedVertices()
	assert.Len(unattached, 3)
	assert.True(unattached[vertices[2]])
	assert.True(unattached[vertices[3]])
	assert.True(unattached[vertices[5]])
}

func TestUpdate_DisparityClonable(t *testing.T) {
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
	c := circuit.NewDisparityClonable(vertices, model2d.BuildPerimiter)

	assert.Len(c.GetAttachedVertices(), 5)
	assert.Len(c.GetAttachedEdges(), 5)
	assert.Len(c.GetUnattachedVertices(), 3)

	c.Update(c.FindNextVertexAndEdge())
	edges := c.GetAttachedEdges()
	attached := c.GetAttachedVertices()

	assert.Len(attached, 6)
	assert.Len(edges, 6)
	assert.Len(c.GetUnattachedVertices(), 2)

	assert.True(vertices[0].EdgeTo(vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(vertices[7].EdgeTo(vertices[3]).Equals(c.GetAttachedEdges()[1]))
	assert.True(vertices[3].EdgeTo(vertices[6]).Equals(c.GetAttachedEdges()[2]))
	assert.True(vertices[6].EdgeTo(vertices[4]).Equals(c.GetAttachedEdges()[3]))
	assert.True(vertices[4].EdgeTo(vertices[1]).Equals(c.GetAttachedEdges()[4]))
	assert.True(vertices[1].EdgeTo(vertices[0]).Equals(c.GetAttachedEdges()[5]))

	assert.Equal(model2d.NewVertex2D(15, -15), attached[1])
	assert.Equal(model2d.NewVertex2D(3, 0), attached[2])
	assert.Equal(model2d.NewVertex2D(9, 6), attached[3])

	c.Update(c.FindNextVertexAndEdge())
	edges = c.GetAttachedEdges()
	attached = c.GetAttachedVertices()

	assert.Len(attached, 7)
	assert.Len(edges, 7)
	assert.Len(c.GetUnattachedVertices(), 1)

	assert.True(vertices[0].EdgeTo(vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(vertices[7].EdgeTo(vertices[3]).Equals(c.GetAttachedEdges()[1]))
	assert.True(vertices[3].EdgeTo(vertices[5]).Equals(c.GetAttachedEdges()[2]))
	assert.True(vertices[5].EdgeTo(vertices[6]).Equals(c.GetAttachedEdges()[3]))
	assert.True(vertices[6].EdgeTo(vertices[4]).Equals(c.GetAttachedEdges()[4]))
	assert.True(vertices[4].EdgeTo(vertices[1]).Equals(c.GetAttachedEdges()[5]))
	assert.True(vertices[1].EdgeTo(vertices[0]).Equals(c.GetAttachedEdges()[6]))

	assert.Equal(model2d.NewVertex2D(3, 0), attached[2])
	assert.Equal(model2d.NewVertex2D(8, 5), attached[3])
	assert.Equal(model2d.NewVertex2D(9, 6), attached[4])

	c.Update(c.FindNextVertexAndEdge())
	edges = c.GetAttachedEdges()
	attached = c.GetAttachedVertices()

	assert.Len(attached, 8)
	assert.Len(edges, 8)
	assert.Len(c.GetUnattachedVertices(), 0)

	assert.True(vertices[0].EdgeTo(vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(vertices[7].EdgeTo(vertices[2]).Equals(c.GetAttachedEdges()[1]))
	assert.True(vertices[2].EdgeTo(vertices[3]).Equals(c.GetAttachedEdges()[2]))
	assert.True(vertices[3].EdgeTo(vertices[5]).Equals(c.GetAttachedEdges()[3]))
	assert.True(vertices[5].EdgeTo(vertices[6]).Equals(c.GetAttachedEdges()[4]))
	assert.True(vertices[6].EdgeTo(vertices[4]).Equals(c.GetAttachedEdges()[5]))
	assert.True(vertices[4].EdgeTo(vertices[1]).Equals(c.GetAttachedEdges()[6]))
	assert.True(vertices[1].EdgeTo(vertices[0]).Equals(c.GetAttachedEdges()[7]))

	assert.Equal(model2d.NewVertex2D(0, 0), attached[2])
	assert.Equal(model2d.NewVertex2D(3, 0), attached[3])

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.GetAttachedVertices(), 8)
	assert.Len(c.GetAttachedEdges(), 8)
	assert.Len(c.GetUnattachedVertices(), 0)
}

func TestUpdate_ShouldNotRemoveAttachedInteriorPointFromPerimeterIfNewEdgeIsCloserThanPreviousEdge_DisparityClonable(t *testing.T) {
	assert := assert.New(t)
	vertices := model2d.DeduplicateVertices([]model.CircuitVertex{
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(4.7, 2.0),
		model2d.NewVertex2D(5.0, 2.25),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6.0, 2.5),
		model2d.NewVertex2D(10, 0),
	})
	c := circuit.NewDisparityClonable(vertices, model2d.BuildPerimiter)

	c.Update(c.FindNextVertexAndEdge())
	attached := c.GetAttachedVertices()
	assert.Len(attached, 4)
	assert.Equal(vertices[0], attached[0])
	assert.Equal(vertices[2], attached[1])
	assert.Equal(vertices[5], attached[2])
	assert.Equal(vertices[3], attached[3])

	c.Update(c.FindNextVertexAndEdge())
	attached = c.GetAttachedVertices()
	assert.Len(attached, 5)
	assert.Equal(vertices[0], attached[0])
	assert.Equal(vertices[2], attached[1])
	assert.Equal(vertices[4], attached[2])
	assert.Equal(vertices[5], attached[3])
	assert.Equal(vertices[3], attached[4])

	c.Update(c.FindNextVertexAndEdge())
	attached = c.GetAttachedVertices()
	assert.Len(attached, 6)
	assert.Equal(vertices[0], attached[0])
	assert.Equal(vertices[1], attached[1])
	assert.Equal(vertices[2], attached[2])
	assert.Equal(vertices[4], attached[3])
	assert.Equal(vertices[5], attached[4])
	assert.Equal(vertices[3], attached[5])

	c.Update(c.FindNextVertexAndEdge())
	attached = c.GetAttachedVertices()
	assert.Len(attached, 6)
	assert.Equal(vertices[0], attached[0])
	assert.Equal(vertices[1], attached[1])
	assert.Equal(vertices[2], attached[2])
	assert.Equal(vertices[4], attached[3])
	assert.Equal(vertices[5], attached[4])
	assert.Equal(vertices[3], attached[5])

	v, _ := c.FindNextVertexAndEdge()
	assert.Nil(v)
}
