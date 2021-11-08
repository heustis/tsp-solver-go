package model2d

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerimeter_Heap(t *testing.T) {
	assert := assert.New(t)
	circuit := &HeapableCircuit2D{
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

	l := 0.0
	for _, edge := range circuit.circuitEdges {
		l += edge.GetLength()
	}
	assert.InDelta(l, circuit.GetLength(), model.Threshold)

	assert.Len(circuit.unattachedVertices, 3)
	assert.True(circuit.unattachedVertices[circuit.Vertices[2]])
	assert.True(circuit.unattachedVertices[circuit.Vertices[3]])
	assert.True(circuit.unattachedVertices[circuit.Vertices[5]])
	assert.Equal(circuit.unattachedVertices, circuit.GetUnattachedVertices())

	assert.InDelta(circuit.length+0.763503994948632, circuit.GetLengthWithNext(), model.Threshold)

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

func TestCloneAndUpdate(t *testing.T) {
	assert := assert.New(t)

	circuit := &HeapableCircuit2D{
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

	clone := circuit.CloneAndUpdate().(*HeapableCircuit2D)

	assert.Len(circuit.unattachedVertices, 3)
	assert.Len(circuit.circuit, 5)
	assert.Len(circuit.circuitEdges, 5)
	assert.Equal(14, circuit.closestEdges.Len())

	assert.Equal(&heapDistanceToEdge{
		edge:     circuit.circuitEdges[2],
		vertex:   circuit.Vertices[5].(*Vertex2D),
		distance: 1.628650237136812,
	}, circuit.closestEdges.Peek())

	assert.Len(clone.unattachedVertices, 2)
	assert.Len(clone.circuit, 6)
	assert.Len(clone.circuitEdges, 6)
	assert.Equal(12, clone.closestEdges.Len())

	assert.Equal(&heapDistanceToEdge{
		edge:     clone.circuitEdges[1],
		vertex:   clone.Vertices[3].(*Vertex2D),
		distance: 5.09082042374693,
	}, clone.closestEdges.Peek())

	assert.InDelta(circuit.GetLength()+0.763503994948632, clone.GetLength(), model.Threshold)

	// Validate that cloning a clone does not affect the original circuit
	cloneOfClone := clone.CloneAndUpdate().(*HeapableCircuit2D)

	assert.Len(circuit.unattachedVertices, 3)
	assert.Len(circuit.circuit, 5)
	assert.Len(circuit.circuitEdges, 5)
	assert.Equal(14, circuit.closestEdges.Len())

	assert.Equal(&heapDistanceToEdge{
		edge:     circuit.circuitEdges[2],
		vertex:   circuit.Vertices[5].(*Vertex2D),
		distance: 1.628650237136812,
	}, circuit.closestEdges.Peek())

	assert.Len(clone.unattachedVertices, 2)
	assert.Len(clone.circuit, 6)
	assert.Len(clone.circuitEdges, 6)
	assert.Equal(11, clone.closestEdges.Len())

	assert.Equal(&heapDistanceToEdge{
		edge:     clone.circuitEdges[5],
		vertex:   clone.Vertices[2].(*Vertex2D),
		distance: 7.9605428386450825,
	}, clone.closestEdges.Peek())

	assert.Len(cloneOfClone.unattachedVertices, 1)
	assert.Len(cloneOfClone.circuit, 7)
	assert.Len(cloneOfClone.circuitEdges, 7)
	assert.Equal(7, cloneOfClone.closestEdges.Len())

	assert.Equal(&heapDistanceToEdge{
		edge:     cloneOfClone.circuitEdges[1],
		vertex:   cloneOfClone.Vertices[2].(*Vertex2D),
		distance: 5.003830723297881,
	}, cloneOfClone.closestEdges.Peek())

	// Validate that cloning a circuit with only one vertex left to attach just updates that object and doesn't create a clone
	assert.Nil(cloneOfClone.CloneAndUpdate())

	assert.Len(cloneOfClone.unattachedVertices, 0)
	assert.Len(cloneOfClone.circuit, 8)
	assert.Len(cloneOfClone.circuitEdges, 8)
	assert.Equal(0, cloneOfClone.closestEdges.Len())

	// Validate that cloning a completed circuit makes no changes
	assert.Nil(cloneOfClone.CloneAndUpdate())

	assert.Len(cloneOfClone.unattachedVertices, 0)
	assert.Len(cloneOfClone.circuit, 8)
	assert.Len(cloneOfClone.circuitEdges, 8)
	assert.Equal(0, cloneOfClone.closestEdges.Len())

}

func TestPrepare_Heap(t *testing.T) {
	assert := assert.New(t)
	circuit := &HeapableCircuit2D{
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

func TestInsertVertex_Heap(t *testing.T) {
	assert := assert.New(t)
	c := &HeapableCircuit2D{
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
