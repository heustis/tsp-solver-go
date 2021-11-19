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

func TestFindShortestPathGreedy_DataFromOldProject(t *testing.T) {
	assert := assert.New(t)

	dataBytes, err := ioutil.ReadFile("../test-data/sample-polygons.json")
	assert.Nil(err)

	data := &testEntries{}
	err = json.Unmarshal(dataBytes, data)
	assert.Nil(err)

	for testIndex, testEntry := range data.Arrays[0:10] {
		// TODO 1: Investigate cause of discrepancy
		//   Error:      	Max difference between 2753.5027007922668 and 2851.703838591683 allowed is 1e-05, but difference was -98.2011377994163
		//   Messages:   	test=0 pathLength=2851.703839 shortestPath=[{"X":-217,"Y":-323},{"X":-2,"Y":-132},{"X":199,"Y":-334},{"X":236,"Y":123},{"X":344,"Y":287},{"X":349,"Y":312},{"X":28,"Y":169},{"X":-169,"Y":325},{"X":-182,"Y":492},{"X":-169,"Y":325},{"X":-248,"Y":329}]

		vertices := make([]model.CircuitVertex, len(testEntry.Points))
		for i, points := range testEntry.Points {
			vertices[i] = model2d.NewVertex2D(points[0], points[1])
		}
		circuit := model.NewCircuitGreedyWithUpdatesImpl(vertices, model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
		FindShortestPathGreedy(circuit)
		shortest := circuit.GetAttachedVertices()
		actual := circuit.GetLength()

		shortestString, err := json.Marshal(shortest)
		assert.Nil(err)
		assert.InDelta(testEntry.Expected, actual, model.Threshold, fmt.Sprintf("test=%d pathLength=%f shortestPath=%s", testIndex, actual, shortestString))
	}
}
