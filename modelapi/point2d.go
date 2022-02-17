package modelapi

import (
	"github.com/heustis/lee-tsp-go/model"
	"github.com/heustis/lee-tsp-go/model2d"
)

// Point2D is the API representation a 2-dimensional point.
// It uses pointers to floats rather than floats, so that the fields can be correctly validated (0.0 is valid, but nil is not).
type Point2D struct {
	X *float64 `json:"x" validate:"required"`
	Y *float64 `json:"y" validate:"required"`
}

// To2D converts an API request into an array of 2-dimensional vertices.
func (api *TspRequest) To2D() []model.CircuitVertex {
	vertices := make([]model.CircuitVertex, len(api.Points2D))

	for i, v := range api.Points2D {
		vertices[i] = &model2d.Vertex2D{
			X: *v.X,
			Y: *v.Y,
		}
	}

	return model2d.DeduplicateVertices(vertices)
}

// ToApiFrom2D converts an array of 2-dimensional vertices into an API response.
func ToApiFrom2D(vertices []model.CircuitVertex) *TspRequest {
	api := &TspRequest{
		Points2D: make([]*Point2D, len(vertices)),
	}

	for i, v := range vertices {
		v2d := v.(*model2d.Vertex2D)
		api.Points2D[i] = &Point2D{
			X: &v2d.X,
			Y: &v2d.Y,
		}
	}
	return api
}
