package solver_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/fealos/lee-tsp-go/solver"
)

func BenchmarkFindShortestPathNP(b *testing.B) {
	b.Run("N=8 Heap", func(b *testing.B) { benchmarkFindShortestPathNPHelper(8, b, solver.FindShortestPathNPHeap) })
	b.Run("N=8 With Checks", func(b *testing.B) { benchmarkFindShortestPathNPHelper(8, b, solver.FindShortestPathNP) })
	b.Run("N=8 No Checks", func(b *testing.B) { benchmarkFindShortestPathNPHelper(8, b, solver.FindShortestPathNPNoChecks) })
	b.Run("N=10 Heap", func(b *testing.B) { benchmarkFindShortestPathNPHelper(10, b, solver.FindShortestPathNPHeap) })
	b.Run("N=10 With Checks", func(b *testing.B) { benchmarkFindShortestPathNPHelper(10, b, solver.FindShortestPathNP) })
	b.Run("N=10 No Checks", func(b *testing.B) { benchmarkFindShortestPathNPHelper(10, b, solver.FindShortestPathNPNoChecks) })
	b.Run("N=12 Heap", func(b *testing.B) { benchmarkFindShortestPathNPHelper(12, b, solver.FindShortestPathNPHeap) })
	b.Run("N=12 With Checks", func(b *testing.B) { benchmarkFindShortestPathNPHelper(12, b, solver.FindShortestPathNP) })
	// b.Run("N=12 No Checks", func(b *testing.B) { benchmarkFindShortestPathNPHelper(12, b, solver.FindShortestPathNPNoChecks) })
}

func benchmarkFindShortestPathNPHelper(numVertices int, b *testing.B, solverFunc func([]model.CircuitVertex) ([]model.CircuitVertex, float64)) {
	for i := 0; i < b.N; i++ {
		vertices := model2d.GenerateVertices2D(numVertices)
		solverFunc(vertices)
	}
}
