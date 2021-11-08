package main

import (
	"encoding/json"
	"fmt"
	"math"
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

	f, err := os.OpenFile("results.scsv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	f.WriteString("np_len;np_nanos;heap_len;heap_nanos;heap_iterations;heap_mc_len;heap_mc_nanos;heap_mc_iterations;np_circuit;heap_circuit;heap_mc_circuit;\r\n")

	numVertices := 10
	numIterations := 1000
	for i := 0; i < numIterations; i++ {
		vertices := generateVertices(numVertices)
		circuitVerices := make([]model.CircuitVertex, numVertices)
		for i, v := range vertices {
			circuitVerices[i] = v
		}

		t0 := time.Now()
		shortestNp, shortestNpLength := solver.FindShortestPathNP(circuitVerices)

		t1 := time.Since(t0)

		heapCircuit := &model2d.HeapableCircuit2D{Vertices: vertices}
		shortestHeap, numIterationsHeap := solver.FindShortestPathHeap(heapCircuit)

		t2 := time.Since(t0)

		heapCircuitMinClones := &model2d.HeapableCircuit2DMinClones{Vertices: vertices}
		shortestHeapMinClones, numIterationsMinClones := solver.FindShortestPathHeap(heapCircuitMinClones)

		t3 := time.Since(t0)

		shortestNpJson, _ := json.Marshal(shortestNp)
		// initialJson, _ := json.Marshal(vertices)
		shortestHeapJson, _ := json.Marshal(shortestHeap.GetAttachedVertices())
		shortestHeapMinClonesJson, _ := json.Marshal(shortestHeapMinClones.GetAttachedVertices())
		f.WriteString(fmt.Sprintf("%f;%d;%f;%d;%d;%f;%d;%d;%s;%s;%s;\r\n",
			shortestNpLength, t1.Nanoseconds(),
			shortestHeap.GetLength(), t2.Nanoseconds()-t1.Nanoseconds(), numIterationsHeap,
			shortestHeapMinClones.GetLength(), t3.Nanoseconds()-t2.Nanoseconds(), numIterationsMinClones,
			string(shortestNpJson),
			// string(initialJson),
			string(shortestHeapJson),
			string(shortestHeapMinClonesJson),
		))

		if math.Abs(shortestHeapMinClones.GetLength()-shortestNpLength) > model.Threshold {
			fmt.Printf("test %d: found mismatched circuits between NP and Heap solutions\n", i)
		}

		shortestHeap.Delete()
		shortestHeapMinClones.Delete()

		if i > 0 && i%10 == 0 {
			fmt.Printf("completed: %d out of %d\r\n", i, numIterations)
		}
	}
}

func generateVertices(size int) []model.CircuitVertex {
	var vertices []model.CircuitVertex
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		vertices = append(vertices, model2d.NewVertex2D(r.Float64()*10000, r.Float64()*10000))
	}
	return vertices
}
