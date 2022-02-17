package graph

import (
	"fmt"
	"strings"

	"github.com/heustis/lee-tsp-go/model"
)

type GraphVertex struct {
	id               string
	adjacentVertices map[*GraphVertex]float64
	paths            map[model.CircuitVertex]*GraphEdge
}

func NewGraphVertex(id string) *GraphVertex {
	g := &GraphVertex{
		id:               id,
		adjacentVertices: make(map[*GraphVertex]float64),
		paths:            make(map[model.CircuitVertex]*GraphEdge),
	}
	return g
}

func (v *GraphVertex) AddAdjacentVertex(other *GraphVertex, distance float64) {
	v.adjacentVertices[other] = distance
}

func (v *GraphVertex) Delete() {
	for k := range v.adjacentVertices {
		delete(v.adjacentVertices, k)
	}
	v.adjacentVertices = nil
	v.DeletePaths()
}

func (v *GraphVertex) DeletePaths() {
	for end, edge := range v.paths {
		edge.Delete()
		delete(v.paths, end)
	}
	v.paths = nil
}

func (v *GraphVertex) DistanceTo(other model.CircuitVertex) float64 {
	return v.EdgeTo(other).GetLength()
}

// EdgeTo creates a new GraphEdge from the start vertex to the end vertex.
// Its complexity is O(n*log(n)), due to needing to find the optimal path, which potentially involves
// checking each vertex in the graph, which are sorted by distance from the start vertex, for early escape.
// If a path cannot be created from the start vertex to the end vertex nil will be returned (the graph is asymmetric, so it is possible to connect only one way).
func (start *GraphVertex) EdgeTo(end model.CircuitVertex) model.CircuitEdge {
	if path := start.paths[end]; path != nil && len(path.path) > 0 {
		return path
	} else {
		visited := make(map[*GraphVertex]bool)

		toVisit := model.NewHeap(func(a interface{}) float64 {
			return a.(*GraphEdge).distance
		})

		current := &GraphEdge{
			path:     []*GraphVertex{start},
			distance: 0.0,
		}
		for okay := true; okay; current, okay = toVisit.PopHeap().(*GraphEdge) {
			currentVertex := current.GetEnd().(*GraphVertex)
			if current.GetEnd() == end {
				start.paths[end] = current
				break
			}

			if !visited[currentVertex] {
				// Only visit vertices once, base on the shortest path to the vertex.
				// Update visited map only with the current vertex, to ensure that it is shortest path to that vertex.
				visited[currentVertex] = true
				for v, dist := range currentVertex.adjacentVertices {
					// Avoid creating unnecessary objects if we've already visted the vertex.
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

		return current
	}
}

func (v *GraphVertex) Equals(other interface{}) bool {
	if otherVertex, okay := other.(*GraphVertex); okay {
		return strings.Compare(v.id, otherVertex.id) == 0
	}
	return false
}

func (v *GraphVertex) GetAdjacentVertices() map[*GraphVertex]float64 {
	return v.adjacentVertices
}

func (v *GraphVertex) GetId() string {
	return v.id
}

func (v *GraphVertex) GetPaths() map[model.CircuitVertex]*GraphEdge {
	return v.paths
}

// InitializePaths sets up the map of the most efficient edges from this vertex to all other vertices in the graph.
// Its complexity is O(n*e*log(n*e)), where n is the number of nodes and e is the average number of edges off of each node.
//   We only visit each node once, however for each node we add each of its connected nodes to a heap to find the shortest path to the next unvisited node.
func (start *GraphVertex) InitializePaths() {
	toVisit := model.NewHeap(func(a interface{}) float64 {
		return a.(*GraphEdge).distance
	})

	startEdge := &GraphEdge{
		path:     []*GraphVertex{start},
		distance: 0.0,
	}
	for current, okay := startEdge, true; okay; current, okay = toVisit.PopHeap().(*GraphEdge) {
		currentVertex := current.GetEnd().(*GraphVertex)
		// Only visit vertices once, base on the shortest path to the vertex.
		if _, visited := start.paths[currentVertex]; !visited {
			start.paths[currentVertex] = current
			for v, dist := range currentVertex.adjacentVertices {
				// Don't push vertices with a shorter path, also prevent looping between adjacent vertices.
				// Note: we don't update visited yet, since there may be a vertex that is farther than the current vertex, but with a shorter connection to 'v'.
				if _, nextVisited := start.paths[v]; !nextVisited {
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
}

func (v *GraphVertex) String() string {
	s := fmt.Sprintf(`{"id":"%s","adjacentVertices":{`, v.id)
	// Sort the adjacent vertices by distance to ensure consistent string representation of the vertex.
	h := model.NewHeap(func(a interface{}) float64 {
		return a.(*StringVertex).distance
	})
	defer h.Delete()
	for adj, dist := range v.adjacentVertices {
		h.PushHeap(&StringVertex{
			id:       adj.id,
			distance: dist,
		})
	}
	isFirst := true
	for v, okay := h.PopHeap().(*StringVertex); okay; v, okay = h.PopHeap().(*StringVertex) {
		if !isFirst {
			s += ","
		}
		isFirst = false
		s += fmt.Sprintf(`"%s":%g`, v.id, v.distance)
	}
	s += "}}"
	return s
}

type StringVertex struct {
	id       string
	distance float64
}

var _ model.CircuitVertex = (*GraphVertex)(nil)
