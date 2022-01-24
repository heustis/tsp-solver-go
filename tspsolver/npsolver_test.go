package tspsolver_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/fealos/lee-tsp-go/tspmodel"
	"github.com/fealos/lee-tsp-go/tspmodel2d"
	"github.com/fealos/lee-tsp-go/tspsolver"
	"github.com/stretchr/testify/assert"
)

type testEntry struct {
	Expected float64     `json:"expected"`
	Points   [][]float64 `json:"points"`
}

type testEntries struct {
	Arrays []testEntry `json:"arrays"`
}

func TestFindShortestPathNPWithChecks_Square(t *testing.T) {
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

	vertices := make([]tspmodel.CircuitVertex, len(testEntry.Points))
	for i, points := range testEntry.Points {
		vertices[i] = &tspmodel2d.Vertex2D{
			X: points[0],
			Y: points[1],
		}
	}

	shortest, actual := tspsolver.FindShortestPathNPWithChecks(vertices)
	shortestString, err := json.Marshal(shortest)
	assert.Nil(err)
	assert.InDelta(testEntry.Expected, actual, 0.00001, fmt.Sprintf("pathLength=%f shortestPath=%s", actual, shortestString))

	shortestNoChecks, actualNoChecks := tspsolver.FindShortestPathNPNoChecks(vertices)
	shortestStringNoChecks, err := json.Marshal(shortestNoChecks)
	assert.Nil(err)
	assert.InDelta(testEntry.Expected, actualNoChecks, 0.00001, fmt.Sprintf("pathLength=%f shortestPath=%s", actualNoChecks, shortestStringNoChecks))
}

func TestFindShortestPathNPWithChecks_Concave(t *testing.T) {
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

	vertices := make([]tspmodel.CircuitVertex, len(testEntry.Points))
	for i, points := range testEntry.Points {
		vertices[i] = &tspmodel2d.Vertex2D{
			X: points[0],
			Y: points[1],
		}
	}

	shortest, actual := tspsolver.FindShortestPathNPWithChecks(vertices)
	shortestString, err := json.Marshal(shortest)
	assert.Nil(err)
	assert.InDelta(testEntry.Expected, actual, 0.00001, fmt.Sprintf("pathLength=%f shortestPath=%s", actual, shortestString))

	shortestNoChecks, actualNoChecks := tspsolver.FindShortestPathNPNoChecks(vertices)
	shortestStringNoChecks, err := json.Marshal(shortestNoChecks)
	assert.Nil(err)
	assert.InDelta(testEntry.Expected, actualNoChecks, 0.00001, fmt.Sprintf("pathLength=%f shortestPath=%s", actualNoChecks, shortestStringNoChecks))
}

func TestFindShortestPathNPWithChecks_DataFromOldProject(t *testing.T) {
	assert := assert.New(t)

	dataBytes, err := ioutil.ReadFile("../test-data/sample-polygons.json")
	assert.Nil(err)

	data := &testEntries{}
	err = json.Unmarshal(dataBytes, data)
	assert.Nil(err, "")

	for testIndex, testEntry := range data.Arrays[0:10] {
		vertices := make([]tspmodel.CircuitVertex, len(testEntry.Points))
		for i, points := range testEntry.Points {
			vertices[i] = &tspmodel2d.Vertex2D{
				X: points[0],
				Y: points[1],
			}
		}
		shortest, actual := tspsolver.FindShortestPathNPWithChecks(vertices)
		shortestBytes, err := json.Marshal(shortest)
		assert.Nil(err)
		shortestString := string(shortestBytes)
		assert.InDelta(testEntry.Expected, actual, tspmodel.Threshold, fmt.Sprintf("test=%d pathLength=%f shortestPath=%s", testIndex, actual, shortestString))

		shortestNoChecks, actualNoChecks := tspsolver.FindShortestPathNPNoChecks(vertices)
		shortestStringNoChecks, err := json.Marshal(shortestNoChecks)
		assert.Nil(err)
		assert.InDelta(testEntry.Expected, actualNoChecks, 0.00001, fmt.Sprintf("pathLength=%f shortestPath=%s", actualNoChecks, shortestStringNoChecks))
	}
}
