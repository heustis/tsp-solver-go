package circuit_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerimeter_Heap(t *testing.T) {
	assert := assert.New(t)
	c := circuit.NewClonableCircuitImpl(model2d.DeduplicateVertices([]model.CircuitVertex{
		// Note: the circuit is sorted by DeduplicateVertices(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}), model2d.BuildPerimiter)
	c.CloneOnFirstAttach = true

	assert.Len(c.Vertices, 8)
	assert.Len(c.GetAttachedVertices(), 5)
	assert.Equal(model2d.NewVertex2D(-15, -15), c.GetAttachedVertices()[0])
	assert.Equal(model2d.NewVertex2D(15, -15), c.GetAttachedVertices()[1])
	assert.Equal(model2d.NewVertex2D(9, 6), c.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(3, 13), c.GetAttachedVertices()[3])
	assert.Equal(model2d.NewVertex2D(-7, 6), c.GetAttachedVertices()[4])
	assert.Equal(c.GetAttachedVertices(), c.GetAttachedVertices())

	assert.Len(c.GetAttachedEdges(), 5)
	assert.True(c.Vertices[0].EdgeTo(c.Vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(c.Vertices[7].EdgeTo(c.Vertices[6]).Equals(c.GetAttachedEdges()[1]))
	assert.True(c.Vertices[6].EdgeTo(c.Vertices[4]).Equals(c.GetAttachedEdges()[2]))
	assert.True(c.Vertices[4].EdgeTo(c.Vertices[1]).Equals(c.GetAttachedEdges()[3]))
	assert.True(c.Vertices[1].EdgeTo(c.Vertices[0]).Equals(c.GetAttachedEdges()[4]))

	l := 0.0
	for _, edge := range c.GetAttachedEdges() {
		l += edge.GetLength()
	}
	assert.InDelta(l, c.GetLength(), model.Threshold)

	assert.Len(c.GetUnattachedVertices(), 3)
	assert.True(c.GetUnattachedVertices()[c.Vertices[2]])
	assert.True(c.GetUnattachedVertices()[c.Vertices[3]])
	assert.True(c.GetUnattachedVertices()[c.Vertices[5]])
	assert.Equal(c.GetUnattachedVertices(), c.GetUnattachedVertices())

	assert.InDelta(c.GetLength()+0.763503994948632, c.GetLengthWithNext(), model.Threshold)

	assert.Equal(15, c.GetClosestEdges().Len())
	assert.Equal(&model.DistanceToEdge{
		Edge:     c.GetAttachedEdges()[1],
		Vertex:   c.Vertices[5],
		Distance: 0.763503994948632,
	}, c.GetClosestEdges().PopHeap())
	assert.Equal(&model.DistanceToEdge{
		Edge:     c.GetAttachedEdges()[2],
		Vertex:   c.Vertices[5],
		Distance: 1.628650237136812,
	}, c.GetClosestEdges().PopHeap())
	assert.Equal(&model.DistanceToEdge{
		Edge:     c.GetAttachedEdges()[1],
		Vertex:   c.Vertices[3],
		Distance: 5.854324418695558,
	}, c.GetClosestEdges().PopHeap())
	assert.Equal(&model.DistanceToEdge{
		Edge:     c.GetAttachedEdges()[4],
		Vertex:   c.Vertices[2],
		Distance: 7.9605428386450825,
	}, c.GetClosestEdges().PopHeap())
}

func TestCloneAndUpdate(t *testing.T) {
	assert := assert.New(t)

	c := circuit.NewClonableCircuitImpl(model2d.DeduplicateVertices([]model.CircuitVertex{
		// Note: the circuit is sorted by DeduplicateVertices(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}), model2d.BuildPerimiter)
	c.CloneOnFirstAttach = true

	clone := c.CloneAndUpdate().(*circuit.ClonableCircuitImpl)

	assert.Len(c.GetUnattachedVertices(), 3)
	assert.Len(c.GetAttachedVertices(), 5)
	assert.Len(c.GetAttachedEdges(), 5)
	assert.Equal(14, c.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Edge:     c.GetAttachedEdges()[2],
		Vertex:   c.Vertices[5],
		Distance: 1.628650237136812,
	}, c.GetClosestEdges().Peek())

	assert.Len(clone.GetUnattachedVertices(), 2)
	assert.Len(clone.GetAttachedVertices(), 6)
	assert.Len(clone.GetAttachedEdges(), 6)
	assert.Equal(12, clone.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Edge:     clone.GetAttachedEdges()[1],
		Vertex:   clone.Vertices[3],
		Distance: 5.09082042374693,
	}, clone.GetClosestEdges().Peek())

	assert.InDelta(c.GetLength()+0.763503994948632, clone.GetLength(), model.Threshold)

	// Validate that cloning a clone does not affect the original c
	cloneOfClone := clone.CloneAndUpdate().(*circuit.ClonableCircuitImpl)

	assert.Len(c.GetUnattachedVertices(), 3)
	assert.Len(c.GetAttachedVertices(), 5)
	assert.Len(c.GetAttachedEdges(), 5)
	assert.Equal(14, c.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Edge:     c.GetAttachedEdges()[2],
		Vertex:   c.Vertices[5],
		Distance: 1.628650237136812,
	}, c.GetClosestEdges().Peek())

	assert.Len(clone.GetUnattachedVertices(), 2)
	assert.Len(clone.GetAttachedVertices(), 6)
	assert.Len(clone.GetAttachedEdges(), 6)
	assert.Equal(11, clone.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Edge:     clone.GetAttachedEdges()[5],
		Vertex:   clone.Vertices[2],
		Distance: 7.9605428386450825,
	}, clone.GetClosestEdges().Peek())

	assert.Len(cloneOfClone.GetUnattachedVertices(), 1)
	assert.Len(cloneOfClone.GetAttachedVertices(), 7)
	assert.Len(cloneOfClone.GetAttachedEdges(), 7)
	assert.Equal(7, cloneOfClone.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Edge:     cloneOfClone.GetAttachedEdges()[1],
		Vertex:   cloneOfClone.Vertices[2],
		Distance: 5.003830723297881,
	}, cloneOfClone.GetClosestEdges().Peek())

	// Validate that cloning a c with only one vertex left to attach just updates that object and doesn't create a clone
	assert.Nil(cloneOfClone.CloneAndUpdate())

	assert.Len(cloneOfClone.GetUnattachedVertices(), 0)
	assert.Len(cloneOfClone.GetAttachedVertices(), 8)
	assert.Len(cloneOfClone.GetAttachedEdges(), 8)
	assert.Equal(0, cloneOfClone.GetClosestEdges().Len())

	// Validate that cloning a completed c makes no changes
	assert.Nil(cloneOfClone.CloneAndUpdate())

	assert.Len(cloneOfClone.GetUnattachedVertices(), 0)
	assert.Len(cloneOfClone.GetAttachedVertices(), 8)
	assert.Len(cloneOfClone.GetAttachedEdges(), 8)
	assert.Equal(0, cloneOfClone.GetClosestEdges().Len())
}

func TestDelete_Heap(t *testing.T) {
	assert := assert.New(t)

	c := circuit.NewClonableCircuitImpl(model2d.DeduplicateVertices([]model.CircuitVertex{
		// Note: the circuit is sorted by DeduplicateVertices(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}), model2d.BuildPerimiter)
	c.CloneOnFirstAttach = true

	clone := c.CloneAndUpdate().(*circuit.ClonableCircuitImpl)
	cloneOfClone := clone.CloneAndUpdate().(*circuit.ClonableCircuitImpl)

	c.Delete()
	assert.Len(c.GetUnattachedVertices(), 0)
	assert.Len(c.GetAttachedVertices(), 0)
	assert.Nil(c.GetAttachedEdges())
	assert.Nil(c.GetClosestEdges())
	assert.Nil(c.Vertices)

	assert.Len(clone.GetUnattachedVertices(), 2)
	assert.Len(clone.GetAttachedEdges(), 6)
	assert.Len(clone.GetAttachedVertices(), 6)
	assert.NotNil(clone.GetClosestEdges())
	assert.Equal(11, clone.GetClosestEdges().Len())
	assert.Len(clone.Vertices, 8)

	assert.Len(cloneOfClone.GetUnattachedVertices(), 1)
	assert.Len(cloneOfClone.GetAttachedEdges(), 7)
	assert.Len(cloneOfClone.GetAttachedVertices(), 7)
	assert.NotNil(cloneOfClone.GetClosestEdges())
	assert.Equal(7, cloneOfClone.GetClosestEdges().Len())
	assert.Len(cloneOfClone.Vertices, 8)

	cloneOfClone.Delete()
	assert.Len(cloneOfClone.GetUnattachedVertices(), 0)
	assert.Len(cloneOfClone.GetAttachedVertices(), 0)
	assert.Nil(cloneOfClone.GetAttachedEdges())
	assert.Nil(cloneOfClone.GetClosestEdges())
	assert.Nil(cloneOfClone.Vertices)

	assert.Len(clone.GetUnattachedVertices(), 2)
	assert.Len(clone.GetAttachedEdges(), 6)
	assert.Len(clone.GetAttachedVertices(), 6)
	assert.NotNil(clone.GetClosestEdges())
	assert.Equal(11, clone.GetClosestEdges().Len())
	assert.Len(clone.Vertices, 8)
	clone.Delete()
}
