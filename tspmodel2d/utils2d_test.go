package tspmodel2d_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/tspmodel"
	"github.com/fealos/lee-tsp-go/tspmodel2d"
	"github.com/stretchr/testify/assert"
)

func TestDeduplicateVertices2D(t *testing.T) {
	assert := assert.New(t)

	init := []tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(-15, -15),
		tspmodel2d.NewVertex2D(0, 0),
		tspmodel2d.NewVertex2D(15, -15),
		tspmodel2d.NewVertex2D(-15-tspmodel.Threshold/3.0, -15),
		tspmodel2d.NewVertex2D(0, 0),
		tspmodel2d.NewVertex2D(-15+tspmodel.Threshold/3.0, -15-tspmodel.Threshold/3.0),
		tspmodel2d.NewVertex2D(3, 0),
		tspmodel2d.NewVertex2D(3, 13),
		tspmodel2d.NewVertex2D(7, 6),
		tspmodel2d.NewVertex2D(-7, 6),
	}

	actual := tspmodel2d.DeduplicateVertices(init)
	assert.ElementsMatch([]*tspmodel2d.Vertex2D{
		tspmodel2d.NewVertex2D(-15+tspmodel.Threshold/3.0, -15-tspmodel.Threshold/3.0),
		tspmodel2d.NewVertex2D(-7, 6),
		tspmodel2d.NewVertex2D(0, 0),
		tspmodel2d.NewVertex2D(3, 0),
		tspmodel2d.NewVertex2D(3, 13),
		tspmodel2d.NewVertex2D(7, 6),
		tspmodel2d.NewVertex2D(15, -15),
	}, actual)
}

func TestGenerateVertices(t *testing.T) {
	assert := assert.New(t)

	for i := 3; i < 15; i++ {
		assert.Len(tspmodel2d.GenerateVertices(i), i)
	}
}
