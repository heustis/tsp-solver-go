package main

import (
	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/fealos/lee-tsp-go/solver"
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
	// 	circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
	// 		c, _, _ := solver.FindShortestPathHeap(model.NewHeapableCircuitMinClones(cv, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{}))
	// 		return c.(model.Circuit)
	// 	},
	// })

	// Do convex concave by edge first, since it is the most comprehensive convex concave algorithm, so that the other algorithms are compared to it.
	// circuits = append(circuits, &NamedCircuit{
	// 	name: "convex_concave_byedge_withupdates",
	// 	circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
	// 		c := circuit.NewConvexConcaveByEdge(cv, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{}, true)
	// 		solver.FindShortestPathGreedy(c)
	// 		return c
	// 	},
	// })

	circuits = append(circuits, &NamedCircuit{
		name: "convex_concave_byedge",
		circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
			c := circuit.NewConvexConcaveByEdge(cv, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{}, false)
			solver.FindShortestPathGreedy(c)
			return c
		},
	})

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "convex_concave_withupdates",
	// 	circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
	// 		circuit := circuit.NewConvexConcave(cv, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{}, true)
	// 		solver.FindShortestPathGreedy(circuit)
	// 		return circuit
	// 	},
	// })

	circuits = append(circuits, &NamedCircuit{
		name: "convex_concave",
		circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
			c := circuit.NewConvexConcave(cv, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{}, false)
			solver.FindShortestPathGreedy(c)
			return c
		},
	})

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "convex_concave_weighted_edge",
	// 	circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
	// 		c := experimental.NewConvexConcaveWeightedEdges(cv, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
	// 		solver.FindShortestPathGreedy(c)
	// 		return c
	// 	},
	// })

	ComparePerformance("results_2d_comp_convex_concave_3.tsv", &NumVertices{initVertices: 100, incrementVertices: 100, maxVertices: 2000, numIterations: 100}, circuits, model2d.GenerateVertices2D)
}
