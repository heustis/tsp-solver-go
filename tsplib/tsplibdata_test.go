package tsplib_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/fealos/lee-tsp-go/solver"
	"github.com/fealos/lee-tsp-go/tsplib"
	"github.com/stretchr/testify/assert"
)

func TestNewData(t *testing.T) {
	assert := assert.New(t)

	data, err := tsplib.NewData("../test-data/tsplib/not_a_file.tsp")
	assert.NotNil(err)
	assert.Nil(data)

	data, err = tsplib.NewData("../test-data/tsplib/a280.tsp.malformed.x")
	assert.NotNil(err)
	assert.Nil(data)

	data, err = tsplib.NewData("../test-data/tsplib/a280.tsp.malformed.y")
	assert.NotNil(err)
	assert.Nil(data)

	data, err = tsplib.NewData("../test-data/tsplib/a280.tsp")
	assert.Nil(err)
	assert.NotNil(data)

	assert.Equal(`a280`, data.GetName())
	assert.Equal(`drilling problem (Ludwig)`, data.GetComment())
	assert.Equal(280, data.GetNumPoints())

	vertices := data.GetVertices()
	assert.Len(vertices, 280)

	assert.Equal(model2d.NewVertex2D(288, 149), vertices[0])
	assert.Equal(model2d.NewVertex2D(288, 129), vertices[1])
	assert.Equal(model2d.NewVertex2D(32, 129), vertices[55])
	assert.Equal(model2d.NewVertex2D(280, 133), vertices[279])

	bestRoute := data.GetBestRoute()
	assert.Len(bestRoute, 280)
	assert.Equal(model2d.NewVertex2D(288, 149), bestRoute[0])
	assert.Equal(model2d.NewVertex2D(288, 129), bestRoute[1])
	assert.Equal(model2d.NewVertex2D(288, 109), bestRoute[2])
	assert.Equal(model2d.NewVertex2D(270, 133), bestRoute[278])
	assert.Equal(model2d.NewVertex2D(280, 133), bestRoute[279])

	assert.InDelta(2586.7696475631606, data.GetBestRouteLength(), model.Threshold)
}

func TestSolveAndCompare(t *testing.T) {
	assert := assert.New(t)

	data, err := tsplib.NewData("../test-data/tsplib/a280.tsp")
	assert.Nil(err)
	assert.NotNil(data)

	err = data.SolveAndCompare(func(cv []model.CircuitVertex) model.Circuit {
		c := circuit.NewConvexConcave(data.GetVertices(), model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{}, false)
		solver.FindShortestPathGreedy(c)
		return c
	})

	assert.Nil(err)
}
