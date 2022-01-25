package solver_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/fealos/lee-tsp-go/solver"
	"github.com/stretchr/testify/assert"
)

func TestFindShortestPathCircuit_DataFromOldProject(t *testing.T) {
	assert := assert.New(t)

	dataBytes, err := ioutil.ReadFile("../test-data/sample-polygons.json")
	assert.Nil(err)

	data := &testEntries{}
	err = json.Unmarshal(dataBytes, data)
	assert.Nil(err)

	for testIndex, testEntry := range data.Arrays[0:10] {
		vertices := make([]model.CircuitVertex, len(testEntry.Points))
		for i, points := range testEntry.Points {
			vertices[i] = model2d.NewVertex2D(points[0], points[1])
		}
		cir := circuit.NewConvexConcave(vertices, model2d.DeduplicateVertices, model2d.BuildPerimiter, true)
		solver.FindShortestPathCircuit(cir)
		shortest := cir.GetAttachedVertices()
		actual := cir.GetLength()

		shortestString, err := json.Marshal(shortest)
		assert.Nil(err)
		// The greedy approximations will not perfectly solve these circuits; assert that they are within 10% of the optimal solution.
		assert.Greater(testEntry.Expected*1.1, actual, fmt.Sprintf("test=%d pathLength=%f shortestPath=%s", testIndex, actual, shortestString))
	}
}

func TestFindShortestPathCircuit_Heap_DataFromOldProject(t *testing.T) {
	assert := assert.New(t)

	dataBytes, err := ioutil.ReadFile("../test-data/sample-polygons.json")
	assert.Nil(err)

	data := &testEntries{}
	err = json.Unmarshal(dataBytes, data)
	assert.Nil(err)

	for testIndex, testEntry := range data.Arrays[0:10] {
		vertices := make([]model.CircuitVertex, len(testEntry.Points))
		for i, points := range testEntry.Points {
			vertices[i] = model2d.NewVertex2D(points[0], points[1])
		}
		c := circuit.NewHeapableCircuit(vertices, model2d.DeduplicateVertices, model2d.BuildPerimiter)
		cir := circuit.NewClonableCircuitSolver(c)
		solver.FindShortestPathCircuit(cir)
		shortest := cir.GetAttachedVertices()
		actual := cir.GetLength()

		shortestString, err := json.Marshal(shortest)
		assert.Nil(err)
		assert.InDelta(testEntry.Expected, actual, model.Threshold, fmt.Sprintf("test=%d pathLength=%f shortestPath=%s", testIndex, actual, shortestString))
	}
}

func TestFindShortestPathCircuit_HeapMinClones_DataFromOldProject(t *testing.T) {
	assert := assert.New(t)

	dataBytes, err := ioutil.ReadFile("../test-data/sample-polygons.json")
	assert.Nil(err)

	data := &testEntries{}
	err = json.Unmarshal(dataBytes, data)
	assert.Nil(err)

	for testIndex, testEntry := range data.Arrays[0:10] {
		vertices := make([]model.CircuitVertex, len(testEntry.Points))
		for i, points := range testEntry.Points {
			vertices[i] = model2d.NewVertex2D(points[0], points[1])
		}
		c := circuit.NewHeapableCircuitMinClones(vertices, model2d.DeduplicateVertices, model2d.BuildPerimiter)
		cir := circuit.NewClonableCircuitSolver(c)
		solver.FindShortestPathCircuit(cir)
		shortest := cir.GetAttachedVertices()
		actual := cir.GetLength()

		shortestString, err := json.Marshal(shortest)
		assert.Nil(err)
		assert.InDelta(testEntry.Expected, actual, model.Threshold, fmt.Sprintf("test=%d pathLength=%f shortestPath=%s", testIndex, actual, shortestString))
	}
}