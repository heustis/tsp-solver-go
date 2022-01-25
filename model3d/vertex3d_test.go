package model3d_test

import (
	"math"
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model3d"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		v1       *model3d.Vertex3D
		v2       *model3d.Vertex3D
		expected *model3d.Vertex3D
	}{
		{model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(0.0, 0.0, 0.0)},
		{model3d.NewVertex3D(1.0, 0.0, -1.0), model3d.NewVertex3D(0.0, 1.0, 0.0), model3d.NewVertex3D(1.0, 1.0, -1.0)},
		{model3d.NewVertex3D(1.0, 1.0, 1.0), model3d.NewVertex3D(-1.0, 10.0, 1.0), model3d.NewVertex3D(0.0, 11.0, 2.0)},
		{model3d.NewVertex3D(-4.0, -15.5, 0.09), model3d.NewVertex3D(5.9, 4.5, 1.11), model3d.NewVertex3D(1.9, -11.0, 1.2)},
	}

	for i, tc := range testCases {
		actual := tc.v1.Add(tc.v2)

		assert.InDelta(tc.expected.X, actual.X, model.Threshold, i)
		assert.InDelta(tc.expected.Y, actual.Y, model.Threshold, i)
		assert.InDelta(tc.expected.Z, actual.Z, model.Threshold, i)

		actualReverse := tc.v2.Add(tc.v1)
		assert.InDelta(tc.expected.X, actualReverse.X, model.Threshold, i)
		assert.InDelta(tc.expected.Y, actualReverse.Y, model.Threshold, i)
		assert.InDelta(tc.expected.Z, actualReverse.Z, model.Threshold, i)
	}
}

func TestDistanceTo(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		v1       *model3d.Vertex3D
		v2       *model3d.Vertex3D
		expected float64
	}{
		{model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(0.0, 0.0, 0.0), 0.0},
		{model3d.NewVertex3D(1.0, 0.0, 0.0), model3d.NewVertex3D(0.0, 0.0, 0.0), 1.0},
		{model3d.NewVertex3D(0.0, 1.0, 0.0), model3d.NewVertex3D(0.0, 0.0, 0.0), 1.0},
		{model3d.NewVertex3D(0.0, 0.0, 1.0), model3d.NewVertex3D(0.0, 0.0, 0.0), 1.0},
		{model3d.NewVertex3D(1.0, 1.0, 1.0), model3d.NewVertex3D(10.0, -1.0, 1.0), math.Sqrt(81 + 4 + 0)},
		{model3d.NewVertex3D(1.0, 1.0, 1.0), model3d.NewVertex3D(10.0, -1.0, 10.0), math.Sqrt(81 + 4 + 81)},
		{model3d.NewVertex3D(-4.0, 0.0, -2.0), model3d.NewVertex3D(0.0, 4.0, 2.0), math.Sqrt(16 + 16 + 16)},
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

	edgeAsc := model3d.NewVertex3D(1.0, 2.0, 1.0).EdgeTo(model3d.NewVertex3D(10.0, 5.0, 1.0)).(*model3d.Edge3D)
	edgeAscReverse := edgeAsc.GetEnd().EdgeTo(edgeAsc.GetStart()).(*model3d.Edge3D)

	edgeDesc := model3d.NewVertex3D(-4.0, 5.0, 1.0).EdgeTo(model3d.NewVertex3D(6.0, -5.0, 1.0)).(*model3d.Edge3D)
	edgeDescReverse := edgeDesc.GetEnd().EdgeTo(edgeDesc.GetStart()).(*model3d.Edge3D)

	testCases := []struct {
		v            *model3d.Vertex3D
		expectedAsc  float64
		expectedDesc float64
	}{
		{v: &model3d.Vertex3D{X: 1.0, Y: 2.0, Z: 1.0}, expectedAsc: 0.0, expectedDesc: 1.4142135623},
		{v: &model3d.Vertex3D{X: 0.0, Y: 0.0, Z: 1.0}, expectedAsc: 1.58113883, expectedDesc: 0.7071067811865},
	}

	for i, tc := range testCases {
		assert.InDelta(tc.v.DistanceToEdge(edgeAsc), tc.expectedAsc, model.Threshold, i)
		assert.InDelta(tc.v.DistanceToEdge(edgeAscReverse), tc.expectedAsc, model.Threshold, i)

		assert.InDelta(tc.v.DistanceToEdge(edgeDesc), tc.expectedDesc, model.Threshold, i)
		assert.InDelta(tc.v.DistanceToEdge(edgeDescReverse), tc.expectedDesc, model.Threshold, i)
	}
}

func TestDotProduct(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		v1       *model3d.Vertex3D
		v2       *model3d.Vertex3D
		expected float64
	}{
		{model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(0.0, 0.0, 0.0), 0.0},
		{model3d.NewVertex3D(1.0, 0.0, -1.0), model3d.NewVertex3D(0.0, 1.0, 0.0), 0.0},
		{model3d.NewVertex3D(1.0, 1.0, 1.0), model3d.NewVertex3D(-1.0, 10.0, 1.0), 10},
		{model3d.NewVertex3D(-4.0, -15.5, 0.09), model3d.NewVertex3D(5.9, -4.5, 1.11), -23.6 + 69.75 + 0.0999},
	}

	for i, tc := range testCases {
		actual := tc.v1.DotProduct(tc.v2)
		assert.InDelta(tc.expected, actual, model.Threshold, i)

		actualReverse := tc.v2.DotProduct(tc.v1)
		assert.InDelta(tc.expected, actualReverse, model.Threshold, i)
	}
}

func TestEquals_Vertex3D(t *testing.T) {
	assert := assert.New(t)

	vertex := model3d.NewVertex3D(-1.1, 0.54321, -3.0)

	var nilVertex *model3d.Vertex3D = nil

	assert.True(vertex.Equals(vertex))
	assert.False(vertex.Equals(nil))
	assert.False(vertex.Equals((*model3d.Edge3D)(nil)))
	assert.False(nilVertex.Equals(vertex))
	assert.True(nilVertex.Equals(nil))
	assert.True(nilVertex.Equals((*model3d.Vertex3D)(nil)))
	assert.True(vertex.Equals(model3d.NewVertex3D(-1.1, 0.54321, -3.0)))
	assert.True(vertex.Equals(model3d.NewVertex3D(-1.1, 0.54321, -3.0+(model.Threshold/10.0))))
	assert.False(vertex.Equals(model3d.NewVertex3D(1.1, 0.54321, -3.0)))
	assert.False(vertex.Equals(model3d.NewVertex3D(-1.1, 1.54321, -3.0)))
	assert.False(vertex.Equals(model3d.NewVertex3D(-1.1, 1.54321, -3.00001)))
}

func TestMultiply(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		v1       *model3d.Vertex3D
		scalar   float64
		expected *model3d.Vertex3D
	}{
		{model3d.NewVertex3D(0.0, 0.0, 0.0), 0, model3d.NewVertex3D(0.0, 0.0, 0.0)},
		{model3d.NewVertex3D(0.0, 0.0, 0.0), 10.0, model3d.NewVertex3D(0.0, 0.0, 0.0)},
		{model3d.NewVertex3D(1.0, 1.0, -1.0), -2.2, model3d.NewVertex3D(-2.2, -2.2, 2.2)},
		{model3d.NewVertex3D(2.0, 4.0, 5.0), 0.25, model3d.NewVertex3D(0.5, 1.0, 1.25)},
		{model3d.NewVertex3D(-4.0, -15.5, 0.09), 3, model3d.NewVertex3D(-12, -46.5, .27)},
	}

	for i, tc := range testCases {
		actual := tc.v1.Multiply(tc.scalar)

		assert.InDelta(tc.expected.X, actual.X, model.Threshold, i)
		assert.InDelta(tc.expected.Y, actual.Y, model.Threshold, i)
		assert.InDelta(tc.expected.Z, actual.Z, model.Threshold, i)
	}
}

func TestProjectToEdge(t *testing.T) {
	assert := assert.New(t)

	edge := model3d.NewVertex3D(1.0, 2.0, 3.0).EdgeTo(model3d.NewVertex3D(10.0, 5.0, 2.5)).(*model3d.Edge3D)
	edgeReverse := edge.GetEnd().EdgeTo(edge.GetStart()).(*model3d.Edge3D)

	testCases := []struct {
		v        *model3d.Vertex3D
		expected *model3d.Vertex3D
	}{
		{v: &model3d.Vertex3D{X: 1.0, Y: 2.0, Z: 3.0}, expected: &model3d.Vertex3D{X: 1.0, Y: 2.0, Z: 3.0}},
		{v: &model3d.Vertex3D{X: 0.0, Y: 0.0, Z: 0.0}, expected: &model3d.Vertex3D{X: -0.34626038781163415, Y: 1.551246537396122, Z: 3.074792243767313}},
	}

	for i, tc := range testCases {
		actual := tc.v.ProjectToEdge(edge)
		assert.InDelta(tc.expected.X, actual.X, model.Threshold, i)
		assert.InDelta(tc.expected.Y, actual.Y, model.Threshold, i)
		assert.InDelta(tc.expected.Z, actual.Z, model.Threshold, i)

		actualReverse := tc.v.ProjectToEdge(edgeReverse)
		assert.InDelta(tc.expected.Z, actualReverse.Z, model.Threshold, i)
		assert.InDelta(tc.expected.Z, actualReverse.Z, model.Threshold, i)
	}
}

func TestString_Vertex3D(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(`{"x":1,"y":-2.3,"z":-1.4}`, model3d.NewVertex3D(1.0, -2.3, -1.4).String())
	assert.Equal(`{"x":0.000000099,"y":-2.30000001,"z":-1.456789123}`, model3d.NewVertex3D(0.000000099, -2.30000001, -1.456789123).String())
}

func TestSubtract(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		v1       *model3d.Vertex3D
		v2       *model3d.Vertex3D
		expected *model3d.Vertex3D
	}{
		{model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(0.0, 0.0, 0.0)},
		{model3d.NewVertex3D(1.0, 0.0, -1.0), model3d.NewVertex3D(0.0, 1.0, 0.0), model3d.NewVertex3D(1.0, -1.0, -1.0)},
		{model3d.NewVertex3D(1.0, 1.0, 1.0), model3d.NewVertex3D(-1.0, 10.0, 1.0), model3d.NewVertex3D(2.0, -9.0, 0.0)},
		{model3d.NewVertex3D(-4.0, -15.5, 0.09), model3d.NewVertex3D(5.9, 4.5, 1.11), model3d.NewVertex3D(-9.9, -20.0, -1.02)},
	}

	for i, tc := range testCases {
		actual := tc.v1.Subtract(tc.v2)

		assert.InDelta(tc.expected.X, actual.X, model.Threshold, i)
		assert.InDelta(tc.expected.Y, actual.Y, model.Threshold, i)
		assert.InDelta(tc.expected.Z, actual.Z, model.Threshold, i)
	}
}
