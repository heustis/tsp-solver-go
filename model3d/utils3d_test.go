package model3d_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model3d"
	"github.com/stretchr/testify/assert"
)

func TestDeduplicateVertices3D(t *testing.T) {
	assert := assert.New(t)

	init := []model.CircuitVertex{
		model3d.NewVertex3D(-15, -15, -5.0),
		model3d.NewVertex3D(0, 0, model.Threshold/9.0),
		model3d.NewVertex3D(15, -15, -5.0),
		model3d.NewVertex3D(-15-model.Threshold/3.0, -15, -5),
		model3d.NewVertex3D(0, 0, 0.0),
		model3d.NewVertex3D(0, model.Threshold*2, 0.0),
		model3d.NewVertex3D(-15+model.Threshold/3.0, -15-model.Threshold/3.0, -5+model.Threshold/4),
		model3d.NewVertex3D(3, 0, 3),
		model3d.NewVertex3D(3, 13, 4),
		model3d.NewVertex3D(7, 6, 5),
		model3d.NewVertex3D(-7, 6, 5),
	}

	actual := model3d.DeduplicateVertices3D(init)
	assert.ElementsMatch([]*model3d.Vertex3D{
		model3d.NewVertex3D(-15, -15, -5),
		model3d.NewVertex3D(-7, 6, 5),
		model3d.NewVertex3D(0, 0, model.Threshold/9.0),
		model3d.NewVertex3D(0, model.Threshold*2, 0.0),
		model3d.NewVertex3D(3, 0, 3),
		model3d.NewVertex3D(3, 13, 4),
		model3d.NewVertex3D(7, 6, 5),
		model3d.NewVertex3D(15, -15, -5),
	}, actual)
}
