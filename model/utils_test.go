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

func TestDeleteIndexInt(t *testing.T) {
	assert := assert.New(t)

	ints := []int{
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
	}

	ints = model.DeleteIndexInt(ints, 0)
	assert.Equal([]int{
		2,
		3,
		4,
		5,
		6,
		7,
		8,
	}, ints)

	ints = model.DeleteIndexInt(ints, 99)
	assert.Equal([]int{
		2,
		3,
		4,
		5,
		6,
		7,
	}, ints)

	ints = model.DeleteIndexInt(ints, -5)
	assert.Equal([]int{
		3,
		4,
		5,
		6,
		7,
	}, ints)

	ints = model.DeleteIndexInt(ints, 3)
	assert.Equal([]int{
		3,
		4,
		5,
		7,
	}, ints)

	ints = model.DeleteIndexInt(ints, 0)
	ints = model.DeleteIndexInt(ints, 0)
	ints = model.DeleteIndexInt(ints, 0)
	assert.Len(ints, 1)
	ints = model.DeleteIndexInt(ints, 0)
	assert.Len(ints, 0)
	ints = model.DeleteIndexInt(ints, 0)
	assert.Len(ints, 0)
}
