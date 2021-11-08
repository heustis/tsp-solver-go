package model3d

import (
	"math"
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		v1       *Vertex3D
		v2       *Vertex3D
		expected *Vertex3D
	}{
		{NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(0.0, 0.0, 0.0)},
		{NewVertex3D(1.0, 0.0, -1.0), NewVertex3D(0.0, 1.0, 0.0), NewVertex3D(1.0, 1.0, -1.0)},
		{NewVertex3D(1.0, 1.0, 1.0), NewVertex3D(-1.0, 10.0, 1.0), NewVertex3D(0.0, 11.0, 2.0)},
		{NewVertex3D(-4.0, -15.5, 0.09), NewVertex3D(5.9, 4.5, 1.11), NewVertex3D(1.9, -11.0, 1.2)},
	}

	for i, tc := range testCases {
		actual := tc.v1.add(tc.v2)

		assert.InDelta(tc.expected.X, actual.X, model.Threshold, i)
		assert.InDelta(tc.expected.Y, actual.Y, model.Threshold, i)
		assert.InDelta(tc.expected.Z, actual.Z, model.Threshold, i)

		actualReverse := tc.v2.add(tc.v1)
		assert.InDelta(tc.expected.X, actualReverse.X, model.Threshold, i)
		assert.InDelta(tc.expected.Y, actualReverse.Y, model.Threshold, i)
		assert.InDelta(tc.expected.Z, actualReverse.Z, model.Threshold, i)
	}
}

func TestDistanceTo(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		v1       *Vertex3D
		v2       *Vertex3D
		expected float64
	}{
		{NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(0.0, 0.0, 0.0), 0.0},
		{NewVertex3D(1.0, 0.0, 0.0), NewVertex3D(0.0, 0.0, 0.0), 1.0},
		{NewVertex3D(0.0, 1.0, 0.0), NewVertex3D(0.0, 0.0, 0.0), 1.0},
		{NewVertex3D(0.0, 0.0, 1.0), NewVertex3D(0.0, 0.0, 0.0), 1.0},
		{NewVertex3D(1.0, 1.0, 1.0), NewVertex3D(10.0, -1.0, 1.0), math.Sqrt(81 + 4 + 0)},
		{NewVertex3D(1.0, 1.0, 1.0), NewVertex3D(10.0, -1.0, 10.0), math.Sqrt(81 + 4 + 81)},
		{NewVertex3D(-4.0, 0.0, -2.0), NewVertex3D(0.0, 4.0, 2.0), math.Sqrt(16 + 16 + 16)},
	}

	for i, tc := range testCases {
		dist := tc.v1.DistanceTo(tc.v2)
		assert.InDelta(tc.expected, dist, model.Threshold, i, tc, dist)

		distReverse := tc.v2.DistanceTo(tc.v1)
		assert.InDelta(tc.expected, distReverse, model.Threshold, i, tc, distReverse)
	}
}

func TestDistanceToEdge(t *testing.T) {
	assert := assert.New(t)

	edgeAsc := NewEdge3D(NewVertex3D(1.0, 2.0, 1.0), NewVertex3D(10.0, 5.0, 1.0))
	edgeAscReverse := NewEdge3D(edgeAsc.End, edgeAsc.Start)

	edgeDesc := NewEdge3D(NewVertex3D(-4.0, 5.0, 1.0), NewVertex3D(6.0, -5.0, 1.0))
	edgeDescReverse := NewEdge3D(edgeDesc.End, edgeDesc.Start)

	testCases := []struct {
		v            *Vertex3D
		expectedAsc  float64
		expectedDesc float64
	}{
		{v: &Vertex3D{X: 1.0, Y: 2.0, Z: 1.0}, expectedAsc: 0.0, expectedDesc: 1.4142135623},
		{v: &Vertex3D{X: 0.0, Y: 0.0, Z: 1.0}, expectedAsc: 1.58113883, expectedDesc: 0.7071067811865},
	}

	for i, tc := range testCases {
		assert.InDelta(tc.v.distanceToEdge(edgeAsc), tc.expectedAsc, model.Threshold, i)
		assert.InDelta(tc.v.distanceToEdge(edgeAscReverse), tc.expectedAsc, model.Threshold, i)

		assert.InDelta(tc.v.distanceToEdge(edgeDesc), tc.expectedDesc, model.Threshold, i)
		assert.InDelta(tc.v.distanceToEdge(edgeDescReverse), tc.expectedDesc, model.Threshold, i)
	}
}

func TestEquals_Vertex3D(t *testing.T) {
	assert := assert.New(t)

	vertex := NewVertex3D(-1.1, 0.54321, -3.0)

	var nilVertex *Vertex3D = nil

	assert.True(vertex.Equals(vertex))
	assert.False(vertex.Equals(nil))
	assert.False(vertex.Equals((*Edge3D)(nil)))
	assert.False(nilVertex.Equals(vertex))
	assert.True(nilVertex.Equals(nil))
	assert.True(nilVertex.Equals((*Vertex3D)(nil)))
	assert.True(vertex.Equals(NewVertex3D(-1.1, 0.54321, -3.0)))
	assert.True(vertex.Equals(NewVertex3D(-1.1, 0.54321, -3.0+(model.Threshold/10.0))))
	assert.False(vertex.Equals(NewVertex3D(1.1, 0.54321, -3.0)))
	assert.False(vertex.Equals(NewVertex3D(-1.1, 1.54321, -3.0)))
	assert.False(vertex.Equals(NewVertex3D(-1.1, 1.54321, -3.00001)))
}

func TestFindClosestEdge(t *testing.T) {
	assert := assert.New(t)

	points := []*Vertex3D{
		{X: 0.0, Y: 0.0},
		{X: 0.0, Y: 1.0},
		{X: 1.0, Y: 1.0},
		{X: 0.7, Y: 0.5},
		{X: 1.0, Y: 0.0},
	}

	edges := []model.CircuitEdge{
		NewEdge3D(points[0], points[1]), //0 = 0.0,0.0 -> 0.0,1.0
		NewEdge3D(points[1], points[2]), //1 = 0.0,1.0 -> 1.0,1.0
		NewEdge3D(points[2], points[3]), //2 = 1.0,1.0 -> 0.7,0.5
		NewEdge3D(points[3], points[4]), //3 = 0.7,0.5 -> 1.0,0.0
		NewEdge3D(points[4], points[0]), //4 = 1.0,0.0 -> 0.0,0.0
	}

	testCases := []struct {
		v        *Vertex3D
		expected model.CircuitEdge
	}{
		{v: &Vertex3D{X: 0.0, Y: 0.0, Z: 0.0}, expected: edges[0]},
		{v: &Vertex3D{X: 0.5, Y: 0.0, Z: 0.0}, expected: edges[4]},
		{v: &Vertex3D{X: 0.5, Y: 0.5, Z: 0.0}, expected: edges[2]},
		{v: &Vertex3D{X: 0.5, Y: 0.6, Z: 0.0}, expected: edges[1]},
		{v: &Vertex3D{X: 0.6, Y: 0.6, Z: 0.0}, expected: edges[2]},
		{v: &Vertex3D{X: 0.5, Y: 0.4, Z: 0.0}, expected: edges[4]},
		{v: &Vertex3D{X: 0.6, Y: 0.4, Z: 0.0}, expected: edges[3]},
		{v: &Vertex3D{X: 0.2, Y: 0.1, Z: 0.0}, expected: edges[4]},
	}

	for i, tc := range testCases {
		assert.Equal(tc.v.FindClosestEdge(edges), tc.expected, i)
	}
}

func TestFindClosestEdge_ShouldReturnNilIfListIsEmpty(t *testing.T) {
	assert := assert.New(t)

	v := &Vertex3D{}

	assert.Nil(v.FindClosestEdge([]model.CircuitEdge{}))
}

func TestIsEdgeCloser(t *testing.T) {
	assert := assert.New(t)

	v := NewVertex3D(10.0, 10.0, 0.0)

	testCases := []struct {
		candiate *Edge3D
		current  *Edge3D
		expected bool
	}{
		{NewEdge3D(NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(20.0, 0.0, 0.0)), NewEdge3D(NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(20.0, 20.0, 0.0)), false},
		{NewEdge3D(NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(20.0, 20.0, 0.0)), NewEdge3D(NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(20.0, 0.0, 0.0)), true},
		{NewEdge3D(NewVertex3D(0.0, 20.0, 0.0), NewVertex3D(20.0, 0.0, 0.0)), NewEdge3D(NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(20.0, 20.0, 0.0)), false},
		{NewEdge3D(NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(20.0, 20.0, 0.0)), NewEdge3D(NewVertex3D(0.0, 20.0, 0.0), NewVertex3D(20.0, 0.0, 0.0)), false},
		{NewEdge3D(NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(21.0, 0.0, 0.0)), NewEdge3D(NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(20.0, 0.0, 0.0)), true},
		{NewEdge3D(NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(18.0, 0.0, 0.0)), NewEdge3D(NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(20.0, 0.0, 0.0)), false},
		{NewEdge3D(NewVertex3D(2.0, 0.0, 0.0), NewVertex3D(22.0, 0.0, 0.0)), NewEdge3D(NewVertex3D(4.0, 0.0, 0.0), NewVertex3D(24.0, 0.0, 0.0)), true},
		{NewEdge3D(NewVertex3D(2.0, 0.0, 0.0), NewVertex3D(22.0, 0.0, 0.0)), NewEdge3D(NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(20.0, 0.0, 0.0)), false},
	}

	for i, tc := range testCases {
		assert.Equal(tc.expected, v.IsEdgeCloser(tc.candiate, tc.current), i)
	}
}

func TestProjectToEdge(t *testing.T) {
	assert := assert.New(t)

	edge := NewEdge3D(NewVertex3D(1.0, 2.0, 3.0), NewVertex3D(10.0, 5.0, 2.5))
	edgeReverse := NewEdge3D(edge.End, edge.Start)

	testCases := []struct {
		v        *Vertex3D
		expected *Vertex3D
	}{
		{v: &Vertex3D{X: 1.0, Y: 2.0, Z: 3.0}, expected: &Vertex3D{X: 1.0, Y: 2.0, Z: 3.0}},
		{v: &Vertex3D{X: 0.0, Y: 0.0, Z: 0.0}, expected: &Vertex3D{X: -0.34626038781163415, Y: 1.551246537396122, Z: 3.074792243767313}},
	}

	for i, tc := range testCases {
		actual := tc.v.projectToEdge(edge)
		assert.InDelta(tc.expected.X, actual.X, model.Threshold, i)
		assert.InDelta(tc.expected.Y, actual.Y, model.Threshold, i)
		assert.InDelta(tc.expected.Z, actual.Z, model.Threshold, i)

		actualReverse := tc.v.projectToEdge(edgeReverse)
		assert.InDelta(tc.expected.Z, actualReverse.Z, model.Threshold, i)
		assert.InDelta(tc.expected.Z, actualReverse.Z, model.Threshold, i)
	}
}

func TestSubtract(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		v1       *Vertex3D
		v2       *Vertex3D
		expected *Vertex3D
	}{
		{NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(0.0, 0.0, 0.0), NewVertex3D(0.0, 0.0, 0.0)},
		{NewVertex3D(1.0, 0.0, -1.0), NewVertex3D(0.0, 1.0, 0.0), NewVertex3D(1.0, -1.0, -1.0)},
		{NewVertex3D(1.0, 1.0, 1.0), NewVertex3D(-1.0, 10.0, 1.0), NewVertex3D(2.0, -9.0, 0.0)},
		{NewVertex3D(-4.0, -15.5, 0.09), NewVertex3D(5.9, 4.5, 1.11), NewVertex3D(-9.9, -20.0, -1.02)},
	}

	for i, tc := range testCases {
		actual := tc.v1.subtract(tc.v2)

		assert.InDelta(tc.expected.X, actual.X, model.Threshold, i)
		assert.InDelta(tc.expected.Y, actual.Y, model.Threshold, i)
		assert.InDelta(tc.expected.Z, actual.Z, model.Threshold, i)
	}
}

func TestToString_Vertex3D(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(`{"x":1,"y":-2.3,"z":-1.4}`, NewVertex3D(1.0, -2.3, -1.4).ToString())
}
