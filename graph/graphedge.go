package graph

import (
	"fmt"

	"github.com/heustis/tsp-solver-go/model"
)

type GraphEdge struct {
	path     []*GraphVertex
	distance float64
}

func (e *GraphEdge) Delete() {
	// Note: Deleting a graph edge should not delete the graph itself, use Graph.Delete() to delete all GraphVertexes.
	e.path = nil
}

func (e *GraphEdge) DistanceIncrease(vertex model.CircuitVertex) float64 {
	a, b := e.Split(vertex)
	return a.GetLength() + b.GetLength() - e.GetLength()
}

func (e *GraphEdge) Equals(other interface{}) bool {
	if otherEdge, okay := other.(*GraphEdge); okay {
		if len(e.path) != len(otherEdge.path) {
			return false
		}
		for i := 0; i < len(e.path); i++ {
			if e.path[i] != otherEdge.path[i] {
				return false
			}
		}
	}
	return true
}

func (e *GraphEdge) GetEnd() model.CircuitVertex {
	if lastIndex := len(e.path) - 1; lastIndex < 0 {
		return nil
	} else {
		return e.path[lastIndex]
	}
}

func (e *GraphEdge) GetLength() float64 {
	return e.distance
}

func (e *GraphEdge) GetPath() []*GraphVertex {
	return e.path
}

func (e *GraphEdge) GetStart() model.CircuitVertex {
	return e.path[0]
}

func (e *GraphEdge) Intersects(other model.CircuitEdge) bool {
	otherGraphEdge := other.(*GraphEdge)
	// Using a map enables the complexity of Intersects() to be O(n)
	graphIds := make(map[string]bool)
	for _, v := range e.path {
		graphIds[v.id] = true
	}
	for _, o := range otherGraphEdge.path {
		if graphIds[o.id] {
			return true
		}
	}
	return false
}

func (e *GraphEdge) Merge(other model.CircuitEdge) model.CircuitEdge {
	return e.GetStart().(*GraphVertex).EdgeTo(other.GetEnd())
}

func (e *GraphEdge) Split(vertex model.CircuitVertex) (model.CircuitEdge, model.CircuitEdge) {
	return e.GetStart().EdgeTo(vertex), vertex.EdgeTo(e.GetEnd())
}

func (e *GraphEdge) String() string {
	s := "["
	isFirst := true
	for _, v := range e.path {
		if !isFirst {
			s += ","
		}
		isFirst = false
		s += fmt.Sprintf(`"%s"`, v.id)
	}
	s += "]"
	return s
}

var _ model.CircuitEdge = (*GraphEdge)(nil)
