package experimental_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/circuit/experimental"
	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerimeter_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)
	c := experimental.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}, model2d.DeduplicateVertices, model2d.BuildPerimiter)
	c.Prepare()
	c.BuildPerimiter()

	assert.Len(c.Vertices, 8)

	assert.Len(c.GetAttachedVertices(), 5)
	assert.Equal(model2d.NewVertex2D(-15, -15), c.GetAttachedVertices()[0])
	assert.Equal(model2d.NewVertex2D(15, -15), c.GetAttachedVertices()[1])
	assert.Equal(model2d.NewVertex2D(9, 6), c.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(3, 13), c.GetAttachedVertices()[3])
	assert.Equal(model2d.NewVertex2D(-7, 6), c.GetAttachedVertices()[4])
	assert.Equal(c.GetAttachedVertices(), c.GetAttachedVertices())

	assert.Len(c.GetAttachedEdges(), 5)
	assert.Equal(c.Vertices[0].EdgeTo(c.Vertices[7]), c.GetAttachedEdges()[0])
	assert.Equal(c.Vertices[7].EdgeTo(c.Vertices[6]), c.GetAttachedEdges()[1])
	assert.Equal(c.Vertices[6].EdgeTo(c.Vertices[4]), c.GetAttachedEdges()[2])
	assert.Equal(c.Vertices[4].EdgeTo(c.Vertices[1]), c.GetAttachedEdges()[3])
	assert.Equal(c.Vertices[1].EdgeTo(c.Vertices[0]), c.GetAttachedEdges()[4])

	expectedLength := 0.0
	for _, edge := range c.GetAttachedEdges() {
		expectedLength += edge.GetLength()
	}
	assert.InDelta(expectedLength, c.GetLength(), model.Threshold)

	assert.Len(c.GetUnattachedVertices(), 3)
	assert.True(c.GetUnattachedVertices()[c.Vertices[2]])
	assert.True(c.GetUnattachedVertices()[c.Vertices[3]])
	assert.True(c.GetUnattachedVertices()[c.Vertices[5]])

	assert.InDelta(95.738634795112368+0.763503994948632, c.GetLengthWithNext(), model.Threshold)

	assert.Equal(9, c.GetClosestEdges().Len())
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

func TestAttachAndMove_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)

	c := experimental.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(0, 0),   // Index 0 after sorting
		model2d.NewVertex2D(0, 3),   // Index 1 after sorting
		model2d.NewVertex2D(3, 3),   // Index 5 after sorting
		model2d.NewVertex2D(3, 0),   // Index 4 after sorting
		model2d.NewVertex2D(1, 1),   // Index 2 after sorting
		model2d.NewVertex2D(1, 2.1), // Index 3 after sorting
	}, model2d.DeduplicateVertices, model2d.BuildPerimiter)
	c.Prepare()
	c.BuildPerimiter()

	testVert := c.GetClosestEdges().PopHeap().(*model.DistanceToEdge)
	c.AttachVertex(testVert)

	assert.InDelta(12.5385336246535019133711157158298, c.GetLength(), model.Threshold)
	assert.Equal(5, c.GetClosestEdges().Len())

	c.MoveVertex(&model.DistanceToEdge{
		Vertex:   testVert.Vertex,
		Edge:     c.GetAttachedEdges()[1],
		Distance: -1.0,
	})

	assert.Len(c.GetUnattachedVertices(), 1)
	assert.True(c.GetUnattachedVertices()[c.Vertices[2]])

	assert.Equal(11.5385336246535019133711157158298, c.GetLength())
	assert.Equal(3, c.GetClosestEdges().Len())
}

func TestAttachAndMoveIndexZero_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)

	c := experimental.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(0, 0),   // Index 0 after sorting
		model2d.NewVertex2D(0, 3),   // Index 1 after sorting
		model2d.NewVertex2D(3, 3),   // Index 5 after sorting
		model2d.NewVertex2D(3, 0),   // Index 4 after sorting
		model2d.NewVertex2D(1, 1),   // Index 2 after sorting
		model2d.NewVertex2D(1, 2.1), // Index 3 after sorting
	}, model2d.DeduplicateVertices, model2d.BuildPerimiter)
	c.Prepare()
	c.BuildPerimiter()

	testVert := &model.DistanceToEdge{
		Vertex:   c.Vertices[2],
		Edge:     c.GetAttachedEdges()[0],
		Distance: 1.0,
	}
	c.AttachVertex(testVert)
	assert.Equal([]model.CircuitEdge{
		c.Vertices[4].EdgeTo(c.Vertices[2]),
		c.Vertices[2].EdgeTo(c.Vertices[5]),
		c.Vertices[5].EdgeTo(c.Vertices[1]),
		c.Vertices[1].EdgeTo(c.Vertices[0]),
		c.Vertices[0].EdgeTo(c.Vertices[4]),
	}, c.GetAttachedEdges())
	assert.Equal([]model.CircuitVertex{
		c.Vertices[4],
		c.Vertices[2],
		c.Vertices[5],
		c.Vertices[1],
		c.Vertices[0],
	}, c.GetAttachedVertices())

	c.MoveVertex(&model.DistanceToEdge{
		Vertex:   c.Vertices[2],
		Edge:     c.GetAttachedEdges()[4],
		Distance: 2.25,
	})
	assert.Equal([]model.CircuitEdge{
		c.Vertices[4].EdgeTo(c.Vertices[5]),
		c.Vertices[5].EdgeTo(c.Vertices[1]),
		c.Vertices[1].EdgeTo(c.Vertices[0]),
		c.Vertices[0].EdgeTo(c.Vertices[2]),
		c.Vertices[2].EdgeTo(c.Vertices[4]),
	}, c.GetAttachedEdges())
	assert.Equal([]model.CircuitVertex{
		c.Vertices[4],
		c.Vertices[5],
		c.Vertices[1],
		c.Vertices[0],
		c.Vertices[2],
	}, c.GetAttachedVertices())

	c.MoveVertex(&model.DistanceToEdge{
		Vertex:   c.Vertices[2],
		Edge:     c.GetAttachedEdges()[0],
		Distance: 3.5,
	})
	assert.Equal([]model.CircuitEdge{
		c.Vertices[4].EdgeTo(c.Vertices[2]),
		c.Vertices[2].EdgeTo(c.Vertices[5]),
		c.Vertices[5].EdgeTo(c.Vertices[1]),
		c.Vertices[1].EdgeTo(c.Vertices[0]),
		c.Vertices[0].EdgeTo(c.Vertices[4]),
	}, c.GetAttachedEdges())
	assert.Equal([]model.CircuitVertex{
		c.Vertices[4],
		c.Vertices[2],
		c.Vertices[5],
		c.Vertices[1],
		c.Vertices[0],
	}, c.GetAttachedVertices())

	assert.Len(c.GetUnattachedVertices(), 1)
	assert.True(c.GetUnattachedVertices()[c.Vertices[3]])

	assert.Equal(12.0+6.75, c.GetLength()) // 6.75 comes from the distances in each DistanceToEdge
	assert.Equal(4, c.GetClosestEdges().Len())
}

func TestAttachAndMoveLastIndex_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)

	c := experimental.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(0, 0),   // Index 0 after sorting
		model2d.NewVertex2D(0, 3),   // Index 1 after sorting
		model2d.NewVertex2D(3, 3),   // Index 5 after sorting
		model2d.NewVertex2D(3, 0),   // Index 4 after sorting
		model2d.NewVertex2D(1, 1),   // Index 2 after sorting
		model2d.NewVertex2D(1, 2.1), // Index 3 after sorting
	}, model2d.DeduplicateVertices, model2d.BuildPerimiter)
	c.Prepare()
	c.BuildPerimiter()

	testVert := &model.DistanceToEdge{
		Vertex:   c.Vertices[2],
		Edge:     c.GetAttachedEdges()[3],
		Distance: 0.5,
	}
	c.AttachVertex(testVert)
	assert.Equal([]model.CircuitEdge{
		c.Vertices[4].EdgeTo(c.Vertices[5]),
		c.Vertices[5].EdgeTo(c.Vertices[1]),
		c.Vertices[1].EdgeTo(c.Vertices[0]),
		c.Vertices[0].EdgeTo(c.Vertices[2]),
		c.Vertices[2].EdgeTo(c.Vertices[4]),
	}, c.GetAttachedEdges())
	assert.Equal([]model.CircuitVertex{
		c.Vertices[4],
		c.Vertices[5],
		c.Vertices[1],
		c.Vertices[0],
		c.Vertices[2],
	}, c.GetAttachedVertices())

	c.MoveVertex(&model.DistanceToEdge{
		Vertex:   c.Vertices[2],
		Edge:     c.GetAttachedEdges()[2],
		Distance: 0.25,
	})
	assert.Equal([]model.CircuitEdge{
		c.Vertices[4].EdgeTo(c.Vertices[5]),
		c.Vertices[5].EdgeTo(c.Vertices[1]),
		c.Vertices[1].EdgeTo(c.Vertices[2]),
		c.Vertices[2].EdgeTo(c.Vertices[0]),
		c.Vertices[0].EdgeTo(c.Vertices[4]),
	}, c.GetAttachedEdges())
	assert.Equal([]model.CircuitVertex{
		c.Vertices[4],
		c.Vertices[5],
		c.Vertices[1],
		c.Vertices[2],
		c.Vertices[0],
	}, c.GetAttachedVertices())

	assert.Len(c.GetUnattachedVertices(), 1)
	assert.True(c.GetUnattachedVertices()[c.Vertices[3]])

	assert.Equal(12.0+0.75, c.GetLength())
	assert.Equal(4, c.GetClosestEdges().Len())
}

func TestAttachShouldPanicIfEdgeIsNotInCircuit_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)

	c := experimental.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}, model2d.DeduplicateVertices, model2d.BuildPerimiter)
	c.Prepare()
	c.BuildPerimiter()

	assert.Panics(func() {
		c.AttachVertex(&model.DistanceToEdge{
			Vertex:   c.Vertices[0],
			Edge:     model2d.NewVertex2D(3, 3).EdgeTo(model2d.NewVertex2D(5, 5)),
			Distance: 1.234,
		})
	})
}

func TestMoveShouldPanicIfEdgeIsNotInCircuit_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)

	c := experimental.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}, model2d.DeduplicateVertices, model2d.BuildPerimiter)
	c.Prepare()
	c.BuildPerimiter()

	assert.Panics(func() {
		c.MoveVertex(&model.DistanceToEdge{
			Vertex:   c.Vertices[0],
			Edge:     model2d.NewVertex2D(3, 3).EdgeTo(model2d.NewVertex2D(5, 5)),
			Distance: 1.234,
		})
	})
}

func TestCloneAndUpdate_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)

	c := experimental.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}, model2d.DeduplicateVertices, model2d.BuildPerimiter)
	c.Prepare()
	c.BuildPerimiter()

	assert.True(c.GetUnattachedVertices()[c.Vertices[2]])
	assert.True(c.GetUnattachedVertices()[c.Vertices[3]])
	assert.True(c.GetUnattachedVertices()[c.Vertices[5]])

	// Index 5 should attach to edge 15,-15 -> 9,6
	assert.Nil(c.CloneAndUpdate())
	assert.Len(c.GetUnattachedVertices(), 2)
	assert.False(c.GetUnattachedVertices()[c.Vertices[5]])
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Len(c.GetAttachedEdges(), 6)
	assert.Equal(8, c.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Edge:     c.GetAttachedEdges()[3],
		Vertex:   c.Vertices[5],
		Distance: 0.8651462421881799,
	}, c.GetClosestEdges().Peek())

	// Index 5 should attach to edge 9,6 -> 3,13, this requires cloning since index 5 is already attached.
	clone := c.CloneAndUpdate().(*experimental.HeapableCircuitMinClonesLimited)
	assert.Len(c.GetUnattachedVertices(), 2)
	assert.Len(clone.GetUnattachedVertices(), 2)
	assert.False(c.GetUnattachedVertices()[c.Vertices[5]])
	assert.False(clone.GetUnattachedVertices()[c.Vertices[5]])
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Len(clone.GetAttachedVertices(), 6)
	assert.Len(c.GetAttachedEdges(), 6)
	assert.Len(clone.GetAttachedEdges(), 6)
	assert.Equal(7, c.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Edge:     c.GetAttachedEdges()[1],
		Vertex:   c.Vertices[3],
		Distance: 5.09082042374693,
	}, c.GetClosestEdges().Peek())

	assert.Equal(&model.DistanceToEdge{
		Edge:     clone.GetAttachedEdges()[1],
		Vertex:   clone.Vertices[3],
		Distance: 5.854324418695558,
	}, clone.GetClosestEdges().Peek())

	// Index 3 should attach to edge 1, no cloning required
	assert.Nil(c.CloneAndUpdate())
	assert.Nil(clone.CloneAndUpdate())
	assert.Len(c.GetUnattachedVertices(), 1)
	assert.Equal(6, c.GetClosestEdges().Len())
	assert.Len(clone.GetUnattachedVertices(), 1)
	assert.Equal(6, clone.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Edge:     c.GetAttachedEdges()[1],
		Vertex:   c.Vertices[2],
		Distance: 5.003830723297881,
	}, c.GetClosestEdges().Peek())

	assert.Equal(&model.DistanceToEdge{
		Edge:     clone.GetAttachedEdges()[4],
		Vertex:   clone.Vertices[3],
		Distance: 4.782762261113314,
	}, clone.GetClosestEdges().Peek())

	// Index 2 should attach to edge 1, no cloning required
	assert.Nil(c.CloneAndUpdate())
	// Index 3 should move to edge 4, cloning required
	cloneOfClone, okay := clone.CloneAndUpdate().(*experimental.HeapableCircuitMinClonesLimited)
	assert.True(okay)

	assert.Len(c.GetUnattachedVertices(), 0)
	assert.Len(clone.GetUnattachedVertices(), 1)
	assert.Len(cloneOfClone.GetUnattachedVertices(), 1)

	// Index 2 should attach to edge 1, no cloning required
	assert.Equal(&model.DistanceToEdge{
		Edge:     clone.GetAttachedEdges()[1],
		Vertex:   clone.Vertices[2],
		Distance: 5.003830723297881,
	}, clone.GetClosestEdges().Peek())
	assert.Nil(clone.CloneAndUpdate())

	// Index 2 should move to edge 7, cloning required
	assert.Equal(&model.DistanceToEdge{
		Edge:     clone.GetAttachedEdges()[7],
		Vertex:   clone.Vertices[2],
		Distance: 2.9567121153472016,
	}, clone.GetClosestEdges().Peek())
	assert.NotNil(clone.CloneAndUpdate())

	// Index 2 should move to edge 6, cloning required
	assert.Equal(&model.DistanceToEdge{
		Edge:     cloneOfClone.GetAttachedEdges()[6],
		Vertex:   cloneOfClone.Vertices[2],
		Distance: 7.9605428386450825,
	}, cloneOfClone.GetClosestEdges().Peek())
	assert.Nil(cloneOfClone.CloneAndUpdate())
}

func TestCloneAndUpdate_HeapMinClonesLimited_Distances(t *testing.T) {
	assert := assert.New(t)

	c := experimental.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(0, 0),   // Index 0 after sorting
		model2d.NewVertex2D(0, 3),   // Index 1 after sorting
		model2d.NewVertex2D(3, 3),   // Index 5 after sorting
		model2d.NewVertex2D(3, 0),   // Index 4 after sorting
		model2d.NewVertex2D(1, 1),   // Index 2 after sorting
		model2d.NewVertex2D(1, 2.1), // Index 3 after sorting
	}, model2d.DeduplicateVertices, model2d.BuildPerimiter)
	c.Prepare()
	c.BuildPerimiter()

	assert.Len(c.GetUnattachedVertices(), 2)
	assert.True(c.GetUnattachedVertices()[c.Vertices[2]])
	assert.True(c.GetUnattachedVertices()[c.Vertices[3]])

	assert.Equal(12.0, c.GetLength())
	assert.Equal(6, c.GetClosestEdges().Len())

	// No clone on first attachment - vertex {1,2.1} to edge {0,3}->{3,3}
	assert.Nil(c.CloneAndUpdate())

	assert.Len(c.GetUnattachedVertices(), 1)
	assert.True(c.GetUnattachedVertices()[c.Vertices[2]])

	assert.InDelta(12.5385336246535019133711157158298, c.GetLength(), model.Threshold)
	assert.Equal(5, c.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Distance: 0.1327694499764709,
		Vertex:   c.Vertices[3],
		Edge:     c.GetAttachedEdges()[3],
	}, c.GetClosestEdges().Peek())

	// Clone on second attachment - vertex {1,2.1} to edge {0,0}->{0,3}
	clone := c.CloneAndUpdate().(*experimental.HeapableCircuitMinClonesLimited)

	// Validate that the first c is unchanged.
	assert.Len(c.GetUnattachedVertices(), 1)
	assert.True(c.GetUnattachedVertices()[c.Vertices[2]])

	assert.InDelta(12.5385336246535019133711157158298, c.GetLength(), model.Threshold)
	assert.Equal(4, c.GetClosestEdges().Len())

	// Validate that the clone is updated correctly.
	assert.Len(clone.GetUnattachedVertices(), 1)
	assert.True(clone.GetUnattachedVertices()[clone.Vertices[2]])

	assert.InDelta(12.6713030746299724753709331208575, clone.GetLength(), model.Threshold)
	assert.Equal(4, clone.GetClosestEdges().Len())

	// No clone on third update of c - vertex {1,1} to edge {0,0}->{0,3} or to edge {3,0}->{0,0}
	assert.Nil(c.CloneAndUpdate())

	assert.Len(c.GetUnattachedVertices(), 0)

	assert.InDelta(13.1888151645263866585819781087708, c.GetLength(), model.Threshold)
	assert.Equal(3, c.GetClosestEdges().Len())

	// No clone on first update of clone - vertex {1,1} to edge {0,0}->{1,2.1}
	assert.Nil(clone.CloneAndUpdate())

	assert.Len(clone.GetUnattachedVertices(), 0)

	assert.InDelta(12.85957596708046608051799727042668, clone.GetLength(), model.Threshold)
	assert.Equal(3, clone.GetClosestEdges().Len())
}

func TestDelete_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)

	c := experimental.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}, model2d.DeduplicateVertices, model2d.BuildPerimiter)
	c.Prepare()
	c.BuildPerimiter()

	c.CloneAndUpdate() //No clone
	clone := c.CloneAndUpdate().(*experimental.HeapableCircuitMinClonesLimited)
	clone.CloneAndUpdate() //No clone
	cloneOfClone := clone.CloneAndUpdate().(*experimental.HeapableCircuitMinClonesLimited)

	c.Delete()
	assert.Len(c.GetUnattachedVertices(), 0)
	assert.Nil(c.GetAttachedEdges())
	assert.NotNil(c.GetAttachedVertices())
	assert.Len(c.GetAttachedVertices(), 0)
	assert.Nil(c.GetClosestEdges())
	assert.Nil(c.Vertices)

	assert.Len(clone.GetUnattachedVertices(), 1)
	assert.Len(clone.GetAttachedEdges(), 7)
	assert.Len(clone.GetAttachedVertices(), 7)
	assert.NotNil(clone.GetClosestEdges())
	assert.Equal(5, clone.GetClosestEdges().Len())
	assert.Len(clone.Vertices, 8)

	assert.Len(cloneOfClone.GetUnattachedVertices(), 1)
	assert.Len(cloneOfClone.GetAttachedEdges(), 7)
	assert.Len(cloneOfClone.GetAttachedVertices(), 7)
	assert.NotNil(cloneOfClone.GetClosestEdges())
	assert.Equal(3, cloneOfClone.GetClosestEdges().Len())
	assert.Len(cloneOfClone.Vertices, 8)

	cloneOfClone.Delete()
	assert.Len(cloneOfClone.GetUnattachedVertices(), 0)
	assert.Nil(cloneOfClone.GetAttachedEdges())
	assert.NotNil(cloneOfClone.GetAttachedVertices())
	assert.Len(cloneOfClone.GetAttachedVertices(), 0)
	assert.Nil(cloneOfClone.GetClosestEdges())
	assert.Nil(cloneOfClone.Vertices)

	assert.Len(clone.GetUnattachedVertices(), 1)
	assert.Len(clone.GetAttachedEdges(), 7)
	assert.Len(clone.GetAttachedVertices(), 7)
	assert.NotNil(clone.GetClosestEdges())
	assert.Equal(5, clone.GetClosestEdges().Len())
	assert.Len(clone.Vertices, 8)
	clone.Delete()
}

func TestPrepare_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)
	c := experimental.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
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
	}, model2d.DeduplicateVertices, model2d.BuildPerimiter)

	c.Prepare()

	assert.NotNil(c.Vertices)
	assert.Len(c.Vertices, 7)
	assert.ElementsMatch(c.Vertices, []*model2d.Vertex2D{
		model2d.NewVertex2D(-15+model.Threshold/3.0, -15-model.Threshold/3.0),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(7, 6),
		model2d.NewVertex2D(15, -15),
	})

	// assert.NotNil(c.GetConvexVertices())
	// assert.Len(c.GetConvexVertices(), 0)

	assert.NotNil(c.GetUnattachedVertices())
	assert.Len(c.GetUnattachedVertices(), 0)

	assert.Equal(0.0, c.GetLength())
	assert.Equal(0.0, c.GetLength())
	assert.Equal(0.0, c.GetLengthWithNext())

	assert.NotNil(c.GetClosestEdges())
	assert.Equal(0, c.GetClosestEdges().Len())

	assert.NotNil(c.GetAttachedVertices())
	assert.Len(c.GetAttachedVertices(), 0)

	assert.NotNil(c.GetAttachedEdges())
	assert.Len(c.GetAttachedEdges(), 0)
}
