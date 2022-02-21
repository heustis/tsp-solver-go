package main

import (
	"github.com/heustis/tsp-solver-go/circuit"
	"github.com/heustis/tsp-solver-go/model"
	"github.com/heustis/tsp-solver-go/model3d"
	"github.com/heustis/tsp-solver-go/solver"
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
		circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
			c, pathLength := solver.FindShortestPathNPWithChecks(cv)
			return &circuit.CompletedCircuit{
				Circuit: c,
				Length:  pathLength,
			}
		},
	})

	circuits = append(circuits, &NamedCircuit{
		name: "np_heap",
		circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
			c, pathLength := solver.FindShortestPathNPHeap(cv)
			return &circuit.CompletedCircuit{
				Circuit: c,
				Length:  pathLength,
			}
		},
	})

	circuits = append(circuits, &NamedCircuit{
		name: "heap",
		circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
			cImpl := circuit.NewClonableCircuitImpl(cv, model3d.BuildPerimiter)
			cImpl.CloneOnFirstAttach = true
			c := circuit.NewClonableCircuitSolver(cImpl)
			solver.FindShortestPathCircuit(c)
			return c
		},
	})

	circuits = append(circuits, &NamedCircuit{
		name: "heap_mc",
		circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
			c := circuit.NewClonableCircuitSolver(circuit.NewClonableCircuitImpl(cv, model3d.BuildPerimiter))
			solver.FindShortestPathCircuit(c)
			return c
		},
	})

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "convex_concave_byedge_withupdates",
	// 	circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
	// 		c := circuit.NewConvexConcaveByEdge(cv, model3d.BuildPerimiter, true)
	// 		solver.FindShortestPathCircuit(c)
	// 		return c
	// 	},
	// })

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "convex_concave_byedge",
	// 	circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
	// 		c := circuit.NewConvexConcaveByEdge(cv, model3d.BuildPerimiter, false)
	// 		solver.FindShortestPathCircuit(c)
	// 		return c
	// 	},
	// })

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "convex_concave_withupdates",
	// 	circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
	// 		c := circuit.NewConvexConcave(cv, model3d.BuildPerimiter, true)
	// 		solver.FindShortestPathCircuit(c)
	// 		return c
	// 	},
	// })

	circuits = append(circuits, &NamedCircuit{
		name: "convex_concave",
		circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
			c := circuit.NewConvexConcave(cv, model3d.BuildPerimiter, false)
			solver.FindShortestPathCircuit(c)
			return c
		},
	})

	ComparePerformance("results_3d_comp_np_1.tsv", &NumVertices{initVertices: 7, incrementVertices: 1, maxVertices: 15, numIterations: 100}, circuits, func(size int) []model.CircuitVertex {
		return model.DeduplicateVerticesNoSorting(model3d.GenerateVertices(size))
	})
}
