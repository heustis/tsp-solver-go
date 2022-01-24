package tspmodel_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsBetween(t *testing.T) {
	assert := assert.New(t)

	assert.True(tspmodel.IsBetween(5, 0, 10))
	assert.True(tspmodel.IsBetween(5, -5, 10))
	assert.True(tspmodel.IsBetween(5, -5, 5))
	assert.True(tspmodel.IsBetween(5, 5, 5))
	assert.False(tspmodel.IsBetween(5, 5.1, 5.6))
	assert.False(tspmodel.IsBetween(5, -5, -5))
	assert.False(tspmodel.IsBetween(5, 0, -10))
	assert.False(tspmodel.IsBetween(5, -10, 0))
}
