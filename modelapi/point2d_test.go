package modelapi_test

import (
	"testing"

	"github.com/heustis/lee-tsp-go/model2d"
	"github.com/heustis/lee-tsp-go/modelapi"
	"github.com/stretchr/testify/assert"
)

func TestTo2D_And_ToApiFrom2D(t *testing.T) {
	assert := assert.New(t)

	initVertices := model2d.DeduplicateVertices(model2d.GenerateVertices(30))
	initLen := len(initVertices)

	apiInit := modelapi.ToApiFrom2D(initVertices)
	assert.Len(apiInit.Points2D, initLen)
	assert.Len(apiInit.Points3D, 0)
	assert.Len(apiInit.PointsGraph, 0)
	assert.Len(apiInit.Algorithms, 0)

	for i := 0; i < initLen; i++ {
		assert.Equal(initVertices[i].(*model2d.Vertex2D).X, *apiInit.Points2D[i].X)
		assert.Equal(initVertices[i].(*model2d.Vertex2D).Y, *apiInit.Points2D[i].Y)
	}

	modelVertices2 := apiInit.To2D()
	assert.Equal(initVertices, modelVertices2)

	api2 := modelapi.ToApiFrom2D(modelVertices2)
	assert.Equal(apiInit, api2)
}
