package model_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
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

func TestIndexOfEdge(t *testing.T) {
	assert := assert.New(t)

	vertices := []*model2d.Vertex2D{
		model2d.NewVertex2D(1, 1),
		model2d.NewVertex2D(2, 2),
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
		model2d.NewVertex2D(8, 8),
	}

	edges := []model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[0]),
	}

	assert.Equal(-1, model.IndexOfEdge(edges, model2d.NewEdge2D(vertices[0], vertices[2])))
	assert.Equal(0, model.IndexOfEdge(edges, model2d.NewEdge2D(vertices[0], vertices[1])))
	assert.Equal(7, model.IndexOfEdge(edges, model2d.NewEdge2D(vertices[7], vertices[0])))
	assert.Equal(7, model.IndexOfEdge(edges, edges[7]))
	assert.Equal(-1, model.IndexOfEdge(edges, model2d.NewEdge2D(vertices[0], vertices[7])))
	assert.Equal(7, model.IndexOfEdge(edges, model2d.NewEdge2D(model2d.NewVertex2D(8, 8), model2d.NewVertex2D(1, 1))))
	assert.Equal(-1, model.IndexOfEdge(edges, model2d.NewEdge2D(model2d.NewVertex2D(7.998, 8), model2d.NewVertex2D(1, 1))))
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

func TestMergeEdges(t *testing.T) {
	assert := assert.New(t)

	vertices := []*model2d.Vertex2D{
		model2d.NewVertex2D(1, 1),
		model2d.NewVertex2D(2, 2),
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
		model2d.NewVertex2D(8, 8),
	}

	edges := []model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[0]),
	}

	edges, a, b := model.MergeEdges(edges, 1)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[0]),
	}, edges)
	assert.Equal(model2d.NewEdge2D(vertices[0], vertices[1]), a)
	assert.Equal(model2d.NewEdge2D(vertices[1], vertices[2]), b)

	edges, a, b = model.MergeEdges(edges, 0)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[2]),
	}, edges)
	assert.Equal(model2d.NewEdge2D(vertices[7], vertices[0]), a)
	assert.Equal(model2d.NewEdge2D(vertices[0], vertices[2]), b)

	edges, a, b = model.MergeEdges(edges, 15)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[2]),
	}, edges)
	assert.Equal(model2d.NewEdge2D(vertices[6], vertices[7]), a)
	assert.Equal(model2d.NewEdge2D(vertices[7], vertices[2]), b)

	edges, a, b = model.MergeEdges(edges, 1)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[2], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[2]),
	}, edges)
	assert.Equal(model2d.NewEdge2D(vertices[2], vertices[3]), a)
	assert.Equal(model2d.NewEdge2D(vertices[3], vertices[4]), b)

	edges, _, _ = model.MergeEdges(edges, 0)
	edges, _, _ = model.MergeEdges(edges, 0)
	assert.Len(edges, 2)
	edges, _, _ = model.MergeEdges(edges, 0)
	assert.Len(edges, 1)
	edges, _, _ = model.MergeEdges(edges, 0)
	assert.Len(edges, 0)
	edges, _, _ = model.MergeEdges(edges, 0)
	assert.Len(edges, 0)
}

func TestMergeEdges2(t *testing.T) {
	assert := assert.New(t)

	vertices := []*model2d.Vertex2D{
		model2d.NewVertex2D(1, 1),
		model2d.NewVertex2D(2, 2),
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
		model2d.NewVertex2D(8, 8),
	}

	initEdges := []model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[0]),
	}

	edges, a, b, merged := model.MergeEdges2(initEdges, vertices[1])
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[0]),
	}, edges)
	assert.Equal(model2d.NewEdge2D(vertices[0], vertices[1]), a)
	assert.Equal(model2d.NewEdge2D(vertices[1], vertices[2]), b)
	assert.Equal(model2d.NewEdge2D(vertices[0], vertices[2]), merged)

	edges, a, b, merged = model.MergeEdges2(edges, vertices[0])
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[2]),
	}, edges)
	assert.Equal(model2d.NewEdge2D(vertices[7], vertices[0]), a)
	assert.Equal(model2d.NewEdge2D(vertices[0], vertices[2]), b)
	assert.Equal(model2d.NewEdge2D(vertices[7], vertices[2]), merged)

	edges, a, b, merged = model.MergeEdges2(edges, vertices[7])
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[2]),
	}, edges)
	assert.Equal(model2d.NewEdge2D(vertices[6], vertices[7]), a)
	assert.Equal(model2d.NewEdge2D(vertices[7], vertices[2]), b)
	assert.Equal(model2d.NewEdge2D(vertices[6], vertices[2]), merged)

	edges, a, b, merged = model.MergeEdges2(edges, vertices[3])
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[2], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[2]),
	}, edges)
	assert.Equal(model2d.NewEdge2D(vertices[2], vertices[3]), a)
	assert.Equal(model2d.NewEdge2D(vertices[3], vertices[4]), b)
	assert.Equal(model2d.NewEdge2D(vertices[2], vertices[4]), merged)

	edges, _, _, _ = model.MergeEdges2(edges, vertices[4])
	edges, _, _, _ = model.MergeEdges2(edges, vertices[5])
	assert.Len(edges, 2)
	edges, _, _, _ = model.MergeEdges2(edges, vertices[6])
	assert.Len(edges, 1)
	edges, _, _, _ = model.MergeEdges2(edges, vertices[2])
	assert.Len(edges, 0)
	edges, _, _, _ = model.MergeEdges2(edges, vertices[2])
	assert.Len(edges, 0)
}

func TestSplitEdge(t *testing.T) {
	assert := assert.New(t)

	vertices := []*model2d.Vertex2D{
		model2d.NewVertex2D(1, 1),
		model2d.NewVertex2D(2, 2),
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
		model2d.NewVertex2D(8, 8),
	}

	edges := []model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
	}

	var index int
	edges, index = model.SplitEdge(edges, model2d.NewEdge2D(vertices[0], vertices[7]), vertices[6])
	assert.Equal(-1, index)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
	}, edges)

	edges, index = model.SplitEdge(edges, edges[0], vertices[6])
	assert.Equal(0, index)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
	}, edges)

	edges, index = model.SplitEdge(edges, edges[5], vertices[7])
	assert.Equal(5, index)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[5]),
	}, edges)
}

func TestSplitEdge2(t *testing.T) {
	assert := assert.New(t)

	vertices := []*model2d.Vertex2D{
		model2d.NewVertex2D(1, 1),
		model2d.NewVertex2D(2, 2),
		model2d.NewVertex2D(3, 3),
		model2d.NewVertex2D(4, 4),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6, 6),
		model2d.NewVertex2D(7, 7),
		model2d.NewVertex2D(8, 8),
	}

	edges := []model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
	}

	var index int
	edges, index = model.SplitEdge2(edges, model2d.NewEdge2D(vertices[0], vertices[7]), vertices[6])
	assert.Equal(-1, index)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
	}, edges)

	edges, index = model.SplitEdge2(edges, edges[0], vertices[6])
	assert.Equal(0, index)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
	}, edges)

	edges, index = model.SplitEdge2(edges, edges[5], vertices[7])
	assert.Equal(5, index)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[5]),
	}, edges)
}
