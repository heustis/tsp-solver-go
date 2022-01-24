package tspmodel_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/tspmodel"
	"github.com/fealos/lee-tsp-go/tspmodel2d"
	"github.com/fealos/lee-tsp-go/tspmodel3d"
	"github.com/stretchr/testify/assert"
)

func TestGetDistanceToEdgeForHeap(t *testing.T) {
	assert := assert.New(t)

	assert.Panics(func() { tspmodel.GetDistanceToEdgeForHeap(3.21) })
	assert.Panics(func() { tspmodel.GetDistanceToEdgeForHeap(tspmodel3d.NewVertex3D(1, 2, 3)) })
	assert.Equal(0.0, tspmodel.GetDistanceToEdgeForHeap(&tspmodel.DistanceToEdge{Distance: 0.0}))
	assert.Equal(15.67, tspmodel.GetDistanceToEdgeForHeap(&tspmodel.DistanceToEdge{Distance: 15.67}))
	assert.Equal(-1.11, tspmodel.GetDistanceToEdgeForHeap(&tspmodel.DistanceToEdge{Distance: -1.11}))
}

func TestHasVertex(t *testing.T) {
	assert := assert.New(t)

	v := tspmodel2d.NewVertex2D(123.45, 678.9)

	d := &tspmodel.DistanceToEdge{
		Vertex: v,
	}

	assert.True(d.HasVertex(&tspmodel.DistanceToEdge{
		Vertex: v,
	}))

	assert.False(d.HasVertex(&tspmodel.DistanceToEdge{
		Vertex: tspmodel2d.NewVertex2D(1.23, 4.56),
	}))

	assert.False(d.HasVertex(v))
	assert.False(d.HasVertex(nil))
}

func TestString_DistanceToEdge(t *testing.T) {
	assert := assert.New(t)

	d := &tspmodel.DistanceToEdge{
		Vertex:   tspmodel2d.NewVertex2D(123.45, 678.9),
		Edge:     tspmodel2d.NewEdge2D(tspmodel2d.NewVertex2D(5.15, 0.13), tspmodel2d.NewVertex2D(1000.3, 1100.25)),
		Distance: 567.89000001,
	}

	assert.Equal(`{"vertex":{"x":123.45,"y":678.9},"edge":{"start":{"x":5.15,"y":0.13},"end":{"x":1000.3,"y":1100.25}},"distance":567.89000001}`, d.String())
}
