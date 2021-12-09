package model_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/fealos/lee-tsp-go/model3d"
	"github.com/stretchr/testify/assert"
)

func TestDeleteVertex(t *testing.T) {
	assert := assert.New(t)

	vertices := []model.CircuitVertex{
		model2d.NewVertex2D(1, 1),
		model2d.NewVertex2D(2, 2),
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
		model2d.NewVertex2D(8, 8),
	}

	vertices = model.DeleteVertex(vertices, 0)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(2, 2),
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
		model2d.NewVertex2D(8, 8),
	}, vertices)

	vertices = model.DeleteVertex(vertices, 99)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(2, 2),
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
	}, vertices)

	vertices = model.DeleteVertex(vertices, -5)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
	}, vertices)

	vertices = model.DeleteVertex(vertices, 3)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(7, 7),
	}, vertices)

	vertices = model.DeleteVertex(vertices, 0)
	vertices = model.DeleteVertex(vertices, 0)
	vertices = model.DeleteVertex(vertices, 0)
	assert.Len(vertices, 1)
	vertices = model.DeleteVertex(vertices, 0)
	assert.Len(vertices, 0)
	vertices = model.DeleteVertex(vertices, 0)
	assert.Len(vertices, 0)
}

func TestDeleteVertex2(t *testing.T) {
	assert := assert.New(t)

	initVertices := []model.CircuitVertex{
		model2d.NewVertex2D(1, 1),
		model2d.NewVertex2D(2, 2),
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
		model2d.NewVertex2D(8, 8),
	}

	vertices := model.DeleteVertex2(initVertices, initVertices[0])
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(2, 2),
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
		model2d.NewVertex2D(8, 8),
	}, vertices)

	vertices = model.DeleteVertex2(vertices, initVertices[7])
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(2, 2),
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
	}, vertices)

	vertices = model.DeleteVertex2(vertices, initVertices[1])
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
	}, vertices)

	vertices = model.DeleteVertex2(vertices, initVertices[5])
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(7, 7),
	}, vertices)

	// No change on deleting an element that is not in the array
	vertices = model.DeleteVertex2(vertices, initVertices[5])
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(7, 7),
	}, vertices)

	vertices = model.DeleteVertex2(vertices, initVertices[2])
	vertices = model.DeleteVertex2(vertices, initVertices[3])
	vertices = model.DeleteVertex2(vertices, initVertices[4])
	assert.Len(vertices, 1)
	vertices = model.DeleteVertex2(vertices, initVertices[6])
	assert.Len(vertices, 0)
	vertices = model.DeleteVertex2(vertices, initVertices[7])
	assert.Len(vertices, 0)
}

func TestFindClosestEdge_2D(t *testing.T) {
	assert := assert.New(t)

	points := []*model2d.Vertex2D{
		{X: 0.0, Y: 0.0},
		{X: 0.0, Y: 1.0},
		{X: 1.0, Y: 1.0},
		{X: 0.7, Y: 0.5},
		{X: 1.0, Y: 0.0},
	}

	edges := []model.CircuitEdge{
		model2d.NewEdge2D(points[0], points[1]), //0 = 0.0,0.0 -> 0.0,1.0
		model2d.NewEdge2D(points[1], points[2]), //1 = 0.0,1.0 -> 1.0,1.0
		model2d.NewEdge2D(points[2], points[3]), //2 = 1.0,1.0 -> 0.7,0.5
		model2d.NewEdge2D(points[3], points[4]), //3 = 0.7,0.5 -> 1.0,0.0
		model2d.NewEdge2D(points[4], points[0]), //4 = 1.0,0.0 -> 0.0,0.0
	}

	testCases := []struct {
		v        *model2d.Vertex2D
		expected model.CircuitEdge
	}{
		{v: &model2d.Vertex2D{X: 0.0, Y: 0.0}, expected: edges[0]},
		{v: &model2d.Vertex2D{X: 0.5, Y: 0.0}, expected: edges[4]},
		{v: &model2d.Vertex2D{X: 0.5, Y: 0.5}, expected: edges[2]},
		{v: &model2d.Vertex2D{X: 0.5, Y: 0.6}, expected: edges[1]},
		{v: &model2d.Vertex2D{X: 0.6, Y: 0.6}, expected: edges[2]},
		{v: &model2d.Vertex2D{X: 0.5, Y: 0.4}, expected: edges[4]},
		{v: &model2d.Vertex2D{X: 0.6, Y: 0.4}, expected: edges[3]},
		{v: &model2d.Vertex2D{X: 0.2, Y: 0.1}, expected: edges[4]},
	}

	for i, tc := range testCases {
		assert.Equal(model.FindClosestEdge(tc.v, edges), tc.expected, i)
	}
}

func TestFindClosestEdge_2D_ShouldReturnNilIfListIsEmpty(t *testing.T) {
	assert := assert.New(t)

	v := &model2d.Vertex2D{}

	assert.Nil(model.FindClosestEdge(v, []model.CircuitEdge{}))
}

func TestFindClosestEdge_3D(t *testing.T) {
	assert := assert.New(t)

	points := []*model3d.Vertex3D{
		{X: 0.0, Y: 0.0},
		{X: 0.0, Y: 1.0},
		{X: 1.0, Y: 1.0},
		{X: 0.7, Y: 0.5},
		{X: 1.0, Y: 0.0},
	}

	edges := []model.CircuitEdge{
		model3d.NewEdge3D(points[0], points[1]), //0 = 0.0,0.0 -> 0.0,1.0
		model3d.NewEdge3D(points[1], points[2]), //1 = 0.0,1.0 -> 1.0,1.0
		model3d.NewEdge3D(points[2], points[3]), //2 = 1.0,1.0 -> 0.7,0.5
		model3d.NewEdge3D(points[3], points[4]), //3 = 0.7,0.5 -> 1.0,0.0
		model3d.NewEdge3D(points[4], points[0]), //4 = 1.0,0.0 -> 0.0,0.0
	}

	testCases := []struct {
		v        *model3d.Vertex3D
		expected model.CircuitEdge
	}{
		{v: &model3d.Vertex3D{X: 0.0, Y: 0.0, Z: 0.0}, expected: edges[0]},
		{v: &model3d.Vertex3D{X: 0.5, Y: 0.0, Z: 0.0}, expected: edges[4]},
		{v: &model3d.Vertex3D{X: 0.5, Y: 0.5, Z: 0.0}, expected: edges[2]},
		{v: &model3d.Vertex3D{X: 0.5, Y: 0.6, Z: 0.0}, expected: edges[1]},
		{v: &model3d.Vertex3D{X: 0.6, Y: 0.6, Z: 0.0}, expected: edges[2]},
		{v: &model3d.Vertex3D{X: 0.5, Y: 0.4, Z: 0.0}, expected: edges[4]},
		{v: &model3d.Vertex3D{X: 0.6, Y: 0.4, Z: 0.0}, expected: edges[3]},
		{v: &model3d.Vertex3D{X: 0.2, Y: 0.1, Z: 0.0}, expected: edges[4]},
	}

	for i, tc := range testCases {
		assert.Equal(model.FindClosestEdge(tc.v, edges), tc.expected, i)
	}
}

func TestFindClosestEdge_3D_ShouldReturnNilIfListIsEmpty(t *testing.T) {
	assert := assert.New(t)

	v := &model3d.Vertex3D{}

	assert.Nil(model.FindClosestEdge(v, []model.CircuitEdge{}))
}

func TestIndexOfVertex(t *testing.T) {
	assert := assert.New(t)

	initVertices := []model.CircuitVertex{
		model2d.NewVertex2D(1, 1),
		model2d.NewVertex2D(2, 2),
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
		model2d.NewVertex2D(8, 8),
	}

	assert.Equal(-1, model.IndexOfVertex([]model.CircuitVertex{
		model2d.NewVertex2D(2, 2),
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
		model2d.NewVertex2D(8, 8),
	}, initVertices[0]))

	assert.Equal(3, model.IndexOfVertex([]model.CircuitVertex{
		model2d.NewVertex2D(2, 2),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(8, 8),
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
		model2d.NewVertex2D(4, 4),
	}, initVertices[2]))

	assert.Equal(6, model.IndexOfVertex([]model.CircuitVertex{
		model2d.NewVertex2D(2, 2),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(8, 8),
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
		model2d.NewVertex2D(4, 4),
	}, initVertices[3]))

	assert.Equal(0, model.IndexOfVertex([]model.CircuitVertex{
		model2d.NewVertex2D(2, 2),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(8, 8),
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
		model2d.NewVertex2D(4, 4),
	}, initVertices[1]))
}

func TestInsertVertex_Heap(t *testing.T) {
	assert := assert.New(t)
	circuit := []model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
	}

	circuit = model.InsertVertex(circuit, 0, model2d.NewVertex2D(5, 5))
	assert.Len(circuit, 4)
	assert.Equal(model2d.NewVertex2D(5, 5), circuit[0])
	assert.Equal(model2d.NewVertex2D(-15, -15), circuit[1])
	assert.Equal(model2d.NewVertex2D(0, 0), circuit[2])
	assert.Equal(model2d.NewVertex2D(15, -15), circuit[3])

	circuit = model.InsertVertex(circuit, 4, model2d.NewVertex2D(-5, -5))
	assert.Len(circuit, 5)
	assert.Equal(model2d.NewVertex2D(5, 5), circuit[0])
	assert.Equal(model2d.NewVertex2D(-15, -15), circuit[1])
	assert.Equal(model2d.NewVertex2D(0, 0), circuit[2])
	assert.Equal(model2d.NewVertex2D(15, -15), circuit[3])
	assert.Equal(model2d.NewVertex2D(-5, -5), circuit[4])

	circuit = model.InsertVertex(circuit, 2, model2d.NewVertex2D(1, -5))
	assert.Len(circuit, 6)
	assert.Equal(model2d.NewVertex2D(5, 5), circuit[0])
	assert.Equal(model2d.NewVertex2D(-15, -15), circuit[1])
	assert.Equal(model2d.NewVertex2D(1, -5), circuit[2])
	assert.Equal(model2d.NewVertex2D(0, 0), circuit[3])
	assert.Equal(model2d.NewVertex2D(15, -15), circuit[4])
	assert.Equal(model2d.NewVertex2D(-5, -5), circuit[5])
}

func TestIsEdgeCloser_2D(t *testing.T) {
	assert := assert.New(t)

	v := model2d.NewVertex2D(10.0, 10.0)

	testCases := []struct {
		candiate *model2d.Edge2D
		current  *model2d.Edge2D
		expected bool
	}{
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(20.0, 0.0)), model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(20.0, 20.0)), false},
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(20.0, 20.0)), model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(20.0, 0.0)), true},
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 20.0), model2d.NewVertex2D(20.0, 0.0)), model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(20.0, 20.0)), false},
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(20.0, 20.0)), model2d.NewEdge2D(model2d.NewVertex2D(0.0, 20.0), model2d.NewVertex2D(20.0, 0.0)), false},
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(21.0, 0.0)), model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(20.0, 0.0)), true},
		{model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(18.0, 0.0)), model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(20.0, 0.0)), false},
		{model2d.NewEdge2D(model2d.NewVertex2D(2.0, 0.0), model2d.NewVertex2D(22.0, 0.0)), model2d.NewEdge2D(model2d.NewVertex2D(4.0, 0.0), model2d.NewVertex2D(24.0, 0.0)), true},
		{model2d.NewEdge2D(model2d.NewVertex2D(2.0, 0.0), model2d.NewVertex2D(22.0, 0.0)), model2d.NewEdge2D(model2d.NewVertex2D(0.0, 0.0), model2d.NewVertex2D(20.0, 0.0)), false},
	}

	for i, tc := range testCases {
		assert.Equal(tc.expected, model.IsEdgeCloser(v, tc.candiate, tc.current), i)
	}
}

func TestIsEdgeCloser_3D(t *testing.T) {
	assert := assert.New(t)

	v := model3d.NewVertex3D(10.0, 10.0, 0.0)

	testCases := []struct {
		candiate *model3d.Edge3D
		current  *model3d.Edge3D
		expected bool
	}{
		{model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(20.0, 0.0, 0.0)), model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(20.0, 20.0, 0.0)), false},
		{model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(20.0, 20.0, 0.0)), model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(20.0, 0.0, 0.0)), true},
		{model3d.NewEdge3D(model3d.NewVertex3D(0.0, 20.0, 0.0), model3d.NewVertex3D(20.0, 0.0, 0.0)), model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(20.0, 20.0, 0.0)), false},
		{model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(20.0, 20.0, 0.0)), model3d.NewEdge3D(model3d.NewVertex3D(0.0, 20.0, 0.0), model3d.NewVertex3D(20.0, 0.0, 0.0)), false},
		{model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(21.0, 0.0, 0.0)), model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(20.0, 0.0, 0.0)), true},
		{model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(18.0, 0.0, 0.0)), model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(20.0, 0.0, 0.0)), false},
		{model3d.NewEdge3D(model3d.NewVertex3D(2.0, 0.0, 0.0), model3d.NewVertex3D(22.0, 0.0, 0.0)), model3d.NewEdge3D(model3d.NewVertex3D(4.0, 0.0, 0.0), model3d.NewVertex3D(24.0, 0.0, 0.0)), true},
		{model3d.NewEdge3D(model3d.NewVertex3D(2.0, 0.0, 0.0), model3d.NewVertex3D(22.0, 0.0, 0.0)), model3d.NewEdge3D(model3d.NewVertex3D(0.0, 0.0, 0.0), model3d.NewVertex3D(20.0, 0.0, 0.0)), false},
	}

	for i, tc := range testCases {
		assert.Equal(tc.expected, model.IsEdgeCloser(v, tc.candiate, tc.current), i)
	}
}