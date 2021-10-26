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

func TestFindShortestPathHeap_DataFromOldProject(t *testing.T) {
	assert := assert.New(t)

	dataBytes, err := ioutil.ReadFile("../test-data/sample-polygons.json")
	assert.Nil(err)

	data := &testEntries{}
	err = json.Unmarshal(dataBytes, data)
	assert.Nil(err)

	for testIndex, testEntry := range data.Arrays[0:10] {
		vertices := make([]*model2d.Vertex2D, len(testEntry.Points))
		for i, points := range testEntry.Points {
			vertices[i] = model2d.NewVertex2D(points[0], points[1])
		}
		circuit := &model2d.HeapableCircuit2D{
			Vertices: vertices,
		}
		result, _ := FindShortestPathHeap(circuit)
		shortest := result.GetAttachedVertices()
		actual := result.GetLength()

		shortestString, err := json.Marshal(shortest)
		assert.Nil(err)
		assert.InDelta(testEntry.Expected, actual, model.Threshold, fmt.Sprintf("test=%d pathLength=%f shortestPath=%s", testIndex, actual, shortestString))
	}
}

func TestFindShortestPathHeap_MinClones_DataFromOldProject(t *testing.T) {
	assert := assert.New(t)

	dataBytes, err := ioutil.ReadFile("../test-data/sample-polygons.json")
	assert.Nil(err)

	data := &testEntries{}
	err = json.Unmarshal(dataBytes, data)
	assert.Nil(err)

	for testIndex, testEntry := range data.Arrays[0:10] {
		vertices := make([]*model2d.Vertex2D, len(testEntry.Points))
		for i, points := range testEntry.Points {
			vertices[i] = model2d.NewVertex2D(points[0], points[1])
		}
		circuit := &model2d.HeapableCircuit2DMinClones{
			Vertices: vertices,
		}
		result, _ := FindShortestPathHeap(circuit)
		shortest := result.GetAttachedVertices()
		actual := result.GetLength()

		shortestString, err := json.Marshal(shortest)
		assert.Nil(err)
		assert.InDelta(testEntry.Expected, actual, model.Threshold, fmt.Sprintf("test=%d pathLength=%f shortestPath=%s", testIndex, actual, shortestString))
	}
}
