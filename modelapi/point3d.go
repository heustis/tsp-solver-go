package modelapi

import (
	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model3d"
)

// Point3D is the API representation a 3-dimensional point.
// It uses pointers to floats rather than floats, so that the fields can be correctly validated (0.0 is valid, but nil is not).
type Point3D struct {
	X *float64 `json:"x" validate:"required"`
	Y *float64 `json:"y" validate:"required"`
	Z *float64 `json:"z" validate:"required"`
}

// To3D converts an API request into an array of 3-dimensional vertices.
func (api *TspRequest) To3D() []model.CircuitVertex {
	vertices := make([]model.CircuitVertex, len(api.Points3D))

	for i, v := range api.Points3D {
		vertices[i] = &model3d.Vertex3D{
			X: *v.X,
			Y: *v.Y,
			Z: *v.Z,
		}
	}

	return model.DeduplicateVerticesNoSorting(vertices)
}

// ToApiFrom3D converts an array of 3-dimensional vertices into an API response.
func ToApiFrom3D(vertices []model.CircuitVertex) *TspRequest {
	api := &TspRequest{
		Points3D: make([]*Point3D, len(vertices)),
	}

	for i, v := range vertices {
		v3d := v.(*model3d.Vertex3D)
		api.Points3D[i] = &Point3D{
			X: &v3d.X,
			Y: &v3d.Y,
			Z: &v3d.Z,
		}
	}
	return api
}
