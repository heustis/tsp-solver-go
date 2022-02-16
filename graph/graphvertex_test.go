package graph_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/graph"
	"github.com/fealos/lee-tsp-go/modelapi"
	"github.com/stretchr/testify/assert"
)

func TestDelete_ShouldClearAdjacentVertices(t *testing.T) {
	assert := assert.New(t)

	gen := &graph.GraphGenerator{
		MaxEdges:    4,
		MinEdges:    2,
		NumVertices: uint32(10),
	}

	g := gen.Create()
	assert.NotNil(g)
	defer g.Delete()

	for _, v := range g.GetVertices() {
		v.Delete()
		assert.Nil(v.GetAdjacentVertices())
	}
}

func TestDistanceTo_ShouldReturnShortestPath(t *testing.T) {
	assert := assert.New(t)

	g := createTestGraphSymmetric()
	defer g.Delete()

	expectedDistances := [][]float64{
		{0.0, 10.0, 60.0, 170.0, 70.0},
		{10, 0.0, 50.0, 160.0, 60.0},
		{60.0, 50.0, 0.0, 110.0, 10.0},
		{170.0, 160.0, 110.0, 0.0, 100.0},
		{70.0, 60.0, 10.0, 100.0, 0.0},
	}

	for start := 0; start < 5; start++ {
		for end := 0; end < 5; end++ {
			distance := g.GetVertices()[start].DistanceTo(g.GetVertices()[end])
			assert.Equal(expectedDistances[start][end], distance)
		}
	}
}

func TestEdgeTo_ShouldReturnShortestPath(t *testing.T) {
	assert := assert.New(t)

	g := createTestGraphSymmetric()
	defer g.Delete()

	expectedDistances := [][]float64{
		{0.0, 10.0, 60.0, 170.0, 70.0},
		{10, 0.0, 50.0, 160.0, 60.0},
		{60.0, 50.0, 0.0, 110.0, 10.0},
		{170.0, 160.0, 110.0, 0.0, 100.0},
		{70.0, 60.0, 10.0, 100.0, 0.0},
	}

	expectedLengths := [][]int{
		{1, 2, 3, 5, 4},
		{2, 1, 2, 4, 3},
		{3, 2, 1, 3, 2},
		{5, 4, 3, 1, 2},
		{4, 3, 2, 2, 1},
	}

	for start := 0; start < 5; start++ {
		for end := 0; end < 5; end++ {
			edge, okay := g.GetVertices()[start].EdgeTo(g.GetVertices()[end]).(*graph.GraphEdge)
			assert.True(okay)
			assert.NotNil(edge)
			assert.Equal(g.GetVertices()[start], edge.GetStart())
			assert.Equal(g.GetVertices()[end], edge.GetEnd())
			assert.Equal(expectedDistances[start][end], edge.GetLength())
			assert.Len(edge.GetPath(), expectedLengths[start][end])
			edge.Delete()
			assert.Nil(edge.GetPath())
		}
	}
}

func TestEquals_ShouldCompareIds(t *testing.T) {
	assert := assert.New(t)

	g := createTestGraphSymmetric()
	defer g.Delete()

	for i := 0; i < 5; i++ {
		current := g.GetVertices()[i]

		g2 := modelapi.ToApiFromGraph(g).ToGraph()
		assert.True(current.Equals(g2.GetVertices()[i]))
		g2.Delete()

		assert.False(current.Equals(nil))
		assert.False(current.Equals(g))
		assert.False(current.Equals(g.GetVertices()))

		for j := 0; j < 5; j++ {
			other := g.GetVertices()[j]
			if i == j {
				assert.True(current.Equals(other))
			} else {
				assert.False(current.Equals(other))
			}
		}
	}
}

func TestPathToAll_ShouldProduceOptimalPaths(t *testing.T) {
	assert := assert.New(t)

	g := createTestGraphSymmetric()
	defer g.Delete()

	expectedLengths := [][]int{
		{1, 2, 3, 5, 4},
		{2, 1, 2, 4, 3},
		{3, 2, 1, 3, 2},
		{5, 4, 3, 1, 2},
		{4, 3, 2, 2, 1},
	}

	expectedDistances := [][]float64{
		{0.0, 10.0, 60.0, 170.0, 70.0},
		{10, 0.0, 50.0, 160.0, 60.0},
		{60.0, 50.0, 0.0, 110.0, 10.0},
		{170.0, 160.0, 110.0, 0.0, 100.0},
		{70.0, 60.0, 10.0, 100.0, 0.0},
	}

	for start := 0; start < 5; start++ {
		paths := g.GetVertices()[start].GetPaths()
		for end := 0; end < 5; end++ {
			edge := paths[g.GetVertices()[end]]
			assert.NotNil(edge)
			assert.Equal(g.GetVertices()[start], edge.GetStart())
			assert.Equal(g.GetVertices()[end], edge.GetEnd())
			assert.Equal(expectedDistances[start][end], edge.GetLength())
			assert.Len(edge.GetPath(), expectedLengths[start][end])
			edge.Delete()
			assert.Nil(edge.GetPath())
		}

		for v := range paths {
			delete(paths, v)
		}
	}
}

func createTestGraphSymmetric() *graph.Graph {
	api := &modelapi.TspRequest{
		PointsGraph: []*modelapi.PointGraph{
			{
				Id: "a",
				Neighbors: []modelapi.PointGraphNeighbor{
					{Id: "b", Distance: 10},
					{Id: "c", Distance: 100},
					{Id: "d", Distance: 1000},
				},
			},
			{
				Id: "b",
				Neighbors: []modelapi.PointGraphNeighbor{
					{Id: "a", Distance: 10},
					{Id: "c", Distance: 50},
					{Id: "e", Distance: 1000},
				},
			},
			{
				Id: "c",
				Neighbors: []modelapi.PointGraphNeighbor{
					{Id: "a", Distance: 100},
					{Id: "b", Distance: 50},
					{Id: "d", Distance: 500},
					{Id: "e", Distance: 10},
				},
			},
			{
				Id: "d",
				Neighbors: []modelapi.PointGraphNeighbor{
					{Id: "a", Distance: 1000},
					{Id: "c", Distance: 500},
					{Id: "e", Distance: 100},
				},
			},
			{
				Id: "e",
				Neighbors: []modelapi.PointGraphNeighbor{
					{Id: "b", Distance: 1000},
					{Id: "c", Distance: 10},
					{Id: "d", Distance: 100},
				},
			},
		},
	}

	return api.ToGraph()
}
