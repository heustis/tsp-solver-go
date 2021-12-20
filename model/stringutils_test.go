package model_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestToString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(`1.234`, model.ToString(1.234))
	assert.Equal(`"test"`, model.ToString("test"))

	d := &model.DistanceToEdge{
		Vertex:   model2d.NewVertex2D(123.45, 678.9),
		Edge:     model2d.NewEdge2D(model2d.NewVertex2D(5.15, 0.13), model2d.NewVertex2D(1000.3, 1100.25)),
		Distance: 567.89,
	}
	assert.Equal(`{"vertex":{"x":123.45,"y":678.9},"edge":{"start":{"x":5.15,"y":0.13},"end":{"x":1000.3,"y":1100.25}},"distance":567.89}`, model.ToString(d))

	type testStruct struct {
		Foo   string `json:"bar"`
		Other int    `json:"other"`
	}

	assert.Equal(`{"bar":"test data","other":567}`, model.ToString(&testStruct{Foo: "test data", Other: 567}))
}
