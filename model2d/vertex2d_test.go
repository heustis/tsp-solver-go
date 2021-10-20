package model2d

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		x1        float64
		y1        float64
		x2        float64
		y2        float64
		expectedX float64
		expectedY float64
	}{
		{0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		{1.0, 0.0, 0.0, 0.0, 1.0, 0.0},
		{1.0, 1.0, 10.0, -1.0, 11.0, 0.0},
		{-4.0, -15.5, 5.9, 4.5, 1.9, -11.0},
	}

	for i, tc := range testCases {
		v1 := NewVertex2D(tc.x1, tc.y1)
		v2 := NewVertex2D(tc.x2, tc.y2)
		actual := v1.add(v2)

		assert.InDelta(tc.expectedX, actual.X, model.Threshold, i)
		assert.InDelta(tc.expectedY, actual.Y, model.Threshold, i)

		actualReverse := v2.add(v1)
		assert.InDelta(tc.expectedX, actualReverse.X, model.Threshold, i)
		assert.InDelta(tc.expectedY, actualReverse.Y, model.Threshold, i)
	}
}

func TestDistanceTo(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		x1       float64
		y1       float64
		x2       float64
		y2       float64
		expected float64
	}{
		{0.0, 0.0, 0.0, 0.0, 0.0},
		{1.0, 0.0, 0.0, 0.0, 1.0},
		{0.0, 1.0, 0.0, 0.0, 1.0},
		{0.0, 0.0, 1.0, 0.0, 1.0},
		{0.0, 0.0, 0.0, 1.0, 1.0},
		{1.0, 1.0, 10.0, -1.0, 9.2195444572928873},
		{-4.0, 0.0, 0.0, 4.0, 5.65685424949238},
	}

	for i, tc := range testCases {
		v1 := NewVertex2D(tc.x1, tc.y1)
		v2 := NewVertex2D(tc.x2, tc.y2)
		dist := v1.DistanceTo(v2)

		assert.InDelta(tc.expected, dist, model.Threshold, i, tc, dist)

		distReverse := v2.DistanceTo(v1)
		assert.InDelta(tc.expected, distReverse, model.Threshold, i, tc, distReverse)
	}
}

func TestDistanceToEdge(t *testing.T) {
	assert := assert.New(t)

	edgeAsc := NewEdge2D(NewVertex2D(1.0, 2.0), NewVertex2D(10.0, 5.0))
	edgeAscReverse := NewEdge2D(edgeAsc.End, edgeAsc.Start)

	edgeDesc := NewEdge2D(NewVertex2D(-4.0, 5.0), NewVertex2D(6.0, -5.0))
	edgeDescReverse := NewEdge2D(edgeDesc.End, edgeDesc.Start)

	testCases := []struct {
		v            *Vertex2D
		expectedAsc  float64
		expectedDesc float64
	}{
		{v: &Vertex2D{X: 1.0, Y: 2.0}, expectedAsc: 0.0, expectedDesc: 1.4142135623},
		{v: &Vertex2D{X: 0.0, Y: 0.0}, expectedAsc: 1.58113883, expectedDesc: 0.7071067811865},
	}

	for i, tc := range testCases {
		assert.InDelta(tc.v.distanceToEdge(edgeAsc), tc.expectedAsc, model.Threshold, i)
		assert.InDelta(tc.v.distanceToEdge(edgeAscReverse), tc.expectedAsc, model.Threshold, i)

		assert.InDelta(tc.v.distanceToEdge(edgeDesc), tc.expectedDesc, model.Threshold, i)
		assert.InDelta(tc.v.distanceToEdge(edgeDescReverse), tc.expectedDesc, model.Threshold, i)
	}
}

func TestFindClosestEdge(t *testing.T) {
	assert := assert.New(t)

	points := []*Vertex2D{
		{X: 0.0, Y: 0.0},
		{X: 0.0, Y: 1.0},
		{X: 1.0, Y: 1.0},
		{X: 0.7, Y: 0.5},
		{X: 1.0, Y: 0.0},
	}

	edges := []model.CircuitEdge{
		NewEdge2D(points[0], points[1]), //0 = 0.0,0.0 -> 0.0,1.0
		NewEdge2D(points[1], points[2]), //1 = 0.0,1.0 -> 1.0,1.0
		NewEdge2D(points[2], points[3]), //2 = 1.0,1.0 -> 0.7,0.5
		NewEdge2D(points[3], points[4]), //3 = 0.7,0.5 -> 1.0,0.0
		NewEdge2D(points[4], points[0]), //4 = 1.0,0.0 -> 0.0,0.0
	}

	testCases := []struct {
		v        *Vertex2D
		expected model.CircuitEdge
	}{
		{v: &Vertex2D{X: 0.0, Y: 0.0}, expected: edges[0]},
		{v: &Vertex2D{X: 0.5, Y: 0.0}, expected: edges[4]},
		{v: &Vertex2D{X: 0.5, Y: 0.5}, expected: edges[2]},
		{v: &Vertex2D{X: 0.5, Y: 0.6}, expected: edges[1]},
		{v: &Vertex2D{X: 0.6, Y: 0.6}, expected: edges[2]},
		{v: &Vertex2D{X: 0.5, Y: 0.4}, expected: edges[4]},
		{v: &Vertex2D{X: 0.6, Y: 0.4}, expected: edges[3]},
		{v: &Vertex2D{X: 0.2, Y: 0.1}, expected: edges[4]},
	}

	for i, tc := range testCases {
		assert.Equal(tc.v.FindClosestEdge(edges), tc.expected, i)
	}
}

func TestFindClosestEdge_ShouldReturnNilIfListIsEmpty(t *testing.T) {
	assert := assert.New(t)

	v := &Vertex2D{}

	assert.Nil(v.FindClosestEdge([]model.CircuitEdge{}))
}

func TestIsEdgeCloser(t *testing.T) {
	assert := assert.New(t)

	v := NewVertex2D(10.0, 10.0)

	testCases := []struct {
		candiate *Edge2D
		current  *Edge2D
		expected bool
	}{
		{NewEdge2D(NewVertex2D(0.0, 0.0), NewVertex2D(20.0, 0.0)), NewEdge2D(NewVertex2D(0.0, 0.0), NewVertex2D(20.0, 20.0)), false},
		{NewEdge2D(NewVertex2D(0.0, 0.0), NewVertex2D(20.0, 20.0)), NewEdge2D(NewVertex2D(0.0, 0.0), NewVertex2D(20.0, 0.0)), true},
		{NewEdge2D(NewVertex2D(0.0, 20.0), NewVertex2D(20.0, 0.0)), NewEdge2D(NewVertex2D(0.0, 0.0), NewVertex2D(20.0, 20.0)), false},
		{NewEdge2D(NewVertex2D(0.0, 0.0), NewVertex2D(20.0, 20.0)), NewEdge2D(NewVertex2D(0.0, 20.0), NewVertex2D(20.0, 0.0)), false},
		{NewEdge2D(NewVertex2D(0.0, 0.0), NewVertex2D(21.0, 0.0)), NewEdge2D(NewVertex2D(0.0, 0.0), NewVertex2D(20.0, 0.0)), true},
		{NewEdge2D(NewVertex2D(0.0, 0.0), NewVertex2D(18.0, 0.0)), NewEdge2D(NewVertex2D(0.0, 0.0), NewVertex2D(20.0, 0.0)), false},
		{NewEdge2D(NewVertex2D(2.0, 0.0), NewVertex2D(22.0, 0.0)), NewEdge2D(NewVertex2D(4.0, 0.0), NewVertex2D(24.0, 0.0)), true},
		{NewEdge2D(NewVertex2D(2.0, 0.0), NewVertex2D(22.0, 0.0)), NewEdge2D(NewVertex2D(0.0, 0.0), NewVertex2D(20.0, 0.0)), false},
	}

	for i, tc := range testCases {
		assert.Equal(tc.expected, v.IsEdgeCloser(tc.candiate, tc.current), i)
	}
}

func TestProjectToEdge(t *testing.T) {
	assert := assert.New(t)

	edge := NewEdge2D(NewVertex2D(1.0, 2.0), NewVertex2D(10.0, 5.0))
	edgeReverse := NewEdge2D(edge.End, edge.Start)

	testCases := []struct {
		v         *Vertex2D
		expectedX float64
		expectedY float64
	}{
		{v: &Vertex2D{X: 1.0, Y: 2.0}, expectedX: 1.0, expectedY: 2.0},
		{v: &Vertex2D{X: 0.0, Y: 0.0}, expectedX: -.5, expectedY: 1.5},
	}

	for i, tc := range testCases {
		actual := tc.v.projectToEdge(edge)
		assert.InDelta(tc.expectedX, actual.X, model.Threshold, i)
		assert.InDelta(tc.expectedY, actual.Y, model.Threshold, i)

		actualReverse := tc.v.projectToEdge(edgeReverse)
		assert.InDelta(tc.expectedX, actualReverse.X, model.Threshold, i)
		assert.InDelta(tc.expectedY, actualReverse.Y, model.Threshold, i)
	}
}

func TestIsLeftOf(t *testing.T) {
	assert := assert.New(t)

	edge := NewEdge2D(NewVertex2D(1.0, 2.0), NewVertex2D(10.0, 5.0))

	testCases := []struct {
		v             *Vertex2D
		expectedLeft  bool
		expectedRight bool
	}{
		{v: &Vertex2D{X: 1.0, Y: 2.0}, expectedLeft: false, expectedRight: false},
		{v: &Vertex2D{X: -2.0, Y: 1.0}, expectedLeft: false, expectedRight: false},
		{v: &Vertex2D{X: 0.0, Y: 0.0}, expectedLeft: false, expectedRight: true},
		{v: &Vertex2D{X: 1.0, Y: 10.0}, expectedLeft: true, expectedRight: false},
	}

	for i, tc := range testCases {
		assert.Equal(tc.expectedLeft, tc.v.isLeftOf(edge), i)
		assert.Equal(tc.expectedRight, tc.v.isRightOf(edge), i)
	}
}

func TestSubtract(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		x1        float64
		y1        float64
		x2        float64
		y2        float64
		expectedX float64
		expectedY float64
	}{
		{0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		{1.0, 0.0, 0.0, 0.0, 1.0, 0.0},
		{1.0, 1.0, 10.0, -1.0, -9.0, 2.0},
		{-4.0, -15.5, 5.9, 4.5, -9.9, -20.0},
	}

	for i, tc := range testCases {
		v1 := NewVertex2D(tc.x1, tc.y1)
		v2 := NewVertex2D(tc.x2, tc.y2)
		actual := v1.subtract(v2)

		assert.InDelta(tc.expectedX, actual.X, model.Threshold, i)
		assert.InDelta(tc.expectedY, actual.Y, model.Threshold, i)
	}
}
