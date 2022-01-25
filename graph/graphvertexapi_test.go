package graph_test

import (
	"strings"
	"testing"

	"github.com/fealos/lee-tsp-go/graph"
	"github.com/stretchr/testify/assert"
)

func TestToGraph_ShouldPreserveGraphData(t *testing.T) {
	assert := assert.New(t)

	gen := &graph.GraphGenerator{
		MaxEdges:    4,
		MinEdges:    2,
		NumVertices: uint32(10),
	}

	g := gen.Create()
	assert.NotNil(g)
	defer g.Delete()

	api := g.ToApi()
	assert.NotNil(api)

	g2 := api.ToGraph()
	assert.NotNil(g2)
	defer g2.Delete()

	assert.Len(api.Vertices, len(g.Vertices))
	assert.Len(g2.Vertices, len(g.Vertices))

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

		var g2Match *graph.GraphVertex
		for _, other := range g2.Vertices {
			if strings.Compare(v.GetId(), other.GetId()) == 0 {
				g2Match = other
				break
			}
		}
		assert.NotNil(g2Match)
		assert.Len(g2Match.GetAdjacentVertices(), len(v.GetAdjacentVertices()))

		for dest, dist := range v.GetAdjacentVertices() {
			matchDist, okay := match.AdjacentVertices[dest.GetId()]
			assert.True(okay)
			assert.Equal(dist, matchDist)

			_, okay = g2Match.GetAdjacentVertices()[dest]
			assert.False(okay)

			for g2Dest, g2Dist := range g2Match.GetAdjacentVertices() {
				if strings.Compare(dest.GetId(), g2Dest.GetId()) == 0 {
					okay = true
					assert.Equal(dist, g2Dist)
				}
			}
			assert.True(okay)
		}
	}
}
