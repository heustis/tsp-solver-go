package tspmodel_test

import (
	"container/list"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/fealos/lee-tsp-go/tspmodel2d"
	"github.com/stretchr/testify/assert"
)

func BenchmarkMergeEdges(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))

	numVertices := int(math.Max(10000000, float64(b.N)*100.0))
	vertices := tspmodel2d.GenerateVertices(numVertices)
	edges := make([]tspmodel.CircuitEdge, numVertices-1)
	edges2 := make([]tspmodel.CircuitEdge, numVertices-1)
	edges3 := make([]tspmodel.CircuitEdge, numVertices-1)
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
		tspmodel.MergeEdgesByIndex(edges, r.Intn(len(edges)))
	})

	b.Run("MergeEdgesByVertex", func(b *testing.B) {
		tspmodel.MergeEdgesByVertex(edges2, edges2[r.Intn(len(edges2))].GetStart())
	})

	b.Run("MergeEdgesCopy", func(b *testing.B) {
		tspmodel.MergeEdgesCopy(edges3, edges3[r.Intn(len(edges3))].GetStart())
	})

	b.Run("MergeEdgesList", func(b *testing.B) {
		tspmodel.MergeEdgesList(edgesList, vertices[r.Intn(len(vertices))])
	})
}

func BenchmarkSplitEdges(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))

	numVertices := int(math.Max(10000000, float64(b.N)*1000.0))
	vertices := tspmodel2d.GenerateVertices(numVertices)
	edges := make([]tspmodel.CircuitEdge, numVertices-1)
	edges2 := make([]tspmodel.CircuitEdge, numVertices-1)
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
		tspmodel.SplitEdge(edges, edges[r.Intn(len(edges))], vertices[r.Intn(len(vertices))])
	})

	b.Run("SplitEdgeCopy", func(b *testing.B) {
		tspmodel.SplitEdgeCopy(edges2, edges2[r.Intn(len(edges))], vertices[r.Intn(len(vertices))])
	})

	b.Run("SplitEdgeList", func(b *testing.B) {
		edgeIndex := r.Intn(len(vertices) - 1)
		tspmodel.SplitEdgeList(edgesList, vertices[edgeIndex].EdgeTo(vertices[edgeIndex+1]), vertices[r.Intn(len(vertices))])
	})
}

func TestIndexOfEdge(t *testing.T) {
	assert := assert.New(t)

	vertices := []*tspmodel2d.Vertex2D{
		tspmodel2d.NewVertex2D(1, 1),
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(8, 8),
	}

	edges := []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[0]),
	}

	assert.Equal(-1, tspmodel.IndexOfEdge(edges, tspmodel2d.NewEdge2D(vertices[0], vertices[2])))
	assert.Equal(0, tspmodel.IndexOfEdge(edges, tspmodel2d.NewEdge2D(vertices[0], vertices[1])))
	assert.Equal(7, tspmodel.IndexOfEdge(edges, tspmodel2d.NewEdge2D(vertices[7], vertices[0])))
	assert.Equal(7, tspmodel.IndexOfEdge(edges, edges[7]))
	assert.Equal(-1, tspmodel.IndexOfEdge(edges, tspmodel2d.NewEdge2D(vertices[0], vertices[7])))
	assert.Equal(7, tspmodel.IndexOfEdge(edges, tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(8, 8), tspmodel2d.NewVertex2D(1, 1))))
	assert.Equal(-1, tspmodel.IndexOfEdge(edges, tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(7.998, 8), tspmodel2d.NewVertex2D(1, 1))))
}

func TestMergeEdgesByIndex(t *testing.T) {
	assert := assert.New(t)

	vertices := []*tspmodel2d.Vertex2D{
		tspmodel2d.NewVertex2D(1, 1),
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(8, 8),
	}

	edges := []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[0]),
	}

	edges, a, b := tspmodel.MergeEdgesByIndex(edges, 1)
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[0]),
	}, edges)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[0], vertices[1]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[1], vertices[2]), b)

	edges, a, b = tspmodel.MergeEdgesByIndex(edges, 0)
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[2]),
	}, edges)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[7], vertices[0]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[0], vertices[2]), b)

	edges, a, b = tspmodel.MergeEdgesByIndex(edges, 15)
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[2]),
	}, edges)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[6], vertices[7]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[7], vertices[2]), b)

	edges, a, b = tspmodel.MergeEdgesByIndex(edges, 1)
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[2], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[2]),
	}, edges)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[2], vertices[3]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[3], vertices[4]), b)

	edges, _, _ = tspmodel.MergeEdgesByIndex(edges, 0)
	edges, _, _ = tspmodel.MergeEdgesByIndex(edges, 0)
	assert.Len(edges, 2)
	edges, _, _ = tspmodel.MergeEdgesByIndex(edges, 0)
	assert.Len(edges, 1)
	edges, _, _ = tspmodel.MergeEdgesByIndex(edges, 0)
	assert.Len(edges, 0)
	edges, _, _ = tspmodel.MergeEdgesByIndex(edges, 0)
	assert.Len(edges, 0)
}

func TestMergeEdgesByVertex(t *testing.T) {
	assert := assert.New(t)

	vertices := []*tspmodel2d.Vertex2D{
		tspmodel2d.NewVertex2D(1, 1),
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(8, 8),
	}

	initEdges := []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[0]),
	}

	edges, a, b, merged := tspmodel.MergeEdgesByVertex(initEdges, vertices[1])
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[0]),
	}, edges)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[0], vertices[1]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[1], vertices[2]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[0], vertices[2]), merged)

	edges, a, b, merged = tspmodel.MergeEdgesByVertex(edges, tspmodel2d.NewVertex2D(10, 10))
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[0]),
	}, edges)
	assert.Nil(a)
	assert.Nil(b)
	assert.Nil(merged)

	edges, a, b, merged = tspmodel.MergeEdgesByVertex(edges, vertices[0])
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[2]),
	}, edges)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[7], vertices[0]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[0], vertices[2]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[7], vertices[2]), merged)

	edges, a, b, merged = tspmodel.MergeEdgesByVertex(edges, vertices[7])
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[2]),
	}, edges)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[6], vertices[7]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[7], vertices[2]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[6], vertices[2]), merged)

	edges, a, b, merged = tspmodel.MergeEdgesByVertex(edges, vertices[3])
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[2], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[2]),
	}, edges)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[2], vertices[3]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[3], vertices[4]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[2], vertices[4]), merged)

	edges, _, _, _ = tspmodel.MergeEdgesByVertex(edges, vertices[4])
	edges, _, _, _ = tspmodel.MergeEdgesByVertex(edges, vertices[5])
	assert.Len(edges, 2)
	edges, _, _, _ = tspmodel.MergeEdgesByVertex(edges, vertices[6])
	assert.Len(edges, 1)
	edges, _, _, _ = tspmodel.MergeEdgesByVertex(edges, vertices[2])
	assert.Len(edges, 1)
	edges, _, _, _ = tspmodel.MergeEdgesByVertex(edges, vertices[2])
	assert.Len(edges, 1)
}

func TestMergeEdgesCopy(t *testing.T) {
	assert := assert.New(t)

	vertices := []*tspmodel2d.Vertex2D{
		tspmodel2d.NewVertex2D(1, 1),
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(8, 8),
	}

	initEdges := []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[0]),
	}

	edges, a, b, merged := tspmodel.MergeEdgesCopy(initEdges, vertices[1])
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[0]),
	}, edges)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[0], vertices[1]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[1], vertices[2]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[0], vertices[2]), merged)

	edges, a, b, merged = tspmodel.MergeEdgesCopy(edges, tspmodel2d.NewVertex2D(10, 10))
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[0]),
	}, edges)
	assert.Nil(a)
	assert.Nil(b)
	assert.Nil(merged)

	edges, a, b, merged = tspmodel.MergeEdgesCopy(edges, vertices[0])
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[2]),
	}, edges)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[7], vertices[0]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[0], vertices[2]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[7], vertices[2]), merged)

	edges, a, b, merged = tspmodel.MergeEdgesCopy(edges, vertices[7])
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[2]),
	}, edges)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[6], vertices[7]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[7], vertices[2]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[6], vertices[2]), merged)

	edges, a, b, merged = tspmodel.MergeEdgesCopy(edges, vertices[3])
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[2], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[2]),
	}, edges)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[2], vertices[3]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[3], vertices[4]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[2], vertices[4]), merged)

	edges, _, _, _ = tspmodel.MergeEdgesCopy(edges, vertices[4])
	edges, _, _, _ = tspmodel.MergeEdgesCopy(edges, vertices[5])
	assert.Len(edges, 2)
	edges, _, _, _ = tspmodel.MergeEdgesCopy(edges, vertices[6])
	assert.Len(edges, 1)
	edges, _, _, _ = tspmodel.MergeEdgesCopy(edges, vertices[2])
	assert.Len(edges, 1)
	edges, _, _, _ = tspmodel.MergeEdgesCopy(edges, vertices[2])
	assert.Len(edges, 1)
}

func TestMergeEdgesList(t *testing.T) {
	assert := assert.New(t)

	vertices := []*tspmodel2d.Vertex2D{
		tspmodel2d.NewVertex2D(1, 1),
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(8, 8),
	}

	edges := list.New()
	edges.PushBack(tspmodel2d.NewEdge2D(vertices[0], vertices[1]))
	edges.PushBack(tspmodel2d.NewEdge2D(vertices[1], vertices[2]))
	edges.PushBack(tspmodel2d.NewEdge2D(vertices[2], vertices[3]))
	edges.PushBack(tspmodel2d.NewEdge2D(vertices[3], vertices[4]))
	edges.PushBack(tspmodel2d.NewEdge2D(vertices[4], vertices[5]))
	edges.PushBack(tspmodel2d.NewEdge2D(vertices[5], vertices[6]))
	edges.PushBack(tspmodel2d.NewEdge2D(vertices[6], vertices[7]))
	edges.PushBack(tspmodel2d.NewEdge2D(vertices[7], vertices[0]))

	a, b, merged := tspmodel.MergeEdgesList(edges, vertices[1])
	assert.Equal(7, edges.Len())
	expected := []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[0]),
	}
	for i, testNode := 0, edges.Front(); testNode != nil; i, testNode = i+1, testNode.Next() {
		assert.Equal(expected[i], testNode.Value)
	}
	assert.Equal(tspmodel2d.NewEdge2D(vertices[0], vertices[1]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[1], vertices[2]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[0], vertices[2]), merged.Value)

	a, b, merged = tspmodel.MergeEdgesList(edges, tspmodel2d.NewVertex2D(10, 10))
	assert.Equal(7, edges.Len())
	expected = []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[0]),
	}
	for i, testNode := 0, edges.Front(); testNode != nil; i, testNode = i+1, testNode.Next() {
		assert.Equal(expected[i], testNode.Value)
	}
	assert.Nil(a)
	assert.Nil(b)
	assert.Nil(merged)

	a, b, merged = tspmodel.MergeEdgesList(edges, vertices[0])
	assert.Equal(6, edges.Len())
	expected = []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[2]),
	}
	for i, testNode := 0, edges.Front(); testNode != nil; i, testNode = i+1, testNode.Next() {
		assert.Equal(expected[i], testNode.Value, i)
	}
	assert.Equal(tspmodel2d.NewEdge2D(vertices[7], vertices[0]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[0], vertices[2]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[7], vertices[2]), merged.Value)

	a, b, merged = tspmodel.MergeEdgesList(edges, vertices[7])
	assert.Equal(5, edges.Len())
	expected = []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[2]),
	}
	for i, testNode := 0, edges.Front(); testNode != nil; i, testNode = i+1, testNode.Next() {
		assert.Equal(expected[i], testNode.Value)
	}
	assert.Equal(tspmodel2d.NewEdge2D(vertices[6], vertices[7]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[7], vertices[2]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[6], vertices[2]), merged.Value)

	a, b, merged = tspmodel.MergeEdgesList(edges, vertices[3])
	assert.Equal(4, edges.Len())
	expected = []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[2], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[2]),
	}
	for i, testNode := 0, edges.Front(); testNode != nil; i, testNode = i+1, testNode.Next() {
		assert.Equal(expected[i], testNode.Value)
	}
	assert.Equal(tspmodel2d.NewEdge2D(vertices[2], vertices[3]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[3], vertices[4]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[2], vertices[4]), merged.Value)

	_, _, _ = tspmodel.MergeEdgesList(edges, vertices[4])
	_, _, _ = tspmodel.MergeEdgesList(edges, vertices[5])
	assert.Equal(2, edges.Len())
	_, _, _ = tspmodel.MergeEdgesList(edges, vertices[6])
	assert.Equal(1, edges.Len())
	_, _, _ = tspmodel.MergeEdgesList(edges, vertices[2])
	assert.Equal(1, edges.Len())
	_, _, _ = tspmodel.MergeEdgesList(edges, vertices[2])
	assert.Equal(1, edges.Len())
}

func TestMoveVertex(t *testing.T) {
	assert := assert.New(t)

	vertices := []*tspmodel2d.Vertex2D{
		tspmodel2d.NewVertex2D(1, 1),
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(8, 8),
	}

	edges := []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[0]),
	}

	edges, a, b, c := tspmodel.MoveVertex(edges, vertices[0], tspmodel2d.NewEdge2D(vertices[3], vertices[4]))
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[0]),
		tspmodel2d.NewEdge2D(vertices[0], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[1]),
	}, edges)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[7], vertices[1]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[3], vertices[0]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[0], vertices[4]), c)

	edges, a, b, c = tspmodel.MoveVertex(edges, vertices[7], tspmodel2d.NewEdge2D(vertices[1], vertices[2]))
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[1], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[0]),
		tspmodel2d.NewEdge2D(vertices[0], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[1]),
	}, edges)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[6], vertices[1]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[1], vertices[7]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[7], vertices[2]), c)

	edges, a, b, c = tspmodel.MoveVertex(edges, vertices[4], tspmodel2d.NewEdge2D(vertices[6], vertices[1]))
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[1], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[0]),
		tspmodel2d.NewEdge2D(vertices[0], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[1]),
	}, edges)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[0], vertices[5]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[6], vertices[4]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[4], vertices[1]), c)

	edges, a, b, c = tspmodel.MoveVertex(edges, tspmodel2d.NewVertex2D(9, 9), tspmodel2d.NewEdge2D(vertices[6], vertices[1]))
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[1], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[0]),
		tspmodel2d.NewEdge2D(vertices[0], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[1]),
	}, edges)
	assert.Nil(a)
	assert.Nil(b)
	assert.Nil(c)

	edges, a, b, c = tspmodel.MoveVertex(edges, vertices[7], tspmodel2d.NewEdge2D(vertices[6], vertices[1]))
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[1], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[0]),
		tspmodel2d.NewEdge2D(vertices[0], vertices[5]),
		tspmodel2d.NewEdge2D(vertices[5], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[1]),
	}, edges)
	assert.Nil(a)
	assert.Nil(b)
	assert.Nil(c)

	edges2 := []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[1], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[1]),
	}
	edges2, a, b, c = tspmodel.MoveVertex(edges2, vertices[7], tspmodel2d.NewEdge2D(vertices[7], vertices[1]))
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[1], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[1]),
	}, edges2)
	assert.Nil(a)
	assert.Nil(b)
	assert.Nil(c)

	edges2 = []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[1], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[1]),
	}
	edges2, a, b, c = tspmodel.MoveVertex(edges2, vertices[7], tspmodel2d.NewEdge2D(vertices[2], vertices[1]))
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[1]),
	}, edges2)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[1], vertices[2]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[2], vertices[7]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[7], vertices[1]), c)

	edges2, a, b, c = tspmodel.MoveVertex(edges2, vertices[1], tspmodel2d.NewEdge2D(vertices[2], vertices[7]))
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[2], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[2]),
	}, edges2)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[7], vertices[2]), a)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[2], vertices[1]), b)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[1], vertices[7]), c)
}

func TestSplitEdge(t *testing.T) {
	assert := assert.New(t)

	vertices := []*tspmodel2d.Vertex2D{
		tspmodel2d.NewVertex2D(1, 1),
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(8, 8),
	}

	edges := []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
	}

	var index int
	edges, index = tspmodel.SplitEdge(edges, tspmodel2d.NewEdge2D(vertices[0], vertices[7]), vertices[6])
	assert.Equal(-1, index)
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
	}, edges)

	edges, index = tspmodel.SplitEdge(edges, edges[0], vertices[6])
	assert.Equal(0, index)
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
	}, edges)

	edges, index = tspmodel.SplitEdge(edges, edges[5], vertices[7])
	assert.Equal(5, index)
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[5]),
	}, edges)
}

func TestSplitEdgeCopy(t *testing.T) {
	assert := assert.New(t)

	vertices := []*tspmodel2d.Vertex2D{
		tspmodel2d.NewVertex2D(1, 1),
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(8, 8),
	}

	edges := []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
	}

	var index int
	edges, index = tspmodel.SplitEdgeCopy(edges, tspmodel2d.NewEdge2D(vertices[0], vertices[7]), vertices[6])
	assert.Equal(-1, index)
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
	}, edges)

	edges, index = tspmodel.SplitEdgeCopy(edges, edges[0], vertices[6])
	assert.Equal(0, index)
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
	}, edges)

	edges, index = tspmodel.SplitEdgeCopy(edges, edges[5], vertices[7])
	assert.Equal(5, index)
	assert.Equal([]tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[5]),
	}, edges)
}

func TestSplitEdgeList(t *testing.T) {
	assert := assert.New(t)

	vertices := []*tspmodel2d.Vertex2D{
		tspmodel2d.NewVertex2D(1, 1),
		tspmodel2d.NewVertex2D(2, 2),
		tspmodel2d.NewVertex2D(3, 3),
		tspmodel2d.NewVertex2D(4, 4),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6, 6),
		tspmodel2d.NewVertex2D(7, 7),
		tspmodel2d.NewVertex2D(8, 8),
	}

	edges := list.New()
	edges.PushBack(tspmodel2d.NewEdge2D(vertices[0], vertices[1]))
	edges.PushBack(tspmodel2d.NewEdge2D(vertices[1], vertices[2]))
	edges.PushBack(tspmodel2d.NewEdge2D(vertices[2], vertices[3]))
	edges.PushBack(tspmodel2d.NewEdge2D(vertices[3], vertices[4]))
	edges.PushBack(tspmodel2d.NewEdge2D(vertices[4], vertices[5]))

	newEdge := tspmodel.SplitEdgeList(edges, tspmodel2d.NewEdge2D(vertices[0], vertices[7]), vertices[6])
	assert.Nil(newEdge)
	expected := []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
	}
	for i, testNode := 0, edges.Front(); testNode != nil; i, testNode = i+1, testNode.Next() {
		assert.Equal(expected[i], testNode.Value, i)
	}

	newEdge = tspmodel.SplitEdgeList(edges, tspmodel2d.NewEdge2D(vertices[0], vertices[1]), vertices[6])
	assert.NotNil(newEdge)
	assert.NotNil(newEdge.Prev())
	assert.Equal(tspmodel2d.NewEdge2D(vertices[0], vertices[6]), newEdge.Prev().Value)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[6], vertices[1]), newEdge.Value)
	expected = []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[5]),
	}
	for i, testNode := 0, edges.Front(); testNode != nil; i, testNode = i+1, testNode.Next() {
		assert.Equal(expected[i], testNode.Value, i)
	}

	newEdge = tspmodel.SplitEdgeList(edges, tspmodel2d.NewEdge2D(vertices[4], vertices[5]), vertices[7])
	assert.NotNil(newEdge)
	assert.NotNil(newEdge.Prev())
	assert.Equal(tspmodel2d.NewEdge2D(vertices[4], vertices[7]), newEdge.Prev().Value)
	assert.Equal(tspmodel2d.NewEdge2D(vertices[7], vertices[5]), newEdge.Value)
	expected = []tspmodel.CircuitEdge{
		tspmodel2d.NewEdge2D(vertices[0], vertices[6]),
		tspmodel2d.NewEdge2D(vertices[6], vertices[1]),
		tspmodel2d.NewEdge2D(vertices[1], vertices[2]),
		tspmodel2d.NewEdge2D(vertices[2], vertices[3]),
		tspmodel2d.NewEdge2D(vertices[3], vertices[4]),
		tspmodel2d.NewEdge2D(vertices[4], vertices[7]),
		tspmodel2d.NewEdge2D(vertices[7], vertices[5]),
	}
	for i, testNode := 0, edges.Front(); testNode != nil; i, testNode = i+1, testNode.Next() {
		assert.Equal(expected[i], testNode.Value, i)
	}
}
