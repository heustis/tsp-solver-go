package experimental_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/circuit/experimental"
	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestNewConvexConcaveWeightedEdges(t *testing.T) {
	assert := assert.New(t)
	circuit := experimental.NewConvexConcaveWeightedEdges(model2d.DeduplicateVertices([]model.CircuitVertex{
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
	).(*experimental.ConvexConcaveWeightedEdges)

	assert.Len(circuit.Vertices, 8)

	assert.Len(circuit.GetAttachedVertices(), 5)
	assert.Equal(model2d.NewVertex2D(-15, -15), circuit.GetAttachedVertices()[0])
	assert.Equal(model2d.NewVertex2D(15, -15), circuit.GetAttachedVertices()[1])
	assert.Equal(model2d.NewVertex2D(9, 6), circuit.GetAttachedVertices()[2])
	assert.Equal(model2d.NewVertex2D(3, 13), circuit.GetAttachedVertices()[3])
	assert.Equal(model2d.NewVertex2D(-7, 6), circuit.GetAttachedVertices()[4])

	assert.InDelta(95.73863479511238, circuit.GetLength(), model.Threshold)

	edges := circuit.GetAttachedEdges()
	assert.Len(edges, 5)
	assert.Equal(circuit.Vertices[0].EdgeTo(circuit.Vertices[7]), edges[0])
	assert.Equal(circuit.Vertices[7].EdgeTo(circuit.Vertices[6]), edges[1])
	assert.Equal(circuit.Vertices[6].EdgeTo(circuit.Vertices[4]), edges[2])
	assert.Equal(circuit.Vertices[4].EdgeTo(circuit.Vertices[1]), edges[3])
	assert.Equal(circuit.Vertices[1].EdgeTo(circuit.Vertices[0]), edges[4])

	assert.Len(circuit.GetUnattachedVertices(), 3)
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[2]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[3]])
	assert.True(circuit.GetUnattachedVertices()[circuit.Vertices[5]])

	closestVertices := circuit.GetClosestVertices()
	assert.Len(closestVertices, 5)
	assert.NotNil(closestVertices[edges[0]])
	assert.Len(closestVertices[edges[0]].GetClosestVertices(), 3)
	assert.InDelta(12.081874046685233, closestVertices[edges[0]].GetDistance(), model.Threshold)

	assert.NotNil(closestVertices[edges[1]])
	assert.Len(closestVertices[edges[1]].GetClosestVertices(), 3)
	assert.InDelta(3.119024051416561, closestVertices[edges[1]].GetDistance(), model.Threshold)

	assert.NotNil(closestVertices[edges[2]])
	assert.Len(closestVertices[edges[2]].GetClosestVertices(), 3)
	assert.InDelta(5.748106026958004, closestVertices[edges[2]].GetDistance(), model.Threshold)

	assert.NotNil(closestVertices[edges[3]])
	assert.Len(closestVertices[edges[3]].GetClosestVertices(), 3)
	assert.InDelta(9.799425448261324, closestVertices[edges[3]].GetDistance(), model.Threshold)

	assert.NotNil(closestVertices[edges[4]])
	assert.Len(closestVertices[edges[4]].GetClosestVertices(), 3)
	assert.InDelta(10.015457439162253, closestVertices[edges[4]].GetDistance(), model.Threshold)
}

func TestUpdate_ConvexConcaveWeightedEdges(t *testing.T) {
	assert := assert.New(t)
	circuit := experimental.NewConvexConcaveWeightedEdges(model2d.DeduplicateVertices([]model.CircuitVertex{
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
	).(*experimental.ConvexConcaveWeightedEdges)

	assert.Len(circuit.Vertices, 8)
	assert.Len(circuit.GetAttachedVertices(), 5)
	assert.Len(circuit.GetAttachedEdges(), 5)
	assert.Len(circuit.GetUnattachedVertices(), 3)
	assert.Len(circuit.GetClosestVertices(), 5)
	for e, c := range circuit.GetClosestVertices() {
		assert.Len(c.GetClosestVertices(), 3, e)
	}

	circuit.Update(circuit.FindNextVertexAndEdge())

	assert.Len(circuit.Vertices, 8)
	assert.Len(circuit.GetAttachedVertices(), 6)
	assert.Len(circuit.GetAttachedEdges(), 6)
	assert.Len(circuit.GetUnattachedVertices(), 2)
	assert.Len(circuit.GetClosestVertices(), 6)
	for e, c := range circuit.GetClosestVertices() {
		assert.Len(c.GetClosestVertices(), 2, e)
	}

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
	assert.Len(circuit.GetUnattachedVertices(), 1)
	assert.Len(circuit.GetClosestVertices(), 7)
	for e, c := range circuit.GetClosestVertices() {
		assert.Len(c.GetClosestVertices(), 1, e)
	}

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
	assert.Len(circuit.GetUnattachedVertices(), 0)
	assert.Len(circuit.GetClosestVertices(), 8)
	for e, c := range circuit.GetClosestVertices() {
		assert.Len(c.GetClosestVertices(), 0, e)
	}

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
	assert.Len(circuit.GetUnattachedVertices(), 0)
	assert.Len(circuit.GetClosestVertices(), 8)
	for e, c := range circuit.GetClosestVertices() {
		assert.Len(c.GetClosestVertices(), 0, e)
	}
}

func TestUpdate_ShouldNotRemoveAttachedInteriorPointFromPerimeterIfNewEdgeIsCloserThanPreviousEdge_ConvexConcave_WeightedEdges(t *testing.T) {
	assert := assert.New(t)
	circuit := experimental.NewConvexConcaveWeightedEdges(model2d.DeduplicateVertices([]model.CircuitVertex{
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(4.7, 2.0),
		model2d.NewVertex2D(5.0, 2.25),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6.0, 2.5),
		model2d.NewVertex2D(10, 0),
	}),
		model2d.BuildPerimiter,
	).(*experimental.ConvexConcaveWeightedEdges)

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
	assert.Len(circuit.GetAttachedVertices(), 6)
	assert.Equal(circuit.Vertices[0], circuit.GetAttachedVertices()[0])
	assert.Equal(circuit.Vertices[1], circuit.GetAttachedVertices()[1])
	assert.Equal(circuit.Vertices[2], circuit.GetAttachedVertices()[2])
	assert.Equal(circuit.Vertices[5], circuit.GetAttachedVertices()[3])
	assert.Equal(circuit.Vertices[4], circuit.GetAttachedVertices()[4])
	assert.Equal(circuit.Vertices[3], circuit.GetAttachedVertices()[5])

	circuit.Update(circuit.FindNextVertexAndEdge())
	assert.Len(circuit.GetAttachedVertices(), 6)
	assert.Equal(circuit.Vertices[0], circuit.GetAttachedVertices()[0])
	assert.Equal(circuit.Vertices[1], circuit.GetAttachedVertices()[1])
	assert.Equal(circuit.Vertices[2], circuit.GetAttachedVertices()[2])
	assert.Equal(circuit.Vertices[5], circuit.GetAttachedVertices()[3])
	assert.Equal(circuit.Vertices[4], circuit.GetAttachedVertices()[4])
	assert.Equal(circuit.Vertices[3], circuit.GetAttachedVertices()[5])

	v, _ := circuit.FindNextVertexAndEdge()
	assert.Nil(v)
}

func TestUpdate_ConvexConcave_WeightedEdges_ShouldPanicIfEdgeIsNotInCircuit(t *testing.T) {
	assert := assert.New(t)
	circuit := experimental.NewConvexConcaveWeightedEdges(model2d.DeduplicateVertices([]model.CircuitVertex{
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
	).(*experimental.ConvexConcaveWeightedEdges)

	assert.Len(circuit.Vertices, 8)
	assert.Len(circuit.GetAttachedVertices(), 5)
	assert.Len(circuit.GetAttachedEdges(), 5)
	assert.Len(circuit.GetUnattachedVertices(), 3)
	assert.Len(circuit.GetClosestVertices(), 5)

	assert.Panics(func() { circuit.Update(circuit.Vertices[2], circuit.Vertices[0].EdgeTo(circuit.Vertices[4])) })
}
