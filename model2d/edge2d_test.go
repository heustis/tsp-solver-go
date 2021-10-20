package model2d

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistanceIncrease(t *testing.T) {
	assert := assert.New(t)

	edge := NewEdge2D(NewVertex2D(0.0, 0.0), NewVertex2D(1.0, 0.0))

	testCases := []struct {
		v        *Vertex2D
		expected float64
	}{
		{NewVertex2D(0.3, 0.0), 0.0},
		{NewVertex2D(1.0, 1.0), math.Sqrt2},
		{NewVertex2D(0.7, 0.5), (math.Sqrt(0.74) + math.Sqrt(0.34)) - 1},
	}

	for i, tc := range testCases {
		assert.InDelta(tc.expected, edge.DistanceIncrease(tc.v), 0.0000001, i)
	}
}

func TestSplit(t *testing.T) {
	assert := assert.New(t)

	for i := 0; i < 10; i++ {
		start := NewVertex2D(rand.Float64(), 1+rand.Float64())
		end := NewVertex2D(rand.Float64(), 3+rand.Float64())
		edge := NewEdge2D(start, end)
		assert.Equal(start, edge.GetStart())
		assert.Equal(end, edge.GetEnd())
		assert.Equal(start.DistanceTo(end), edge.GetLength())

		mid := NewVertex2D(rand.Float64(), 2+rand.Float64())
		e1, e2 := edge.Split(mid)

		assert.Equal(start, e1.GetStart())
		assert.Equal(mid, e1.GetEnd())
		assert.Equal(start.DistanceTo(mid), e1.GetLength())

		assert.Equal(mid, e2.GetStart())
		assert.Equal(end, e2.GetEnd())
		assert.Equal(mid.DistanceTo(end), e2.GetLength())
	}
}
