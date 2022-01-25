package graph_test

import (
	"strings"
	"testing"

	"github.com/fealos/lee-tsp-go/graph"
	"github.com/stretchr/testify/assert"
)

func TestPathToAllFromAll_ShouldCreateAnEdgeFromEveryVertexToEveryOtherVertex(t *testing.T) {
	assert := assert.New(t)

	gen := &graph.GraphGenerator{
		MaxEdges:    6,
		MinEdges:    2,
		NumVertices: uint32(20),
	}

	g := gen.Create()
	assert.NotNil(g)
	defer g.Delete()

	graphMap := g.PathToAllFromAll()
	assert.Len(graphMap, 20)

	for _, start := range g.Vertices {
		edgeMap, okay := graphMap[start]
		assert.True(okay)
		assert.Len(edgeMap, 20)

		for _, destination := range g.Vertices {
			edge, okay2 := edgeMap[destination]
			assert.True(okay2)
			assert.Equal(start, edge.GetStart())
			assert.Equal(destination, edge.GetEnd())
		}
	}

	for k, v := range graphMap {
		delete(graphMap, k)
		for k2 := range v {
			delete(v, k2)
		}
	}
}

func TestToApi_ShouldPreserveGraphData(t *testing.T) {
	assert := assert.New(t)

	gen := &graph.GraphGenerator{
		MaxEdges:    5,
		MinEdges:    3,
		NumVertices: uint32(15),
	}

	g := gen.Create()
	assert.NotNil(g)
	defer g.Delete()

	api := g.ToApi()
	assert.NotNil(api)

	assert.Len(api.Vertices, len(g.Vertices))

	for _, v := range g.Vertices {
		var match *graph.GraphVertexApi
		for _, other := range api.Vertices {
			if strings.Compare(v.GetId(), other.Id) == 0 {
				match = other
				break
			}
		}
		assert.NotNil(match)
		assert.Len(match.AdjacentVertices, len(v.GetAdjacentVertices()))

		for dest, dist := range v.GetAdjacentVertices() {
			matchDist, okay := match.AdjacentVertices[dest.GetId()]
			assert.True(okay)
			assert.Equal(dist, matchDist)
		}
	}
}
