package model_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerimeter_HeapMinClones(t *testing.T) {
	assert := assert.New(t)
	circuit := model.CreateHeapableCircuitMinClones([]model.CircuitVertex{
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
	assert.Equal(model2d.NewEdge2D(circuit.Vertices[0].(*model2d.Vertex2D), circuit.Vertices[7].(*model2d.Vertex2D)), circuit.GetAttachedEdges()[0])
	assert.Equal(model2d.NewEdge2D(circuit.Vertices[7].(*model2d.Vertex2D), circuit.Vertices[6].(*model2d.Vertex2D)), circuit.GetAttachedEdges()[1])
	assert.Equal(model2d.NewEdge2D(circuit.Vertices[6].(*model2d.Vertex2D), circuit.Vertices[4].(*model2d.Vertex2D)), circuit.GetAttachedEdges()[2])
	assert.Equal(model2d.NewEdge2D(circuit.Vertices[4].(*model2d.Vertex2D), circuit.Vertices[1].(*model2d.Vertex2D)), circuit.GetAttachedEdges()[3])
	assert.Equal(model2d.NewEdge2D(circuit.Vertices[1].(*model2d.Vertex2D), circuit.Vertices[0].(*model2d.Vertex2D)), circuit.GetAttachedEdges()[4])

	expectedLength := 0.0
	for _, edge := range circuit.GetAttachedEdges() {
		expectedLength += edge.GetLength()
	}
	assert.InDelta(expectedLength, circuit.GetLength(), model.Threshold)

	assert.NotNil(circuit.GetConvexVertices())
	assert.Len(circuit.GetConvexVertices(), 5)
	assert.True(circuit.GetConvexVertices()[circuit.Vertices[0]])
	assert.True(circuit.GetConvexVertices()[circuit.Vertices[1]])
	assert.True(circuit.GetConvexVertices()[circuit.Vertices[4]])
	assert.True(circuit.GetConvexVertices()[circuit.Vertices[6]])
	assert.True(circuit.GetConvexVertices()[circuit.Vertices[7]])

	assert.Len(circuit.GetUnattachedVertices(), 3)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[3]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[5]])

	assert.InDelta(95.738634795112368+0.763503994948632, circuit.GetLengthWithNext(), model.Threshold)

	assert.Equal(15, circuit.GetClosestEdges().Len())
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

func TestAttachAndDetach(t *testing.T) {
	assert := assert.New(t)

	circuit := model.CreateHeapableCircuitMinClones([]model.CircuitVertex{
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
	circuit.DetachVertex(testVert.Vertex)

	assert.Len(circuit.GetUnattachedVertices(), 2)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[3]])

	assert.Equal(12.0, circuit.GetLength())
	assert.Equal(7, circuit.GetClosestEdges().Len())

	circuit.AttachVertex(testVert)
	assert.Len(circuit.GetUnattachedVertices(), 1)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])

	assert.InDelta(12.5385336246535019133711157158298, circuit.GetLength(), model.Threshold)
	assert.Equal(8, circuit.GetClosestEdges().Len())
}

func TestAttachAndDetachIndexZero(t *testing.T) {
	assert := assert.New(t)

	circuit := model.CreateHeapableCircuitMinClones([]model.CircuitVertex{
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
		Distance: circuit.GetAttachedEdges()[0].DistanceIncrease(circuit.Vertices[2]),
	}
	circuit.AttachVertex(testVert)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(circuit.Vertices[4].(*model2d.Vertex2D), circuit.Vertices[2].(*model2d.Vertex2D)),
		model2d.NewEdge2D(circuit.Vertices[2].(*model2d.Vertex2D), circuit.Vertices[5].(*model2d.Vertex2D)),
		model2d.NewEdge2D(circuit.Vertices[5].(*model2d.Vertex2D), circuit.Vertices[1].(*model2d.Vertex2D)),
		model2d.NewEdge2D(circuit.Vertices[1].(*model2d.Vertex2D), circuit.Vertices[0].(*model2d.Vertex2D)),
		model2d.NewEdge2D(circuit.Vertices[0].(*model2d.Vertex2D), circuit.Vertices[4].(*model2d.Vertex2D)),
	}, circuit.GetAttachedEdges())
	assert.Equal([]model.CircuitVertex{
		circuit.Vertices[4],
		circuit.Vertices[2],
		circuit.Vertices[5],
		circuit.Vertices[1],
		circuit.Vertices[0],
	}, circuit.GetAttachedVertices())

	circuit.DetachVertex(testVert.Vertex)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(circuit.Vertices[4].(*model2d.Vertex2D), circuit.Vertices[5].(*model2d.Vertex2D)),
		model2d.NewEdge2D(circuit.Vertices[5].(*model2d.Vertex2D), circuit.Vertices[1].(*model2d.Vertex2D)),
		model2d.NewEdge2D(circuit.Vertices[1].(*model2d.Vertex2D), circuit.Vertices[0].(*model2d.Vertex2D)),
		model2d.NewEdge2D(circuit.Vertices[0].(*model2d.Vertex2D), circuit.Vertices[4].(*model2d.Vertex2D)),
	}, circuit.GetAttachedEdges())
	assert.Equal([]model.CircuitVertex{
		circuit.Vertices[4],
		circuit.Vertices[5],
		circuit.Vertices[1],
		circuit.Vertices[0],
	}, circuit.GetAttachedVertices())

	assert.Len(circuit.GetUnattachedVertices(), 2)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[3]])

	assert.Equal(12.0, circuit.GetLength())
	assert.Equal(7, circuit.GetClosestEdges().Len())

	circuit.AttachVertex(testVert)
	assert.Len(circuit.GetUnattachedVertices(), 1)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[3]])

	assert.InDelta(12+testVert.Distance, circuit.GetLength(), model.Threshold)
	assert.Equal(8, circuit.GetClosestEdges().Len())
}

func TestAttachAndDetachLastIndex(t *testing.T) {
	assert := assert.New(t)

	circuit := model.CreateHeapableCircuitMinClones([]model.CircuitVertex{
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
		Distance: circuit.GetAttachedEdges()[3].DistanceIncrease(circuit.Vertices[2]),
	}
	circuit.AttachVertex(testVert)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(circuit.Vertices[4].(*model2d.Vertex2D), circuit.Vertices[5].(*model2d.Vertex2D)),
		model2d.NewEdge2D(circuit.Vertices[5].(*model2d.Vertex2D), circuit.Vertices[1].(*model2d.Vertex2D)),
		model2d.NewEdge2D(circuit.Vertices[1].(*model2d.Vertex2D), circuit.Vertices[0].(*model2d.Vertex2D)),
		model2d.NewEdge2D(circuit.Vertices[0].(*model2d.Vertex2D), circuit.Vertices[2].(*model2d.Vertex2D)),
		model2d.NewEdge2D(circuit.Vertices[2].(*model2d.Vertex2D), circuit.Vertices[4].(*model2d.Vertex2D)),
	}, circuit.GetAttachedEdges())
	assert.Equal([]model.CircuitVertex{
		circuit.Vertices[4],
		circuit.Vertices[5],
		circuit.Vertices[1],
		circuit.Vertices[0],
		circuit.Vertices[2],
	}, circuit.GetAttachedVertices())

	circuit.DetachVertex(testVert.Vertex)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(circuit.Vertices[4].(*model2d.Vertex2D), circuit.Vertices[5].(*model2d.Vertex2D)),
		model2d.NewEdge2D(circuit.Vertices[5].(*model2d.Vertex2D), circuit.Vertices[1].(*model2d.Vertex2D)),
		model2d.NewEdge2D(circuit.Vertices[1].(*model2d.Vertex2D), circuit.Vertices[0].(*model2d.Vertex2D)),
		model2d.NewEdge2D(circuit.Vertices[0].(*model2d.Vertex2D), circuit.Vertices[4].(*model2d.Vertex2D)),
	}, circuit.GetAttachedEdges())
	assert.Equal([]model.CircuitVertex{
		circuit.Vertices[4],
		circuit.Vertices[5],
		circuit.Vertices[1],
		circuit.Vertices[0],
	}, circuit.GetAttachedVertices())

	assert.Len(circuit.GetUnattachedVertices(), 2)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[3]])

	assert.Equal(12.0, circuit.GetLength())
	assert.Equal(7, circuit.GetClosestEdges().Len())

	circuit.AttachVertex(testVert)
	assert.Len(circuit.GetUnattachedVertices(), 1)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[3]])

	assert.InDelta(12+testVert.Distance, circuit.GetLength(), model.Threshold)
	assert.Equal(8, circuit.GetClosestEdges().Len())
}

func TestAttachShouldPanicIfEdgeIsNotInCircuit(t *testing.T) {
	assert := assert.New(t)

	circuit := model.CreateHeapableCircuitMinClones([]model.CircuitVertex{
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
			Edge:     model2d.NewEdge2D(model2d.NewVertex2D(3, 3), model2d.NewVertex2D(5, 5)),
			Distance: 1.234,
		})
	})
}

func TestCloneAndUpdate_HeapMinClones(t *testing.T) {
	assert := assert.New(t)

	circuit := model.CreateHeapableCircuitMinClones([]model.CircuitVertex{
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
	assert.Equal(16, circuit.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Edge:     circuit.GetAttachedEdges()[3],
		Vertex:   circuit.Vertices[5],
		Distance: 0.8651462421881799,
	}, circuit.GetClosestEdges().Peek())

	// Index 5 should attach to edge 9,6 -> 3,13, this requires cloning since index 5 is already attached.
	clone := circuit.CloneAndUpdate().(*model.HeapableCircuitMinClones)
	assert.Len(circuit.GetUnattachedVertices(), 2)
	assert.Len(clone.GetUnattachedVertices(), 2)
	assert.False(circuit.GetUnattachedVertices()[circuit.Vertices[5]])
	assert.False(clone.GetUnattachedVertices()[circuit.Vertices[5]])
	assert.Len(circuit.GetAttachedVertices(), 6)
	assert.Len(clone.GetAttachedVertices(), 6)
	assert.Len(circuit.GetAttachedEdges(), 6)
	assert.Len(clone.GetAttachedEdges(), 6)
	assert.Equal(15, circuit.GetClosestEdges().Len())

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
	assert.Equal(15, circuit.GetClosestEdges().Len())
	assert.Len(clone.GetUnattachedVertices(), 1)
	assert.Equal(12, clone.GetClosestEdges().Len())

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
	cloneOfClone, okay := clone.CloneAndUpdate().(*model.HeapableCircuitMinClones)
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

	// Index 2 should move to edge 3, cloning required
	assert.Equal(&model.DistanceToEdge{
		Edge:     clone.GetAttachedEdges()[3],
		Vertex:   clone.Vertices[2],
		Distance: 0.32754172885551824,
	}, clone.GetClosestEdges().Peek())
	assert.NotNil(clone.CloneAndUpdate())

	// Index 2 should move to edge 4, cloning required
	assert.Equal(&model.DistanceToEdge{
		Edge:     cloneOfClone.GetAttachedEdges()[4],
		Vertex:   cloneOfClone.Vertices[2],
		Distance: 3.341664064126334,
	}, cloneOfClone.GetClosestEdges().Peek())
	assert.Nil(cloneOfClone.CloneAndUpdate())
}

func TestCloneAndUpdate_HeapMinClones_Distances(t *testing.T) {
	assert := assert.New(t)

	circuit := model.CreateHeapableCircuitMinClones([]model.CircuitVertex{
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
	assert.Equal(8, circuit.GetClosestEdges().Len())

	// No clone on first attachment - vertex {1,2.1} to edge {0,3}->{3,3}
	assert.Nil(circuit.CloneAndUpdate())

	assert.Len(circuit.GetUnattachedVertices(), 1)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])

	assert.InDelta(12.5385336246535019133711157158298, circuit.GetLength(), model.Threshold)
	assert.Equal(8, circuit.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Distance: 0.1327694499764709,
		Vertex:   circuit.Vertices[3],
		Edge:     circuit.GetAttachedEdges()[3],
	}, circuit.GetClosestEdges().Peek())

	// Clone on second attachment - vertex {1,2.1} to edge {0,0}->{0,3}
	clone := circuit.CloneAndUpdate().(*model.HeapableCircuitMinClones)

	// Validate that the first circuit is unchanged.
	assert.Len(circuit.GetUnattachedVertices(), 1)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])

	assert.InDelta(12.5385336246535019133711157158298, circuit.GetLength(), model.Threshold)
	assert.Equal(7, circuit.GetClosestEdges().Len())

	// Validate that the clone is updated correctly.
	assert.Len(clone.GetUnattachedVertices(), 1)
	assert.True(clone.GetUnattachedVertices()[clone.Vertices[2]])

	assert.InDelta(12.6713030746299724753709331208575, clone.GetLength(), model.Threshold)
	assert.Equal(5, clone.GetClosestEdges().Len())

	// No clone on third update of circuit - vertex {1,1} to edge {0,0}->{0,3} or to edge {3,0}->{0,0}
	assert.Nil(circuit.CloneAndUpdate())

	assert.Len(circuit.GetUnattachedVertices(), 0)

	assert.InDelta(13.1888151645263866585819781087708, circuit.GetLength(), model.Threshold)
	assert.Equal(6, circuit.GetClosestEdges().Len())

	// No clone on first update of clone - vertex {1,1} to edge {0,0}->{1,2.1}
	assert.Nil(clone.CloneAndUpdate())

	assert.Len(clone.GetUnattachedVertices(), 0)

	assert.InDelta(12.85957596708046608051799727042668, clone.GetLength(), model.Threshold)
	assert.Equal(4, clone.GetClosestEdges().Len())
}

func TestDelete_HeapMinClones(t *testing.T) {
	assert := assert.New(t)

	circuit := model.CreateHeapableCircuitMinClones([]model.CircuitVertex{
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
	clone := circuit.CloneAndUpdate().(*model.HeapableCircuitMinClones)
	clone.CloneAndUpdate() //No clone
	cloneOfClone := clone.CloneAndUpdate().(*model.HeapableCircuitMinClones)

	circuit.Delete()
	assert.Len(circuit.GetUnattachedVertices(), 0)
	assert.Nil(circuit.GetAttachedEdges())
	assert.Nil(circuit.GetAttachedVertices())
	assert.Nil(circuit.GetClosestEdges())
	assert.Nil(circuit.Vertices)

	assert.Len(clone.GetUnattachedVertices(), 1)
	assert.Len(clone.GetAttachedEdges(), 7)
	assert.Len(clone.GetAttachedVertices(), 7)
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
	assert.Nil(cloneOfClone.GetAttachedEdges())
	assert.Nil(cloneOfClone.GetAttachedVertices())
	assert.Nil(cloneOfClone.GetClosestEdges())
	assert.Nil(cloneOfClone.Vertices)

	assert.Len(clone.GetUnattachedVertices(), 1)
	assert.Len(clone.GetAttachedEdges(), 7)
	assert.Len(clone.GetAttachedVertices(), 7)
	assert.NotNil(clone.GetClosestEdges())
	assert.Equal(11, clone.GetClosestEdges().Len())
	assert.Len(clone.Vertices, 8)
	clone.Delete()
}

func TestPrepare_HeapMinClones(t *testing.T) {
	assert := assert.New(t)
	circuit := model.CreateHeapableCircuitMinClones([]model.CircuitVertex{
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

	assert.NotNil(circuit.GetConvexVertices())
	assert.Len(circuit.GetConvexVertices(), 0)

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

func TestSolve_HeapMinClones(t *testing.T) {
	assert := assert.New(t)

	testData := []struct {
		len            int
		vertices       string
		expectedLength float64
	}{
		// {
		// 	len:            10,
		// 	vertices:       `[{"x":449.0904385101078,"y":1163.6150486330282},{"x":2846.191007802421,"y":5564.992099820763},{"x":2961.377236141765,"y":2035.3395220031912},{"x":3102.158315380755,"y":6577.361560477899},{"x":3404.4109094786213,"y":3582.9840359952696},{"x":4111.1661756679205,"y":4054.9949515059243},{"x":4867.500584018192,"y":1366.616080844533},{"x":6533.580847366987,"y":4302.364845399589},{"x":8087.6596916825065,"y":5719.191723935985},{"x":8369.042936423306,"y":3826.201241779603}]`,
		// 	expectedLength: 24606.92092681067,
		// },
		// {
		// 	len:            10,
		// 	vertices:       `[{"x":5484.54767627217,"y":6102.141372143685},{"x":6028.193790687806,"y":3510.0105605018352},{"x":5707.958405221888,"y":1186.068951762566},{"x":7735.627198076895,"y":3632.719377795526},{"x":9568.080830249783,"y":4069.3575617048177},{"x":7737.393316180935,"y":8833.192379589624},{"x":6862.826656822809,"y":8261.45375922393},{"x":1504.8308948639935,"y":8612.33378158451},{"x":2317.672614461242,"y":8243.607064427804},{"x":3531.7754803836497,"y":6985.98682680876}]`,
		// 	expectedLength: 26357.758795626996,
		// },
		// {
		// 	len:            10,
		// 	vertices:       `[{"x":9191.254586900677,"y":9795.917503309616},{"x":3962.4333635960766,"y":9515.470044522293},{"x":843.3628316184105,"y":9809.035571758466},{"x":4859.924914952512,"y":7739.423447671773},{"x":6331.8517122526655,"y":6534.273967600264},{"x":5698.23063956672,"y":5465.990541535624},{"x":3493.995534094023,"y":2544.0574075888912},{"x":6346.678944926111,"y":4429.945106452629},{"x":7911.990559441331,"y":5629.649385235076},{"x":9790.128451765924,"y":4907.839724553531}]`,
		// 	expectedLength: 32020.69626589124,
		// },
		// TODO - debug
		// {
		// 	len:            10,
		// 	vertices:       `[{"x":2760.6690690740065,"y":3849.7666200324197},{"x":3857.2452017168252,"y":5455.857006708013},{"x":5383.324285557034,"y":6126.133439885709},{"x":9621.448154212852,"y":3209.9965998551984},{"x":7008.046221229183,"y":8290.39332435957},{"x":3701.482778524058,"y":9357.256620736556},{"x":2517.4047824536633,"y":9337.703635201875},{"x":1194.001834274779,"y":7662.718457736367},{"x":665.3346612768264,"y":724.4366074426538},{"x":2285.2883009226625,"y":597.3752679309797}]`,
		// 	expectedLength: 33132.797166,
		// },
	}

	for i, t := range testData {
		vertices2d := make([]*model2d.Vertex2D, t.len)
		err := json.Unmarshal([]byte(t.vertices), &vertices2d)

		vertices := make([]model.CircuitVertex, t.len)
		for i, v2d := range vertices2d {
			vertices[i] = v2d
		}

		assert.Nil(err, "Failed to unmarshal vertices for test=", i)
		actual := solveWithLogging_HeapMinClones(&model.HeapableCircuitMinClones{
			Vertices: vertices,
		})
		assert.InDelta(t.expectedLength, actual.GetLength(), model.Threshold)
	}
}

func solveWithLogging_HeapMinClones(circuit *model.HeapableCircuitMinClones) *model.HeapableCircuitMinClones {
	circuit.Prepare()
	circuit.BuildPerimiter()

	circuitHeap := model.NewHeap(func(a interface{}) float64 {
		return a.(*model.HeapableCircuitMinClones).GetLengthWithNext()
	})
	circuitHeap.PushHeap(circuit)

	next := circuitHeap.PopHeap().(*model.HeapableCircuitMinClones)
	for i := 0; len(next.GetUnattachedVertices()) > 0 || next.GetLengthWithNext() < next.GetLength(); next = circuitHeap.PopHeap().(*model.HeapableCircuitMinClones) {
		toAttach := next.GetClosestEdges().Peek()
		clone := next.CloneAndUpdate()
		circuitHeap.PushHeap(next)
		if clone != nil {
			circuitBytes, _ := json.Marshal(clone.(*model.HeapableCircuitMinClones).GetAttachedVertices())
			fmt.Printf("Step %d: Created clone=%p from existing=%p with \n\ttoAttach=%s\n\tcircuit=%s\n\theap=%s\n", i, clone, next, toAttach.(*model.DistanceToEdge).ToString(), string(circuitBytes), clone.(*model.HeapableCircuitMinClones).GetClosestEdges().ToString())
			circuitHeap.PushHeap(clone)
		} else {
			circuitBytes, _ := json.Marshal(next.GetAttachedVertices())
			fmt.Printf("Step %d: Updated existing=%p with \n\ttoAttach=%s\n\tcircuit=%s\n\theap=%s\n", i, next, toAttach.(*model.DistanceToEdge).ToString(), string(circuitBytes), next.GetClosestEdges().ToString())
		}
		i++
	}

	// clean up the heap and each circuitHeap
	circuitHeap.Delete()

	return next
}
