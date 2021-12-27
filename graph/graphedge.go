package graph

import (
	"fmt"

	"github.com/fealos/lee-tsp-go/model"
)

type GraphEdge struct {
	path     []*GraphVertex
	distance float64
}

// NewGraphEdge creates a new GraphEdge from the start vertex to the end vertex.
// Its complexity is O(n*log(n)), due to needing to find the optimal path, which potentially involves
// checking each vertex in the graph, which are sorted by distance from the start vertex, for early escape.
// If a path cannot be created from the start vertex to the end vertex nil will be returned (the graph is asymmetric, so it is possible to connect only one way).
func NewGraphEdge(start *GraphVertex, end *GraphVertex) *GraphEdge {
	visited := make(map[*GraphVertex]bool)

	toVisit := model.NewHeap(func(a interface{}) float64 {
		return a.(*GraphEdge).distance
	})

	var toReturn *GraphEdge = nil
	startEdge := &GraphEdge{
		path:     []*GraphVertex{start},
		distance: 0.0,
	}
	for current, okay := startEdge, true; okay; current, okay = toVisit.PopHeap().(*GraphEdge) {
		currentVertex := current.GetEnd().(*GraphVertex)
		if currentVertex == end {
			toReturn = current
			break
		} else if !visited[currentVertex] {
			// Only visit vertices once, base on the shortest path to the vertex.
			visited[currentVertex] = true
			for v, dist := range currentVertex.adjacentVertices {
				// Don't push vertices with a shorter path, also prevent looping between adjacent vertices.
				// Note: we don't update visited yet, since there may be a vertex that is farther than the current vertex, but with a shorter connection to 'v'.
				if !visited[v] {
					next := &GraphEdge{
						path:     make([]*GraphVertex, len(current.path), len(current.path)+1),
						distance: current.distance + dist,
					}
					copy(next.path, current.path)
					next.path = append(next.path, v)
					toVisit.PushHeap(next)
				}
			}
		}
		current.Delete()
	}

	toVisit.Delete()
	return toReturn
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
	for _, v := range e.path {
		for _, o := range otherGraphEdge.path {
			if v == o {
				return true
			}
		}
	}
	return false
}

func (e *GraphEdge) Merge(other model.CircuitEdge) model.CircuitEdge {
	return NewGraphEdge(e.GetStart().(*GraphVertex), other.GetEnd().(*GraphVertex))
}

func (e *GraphEdge) Split(vertex model.CircuitVertex) (model.CircuitEdge, model.CircuitEdge) {
	return e.GetStart().EdgeTo(vertex), vertex.EdgeTo(e.GetEnd())
}

func (e *GraphEdge) ToString() string {
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
