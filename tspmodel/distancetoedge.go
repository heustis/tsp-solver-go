package tspmodel

import (
	"fmt"
	"strconv"
)

// DistanceToEdge is used to cache the distance (typically the distance increase) from a vertex to an edge.
type DistanceToEdge struct {
	Vertex   CircuitVertex
	Edge     CircuitEdge
	Distance float64
}

// GetDistanceToEdgeForHeap is used to retrieve the Distance field from a DistanceToEdge in a heap,
// so that DistanceToEdges can be sorted from closest to farthest.
func GetDistanceToEdgeForHeap(a interface{}) float64 {
	if dist, okay := a.(*DistanceToEdge); okay {
		return dist.Distance
	}
	panic(fmt.Sprintf("Received non-DistanceToEdge object=%v", a))
}

// HasVertex returns true if the supplied interface is another DistanceToEdge, and it has the same Vertex as this DistanceToEdge.
// If is used to remove other DistanceToEdges from a heap after the vertex in this DistanceToEdge was attached to its edge.
func (h *DistanceToEdge) HasVertex(x interface{}) bool {
	dist, okay := x.(*DistanceToEdge)
	return okay && dist.Vertex == h.Vertex
}

// String converts the DistanceToEdge to a json string.
func (h *DistanceToEdge) String() string {
	distString := strconv.FormatFloat(h.Distance, 'f', -1, 64)
	return fmt.Sprintf(`{"vertex":%s,"edge":%s,"distance":%s}`, h.Vertex, h.Edge, distString)
}

var _ fmt.Stringer = (*DistanceToEdge)(nil)
