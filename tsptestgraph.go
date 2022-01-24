package main

import (
	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/tspgraph"
	"github.com/fealos/lee-tsp-go/tspmodel"
	"github.com/fealos/lee-tsp-go/tspsolver"
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
	// 	circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
	// 		circuit, pathLength := tspsolver.FindShortestPathNPWithChecks(cv)
	// 		return &circuit.CompletedCircuit{
	// 			Circuit: circuit,
	// 			Length:  pathLength,
	// 		}
	// 	},
	// })

	// circuits = append(circuits, &NamedCircuit{
	// 	name: "np_heap",
	// 	circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
	// 		circuit, pathLength := tspsolver.FindShortestPathNPHeap(cv)
	// 		return &circuit.CompletedCircuit{
	// 			Circuit: circuit,
	// 			Length:  pathLength,
	// 		}
	// 	},
	// })

	circuits = append(circuits, &NamedCircuit{
		name: "greedy",
		circuitFunc: func(cv []tspmodel.CircuitVertex) tspmodel.Circuit {
			v := make([]*tspgraph.GraphVertex, len(cv))
			for i, vertex := range cv {
				v[i] = vertex.(*tspgraph.GraphVertex)
			}
			g := &tspgraph.Graph{
				Vertices: v,
			}
			c := tspgraph.NewGraphCircuit(g)
			defer c.Delete()

			tspsolver.FindShortestPathCircuit(c)

			return &circuit.CompletedCircuit{
				Circuit: c.GetAttachedVertices(),
				Length:  c.GetLength(),
			}
		},
	})

	ComparePerformance("results_graph_perf_greedy_1.tsv", &NumVertices{initVertices: 100, incrementVertices: 50, maxVertices: 1500, numIterations: 100}, circuits, func(size int) []tspmodel.CircuitVertex {
		gen := &tspgraph.GraphGenerator{
			NumVertices: uint32(size),
			MinEdges:    3,
			MaxEdges:    5,
		}

		g := gen.Create()

		cv := make([]tspmodel.CircuitVertex, len(g.Vertices))
		for i, vertex := range g.Vertices {
			cv[i] = vertex
		}
		return cv
	})
}
