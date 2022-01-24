package circuit

import "github.com/fealos/lee-tsp-go/tspmodel"

// ClonableCircuit is a Circuit variant where the circuit may be cloned with each update, depending on the implementation,
// so that each clone represents a different branch of solving the circuit.
type ClonableCircuit interface {
	tspmodel.Deletable

	// BuildPerimeter creates an initial circuit, using the minimum vertices required to fully enclose the other (interior) vertices.
	// For example, when using 2-D points, this constructs a convex polygon such that all points are either vertices or inside the polygon.
	BuildPerimiter()

	// CloneAndUpdate combines updating the shortest clone via FindNextVertexAndEdge and Update so that either it is updated in place, or a clone is created.
	// For example, one approach (O(n!)) is to:
	//  - create a copy (B) of the ClonableCircuit (A),
	//  - update one of the versions so that the vertex closest to its nearest edge is attached to that edge,
	//  - update the other circuit so that it cannot have that vertex attached to that edge.
	CloneAndUpdate() ClonableCircuit

	// FindNextVertexAndEdge determines the next vertex to add to the circuit, along with which edge it should be added to.
	// For example, when using 2-D points, this finds the point with the minimum distance to its nearest edge (returning both that point and edge).
	FindNextVertexAndEdge() (tspmodel.CircuitVertex, tspmodel.CircuitEdge)

	// GetAttachedVertices returns all vertices that have been added to the circuit (either as part of BuildPerimeter or Update).
	// This returns them in the order they should be traversed as part of the circuit (ignoring any unattached vertices).
	GetAttachedVertices() []tspmodel.CircuitVertex

	// GetLength returns the length of the circuit (at the current stage of processing).
	GetLength() float64

	// GetLengthWithNext returns the length of the circuit, if the next cloesest vertex were attached.
	// This allows the tspsolver to prioritize combinations that may reduce the length increase of detached vertices (due to new edges being closer to those vertices).
	GetLengthWithNext() float64

	// GetUnattachedVertices returns the set of vertices that have not been added to the circuit yet. (all of these points are internal to the perimeter)
	GetUnattachedVertices() map[tspmodel.CircuitVertex]bool

	// Prepare may be used by a circuit to pre-compute values that will save time while processing the circuit.
	// Prepare should be called prior to performing any other operations on a circuit.
	Prepare()
}

// ClonableCircuitSolver is a wrapper for a ClonableCircuit and its clones that allows them to match the Circuit interface.
type ClonableCircuitSolver struct {
	circuits      *tspmodel.Heap
	numClones     int
	numIterations int
}

func NewClonableCircuitSolver(initialCircuit ClonableCircuit) tspmodel.Circuit {
	tspsolver := &ClonableCircuitSolver{
		circuits: tspmodel.NewHeap(getClonableLength),
	}
	tspsolver.circuits.PushHeap(initialCircuit)
	return tspsolver
}

func getClonableLength(a interface{}) float64 {
	return a.(ClonableCircuit).GetLengthWithNext()
}

func (c *ClonableCircuitSolver) BuildPerimiter() {
	c.circuits.Peek().(ClonableCircuit).BuildPerimiter()
}

func (c *ClonableCircuitSolver) FindNextVertexAndEdge() (tspmodel.CircuitVertex, tspmodel.CircuitEdge) {
	return c.circuits.Peek().(ClonableCircuit).FindNextVertexAndEdge()
}

func (c *ClonableCircuitSolver) GetAttachedVertices() []tspmodel.CircuitVertex {
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

func (c *ClonableCircuitSolver) GetUnattachedVertices() map[tspmodel.CircuitVertex]bool {
	return c.circuits.Peek().(ClonableCircuit).GetUnattachedVertices()
}

func (c *ClonableCircuitSolver) Prepare() {
	c.circuits.Peek().(ClonableCircuit).Prepare()
	c.numIterations = 0
	c.numClones = 0
}

func (c *ClonableCircuitSolver) Update(vertexToAdd tspmodel.CircuitVertex, edgeToSplit tspmodel.CircuitEdge) {
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
		c.circuits = tspmodel.NewHeap(getClonableLength)
		c.circuits.PushHeap(result)
	}
}

var _ tspmodel.Circuit = (*ClonableCircuitSolver)(nil)
