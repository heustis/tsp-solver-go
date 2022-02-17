package model3d_test

import (
	"testing"

	"github.com/heustis/lee-tsp-go/model3d"
	"github.com/stretchr/testify/assert"
)

func TestGenerateVertices(t *testing.T) {
	assert := assert.New(t)

	for i := 3; i < 15; i++ {
		assert.Len(model3d.GenerateVertices(i), i)
	}
}
