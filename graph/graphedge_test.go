package graph_test

import (
	"testing"

	"github.com/heustis/tsp-solver-go/graph"
	"github.com/heustis/tsp-solver-go/model"
	"github.com/stretchr/testify/assert"
)

func TestDistanceIncrease_ShouldReturnTheAmountThePathWillIncreaseByInsertingTheSuppliedVertex(t *testing.T) {
	assert := assert.New(t)

	g := createTestGraphSymmetric()
	defer g.Delete()

	edgeAToC := g.GetVertices()[0].EdgeTo(g.GetVertices()[2])
	assert.InDelta(0.0, edgeAToC.DistanceIncrease(g.GetVertices()[0]), model.Threshold)
	assert.InDelta(0.0, edgeAToC.DistanceIncrease(g.GetVertices()[1]), model.Threshold)
	assert.InDelta(0.0, edgeAToC.DistanceIncrease(g.GetVertices()[2]), model.Threshold)
	assert.InDelta(170.0+110.0-60.0, edgeAToC.DistanceIncrease(g.GetVertices()[3]), model.Threshold)
	assert.InDelta(70.0+10.0-60.0, edgeAToC.DistanceIncrease(g.GetVertices()[4]), model.Threshold)
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

	edge0To7 := g.GetVertices()[0].EdgeTo(g.GetVertices()[7]) // B -> F
	edge0To5 := g.GetVertices()[0].EdgeTo(g.GetVertices()[5]) // B -> J
	edge7To9 := g.GetVertices()[7].EdgeTo(g.GetVertices()[9]) // F -> G -> D
	edge9To2 := g.GetVertices()[9].EdgeTo(g.GetVertices()[2]) // D -> C
	edge1To8 := g.GetVertices()[1].EdgeTo(g.GetVertices()[8]) // H -> G -> D -> A

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
				edge := g.GetVertices()[start].EdgeTo(g.GetVertices()[end])
				assert.NotNil(edge)
				assert.Equal(g.GetVertices()[start], edge.GetStart())
				assert.Equal(g.GetVertices()[end], edge.GetEnd())

				a, b := edge.Split(g.GetVertices()[split])
				assert.Equal(g.GetVertices()[start], a.GetStart())
				assert.Equal(g.GetVertices()[split], a.GetEnd())
				if split != end {
					assert.False(edge.Equals(a))
				} else {
					assert.True(edge.Equals(a))
				}

				assert.Equal(g.GetVertices()[split], b.GetStart())
				assert.Equal(g.GetVertices()[end], b.GetEnd())
				if split != start {
					assert.False(edge.Equals(b))
				} else {
					assert.True(edge.Equals(b))
				}

				merged := a.Merge(b)
				assert.True(edge.Equals(merged))

				mergedReverse := b.Merge(a)
				assert.Equal(g.GetVertices()[split], mergedReverse.GetStart())
				assert.Equal(g.GetVertices()[split], mergedReverse.GetEnd())
				if start == end && (split == start || split == end) {
					assert.True(edge.Equals(mergedReverse))
				} else {
					assert.False(edge.Equals(mergedReverse))
				}

				edge.(*graph.GraphEdge).Delete()
				a.(*graph.GraphEdge).Delete()
				b.(*graph.GraphEdge).Delete()
				merged.(*graph.GraphEdge).Delete()
				mergedReverse.(*graph.GraphEdge).Delete()
			}
		}
	}
}

func TestString_ShouldReturnTheEdgeAsAStringArrayOfIds(t *testing.T) {
	assert := assert.New(t)

	g := createTestGraphSymmetric()
	defer g.Delete()

	edgeAToD := g.GetVertices()[0].EdgeTo(g.GetVertices()[3])
	defer edgeAToD.(*graph.GraphEdge).Delete()

	assert.Equal(`["a","b","c","e","d"]`, edgeAToD.String())

	edgeAToA := g.GetVertices()[0].EdgeTo(g.GetVertices()[0])
	defer edgeAToA.(*graph.GraphEdge).Delete()

	assert.Equal(`["a"]`, edgeAToA.String())
}
