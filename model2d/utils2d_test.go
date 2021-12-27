package model2d_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestDeduplicateVertices2D(t *testing.T) {
	assert := assert.New(t)

	init := []model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(-15-model.Threshold/3.0, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(-15+model.Threshold/3.0, -15-model.Threshold/3.0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(7, 6),
		model2d.NewVertex2D(-7, 6),
	}

	actual := model2d.DeduplicateVertices(init)
	assert.ElementsMatch([]*model2d.Vertex2D{
		model2d.NewVertex2D(-15+model.Threshold/3.0, -15-model.Threshold/3.0),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(7, 6),
		model2d.NewVertex2D(15, -15),
	}, actual)
}

func TestGenerateVertices2D(t *testing.T) {
	assert := assert.New(t)

	for i := 3; i < 15; i++ {
		assert.Len(model2d.GenerateVertices2D(i), i)
	}
}
