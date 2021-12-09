package model3d

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDistanceIncrease(t *testing.T) {
	assert := assert.New(t)

	edge := NewEdge3D(NewVertex3D(-1.0, 0.0, -3.0), NewVertex3D(1.0, 2.0, -1.0))

	testCases := []struct {
		v        *Vertex3D
		expected float64
	}{
		{NewVertex3D(0.5, 1.5, -1.5), 0.0},
		{NewVertex3D(1.0, 1.0, -2.0), math.Sqrt(4+1+1) + math.Sqrt(0+1+1) - math.Sqrt(12)},
		{NewVertex3D(0.7, -0.5, 0.3), math.Sqrt(1.7*1.7+0.5*0.5+3.3*3.3) + math.Sqrt(0.3*0.3+2.5*2.5+1.3*1.3) - math.Sqrt(12)},
	}

	for i, tc := range testCases {
		assert.InDelta(tc.expected, edge.DistanceIncrease(tc.v), 0.0000001, i)
	}
}

func TestDistanceIncrease2D(t *testing.T) {
	assert := assert.New(t)

	edge := NewEdge3D(NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(1.0, 0.0, 0.0))

	testCases := []struct {
		v        *Vertex3D
		expected float64
	}{
		{NewVertex3D(0.3, 0.0, 0.0), 0.0},
		{NewVertex3D(1.0, 1.0, 0.0), math.Sqrt2},
		{NewVertex3D(0.7, 0.5, 0.0), (math.Sqrt(0.74) + math.Sqrt(0.34)) - 1},
	}

	for i, tc := range testCases {
		assert.InDelta(tc.expected, edge.DistanceIncrease(tc.v), 0.0000001, i)
	}
}

func TestEquals_Edge3D(t *testing.T) {
	assert := assert.New(t)

	edge := NewEdge3D(NewVertex3D(-1.0, 0.0, -3.0), NewVertex3D(1.0, 2.0, -1.0))

	var nilEdge *Edge3D = nil

	assert.True(edge.Equals(edge))
	assert.False(edge.Equals(nil))
	assert.False(edge.Equals((*Edge3D)(nil)))
	assert.False(nilEdge.Equals(edge))
	assert.True(nilEdge.Equals(nil))
	assert.True(nilEdge.Equals((*Edge3D)(nil)))
	assert.True(edge.Equals(NewEdge3D(NewVertex3D(-1.0, 0.0, -3.0), NewVertex3D(1.0, 2.0, -1.0))))
	assert.False(edge.Equals(NewEdge3D(NewVertex3D(-1.0, 0.0, -2.0), NewVertex3D(1.0, 2.0, -1.0))))
	assert.False(edge.Equals(NewEdge3D(NewVertex3D(-1.0, 0.0, -3.0), NewVertex3D(2.0, 2.0, -1.0))))
}

func TestGetVector(t *testing.T) {
	assert := assert.New(t)

	edge := NewEdge3D(NewVertex3D(-1.0, 0.0, -4.0), NewVertex3D(1.0, -2.0, -1.0))
	vector := edge.GetVector()
	assert.InDelta(2.0/math.Sqrt(17), vector.X, 0.0000001)
	assert.InDelta(-2.0/math.Sqrt(17), vector.Y, 0.0000001)
	assert.InDelta(3.0/math.Sqrt(17), vector.Z, 0.0000001)
}

func TestSplit(t *testing.T) {
	assert := assert.New(t)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		start := NewVertex3D(r.Float64(), 1+r.Float64(), r.Float64())
		end := NewVertex3D(r.Float64(), 3+r.Float64(), r.Float64())
		edge := NewEdge3D(start, end)
		assert.Equal(start, edge.GetStart())
		assert.Equal(end, edge.GetEnd())
		assert.Equal(start.DistanceTo(end), edge.GetLength())

		mid := NewVertex3D(r.Float64(), 2+r.Float64(), r.Float64())
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

func TestToString_Edge3D(t *testing.T) {
	assert := assert.New(t)

	edge := NewEdge3D(NewVertex3D(-1.2, 0.0, -4.0), NewVertex3D(1.0, -2.3, -1.4))
	assert.Equal(`{"start":{"x":-1.2,"y":0,"z":-4},"end":{"x":1,"y":-2.3,"z":-1.4}}`, edge.ToString())
}
