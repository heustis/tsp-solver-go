package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/fealos/lee-tsp-go/solver"
)

func main() {
	// args := os.Args
	// if len(args) != 2 {
	// 	panic("Usage: " + args[0] + " <number_of_vertices>")
	// }

	// numVertices, err := strconv.Atoi(args[1])
	// if err != nil || numVertices < 3 {
	// 	panic("number_of_vertices must be an integer greater than 2")
	// }

	f, err := os.OpenFile("results_2d_comp.tsv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// cpuProfileFile, err := os.Create("tsptest_cpu.prof")
	// if err != nil {
	// 	panic(err)
	// }
	// defer cpuProfileFile.Close()

	// memProfileFile, err := os.Create("tsptest_mem.prof")
	// if err != nil {
	// 	panic(err)
	// }
	// defer memProfileFile.Close()

	// err = pprof.StartCPUProfile(cpuProfileFile)
	// if err != nil {
	// 	panic(err)
	// }
	// defer pprof.StopCPUProfile()

	// f.WriteString("np_len;np_nanos;heap_len;heap_nanos;heap_iterations;heap_clones;heap_mc_len;heap_mc_nanos;heap_mc_iterations;heap_mc_clones;np_circuit;heap_circuit;heap_mc_circuit;\r\n")
	f.WriteString("num_vertices\theap_mc_len\theap_mc_nanos\tgready_len\tgready_len_perc\tgreedy_nanos\t\r\n")
	// f.WriteString("greedy_len\tgreedy_nanos\r\n")

	numIterations := 50
	for numVertices := 28; numVertices <= 28; numVertices += 5 {
		for i := 0; i < numIterations; i++ {
			vertices := generateVertices(numVertices)
			circuitVerices := make([]model.CircuitVertex, numVertices)
			copy(circuitVerices, vertices)

			t0 := time.Now()
			// shortestNp, shortestNpLength := solver.FindShortestPathNPHeap(circuitVerices)

			// t1 := time.Since(t0)

			// shortestHeap, numIterationsHeap, numClonesHeap := solver.FindShortestPathHeap(model.CreateHeapableCircuitImpl(vertices, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{}))
			// shortestHeap, _, _ := solver.FindShortestPathHeap(model.CreateHeapableCircuitImpl(vertices, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{}))

			// t2 := time.Since(t0)

			// shortestHeapMinClones, numIterationsMinClones, numClonesMinHeap := solver.FindShortestPathHeap(model.CreateHeapableCircuitMinClones(vertices, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{}))
			shortestHeapMinClones, _, _ := solver.FindShortestPathHeap(model.CreateHeapableCircuitMinClones(vertices, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{}))

			t3 := time.Since(t0)

			greedyCircuit := model.NewCircuitGreedyImpl(vertices, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
			solver.FindShortestPathGreedy(greedyCircuit)

			t4 := time.Since(t0)

			// greedyCircuitWithUpdates := model.NewCircuitGreedyWithUpdatesImpl(vertices, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
			// solver.FindShortestPathGreedy(greedyCircuitWithUpdates)

			// t5 := time.Since(t0)

			// shortestNpJson, _ := json.Marshal(shortestNp)
			// initialJson, _ := json.Marshal(vertices)
			// shortestHeapJson, _ := json.Marshal(shortestHeap.GetAttachedVertices())
			// shortestHeapMinClonesJson, _ := json.Marshal(shortestHeapMinClones.GetAttachedVertices())
			// shortestGreedyCircuitJson, _ := json.Marshal(greedyCircuit.GetAttachedVertices())
			// shortestGreedyCircuitWithUpdatesJson, _ := json.Marshal(greedyCircuitWithUpdates.GetAttachedVertices())

			f.WriteString(fmt.Sprintf("%d\t%f\t%d\t%f\t%f\t%d\r\n",
				numVertices,
				// shortestNpLength, t1.Nanoseconds(),
				// shortestHeap.GetLength(), t2.Nanoseconds(), //-t1.Nanoseconds(), numIterationsHeap, numClonesHeap,
				shortestHeapMinClones.GetLength(), t3.Nanoseconds(), //-t2.Nanoseconds(), //numIterationsMinClones, numClonesMinHeap,
				greedyCircuit.GetLength(), shortestHeapMinClones.GetLength()/greedyCircuit.GetLength(), t4.Nanoseconds()-t3.Nanoseconds(), //-t3.Nanoseconds(), shortestHeapMinClones.GetLength()/greedyCircuit.GetLength(),
				// greedyCircuitWithUpdates.GetLength(), t5.Nanoseconds()-t4.Nanoseconds(), //shortestHeapMinClones.GetLength()/greedyCircuitWithUpdates.GetLength(),

				// string(shortestNpJson),
				// string(initialJson),
				// string(shortestHeapJson),
				// string(shortestHeapMinClonesJson),
				// string(shortestGreedyCircuitJson),
				// string(shortestGreedyCircuitWithUpdatesJson),
			))

			// if math.Abs(shortestHeap.GetLength()-shortestHeapMinClones.GetLength()) > model.Threshold {
			// 	fmt.Printf("test %d: found mismatched circuits between Heap (Full) and Heap (Min Clones) solutions\n", i)
			// }

			// if math.Abs(shortestHeapMinClones.GetLength()-shortestNpLength) > model.Threshold {
			// 	fmt.Printf("test %d: found mismatched circuits between NP and Heap Min Clones solutions\n", i)
			// }

			// shortestHeap.Delete()
			shortestHeapMinClones.Delete()

			if i > 0 && i%10 == 0 {
				fmt.Printf("completed: %d out of %d\r\n", i, numIterations)
			}
		}
	}
	// pprof.WriteHeapProfile(memProfileFile)
}

func generateVertices(size int) []model.CircuitVertex {
	var vertices []model.CircuitVertex
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		vertices = append(vertices, model2d.NewVertex2D(r.Float64()*10000, r.Float64()*10000))
	}
	return vertices
}
