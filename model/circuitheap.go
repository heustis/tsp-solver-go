package model

// HeapableCircuit extends Circuit to facilitate its use in a min heap,
// where each entry is a copy of the circuit at a different stage of completion
type HeapableCircuit interface {
	Deletable

	// BuildPerimeter creates an initial circuit, using the minimum vertices required to fully enclose the other (interior) vertices.
	// For example, when using 2-D points, this constructs a convex polygon such that all points are either vertices or inside the polygon.
	// This is the second step in the TSP solver.
	BuildPerimiter()

	// CloneAndUpdate creates a copy of this circuit updated so that the vertex closest to its nearest edge is attached to that edge.
	// This circuit then is prevented from having that vertex being added to that edge.
	//
	// Note: If this circuit has only two edges that the vertex can connect to, the clone will attach it to the cloeset edge and this circuit will attach it to the remaining edge.
	CloneAndUpdate() HeapableCircuit

	// Delete cleans up this heap to prevent memory leaks.
	Delete()

	// GetAttachedVertices returns all vertices that have been added to the circuit (either as part of BuildPerimeter or Update).
	// This returns them in the order they should be traversed as part of the circuit (ignoring any unattached vertices).
	GetAttachedVertices() []CircuitVertex

	// GetLength returns the length of the circuit (at the current stage of processing).
	GetLength() float64

	// GetLengthWithNext returns the length of the circuit, if the next cloesest vertex were attached.
	// This allows the heap to prioritize combinations that may reduce the length increase of detached vertices (due to new edges being closer to those vertices).
	GetLengthWithNext() float64

	// GetUnattachedVertices returns the set of vertices that have not been added to the circuit yet. (all of these points are internal to the perimeter)
	GetUnattachedVertices() map[CircuitVertex]bool

	// Prepare is a method that some implementation may use to pre-compute some values to save processing time when computing the optimal circuit.
	// This is the first step in the TSP solver.
	Prepare()
}
