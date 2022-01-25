package model2d_test

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestDistanceIncrease(t *testing.T) {
	assert := assert.New(t)

	edge := model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(1.0, 0.0))

	testCases := []struct {
		v        *model2d.Vertex2D
		expected float64
	}{
		{model2d.NewVertex2D(0.3, 0.0), 0.0},
		{model2d.NewVertex2D(1.0, 1.0), math.Sqrt2},
		{model2d.NewVertex2D(0.7, 0.5), (math.Sqrt(0.74) + math.Sqrt(0.34)) - 1},
	}

	for i, tc := range testCases {
		assert.InDelta(tc.expected, edge.DistanceIncrease(tc.v), 0.0000001, i)
	}
}

func TestEquals_Edge2D(t *testing.T) {
	assert := assert.New(t)

	edge := model2d.NewEdge2D(model2d.NewVertex2D(-3.0, -4.5), model2d.NewVertex2D(1.1, 2.0))
	assert.True(edge.Equals(edge))
	assert.True(edge.Equals(model2d.NewEdge2D(model2d.NewVertex2D(-3.0, -4.5), model2d.NewVertex2D(1.1, 2.0))))
	assert.False(edge.Equals(model2d.NewEdge2D(model2d.NewVertex2D(3.0, -4.5), model2d.NewVertex2D(1.1, 2.0))))
	assert.False(edge.Equals(&model2d.Vertex2D{X: -1.0, Y: -2.0}))
	assert.False(edge.Equals(nil))
}

func TestGetVector(t *testing.T) {
	assert := assert.New(t)

	sqrt2Inv := 1.0 / math.Sqrt2
	sqrt8 := math.Sqrt(8.0)

	testCases := []struct {
		edge      *model2d.Edge2D
		expectedX float64
		expectedY float64
	}{
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(1.0, 1.0)), sqrt2Inv, sqrt2Inv},
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(1.0, -1.0)), sqrt2Inv, -sqrt2Inv},
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(-1.0, 1.0)), -sqrt2Inv, sqrt2Inv},
		{model2d.NewEdge2D(model2d.NewVertex2D(-1.0, 1.0), model2d.NewVertex2D(1.0, -1.0)), 2.0 / sqrt8, -2.0 / sqrt8},
		{model2d.NewEdge2D(model2d.NewVertex2D(15.5, 12.0), model2d.NewVertex2D(10.3, 7.25)), -5.2 / math.Sqrt(5.2*5.2+4.75*4.75), -4.75 / math.Sqrt(5.2*5.2+4.75*4.75)},
	}

	for i, tc := range testCases {
		assert.InDelta(tc.expectedX, tc.edge.GetVector().X, model.Threshold, i)
		assert.InDelta(tc.expectedY, tc.edge.GetVector().Y, model.Threshold, i)

		reverse := tc.edge.End.EdgeTo(tc.edge.Start).(*model2d.Edge2D)
		assert.InDelta(-tc.expectedX, reverse.GetVector().X, model.Threshold, i)
		assert.InDelta(-tc.expectedY, reverse.GetVector().Y, model.Threshold, i)
	}
}

func TestIntersects_Edge2D(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		edgeA    *model2d.Edge2D
		edgeB    *model2d.Edge2D
		expected bool
	}{
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(1.0, 1.0)), model2d.NewEdge2D(model2d.NewVertex2D(0.0, -1.0), model2d.NewVertex2D(1.0, 0.0)), false},
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(1.0, 1.0)), model2d.NewEdge2D(model2d.NewVertex2D(-1.0, -1.0), model2d.NewVertex2D(1.0, 0.0)), false},
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(1.0, 1.0)), model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(1.0, 0.0)), true},
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(1.0, 1.0)), model2d.NewEdge2D(model2d.NewVertex2D(0.5, 0.5), model2d.NewVertex2D(-1.5, -1.5)), true},
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(1.0, 1.0)), model2d.NewEdge2D(model2d.NewVertex2D(0.0, 1.0), model2d.NewVertex2D(1.0, 0.0)), true},
		{model2d.NewEdge2D(model2d.NewVertex2D(15.5, 12.0), model2d.NewVertex2D(10.3, 7.25)), model2d.NewEdge2D(model2d.NewVertex2D(5.5, 28.0), model2d.NewVertex2D(19.6, 3.0)), true},
	}

	for i, tc := range testCases {
		assert.Equal(tc.expected, tc.edgeA.Intersects(tc.edgeB), i)
		assert.Equal(tc.expected, tc.edgeB.Intersects(tc.edgeA), i)

		edgeAReverse := model2d.NewEdge2D(tc.edgeA.End, tc.edgeA.Start)
		edgeBReverse := model2d.NewEdge2D(tc.edgeB.End, tc.edgeB.Start)
		assert.Equal(tc.expected, edgeAReverse.Intersects(tc.edgeB), i)
		assert.Equal(tc.expected, edgeAReverse.Intersects(edgeBReverse), i)
		assert.Equal(tc.expected, edgeBReverse.Intersects(tc.edgeA), i)
		assert.Equal(tc.expected, edgeBReverse.Intersects(edgeAReverse), i)
	}
}

func TestMerge_Edge2D(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		edgeA *model2d.Edge2D
		edgeB *model2d.Edge2D
	}{
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(1.0, 1.0)), model2d.NewEdge2D(model2d.NewVertex2D(0.0, -1.0), model2d.NewVertex2D(1.0, 0.0))},
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(1.0, 1.0)), model2d.NewEdge2D(model2d.NewVertex2D(-1.0, -1.0), model2d.NewVertex2D(1.0, 0.0))},
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(1.0, 1.0)), model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(1.0, 0.0))},
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(1.0, 1.0)), model2d.NewEdge2D(model2d.NewVertex2D(0.5, 0.5), model2d.NewVertex2D(-1.5, -1.5))},
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(1.0, 1.0)), model2d.NewEdge2D(model2d.NewVertex2D(0.0, 1.0), model2d.NewVertex2D(1.0, 0.0))},
		{model2d.NewEdge2D(model2d.NewVertex2D(15.5, 12.0), model2d.NewVertex2D(10.3, 7.25)), model2d.NewEdge2D(model2d.NewVertex2D(5.5, 28.0), model2d.NewVertex2D(19.6, 3.0))},
	}

	for _, tc := range testCases {
		merged := tc.edgeA.Merge(tc.edgeB)
		assert.Equal(tc.edgeA.GetStart(), merged.GetStart())
		assert.Equal(tc.edgeB.GetEnd(), merged.GetEnd())

		mergedReverse := tc.edgeB.Merge(tc.edgeA)
		assert.Equal(tc.edgeB.GetStart(), mergedReverse.GetStart())
		assert.Equal(tc.edgeA.GetEnd(), mergedReverse.GetEnd())
	}
}

func TestSplit_Edge2D(t *testing.T) {
	assert := assert.New(t)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		start := model2d.NewVertex2D(r.Float64(), 1+r.Float64())
		end := model2d.NewVertex2D(r.Float64(), 3+r.Float64())
		edge := model2d.NewEdge2D(start, end)
		assert.Equal(start, edge.GetStart())
		assert.Equal(end, edge.GetEnd())
		assert.Equal(start.DistanceTo(end), edge.GetLength())

		mid := model2d.NewVertex2D(r.Float64(), 2+r.Float64())
		e1, e2 := edge.Split(mid)

		assert.Equal(start, e1.GetStart())
		assert.Equal(mid, e1.GetEnd())
		assert.Equal(start.DistanceTo(mid), e1.GetLength())

		assert.Equal(mid, e2.GetStart())
		assert.Equal(end, e2.GetEnd())
		assert.Equal(mid.DistanceTo(end), e2.GetLength())
	}
}

func TestString_Edge2D(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		edge     *model2d.Edge2D
		expected string
	}{
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(1.0, 1.0)), `{"start":{"x":0,"y":0},"end":{"x":1,"y":1}}`},
		{model2d.NewEdge2D(model2d.NewVertex2D(0.000000000002345, 0.0), model2d.NewVertex2D(1.0, 1.0)), `{"start":{"x":0.000000000002345,"y":0},"end":{"x":1,"y":1}}`},
		{model2d.NewEdge2D(model2d.NewVertex2D(-0.123, 0.00001), model2d.NewVertex2D(123.45, -1.987)), `{"start":{"x":-0.123,"y":0.00001},"end":{"x":123.45,"y":-1.987}}`},
		{model2d.NewEdge2D(model2d.NewVertex2D(15.5, 12.0), model2d.NewVertex2D(10.3, 7.25)), `{"start":{"x":15.5,"y":12},"end":{"x":10.3,"y":7.25}}`},
	}

	for i, tc := range testCases {
		assert.Equal(tc.expected, tc.edge.String(), i)
	}
}
