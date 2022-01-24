package main

import (
	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/tspmodel"
	"github.com/fealos/lee-tsp-go/tspmodel2d"
	"github.com/fealos/lee-tsp-go/tspsolver"
)

func ComparePerformance2d() {
	// args := os.Args
	// if len(args) != 2 {
	// 	panic("Usage: " + args[0] + " <number_of_vertices>")
	// }

	// numVertices, err := strconv.Atoi(args[1])
	// if err != nil || numVertices < 3 {
	// 	panic("number_of_vertices must be an integer greater than 2")
	// }

	circuits := []*NamedCircuit{}
	// circuits = append(circuits, &NamedCircuit{
	// 	name: "heap_mc",
	// 	circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
	// 		c := circuit.NewClonableCircuitSolver(circuit.NewHeapableCircuitMinClones(cv, tspmodel2d.DeduplicateVertices, tspmodel2d.BuildPerimiter))
	// 		tspsolver.FindShortestPathCircuit(c)
	// 		return c
	// 	},
	// })

	// Do convex concave by edge first, since it is the most comprehensive convex concave algorithm, so that the other algorithms are compared to it.
	// circuits = append(circuits, &NamedCircuit{
	// 	name: "convex_concave_byedge_withupdates",
	// 	circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
	// 		c := circuit.NewConvexConcaveByEdge(cv, tspmodel2d.DeduplicateVertices, tspmodel2d.BuildPerimiter, true)
	// 		tspsolver.FindShortestPathCircuit(c)
	// 		return c
	// 	},
	// })

	circuits = append(circuits, &NamedCircuit{
		name: "convex_concave_byedge",
		circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
			c := circuit.NewConvexConcaveByEdge(cv, tspmodel2d.DeduplicateVertices, tspmodel2d.BuildPerimiter, false)
			tspsolver.FindShortestPathCircuit(c)
			return c
		},
	})

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "convex_concave_withupdates",
	// 	circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
	// 		circuit := circuit.NewConvexConcave(cv, tspmodel2d.DeduplicateVertices, tspmodel2d.BuildPerimiter, true)
	// 		tspsolver.FindShortestPathCircuit(circuit)
	// 		return circuit
	// 	},
	// })

	circuits = append(circuits, &NamedCircuit{
		name: "convex_concave",
		circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
			c := circuit.NewConvexConcave(cv, tspmodel2d.DeduplicateVertices, tspmodel2d.BuildPerimiter, false)
			tspsolver.FindShortestPathCircuit(c)
			return c
		},
	})

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "convex_concave_weighted_edge",
	// 	circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
	// 		c := experimental.NewConvexConcaveWeightedEdges(cv, tspmodel2d.DeduplicateVertices, tspmodel2d.BuildPerimiter)
	// 		tspsolver.FindShortestPathCircuit(c)
	// 		return c
	// 	},
	// })

	ComparePerformance("results_2d_comp_convex_concave_3.tsv", &NumVertices{initVertices: 100, incrementVertices: 100, maxVertices: 2000, numIterations: 100}, circuits, tspmodel2d.GenerateVertices)
}
