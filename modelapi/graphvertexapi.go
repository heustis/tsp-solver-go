package modelapi

import "github.com/fealos/lee-tsp-go/graph"

type GraphVertexApi struct {
	Id               string             `json:"id"`
	AdjacentVertices map[string]float64 `json:"adjacentVertices"`
}

type GraphApi struct {
	Vertices []*GraphVertexApi `json:"vertices"`
}

func (api *GraphApi) ToGraph() *graph.Graph {
	vertices := []*graph.GraphVertex{}

	// This map deduplicates vertices (by ID), and prevents repeat processing of vertices.
	vertexMap := make(map[string]*graph.GraphVertex)

	for _, vApi := range api.Vertices {
		var v *graph.GraphVertex

		// Ensure each vertex is created only once; re-use the vertex if it was created while processing adjacent vertices of an earlier vertex.
		if existing, okay := vertexMap[vApi.Id]; okay {
			v = existing
		} else {
			v = graph.NewGraphVertex(vApi.Id)
			vertexMap[vApi.Id] = v
		}
		vertices = append(vertices, v)

		// Create one vertex for each adjacent vertex, unless that vertex already exists, in which case re-use it.
		for adjId, dist := range vApi.AdjacentVertices {
			if adj, okay := vertexMap[adjId]; okay {
				v.AddAdjacentVertex(adj, dist)
			} else {
				adj = graph.NewGraphVertex(adjId)
				vertexMap[adjId] = adj
				v.AddAdjacentVertex(adj, dist)
			}
		}
	}

	for k := range vertexMap {
		delete(vertexMap, k)
	}

	return graph.NewGraph(vertices)
}

func ToApiGraph(g *graph.Graph) *GraphApi {
	api := &GraphApi{
		Vertices: []*GraphVertexApi{},
	}

	for _, v := range g.GetVertices() {
		vApi := &GraphVertexApi{
			Id:               v.GetId(),
			AdjacentVertices: make(map[string]float64),
		}

		for adj, distance := range v.GetAdjacentVertices() {
			vApi.AdjacentVertices[adj.GetId()] = distance
		}

		api.Vertices = append(api.Vertices, vApi)
	}
	return api
}
