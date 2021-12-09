package main

import (
	"math/rand"
	"time"

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
	// 		circuit, _, _ := solver.FindShortestPathHeap(model.NewHeapableCircuitMinClones(cv, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{}))
	// 		return circuit.(model.Circuit)
	// 	},
	// })

	// Do greedy by edge first, since it is the most comprehensive greedy algorithm, so that the other algorithms are compared to it.
	// circuits = append(circuits, &NamedCircuit{
	// 	name: "greedy_byedge_withupdates",
	// 	circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
	// 		circuit := model.NewCircuitGreedyByEdgeWithUpdatesImpl(cv, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
	// 		solver.FindShortestPathGreedy(circuit)
	// 		return circuit
	// 	},
	// })

	circuits = append(circuits, &NamedCircuit{
		name: "greedy_byedge",
		circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
			circuit := model.NewCircuitGreedyByEdgeImpl(cv, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
			solver.FindShortestPathGreedy(circuit)
			return circuit
		},
	})

	circuits = append(circuits, &NamedCircuit{
		name: "greedy",
		circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
			circuit := model.NewCircuitGreedyImpl(cv, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
			solver.FindShortestPathGreedy(circuit)
			return circuit
		},
	})

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "greedy_weighted_edge",
	// 	circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
	// 		circuit := model.NewCircuitGreedyWeightedEdgeImpl(cv, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
	// 		solver.FindShortestPathGreedy(circuit)
	// 		return circuit
	// 	},
	// })

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "greedy_withupdates",
	// 	circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
	// 		circuit := model.NewCircuitGreedyWithUpdatesImpl(cv, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
	// 		solver.FindShortestPathGreedy(circuit)
	// 		return circuit
	// 	},
	// })

	ComparePerformance("results_2d_comp_greedy_3.tsv", &NumVertices{initVertices: 100, incrementVertices: 100, maxVertices: 2000, numIterations: 100}, circuits, GenerateVertices2d)
}

func GenerateVertices2d(size int) []model.CircuitVertex {
	var vertices []model.CircuitVertex
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		vertices = append(vertices, model2d.NewVertex2D(r.Float64()*10000, r.Float64()*10000))
	}
	return vertices
}
