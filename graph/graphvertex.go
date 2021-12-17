package graph

import (
	"fmt"
	"strings"

	"github.com/fealos/lee-tsp-go/model"
)

type GraphVertex struct {
	id               string
	adjacentVertices map[*GraphVertex]float64
}

func (g *GraphVertex) Delete() {
	for k := range g.adjacentVertices {
		delete(g.adjacentVertices, k)
	}
	g.adjacentVertices = nil
}

func (g *GraphVertex) DistanceTo(other model.CircuitVertex) float64 {
	return g.EdgeTo(other).GetLength()
}

func (g *GraphVertex) EdgeTo(other model.CircuitVertex) model.CircuitEdge {
	return NewGraphEdge(g, other.(*GraphVertex))
}

func (g *GraphVertex) Equals(other interface{}) bool {
	if otherVertex, okay := other.(*GraphVertex); okay {
		return strings.Compare(g.id, otherVertex.id) == 0
	}
	return false
}

func (g *GraphVertex) GetAdjacentVertices() map[*GraphVertex]float64 {
	return g.adjacentVertices
}

func (g *GraphVertex) GetId() string {
	return g.id
}

// PathToAll creates a map of the most efficient edges from this vertex to all other vertices in the graph.
// Its complexity is O(n*e*log(n*e)), where n is the number of nodes and e is the average number of edges off of each node.
//   We only visit each node once, however for each node we add each of its connected nodes to a heap to find the shortest path to the next unvisited node.
func (g *GraphVertex) PathToAll() map[model.CircuitVertex]model.CircuitEdge {
	edges := make(map[model.CircuitVertex]model.CircuitEdge)

	toVisit := model.NewHeap(func(a interface{}) float64 {
		return a.(*GraphEdge).distance
	})

	startEdge := &GraphEdge{
		path:     []*GraphVertex{g},
		distance: 0.0,
	}
	for current, okay := startEdge, true; okay; current, okay = toVisit.PopHeap().(*GraphEdge) {
		currentVertex := current.GetEnd().(*GraphVertex)
		// Only visit vertices once, base on the shortest path to the vertex.
		if _, visited := edges[currentVertex]; !visited {
			edges[currentVertex] = current
			for v, dist := range currentVertex.adjacentVertices {
				// Don't push vertices with a shorter path, also prevent looping between adjacent vertices.
				// Note: we don't update visited yet, since there may be a vertex that is farther than the current vertex, but with a shorter connection to 'v'.
				if _, nextVisited := edges[v]; !nextVisited {
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
	}

	return edges
}

func (g *GraphVertex) ToString() string {
	s := fmt.Sprintf(`{"id":"%s","adjacentVertices":{`, g.id)
	// Sort the adjacent vertices by distance to ensure consistent string representation of the vertex.
	h := model.NewHeap(func(a interface{}) float64 {
		return a.(*toStringVertex).distance
	})
	defer h.Delete()
	for adj, dist := range g.adjacentVertices {
		h.PushHeap(&toStringVertex{
			id:       adj.id,
			distance: dist,
		})
	}
	isFirst := true
	for v, okay := h.PopHeap().(*toStringVertex); okay; v, okay = h.PopHeap().(*toStringVertex) {
		if !isFirst {
			s += ","
		}
		isFirst = false
		s += fmt.Sprintf(`"%s":%g`, v.id, v.distance)
	}
	s += "}}"
	return s
}

type toStringVertex struct {
	id       string
	distance float64
}

var _ model.CircuitVertex = (*GraphVertex)(nil)
