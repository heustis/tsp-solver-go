package main

import (
	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/graph"
	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/solver"
)

func ComparePerformanceGraph() {
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

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "np",
	// 	circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
	// 		circuit, pathLength := solver.FindShortestPathNPWithChecks(cv)
	// 		return &circuit.CompletedCircuit{
	// 			Circuit: circuit,
	// 			Length:  pathLength,
	// 		}
	// 	},
	// })

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "np_heap",
	// 	circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
	// 		circuit, pathLength := solver.FindShortestPathNPHeap(cv)
	// 		return &circuit.CompletedCircuit{
	// 			Circuit: circuit,
	// 			Length:  pathLength,
	// 		}
	// 	},
	// })

	circuits = append(circuits, &NamedCircuit{
		name: "greedy",
		circuitFunc: func(cv []model.CircuitVertex) model.Circuit {
			v := make([]*graph.GraphVertex, len(cv))
			for i, vertex := range cv {
				v[i] = vertex.(*graph.GraphVertex)
			}
			g := &graph.Graph{
				Vertices: v,
			}
			c := graph.NewGraphCircuit(g)
			defer c.Delete()

			solver.FindShortestPathCircuit(c)

			return &circuit.CompletedCircuit{
				Circuit: c.GetAttachedVertices(),
				Length:  c.GetLength(),
			}
		},
	})

	ComparePerformance("results_graph_perf_greedy_1.tsv", &NumVertices{initVertices: 100, incrementVertices: 50, maxVertices: 1500, numIterations: 100}, circuits, func(size int) []model.CircuitVertex {
		gen := &graph.GraphGenerator{
			NumVertices: uint32(size),
			MinEdges:    3,
			MaxEdges:    5,
		}

		g := gen.Create()

		cv := make([]model.CircuitVertex, len(g.Vertices))
		for i, vertex := range g.Vertices {
			cv[i] = vertex
		}
		return cv
	})
}
