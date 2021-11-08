package model2d

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerimeter_HeapMinClones(t *testing.T) {
	assert := assert.New(t)
	circuit := &HeapableCircuit2DMinClones{
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

	assert.NotNil(circuit.convexVertices)
	assert.Len(circuit.convexVertices, 5)
	assert.True(circuit.convexVertices[circuit.Vertices[0]])
	assert.True(circuit.convexVertices[circuit.Vertices[1]])
	assert.True(circuit.convexVertices[circuit.Vertices[4]])
	assert.True(circuit.convexVertices[circuit.Vertices[6]])
	assert.True(circuit.convexVertices[circuit.Vertices[7]])

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

func TestAttachAndDetach(t *testing.T) {
	assert := assert.New(t)

	circuit := &HeapableCircuit2DMinClones{
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

	testVert := circuit.closestEdges.PopHeap().(*heapDistanceToEdge)
	circuit.attachVertex(testVert)
	circuit.detachVertex(testVert.vertex)

	assert.Len(circuit.unattachedVertices, 2)
	assert.True(circuit.unattachedVertices[circuit.Vertices[2]])
	assert.True(circuit.unattachedVertices[circuit.Vertices[3]])

	// assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[2]])
	// assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[3]])

	assert.Equal(12.0, circuit.length)
	assert.Equal(7, circuit.closestEdges.Len())

	circuit.attachVertex(testVert)
	assert.Len(circuit.unattachedVertices, 1)
	assert.True(circuit.unattachedVertices[circuit.Vertices[2]])

	// assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[2]])
	// assert.InDelta(0.5385336246535019133711157158298, circuit.distanceIncreases[circuit.Vertices[3]], model.Threshold)

	assert.InDelta(12.5385336246535019133711157158298, circuit.length, model.Threshold)
	assert.Equal(8, circuit.closestEdges.Len())
}

func TestAttachAndDetachIndexZero(t *testing.T) {
	assert := assert.New(t)

	circuit := &HeapableCircuit2DMinClones{
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

	testVert := &heapDistanceToEdge{
		vertex:   circuit.Vertices[2].(*Vertex2D),
		edge:     circuit.circuitEdges[0],
		distance: circuit.circuitEdges[0].DistanceIncrease(circuit.Vertices[2]),
	}
	circuit.attachVertex(testVert)
	assert.Equal([]model.CircuitEdge{
		NewEdge2D(circuit.Vertices[4].(*Vertex2D), circuit.Vertices[2].(*Vertex2D)),
		NewEdge2D(circuit.Vertices[2].(*Vertex2D), circuit.Vertices[5].(*Vertex2D)),
		NewEdge2D(circuit.Vertices[5].(*Vertex2D), circuit.Vertices[1].(*Vertex2D)),
		NewEdge2D(circuit.Vertices[1].(*Vertex2D), circuit.Vertices[0].(*Vertex2D)),
		NewEdge2D(circuit.Vertices[0].(*Vertex2D), circuit.Vertices[4].(*Vertex2D)),
	}, circuit.circuitEdges)
	assert.Equal([]model.CircuitVertex{
		circuit.Vertices[4],
		circuit.Vertices[2],
		circuit.Vertices[5],
		circuit.Vertices[1],
		circuit.Vertices[0],
	}, circuit.circuit)

	circuit.detachVertex(testVert.vertex)
	assert.Equal([]model.CircuitEdge{
		NewEdge2D(circuit.Vertices[4].(*Vertex2D), circuit.Vertices[5].(*Vertex2D)),
		NewEdge2D(circuit.Vertices[5].(*Vertex2D), circuit.Vertices[1].(*Vertex2D)),
		NewEdge2D(circuit.Vertices[1].(*Vertex2D), circuit.Vertices[0].(*Vertex2D)),
		NewEdge2D(circuit.Vertices[0].(*Vertex2D), circuit.Vertices[4].(*Vertex2D)),
	}, circuit.circuitEdges)
	assert.Equal([]model.CircuitVertex{
		circuit.Vertices[4],
		circuit.Vertices[5],
		circuit.Vertices[1],
		circuit.Vertices[0],
	}, circuit.circuit)

	assert.Len(circuit.unattachedVertices, 2)
	assert.True(circuit.unattachedVertices[circuit.Vertices[2]])
	assert.True(circuit.unattachedVertices[circuit.Vertices[3]])

	// assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[2]])
	// assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[3]])

	assert.Equal(12.0, circuit.length)
	assert.Equal(7, circuit.closestEdges.Len())

	circuit.attachVertex(testVert)
	assert.Len(circuit.unattachedVertices, 1)
	assert.True(circuit.unattachedVertices[circuit.Vertices[3]])

	// assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[3]])
	// assert.InDelta(testVert.distance, circuit.distanceIncreases[circuit.Vertices[2]], model.Threshold)

	assert.InDelta(12+testVert.distance, circuit.length, model.Threshold)
	assert.Equal(8, circuit.closestEdges.Len())
}

func TestAttachAndDetachLastIndex(t *testing.T) {
	assert := assert.New(t)

	circuit := &HeapableCircuit2DMinClones{
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

	testVert := &heapDistanceToEdge{
		vertex:   circuit.Vertices[2].(*Vertex2D),
		edge:     circuit.circuitEdges[3],
		distance: circuit.circuitEdges[3].DistanceIncrease(circuit.Vertices[2]),
	}
	circuit.attachVertex(testVert)
	assert.Equal([]model.CircuitEdge{
		NewEdge2D(circuit.Vertices[4].(*Vertex2D), circuit.Vertices[5].(*Vertex2D)),
		NewEdge2D(circuit.Vertices[5].(*Vertex2D), circuit.Vertices[1].(*Vertex2D)),
		NewEdge2D(circuit.Vertices[1].(*Vertex2D), circuit.Vertices[0].(*Vertex2D)),
		NewEdge2D(circuit.Vertices[0].(*Vertex2D), circuit.Vertices[2].(*Vertex2D)),
		NewEdge2D(circuit.Vertices[2].(*Vertex2D), circuit.Vertices[4].(*Vertex2D)),
	}, circuit.circuitEdges)
	assert.Equal([]model.CircuitVertex{
		circuit.Vertices[4],
		circuit.Vertices[5],
		circuit.Vertices[1],
		circuit.Vertices[0],
		circuit.Vertices[2],
	}, circuit.circuit)

	circuit.detachVertex(testVert.vertex)
	assert.Equal([]model.CircuitEdge{
		NewEdge2D(circuit.Vertices[4].(*Vertex2D), circuit.Vertices[5].(*Vertex2D)),
		NewEdge2D(circuit.Vertices[5].(*Vertex2D), circuit.Vertices[1].(*Vertex2D)),
		NewEdge2D(circuit.Vertices[1].(*Vertex2D), circuit.Vertices[0].(*Vertex2D)),
		NewEdge2D(circuit.Vertices[0].(*Vertex2D), circuit.Vertices[4].(*Vertex2D)),
	}, circuit.circuitEdges)
	assert.Equal([]model.CircuitVertex{
		circuit.Vertices[4],
		circuit.Vertices[5],
		circuit.Vertices[1],
		circuit.Vertices[0],
	}, circuit.circuit)

	assert.Len(circuit.unattachedVertices, 2)
	assert.True(circuit.unattachedVertices[circuit.Vertices[2]])
	assert.True(circuit.unattachedVertices[circuit.Vertices[3]])

	// assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[2]])
	// assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[3]])

	assert.Equal(12.0, circuit.length)
	assert.Equal(7, circuit.closestEdges.Len())

	circuit.attachVertex(testVert)
	assert.Len(circuit.unattachedVertices, 1)
	assert.True(circuit.unattachedVertices[circuit.Vertices[3]])

	// assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[3]])
	// assert.InDelta(testVert.distance, circuit.distanceIncreases[circuit.Vertices[2]], model.Threshold)

	assert.InDelta(12+testVert.distance, circuit.length, model.Threshold)
	assert.Equal(8, circuit.closestEdges.Len())
}

func TestCloneAndUpdate_HeapMinClones(t *testing.T) {
	assert := assert.New(t)

	circuit := &HeapableCircuit2DMinClones{
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
	clone := circuit.CloneAndUpdate().(*HeapableCircuit2DMinClones)
	assert.Len(circuit.unattachedVertices, 2)
	assert.Len(clone.unattachedVertices, 2)
	assert.False(circuit.unattachedVertices[circuit.Vertices[5]])
	assert.False(clone.unattachedVertices[circuit.Vertices[5]])
	assert.Len(circuit.circuit, 6)
	assert.Len(clone.circuit, 6)
	assert.Len(circuit.circuitEdges, 6)
	assert.Len(clone.circuitEdges, 6)
	assert.Equal(12, clone.closestEdges.Len())

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
	assert.Len(clone.unattachedVertices, 1)
	assert.Equal(12, clone.closestEdges.Len())

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
	cloneOfClone, okay := clone.CloneAndUpdate().(*HeapableCircuit2DMinClones)
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

	// Index 2 should move to edge 3, cloning required
	assert.Equal(&heapDistanceToEdge{
		edge:     clone.circuitEdges[3],
		vertex:   clone.Vertices[2].(*Vertex2D),
		distance: 0.32754172885551824,
	}, clone.closestEdges.Peek())
	assert.NotNil(clone.CloneAndUpdate())

	// Index 2 should move to edge 4, cloning required
	assert.Equal(&heapDistanceToEdge{
		edge:     cloneOfClone.circuitEdges[4],
		vertex:   cloneOfClone.Vertices[2].(*Vertex2D),
		distance: 3.341664064126334,
	}, cloneOfClone.closestEdges.Peek())
	assert.Nil(cloneOfClone.CloneAndUpdate())
}

func TestCloneAndUpdate_HeapMinClones_Distances(t *testing.T) {
	assert := assert.New(t)

	circuit := &HeapableCircuit2DMinClones{
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

	// assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[2]])
	// assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[3]])

	assert.Equal(12.0, circuit.length)
	assert.Equal(8, circuit.closestEdges.Len())

	// No clone on first attachment - vertex {1,2.1} to edge {0,3}->{3,3}
	assert.Nil(circuit.CloneAndUpdate())

	assert.Len(circuit.unattachedVertices, 1)
	assert.True(circuit.unattachedVertices[circuit.Vertices[2]])

	// assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[2]])
	// assert.InDelta(0.5385336246535019133711157158298, circuit.distanceIncreases[circuit.Vertices[3]], model.Threshold)

	assert.InDelta(12.5385336246535019133711157158298, circuit.length, model.Threshold)
	assert.Equal(8, circuit.closestEdges.Len())

	assert.Equal(&heapDistanceToEdge{
		distance: 0.1327694499764709,
		vertex:   circuit.Vertices[3].(*Vertex2D),
		edge:     circuit.circuitEdges[3],
	}, circuit.closestEdges.Peek())

	// Clone on second attachment - vertex {1,2.1} to edge {0,0}->{0,3}
	clone := circuit.CloneAndUpdate().(*HeapableCircuit2DMinClones)

	// Validate that the first circuit is unchanged.
	assert.Len(circuit.unattachedVertices, 1)
	assert.True(circuit.unattachedVertices[circuit.Vertices[2]])

	// assert.Equal(0.0, circuit.distanceIncreases[circuit.Vertices[2]])
	// assert.InDelta(0.5385336246535019133711157158298, circuit.distanceIncreases[circuit.Vertices[3]], model.Threshold)

	assert.InDelta(12.5385336246535019133711157158298, circuit.length, model.Threshold)
	// assert.Equal(5, circuit.closestEdges.Len())

	// Validate that the clone is updated correctly.
	assert.Len(clone.unattachedVertices, 1)
	assert.True(clone.unattachedVertices[clone.Vertices[2]])

	// assert.Equal(0.0, clone.distanceIncreases[clone.Vertices[2]])
	// assert.InDelta(0.6713030746299724753709331208575, clone.distanceIncreases[clone.Vertices[3]], model.Threshold)

	assert.InDelta(12.6713030746299724753709331208575, clone.length, model.Threshold)
	assert.Equal(5, clone.closestEdges.Len())

	// No clone on third update of circuit - vertex {1,1} to edge {0,0}->{0,3} or to edge {3,0}->{0,0}
	assert.Nil(circuit.CloneAndUpdate())

	assert.Len(circuit.unattachedVertices, 0)

	// 1.4142135623730950488016887242097 + 2.2360679774997896964091736687313 = 3.650281539872884745210862392941
	// assert.InDelta(0.650281539872884745210862392941, circuit.distanceIncreases[circuit.Vertices[2]], model.Threshold)
	// assert.InDelta(0.5385336246535019133711157158298, circuit.distanceIncreases[circuit.Vertices[3]], model.Threshold)

	assert.InDelta(13.1888151645263866585819781087708, circuit.length, model.Threshold)
	// assert.Equal(4, circuit.closestEdges.Len())

	// No clone on first update of clone - vertex {1,1} to edge {0,0}->{1,2.1}
	assert.Nil(clone.CloneAndUpdate())

	assert.Len(clone.unattachedVertices, 0)

	// assert.InDelta(0.18827289245049360514706414956918, clone.distanceIncreases[circuit.Vertices[2]], model.Threshold)
	// assert.InDelta(0.20929442720758118, clone.distanceIncreases[circuit.Vertices[3]], model.Threshold)

	assert.InDelta(12.85957596708046608051799727042668, clone.length, model.Threshold)
	assert.Equal(4, clone.closestEdges.Len())
}

func TestPrepare_HeapMinClones(t *testing.T) {
	assert := assert.New(t)
	circuit := &HeapableCircuit2DMinClones{
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
	assert.ElementsMatch(circuit.Vertices, []*Vertex2D{
		NewVertex2D(-15+model.Threshold/3.0, -15-model.Threshold/3.0),
		NewVertex2D(-7, 6),
		NewVertex2D(0, 0),
		NewVertex2D(3, 0),
		NewVertex2D(3, 13),
		NewVertex2D(7, 6),
		NewVertex2D(15, -15),
	})

	assert.NotNil(circuit.convexVertices)
	assert.Len(circuit.convexVertices, 0)

	assert.NotNil(circuit.unattachedVertices)
	assert.Len(circuit.unattachedVertices, 0)

	assert.Equal(0.0, circuit.length)
	assert.Equal(0.0, circuit.GetLength())
	assert.Equal(0.0, circuit.GetLengthWithNext())

	assert.NotNil(circuit.closestEdges)
	assert.Equal(0, circuit.closestEdges.Len())

	assert.NotNil(circuit.circuit)
	assert.Len(circuit.circuit, 0)

	assert.NotNil(circuit.circuitEdges)
	assert.Len(circuit.circuitEdges, 0)
}

func TestInsertVertex_HeapMinClones(t *testing.T) {
	assert := assert.New(t)
	c := &HeapableCircuit2DMinClones{
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
		vertices2d := make([]*Vertex2D, t.len)
		err := json.Unmarshal([]byte(t.vertices), &vertices2d)

		vertices := make([]model.CircuitVertex, t.len)
		for i, v2d := range vertices2d {
			vertices[i] = v2d
		}

		assert.Nil(err, "Failed to unmarshal vertices for test=", i)
		actual := solveWithLogging_HeapMinClones(&HeapableCircuit2DMinClones{
			Vertices: vertices,
		})
		assert.InDelta(t.expectedLength, actual.length, model.Threshold)
	}
}

func solveWithLogging_HeapMinClones(circuit *HeapableCircuit2DMinClones) *HeapableCircuit2DMinClones {
	circuit.Prepare()
	circuit.BuildPerimiter()

	circuitHeap := model.NewHeap(func(a interface{}) float64 {
		return a.(*HeapableCircuit2DMinClones).GetLengthWithNext()
	})
	circuitHeap.PushHeap(circuit)

	next := circuitHeap.PopHeap().(*HeapableCircuit2DMinClones)
	for i := 0; len(next.GetUnattachedVertices()) > 0 || next.GetLengthWithNext() < next.GetLength(); next = circuitHeap.PopHeap().(*HeapableCircuit2DMinClones) {
		toAttach := next.closestEdges.Peek()
		clone := next.CloneAndUpdate()
		circuitHeap.PushHeap(next)
		if clone != nil {
			circuitBytes, _ := json.Marshal(clone.(*HeapableCircuit2DMinClones).circuit)
			fmt.Printf("Step %d: Created clone=%p from existing=%p with \n\ttoAttach=%s\n\tcircuit=%s\n\theap=%s\n", i, clone, next, toAttach.(*heapDistanceToEdge).ToString(), string(circuitBytes), clone.(*HeapableCircuit2DMinClones).closestEdges.ToString())
			circuitHeap.PushHeap(clone)
		} else {
			circuitBytes, _ := json.Marshal(next.circuit)
			fmt.Printf("Step %d: Updated existing=%p with \n\ttoAttach=%s\n\tcircuit=%s\n\theap=%s\n", i, next, toAttach.(*heapDistanceToEdge).ToString(), string(circuitBytes), next.closestEdges.ToString())
		}
		i++
	}

	// clean up the heap and each circuitHeap
	circuitHeap.Delete()

	return next
}
