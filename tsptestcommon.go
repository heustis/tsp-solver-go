package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fealos/lee-tsp-go/model"
)

func ComparePerformance(fileName string, verticesConfig *NumVertices, circuits []*NamedCircuit, verticesFunc func(size int) []model.CircuitVertex) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
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

	f.WriteString("num_vertices\t")
	for _, c := range circuits {
		fmt.Fprintf(f, "%s_len\t%s_perc\t%s_nanos\t", c.name, c.name, c.name)
	}
	f.WriteString("\r\n")

	for numVertices := verticesConfig.initVertices; numVertices <= verticesConfig.maxVertices; numVertices += verticesConfig.incrementVertices {
		for i := 0; i < verticesConfig.numIterations; i++ {
			vertices := verticesFunc(numVertices)
			circuitVerices := make([]model.CircuitVertex, numVertices)
			copy(circuitVerices, vertices)

			fmt.Fprintf(f, "%d\t", numVertices)
			minLen := -1.0
			for _, circuit := range circuits {
				t0 := time.Now()
				c := circuit.circuitFunc(vertices)
				t1 := time.Since(t0)

				circuitLen := c.GetLength()
				if minLen < 0 {
					minLen = circuitLen
				}
				// circuitJson, _ := json.Marshal(c.GetAttachedVertices())
				fmt.Fprintf(f, "%f\t%f\t%d\t", circuitLen, minLen/circuitLen, t1.Nanoseconds()) //, string(circuitJson))

				// if math.Abs(circuitLen-minLen) > model.Threshold {
				// 	fmt.Printf("test %d-%d: found mismatched circuits between %s and min solution\n", numVertices, i, circuitName)
				// }

				if d, okay := c.(model.Deletable); okay {
					d.Delete()
				}
			}
			f.WriteString("\r\n")

			if i > 0 && i%10 == 0 {
				fmt.Printf("completed: %d out of %d in %d vertices\r\n", i, verticesConfig.numIterations, numVertices)
			}
		}
	}
	// pprof.WriteHeapProfile(memProfileFile)
}

type NamedCircuit struct {
	name        string
	circuitFunc func([]model.CircuitVertex) model.Circuit
}

type NumVertices struct {
	initVertices      int
	incrementVertices int
	maxVertices       int
	numIterations     int
}
