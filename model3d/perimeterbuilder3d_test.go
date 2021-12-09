package model3d_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model3d"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerimeter_ShouldProduceSameResultAs2D(t *testing.T) {
	assert := assert.New(t)
	pb := &model3d.PerimeterBuilder3D{}
	vertices := []model.CircuitVertex{
		// Note: 3d circuits are not sorted.
		model3d.NewVertex3D(-15, -15, 0),
		model3d.NewVertex3D(0, 0, 0),
		model3d.NewVertex3D(15, -15, 0),
		model3d.NewVertex3D(3, 0, 0),
		model3d.NewVertex3D(3, 13, 0),
		model3d.NewVertex3D(8, 5, 0),
		model3d.NewVertex3D(9, 6, 0),
		model3d.NewVertex3D(-7, 6, 0),
	}

	vertices = model3d.DeduplicateVertices3D(vertices)

	circuitEdges, unattachedVertices := pb.BuildPerimiter(vertices)

	assert.Len(vertices, 8)

	assert.Len(circuitEdges, 5)
	assert.Equal(model3d.NewEdge3D(vertices[0].(*model3d.Vertex3D), vertices[2].(*model3d.Vertex3D)), circuitEdges[0])
	assert.Equal(model3d.NewEdge3D(vertices[2].(*model3d.Vertex3D), vertices[6].(*model3d.Vertex3D)), circuitEdges[1])
	assert.Equal(model3d.NewEdge3D(vertices[6].(*model3d.Vertex3D), vertices[4].(*model3d.Vertex3D)), circuitEdges[2])
	assert.Equal(model3d.NewEdge3D(vertices[4].(*model3d.Vertex3D), vertices[7].(*model3d.Vertex3D)), circuitEdges[3])
	assert.Equal(model3d.NewEdge3D(vertices[7].(*model3d.Vertex3D), vertices[0].(*model3d.Vertex3D)), circuitEdges[4])

	assert.Len(unattachedVertices, 3)
	assert.True(unattachedVertices[vertices[1]])
	assert.True(unattachedVertices[vertices[3]])
	assert.True(unattachedVertices[vertices[5]])
}

func TestBuildPerimeter3D(t *testing.T) {
	assert := assert.New(t)
	pb := &model3d.PerimeterBuilder3D{}
	vertices := []model.CircuitVertex{
		// Note: 3d circuits are not sorted.
		model3d.NewVertex3D(-15, -15, 0),
		model3d.NewVertex3D(0, 0, 13),
		model3d.NewVertex3D(15, -15, 5),
		model3d.NewVertex3D(3, 0, -3),
		model3d.NewVertex3D(3, 13, -16),
		model3d.NewVertex3D(8, 5, 4),
		model3d.NewVertex3D(9, 6, -10),
		model3d.NewVertex3D(-7, 6, -8),
	}

	circuitEdges, unattachedVertices := pb.BuildPerimiter(vertices)

	assert.Len(vertices, 8)

	assert.Len(circuitEdges, 6)
	assert.Equal(model3d.NewEdge3D(vertices[0].(*model3d.Vertex3D), vertices[2].(*model3d.Vertex3D)), circuitEdges[0])
	assert.Equal(model3d.NewEdge3D(vertices[2].(*model3d.Vertex3D), vertices[5].(*model3d.Vertex3D)), circuitEdges[1])

	// (8, 5, 4) may be an issue - if it is, here are 3 options to test:
	// 1) When a point is added check adjacent points to see if they could now be interior relative to the edge that was just created,
	//    i.e. after adding (9, 6, -10) does (8, 5, 4) become closer to middle than its projection to the edge between (15, -15, 5) and (9, 6, -10)?
	// 2) Can check angles using cross product to see if edge would become concave (will not help this case, may be required for other cases)
	// 3) Use PCA (principle component analysis) to create custom 3-D axis based on dimensions with the most variance.
	//    Project to 2 dimensions (track 2d-to-3d vertices), find perimeter using 2D approach, and convert back to 3D.
	assert.Equal(model3d.NewEdge3D(vertices[5].(*model3d.Vertex3D), vertices[6].(*model3d.Vertex3D)), circuitEdges[2])

	assert.Equal(model3d.NewEdge3D(vertices[6].(*model3d.Vertex3D), vertices[4].(*model3d.Vertex3D)), circuitEdges[3])
	assert.Equal(model3d.NewEdge3D(vertices[4].(*model3d.Vertex3D), vertices[7].(*model3d.Vertex3D)), circuitEdges[4])
	assert.Equal(model3d.NewEdge3D(vertices[7].(*model3d.Vertex3D), vertices[0].(*model3d.Vertex3D)), circuitEdges[5])

	assert.Len(unattachedVertices, 2)
	assert.True(unattachedVertices[vertices[1]])
	assert.True(unattachedVertices[vertices[3]])
}
