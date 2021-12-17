package graph_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/graph"
	"github.com/fealos/lee-tsp-go/model"
	"github.com/stretchr/testify/assert"
)

func TestNewGraphEdge_ShouldReturnShortestPath(t *testing.T) {
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
			edge := graph.NewGraphEdge(g.Vertices[start], g.Vertices[end])
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

func TestDistanceIncrease_ShouldReturnTheAmountThePathWillIncreaseByInsertingTheSuppliedVertex(t *testing.T) {
	assert := assert.New(t)

	g := createTestGraphSymmetric()
	defer g.Delete()

	edgeAToC := graph.NewGraphEdge(g.Vertices[0], g.Vertices[2])
	assert.InDelta(0.0, edgeAToC.DistanceIncrease(g.Vertices[0]), model.Threshold)
	assert.InDelta(0.0, edgeAToC.DistanceIncrease(g.Vertices[1]), model.Threshold)
	assert.InDelta(0.0, edgeAToC.DistanceIncrease(g.Vertices[2]), model.Threshold)
	assert.InDelta(170.0+110.0-60.0, edgeAToC.DistanceIncrease(g.Vertices[3]), model.Threshold)
	assert.InDelta(70.0+10.0-60.0, edgeAToC.DistanceIncrease(g.Vertices[4]), model.Threshold)
}

func TestIntersects_ShouldReturnTrueIfTheEdgesContainACommonPoint(t *testing.T) {
	assert := assert.New(t)

	seed := int64(1)
	gen := &graph.GraphGenerator{
		EnableAsymetricDistances: true,
		MaxEdges:                 6,
		MinEdges:                 3,
		NumVertices:              10,
		Seed:                     &seed,
	}

	g := gen.Create()
	defer g.Delete()

	edge0To7 := graph.NewGraphEdge(g.Vertices[0], (g.Vertices[7])) // B -> F
	edge0To5 := graph.NewGraphEdge(g.Vertices[0], (g.Vertices[5])) // B -> J
	edge7To9 := graph.NewGraphEdge(g.Vertices[7], (g.Vertices[9])) // F -> G -> D
	edge9To2 := graph.NewGraphEdge(g.Vertices[9], (g.Vertices[2])) // D -> C
	edge1To8 := graph.NewGraphEdge(g.Vertices[1], (g.Vertices[8])) // H -> G -> D -> A

	assert.NotNil(edge0To7)
	assert.NotNil(edge0To5)
	assert.NotNil(edge7To9)
	assert.NotNil(edge9To2)
	assert.NotNil(edge1To8)

	assert.True(edge0To7.Intersects(edge0To7))
	assert.True(edge0To7.Intersects(edge0To5))
	assert.True(edge0To7.Intersects(edge7To9))
	assert.False(edge0To7.Intersects(edge9To2))
	assert.False(edge0To7.Intersects(edge1To8))

	assert.True(edge0To5.Intersects(edge0To7))
	assert.True(edge0To5.Intersects(edge0To5))
	assert.False(edge0To5.Intersects(edge7To9))
	assert.False(edge0To5.Intersects(edge9To2))
	assert.False(edge0To5.Intersects(edge1To8))

	assert.True(edge7To9.Intersects(edge0To7))
	assert.False(edge7To9.Intersects(edge0To5))
	assert.True(edge7To9.Intersects(edge7To9))
	assert.True(edge7To9.Intersects(edge9To2))
	assert.True(edge7To9.Intersects(edge1To8))

	assert.False(edge1To8.Intersects(edge0To7))
	assert.False(edge1To8.Intersects(edge0To5))
	assert.True(edge1To8.Intersects(edge7To9))
	assert.True(edge1To8.Intersects(edge9To2))
	assert.True(edge1To8.Intersects(edge1To8))
}

func TestMergeAndSplit(t *testing.T) {
	assert := assert.New(t)

	g := createTestGraphSymmetric()
	defer g.Delete()

	for start := 0; start < 5; start++ {
		for end := 0; end < 5; end++ {
			for split := 0; split < 5; split++ {
				edge := graph.NewGraphEdge(g.Vertices[start], g.Vertices[end])
				assert.NotNil(edge)
				assert.Equal(g.Vertices[start], edge.GetStart())
				assert.Equal(g.Vertices[end], edge.GetEnd())

				a, b := edge.Split(g.Vertices[split])
				assert.Equal(g.Vertices[start], a.GetStart())
				assert.Equal(g.Vertices[split], a.GetEnd())
				assert.False(edge.Equals(a))

				assert.Equal(g.Vertices[split], b.GetStart())
				assert.Equal(g.Vertices[end], b.GetEnd())
				assert.False(edge.Equals(b))

				merged := a.Merge(b)
				assert.True(edge.Equals(merged))

				mergedReverse := b.Merge(a)
				assert.Equal(g.Vertices[split], mergedReverse.GetStart())
				assert.Equal(g.Vertices[split], mergedReverse.GetEnd())
				assert.False(edge.Equals(mergedReverse))

				edge.Delete()
				a.(*graph.GraphEdge).Delete()
				b.(*graph.GraphEdge).Delete()
				merged.(*graph.GraphEdge).Delete()
				mergedReverse.(*graph.GraphEdge).Delete()
			}
		}
	}
}

func TestToString_ShouldReturnTheEdgeAsAStringArrayOfIds(t *testing.T) {
	assert := assert.New(t)

	g := createTestGraphSymmetric()
	defer g.Delete()

	edgeAToD := graph.NewGraphEdge(g.Vertices[0], g.Vertices[3])
	defer edgeAToD.Delete()

	assert.Equal(`["a","b","c","e","d"]`, edgeAToD.ToString())

	edgeAToA := graph.NewGraphEdge(g.Vertices[0], g.Vertices[0])
	defer edgeAToA.Delete()

	assert.Equal(`["a"]`, edgeAToA.ToString())
}
