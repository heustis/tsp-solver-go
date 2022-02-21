package modelapi

import "github.com/heustis/tsp-solver-go/graph"

// PointGraph is the API representation of a single point in a graph.
// It references its neighbors by name, in an array, to avoid circular references and have consistent field names in its JSON representation.
type PointGraph struct {
	Id string `json:"id" validate:"required,min=1"`
	// Validator/v10 does not support `unique` with nil values in the array, see validate_test.go, so the array does not use pointers.
	// Once that is supported Neighbors can be converted to []*PointGraphNeighbor.
	Neighbors []PointGraphNeighbor `json:"neighbors" validate:"required,min=1,unique=Id,dive,required"`
}

// PointGraphNeighbor is a neighboring point to a PointGraph point.
// Its id must correspond to the id of a point in the request's array of PointGraphs.
// The distance is the distance from the PointGraph point to the point with the id, this may be asymmetrical.
type PointGraphNeighbor struct {
	Id       string  `json:"id" validate:"required,min=1"`
	Distance float64 `json:"distance" validate:"required,min=0"`
}

// ToGraph converts an API request into a graph.
func (api *TspRequest) ToGraph() *graph.Graph {
	vertices := []*graph.GraphVertex{}

	// This map deduplicates vertices (by ID), and prevents repeat processing of vertices.
	vertexMap := make(map[string]*graph.GraphVertex)

	for _, vApi := range api.PointsGraph {
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
		for _, n := range vApi.Neighbors {
			if adj, okay := vertexMap[n.Id]; okay {
				v.AddAdjacentVertex(adj, n.Distance)
			} else {
				adj = graph.NewGraphVertex(n.Id)
				vertexMap[n.Id] = adj
				v.AddAdjacentVertex(adj, n.Distance)
			}
		}
	}

	for k := range vertexMap {
		delete(vertexMap, k)
	}

	return graph.NewGraph(vertices)
}

// ToApiFromGraph converts a graph into an API response.
func ToApiFromGraph(g *graph.Graph) *TspRequest {
	api := &TspRequest{
		PointsGraph: []*PointGraph{},
	}

	for _, v := range g.GetVertices() {
		vApi := &PointGraph{
			Id:        v.GetId(),
			Neighbors: make([]PointGraphNeighbor, 0, len(v.GetAdjacentVertices())),
		}

		for adj, distance := range v.GetAdjacentVertices() {
			vApi.Neighbors = append(vApi.Neighbors, PointGraphNeighbor{
				Id:       adj.GetId(),
				Distance: distance,
			})
		}

		api.PointsGraph = append(api.PointsGraph, vApi)
	}
	return api
}
