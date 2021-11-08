package model2d

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerimeter_HeapMinClonesMove(t *testing.T) {
	assert := assert.New(t)
	circuit := &HeapableCircuit2DMinClonesMove{
		Vertices: []model.CircuitVertex{
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
	assert.Equal(NewEdge2D(circuit.Vertices[0].(*Vertex2D), circuit.Vertices[7].(*Vertex2D)), circuit.circuitEdges[0])
	assert.Equal(NewEdge2D(circuit.Vertices[7].(*Vertex2D), circuit.Vertices[6].(*Vertex2D)), circuit.circuitEdges[1])
	assert.Equal(NewEdge2D(circuit.Vertices[6].(*Vertex2D), circuit.Vertices[4].(*Vertex2D)), circuit.circuitEdges[2])
	assert.Equal(NewEdge2D(circuit.Vertices[4].(*Vertex2D), circuit.Vertices[1].(*Vertex2D)), circuit.circuitEdges[3])
	assert.Equal(NewEdge2D(circuit.Vertices[1].(*Vertex2D), circuit.Vertices[0].(*Vertex2D)), circuit.circuitEdges[4])

	expectedLength := 0.0
	for _, edge := range circuit.circuitEdges {
		expectedLength += edge.GetLength()
	}
	assert.InDelta(expectedLength, circuit.GetLength(), model.Threshold)

	assert.Len(circuit.unattachedVertices, 3)
	assert.True(circuit.unattachedVertices[circuit.Vertices[2]])
	assert.True(circuit.unattachedVertices[circuit.Vertices[3]])
	assert.True(circuit.unattachedVertices[circuit.Vertices[5]])
	assert.Equal(circuit.unattachedVertices, circuit.GetUnattachedVertices())

	assert.InDelta(95.738634795112368+0.763503994948632, circuit.GetLengthWithNext(), model.Threshold)

	assert.Equal(15, circuit.closestEdges.Len())
	assert.Equal(&heapDistanceToEdge{
		edge:     circuit.circuitEdges[1],
		vertex:   circuit.Vertices[5].(*Vertex2D),
		distance: 0.763503994948632,
	}, circuit.closestEdges.PopHeap())
	assert.Equal(&heapDistanceToEdge{
		edge:     circuit.circuitEdges[2],
		vertex:   circuit.Vertices[5].(*Vertex2D),
		distance: 1.628650237136812,
	}, circuit.closestEdges.PopHeap())
	assert.Equal(&heapDistanceToEdge{
		edge:     circuit.circuitEdges[1],
		vertex:   circuit.Vertices[3].(*Vertex2D),
		distance: 5.854324418695558,
	}, circuit.closestEdges.PopHeap())
	assert.Equal(&heapDistanceToEdge{
		edge:     circuit.circuitEdges[4],
		vertex:   circuit.Vertices[2].(*Vertex2D),
		distance: 7.9605428386450825,
	}, circuit.closestEdges.PopHeap())
}

func TestCloneAndUpdate_HeapMinClonesMove(t *testing.T) {
	assert := assert.New(t)

	circuit := &HeapableCircuit2DMinClonesMove{
		Vertices: []model.CircuitVertex{
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

	assert.True(circuit.unattachedVertices[circuit.Vertices[2]])
	assert.True(circuit.unattachedVertices[circuit.Vertices[3]])
	assert.True(circuit.unattachedVertices[circuit.Vertices[5]])

	// Index 5 should attach to edge 15,-15 -> 9,6
	assert.Nil(circuit.CloneAndUpdate())
	assert.Len(circuit.unattachedVertices, 2)
	assert.False(circuit.unattachedVertices[circuit.Vertices[5]])
	assert.Len(circuit.circuit, 6)
	assert.Len(circuit.circuitEdges, 6)
	assert.Equal(16, circuit.closestEdges.Len())

	assert.Equal(&heapDistanceToEdge{
		edge:     circuit.circuitEdges[3],
		vertex:   circuit.Vertices[5].(*Vertex2D),
		distance: 0.8651462421881799,
	}, circuit.closestEdges.Peek())

	// Index 5 should attach to edge 9,6 -> 3,13, this requires cloning since index 5 is already attached.
	clone := circuit.CloneAndUpdate().(*HeapableCircuit2DMinClonesMove)
	assert.Len(circuit.unattachedVertices, 2)
	assert.Len(clone.unattachedVertices, 2)
	assert.False(circuit.unattachedVertices[circuit.Vertices[5]])
	assert.False(clone.unattachedVertices[circuit.Vertices[5]])
	assert.Len(circuit.circuit, 6)
	assert.Len(clone.circuit, 6)
	assert.Len(circuit.circuitEdges, 6)
	assert.Len(clone.circuitEdges, 6)
	assert.Equal(15, circuit.closestEdges.Len())
	assert.Equal(15, clone.closestEdges.Len())

	assert.Equal(&heapDistanceToEdge{
		edge:     circuit.circuitEdges[1],
		vertex:   circuit.Vertices[3].(*Vertex2D),
		distance: 5.09082042374693,
	}, circuit.closestEdges.Peek())

	assert.Equal(&heapDistanceToEdge{
		edge:     clone.circuitEdges[1],
		vertex:   clone.Vertices[3].(*Vertex2D),
		distance: 5.854324418695558,
	}, clone.closestEdges.Peek())

	// Index 3 should attach to edge 1, no cloning required
	assert.Nil(circuit.CloneAndUpdate())
	assert.Nil(clone.CloneAndUpdate())
	assert.Len(circuit.unattachedVertices, 1)
	assert.Equal(15, circuit.closestEdges.Len())
	assert.Len(clone.unattachedVertices, 1)
	assert.Equal(15, clone.closestEdges.Len())

	assert.Equal(&heapDistanceToEdge{
		edge:     circuit.circuitEdges[1],
		vertex:   circuit.Vertices[2].(*Vertex2D),
		distance: 5.003830723297881,
	}, circuit.closestEdges.Peek())

	assert.Equal(&heapDistanceToEdge{
		edge:     clone.circuitEdges[4],
		vertex:   clone.Vertices[3].(*Vertex2D),
		distance: 4.782762261113314,
	}, clone.closestEdges.Peek())

	// Index 2 should attach to edge 1, no cloning required
	assert.Nil(circuit.CloneAndUpdate())
	// Index 3 should move to edge 4, cloning required
	cloneOfClone, okay := clone.CloneAndUpdate().(*HeapableCircuit2DMinClonesMove)
	assert.True(okay)

	assert.Len(circuit.unattachedVertices, 0)
	assert.Len(clone.unattachedVertices, 1)
	assert.Len(cloneOfClone.unattachedVertices, 1)

	// Index 2 should attach to edge 1, no cloning required
	assert.Equal(&heapDistanceToEdge{
		edge:     clone.circuitEdges[1],
		vertex:   clone.Vertices[2].(*Vertex2D),
		distance: 5.003830723297881,
	}, clone.closestEdges.Peek())
	assert.Nil(clone.CloneAndUpdate())

	// Index 3 should move to edge 5, cloning required
	assert.Equal(&heapDistanceToEdge{
		edge:     cloneOfClone.circuitEdges[5],
		vertex:   cloneOfClone.Vertices[3].(*Vertex2D),
		distance: 1.818261494148027,
	}, cloneOfClone.closestEdges.Peek())
	assert.NotNil(cloneOfClone.CloneAndUpdate())
}

func TestCloneAndUpdate_HeapMinClonesMove_Distances(t *testing.T) {
	assert := assert.New(t)

	circuit := &HeapableCircuit2DMinClonesMove{
		Vertices: []model.CircuitVertex{
			// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
			NewVertex2D(0, 0),   // Index 0 after sorting
			NewVertex2D(0, 3),   // Index 1 after sorting
			NewVertex2D(3, 3),   // Index 5 after sorting
			NewVertex2D(3, 0),   // Index 4 after sorting
			NewVertex2D(1, 1),   // Index 2 after sorting
			NewVertex2D(1, 2.1), // Index 3 after sorting
		},
	}
	circuit.Prepare()
	circuit.BuildPerimiter()

	assert.Len(circuit.unattachedVertices, 2)
	assert.True(circuit.unattachedVertices[circuit.Vertices[2]])
	assert.True(circuit.unattachedVertices[circuit.Vertices[3]])

	assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[2]])
	assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[3]])

	assert.Equal(12.0, circuit.length)
	assert.Equal(8, circuit.closestEdges.Len())

	// No clone on first attachment - vertex {1,2.1} to edge {0,3}->{3,3}
	assert.Nil(circuit.CloneAndUpdate())

	assert.Len(circuit.unattachedVertices, 1)
	assert.True(circuit.unattachedVertices[circuit.Vertices[2]])

	assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[2]])
	assert.InDelta(0.5385336246535019133711157158298, circuit.distanceIncreases[circuit.Vertices[3]], model.Threshold)

	assert.InDelta(12.5385336246535019133711157158298, circuit.length, model.Threshold)
	assert.Equal(8, circuit.closestEdges.Len())

	assert.Equal(&heapDistanceToEdge{
		distance: 0.1327694499764709,
		vertex:   circuit.Vertices[3].(*Vertex2D),
		edge:     circuit.circuitEdges[3],
	}, circuit.closestEdges.Peek())

	// Clone on second attachment - vertex {1,2.1} to edge {0,0}->{0,3}
	clone := circuit.CloneAndUpdate().(*HeapableCircuit2DMinClonesMove)

	// Validate that the first circuit is unchanged.
	assert.Len(circuit.unattachedVertices, 1)
	assert.True(circuit.unattachedVertices[circuit.Vertices[2]])

	assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[2]])
	assert.InDelta(0.5385336246535019133711157158298, circuit.distanceIncreases[circuit.Vertices[3]], model.Threshold)

	assert.InDelta(12.5385336246535019133711157158298, circuit.length, model.Threshold)
	assert.Equal(7, circuit.closestEdges.Len())

	// Validate that the clone is updated correctly.
	assert.Len(clone.unattachedVertices, 1)
	assert.True(clone.unattachedVertices[clone.Vertices[2]])

	assert.Equal(0.0, clone.distanceIncreases[clone.Vertices[2]])
	assert.InDelta(0.6713030746299724753709331208575, clone.distanceIncreases[clone.Vertices[3]], model.Threshold)

	assert.InDelta(12.6713030746299724753709331208575, clone.length, model.Threshold)
	assert.Equal(7, clone.closestEdges.Len())

	// No clone on third update of circuit - vertex {1,1} to edge {0,0}->{0,3} or to edge {3,0}->{0,0}
	assert.Nil(circuit.CloneAndUpdate())

	assert.Len(circuit.unattachedVertices, 0)

	// 1.4142135623730950488016887242097 + 2.2360679774997896964091736687313 = 3.650281539872884745210862392941
	assert.InDelta(0.650281539872884745210862392941, circuit.distanceIncreases[circuit.Vertices[2]], model.Threshold)
	assert.InDelta(0.5385336246535019133711157158298, circuit.distanceIncreases[circuit.Vertices[3]], model.Threshold)

	assert.InDelta(13.1888151645263866585819781087708, circuit.length, model.Threshold)
	assert.Equal(7, circuit.closestEdges.Len())

	// No clone on first update of clone - vertex {1,1} to edge {0,0}->{1,2.1}
	assert.Nil(clone.CloneAndUpdate())

	assert.Len(clone.unattachedVertices, 0)

	assert.InDelta(0.18827289245049360514706414956918, clone.distanceIncreases[circuit.Vertices[2]], model.Threshold)
	assert.InDelta(0.20929442720758118, clone.distanceIncreases[circuit.Vertices[3]], model.Threshold)

	assert.InDelta(12.85957596708046608051799727042668, clone.length, model.Threshold)
	assert.Equal(6, clone.closestEdges.Len())
}

func TestPrepare_HeapMinClonesMove(t *testing.T) {
	assert := assert.New(t)
	circuit := &HeapableCircuit2DMinClonesMove{
		Vertices: []model.CircuitVertex{
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
	assert.ElementsMatch(circuit.Vertices, []model.CircuitVertex{
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

	assert.Equal(0.0, circuit.length)
	assert.Equal(0.0, circuit.GetLength())
	assert.Equal(0.0, circuit.GetLengthWithNext())

	assert.NotNil(circuit.closestEdges)
	assert.Equal(0, circuit.closestEdges.Len())

	assert.NotNil(circuit.circuit)
	assert.Len(circuit.circuit, 0)

	assert.NotNil(circuit.circuitEdges)
	assert.Len(circuit.circuitEdges, 0)

	assert.NotNil(circuit.midpoint)
	assert.InDelta((6.0+model.Threshold/3.0)/7.0, circuit.midpoint.X, model.Threshold)
	assert.InDelta(-5.0/7.0, circuit.midpoint.Y, model.Threshold)
}

func TestInsertVertex_HeapMinClonesMove(t *testing.T) {
	assert := assert.New(t)
	c := &HeapableCircuit2DMinClonesMove{
		Vertices: []model.CircuitVertex{
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
