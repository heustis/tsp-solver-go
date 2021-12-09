package graph

import "github.com/fealos/lee-tsp-go/model"

type Graph struct {
	Vertices []*GraphVertex
}

func (g *Graph) Delete() {
	for _, v := range g.Vertices {
		v.Delete()
	}
	g.Vertices = nil
}

// PathToAllFromAll creates a map of the most efficient edges from each vertex to all other vertices in the graph.
// Its complexity is O(n*n*e*log(n*e)), where n is the number of nodes and e is the average number of edges off of each node.
func (g *Graph) PathToAllFromAll() map[model.CircuitVertex]map[model.CircuitVertex]model.CircuitEdge {
	edges := make(map[model.CircuitVertex]map[model.CircuitVertex]model.CircuitEdge)
	for _, v := range g.Vertices {
		edges[v] = v.PathToAll()
	}
	return edges
}

func (g *Graph) ToApi() *GraphApi {
	api := &GraphApi{
		Vertices: []*GraphVertexApi{},
	}

	for _, v := range g.Vertices {
		vApi := &GraphVertexApi{
			Id:               v.id,
			AdjacentVertices: make(map[string]float64),
		}

		for adj, distance := range v.adjacentVertices {
			vApi.AdjacentVertices[adj.id] = distance
		}

		api.Vertices = append(api.Vertices, vApi)
	}
	return api
}

func (g *Graph) ToString() string {
	s := `{"vertices":[`

	isFirst := true
	for _, v := range g.Vertices {
		if !isFirst {
			s += ","
		}
		isFirst = false
		s += v.ToString()
	}

	s += `]}`
	return s
}
