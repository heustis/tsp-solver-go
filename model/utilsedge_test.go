package model_test

import (
	"container/list"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func BenchmarkMergeEdges(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	numVertices := int(math.Max(10000000, float64(b.N)*100.0))
	vertices := model2d.GenerateVertices(numVertices)
	edges := make([]model.CircuitEdge, numVertices-1)
	edges2 := make([]model.CircuitEdge, numVertices-1)
	edges3 := make([]model.CircuitEdge, numVertices-1)
	edgesList := list.New()

	for i := 1; i < numVertices; i++ {
		edge := vertices[i-1].EdgeTo(vertices[i])
		edges[i-1] = edge
		edges2[i-1] = edge
		edges3[i-1] = edge
		edgesList.PushBack(edge)
	}

	// BenchmarkMergeEdgesByIndex/MergeEdgesByIndex-16         	1000000000	         0.007999 ns/op	       0 B/op	       0 allocs/op
	// BenchmarkMergeEdgesByIndex/MergeEdgesByVertex-16        	1000000000	         0.01300 ns/op	       0 B/op	       0 allocs/op
	// BenchmarkMergeEdgesByIndex/MergeEdgesCopy-16            	1000000000	         0.05601 ns/op	       0 B/op	       0 allocs/op
	// BenchmarkMergeEdgesByIndex/MergeEdgesList-16            	1000000000	         0.05401 ns/op	       0 B/op	       0 allocs/op

	b.Run("MergeEdgesByIndex", func(b *testing.B) {
		model.MergeEdgesByIndex(edges, r.Intn(len(edges)))
	})

	b.Run("MergeEdgesByVertex", func(b *testing.B) {
		model.MergeEdgesByVertex(edges2, edges2[r.Intn(len(edges2))].GetStart())
	})

	b.Run("MergeEdgesCopy", func(b *testing.B) {
		model.MergeEdgesCopy(edges3, edges3[r.Intn(len(edges3))].GetStart())
	})

	b.Run("MergeEdgesList", func(b *testing.B) {
		model.MergeEdgesList(edgesList, vertices[r.Intn(len(vertices))])
	})
}

func BenchmarkSplitEdges(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	numVertices := int(math.Max(10000000, float64(b.N)*1000.0))
	vertices := model2d.GenerateVertices(numVertices)
	edges := make([]model.CircuitEdge, numVertices-1)
	edges2 := make([]model.CircuitEdge, numVertices-1)
	edgesList := list.New()

	for i := 1; i < numVertices; i++ {
		edge := vertices[i-1].EdgeTo(vertices[i])
		edges[i-1] = edge
		edges2[i-1] = edge
		edgesList.PushBack(edge)
	}

	// BenchmarkSplitEdges/SplitEdge-16         	1000000000	         0.06982 ns/op	       0 B/op	       0 allocs/op
	// BenchmarkSplitEdges/SplitEdgeCopy-16     	1000000000	         0.08602 ns/op	       0 B/op	       0 allocs/op
	// BenchmarkSplitEdges/MergeEdgesList-16    	1000000000	         0.1195 ns/op	       0 B/op	       0 allocs/op

	b.Run("SplitEdge", func(b *testing.B) {
		model.SplitEdge(edges, edges[r.Intn(len(edges))], vertices[r.Intn(len(vertices))])
	})

	b.Run("SplitEdgeCopy", func(b *testing.B) {
		model.SplitEdgeCopy(edges2, edges2[r.Intn(len(edges))], vertices[r.Intn(len(vertices))])
	})

	b.Run("SplitEdgeList", func(b *testing.B) {
		edgeIndex := r.Intn(len(vertices) - 1)
		model.SplitEdgeList(edgesList, vertices[edgeIndex].EdgeTo(vertices[edgeIndex+1]), vertices[r.Intn(len(vertices))])
	})
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

func TestMergeEdgesByIndex(t *testing.T) {
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

	edges, a, b := model.MergeEdgesByIndex(edges, 1)
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

	edges, a, b = model.MergeEdgesByIndex(edges, 0)
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

	edges, a, b = model.MergeEdgesByIndex(edges, 15)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[2]),
	}, edges)
	assert.Equal(model2d.NewEdge2D(vertices[6], vertices[7]), a)
	assert.Equal(model2d.NewEdge2D(vertices[7], vertices[2]), b)

	edges, a, b = model.MergeEdgesByIndex(edges, 1)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[2], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[2]),
	}, edges)
	assert.Equal(model2d.NewEdge2D(vertices[2], vertices[3]), a)
	assert.Equal(model2d.NewEdge2D(vertices[3], vertices[4]), b)

	edges, _, _ = model.MergeEdgesByIndex(edges, 0)
	edges, _, _ = model.MergeEdgesByIndex(edges, 0)
	assert.Len(edges, 2)
	edges, _, _ = model.MergeEdgesByIndex(edges, 0)
	assert.Len(edges, 1)
	edges, _, _ = model.MergeEdgesByIndex(edges, 0)
	assert.Len(edges, 0)
	edges, _, _ = model.MergeEdgesByIndex(edges, 0)
	assert.Len(edges, 0)
}

func TestMergeEdgesByVertex(t *testing.T) {
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

	edges, a, b, merged := model.MergeEdgesByVertex(initEdges, vertices[1])
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

	edges, a, b, merged = model.MergeEdgesByVertex(edges, model2d.NewVertex2D(10, 10))
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

	edges, a, b, merged = model.MergeEdgesByVertex(edges, vertices[0])
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

	edges, a, b, merged = model.MergeEdgesByVertex(edges, vertices[7])
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

	edges, a, b, merged = model.MergeEdgesByVertex(edges, vertices[3])
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[2], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[2]),
	}, edges)
	assert.Equal(model2d.NewEdge2D(vertices[2], vertices[3]), a)
	assert.Equal(model2d.NewEdge2D(vertices[3], vertices[4]), b)
	assert.Equal(model2d.NewEdge2D(vertices[2], vertices[4]), merged)

	edges, _, _, _ = model.MergeEdgesByVertex(edges, vertices[4])
	edges, _, _, _ = model.MergeEdgesByVertex(edges, vertices[5])
	assert.Len(edges, 2)
	edges, _, _, _ = model.MergeEdgesByVertex(edges, vertices[6])
	assert.Len(edges, 1)
	edges, _, _, _ = model.MergeEdgesByVertex(edges, vertices[2])
	assert.Len(edges, 1)
	edges, _, _, _ = model.MergeEdgesByVertex(edges, vertices[2])
	assert.Len(edges, 1)
}

func TestMergeEdgesCopy(t *testing.T) {
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

	edges, a, b, merged := model.MergeEdgesCopy(initEdges, vertices[1])
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

	edges, a, b, merged = model.MergeEdgesCopy(edges, model2d.NewVertex2D(10, 10))
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

	edges, a, b, merged = model.MergeEdgesCopy(edges, vertices[0])
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

	edges, a, b, merged = model.MergeEdgesCopy(edges, vertices[7])
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

	edges, a, b, merged = model.MergeEdgesCopy(edges, vertices[3])
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[2], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
		model2d.NewEdge2D(vertices[5], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[2]),
	}, edges)
	assert.Equal(model2d.NewEdge2D(vertices[2], vertices[3]), a)
	assert.Equal(model2d.NewEdge2D(vertices[3], vertices[4]), b)
	assert.Equal(model2d.NewEdge2D(vertices[2], vertices[4]), merged)

	edges, _, _, _ = model.MergeEdgesCopy(edges, vertices[4])
	edges, _, _, _ = model.MergeEdgesCopy(edges, vertices[5])
	assert.Len(edges, 2)
	edges, _, _, _ = model.MergeEdgesCopy(edges, vertices[6])
	assert.Len(edges, 1)
	edges, _, _, _ = model.MergeEdgesCopy(edges, vertices[2])
	assert.Len(edges, 1)
	edges, _, _, _ = model.MergeEdgesCopy(edges, vertices[2])
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

func TestSplitEdgeCopy(t *testing.T) {
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
	edges, index = model.SplitEdgeCopy(edges, model2d.NewEdge2D(vertices[0], vertices[7]), vertices[6])
	assert.Equal(-1, index)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
	}, edges)

	edges, index = model.SplitEdgeCopy(edges, edges[0], vertices[6])
	assert.Equal(0, index)
	assert.Equal([]model.CircuitEdge{
		model2d.NewEdge2D(vertices[0], vertices[6]),
		model2d.NewEdge2D(vertices[6], vertices[1]),
		model2d.NewEdge2D(vertices[1], vertices[2]),
		model2d.NewEdge2D(vertices[2], vertices[3]),
		model2d.NewEdge2D(vertices[3], vertices[4]),
		model2d.NewEdge2D(vertices[4], vertices[5]),
	}, edges)

	edges, index = model.SplitEdgeCopy(edges, edges[5], vertices[7])
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
