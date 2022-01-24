package main

import (
	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/tspmodel"
	"github.com/fealos/lee-tsp-go/tspmodel3d"
	"github.com/fealos/lee-tsp-go/tspsolver"
)

func ComparePerformance3d() {
	// args := os.Args
	// if len(args) != 2 {
	// 	panic("Usage: " + args[0] + " <number_of_vertices>")
	// }

	// numVertices, err := strconv.Atoi(args[1])
	// if err != nil || numVertices < 3 {
	// 	panic("number_of_vertices must be an integer greater than 2")
	// }

	// Order circuits from most accurate to least accurate, so that ComparePerformance can create accurate percentages.
	circuits := []*NamedCircuit{}
	circuits = append(circuits, &NamedCircuit{
		name: "np",
		circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
			c, pathLength := tspsolver.FindShortestPathNPWithChecks(cv)
			return &circuit.CompletedCircuit{
				Circuit: c,
				Length:  pathLength,
			}
		},
	})

	circuits = append(circuits, &NamedCircuit{
		name: "np_heap",
		circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
			c, pathLength := tspsolver.FindShortestPathNPHeap(cv)
			return &circuit.CompletedCircuit{
				Circuit: c,
				Length:  pathLength,
			}
		},
	})

	circuits = append(circuits, &NamedCircuit{
		name: "heap",
		circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
			c := circuit.NewClonableCircuitSolver(circuit.NewHeapableCircuit(cv, tspmodel.DeduplicateVerticesNoSorting, tspmodel3d.BuildPerimiter))
			tspsolver.FindShortestPathCircuit(c)
			return c
		},
	})

	circuits = append(circuits, &NamedCircuit{
		name: "heap_mc",
		circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
			c := circuit.NewClonableCircuitSolver(circuit.NewHeapableCircuitMinClones(cv, tspmodel.DeduplicateVerticesNoSorting, tspmodel3d.BuildPerimiter))
			tspsolver.FindShortestPathCircuit(c)
			return c
		},
	})

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "convex_concave_byedge_withupdates",
	// 	circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
	// 		c := circuit.NewConvexConcaveByEdge(cv, tspmodel.DeduplicateVerticesNoSorting, tspmodel3d.BuildPerimiter, true)
	// 		tspsolver.FindShortestPathCircuit(c)
	// 		return c
	// 	},
	// })

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "convex_concave_byedge",
	// 	circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
	// 		c := circuit.NewConvexConcaveByEdge(cv, tspmodel.DeduplicateVerticesNoSorting, tspmodel3d.BuildPerimiter, false)
	// 		tspsolver.FindShortestPathCircuit(c)
	// 		return c
	// 	},
	// })

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "convex_concave_withupdates",
	// 	circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
	// 		c := circuit.NewConvexConcave(cv, tspmodel.DeduplicateVerticesNoSorting, tspmodel3d.BuildPerimiter, true)
	// 		tspsolver.FindShortestPathCircuit(c)
	// 		return c
	// 	},
	// })

	circuits = append(circuits, &NamedCircuit{
		name: "convex_concave",
		circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
			c := circuit.NewConvexConcave(cv, tspmodel.DeduplicateVerticesNoSorting, tspmodel3d.BuildPerimiter, false)
			tspsolver.FindShortestPathCircuit(c)
			return c
		},
	})

	ComparePerformance("results_3d_comp_np_1.tsv", &NumVertices{initVertices: 7, incrementVertices: 1, maxVertices: 15, numIterations: 100}, circuits, tspmodel3d.GenerateVertices)
}
