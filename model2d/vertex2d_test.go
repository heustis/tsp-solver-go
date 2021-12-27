package model2d_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
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
		v1 := model2d.NewVertex2D(tc.x1, tc.y1)
		v2 := model2d.NewVertex2D(tc.x2, tc.y2)
		actual := v1.Add(v2)

		assert.InDelta(tc.expectedX, actual.X, model.Threshold, i)
		assert.InDelta(tc.expectedY, actual.Y, model.Threshold, i)

		actualReverse := v2.Add(v1)
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
		v1 := model2d.NewVertex2D(tc.x1, tc.y1)
		v2 := model2d.NewVertex2D(tc.x2, tc.y2)
		dist := v1.DistanceTo(v2)

		assert.InDelta(tc.expected, dist, model.Threshold, i, tc, dist)

		distReverse := v2.DistanceTo(v1)
		assert.InDelta(tc.expected, distReverse, model.Threshold, i, tc, distReverse)
	}
}

func TestDistanceToEdge(t *testing.T) {
	assert := assert.New(t)

	edgeAsc := model2d.NewEdge2D(model2d.NewVertex2D(1.0, 2.0), model2d.NewVertex2D(10.0, 5.0))
	edgeAscReverse := model2d.NewEdge2D(edgeAsc.End, edgeAsc.Start)

	edgeDesc := model2d.NewEdge2D(model2d.NewVertex2D(-4.0, 5.0), model2d.NewVertex2D(6.0, -5.0))
	edgeDescReverse := model2d.NewEdge2D(edgeDesc.End, edgeDesc.Start)

	testCases := []struct {
		v            *model2d.Vertex2D
		expectedAsc  float64
		expectedDesc float64
	}{
		{v: &model2d.Vertex2D{X: 1.0, Y: 2.0}, expectedAsc: 0.0, expectedDesc: 1.4142135623},
		{v: &model2d.Vertex2D{X: 0.0, Y: 0.0}, expectedAsc: 1.58113883, expectedDesc: 0.7071067811865},
	}

	for i, tc := range testCases {
		assert.InDelta(tc.v.DistanceToEdge(edgeAsc), tc.expectedAsc, model.Threshold, i)
		assert.InDelta(tc.v.DistanceToEdge(edgeAscReverse), tc.expectedAsc, model.Threshold, i)

		assert.InDelta(tc.v.DistanceToEdge(edgeDesc), tc.expectedDesc, model.Threshold, i)
		assert.InDelta(tc.v.DistanceToEdge(edgeDescReverse), tc.expectedDesc, model.Threshold, i)
	}
}

func TestEdgeTo(t *testing.T) {
	assert := assert.New(t)

	edgeAsc := model2d.NewEdge2D(model2d.NewVertex2D(1.0, 2.0), model2d.NewVertex2D(10.0, 5.0))
	edgeAscReverse := model2d.NewEdge2D(edgeAsc.End, edgeAsc.Start)
	edgeToAsc := edgeAsc.Start.EdgeTo(edgeAsc.End)
	edgeToAscReverse := edgeAsc.End.EdgeTo(edgeAsc.Start)

	assert.Equal(edgeAsc, edgeToAsc)
	assert.Equal(edgeAscReverse, edgeToAscReverse)
}

func TestEquals_Vertex2D(t *testing.T) {
	assert := assert.New(t)

	v1 := &model2d.Vertex2D{X: 1.0, Y: 2.0}
	assert.True(v1.Equals(v1))
	assert.True(v1.Equals(&model2d.Vertex2D{X: 1.0, Y: 2.0}))
	assert.False(v1.Equals(&model2d.Vertex2D{X: 2.0, Y: 1.0}))
	assert.False(v1.Equals(&model2d.Vertex2D{X: 1.0, Y: -2.0}))
	assert.False(v1.Equals(&model2d.Vertex2D{X: -1.0, Y: 2.0}))
	assert.False(v1.Equals(&model2d.Vertex2D{X: -1.0, Y: -2.0}))
	assert.False(v1.Equals(nil))
}

func TestProjectToEdge(t *testing.T) {
	assert := assert.New(t)

	edge := model2d.NewEdge2D(model2d.NewVertex2D(1.0, 2.0), model2d.NewVertex2D(10.0, 5.0))
	edgeReverse := model2d.NewEdge2D(edge.End, edge.Start)

	testCases := []struct {
		v         *model2d.Vertex2D
		expectedX float64
		expectedY float64
	}{
		{v: &model2d.Vertex2D{X: 1.0, Y: 2.0}, expectedX: 1.0, expectedY: 2.0},
		{v: &model2d.Vertex2D{X: 0.0, Y: 0.0}, expectedX: -.5, expectedY: 1.5},
	}

	for i, tc := range testCases {
		actual := tc.v.ProjectToEdge(edge)
		assert.InDelta(tc.expectedX, actual.X, model.Threshold, i)
		assert.InDelta(tc.expectedY, actual.Y, model.Threshold, i)

		actualReverse := tc.v.ProjectToEdge(edgeReverse)
		assert.InDelta(tc.expectedX, actualReverse.X, model.Threshold, i)
		assert.InDelta(tc.expectedY, actualReverse.Y, model.Threshold, i)
	}
}

func TestIsLeftOf(t *testing.T) {
	assert := assert.New(t)

	edge := model2d.NewEdge2D(model2d.NewVertex2D(1.0, 2.0), model2d.NewVertex2D(10.0, 5.0))

	testCases := []struct {
		v             *model2d.Vertex2D
		expectedLeft  bool
		expectedRight bool
	}{
		{v: &model2d.Vertex2D{X: 1.0, Y: 2.0}, expectedLeft: false, expectedRight: false},
		{v: &model2d.Vertex2D{X: -2.0, Y: 1.0}, expectedLeft: false, expectedRight: false},
		{v: &model2d.Vertex2D{X: 0.0, Y: 0.0}, expectedLeft: false, expectedRight: true},
		{v: &model2d.Vertex2D{X: 1.0, Y: 10.0}, expectedLeft: true, expectedRight: false},
	}

	for i, tc := range testCases {
		assert.Equal(tc.expectedLeft, tc.v.IsLeftOf(edge), i)
		assert.Equal(tc.expectedRight, tc.v.IsRightOf(edge), i)
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
		v1 := model2d.NewVertex2D(tc.x1, tc.y1)
		v2 := model2d.NewVertex2D(tc.x2, tc.y2)
		actual := v1.Subtract(v2)

		assert.InDelta(tc.expectedX, actual.X, model.Threshold, i)
		assert.InDelta(tc.expectedY, actual.Y, model.Threshold, i)
	}
}
