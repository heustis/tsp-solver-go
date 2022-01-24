package tspmodel_test

import (
	"container/list"
	"math/rand"
	"testing"
	"time"

	"github.com/fealos/lee-tsp-go/tspmodel2d"
	"github.com/fealos/lee-tsp-go/tspmodel3d"
	"github.com/stretchr/testify/assert"
)

func BenchmarkFindClosestEdge(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	vertices := tspmodel2d.GenerateVertices(b.N * 10)
	edges, _ := tspmodel2d.BuildPerimiter(vertices)
	edgesList := list.New()
	for _, v := range edges {
		edgesList.PushBack(v)
	}

	//BenchmarkFindClosestEdge/FindClosestEdge-16         	15805778	        75.28 ns/op	       0 B/op	       0 allocs/op
	b.Run("FindClosestEdge", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tspmodel.FindClosestEdge(vertices[r.Intn(len(vertices))], edges)
		}
	})
	//BenchmarkFindClosestEdge/FindClosestEdgeList-16     	 7892976	       149.5 ns/op	       0 B/op	       0 allocs/op
	b.Run("FindClosestEdgeList", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tspmodel.FindClosestEdgeList(vertices[r.Intn(len(vertices))], edgesList)
		}
	})
}

func BenchmarkDeleteVertex(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	//BenchmarkDeleteVertex/DeleteVertex-16         	   31269	    168491 ns/op	    1080 B/op	      10 allocs/op
	b.Run("DeleteVertex", func(b *testing.B) {
		vertices := tspmodel2d.GenerateVertices(b.N * 10)
		for i := 0; i < b.N; i++ {
			vertices = tspmodel.DeleteVertex(vertices, r.Intn(len(vertices)))
		}
	})
	//BenchmarkDeleteVertex/DeleteVertexCopy-16     	   10000	   1328849 ns/op	 1525177 B/op	      11 allocs/op
	b.Run("DeleteVertexCopy", func(b *testing.B) {
		vertices := tspmodel2d.GenerateVertices(b.N * 10)
		for i := 0; i < b.N; i++ {
			vertices = tspmodel.DeleteVertexCopy(vertices, vertices[r.Intn(len(vertices))])
		}
	})
}

func TestDeduplicateVerticesNoSorting(t *testing.T) {
	assert := assert.New(t)

	init := []tspmodel.CircuitVertex{
		tspmodel3d.NewVertex3D(-15, -15, -5.0),
		tspmodel3d.NewVertex3D(0, 0, tspmodel.Threshold/9.0),
		tspmodel3d.NewVertex3D(15, -15, -5.0),
		tspmodel3d.NewVertex3D(-15-tspmodel.Threshold/3.0, -15, -5),
		tspmodel3d.NewVertex3D(0, 0, 0.0),
		tspmodel3d.NewVertex3D(0, tspmodel.Threshold*2, 0.0),
		tspmodel3d.NewVertex3D(-15+tspmodel.Threshold/3.0, -15-tspmodel.Threshold/3.0, -5+tspmodel.Threshold/4),
		tspmodel3d.NewVertex3D(3, 0, 3),
		tspmodel3d.NewVertex3D(3, 13, 4),
		tspmodel3d.NewVertex3D(7, 6, 5),
		tspmodel3d.NewVertex3D(-7, 6, 5),
	}

	actual := tspmodel.DeduplicateVerticesNoSorting(init)
	assert.ElementsMatch([]*tspmodel3d.Vertex3D{
		tspmodel3d.NewVertex3D(-15, -15, -5),
		tspmodel3d.NewVertex3D(-7, 6, 5),
		tspmodel3d.NewVertex3D(0, 0, tspmodel.Threshold/9.0),
		tspmodel3d.NewVertex3D(0, tspmodel.Threshold*2, 0.0),
		tspmodel3d.NewVertex3D(3, 0, 3),
		tspmodel3d.NewVertex3D(3, 13, 4),
		tspmodel3d.NewVertex3D(7, 6, 5),
		tspmodel3d.NewVertex3D(15, -15, -5),
	}, actual)
}

func TestDeleteVertex(t *testing.T) {
	assert := assert.New(t)

	vertices := []tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(1, 1),
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(8, 8),
	}

	vertices = tspmodel.DeleteVertex(vertices, 0)
	assert.Equal([]tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(8, 8),
	}, vertices)

	vertices = tspmodel.DeleteVertex(vertices, 99)
	assert.Equal([]tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
	}, vertices)

	vertices = tspmodel.DeleteVertex(vertices, -5)
	assert.Equal([]tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
	}, vertices)

	vertices = tspmodel.DeleteVertex(vertices, 3)
	assert.Equal([]tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(7, 7),
	}, vertices)

	vertices = tspmodel.DeleteVertex(vertices, 0)
	vertices = tspmodel.DeleteVertex(vertices, 0)
	vertices = tspmodel.DeleteVertex(vertices, 0)
	assert.Len(vertices, 1)
	vertices = tspmodel.DeleteVertex(vertices, 0)
	assert.Len(vertices, 0)
	vertices = tspmodel.DeleteVertex(vertices, 0)
	assert.Len(vertices, 0)
}

func TestDeleteVertexCopy(t *testing.T) {
	assert := assert.New(t)

	initVertices := []tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(1, 1),
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(8, 8),
	}

	vertices := tspmodel.DeleteVertexCopy(initVertices, initVertices[0])
	assert.Equal([]tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(8, 8),
	}, vertices)

	vertices = tspmodel.DeleteVertexCopy(vertices, initVertices[7])
	assert.Equal([]tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
	}, vertices)

	vertices = tspmodel.DeleteVertexCopy(vertices, initVertices[1])
	assert.Equal([]tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
	}, vertices)

	vertices = tspmodel.DeleteVertexCopy(vertices, initVertices[5])
	assert.Equal([]tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(7, 7),
	}, vertices)

	// No change on deleting an element that is not in the array
	vertices = tspmodel.DeleteVertexCopy(vertices, initVertices[5])
	assert.Equal([]tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(7, 7),
	}, vertices)

	vertices = tspmodel.DeleteVertexCopy(vertices, initVertices[2])
	vertices = tspmodel.DeleteVertexCopy(vertices, initVertices[3])
	vertices = tspmodel.DeleteVertexCopy(vertices, initVertices[4])
	assert.Len(vertices, 1)
	vertices = tspmodel.DeleteVertexCopy(vertices, initVertices[6])
	assert.Len(vertices, 0)
	vertices = tspmodel.DeleteVertexCopy(vertices, initVertices[7])
	assert.Len(vertices, 0)
}

func TestFindClosestEdge_2D(t *testing.T) {
	assert := assert.New(t)

	points := []*tspmodel2d.Vertex2D{
		{X: 0.0, Y: 0.0},
		{X: 0.0, Y: 1.0},
		{X: 1.0, Y: 1.0},
		{X: 0.7, Y: 0.5},
		{X: 1.0, Y: 0.0},
	}

	edges := []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(points[0], points[1]), //0 = 0.0,0.0 -> 0.0,1.0
		tspmodel2d.NewEdge2D(points[1], points[2]), //1 = 0.0,1.0 -> 1.0,1.0
		tspmodel2d.NewEdge2D(points[2], points[3]), //2 = 1.0,1.0 -> 0.7,0.5
		tspmodel2d.NewEdge2D(points[3], points[4]), //3 = 0.7,0.5 -> 1.0,0.0
		tspmodel2d.NewEdge2D(points[4], points[0]), //4 = 1.0,0.0 -> 0.0,0.0
	}

	testCases := []struct {
		v        *tspmodel2d.Vertex2D
		expected tspmodel.CircuitEdge
	}{
		{v: &tspmodel2d.Vertex2D{X: 0.0, Y: 0.0}, expected: edges[0]},
		{v: &tspmodel2d.Vertex2D{X: 0.5, Y: 0.0}, expected: edges[4]},
		{v: &tspmodel2d.Vertex2D{X: 0.5, Y: 0.5}, expected: edges[2]},
		{v: &tspmodel2d.Vertex2D{X: 0.5, Y: 0.6}, expected: edges[1]},
		{v: &tspmodel2d.Vertex2D{X: 0.6, Y: 0.6}, expected: edges[2]},
		{v: &tspmodel2d.Vertex2D{X: 0.5, Y: 0.4}, expected: edges[4]},
		{v: &tspmodel2d.Vertex2D{X: 0.6, Y: 0.4}, expected: edges[3]},
		{v: &tspmodel2d.Vertex2D{X: 0.2, Y: 0.1}, expected: edges[4]},
	}

	for i, tc := range testCases {
		assert.Equal(tspmodel.FindClosestEdge(tc.v, edges), tc.expected, i)
	}
}

func TestFindClosestEdge_2D_ShouldReturnNilIfListIsEmpty(t *testing.T) {
	assert := assert.New(t)

	v := &tspmodel2d.Vertex2D{}

	assert.Nil(tspmodel.FindClosestEdge(v, []tspmodel.CircuitEdge{}))
}

func TestFindClosestEdge_3D(t *testing.T) {
	assert := assert.New(t)

	points := []*tspmodel3d.Vertex3D{
		{X: 0.0, Y: 0.0},
		{X: 0.0, Y: 1.0},
		{X: 1.0, Y: 1.0},
		{X: 0.7, Y: 0.5},
		{X: 1.0, Y: 0.0},
	}

	edges := []tspmodel.CircuitEdge{
		tspmodel3d.NewEdge3D(points[0], points[1]), //0 = 0.0,0.0 -> 0.0,1.0
		tspmodel3d.NewEdge3D(points[1], points[2]), //1 = 0.0,1.0 -> 1.0,1.0
		tspmodel3d.NewEdge3D(points[2], points[3]), //2 = 1.0,1.0 -> 0.7,0.5
		tspmodel3d.NewEdge3D(points[3], points[4]), //3 = 0.7,0.5 -> 1.0,0.0
		tspmodel3d.NewEdge3D(points[4], points[0]), //4 = 1.0,0.0 -> 0.0,0.0
	}

	testCases := []struct {
		v        *tspmodel3d.Vertex3D
		expected tspmodel.CircuitEdge
	}{
		{v: &tspmodel3d.Vertex3D{X: 0.0, Y: 0.0, Z: 0.0}, expected: edges[0]},
		{v: &tspmodel3d.Vertex3D{X: 0.5, Y: 0.0, Z: 0.0}, expected: edges[4]},
		{v: &tspmodel3d.Vertex3D{X: 0.5, Y: 0.5, Z: 0.0}, expected: edges[2]},
		{v: &tspmodel3d.Vertex3D{X: 0.5, Y: 0.6, Z: 0.0}, expected: edges[1]},
		{v: &tspmodel3d.Vertex3D{X: 0.6, Y: 0.6, Z: 0.0}, expected: edges[2]},
		{v: &tspmodel3d.Vertex3D{X: 0.5, Y: 0.4, Z: 0.0}, expected: edges[4]},
		{v: &tspmodel3d.Vertex3D{X: 0.6, Y: 0.4, Z: 0.0}, expected: edges[3]},
		{v: &tspmodel3d.Vertex3D{X: 0.2, Y: 0.1, Z: 0.0}, expected: edges[4]},
	}

	for i, tc := range testCases {
		assert.Equal(tspmodel.FindClosestEdge(tc.v, edges), tc.expected, i)
	}
}

func TestFindClosestEdgeList_2D(t *testing.T) {
	assert := assert.New(t)

	points := []*tspmodel2d.Vertex2D{
		{X: 0.0, Y: 0.0},
		{X: 0.0, Y: 1.0},
		{X: 1.0, Y: 1.0},
		{X: 0.7, Y: 0.5},
		{X: 1.0, Y: 0.0},
	}

	edges := list.New()
	edges.PushBack(tspmodel2d.NewEdge2D(points[0], points[1])) //0 = 0.0,0.0 -> 0.0,1.0
	edges.PushBack(tspmodel2d.NewEdge2D(points[1], points[2])) //1 = 0.0,1.0 -> 1.0,1.0
	edges.PushBack(tspmodel2d.NewEdge2D(points[2], points[3])) //2 = 1.0,1.0 -> 0.7,0.5
	edges.PushBack(tspmodel2d.NewEdge2D(points[3], points[4])) //3 = 0.7,0.5 -> 1.0,0.0
	edges.PushBack(tspmodel2d.NewEdge2D(points[4], points[0])) //4 = 1.0,0.0 -> 0.0,0.0

	testCases := []struct {
		v        *tspmodel2d.Vertex2D
		expected tspmodel.CircuitEdge
	}{
		{v: &tspmodel2d.Vertex2D{X: 0.0, Y: 0.0}, expected: tspmodel2d.NewEdge2D(points[0], points[1])},
		{v: &tspmodel2d.Vertex2D{X: 0.5, Y: 0.0}, expected: tspmodel2d.NewEdge2D(points[4], points[0])},
		{v: &tspmodel2d.Vertex2D{X: 0.5, Y: 0.5}, expected: tspmodel2d.NewEdge2D(points[2], points[3])},
		{v: &tspmodel2d.Vertex2D{X: 0.5, Y: 0.6}, expected: tspmodel2d.NewEdge2D(points[1], points[2])},
		{v: &tspmodel2d.Vertex2D{X: 0.6, Y: 0.6}, expected: tspmodel2d.NewEdge2D(points[2], points[3])},
		{v: &tspmodel2d.Vertex2D{X: 0.5, Y: 0.4}, expected: tspmodel2d.NewEdge2D(points[4], points[0])},
		{v: &tspmodel2d.Vertex2D{X: 0.6, Y: 0.4}, expected: tspmodel2d.NewEdge2D(points[3], points[4])},
		{v: &tspmodel2d.Vertex2D{X: 0.2, Y: 0.1}, expected: tspmodel2d.NewEdge2D(points[4], points[0])},
	}

	for i, tc := range testCases {
		assert.Equal(tspmodel.FindClosestEdgeList(tc.v, edges), tc.expected, i)
	}
}

func TestFindClosestEdge_3D_ShouldReturnNilIfListIsEmpty(t *testing.T) {
	assert := assert.New(t)

	v := &tspmodel3d.Vertex3D{}

	assert.Nil(tspmodel.FindClosestEdge(v, []tspmodel.CircuitEdge{}))
}

func TestFindFarthestPoint(t *testing.T) {
	assert := assert.New(t)

	vertices := []tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(-15, -15),
		tspmodel2d.NewVertex2D(0, 0),
		tspmodel2d.NewVertex2D(15, -15),
		tspmodel2d.NewVertex2D(3, 0),
		tspmodel2d.NewVertex2D(3, 13),
		tspmodel2d.NewVertex2D(8, 5),
		tspmodel2d.NewVertex2D(9, 6),
		tspmodel2d.NewVertex2D(-7, 6),
	}

	assert.Equal(vertices[4], tspmodel.FindFarthestPoint(vertices[0], vertices))
	assert.Equal(vertices[0], tspmodel.FindFarthestPoint(vertices[1], vertices))
	assert.Equal(vertices[4], tspmodel.FindFarthestPoint(vertices[2], vertices))
	assert.Equal(vertices[0], tspmodel.FindFarthestPoint(vertices[3], vertices))
	assert.Equal(vertices[0], tspmodel.FindFarthestPoint(vertices[4], vertices))
	assert.Equal(vertices[0], tspmodel.FindFarthestPoint(vertices[5], vertices))
	assert.Equal(vertices[0], tspmodel.FindFarthestPoint(vertices[6], vertices))
	assert.Equal(vertices[2], tspmodel.FindFarthestPoint(vertices[7], vertices))
}

func TestIndexOfVertex(t *testing.T) {
	assert := assert.New(t)

	initVertices := []tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(1, 1),
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(8, 8),
	}

	assert.Equal(-1, tspmodel.IndexOfVertex([]tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(8, 8),
	}, initVertices[0]))

	assert.Equal(3, tspmodel.IndexOfVertex([]tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(8, 8),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(4, 4),
	}, initVertices[2]))

	assert.Equal(6, tspmodel.IndexOfVertex([]tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(8, 8),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(4, 4),
	}, initVertices[3]))

	assert.Equal(0, tspmodel.IndexOfVertex([]tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(8, 8),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(4, 4),
	}, initVertices[1]))
}

func TestInsertVertex(t *testing.T) {
	assert := assert.New(t)
	circuit := []tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(-15, -15),
		tspmodel2d.NewVertex2D(0, 0),
		tspmodel2d.NewVertex2D(15, -15),
	}

	circuit = tspmodel.InsertVertex(circuit, 0, tspmodel2d.NewVertex2D(5, 5))
	assert.Len(circuit, 4)
	assert.Equal(tspmodel2d.NewVertex2D(5, 5), circuit[0])
	assert.Equal(tspmodel2d.NewVertex2D(-15, -15), circuit[1])
	assert.Equal(tspmodel2d.NewVertex2D(0, 0), circuit[2])
	assert.Equal(tspmodel2d.NewVertex2D(15, -15), circuit[3])

	circuit = tspmodel.InsertVertex(circuit, 4, tspmodel2d.NewVertex2D(-5, -5))
	assert.Len(circuit, 5)
	assert.Equal(tspmodel2d.NewVertex2D(5, 5), circuit[0])
	assert.Equal(tspmodel2d.NewVertex2D(-15, -15), circuit[1])
	assert.Equal(tspmodel2d.NewVertex2D(0, 0), circuit[2])
	assert.Equal(tspmodel2d.NewVertex2D(15, -15), circuit[3])
	assert.Equal(tspmodel2d.NewVertex2D(-5, -5), circuit[4])

	circuit = tspmodel.InsertVertex(circuit, 2, tspmodel2d.NewVertex2D(1, -5))
	assert.Len(circuit, 6)
	assert.Equal(tspmodel2d.NewVertex2D(5, 5), circuit[0])
	assert.Equal(tspmodel2d.NewVertex2D(-15, -15), circuit[1])
	assert.Equal(tspmodel2d.NewVertex2D(1, -5), circuit[2])
	assert.Equal(tspmodel2d.NewVertex2D(0, 0), circuit[3])
	assert.Equal(tspmodel2d.NewVertex2D(15, -15), circuit[4])
	assert.Equal(tspmodel2d.NewVertex2D(-5, -5), circuit[5])
}

func TestIsEdgeCloser_2D(t *testing.T) {
	assert := assert.New(t)

	v := tspmodel2d.NewVertex2D(10.0, 10.0)

	testCases := []struct {
		candiate *tspmodel2d.Edge2D
		current  *tspmodel2d.Edge2D
		expected bool
	}{
		{tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(0.0, 0.0), tspmodel2d.NewVertex2D(20.0, 0.0)), tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(0.0, 0.0), tspmodel2d.NewVertex2D(20.0, 20.0)), false},
		{tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(0.0, 0.0), tspmodel2d.NewVertex2D(20.0, 20.0)), tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(0.0, 0.0), tspmodel2d.NewVertex2D(20.0, 0.0)), true},
		{tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(0.0, 20.0), tspmodel2d.NewVertex2D(20.0, 0.0)), tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(0.0, 0.0), tspmodel2d.NewVertex2D(20.0, 20.0)), false},
		{tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(0.0, 0.0), tspmodel2d.NewVertex2D(20.0, 20.0)), tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(0.0, 20.0), tspmodel2d.NewVertex2D(20.0, 0.0)), false},
		{tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(0.0, 0.0), tspmodel2d.NewVertex2D(21.0, 0.0)), tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(0.0, 0.0), tspmodel2d.NewVertex2D(20.0, 0.0)), true},
		{tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(0.0, 0.0), tspmodel2d.NewVertex2D(18.0, 0.0)), tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(0.0, 0.0), tspmodel2d.NewVertex2D(20.0, 0.0)), false},
		{tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(2.0, 0.0), tspmodel2d.NewVertex2D(22.0, 0.0)), tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(4.0, 0.0), tspmodel2d.NewVertex2D(24.0, 0.0)), true},
		{tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(2.0, 0.0), tspmodel2d.NewVertex2D(22.0, 0.0)), tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(0.0, 0.0), tspmodel2d.NewVertex2D(20.0, 0.0)), false},
	}

	for i, tc := range testCases {
		assert.Equal(tc.expected, tspmodel.IsEdgeCloser(v, tc.candiate, tc.current), i)
	}
}

func TestIsEdgeCloser_3D(t *testing.T) {
	assert := assert.New(t)

	v := tspmodel3d.NewVertex3D(10.0, 10.0, 0.0)

	testCases := []struct {
		candiate *tspmodel3d.Edge3D
		current  *tspmodel3d.Edge3D
		expected bool
	}{
		{tspmodel3d.NewEdge3D(tspmodel3d.NewVertex3D(0.0, 0.0, 0.0), tspmodel3d.NewVertex3D(20.0, 0.0, 0.0)), tspmodel3d.NewEdge3D(tspmodel3d.NewVertex3D(0.0, 0.0, 0.0), tspmodel3d.NewVertex3D(20.0, 20.0, 0.0)), false},
		{tspmodel3d.NewEdge3D(tspmodel3d.NewVertex3D(0.0, 0.0, 0.0), tspmodel3d.NewVertex3D(20.0, 20.0, 0.0)), tspmodel3d.NewEdge3D(tspmodel3d.NewVertex3D(0.0, 0.0, 0.0), tspmodel3d.NewVertex3D(20.0, 0.0, 0.0)), true},
		{tspmodel3d.NewEdge3D(tspmodel3d.NewVertex3D(0.0, 20.0, 0.0), tspmodel3d.NewVertex3D(20.0, 0.0, 0.0)), tspmodel3d.NewEdge3D(tspmodel3d.NewVertex3D(0.0, 0.0, 0.0), tspmodel3d.NewVertex3D(20.0, 20.0, 0.0)), false},
		{tspmodel3d.NewEdge3D(tspmodel3d.NewVertex3D(0.0, 0.0, 0.0), tspmodel3d.NewVertex3D(20.0, 20.0, 0.0)), tspmodel3d.NewEdge3D(tspmodel3d.NewVertex3D(0.0, 20.0, 0.0), tspmodel3d.NewVertex3D(20.0, 0.0, 0.0)), false},
		{tspmodel3d.NewEdge3D(tspmodel3d.NewVertex3D(0.0, 0.0, 0.0), tspmodel3d.NewVertex3D(21.0, 0.0, 0.0)), tspmodel3d.NewEdge3D(tspmodel3d.NewVertex3D(0.0, 0.0, 0.0), tspmodel3d.NewVertex3D(20.0, 0.0, 0.0)), true},
		{tspmodel3d.NewEdge3D(tspmodel3d.NewVertex3D(0.0, 0.0, 0.0), tspmodel3d.NewVertex3D(18.0, 0.0, 0.0)), tspmodel3d.NewEdge3D(tspmodel3d.NewVertex3D(0.0, 0.0, 0.0), tspmodel3d.NewVertex3D(20.0, 0.0, 0.0)), false},
		{tspmodel3d.NewEdge3D(tspmodel3d.NewVertex3D(2.0, 0.0, 0.0), tspmodel3d.NewVertex3D(22.0, 0.0, 0.0)), tspmodel3d.NewEdge3D(tspmodel3d.NewVertex3D(4.0, 0.0, 0.0), tspmodel3d.NewVertex3D(24.0, 0.0, 0.0)), true},
		{tspmodel3d.NewEdge3D(tspmodel3d.NewVertex3D(2.0, 0.0, 0.0), tspmodel3d.NewVertex3D(22.0, 0.0, 0.0)), tspmodel3d.NewEdge3D(tspmodel3d.NewVertex3D(0.0, 0.0, 0.0), tspmodel3d.NewVertex3D(20.0, 0.0, 0.0)), false},
	}

	for i, tc := range testCases {
		assert.Equal(tc.expected, tspmodel.IsEdgeCloser(v, tc.candiate, tc.current), i)
	}
}
