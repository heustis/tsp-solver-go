package solver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestFindShortestPathNP_Heap_Square(t *testing.T) {
	assert := assert.New(t)

	testEntry := &testEntry{
		Expected: 16.0,
		Points: [][]float64{
			{0.0, 0.0},
			{4.0, 4.0},
			{4.0, 0.0},
			{0.0, 4.0},
		},
	}

	vertices := make([]model.CircuitVertex, len(testEntry.Points))
	for i, points := range testEntry.Points {
		vertices[i] = &model2d.Vertex2D{
			X: points[0],
			Y: points[1],
		}
	}

	shortest, actual := FindShortestPathNPHeap(vertices)
	shortestString, err := json.Marshal(shortest)
	assert.Nil(err)
	assert.InDelta(testEntry.Expected, actual, 0.00001, fmt.Sprintf("pathLength=%f shortestPath=%s", actual, shortestString))
}

func TestFindShortestPathNP_Heap_Concave(t *testing.T) {
	assert := assert.New(t)

	testEntry := &testEntry{
		Expected: 16.472135955,
		Points: [][]float64{
			{0.0, 0.0},
			{4.0, 4.0},
			{4.0, 0.0},
			{0.0, 4.0},
			{2.0, 1.0},
		},
	}

	vertices := make([]model.CircuitVertex, len(testEntry.Points))
	for i, points := range testEntry.Points {
		vertices[i] = &model2d.Vertex2D{
			X: points[0],
			Y: points[1],
		}
	}

	shortest, actual := FindShortestPathNPHeap(vertices)
	shortestString, err := json.Marshal(shortest)
	assert.Nil(err)
	assert.InDelta(testEntry.Expected, actual, 0.00001, fmt.Sprintf("pathLength=%f shortestPath=%s", actual, shortestString))
}

func TestFindShortestPathNP_Heap_DataFromOldProject(t *testing.T) {
	assert := assert.New(t)

	dataBytes, err := ioutil.ReadFile("../test-data/sample-polygons.json")
	assert.Nil(err)

	data := &testEntries{}
	err = json.Unmarshal(dataBytes, data)
	assert.Nil(err, "")

	for testIndex, testEntry := range data.Arrays[0:10] {
		vertices := make([]model.CircuitVertex, len(testEntry.Points))
		for i, points := range testEntry.Points {
			vertices[i] = &model2d.Vertex2D{
				X: points[0],
				Y: points[1],
			}
		}
		shortest, actual := FindShortestPathNPHeap(vertices)
		shortestBytes, err := json.Marshal(shortest)
		assert.Nil(err)
		shortestString := string(shortestBytes)
		assert.InDelta(testEntry.Expected, actual, model.Threshold, fmt.Sprintf("test=%d pathLength=%f shortestPath=%s", testIndex, actual, shortestString))
	}
}
