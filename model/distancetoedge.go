package model

import (
	"fmt"
)

type DistanceToEdge struct {
	Vertex   CircuitVertex
	Edge     CircuitEdge
	Distance float64
}

func GetDistanceToEdgeForHeap(a interface{}) float64 {
	if dist, okay := a.(*DistanceToEdge); okay {
		return dist.Distance
	}
	panic(fmt.Sprintf("Received non-DistanceToEdge object=%v", a))
}

func (h *DistanceToEdge) HasVertex(x interface{}) bool {
	dist, okay := x.(*DistanceToEdge)
	return okay && dist.Vertex == h.Vertex
}

func (h *DistanceToEdge) ToString() string {
	return fmt.Sprintf(`{"vertex":%s,"edge":%s,"distance":%v}`, h.Vertex.ToString(), h.Edge.ToString(), h.Distance)
}

var _ Printable = (*DistanceToEdge)(nil)
