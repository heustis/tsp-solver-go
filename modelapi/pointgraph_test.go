package modelapi_test

import (
	"strings"
	"testing"

	"github.com/heustis/tsp-solver-go/graph"
	"github.com/heustis/tsp-solver-go/modelapi"
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

	api := modelapi.ToApiFromGraph(g)
	assert.NotNil(api)

	g2 := api.ToGraph()
	assert.NotNil(g2)
	defer g2.Delete()

	assert.Len(api.PointsGraph, len(g.GetVertices()))
	assert.Len(g2.GetVertices(), len(g.GetVertices()))
	assert.Len(api.Points2D, 0)
	assert.Len(api.Points3D, 0)

	for _, v := range g.GetVertices() {
		var match *modelapi.PointGraph
		for _, other := range api.PointsGraph {
			if strings.Compare(v.GetId(), other.Id) == 0 {
				match = other
				break
			}
		}
		assert.NotNil(match)
		assert.Len(match.Neighbors, len(v.GetAdjacentVertices()))

		var g2Match *graph.GraphVertex
		for _, other := range g2.GetVertices() {
			if strings.Compare(v.GetId(), other.GetId()) == 0 {
				g2Match = other
				break
			}
		}
		assert.NotNil(g2Match)
		assert.Len(g2Match.GetAdjacentVertices(), len(v.GetAdjacentVertices()))

		for dest, dist := range v.GetAdjacentVertices() {
			numMatch := 0
			for _, n := range match.Neighbors {
				if strings.Compare(dest.GetId(), n.Id) == 0 {
					numMatch++
					assert.Equal(dist, n.Distance)
				}
			}
			assert.Equal(1, numMatch)

			_, okay := g2Match.GetAdjacentVertices()[dest]
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

func TestToApiGraph_ShouldPreserveGraphData(t *testing.T) {
	assert := assert.New(t)

	gen := &graph.GraphGenerator{
		MaxEdges:    5,
		MinEdges:    3,
		NumVertices: uint32(15),
	}

	g := gen.Create()
	assert.NotNil(g)
	defer g.Delete()

	api := modelapi.ToApiFromGraph(g)
	assert.NotNil(api)

	assert.Len(api.PointsGraph, len(g.GetVertices()))
	assert.Len(api.Points2D, 0)
	assert.Len(api.Points3D, 0)

	for _, v := range g.GetVertices() {
		var match *modelapi.PointGraph
		for _, other := range api.PointsGraph {
			if strings.Compare(v.GetId(), other.Id) == 0 {
				match = other
				break
			}
		}
		assert.NotNil(match)
		assert.Len(match.Neighbors, len(v.GetAdjacentVertices()))

		for dest, dist := range v.GetAdjacentVertices() {
			numMatch := 0
			for _, n := range match.Neighbors {
				if strings.Compare(dest.GetId(), n.Id) == 0 {
					numMatch++
					assert.Equal(dist, n.Distance)
				}
			}
			assert.Equal(1, numMatch)
		}
	}
}
