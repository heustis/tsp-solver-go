package model_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerimeter_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)
	circuit := model.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
	circuit.Prepare()
	circuit.BuildPerimiter()

	assert.Len(circuit.Vertices, 8)

	assert.Len(circuit.GetAttachedVertices(), 5)
	assert.Equal(model2d.NewVertex2D(-15, -15), circuit.GetAttachedVertices()[0])
	assert.Equal(model2d.NewVertex2D(15, -15), circuit.GetAttachedVertices()[1])
	assert.Equal(model2d.NewVertex2D(9, 6), circuit.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(3, 13), circuit.GetAttachedVertices()[3])
	assert.Equal(model2d.NewVertex2D(-7, 6), circuit.GetAttachedVertices()[4])
	assert.Equal(circuit.GetAttachedVertices(), circuit.GetAttachedVertices())

	assert.Len(circuit.GetAttachedEdges(), 5)
	assert.Equal(circuit.Vertices[0].EdgeTo(circuit.Vertices[7]), circuit.GetAttachedEdges()[0])
	assert.Equal(circuit.Vertices[7].EdgeTo(circuit.Vertices[6]), circuit.GetAttachedEdges()[1])
	assert.Equal(circuit.Vertices[6].EdgeTo(circuit.Vertices[4]), circuit.GetAttachedEdges()[2])
	assert.Equal(circuit.Vertices[4].EdgeTo(circuit.Vertices[1]), circuit.GetAttachedEdges()[3])
	assert.Equal(circuit.Vertices[1].EdgeTo(circuit.Vertices[0]), circuit.GetAttachedEdges()[4])

	expectedLength := 0.0
	for _, edge := range circuit.GetAttachedEdges() {
		expectedLength += edge.GetLength()
	}
	assert.InDelta(expectedLength, circuit.GetLength(), model.Threshold)

	assert.Len(circuit.GetUnattachedVertices(), 3)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[3]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[5]])

	assert.InDelta(95.738634795112368+0.763503994948632, circuit.GetLengthWithNext(), model.Threshold)

	assert.Equal(9, circuit.GetClosestEdges().Len())
	assert.Equal(&model.DistanceToEdge{
		Edge:     circuit.GetAttachedEdges()[1],
		Vertex:   circuit.Vertices[5],
		Distance: 0.763503994948632,
	}, circuit.GetClosestEdges().PopHeap())
	assert.Equal(&model.DistanceToEdge{
		Edge:     circuit.GetAttachedEdges()[2],
		Vertex:   circuit.Vertices[5],
		Distance: 1.628650237136812,
	}, circuit.GetClosestEdges().PopHeap())
	assert.Equal(&model.DistanceToEdge{
		Edge:     circuit.GetAttachedEdges()[1],
		Vertex:   circuit.Vertices[3],
		Distance: 5.854324418695558,
	}, circuit.GetClosestEdges().PopHeap())
	assert.Equal(&model.DistanceToEdge{
		Edge:     circuit.GetAttachedEdges()[4],
		Vertex:   circuit.Vertices[2],
		Distance: 7.9605428386450825,
	}, circuit.GetClosestEdges().PopHeap())
}

func TestAttachAndMove_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)

	circuit := model.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(0, 0),   // Index 0 after sorting
		model2d.NewVertex2D(0, 3),   // Index 1 after sorting
		model2d.NewVertex2D(3, 3),   // Index 5 after sorting
		model2d.NewVertex2D(3, 0),   // Index 4 after sorting
		model2d.NewVertex2D(1, 1),   // Index 2 after sorting
		model2d.NewVertex2D(1, 2.1), // Index 3 after sorting
	}, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
	circuit.Prepare()
	circuit.BuildPerimiter()

	testVert := circuit.GetClosestEdges().PopHeap().(*model.DistanceToEdge)
	circuit.AttachVertex(testVert)

	assert.InDelta(12.5385336246535019133711157158298, circuit.GetLength(), model.Threshold)
	assert.Equal(5, circuit.GetClosestEdges().Len())

	circuit.MoveVertex(&model.DistanceToEdge{
		Vertex:   testVert.Vertex,
		Edge:     circuit.GetAttachedEdges()[1],
		Distance: -1.0,
	})

	assert.Len(circuit.GetUnattachedVertices(), 1)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])

	assert.Equal(11.5385336246535019133711157158298, circuit.GetLength())
	assert.Equal(3, circuit.GetClosestEdges().Len())
}

func TestAttachAndMoveIndexZero_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)

	circuit := model.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(0, 0),   // Index 0 after sorting
		model2d.NewVertex2D(0, 3),   // Index 1 after sorting
		model2d.NewVertex2D(3, 3),   // Index 5 after sorting
		model2d.NewVertex2D(3, 0),   // Index 4 after sorting
		model2d.NewVertex2D(1, 1),   // Index 2 after sorting
		model2d.NewVertex2D(1, 2.1), // Index 3 after sorting
	}, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
	circuit.Prepare()
	circuit.BuildPerimiter()

	testVert := &model.DistanceToEdge{
		Vertex:   circuit.Vertices[2],
		Edge:     circuit.GetAttachedEdges()[0],
		Distance: 1.0,
	}
	circuit.AttachVertex(testVert)
	assert.Equal([]model.CircuitEdge{
		circuit.Vertices[4].EdgeTo(circuit.Vertices[2]),
		circuit.Vertices[2].EdgeTo(circuit.Vertices[5]),
		circuit.Vertices[5].EdgeTo(circuit.Vertices[1]),
		circuit.Vertices[1].EdgeTo(circuit.Vertices[0]),
		circuit.Vertices[0].EdgeTo(circuit.Vertices[4]),
	}, circuit.GetAttachedEdges())
	assert.Equal([]model.CircuitVertex{
		circuit.Vertices[4],
		circuit.Vertices[2],
		circuit.Vertices[5],
		circuit.Vertices[1],
		circuit.Vertices[0],
	}, circuit.GetAttachedVertices())

	circuit.MoveVertex(&model.DistanceToEdge{
		Vertex:   circuit.Vertices[2],
		Edge:     circuit.GetAttachedEdges()[4],
		Distance: 2.25,
	})
	assert.Equal([]model.CircuitEdge{
		circuit.Vertices[4].EdgeTo(circuit.Vertices[5]),
		circuit.Vertices[5].EdgeTo(circuit.Vertices[1]),
		circuit.Vertices[1].EdgeTo(circuit.Vertices[0]),
		circuit.Vertices[0].EdgeTo(circuit.Vertices[2]),
		circuit.Vertices[2].EdgeTo(circuit.Vertices[4]),
	}, circuit.GetAttachedEdges())
	assert.Equal([]model.CircuitVertex{
		circuit.Vertices[4],
		circuit.Vertices[5],
		circuit.Vertices[1],
		circuit.Vertices[0],
		circuit.Vertices[2],
	}, circuit.GetAttachedVertices())

	circuit.MoveVertex(&model.DistanceToEdge{
		Vertex:   circuit.Vertices[2],
		Edge:     circuit.GetAttachedEdges()[0],
		Distance: 3.5,
	})
	assert.Equal([]model.CircuitEdge{
		circuit.Vertices[4].EdgeTo(circuit.Vertices[2]),
		circuit.Vertices[2].EdgeTo(circuit.Vertices[5]),
		circuit.Vertices[5].EdgeTo(circuit.Vertices[1]),
		circuit.Vertices[1].EdgeTo(circuit.Vertices[0]),
		circuit.Vertices[0].EdgeTo(circuit.Vertices[4]),
	}, circuit.GetAttachedEdges())
	assert.Equal([]model.CircuitVertex{
		circuit.Vertices[4],
		circuit.Vertices[2],
		circuit.Vertices[5],
		circuit.Vertices[1],
		circuit.Vertices[0],
	}, circuit.GetAttachedVertices())

	assert.Len(circuit.GetUnattachedVertices(), 1)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[3]])

	assert.Equal(12.0+6.75, circuit.GetLength()) // 6.75 comes from the distances in each DistanceToEdge
	assert.Equal(4, circuit.GetClosestEdges().Len())
}

func TestAttachAndMoveLastIndex_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)

	circuit := model.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(0, 0),   // Index 0 after sorting
		model2d.NewVertex2D(0, 3),   // Index 1 after sorting
		model2d.NewVertex2D(3, 3),   // Index 5 after sorting
		model2d.NewVertex2D(3, 0),   // Index 4 after sorting
		model2d.NewVertex2D(1, 1),   // Index 2 after sorting
		model2d.NewVertex2D(1, 2.1), // Index 3 after sorting
	}, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
	circuit.Prepare()
	circuit.BuildPerimiter()

	testVert := &model.DistanceToEdge{
		Vertex:   circuit.Vertices[2],
		Edge:     circuit.GetAttachedEdges()[3],
		Distance: 0.5,
	}
	circuit.AttachVertex(testVert)
	assert.Equal([]model.CircuitEdge{
		circuit.Vertices[4].EdgeTo(circuit.Vertices[5]),
		circuit.Vertices[5].EdgeTo(circuit.Vertices[1]),
		circuit.Vertices[1].EdgeTo(circuit.Vertices[0]),
		circuit.Vertices[0].EdgeTo(circuit.Vertices[2]),
		circuit.Vertices[2].EdgeTo(circuit.Vertices[4]),
	}, circuit.GetAttachedEdges())
	assert.Equal([]model.CircuitVertex{
		circuit.Vertices[4],
		circuit.Vertices[5],
		circuit.Vertices[1],
		circuit.Vertices[0],
		circuit.Vertices[2],
	}, circuit.GetAttachedVertices())

	circuit.MoveVertex(&model.DistanceToEdge{
		Vertex:   circuit.Vertices[2],
		Edge:     circuit.GetAttachedEdges()[2],
		Distance: 0.25,
	})
	assert.Equal([]model.CircuitEdge{
		circuit.Vertices[4].EdgeTo(circuit.Vertices[5]),
		circuit.Vertices[5].EdgeTo(circuit.Vertices[1]),
		circuit.Vertices[1].EdgeTo(circuit.Vertices[2]),
		circuit.Vertices[2].EdgeTo(circuit.Vertices[0]),
		circuit.Vertices[0].EdgeTo(circuit.Vertices[4]),
	}, circuit.GetAttachedEdges())
	assert.Equal([]model.CircuitVertex{
		circuit.Vertices[4],
		circuit.Vertices[5],
		circuit.Vertices[1],
		circuit.Vertices[2],
		circuit.Vertices[0],
	}, circuit.GetAttachedVertices())

	assert.Len(circuit.GetUnattachedVertices(), 1)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[3]])

	assert.Equal(12.0+0.75, circuit.GetLength())
	assert.Equal(4, circuit.GetClosestEdges().Len())
}

func TestAttachShouldPanicIfEdgeIsNotInCircuit_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)

	circuit := model.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
	circuit.Prepare()
	circuit.BuildPerimiter()

	assert.Panics(func() {
		circuit.AttachVertex(&model.DistanceToEdge{
			Vertex:   circuit.Vertices[0],
			Edge:     model2d.NewVertex2D(3, 3).EdgeTo(model2d.NewVertex2D(5, 5)),
			Distance: 1.234,
		})
	})
}

func TestMoveShouldPanicIfEdgeIsNotInCircuit_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)

	circuit := model.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
	circuit.Prepare()
	circuit.BuildPerimiter()

	assert.Panics(func() {
		circuit.MoveVertex(&model.DistanceToEdge{
			Vertex:   circuit.Vertices[0],
			Edge:     model2d.NewVertex2D(3, 3).EdgeTo(model2d.NewVertex2D(5, 5)),
			Distance: 1.234,
		})
	})
}

func TestCloneAndUpdate_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)

	circuit := model.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
	circuit.Prepare()
	circuit.BuildPerimiter()

	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[3]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[5]])

	// Index 5 should attach to edge 15,-15 -> 9,6
	assert.Nil(circuit.CloneAndUpdate())
	assert.Len(circuit.GetUnattachedVertices(), 2)
	assert.False(circuit.GetUnattachedVertices()[circuit.Vertices[5]])
	assert.Len(circuit.GetAttachedVertices(), 6)
	assert.Len(circuit.GetAttachedEdges(), 6)
	assert.Equal(8, circuit.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Edge:     circuit.GetAttachedEdges()[3],
		Vertex:   circuit.Vertices[5],
		Distance: 0.8651462421881799,
	}, circuit.GetClosestEdges().Peek())

	// Index 5 should attach to edge 9,6 -> 3,13, this requires cloning since index 5 is already attached.
	clone := circuit.CloneAndUpdate().(*model.HeapableCircuitMinClonesLimited)
	assert.Len(circuit.GetUnattachedVertices(), 2)
	assert.Len(clone.GetUnattachedVertices(), 2)
	assert.False(circuit.GetUnattachedVertices()[circuit.Vertices[5]])
	assert.False(clone.GetUnattachedVertices()[circuit.Vertices[5]])
	assert.Len(circuit.GetAttachedVertices(), 6)
	assert.Len(clone.GetAttachedVertices(), 6)
	assert.Len(circuit.GetAttachedEdges(), 6)
	assert.Len(clone.GetAttachedEdges(), 6)
	assert.Equal(7, circuit.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Edge:     circuit.GetAttachedEdges()[1],
		Vertex:   circuit.Vertices[3],
		Distance: 5.09082042374693,
	}, circuit.GetClosestEdges().Peek())

	assert.Equal(&model.DistanceToEdge{
		Edge:     clone.GetAttachedEdges()[1],
		Vertex:   clone.Vertices[3],
		Distance: 5.854324418695558,
	}, clone.GetClosestEdges().Peek())

	// Index 3 should attach to edge 1, no cloning required
	assert.Nil(circuit.CloneAndUpdate())
	assert.Nil(clone.CloneAndUpdate())
	assert.Len(circuit.GetUnattachedVertices(), 1)
	assert.Equal(6, circuit.GetClosestEdges().Len())
	assert.Len(clone.GetUnattachedVertices(), 1)
	assert.Equal(6, clone.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Edge:     circuit.GetAttachedEdges()[1],
		Vertex:   circuit.Vertices[2],
		Distance: 5.003830723297881,
	}, circuit.GetClosestEdges().Peek())

	assert.Equal(&model.DistanceToEdge{
		Edge:     clone.GetAttachedEdges()[4],
		Vertex:   clone.Vertices[3],
		Distance: 4.782762261113314,
	}, clone.GetClosestEdges().Peek())

	// Index 2 should attach to edge 1, no cloning required
	assert.Nil(circuit.CloneAndUpdate())
	// Index 3 should move to edge 4, cloning required
	cloneOfClone, okay := clone.CloneAndUpdate().(*model.HeapableCircuitMinClonesLimited)
	assert.True(okay)

	assert.Len(circuit.GetUnattachedVertices(), 0)
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

	circuit := model.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(0, 0),   // Index 0 after sorting
		model2d.NewVertex2D(0, 3),   // Index 1 after sorting
		model2d.NewVertex2D(3, 3),   // Index 5 after sorting
		model2d.NewVertex2D(3, 0),   // Index 4 after sorting
		model2d.NewVertex2D(1, 1),   // Index 2 after sorting
		model2d.NewVertex2D(1, 2.1), // Index 3 after sorting
	}, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
	circuit.Prepare()
	circuit.BuildPerimiter()

	assert.Len(circuit.GetUnattachedVertices(), 2)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[3]])

	assert.Equal(12.0, circuit.GetLength())
	assert.Equal(6, circuit.GetClosestEdges().Len())

	// No clone on first attachment - vertex {1,2.1} to edge {0,3}->{3,3}
	assert.Nil(circuit.CloneAndUpdate())

	assert.Len(circuit.GetUnattachedVertices(), 1)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])

	assert.InDelta(12.5385336246535019133711157158298, circuit.GetLength(), model.Threshold)
	assert.Equal(5, circuit.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Distance: 0.1327694499764709,
		Vertex:   circuit.Vertices[3],
		Edge:     circuit.GetAttachedEdges()[3],
	}, circuit.GetClosestEdges().Peek())

	// Clone on second attachment - vertex {1,2.1} to edge {0,0}->{0,3}
	clone := circuit.CloneAndUpdate().(*model.HeapableCircuitMinClonesLimited)

	// Validate that the first circuit is unchanged.
	assert.Len(circuit.GetUnattachedVertices(), 1)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])

	assert.InDelta(12.5385336246535019133711157158298, circuit.GetLength(), model.Threshold)
	assert.Equal(4, circuit.GetClosestEdges().Len())

	// Validate that the clone is updated correctly.
	assert.Len(clone.GetUnattachedVertices(), 1)
	assert.True(clone.GetUnattachedVertices()[clone.Vertices[2]])

	assert.InDelta(12.6713030746299724753709331208575, clone.GetLength(), model.Threshold)
	assert.Equal(4, clone.GetClosestEdges().Len())

	// No clone on third update of circuit - vertex {1,1} to edge {0,0}->{0,3} or to edge {3,0}->{0,0}
	assert.Nil(circuit.CloneAndUpdate())

	assert.Len(circuit.GetUnattachedVertices(), 0)

	assert.InDelta(13.1888151645263866585819781087708, circuit.GetLength(), model.Threshold)
	assert.Equal(3, circuit.GetClosestEdges().Len())

	// No clone on first update of clone - vertex {1,1} to edge {0,0}->{1,2.1}
	assert.Nil(clone.CloneAndUpdate())

	assert.Len(clone.GetUnattachedVertices(), 0)

	assert.InDelta(12.85957596708046608051799727042668, clone.GetLength(), model.Threshold)
	assert.Equal(3, clone.GetClosestEdges().Len())
}

func TestDelete_HeapMinClonesLimited(t *testing.T) {
	assert := assert.New(t)

	circuit := model.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
	circuit.Prepare()
	circuit.BuildPerimiter()

	circuit.CloneAndUpdate() //No clone
	clone := circuit.CloneAndUpdate().(*model.HeapableCircuitMinClonesLimited)
	clone.CloneAndUpdate() //No clone
	cloneOfClone := clone.CloneAndUpdate().(*model.HeapableCircuitMinClonesLimited)

	circuit.Delete()
	assert.Len(circuit.GetUnattachedVertices(), 0)
	assert.Nil(circuit.GetAttachedEdges())
	assert.NotNil(circuit.GetAttachedVertices())
	assert.Len(circuit.GetAttachedVertices(), 0)
	assert.Nil(circuit.GetClosestEdges())
	assert.Nil(circuit.Vertices)

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
	circuit := model.NewHeapableCircuitMinClonesLimited([]model.CircuitVertex{
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
	}, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})

	circuit.Prepare()

	assert.NotNil(circuit.Vertices)
	assert.Len(circuit.Vertices, 7)
	assert.ElementsMatch(circuit.Vertices, []*model2d.Vertex2D{
		model2d.NewVertex2D(-15+model.Threshold/3.0, -15-model.Threshold/3.0),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(7, 6),
		model2d.NewVertex2D(15, -15),
	})

	// assert.NotNil(circuit.GetConvexVertices())
	// assert.Len(circuit.GetConvexVertices(), 0)

	assert.NotNil(circuit.GetUnattachedVertices())
	assert.Len(circuit.GetUnattachedVertices(), 0)

	assert.Equal(0.0, circuit.GetLength())
	assert.Equal(0.0, circuit.GetLength())
	assert.Equal(0.0, circuit.GetLengthWithNext())

	assert.NotNil(circuit.GetClosestEdges())
	assert.Equal(0, circuit.GetClosestEdges().Len())

	assert.NotNil(circuit.GetAttachedVertices())
	assert.Len(circuit.GetAttachedVertices(), 0)

	assert.NotNil(circuit.GetAttachedEdges())
	assert.Len(circuit.GetAttachedEdges(), 0)
}
