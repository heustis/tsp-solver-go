package model

import "fmt"

// Threshold defines how close two values can be and still be considered identical (e.g. for de-duplicating points).
const Threshold = 0.0000001

// Circuit provides an abstract representation of a set of points (locations, vertices) for the TSP solver to interact with.
// This allows it to ignore whether the implementation is a set of N-dimentional points, a graph, or any other representation of points.
type Circuit interface {
	// BuildPerimeter creates an initial circuit, using the minimum vertices required to fully enclose the other (interior) vertices.
	// For example, when using 2-D points, this constructs a convex polygon such that all points are either vertices or inside the polygon.
	BuildPerimiter()

	// FindNextVertexAndEdge determines the next vertex to add to the circuit, along with which edge it should be added to.
	// For example, when using 2-D points, this finds the point with the minimum distance to its nearest edge (returning both that point and edge).
	FindNextVertexAndEdge() (CircuitVertex, CircuitEdge)

	// GetAttachedVertices returns all vertices that have been added to the circuit (either as part of BuildPerimeter or Update).
	// This returns them in the order they should be traversed as part of the circuit (ignoring any unattached vertices).
	GetAttachedVertices() []CircuitVertex

	// GetLength returns the length of the circuit (at the current stage of processing).
	GetLength() float64

	// GetUnattachedVertices returns the set of vertices that have not been added to the circuit yet. (all of these points are internal to the perimeter)
	GetUnattachedVertices() map[CircuitVertex]bool

	// Prepare may be used by a circuit to pre-compute values that will save time while processing the circuit.
	// Prepare should be called prior to performing any other operations on a circuit.
	Prepare()

	// Update adds the supplied vertex to circuit by splitting the supplied edge and creating two edges with the supplied point as the common vertex of the edges.
	Update(vertexToAdd CircuitVertex, edgeToSplit CircuitEdge)
}

// CircuitVertex provides an abstract representation of a single point (location, vertex) for the TSP solver to interact with.
type CircuitVertex interface {
	Equal
	fmt.Stringer
	// DistanceTo returns the distance between the two vertices; this should always be a positive number.
	DistanceTo(other CircuitVertex) float64
	// EdgeTo creates a new CircuitEdge from this point (start) to the supplied point (end).
	EdgeTo(end CircuitVertex) CircuitEdge
}

// CircuitVertex provides an abstract representation of an edge for the TSP solver to interact with.
type CircuitEdge interface {
	Equal
	fmt.Stringer
	// DistanceIncrease returns the difference in length between the edge
	// and the two edges formed by inserting the vertex between the edge's start and end.
	// For example, if start->end has a length of 5, start->vertex has a length of 3,
	//  and vertex->end has a length of 6, this will return 4 (i.e. 6 + 3 - 5)
	DistanceIncrease(vertex CircuitVertex) float64
	// GetEnd returns the ending point of this edge.
	GetEnd() CircuitVertex
	// GetLength returns the distance between the start and end vertices.
	GetLength() float64
	// GetStart returns the starting point of this edge.
	GetStart() CircuitVertex
	// Intersects checks if the two edges go through at least one identical point.
	// Note: Edges may share multiple points if they are co-linear, or in the use-case of graphs.
	Intersects(other CircuitEdge) bool
	// Merge creates a new edge starting from this edge's start vertex and ending at the supplied edge's end vertex.
	Merge(CircuitEdge) CircuitEdge
	// Split creates two new edges "start-to-vertex" and "vertex-to-end" based on this edge and the supplied vertex.
	Split(vertex CircuitVertex) (CircuitEdge, CircuitEdge)
}

// Deduplicator is a function that takes in an array of vertices, and returns a copy of the array without duplicate points.
type Deduplicator func([]CircuitVertex) []CircuitVertex

// PerimeterBuilder creates an initial circuit, using the minimum vertices required to fully enclose the other (interior) vertices.
// For example, when using 2-D points, this constructs a convex polygon such that all points are either vertices or inside the polygon.
// This returns the perimeter as an array of edges, and a map of unattached vertices.
type PerimeterBuilder func(verticesArg []CircuitVertex) (edges []CircuitEdge, unattached map[CircuitVertex]bool)
