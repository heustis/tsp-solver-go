package model_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/stretchr/testify/assert"
)

func TestIsBetween(t *testing.T) {
	assert := assert.New(t)

	assert.True(model.IsBetween(5, 0, 10))
	assert.True(model.IsBetween(5, -5, 10))
	assert.True(model.IsBetween(5, -5, 5))
	assert.True(model.IsBetween(5, 5, 5))
	assert.False(model.IsBetween(5, 5.1, 5.6))
	assert.False(model.IsBetween(5, -5, -5))
	assert.False(model.IsBetween(5, 0, -10))
	assert.False(model.IsBetween(5, -10, 0))
}
