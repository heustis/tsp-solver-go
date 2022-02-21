package model2d_test

import (
	"fmt"
	"testing"

	"github.com/heustis/tsp-solver-go/model"
	"github.com/heustis/tsp-solver-go/model2d"
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

func TestDotProduct(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		x1       float64
		y1       float64
		x2       float64
		y2       float64
		expected float64
	}{
		{0.0, 0.0, 0.0, 0.0, 0.0},
		{1.5, 0.0, 5.0, 0.5, 7.5},
		{1.0, 2.0, 0.0, 2.25, 4.5},
		{1.0, 1.0, 10.0, -1.0, 9},
		{-4.0, 0.0, 0.0, 4.0, 0.0},
		{-4.0, 2.0, 5.0, 4.0, -12.0},
	}

	for i, tc := range testCases {
		v1 := model2d.NewVertex2D(tc.x1, tc.y1)
		v2 := model2d.NewVertex2D(tc.x2, tc.y2)
		dot := v1.DotProduct(v2)

		assert.InDelta(tc.expected, dot, model.Threshold, fmt.Sprintf("index:%d test:%v actual:%g", i, tc, dot))

		dotReverse := v2.DotProduct(v1)
		assert.InDelta(tc.expected, dotReverse, model.Threshold, fmt.Sprintf("index:%d test:%v actual:%g", i, tc, dotReverse))
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

		assert.InDelta(tc.expected, dist, model.Threshold, fmt.Sprintf("index:%d test:%v actual:%g", i, tc, dist))

		distReverse := v2.DistanceTo(v1)
		assert.InDelta(tc.expected, distReverse, model.Threshold, fmt.Sprintf("index:%d test:%v actual:%g", i, tc, distReverse))
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

func TestMultiply(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		x1        float64
		y1        float64
		scalar    float64
		expectedX float64
		expectedY float64
	}{
		{0.0, 0.0, 0.0, 0.0, 0.0},
		{0.0, 0.0, 10.5, 0.0, 0.0},
		{1.23, 2.34, 0.0, 0.0, 0.0},
		{1.5, 0.0, 5.0, 7.5, 0.0},
		{1.0, 2.0, 0.5, 0.5, 1.0},
		{.456, 1.23, -10.0, -4.56, -12.3},
		{-4.0, -2.2, 2.5, -10.0, -5.5},
		{-3.0, 2.0, -5.0, 15.0, -10.0},
	}

	for i, tc := range testCases {
		v1 := model2d.NewVertex2D(tc.x1, tc.y1)
		mult := v1.Multiply(tc.scalar)

		assert.InDelta(tc.expectedX, mult.X, model.Threshold, fmt.Sprintf("index:%d test:%v actualX:%g", i, tc, mult.X))
		assert.InDelta(tc.expectedY, mult.Y, model.Threshold, fmt.Sprintf("index:%d test:%v actualY:%g", i, tc, mult.Y))
	}
}

func TestPerpendiculars(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		x1            float64
		y1            float64
		expectedLeftX float64
		expectedLeftY float64
	}{
		{0.0, 0.0, 0.0, 0.0},
		{1.5, 0.0, 0.0, 1.5},
		{1.0, 2.5, -2.5, 1.0},
		{-10.0, -1.0, 1.0, -10.0},
	}

	for i, tc := range testCases {
		v1 := model2d.NewVertex2D(tc.x1, tc.y1)
		left := v1.LeftPerpendicular()
		assert.Equal(v1, left.RightPerpendicular(), i)
		assert.Equal(tc.expectedLeftX, left.X, i)
		assert.Equal(tc.expectedLeftY, left.Y, i)
		assert.Equal(0.0, left.DotProduct(v1))

		right := v1.RightPerpendicular()
		assert.Equal(v1, right.LeftPerpendicular(), i)
		assert.Equal(-tc.expectedLeftX, right.X, i)
		assert.Equal(-tc.expectedLeftY, right.Y, i)
		assert.Equal(0.0, right.DotProduct(v1), i)
	}
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

func TestString(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		x1       float64
		y1       float64
		expected string
	}{
		{0.0, 0.0, `{"x":0,"y":0}`},
		{6.0, 5.0, `{"x":6,"y":5}`},
		{1.234, 2.5, `{"x":1.234,"y":2.5}`},
		{-10.005, -1.01234, `{"x":-10.005,"y":-1.01234}`},
		{-0.000005, 1.000000000002345, `{"x":-0.000005,"y":1.000000000002345}`},
	}

	for i, tc := range testCases {
		v1 := model2d.NewVertex2D(tc.x1, tc.y1)
		assert.Equal(tc.expected, v1.String(), i)
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
