package model2d_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerimeter(t *testing.T) {
	assert := assert.New(t)
	vertices := []model.CircuitVertex{
		// Note: the circuit is sorted by DeduplicateVertices(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	}

	vertices = model2d.DeduplicateVertices(vertices)

	circuitEdges, unattachedVertices := model2d.BuildPerimiter(vertices)

	assert.Len(vertices, 8)

	assert.Len(circuitEdges, 5)
	assert.True(model2d.NewEdge2D(vertices[0].(*model2d.Vertex2D), vertices[7].(*model2d.Vertex2D)).Equals(circuitEdges[0]))
	assert.True(model2d.NewEdge2D(vertices[7].(*model2d.Vertex2D), vertices[6].(*model2d.Vertex2D)).Equals(circuitEdges[1]))
	assert.True(model2d.NewEdge2D(vertices[6].(*model2d.Vertex2D), vertices[4].(*model2d.Vertex2D)).Equals(circuitEdges[2]))
	assert.True(model2d.NewEdge2D(vertices[4].(*model2d.Vertex2D), vertices[1].(*model2d.Vertex2D)).Equals(circuitEdges[3]))
	assert.True(model2d.NewEdge2D(vertices[1].(*model2d.Vertex2D), vertices[0].(*model2d.Vertex2D)).Equals(circuitEdges[4]))

	assert.Len(unattachedVertices, 3)
	assert.True(unattachedVertices[vertices[2]])
	assert.True(unattachedVertices[vertices[3]])
	assert.True(unattachedVertices[vertices[5]])
}
