package model_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerimeter_GreedyWithUpdates(t *testing.T) {
	assert := assert.New(t)
	circuit := model.NewCircuitGreedyWithUpdatesImpl([]model.CircuitVertex{
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
	).(*model.CircuitGreedyWithUpdatesImpl)

	circuit.Prepare()
	circuit.BuildPerimiter()

	assert.Len(circuit.Vertices, 8)

	assert.Len(circuit.GetAttachedVertices(), 5)
	assert.Equal(model2d.NewVertex2D(-15, -15), circuit.GetAttachedVertices()[0])
	assert.Equal(model2d.NewVertex2D(15, -15), circuit.GetAttachedVertices()[1])
	assert.Equal(model2d.NewVertex2D(9, 6), circuit.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(3, 13), circuit.GetAttachedVertices()[3])
	assert.Equal(model2d.NewVertex2D(-7, 6), circuit.GetAttachedVertices()[4])

	assert.InDelta(95.73863479511238, circuit.GetLength(), model.Threshold)

	assert.Len(circuit.GetAttachedEdges(), 5)
	assert.Equal(circuit.Vertices[0].EdgeTo(circuit.Vertices[7]), circuit.GetAttachedEdges()[0])
	assert.Equal(circuit.Vertices[7].EdgeTo(circuit.Vertices[6]), circuit.GetAttachedEdges()[1])
	assert.Equal(circuit.Vertices[6].EdgeTo(circuit.Vertices[4]), circuit.GetAttachedEdges()[2])
	assert.Equal(circuit.Vertices[4].EdgeTo(circuit.Vertices[1]), circuit.GetAttachedEdges()[3])
	assert.Equal(circuit.Vertices[1].EdgeTo(circuit.Vertices[0]), circuit.GetAttachedEdges()[4])

	assert.Len(circuit.GetInteriorVertices(), 3)
	assert.True(circuit.GetInteriorVertices()[circuit.Vertices[2]])
	assert.True(circuit.GetInteriorVertices()[circuit.Vertices[3]])
	assert.True(circuit.GetInteriorVertices()[circuit.Vertices[5]])

	assert.Len(circuit.GetUnattachedVertices(), 3)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[3]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[5]])

	assert.Equal(3, circuit.GetClosestEdges().Len())
	assert.Equal(&model.DistanceToEdge{
		Edge:     circuit.GetAttachedEdges()[1],
		Distance: 0.763503994948632,
		Vertex:   circuit.Vertices[5],
	}, circuit.GetClosestEdges().PopHeap())
	assert.Equal(&model.DistanceToEdge{
		Edge:     circuit.GetAttachedEdges()[1],
		Distance: 5.854324418695558,
		Vertex:   circuit.Vertices[3],
	}, circuit.GetClosestEdges().PopHeap())
	assert.Equal(&model.DistanceToEdge{
		Edge:     circuit.GetAttachedEdges()[4],
		Distance: 7.9605428386450825,
		Vertex:   circuit.Vertices[2],
	}, circuit.GetClosestEdges().PopHeap())
}

func TestPrepare_GreedyWithUpdates(t *testing.T) {
	assert := assert.New(t)
	circuit := model.NewCircuitGreedyWithUpdatesImpl([]model.CircuitVertex{
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
	).(*model.CircuitGreedyWithUpdatesImpl)

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

	assert.NotNil(circuit.GetInteriorVertices())
	assert.Len(circuit.GetInteriorVertices(), 0)

	assert.NotNil(circuit.GetClosestEdges())
	assert.Equal(0, circuit.GetClosestEdges().Len())

	assert.NotNil(circuit.GetAttachedVertices())
	assert.Len(circuit.GetAttachedVertices(), 0)

	assert.NotNil(circuit.GetAttachedEdges())
	assert.Len(circuit.GetAttachedEdges(), 0)
}

func TestUpdate_GreedyWithUpdates(t *testing.T) {
	assert := assert.New(t)
	circuit := model.NewCircuitGreedyWithUpdatesImpl([]model.CircuitVertex{
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
	).(*model.CircuitGreedyWithUpdatesImpl)

	circuit.Prepare()
	circuit.BuildPerimiter()

	assert.Len(circuit.Vertices, 8)
	assert.Len(circuit.GetAttachedVertices(), 5)
	assert.Len(circuit.GetAttachedEdges(), 5)
	assert.Len(circuit.GetInteriorVertices(), 3)
	assert.Len(circuit.GetUnattachedVertices(), 3)
	assert.Equal(3, circuit.GetClosestEdges().Len())

	circuit.Update(circuit.FindNextVertexAndEdge())

	assert.Len(circuit.Vertices, 8)
	assert.Len(circuit.GetAttachedVertices(), 6)
	assert.Len(circuit.GetAttachedEdges(), 6)
	assert.Len(circuit.GetInteriorVertices(), 3)
	assert.Len(circuit.GetUnattachedVertices(), 2)
	assert.Equal(2, circuit.GetClosestEdges().Len())

	assert.Equal(circuit.Vertices[0].EdgeTo(circuit.Vertices[7]), circuit.GetAttachedEdges()[0])
	assert.Equal(circuit.Vertices[7].EdgeTo(circuit.Vertices[5]), circuit.GetAttachedEdges()[1])
	assert.Equal(circuit.Vertices[5].EdgeTo(circuit.Vertices[6]), circuit.GetAttachedEdges()[2])
	assert.Equal(circuit.Vertices[6].EdgeTo(circuit.Vertices[4]), circuit.GetAttachedEdges()[3])
	assert.Equal(circuit.Vertices[4].EdgeTo(circuit.Vertices[1]), circuit.GetAttachedEdges()[4])
	assert.Equal(circuit.Vertices[1].EdgeTo(circuit.Vertices[0]), circuit.GetAttachedEdges()[5])

	assert.Equal(model2d.NewVertex2D(8, 5), circuit.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(9, 6), circuit.GetAttachedVertices()[3])

	circuit.Update(circuit.FindNextVertexAndEdge())

	assert.Len(circuit.Vertices, 8)
	assert.Len(circuit.GetAttachedVertices(), 7)
	assert.Len(circuit.GetAttachedEdges(), 7)
	assert.Len(circuit.GetInteriorVertices(), 3)
	assert.Len(circuit.GetUnattachedVertices(), 1)
	assert.Equal(1, circuit.GetClosestEdges().Len())

	assert.Equal(circuit.Vertices[0].EdgeTo(circuit.Vertices[7]), circuit.GetAttachedEdges()[0])
	assert.Equal(circuit.Vertices[7].EdgeTo(circuit.Vertices[3]), circuit.GetAttachedEdges()[1])
	assert.Equal(circuit.Vertices[3].EdgeTo(circuit.Vertices[5]), circuit.GetAttachedEdges()[2])
	assert.Equal(circuit.Vertices[5].EdgeTo(circuit.Vertices[6]), circuit.GetAttachedEdges()[3])
	assert.Equal(circuit.Vertices[6].EdgeTo(circuit.Vertices[4]), circuit.GetAttachedEdges()[4])
	assert.Equal(circuit.Vertices[4].EdgeTo(circuit.Vertices[1]), circuit.GetAttachedEdges()[5])
	assert.Equal(circuit.Vertices[1].EdgeTo(circuit.Vertices[0]), circuit.GetAttachedEdges()[6])

	assert.Equal(model2d.NewVertex2D(3, 0), circuit.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(8, 5), circuit.GetAttachedVertices()[3])
	assert.Equal(model2d.NewVertex2D(9, 6), circuit.GetAttachedVertices()[4])

	circuit.Update(circuit.FindNextVertexAndEdge())

	assert.Len(circuit.Vertices, 8)
	assert.Len(circuit.GetAttachedVertices(), 8)
	assert.Len(circuit.GetAttachedEdges(), 8)
	assert.Len(circuit.GetInteriorVertices(), 3)
	assert.Len(circuit.GetUnattachedVertices(), 0)
	assert.Equal(0, circuit.GetClosestEdges().Len())

	assert.Equal(circuit.Vertices[0].EdgeTo(circuit.Vertices[7]), circuit.GetAttachedEdges()[0])
	assert.Equal(circuit.Vertices[7].EdgeTo(circuit.Vertices[2]), circuit.GetAttachedEdges()[1])
	assert.Equal(circuit.Vertices[2].EdgeTo(circuit.Vertices[3]), circuit.GetAttachedEdges()[2])
	assert.Equal(circuit.Vertices[3].EdgeTo(circuit.Vertices[5]), circuit.GetAttachedEdges()[3])
	assert.Equal(circuit.Vertices[5].EdgeTo(circuit.Vertices[6]), circuit.GetAttachedEdges()[4])
	assert.Equal(circuit.Vertices[6].EdgeTo(circuit.Vertices[4]), circuit.GetAttachedEdges()[5])
	assert.Equal(circuit.Vertices[4].EdgeTo(circuit.Vertices[1]), circuit.GetAttachedEdges()[6])
	assert.Equal(circuit.Vertices[1].EdgeTo(circuit.Vertices[0]), circuit.GetAttachedEdges()[7])

	assert.Equal(model2d.NewVertex2D(0, 0), circuit.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(3, 0), circuit.GetAttachedVertices()[3])

	circuit.Update(circuit.FindNextVertexAndEdge())

	assert.Len(circuit.Vertices, 8)
	assert.Len(circuit.GetAttachedVertices(), 8)
	assert.Len(circuit.GetAttachedEdges(), 8)
	assert.Len(circuit.GetInteriorVertices(), 3)
	assert.Len(circuit.GetUnattachedVertices(), 0)
	assert.Equal(0, circuit.GetClosestEdges().Len())
}

func TestUpdate_ShouldRemoveAttachedInteriorPointFromPerimeterIfNewEdgeIsCloserThanPreviousEdge_GreedyWithUpdates(t *testing.T) {
	assert := assert.New(t)
	circuit := model.NewCircuitGreedyWithUpdatesImpl([]model.CircuitVertex{
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(4.7, 2.0),
		model2d.NewVertex2D(5.0, 2.25),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6.0, 2.5),
		model2d.NewVertex2D(10, 0),
	},
		model2d.DeduplicateVertices,
		&model2d.PerimeterBuilder2D{},
	).(*model.CircuitGreedyWithUpdatesImpl)

	circuit.Prepare()
	circuit.BuildPerimiter()

	circuit.Update(circuit.FindNextVertexAndEdge())
	assert.Len(circuit.GetAttachedVertices(), 4)
	assert.Equal(circuit.Vertices[0], circuit.GetAttachedVertices()[0])
	assert.Equal(circuit.Vertices[5], circuit.GetAttachedVertices()[1])
	assert.Equal(circuit.Vertices[4], circuit.GetAttachedVertices()[2])
	assert.Equal(circuit.Vertices[3], circuit.GetAttachedVertices()[3])

	circuit.Update(circuit.FindNextVertexAndEdge())
	assert.Len(circuit.GetAttachedVertices(), 5)
	assert.Equal(circuit.Vertices[0], circuit.GetAttachedVertices()[0])
	assert.Equal(circuit.Vertices[1], circuit.GetAttachedVertices()[1])
	assert.Equal(circuit.Vertices[5], circuit.GetAttachedVertices()[2])
	assert.Equal(circuit.Vertices[4], circuit.GetAttachedVertices()[3])
	assert.Equal(circuit.Vertices[3], circuit.GetAttachedVertices()[4])

	circuit.Update(circuit.FindNextVertexAndEdge())
	assert.Len(circuit.GetAttachedVertices(), 5)
	assert.Equal(circuit.Vertices[0], circuit.GetAttachedVertices()[0])
	assert.Equal(circuit.Vertices[1], circuit.GetAttachedVertices()[1])
	assert.Equal(circuit.Vertices[2], circuit.GetAttachedVertices()[2])
	assert.Equal(circuit.Vertices[5], circuit.GetAttachedVertices()[3])
	assert.Equal(circuit.Vertices[3], circuit.GetAttachedVertices()[4])

	circuit.Update(circuit.FindNextVertexAndEdge())
	assert.Len(circuit.GetAttachedVertices(), 6)
	assert.Equal(circuit.Vertices[0], circuit.GetAttachedVertices()[0])
	assert.Equal(circuit.Vertices[1], circuit.GetAttachedVertices()[1])
	assert.Equal(circuit.Vertices[2], circuit.GetAttachedVertices()[2])
	assert.Equal(circuit.Vertices[4], circuit.GetAttachedVertices()[3])
	assert.Equal(circuit.Vertices[5], circuit.GetAttachedVertices()[4])
	assert.Equal(circuit.Vertices[3], circuit.GetAttachedVertices()[5])

	v, _ := circuit.FindNextVertexAndEdge()
	assert.Nil(v)
}
