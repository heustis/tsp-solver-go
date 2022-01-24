package tspmodel3d_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/tspmodel3d"
	"github.com/stretchr/testify/assert"
)

func TestGenerateVertices(t *testing.T) {
	assert := assert.New(t)

	for i := 3; i < 15; i++ {
		assert.Len(tspmodel3d.GenerateVertices(i), i)
	}
}
