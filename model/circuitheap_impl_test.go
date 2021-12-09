package model_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerimeter_Heap(t *testing.T) {
	assert := assert.New(t)
	circuit := model.NewHeapableCircuitImpl([]model.CircuitVertex{
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

	assert.Len(circuit.Vertices, 8)

	circuit.BuildPerimiter()
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

	l := 0.0
	for _, edge := range circuit.GetAttachedEdges() {
		l += edge.GetLength()
	}
	assert.InDelta(l, circuit.GetLength(), model.Threshold)

	assert.Len(circuit.GetUnattachedVertices(), 3)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[3]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[5]])
	assert.Equal(circuit.GetUnattachedVertices(), circuit.GetUnattachedVertices())

	assert.InDelta(circuit.GetLength()+0.763503994948632, circuit.GetLengthWithNext(), model.Threshold)

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

func TestCloneAndUpdate(t *testing.T) {
	assert := assert.New(t)

	circuit := model.NewHeapableCircuitImpl([]model.CircuitVertex{
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

	clone := circuit.CloneAndUpdate().(*model.HeapableCircuitImpl)

	assert.Len(circuit.GetUnattachedVertices(), 3)
	assert.Len(circuit.GetAttachedVertices(), 5)
	assert.Len(circuit.GetAttachedEdges(), 5)
	assert.Equal(14, circuit.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Edge:     circuit.GetAttachedEdges()[2],
		Vertex:   circuit.Vertices[5],
		Distance: 1.628650237136812,
	}, circuit.GetClosestEdges().Peek())

	assert.Len(clone.GetUnattachedVertices(), 2)
	assert.Len(clone.GetAttachedVertices(), 6)
	assert.Len(clone.GetAttachedEdges(), 6)
	assert.Equal(12, clone.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Edge:     clone.GetAttachedEdges()[1],
		Vertex:   clone.Vertices[3],
		Distance: 5.09082042374693,
	}, clone.GetClosestEdges().Peek())

	assert.InDelta(circuit.GetLength()+0.763503994948632, clone.GetLength(), model.Threshold)

	// Validate that cloning a clone does not affect the original circuit
	cloneOfClone := clone.CloneAndUpdate().(*model.HeapableCircuitImpl)

	assert.Len(circuit.GetUnattachedVertices(), 3)
	assert.Len(circuit.GetAttachedVertices(), 5)
	assert.Len(circuit.GetAttachedEdges(), 5)
	assert.Equal(14, circuit.GetClosestEdges().Len())

	assert.Equal(&model.DistanceToEdge{
		Edge:     circuit.GetAttachedEdges()[2],
		Vertex:   circuit.Vertices[5],
		Distance: 1.628650237136812,
	}, circuit.GetClosestEdges().Peek())

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

	// Validate that cloning a circuit with only one vertex left to attach just updates that object and doesn't create a clone
	assert.Nil(cloneOfClone.CloneAndUpdate())

	assert.Len(cloneOfClone.GetUnattachedVertices(), 0)
	assert.Len(cloneOfClone.GetAttachedVertices(), 8)
	assert.Len(cloneOfClone.GetAttachedEdges(), 8)
	assert.Equal(0, cloneOfClone.GetClosestEdges().Len())

	// Validate that cloning a completed circuit makes no changes
	assert.Nil(cloneOfClone.CloneAndUpdate())

	assert.Len(cloneOfClone.GetUnattachedVertices(), 0)
	assert.Len(cloneOfClone.GetAttachedVertices(), 8)
	assert.Len(cloneOfClone.GetAttachedEdges(), 8)
	assert.Equal(0, cloneOfClone.GetClosestEdges().Len())
}

func TestDelete_Heap(t *testing.T) {
	assert := assert.New(t)

	circuit := model.NewHeapableCircuitImpl([]model.CircuitVertex{
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

	clone := circuit.CloneAndUpdate().(*model.HeapableCircuitImpl)
	cloneOfClone := clone.CloneAndUpdate().(*model.HeapableCircuitImpl)

	circuit.Delete()
	assert.Len(circuit.GetUnattachedVertices(), 0)
	assert.Nil(circuit.GetAttachedEdges())
	assert.Nil(circuit.GetAttachedVertices())
	assert.Nil(circuit.GetClosestEdges())
	assert.Nil(circuit.Vertices)

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
	assert.Nil(cloneOfClone.GetAttachedEdges())
	assert.Nil(cloneOfClone.GetAttachedVertices())
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

func TestPrepare_Heap(t *testing.T) {
	assert := assert.New(t)
	circuit := model.NewHeapableCircuitImpl([]model.CircuitVertex{
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

	assert.Equal(0.0, circuit.GetLength())
	assert.Equal(0.0, circuit.GetLengthWithNext())

	assert.NotNil(circuit.GetClosestEdges())
	assert.Equal(0, circuit.GetClosestEdges().Len())

	assert.NotNil(circuit.GetAttachedVertices())
	assert.Len(circuit.GetAttachedVertices(), 0)

	assert.NotNil(circuit.GetAttachedEdges())
	assert.Len(circuit.GetAttachedEdges(), 0)
}
