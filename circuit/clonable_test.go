package circuit_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestClonableCircuitSolver(t *testing.T) {
	assert := assert.New(t)

	initVertices := []model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15),                    // Index 0 after sorting
		model2d.NewVertex2D(0, 0),                        // Index 2 after sorting
		model2d.NewVertex2D(15, -15),                     // Index 7 after sorting
		model2d.NewVertex2D(15, -15+model.Threshold/2.0), // Removed by deduplication
		model2d.NewVertex2D(3, 0),                        // Index 3 after sorting
		model2d.NewVertex2D(3, 13),                       // Index 4 after sorting
		model2d.NewVertex2D(8, 5),                        // Index 5 after sorting
		model2d.NewVertex2D(9, 6),                        // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),                       // Index 1 after sorting
	}

	c := circuit.NewClonableCircuitSolver(
		circuit.NewClonableCircuitImpl(model2d.DeduplicateVertices(initVertices), model2d.BuildPerimiter)).(*circuit.ClonableCircuitSolver)

	assert.Len(c.GetAttachedVertices(), 0)
	assert.Len(c.GetUnattachedVertices(), 0)
	assert.Equal(0.0, c.GetLength())
	assert.Equal(0, c.GetNumClones())
	assert.Equal(0, c.GetNumIterations())

	assert.Len(c.GetAttachedVertices(), 5)

	unattached := c.GetUnattachedVertices()
	assert.Len(unattached, 3)

	assert.True(unattached[initVertices[2]])
	assert.True(unattached[initVertices[3]])
	assert.True(unattached[initVertices[5]])

	assert.InDelta(95.73863479511238, c.GetLength(), model.Threshold)
	assert.Equal(0, c.GetNumClones())
	assert.Equal(0, c.GetNumIterations())

	nextVertex, nextEdge := c.FindNextVertexAndEdge()
	assert.Equal(initVertices[5], nextVertex)
	assert.True(initVertices[7].EdgeTo(initVertices[6]).Equals(nextEdge))

	c.Update(nextVertex, nextEdge)
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Len(c.GetUnattachedVertices(), 2)
	assert.InDelta(96.50213879006101, c.GetLength(), model.Threshold)
	assert.Equal(0, c.GetNumClones())
	assert.Equal(1, c.GetNumIterations())

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Len(c.GetUnattachedVertices(), 2)
	assert.InDelta(96.50213879006101, c.GetLength(), model.Threshold)
	assert.Equal(1, c.GetNumClones())
	assert.Equal(2, c.GetNumIterations())

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Len(c.GetUnattachedVertices(), 2)
	assert.InDelta(97.36728503224919, c.GetLength(), model.Threshold)
	assert.Equal(1, c.GetNumClones())
	assert.Equal(3, c.GetNumIterations())

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 7)
	assert.Len(c.GetUnattachedVertices(), 1)
	assert.InDelta(101.59295921380794, c.GetLength(), model.Threshold)
	assert.Equal(1, c.GetNumClones())
	assert.Equal(4, c.GetNumIterations())

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 8)
	assert.Len(c.GetUnattachedVertices(), 0)
	assert.InDelta(106.59678993710583, c.GetLength(), model.Threshold)
	assert.Equal(1, c.GetNumClones())
	assert.Equal(5, c.GetNumIterations())

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 8)
	assert.Len(c.GetUnattachedVertices(), 0)
	assert.InDelta(106.59678993710583, c.GetLength(), model.Threshold)
	assert.Equal(1, c.GetNumClones())
	assert.Equal(5, c.GetNumIterations())
}

func TestClonableCircuitSolver_CloneOnAttach(t *testing.T) {
	//TODO
	assert := assert.New(t)

	initVertices := []model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15),                    // Index 0 after sorting
		model2d.NewVertex2D(0, 0),                        // Index 2 after sorting
		model2d.NewVertex2D(15, -15),                     // Index 7 after sorting
		model2d.NewVertex2D(15, -15+model.Threshold/2.0), // Removed by deduplication
		model2d.NewVertex2D(3, 0),                        // Index 3 after sorting
		model2d.NewVertex2D(3, 13),                       // Index 4 after sorting
		model2d.NewVertex2D(8, 5),                        // Index 5 after sorting
		model2d.NewVertex2D(9, 6),                        // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),                       // Index 1 after sorting
	}

	cImpl := circuit.NewClonableCircuitImpl(model2d.DeduplicateVertices(initVertices), model2d.BuildPerimiter)
	cImpl.CloneOnFirstAttach = true
	c := circuit.NewClonableCircuitSolver(cImpl).(*circuit.ClonableCircuitSolver)

	assert.Len(c.GetAttachedVertices(), 0)
	assert.Len(c.GetUnattachedVertices(), 0)
	assert.Equal(0.0, c.GetLength())
	assert.Equal(0, c.GetNumClones())
	assert.Equal(0, c.GetNumIterations())

	assert.Len(c.GetAttachedVertices(), 5)

	unattached := c.GetUnattachedVertices()
	assert.Len(unattached, 3)

	assert.True(unattached[initVertices[2]])
	assert.True(unattached[initVertices[3]])
	assert.True(unattached[initVertices[5]])

	assert.InDelta(95.73863479511238, c.GetLength(), model.Threshold)
	assert.Equal(0, c.GetNumClones())
	assert.Equal(0, c.GetNumIterations())

	nextVertex, nextEdge := c.FindNextVertexAndEdge()
	assert.Equal(initVertices[5], nextVertex)
	assert.True(initVertices[7].EdgeTo(initVertices[6]).Equals(nextEdge))

	c.Update(nextVertex, nextEdge)
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Len(c.GetUnattachedVertices(), 2)
	assert.InDelta(96.50213879006101, c.GetLength(), model.Threshold)
	assert.Equal(0, c.GetNumClones())
	assert.Equal(1, c.GetNumIterations())

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Len(c.GetUnattachedVertices(), 2)
	assert.InDelta(96.50213879006101, c.GetLength(), model.Threshold)
	assert.Equal(1, c.GetNumClones())
	assert.Equal(2, c.GetNumIterations())

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Len(c.GetUnattachedVertices(), 2)
	assert.InDelta(97.36728503224919, c.GetLength(), model.Threshold)
	assert.Equal(1, c.GetNumClones())
	assert.Equal(3, c.GetNumIterations())

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 7)
	assert.Len(c.GetUnattachedVertices(), 1)
	assert.InDelta(101.59295921380794, c.GetLength(), model.Threshold)
	assert.Equal(1, c.GetNumClones())
	assert.Equal(4, c.GetNumIterations())

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 8)
	assert.Len(c.GetUnattachedVertices(), 0)
	assert.InDelta(106.59678993710583, c.GetLength(), model.Threshold)
	assert.Equal(1, c.GetNumClones())
	assert.Equal(5, c.GetNumIterations())

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 8)
	assert.Len(c.GetUnattachedVertices(), 0)
	assert.InDelta(106.59678993710583, c.GetLength(), model.Threshold)
	assert.Equal(1, c.GetNumClones())
	assert.Equal(5, c.GetNumIterations())
}
func TestClonableCircuitSolver_MaxClones(t *testing.T) {
	//TODO
	assert := assert.New(t)

	initVertices := []model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15),                    // Index 0 after sorting
		model2d.NewVertex2D(0, 0),                        // Index 2 after sorting
		model2d.NewVertex2D(15, -15),                     // Index 7 after sorting
		model2d.NewVertex2D(15, -15+model.Threshold/2.0), // Removed by deduplication
		model2d.NewVertex2D(3, 0),                        // Index 3 after sorting
		model2d.NewVertex2D(3, 13),                       // Index 4 after sorting
		model2d.NewVertex2D(8, 5),                        // Index 5 after sorting
		model2d.NewVertex2D(9, 6),                        // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),                       // Index 1 after sorting
	}

	cImpl := circuit.NewClonableCircuitImpl(model2d.DeduplicateVertices(initVertices), model2d.BuildPerimiter)
	cImpl.CloneOnFirstAttach = true
	c := circuit.NewClonableCircuitSolver(cImpl).(*circuit.ClonableCircuitSolver)
	c.MaxClones = 5

	assert.Len(c.GetAttachedVertices(), 0)
	assert.Len(c.GetUnattachedVertices(), 0)
	assert.Equal(0.0, c.GetLength())
	assert.Equal(0, c.GetNumClones())
	assert.Equal(0, c.GetNumIterations())

	assert.Len(c.GetAttachedVertices(), 5)

	unattached := c.GetUnattachedVertices()
	assert.Len(unattached, 3)

	assert.True(unattached[initVertices[2]])
	assert.True(unattached[initVertices[3]])
	assert.True(unattached[initVertices[5]])

	assert.InDelta(95.73863479511238, c.GetLength(), model.Threshold)
	assert.Equal(0, c.GetNumClones())
	assert.Equal(0, c.GetNumIterations())

	nextVertex, nextEdge := c.FindNextVertexAndEdge()
	assert.Equal(initVertices[5], nextVertex)
	assert.True(initVertices[7].EdgeTo(initVertices[6]).Equals(nextEdge))

	c.Update(nextVertex, nextEdge)
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Len(c.GetUnattachedVertices(), 2)
	assert.InDelta(96.50213879006101, c.GetLength(), model.Threshold)
	assert.Equal(0, c.GetNumClones())
	assert.Equal(1, c.GetNumIterations())

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Len(c.GetUnattachedVertices(), 2)
	assert.InDelta(96.50213879006101, c.GetLength(), model.Threshold)
	assert.Equal(1, c.GetNumClones())
	assert.Equal(2, c.GetNumIterations())

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Len(c.GetUnattachedVertices(), 2)
	assert.InDelta(97.36728503224919, c.GetLength(), model.Threshold)
	assert.Equal(1, c.GetNumClones())
	assert.Equal(3, c.GetNumIterations())

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 7)
	assert.Len(c.GetUnattachedVertices(), 1)
	assert.InDelta(101.59295921380794, c.GetLength(), model.Threshold)
	assert.Equal(1, c.GetNumClones())
	assert.Equal(4, c.GetNumIterations())

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 8)
	assert.Len(c.GetUnattachedVertices(), 0)
	assert.InDelta(106.59678993710583, c.GetLength(), model.Threshold)
	assert.Equal(1, c.GetNumClones())
	assert.Equal(5, c.GetNumIterations())

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 8)
	assert.Len(c.GetUnattachedVertices(), 0)
	assert.InDelta(106.59678993710583, c.GetLength(), model.Threshold)
	assert.Equal(1, c.GetNumClones())
	assert.Equal(5, c.GetNumIterations())
}
