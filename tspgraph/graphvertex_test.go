package tspgraph_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/tspgraph"
	"github.com/stretchr/testify/assert"
)

func TestDelete_ShouldClearAdjacentVertices(t *testing.T) {
	assert := assert.New(t)

	gen := &tspgraph.GraphGenerator{
		MaxEdges:    4,
		MinEdges:    2,
		NumVertices: uint32(10),
	}

	g := gen.Create()
	assert.NotNil(g)
	defer g.Delete()

	for _, v := range g.Vertices {
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
			distance := g.Vertices[start].DistanceTo(g.Vertices[end])
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
			edge, okay := g.Vertices[start].EdgeTo(g.Vertices[end]).(*tspgraph.GraphEdge)
			assert.True(okay)
			assert.NotNil(edge)
			assert.Equal(g.Vertices[start], edge.GetStart())
			assert.Equal(g.Vertices[end], edge.GetEnd())
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
		current := g.Vertices[i]

		g2 := g.ToApi().ToGraph()
		assert.True(current.Equals(g2.Vertices[i]))
		g2.Delete()

		assert.False(current.Equals(nil))
		assert.False(current.Equals(g))
		assert.False(current.Equals(g.Vertices))

		for j := 0; j < 5; j++ {
			other := g.Vertices[j]
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
		paths := g.Vertices[start].PathToAll()
		for end := 0; end < 5; end++ {
			edge, okay := paths[g.Vertices[end]].(*tspgraph.GraphEdge)
			assert.True(okay)
			assert.NotNil(edge)
			assert.Equal(g.Vertices[start], edge.GetStart())
			assert.Equal(g.Vertices[end], edge.GetEnd())
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

func createTestGraphSymmetric() *tspgraph.Graph {
	api := &tspgraph.GraphApi{
		Vertices: []*tspgraph.GraphVertexApi{
			{
				Id:               "a",
				AdjacentVertices: make(map[string]float64),
			},
			{
				Id:               "b",
				AdjacentVertices: make(map[string]float64),
			},
			{
				Id:               "c",
				AdjacentVertices: make(map[string]float64),
			},
			{
				Id:               "d",
				AdjacentVertices: make(map[string]float64),
			},
			{
				Id:               "e",
				AdjacentVertices: make(map[string]float64),
			},
		},
	}
	api.Vertices[0].AdjacentVertices["b"] = 10
	api.Vertices[0].AdjacentVertices["c"] = 100
	api.Vertices[0].AdjacentVertices["d"] = 1000

	api.Vertices[1].AdjacentVertices["a"] = 10
	api.Vertices[1].AdjacentVertices["c"] = 50
	api.Vertices[1].AdjacentVertices["e"] = 1000

	api.Vertices[2].AdjacentVertices["a"] = 100
	api.Vertices[2].AdjacentVertices["b"] = 50
	api.Vertices[2].AdjacentVertices["d"] = 500
	api.Vertices[2].AdjacentVertices["e"] = 10

	api.Vertices[3].AdjacentVertices["a"] = 1000
	api.Vertices[3].AdjacentVertices["c"] = 500
	api.Vertices[3].AdjacentVertices["e"] = 100

	api.Vertices[4].AdjacentVertices["b"] = 1000
	api.Vertices[4].AdjacentVertices["c"] = 10
	api.Vertices[4].AdjacentVertices["d"] = 100

	return api.ToGraph()
}
