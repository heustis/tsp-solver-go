package tspsolver_test

import (
	"fmt"
	"testing"

	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/tspmodel"
	"github.com/fealos/lee-tsp-go/tspmodel2d"
	"github.com/fealos/lee-tsp-go/tspsolver"
)

func BenchmarkFindShortestPathCircuit(b *testing.B) {
	for n := 50; n <= 250; n += 50 {
		b.Run(fmt.Sprintf("N=%d ConcaveConvex", n), func(b *testing.B) {
			benchmarkFindShortestPathCircuitHelper(n, b, func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
				return circuit.NewConvexConcave(cv, tspmodel2d.DeduplicateVertices, tspmodel2d.BuildPerimiter, false)
			})
		})
		b.Run(fmt.Sprintf("N=%d ConcaveConvex With Checks", n), func(b *testing.B) {
			benchmarkFindShortestPathCircuitHelper(n, b, func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
				return circuit.NewConvexConcave(cv, tspmodel2d.DeduplicateVertices, tspmodel2d.BuildPerimiter, true)
			})
		})
		b.Run(fmt.Sprintf("N=%d  ConcaveConvex By Edge", n), func(b *testing.B) {
			benchmarkFindShortestPathCircuitHelper(n, b, func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
				return circuit.NewConvexConcaveByEdge(cv, tspmodel2d.DeduplicateVertices, tspmodel2d.BuildPerimiter, false)
			})
		})
		b.Run(fmt.Sprintf("N=%d ConcaveConvex By Edge With Checks", n), func(b *testing.B) {
			benchmarkFindShortestPathCircuitHelper(n, b, func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
				return circuit.NewConvexConcaveByEdge(cv, tspmodel2d.DeduplicateVertices, tspmodel2d.BuildPerimiter, true)
			})
		})
		b.Run(fmt.Sprintf("N=%d ConcaveConvex Disparity", n), func(b *testing.B) {
			benchmarkFindShortestPathCircuitHelper(n, b, func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
				return circuit.NewConvexConcaveDisparity(cv, tspmodel2d.DeduplicateVertices, tspmodel2d.BuildPerimiter, false)
			})
		})
	}
}

func benchmarkFindShortestPathCircuitHelper(numVertices int, b *testing.B, circuitFunc func([]tspmodel.CircuitVertex) tspmodel.Circuit) {
	for i := 0; i < b.N; i++ {
		vertices := tspmodel2d.GenerateVertices(numVertices)
		cir := circuitFunc(vertices)
		tspsolver.FindShortestPathCircuit(cir)
	}
}
