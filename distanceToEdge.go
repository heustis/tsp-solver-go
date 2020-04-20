package tsp

// DistanceToEdge represents the distance from a vertex to an edge
type DistanceToEdge struct {
	edge     *Edge
	distance float64
}

// GetEdge returns the edge
func (e DistanceToEdge) GetEdge() *Edge {
	return e.edge
}

// GetDistance returns the distance from the vertex to the edge
func (e DistanceToEdge) GetDistance() float64 {
	return e.distance
}

// NewDistanceToEdge creates a distance-to-edge object
func NewDistanceToEdge(edge *Edge, distance float64) *DistanceToEdge {
	return &DistanceToEdge{edge: edge, distance: distance}
}

// IsCloser checks if the supplied edge is closer to the vertex than this DistanceToEdge.
// Returns true if the supplied edge is closer, along with the closest edge.
func (e *DistanceToEdge) IsCloser(edge *Edge, vertex *Vertex) (bool, *DistanceToEdge) {
	if edge.Contains(vertex) {
		return false, e
	}
	distanceToEdge := edge.DistanceIncrease(vertex)
	if distanceToEdge < e.GetDistance() {
		return true, NewDistanceToEdge(edge, distanceToEdge)
	}
	return false, e
}
