package circuit

import "github.com/heustis/tsp-solver-go/model"

// ClonableCircuit is a Circuit variant where the circuit may be cloned with each update, depending on the implementation,
// so that each clone represents a different branch of solving the circuit.
type ClonableCircuit interface {
	model.Deletable

	// CloneAndUpdate combines updating the shortest clone via FindNextVertexAndEdge and Update so that either it is updated in place, or a clone is created.
	// For example, one approach (O(n!)) is to:
	//  - create a copy (B) of the ClonableCircuit (A),
	//  - update one of the versions so that the vertex closest to its nearest edge is attached to that edge,
	//  - update the other circuit so that it cannot have that vertex attached to that edge.
	// If the circuit is updated in place, `nil` should be returned.
	CloneAndUpdate() ClonableCircuit

	// FindNextVertexAndEdge determines the next vertex to add to the circuit, along with which edge it should be added to.
	// For example, when using 2-D points, this finds the point with the minimum distance to its nearest edge (returning both that point and edge).
	FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge)

	// GetAttachedVertices returns all vertices that have been added to the circuit (either as part of BuildPerimeter or Update).
	// This returns them in the order they should be traversed as part of the circuit (ignoring any unattached vertices).
	GetAttachedVertices() []model.CircuitVertex

	// GetLength returns the length of the circuit (at the current stage of processing).
	GetLength() float64

	// GetLengthWithNext returns the length of the circuit, if the next cloesest vertex were attached.
	// This allows the solver to prioritize combinations that may reduce the length increase of detached vertices (due to new edges being closer to those vertices).
	GetLengthWithNext() float64

	// GetUnattachedVertices returns the set of vertices that have not been added to the circuit yet. (all of these points are internal to the perimeter)
	GetUnattachedVertices() map[model.CircuitVertex]bool
}

// ClonableCircuitSolver is a wrapper for a ClonableCircuit and its clones that allows them to match the Circuit interface.
type ClonableCircuitSolver struct {
	maxClones     int
	circuits      *model.Heap
	numClones     int
	numIterations int
}

// NewClonableCircuitSolver creates a ClonableCircuitSolver with the supplied ClonableCircuit as its initial circuit.
func NewClonableCircuitSolver(initialCircuit ClonableCircuit) *ClonableCircuitSolver {
	solver := &ClonableCircuitSolver{
		circuits:      model.NewHeap(getClonableLength),
		numClones:     0,
		numIterations: 0,
		maxClones:     -1,
	}
	solver.circuits.PushHeap(initialCircuit)
	return solver
}

func getClonableLength(a interface{}) float64 {
	return a.(ClonableCircuit).GetLengthWithNext()
}

func (c *ClonableCircuitSolver) FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge) {
	return c.circuits.Peek().(ClonableCircuit).FindNextVertexAndEdge()
}

func (c *ClonableCircuitSolver) GetAttachedVertices() []model.CircuitVertex {
	return c.circuits.Peek().(ClonableCircuit).GetAttachedVertices()
}

func (c *ClonableCircuitSolver) GetLength() float64 {
	return c.circuits.Peek().(ClonableCircuit).GetLength()
}

func (c *ClonableCircuitSolver) GetNumClones() int {
	return c.numClones
}

func (c *ClonableCircuitSolver) GetNumIterations() int {
	return c.numIterations
}

func (c *ClonableCircuitSolver) GetUnattachedVertices() map[model.CircuitVertex]bool {
	return c.circuits.Peek().(ClonableCircuit).GetUnattachedVertices()
}

func (c *ClonableCircuitSolver) SetMaxClones(max int) {
	c.maxClones = max
}

func (c *ClonableCircuitSolver) Update(vertexToAdd model.CircuitVertex, edgeToSplit model.CircuitEdge) {
	if _, isCompleted := c.circuits.Peek().(*CompletedCircuit); isCompleted {
		return
	}

	next := c.circuits.PopHeap().(ClonableCircuit)
	clone := next.CloneAndUpdate()
	c.circuits.PushHeap(next)
	if clone != nil {
		c.circuits.PushHeap(clone)
		c.numClones++
	}
	c.numIterations++

	// Check if the circuit is completed. If so, clean up the heap and clones, so that only the completed circuit remains.
	next = c.circuits.Peek().(ClonableCircuit)
	if len(next.GetUnattachedVertices()) == 0 && next.GetLengthWithNext() >= next.GetLength() {
		result := &CompletedCircuit{
			Circuit: next.GetAttachedVertices(),
			Length:  next.GetLength(),
		}

		// Clean up the heap and each circuitHeap.
		c.circuits.Delete()
		next.Delete()

		// Create a new heap with only the completed circuit in it.
		c.circuits = model.NewHeap(getClonableLength)
		c.circuits.PushHeap(result)
	} else if c.maxClones > 0 && c.numClones > c.maxClones {
		c.trimClones()
	}
}

func (c *ClonableCircuitSolver) trimClones() {
	// trimClones() is only called when there is MaxClones+1 circuits, so we only need to discard the worst circuit.
	worstCircuit := c.circuits.PopHeap().(ClonableCircuit)

	// Prioritize preserving clones that are the closest to completion, so use the length per attached vertex rather than raw length of the circuit.
	worstLength := worstCircuit.GetLengthWithNext() / float64(len(worstCircuit.GetAttachedVertices()))

	retainedCircuits := make([]interface{}, 0, c.maxClones)

	for current := c.circuits.PopHeap(); current != nil; current = c.circuits.PopHeap() {
		currentCircuit := current.(ClonableCircuit)
		// If the current circuit is worse than the previous worst, add the previous worst to the retained circuits, and track the current circuit as the new worst.
		if currentLength := currentCircuit.GetLengthWithNext() / float64(len(currentCircuit.GetAttachedVertices())); currentLength > worstLength {
			retainedCircuits = append(retainedCircuits, worstCircuit)
			worstCircuit = currentCircuit
			worstLength = currentLength
		} else {
			// If the current circuit is better than the previous worst, retain the current circuit.
			retainedCircuits = append(retainedCircuits, currentCircuit)
		}
	}

	c.circuits.PushAll(retainedCircuits...)
}

var _ model.Circuit = (*ClonableCircuitSolver)(nil)
