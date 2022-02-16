package modelapi_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/model3d"
	"github.com/fealos/lee-tsp-go/modelapi"
	"github.com/stretchr/testify/assert"
)

func TestTo3D_And_ToApiFrom3D(t *testing.T) {
	assert := assert.New(t)

	initVertices := model3d.GenerateVertices(30)
	initLen := len(initVertices)

	apiInit := modelapi.ToApiFrom3D(initVertices)
	assert.Len(apiInit.Points3D, initLen)
	assert.Len(apiInit.Points2D, 0)
	assert.Len(apiInit.PointsGraph, 0)
	assert.Len(apiInit.Algorithms, 0)

	for i := 0; i < initLen; i++ {
		assert.Equal(initVertices[i].(*model3d.Vertex3D).X, *apiInit.Points3D[i].X)
		assert.Equal(initVertices[i].(*model3d.Vertex3D).Y, *apiInit.Points3D[i].Y)
		assert.Equal(initVertices[i].(*model3d.Vertex3D).Z, *apiInit.Points3D[i].Z)
	}

	modelVertices2 := apiInit.To3D()
	assert.Equal(initVertices, modelVertices2)

	api2 := modelapi.ToApiFrom3D(modelVertices2)
	assert.Equal(apiInit, api2)
}
