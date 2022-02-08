package solver_test

import (
	"fmt"
	"testing"

	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/fealos/lee-tsp-go/solver"
)

func BenchmarkFindShortestPathCircuit(b *testing.B) {
	for n := 50; n <= 250; n += 50 {
		b.Run(fmt.Sprintf("N=%d ConcaveConvex", n), func(b *testing.B) {
			benchmarkFindShortestPathCircuitHelper(n, b, func(cv []model.CircuitVertex) model.Circuit {
				return circuit.NewConvexConcave(cv, model2d.BuildPerimiter, false)
			})
		})
		b.Run(fmt.Sprintf("N=%d ConcaveConvex With Checks", n), func(b *testing.B) {
			benchmarkFindShortestPathCircuitHelper(n, b, func(cv []model.CircuitVertex) model.Circuit {
				return circuit.NewConvexConcave(cv, model2d.BuildPerimiter, true)
			})
		})
		b.Run(fmt.Sprintf("N=%d  ConcaveConvex By Edge", n), func(b *testing.B) {
			benchmarkFindShortestPathCircuitHelper(n, b, func(cv []model.CircuitVertex) model.Circuit {
				return circuit.NewConvexConcaveByEdge(cv, model2d.BuildPerimiter, false)
			})
		})
		b.Run(fmt.Sprintf("N=%d ConcaveConvex By Edge With Checks", n), func(b *testing.B) {
			benchmarkFindShortestPathCircuitHelper(n, b, func(cv []model.CircuitVertex) model.Circuit {
				return circuit.NewConvexConcaveByEdge(cv, model2d.BuildPerimiter, true)
			})
		})
		b.Run(fmt.Sprintf("N=%d ConcaveConvex Disparity", n), func(b *testing.B) {
			benchmarkFindShortestPathCircuitHelper(n, b, func(cv []model.CircuitVertex) model.Circuit {
				return circuit.NewConvexConcaveDisparity(cv, model2d.BuildPerimiter, false)
			})
		})
	}
}

func benchmarkFindShortestPathCircuitHelper(numVertices int, b *testing.B, circuitFunc func([]model.CircuitVertex) model.Circuit) {
	for i := 0; i < b.N; i++ {
		vertices := model2d.DeduplicateVertices(model2d.GenerateVertices(numVertices))
		cir := circuitFunc(vertices)
		solver.FindShortestPathCircuit(cir)
	}
}
