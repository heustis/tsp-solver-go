package graph

type GraphVertexApi struct {
	Id               string             `json:"id"`
	AdjacentVertices map[string]float64 `json:"adjacentVertices"`
}

type GraphApi struct {
	Vertices []*GraphVertexApi `json:"vertices"`
}

func (api *GraphApi) ToGraph() *Graph {
	g := &Graph{
		Vertices: []*GraphVertex{},
	}

	vertexMap := make(map[string]*GraphVertex)

	for _, vApi := range api.Vertices {
		var v *GraphVertex

		// Ensure each vertex is created only once; re-use the vertex if it was created while processing adjacent vertices of an earlier vertex.
		if existing, okay := vertexMap[vApi.Id]; okay {
			v = existing
		} else {
			v = &GraphVertex{
				id:               vApi.Id,
				adjacentVertices: make(map[*GraphVertex]float64),
			}
			vertexMap[vApi.Id] = v
		}
		g.Vertices = append(g.Vertices, v)

		// Create one vertex for each adjacent vertex, unless that vertex already exists, in which case re-use it.
		for adjId, dist := range vApi.AdjacentVertices {
			if adj, okay := vertexMap[adjId]; okay {
				v.adjacentVertices[adj] = dist
			} else {
				adj = &GraphVertex{
					id:               adjId,
					adjacentVertices: make(map[*GraphVertex]float64),
				}
				vertexMap[adjId] = adj
				v.adjacentVertices[adj] = dist
			}
		}
	}

	for k := range vertexMap {
		delete(vertexMap, k)
	}

	return g
}
