package model2d

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerimeter(t *testing.T) {
	assert := assert.New(t)
	circuit := &Circuit2D{
		Vertices: []*Vertex2D{
			// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
			NewVertex2D(-15, -15), // Index 0 after sorting
			NewVertex2D(0, 0),     // Index 2 after sorting
			NewVertex2D(15, -15),  // Index 7 after sorting
			NewVertex2D(3, 0),     // Index 3 after sorting
			NewVertex2D(3, 13),    // Index 4 after sorting
			NewVertex2D(8, 5),     // Index 5 after sorting
			NewVertex2D(9, 6),     // Index 6 after sorting
			NewVertex2D(-7, 6),    // Index 1 after sorting
		},
	}
	circuit.Prepare()
	circuit.BuildPerimiter()

	assert.Len(circuit.Vertices, 8)

	assert.Len(circuit.circuit, 5)
	assert.Equal(NewVertex2D(-15, -15), circuit.circuit[0])
	assert.Equal(NewVertex2D(15, -15), circuit.circuit[1])
	assert.Equal(NewVertex2D(9, 6), circuit.circuit[2])
	assert.Equal(NewVertex2D(3, 13), circuit.circuit[3])
	assert.Equal(NewVertex2D(-7, 6), circuit.circuit[4])
	assert.Equal(circuit.circuit, circuit.GetAttachedVertices())

	assert.Len(circuit.circuitEdges, 5)
	assert.Equal(NewEdge2D(circuit.Vertices[0], circuit.Vertices[7]), circuit.circuitEdges[0])
	assert.Equal(NewEdge2D(circuit.Vertices[7], circuit.Vertices[6]), circuit.circuitEdges[1])
	assert.Equal(NewEdge2D(circuit.Vertices[6], circuit.Vertices[4]), circuit.circuitEdges[2])
	assert.Equal(NewEdge2D(circuit.Vertices[4], circuit.Vertices[1]), circuit.circuitEdges[3])
	assert.Equal(NewEdge2D(circuit.Vertices[1], circuit.Vertices[0]), circuit.circuitEdges[4])

	assert.Len(circuit.interiorVertices, 3)
	assert.True(circuit.interiorVertices[circuit.Vertices[2]])
	assert.True(circuit.interiorVertices[circuit.Vertices[3]])
	assert.True(circuit.interiorVertices[circuit.Vertices[5]])
	assert.Equal(circuit.interiorVertices, circuit.GetInteriorVertices())

	assert.Len(circuit.unattachedVertices, 3)
	assert.True(circuit.unattachedVertices[circuit.Vertices[2]])
	assert.True(circuit.unattachedVertices[circuit.Vertices[3]])
	assert.True(circuit.unattachedVertices[circuit.Vertices[5]])
	assert.Equal(circuit.unattachedVertices, circuit.GetUnattachedVertices())

	assert.Len(circuit.closestEdges, 3)
	assert.Equal(&distanceToEdge{
		edge:     circuit.circuitEdges[4],
		distance: 7.9605428386450825,
	}, circuit.closestEdges[circuit.Vertices[2]])
	assert.Equal(&distanceToEdge{
		edge:     circuit.circuitEdges[1],
		distance: 5.854324418695558,
	}, circuit.closestEdges[circuit.Vertices[3]])
	assert.Equal(&distanceToEdge{
		edge:     circuit.circuitEdges[1],
		distance: 0.763503994948632,
	}, circuit.closestEdges[circuit.Vertices[5]])
}

func TestFindNextVertexAndEdge(t *testing.T) {
	assert := assert.New(t)
	circuit := &Circuit2D{
		Vertices: []*Vertex2D{
			NewVertex2D(-15, -15),
			NewVertex2D(0, 0),
			NewVertex2D(15, -15),
		},
	}

	circuit.Prepare()
	circuit.BuildPerimiter()

	v, e := circuit.FindNextVertexAndEdge()
	assert.Nil(v)
	assert.Nil(e)

	v1 := NewVertex2D(1, 1)
	circuit.closestEdges[v1] = &distanceToEdge{
		distance: 1.2345,
		edge:     circuit.circuitEdges[2],
	}

	// Verify that the vertex must be in unattached vertices, and is ignored if it is only in closest edges.
	v, e = circuit.FindNextVertexAndEdge()
	assert.Nil(v)
	assert.Nil(e)

	circuit.unattachedVertices[v1] = true
	v, e = circuit.FindNextVertexAndEdge()
	assert.Equal(v1, v)
	assert.Equal(circuit.circuitEdges[2], e)

	// v1 is closest currently, so it should be returned
	v2 := NewVertex2D(-2, 1)
	circuit.unattachedVertices[v2] = true
	circuit.closestEdges[v2] = &distanceToEdge{
		distance: 2.3456,
		edge:     circuit.circuitEdges[1],
	}
	v, e = circuit.FindNextVertexAndEdge()
	assert.Equal(v1, v)
	assert.Equal(circuit.circuitEdges[2], e)

	// update v2 to be closer, so it should be returned
	circuit.closestEdges[v2] = &distanceToEdge{
		distance: 0.3456,
		edge:     circuit.circuitEdges[0],
	}
	v, e = circuit.FindNextVertexAndEdge()
	assert.Equal(v2, v)
	assert.Equal(circuit.circuitEdges[0], e)
}

func TestPrepare(t *testing.T) {
	assert := assert.New(t)
	circuit := &Circuit2D{
		Vertices: []*Vertex2D{
			NewVertex2D(-15, -15),
			NewVertex2D(0, 0),
			NewVertex2D(15, -15),
			NewVertex2D(-15-model.Threshold/3.0, -15),
			NewVertex2D(0, 0),
			NewVertex2D(-15+model.Threshold/3.0, -15-model.Threshold/3.0),
			NewVertex2D(3, 0),
			NewVertex2D(3, 13),
			NewVertex2D(7, 6),
			NewVertex2D(-7, 6),
		},
	}

	circuit.Prepare()

	assert.NotNil(circuit.Vertices)
	assert.Len(circuit.Vertices, 7)
	assert.ElementsMatch(circuit.Vertices, []*Vertex2D{
		NewVertex2D(-15+model.Threshold/3.0, -15-model.Threshold/3.0),
		NewVertex2D(-7, 6),
		NewVertex2D(0, 0),
		NewVertex2D(3, 0),
		NewVertex2D(3, 13),
		NewVertex2D(7, 6),
		NewVertex2D(15, -15),
	})

	assert.NotNil(circuit.unattachedVertices)
	assert.Len(circuit.unattachedVertices, 7)

	assert.NotNil(circuit.interiorVertices)
	assert.Len(circuit.interiorVertices, 0)

	assert.NotNil(circuit.closestEdges)
	assert.Len(circuit.closestEdges, 0)

	assert.NotNil(circuit.circuit)
	assert.Len(circuit.circuit, 0)

	assert.NotNil(circuit.circuitEdges)
	assert.Len(circuit.circuitEdges, 0)

	assert.NotNil(circuit.midpoint)
	assert.InDelta((6.0+model.Threshold/3.0)/7.0, circuit.midpoint.X, model.Threshold)
	assert.InDelta(-5.0/7.0, circuit.midpoint.Y, model.Threshold)
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)
	circuit := &Circuit2D{
		Vertices: []*Vertex2D{
			// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
			NewVertex2D(-15, -15), // Index 0 after sorting
			NewVertex2D(0, 0),     // Index 2 after sorting
			NewVertex2D(15, -15),  // Index 7 after sorting
			NewVertex2D(3, 0),     // Index 3 after sorting
			NewVertex2D(3, 13),    // Index 4 after sorting
			NewVertex2D(8, 5),     // Index 5 after sorting
			NewVertex2D(9, 6),     // Index 6 after sorting
			NewVertex2D(-7, 6),    // Index 1 after sorting
		},
	}
	circuit.Prepare()
	circuit.BuildPerimiter()

	assert.Len(circuit.Vertices, 8)
	assert.Len(circuit.circuit, 5)
	assert.Len(circuit.circuitEdges, 5)
	assert.Len(circuit.interiorVertices, 3)
	assert.Len(circuit.closestEdges, 3)
	assert.Len(circuit.unattachedVertices, 3)

	circuit.Update(circuit.FindNextVertexAndEdge())

	assert.Len(circuit.Vertices, 8)
	assert.Len(circuit.circuit, 6)
	assert.Len(circuit.circuitEdges, 6)
	assert.Len(circuit.interiorVertices, 3)
	assert.Len(circuit.closestEdges, 3)
	assert.Len(circuit.unattachedVertices, 2)

	assert.Equal(NewEdge2D(circuit.Vertices[0], circuit.Vertices[7]), circuit.circuitEdges[0])
	assert.Equal(NewEdge2D(circuit.Vertices[7], circuit.Vertices[5]), circuit.circuitEdges[1])
	assert.Equal(NewEdge2D(circuit.Vertices[5], circuit.Vertices[6]), circuit.circuitEdges[2])
	assert.Equal(NewEdge2D(circuit.Vertices[6], circuit.Vertices[4]), circuit.circuitEdges[3])
	assert.Equal(NewEdge2D(circuit.Vertices[4], circuit.Vertices[1]), circuit.circuitEdges[4])
	assert.Equal(NewEdge2D(circuit.Vertices[1], circuit.Vertices[0]), circuit.circuitEdges[5])

	assert.Equal(NewVertex2D(8, 5), circuit.circuit[2])
	assert.Equal(NewVertex2D(9, 6), circuit.circuit[3])

	circuit.Update(circuit.FindNextVertexAndEdge())

	assert.Len(circuit.Vertices, 8)
	assert.Len(circuit.circuit, 7)
	assert.Len(circuit.circuitEdges, 7)
	assert.Len(circuit.interiorVertices, 3)
	assert.Len(circuit.closestEdges, 3)
	assert.Len(circuit.unattachedVertices, 1)

	assert.Equal(NewEdge2D(circuit.Vertices[0], circuit.Vertices[7]), circuit.circuitEdges[0])
	assert.Equal(NewEdge2D(circuit.Vertices[7], circuit.Vertices[3]), circuit.circuitEdges[1])
	assert.Equal(NewEdge2D(circuit.Vertices[3], circuit.Vertices[5]), circuit.circuitEdges[2])
	assert.Equal(NewEdge2D(circuit.Vertices[5], circuit.Vertices[6]), circuit.circuitEdges[3])
	assert.Equal(NewEdge2D(circuit.Vertices[6], circuit.Vertices[4]), circuit.circuitEdges[4])
	assert.Equal(NewEdge2D(circuit.Vertices[4], circuit.Vertices[1]), circuit.circuitEdges[5])
	assert.Equal(NewEdge2D(circuit.Vertices[1], circuit.Vertices[0]), circuit.circuitEdges[6])

	assert.Equal(NewVertex2D(3, 0), circuit.circuit[2])
	assert.Equal(NewVertex2D(8, 5), circuit.circuit[3])
	assert.Equal(NewVertex2D(9, 6), circuit.circuit[4])

	circuit.Update(circuit.FindNextVertexAndEdge())

	assert.Len(circuit.Vertices, 8)
	assert.Len(circuit.circuit, 8)
	assert.Len(circuit.circuitEdges, 8)
	assert.Len(circuit.interiorVertices, 3)
	assert.Len(circuit.closestEdges, 3)
	assert.Len(circuit.unattachedVertices, 0)

	assert.Equal(NewEdge2D(circuit.Vertices[0], circuit.Vertices[7]), circuit.circuitEdges[0])
	assert.Equal(NewEdge2D(circuit.Vertices[7], circuit.Vertices[2]), circuit.circuitEdges[1])
	assert.Equal(NewEdge2D(circuit.Vertices[2], circuit.Vertices[3]), circuit.circuitEdges[2])
	assert.Equal(NewEdge2D(circuit.Vertices[3], circuit.Vertices[5]), circuit.circuitEdges[3])
	assert.Equal(NewEdge2D(circuit.Vertices[5], circuit.Vertices[6]), circuit.circuitEdges[4])
	assert.Equal(NewEdge2D(circuit.Vertices[6], circuit.Vertices[4]), circuit.circuitEdges[5])
	assert.Equal(NewEdge2D(circuit.Vertices[4], circuit.Vertices[1]), circuit.circuitEdges[6])
	assert.Equal(NewEdge2D(circuit.Vertices[1], circuit.Vertices[0]), circuit.circuitEdges[7])

	assert.Equal(NewVertex2D(0, 0), circuit.circuit[2])
	assert.Equal(NewVertex2D(3, 0), circuit.circuit[3])

	circuit.Update(circuit.FindNextVertexAndEdge())

	assert.Len(circuit.Vertices, 8)
	assert.Len(circuit.circuit, 8)
	assert.Len(circuit.circuitEdges, 8)
	assert.Len(circuit.interiorVertices, 3)
	assert.Len(circuit.closestEdges, 3)
	assert.Len(circuit.unattachedVertices, 0)
}

func TestUpdate_ShouldRemoveAttachedInteriorPointFromPerimeterIfNewEdgeIsCloserThanPreviousEdge(t *testing.T) {
	assert := assert.New(t)
	circuit := &Circuit2D{
		Vertices: []*Vertex2D{
			NewVertex2D(0, 0),
			NewVertex2D(4.7, 2.0),
			NewVertex2D(5.0, 2.25),
			NewVertex2D(5, 5),
			NewVertex2D(6.0, 2.5),
			NewVertex2D(10, 0),
		},
	}

	circuit.Prepare()
	circuit.BuildPerimiter()

	circuit.Update(circuit.FindNextVertexAndEdge())
	assert.Len(circuit.circuit, 4)
	assert.Equal(circuit.Vertices[4], circuit.circuit[2])

	circuit.Update(circuit.FindNextVertexAndEdge())
	assert.Len(circuit.circuit, 5)
	assert.Equal(circuit.Vertices[1], circuit.circuit[1])
	assert.Equal(circuit.Vertices[4], circuit.circuit[3])
	assert.NotNil(circuit.closestEdges[circuit.Vertices[1]])
	assert.Equal(NewEdge2D(circuit.Vertices[0], circuit.Vertices[5]), circuit.closestEdges[circuit.Vertices[1]].edge)
	assert.NotNil(circuit.closestEdges[circuit.Vertices[2]])
	assert.Equal(NewEdge2D(circuit.Vertices[1], circuit.Vertices[5]), circuit.closestEdges[circuit.Vertices[2]].edge)
	assert.NotNil(circuit.closestEdges[circuit.Vertices[4]])
	assert.Equal(NewEdge2D(circuit.Vertices[5], circuit.Vertices[3]), circuit.closestEdges[circuit.Vertices[4]].edge)

	circuit.Update(circuit.FindNextVertexAndEdge())
	assert.Len(circuit.circuit, 5)
	assert.Equal(circuit.Vertices[1], circuit.circuit[1])
	assert.Equal(circuit.Vertices[2], circuit.circuit[2])
	assert.NotNil(circuit.closestEdges[circuit.Vertices[1]])
	assert.Equal(NewEdge2D(circuit.Vertices[0], circuit.Vertices[2]), circuit.closestEdges[circuit.Vertices[1]].edge)
	assert.NotNil(circuit.closestEdges[circuit.Vertices[2]])
	assert.Equal(NewEdge2D(circuit.Vertices[1], circuit.Vertices[5]), circuit.closestEdges[circuit.Vertices[2]].edge)
	assert.NotNil(circuit.closestEdges[circuit.Vertices[4]])
	assert.Equal(NewEdge2D(circuit.Vertices[2], circuit.Vertices[5]), circuit.closestEdges[circuit.Vertices[4]].edge)

	circuit.Update(circuit.FindNextVertexAndEdge())
	assert.Len(circuit.circuit, 6)
	assert.Equal(circuit.Vertices[1], circuit.circuit[1])
	assert.Equal(circuit.Vertices[2], circuit.circuit[2])
	assert.Equal(circuit.Vertices[4], circuit.circuit[3])

	v, _ := circuit.FindNextVertexAndEdge()
	assert.Nil(v)
}

func TestInsertVertex(t *testing.T) {
	assert := assert.New(t)
	c := &Circuit2D{
		Vertices: []*Vertex2D{
			NewVertex2D(-15, -15),
			NewVertex2D(0, 0),
			NewVertex2D(15, -15),
		},
		circuit: []model.CircuitVertex{
			NewVertex2D(-15, -15),
			NewVertex2D(0, 0),
			NewVertex2D(15, -15),
		},
	}

	c.insertVertex(0, NewVertex2D(5, 5))
	assert.Len(c.circuit, 4)
	assert.Equal(NewVertex2D(5, 5), c.circuit[0])
	assert.Equal(NewVertex2D(-15, -15), c.circuit[1])
	assert.Equal(NewVertex2D(0, 0), c.circuit[2])
	assert.Equal(NewVertex2D(15, -15), c.circuit[3])

	c.insertVertex(4, NewVertex2D(-5, -5))
	assert.Len(c.circuit, 5)
	assert.Equal(NewVertex2D(5, 5), c.circuit[0])
	assert.Equal(NewVertex2D(-15, -15), c.circuit[1])
	assert.Equal(NewVertex2D(0, 0), c.circuit[2])
	assert.Equal(NewVertex2D(15, -15), c.circuit[3])
	assert.Equal(NewVertex2D(-5, -5), c.circuit[4])

	c.insertVertex(2, NewVertex2D(1, -5))
	assert.Len(c.circuit, 6)
	assert.Equal(NewVertex2D(5, 5), c.circuit[0])
	assert.Equal(NewVertex2D(-15, -15), c.circuit[1])
	assert.Equal(NewVertex2D(1, -5), c.circuit[2])
	assert.Equal(NewVertex2D(0, 0), c.circuit[3])
	assert.Equal(NewVertex2D(15, -15), c.circuit[4])
	assert.Equal(NewVertex2D(-5, -5), c.circuit[5])
}
