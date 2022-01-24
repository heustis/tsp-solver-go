package tspsolver_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/tspmodel"
	"github.com/fealos/lee-tsp-go/tspmodel2d"
	"github.com/fealos/lee-tsp-go/tspsolver"
)

func BenchmarkFindShortestPathNPWithChecks(b *testing.B) {
	b.Run("N=8 Heap", func(b *testing.B) { benchmarkFindShortestPathNPHelper(8, b, tspsolver.FindShortestPathNPHeap) })
	b.Run("N=8 With Checks", func(b *testing.B) { benchmarkFindShortestPathNPHelper(8, b, tspsolver.FindShortestPathNPWithChecks) })
	b.Run("N=8 No Checks", func(b *testing.B) { benchmarkFindShortestPathNPHelper(8, b, tspsolver.FindShortestPathNPNoChecks) })
	b.Run("N=10 Heap", func(b *testing.B) { benchmarkFindShortestPathNPHelper(10, b, tspsolver.FindShortestPathNPHeap) })
	b.Run("N=10 With Checks", func(b *testing.B) { benchmarkFindShortestPathNPHelper(10, b, tspsolver.FindShortestPathNPWithChecks) })
	b.Run("N=10 No Checks", func(b *testing.B) { benchmarkFindShortestPathNPHelper(10, b, tspsolver.FindShortestPathNPNoChecks) })
	b.Run("N=12 Heap", func(b *testing.B) { benchmarkFindShortestPathNPHelper(12, b, tspsolver.FindShortestPathNPHeap) })
	b.Run("N=12 With Checks", func(b *testing.B) { benchmarkFindShortestPathNPHelper(12, b, tspsolver.FindShortestPathNPWithChecks) })
	// b.Run("N=12 No Checks", func(b *testing.B) { benchmarkFindShortestPathNPHelper(12, b, tspsolver.FindShortestPathNPNoChecks) })
}

func benchmarkFindShortestPathNPHelper(numVertices int, b *testing.B, solverFunc func([]tspmodel.CircuitVertex) ([]tspmodel.CircuitVertex, float64)) {
	for i := 0; i < b.N; i++ {
		vertices := tspmodel2d.GenerateVertices(numVertices)
		solverFunc(vertices)
	}
}
