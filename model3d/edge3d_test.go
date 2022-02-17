package model3d_test

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/heustis/lee-tsp-go/model3d"
	"github.com/stretchr/testify/assert"
)

func TestDistanceIncrease(t *testing.T) {
	assert := assert.New(t)

	edge := model3d.NewEdge3D(model3d.NewVertex3D(-1.0, 0.0, -3.0), model3d.NewVertex3D(1.0, 2.0, -1.0))

	testCases := []struct {
		v        *model3d.Vertex3D
		expected float64
	}{
		{model3d.NewVertex3D(0.5, 1.5, -1.5), 0.0},
		{model3d.NewVertex3D(1.0, 1.0, -2.0), math.Sqrt(4+1+1) + math.Sqrt(0+1+1) - math.Sqrt(12)},
		{model3d.NewVertex3D(0.7, -0.5, 0.3), math.Sqrt(1.7*1.7+0.5*0.5+3.3*3.3) + math.Sqrt(0.3*0.3+2.5*2.5+1.3*1.3) - math.Sqrt(12)},
	}

	for i, tc := range testCases {
		assert.InDelta(tc.expected, edge.DistanceIncrease(tc.v), 0.0000001, i)
	}
}

func TestDistanceIncrease3D(t *testing.T) {
	assert := assert.New(t)

	edge := model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(1.0, 0.0, 0.0))

	testCases := []struct {
		v        *model3d.Vertex3D
		expected float64
	}{
		{model3d.NewVertex3D(0.3, 0.0, 0.0), 0.0},
		{model3d.NewVertex3D(1.0, 1.0, 0.0), math.Sqrt2},
		{model3d.NewVertex3D(0.7, 0.5, 0.0), (math.Sqrt(0.74) + math.Sqrt(0.34)) - 1},
	}

	for i, tc := range testCases {
		assert.InDelta(tc.expected, edge.DistanceIncrease(tc.v), 0.0000001, i)
	}
}

func TestEquals_Edge3D(t *testing.T) {
	assert := assert.New(t)

	edge := model3d.NewEdge3D(model3d.NewVertex3D(-1.0, 0.0, -3.0), model3d.NewVertex3D(1.0, 2.0, -1.0))

	var nilEdge *model3d.Edge3D = nil

	assert.True(edge.Equals(edge))
	assert.False(edge.Equals(nil))
	assert.False(edge.Equals((*model3d.Edge3D)(nil)))
	assert.False(nilEdge.Equals(edge))
	assert.True(nilEdge.Equals(nil))
	assert.True(nilEdge.Equals((*model3d.Edge3D)(nil)))
	assert.True(edge.Equals(model3d.NewEdge3D(model3d.NewVertex3D(-1.0, 0.0, -3.0), model3d.NewVertex3D(1.0, 2.0, -1.0))))
	assert.False(edge.Equals(model3d.NewEdge3D(model3d.NewVertex3D(-1.0, 0.0, -2.0), model3d.NewVertex3D(1.0, 2.0, -1.0))))
	assert.False(edge.Equals(model3d.NewEdge3D(model3d.NewVertex3D(-1.0, 0.0, -3.0), model3d.NewVertex3D(2.0, 2.0, -1.0))))
}

func TestGetVector(t *testing.T) {
	assert := assert.New(t)

	edge := model3d.NewEdge3D(model3d.NewVertex3D(-1.0, 0.0, -4.0), model3d.NewVertex3D(1.0, -2.0, -1.0))
	vector := edge.GetVector()
	assert.InDelta(2.0/math.Sqrt(17), vector.X, 0.0000001)
	assert.InDelta(-2.0/math.Sqrt(17), vector.Y, 0.0000001)
	assert.InDelta(3.0/math.Sqrt(17), vector.Z, 0.0000001)
}

func TestIntersects_Edge3D_ShouldWorkFor2DTestCases(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		edgeA    *model3d.Edge3D
		edgeB    *model3d.Edge3D
		expected bool
	}{
		{model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(1.0, 1.0, 0.0)), model3d.NewEdge3D(model3d.NewVertex3D(0.0, -1.0, 0.0), model3d.NewVertex3D(1.0, 0.0, 0.0)), false},
		{model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(1.0, 1.0, 0.0)), model3d.NewEdge3D(model3d.NewVertex3D(-1.0, -1.0, 0.0), model3d.NewVertex3D(1.0, 0.0, 0.0)), false},
		{model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(1.0, 1.0, 0.0)), model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(1.0, 0.0, 0.0)), true},
		{model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(1.0, 1.0, 0.0)), model3d.NewEdge3D(model3d.NewVertex3D(0.5, 0.5, 0.0), model3d.NewVertex3D(-1.5, -1.5, 0.0)), true},
		{model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(1.0, 1.0, 0.0)), model3d.NewEdge3D(model3d.NewVertex3D(0.0, 1.0, 0.0), model3d.NewVertex3D(1.0, 0.0, 0.0)), true},
		{model3d.NewEdge3D(model3d.NewVertex3D(15.5, 12.0, 0.0), model3d.NewVertex3D(10.3, 7.25, 0.0)), model3d.NewEdge3D(model3d.NewVertex3D(5.5, 28.0, 0.0), model3d.NewVertex3D(19.6, 3.0, 0.0)), true},
	}

	for i, tc := range testCases {
		assert.Equal(tc.expected, tc.edgeA.Intersects(tc.edgeB), i)
		assert.Equal(tc.expected, tc.edgeB.Intersects(tc.edgeA), i)

		edgeAReverse := model3d.NewEdge3D(tc.edgeA.End, tc.edgeA.Start)
		edgeBReverse := model3d.NewEdge3D(tc.edgeB.End, tc.edgeB.Start)
		assert.Equal(tc.expected, edgeAReverse.Intersects(tc.edgeB), i)
		assert.Equal(tc.expected, edgeAReverse.Intersects(edgeBReverse), i)
		assert.Equal(tc.expected, edgeBReverse.Intersects(tc.edgeA), i)
		assert.Equal(tc.expected, edgeBReverse.Intersects(edgeAReverse), i)
	}
}

func TestIntersects_Edge3D_ShouldWorkFor3D(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		edgeA    *model3d.Edge3D
		edgeB    *model3d.Edge3D
		expected bool
	}{
		{model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(1.0, 1.0, 1.0)), model3d.NewEdge3D(model3d.NewVertex3D(1.0, 0.0, 0.0), model3d.NewVertex3D(1.0, 1.0, 0.0)), false},
		{model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(1.0, 1.0, 1.0)), model3d.NewEdge3D(model3d.NewVertex3D(2.0, 2.0, 2.0), model3d.NewVertex3D(3.0, 3.0, 3.0)), false},
		{model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(1.0, 1.0, 1.0)), model3d.NewEdge3D(model3d.NewVertex3D(.75, .75, .75), model3d.NewVertex3D(3.0, 3.0, 3.0)), true},
		{model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(1.0, 1.0, 1.0)), model3d.NewEdge3D(model3d.NewVertex3D(.5, .5, .5), model3d.NewVertex3D(-3.0, 15.0, 30.0)), true},
		{model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(1.0, 1.0, 1.0)), model3d.NewEdge3D(model3d.NewVertex3D(-9.5, 5.5, 13.8333), model3d.NewVertex3D(10.5, -4.5, -12.8333)), true},
	}

	for i, tc := range testCases {
		assert.Equal(tc.expected, tc.edgeA.Intersects(tc.edgeB), i)
		assert.Equal(tc.expected, tc.edgeB.Intersects(tc.edgeA), i)

		edgeAReverse := model3d.NewEdge3D(tc.edgeA.End, tc.edgeA.Start)
		edgeBReverse := model3d.NewEdge3D(tc.edgeB.End, tc.edgeB.Start)
		assert.Equal(tc.expected, edgeAReverse.Intersects(tc.edgeB), i)
		assert.Equal(tc.expected, edgeAReverse.Intersects(edgeBReverse), i)
		assert.Equal(tc.expected, edgeBReverse.Intersects(tc.edgeA), i)
		assert.Equal(tc.expected, edgeBReverse.Intersects(edgeAReverse), i)
	}
}

func TestSplit(t *testing.T) {
	assert := assert.New(t)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		start := model3d.NewVertex3D(r.Float64(), 1+r.Float64(), r.Float64())
		end := model3d.NewVertex3D(r.Float64(), 3+r.Float64(), r.Float64())
		edge := model3d.NewEdge3D(start, end)
		assert.Equal(start, edge.GetStart())
		assert.Equal(end, edge.GetEnd())
		assert.Equal(start.DistanceTo(end), edge.GetLength())

		mid := model3d.NewVertex3D(r.Float64(), 2+r.Float64(), r.Float64())
		e1, e2 := edge.Split(mid)

		assert.Equal(start, e1.GetStart())
		assert.Equal(mid, e1.GetEnd())
		assert.Equal(start.DistanceTo(mid), e1.GetLength())

		assert.Equal(mid, e2.GetStart())
		assert.Equal(end, e2.GetEnd())
		assert.Equal(mid.DistanceTo(end), e2.GetLength())

		merged := e1.Merge(e2)
		assert.True(merged.Equals(edge))

		merged2 := e2.Merge(e1)
		assert.True(merged2.GetStart().Equals(mid))
		assert.True(merged2.GetEnd().Equals(mid))
	}
}

func TestString_Edge3D(t *testing.T) {
	assert := assert.New(t)

	edge := model3d.NewEdge3D(model3d.NewVertex3D(-1.2, 0.0, -4.0), model3d.NewVertex3D(1.0, -2.3, -1.4))
	assert.Equal(`{"start":{"x":-1.2,"y":0,"z":-4},"end":{"x":1,"y":-2.3,"z":-1.4}}`, edge.String())
}
