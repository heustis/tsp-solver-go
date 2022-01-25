package model_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/fealos/lee-tsp-go/model3d"
	"github.com/stretchr/testify/assert"
)

func TestGetDistanceToEdgeForHeap(t *testing.T) {
	assert := assert.New(t)

	assert.Panics(func() { model.GetDistanceToEdgeForHeap(3.21) })
	assert.Panics(func() { model.GetDistanceToEdgeForHeap(model3d.NewVertex3D(1, 2, 3)) })
	assert.Equal(0.0, model.GetDistanceToEdgeForHeap(&model.DistanceToEdge{Distance: 0.0}))
	assert.Equal(15.67, model.GetDistanceToEdgeForHeap(&model.DistanceToEdge{Distance: 15.67}))
	assert.Equal(-1.11, model.GetDistanceToEdgeForHeap(&model.DistanceToEdge{Distance: -1.11}))
}

func TestHasVertex(t *testing.T) {
	assert := assert.New(t)

	v := model2d.NewVertex2D(123.45, 678.9)

	d := &model.DistanceToEdge{
		Vertex: v,
	}

	assert.True(d.HasVertex(&model.DistanceToEdge{
		Vertex: v,
	}))

	assert.False(d.HasVertex(&model.DistanceToEdge{
		Vertex: model2d.NewVertex2D(1.23, 4.56),
	}))

	assert.False(d.HasVertex(v))
	assert.False(d.HasVertex(nil))
}

func TestString_DistanceToEdge(t *testing.T) {
	assert := assert.New(t)

	d := &model.DistanceToEdge{
		Vertex:   model2d.NewVertex2D(123.45, 678.9),
		Edge:     model2d.NewEdge2D(model2d.NewVertex2D(5.15, 0.13), model2d.NewVertex2D(1000.3, 1100.25)),
		Distance: 567.89000001,
	}

	assert.Equal(`{"vertex":{"x":123.45,"y":678.9},"edge":{"start":{"x":5.15,"y":0.13},"end":{"x":1000.3,"y":1100.25}},"distance":567.89000001}`, d.String())
}
