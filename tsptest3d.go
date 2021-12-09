package main

import (
	"math/rand"
	"time"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model3d"
	"github.com/fealos/lee-tsp-go/solver"
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
			circuit, pathLength := solver.FindShortestPathNP(cv)
			return &model.HeapableCircuitComplete{
				Circuit: circuit,
				Length:  pathLength,
			}
		},
	})

	circuits = append(circuits, &NamedCircuit{
		name: "np_heap",
		circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
			circuit, pathLength := solver.FindShortestPathNPHeap(cv)
			return &model.HeapableCircuitComplete{
				Circuit: circuit,
				Length:  pathLength,
			}
		},
	})

	circuits = append(circuits, &NamedCircuit{
		name: "heap",
		circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
			circuit, _, _ := solver.FindShortestPathHeap(model.NewHeapableCircuitImpl(cv, model3d.DeduplicateVertices3D, &model3d.PerimeterBuilder3D{}))
			return circuit.(model.Circuit)
		},
	})

	circuits = append(circuits, &NamedCircuit{
		name: "heap_mc",
		circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
			circuit, _, _ := solver.FindShortestPathHeap(model.NewHeapableCircuitMinClones(cv, model3d.DeduplicateVertices3D, &model3d.PerimeterBuilder3D{}))
			return circuit.(model.Circuit)
		},
	})

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "greedy_byedge_withupdates",
	// 	circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
	// 		circuit := model.NewCircuitGreedyByEdgeWithUpdatesImpl(cv, model3d.DeduplicateVertices3D, &model3d.PerimeterBuilder3D{})
	// 		solver.FindShortestPathGreedy(circuit)
	// 		return circuit
	// 	},
	// })

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "greedy_byedge",
	// 	circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
	// 		circuit := model.NewCircuitGreedyByEdgeImpl(cv, model3d.DeduplicateVertices3D, &model3d.PerimeterBuilder3D{})
	// 		solver.FindShortestPathGreedy(circuit)
	// 		return circuit
	// 	},
	// })

	circuits = append(circuits, &NamedCircuit{
		name: "greedy",
		circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
			circuit := model.NewCircuitGreedyImpl(cv, model3d.DeduplicateVertices3D, &model3d.PerimeterBuilder3D{})
			solver.FindShortestPathGreedy(circuit)
			return circuit
		},
	})

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "greedy_withupdates",
	// 	circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
	// 		circuit := model.NewCircuitGreedyWithUpdatesImpl(cv, model3d.DeduplicateVertices3D, &model3d.PerimeterBuilder3D{})
	// 		solver.FindShortestPathGreedy(circuit)
	// 		return circuit
	// 	},
	// })

	ComparePerformance("results_3d_comp_np_1.tsv", &NumVertices{initVertices: 7, incrementVertices: 1, maxVertices: 15, numIterations: 100}, circuits, GenerateVertices3d)
}

func GenerateVertices3d(size int) []model.CircuitVertex {
	var vertices []model.CircuitVertex
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		vertices = append(vertices, model3d.NewVertex3D(r.Float64()*10000, r.Float64()*10000, r.Float64()*10000))
	}
	return vertices
}
