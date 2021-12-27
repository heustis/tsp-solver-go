package model_test

import (
	"container/list"
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

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

	edges, a, b, merged = model.MergeEdges2(edges, model2d.NewVertex2D(10, 10))
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[0]),
	}, edges)
	assert.Nil(a)
	assert.Nil(b)
	assert.Nil(merged)

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
	assert.Len(edges, 1)
	edges, _, _, _ = model.MergeEdges2(edges, vertices[2])
	assert.Len(edges, 1)
}

func TestMergeEdgesList(t *testing.T) {
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

	edges := list.New()
	edges.PushBack(model2d.NewEdge2D(vertices[0], vertices[1]))
	edges.PushBack(model2d.NewEdge2D(vertices[1], vertices[2]))
	edges.PushBack(model2d.NewEdge2D(vertices[2], vertices[3]))
	edges.PushBack(model2d.NewEdge2D(vertices[3], vertices[4]))
	edges.PushBack(model2d.NewEdge2D(vertices[4], vertices[5]))
	edges.PushBack(model2d.NewEdge2D(vertices[5], vertices[6]))
	edges.PushBack(model2d.NewEdge2D(vertices[6], vertices[7]))
	edges.PushBack(model2d.NewEdge2D(vertices[7], vertices[0]))

	a, b, merged := model.MergeEdgesList(edges, vertices[1])
	assert.Equal(7, edges.Len())
	expected := []model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[0]),
	}
	for i, testNode := 0, edges.Front(); testNode != nil; i, testNode = i+1, testNode.Next() {
		assert.Equal(expected[i], testNode.Value)
	}
	assert.Equal(model2d.NewEdge2D(vertices[0], vertices[1]), a)
	assert.Equal(model2d.NewEdge2D(vertices[1], vertices[2]), b)
	assert.Equal(model2d.NewEdge2D(vertices[0], vertices[2]), merged.Value)

	a, b, merged = model.MergeEdgesList(edges, model2d.NewVertex2D(10, 10))
	assert.Equal(7, edges.Len())
	expected = []model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[0]),
	}
	for i, testNode := 0, edges.Front(); testNode != nil; i, testNode = i+1, testNode.Next() {
		assert.Equal(expected[i], testNode.Value)
	}
	assert.Nil(a)
	assert.Nil(b)
	assert.Nil(merged)

	a, b, merged = model.MergeEdgesList(edges, vertices[0])
	assert.Equal(6, edges.Len())
	expected = []model.CircuitEdge{
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[2]),
	}
	for i, testNode := 0, edges.Front(); testNode != nil; i, testNode = i+1, testNode.Next() {
		assert.Equal(expected[i], testNode.Value, i)
	}
	assert.Equal(model2d.NewEdge2D(vertices[7], vertices[0]), a)
	assert.Equal(model2d.NewEdge2D(vertices[0], vertices[2]), b)
	assert.Equal(model2d.NewEdge2D(vertices[7], vertices[2]), merged.Value)

	a, b, merged = model.MergeEdgesList(edges, vertices[7])
	assert.Equal(5, edges.Len())
	expected = []model.CircuitEdge{
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[2]),
	}
	for i, testNode := 0, edges.Front(); testNode != nil; i, testNode = i+1, testNode.Next() {
		assert.Equal(expected[i], testNode.Value)
	}
	assert.Equal(model2d.NewEdge2D(vertices[6], vertices[7]), a)
	assert.Equal(model2d.NewEdge2D(vertices[7], vertices[2]), b)
	assert.Equal(model2d.NewEdge2D(vertices[6], vertices[2]), merged.Value)

	a, b, merged = model.MergeEdgesList(edges, vertices[3])
	assert.Equal(4, edges.Len())
	expected = []model.CircuitEdge{
		model2d.NewEdge2D(vertices[2], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[2]),
	}
	for i, testNode := 0, edges.Front(); testNode != nil; i, testNode = i+1, testNode.Next() {
		assert.Equal(expected[i], testNode.Value)
	}
	assert.Equal(model2d.NewEdge2D(vertices[2], vertices[3]), a)
	assert.Equal(model2d.NewEdge2D(vertices[3], vertices[4]), b)
	assert.Equal(model2d.NewEdge2D(vertices[2], vertices[4]), merged.Value)

	_, _, _ = model.MergeEdgesList(edges, vertices[4])
	_, _, _ = model.MergeEdgesList(edges, vertices[5])
	assert.Equal(2, edges.Len())
	_, _, _ = model.MergeEdgesList(edges, vertices[6])
	assert.Equal(1, edges.Len())
	_, _, _ = model.MergeEdgesList(edges, vertices[2])
	assert.Equal(1, edges.Len())
	_, _, _ = model.MergeEdgesList(edges, vertices[2])
	assert.Equal(1, edges.Len())
}

func TestMoveVertex(t *testing.T) {
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

	edges, a, b, c := model.MoveVertex(edges, vertices[0], model2d.NewEdge2D(vertices[3], vertices[4]))
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[0]),
		model2d.NewEdge2D(vertices[0], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[1]),
	}, edges)
	assert.Equal(model2d.NewEdge2D(vertices[7], vertices[1]), a)
	assert.Equal(model2d.NewEdge2D(vertices[3], vertices[0]), b)
	assert.Equal(model2d.NewEdge2D(vertices[0], vertices[4]), c)

	edges, a, b, c = model.MoveVertex(edges, vertices[7], model2d.NewEdge2D(vertices[1], vertices[2]))
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[1], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[0]),
		model2d.NewEdge2D(vertices[0], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[1]),
	}, edges)
	assert.Equal(model2d.NewEdge2D(vertices[6], vertices[1]), a)
	assert.Equal(model2d.NewEdge2D(vertices[1], vertices[7]), b)
	assert.Equal(model2d.NewEdge2D(vertices[7], vertices[2]), c)

	edges, a, b, c = model.MoveVertex(edges, vertices[4], model2d.NewEdge2D(vertices[6], vertices[1]))
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[1], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[0]),
		model2d.NewEdge2D(vertices[0], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[1]),
	}, edges)
	assert.Equal(model2d.NewEdge2D(vertices[0], vertices[5]), a)
	assert.Equal(model2d.NewEdge2D(vertices[6], vertices[4]), b)
	assert.Equal(model2d.NewEdge2D(vertices[4], vertices[1]), c)

	edges, a, b, c = model.MoveVertex(edges, model2d.NewVertex2D(9, 9), model2d.NewEdge2D(vertices[6], vertices[1]))
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[1], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[0]),
		model2d.NewEdge2D(vertices[0], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[1]),
	}, edges)
	assert.Nil(a)
	assert.Nil(b)
	assert.Nil(c)

	edges, a, b, c = model.MoveVertex(edges, vertices[7], model2d.NewEdge2D(vertices[6], vertices[1]))
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[1], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[0]),
		model2d.NewEdge2D(vertices[0], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[1]),
	}, edges)
	assert.Nil(a)
	assert.Nil(b)
	assert.Nil(c)

	edges2 := []model.CircuitEdge{
		model2d.NewEdge2D(vertices[1], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[1]),
	}
	edges2, a, b, c = model.MoveVertex(edges2, vertices[7], model2d.NewEdge2D(vertices[7], vertices[1]))
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[1], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[1]),
	}, edges2)
	assert.Nil(a)
	assert.Nil(b)
	assert.Nil(c)

	edges2 = []model.CircuitEdge{
		model2d.NewEdge2D(vertices[1], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[1]),
	}
	edges2, a, b, c = model.MoveVertex(edges2, vertices[7], model2d.NewEdge2D(vertices[2], vertices[1]))
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[1]),
	}, edges2)
	assert.Equal(model2d.NewEdge2D(vertices[1], vertices[2]), a)
	assert.Equal(model2d.NewEdge2D(vertices[2], vertices[7]), b)
	assert.Equal(model2d.NewEdge2D(vertices[7], vertices[1]), c)

	edges2, a, b, c = model.MoveVertex(edges2, vertices[1], model2d.NewEdge2D(vertices[2], vertices[7]))
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[2], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[2]),
	}, edges2)
	assert.Equal(model2d.NewEdge2D(vertices[7], vertices[2]), a)
	assert.Equal(model2d.NewEdge2D(vertices[2], vertices[1]), b)
	assert.Equal(model2d.NewEdge2D(vertices[1], vertices[7]), c)
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

func TestSplitEdgeList(t *testing.T) {
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

	edges := list.New()
	edges.PushBack(model2d.NewEdge2D(vertices[0], vertices[1]))
	edges.PushBack(model2d.NewEdge2D(vertices[1], vertices[2]))
	edges.PushBack(model2d.NewEdge2D(vertices[2], vertices[3]))
	edges.PushBack(model2d.NewEdge2D(vertices[3], vertices[4]))
	edges.PushBack(model2d.NewEdge2D(vertices[4], vertices[5]))

	newEdge := model.SplitEdgeList(edges, model2d.NewEdge2D(vertices[0], vertices[7]), vertices[6])
	assert.Nil(newEdge)
	expected := []model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
	}
	for i, testNode := 0, edges.Front(); testNode != nil; i, testNode = i+1, testNode.Next() {
		assert.Equal(expected[i], testNode.Value, i)
	}

	newEdge = model.SplitEdgeList(edges, model2d.NewEdge2D(vertices[0], vertices[1]), vertices[6])
	assert.NotNil(newEdge)
	assert.NotNil(newEdge.Prev())
	assert.Equal(model2d.NewEdge2D(vertices[0], vertices[6]), newEdge.Prev().Value)
	assert.Equal(model2d.NewEdge2D(vertices[6], vertices[1]), newEdge.Value)
	expected = []model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
	}
	for i, testNode := 0, edges.Front(); testNode != nil; i, testNode = i+1, testNode.Next() {
		assert.Equal(expected[i], testNode.Value, i)
	}

	newEdge = model.SplitEdgeList(edges, model2d.NewEdge2D(vertices[4], vertices[5]), vertices[7])
	assert.NotNil(newEdge)
	assert.NotNil(newEdge.Prev())
	assert.Equal(model2d.NewEdge2D(vertices[4], vertices[7]), newEdge.Prev().Value)
	assert.Equal(model2d.NewEdge2D(vertices[7], vertices[5]), newEdge.Value)
	expected = []model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[7]),
		model2d.NewEdge2D(vertices[7], vertices[5]),
	}
	for i, testNode := 0, edges.Front(); testNode != nil; i, testNode = i+1, testNode.Next() {
		assert.Equal(expected[i], testNode.Value, i)
	}
}
