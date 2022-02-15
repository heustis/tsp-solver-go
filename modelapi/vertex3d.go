package modelapi

// Vertex3D is the API representation a 3-dimensional point.
// It uses pointers to floats rather than floats, so that the fields can be correctly validated (0.0 is valid, but nil is not).
type Vertex3D struct {
	X *float64 `json:"x" validate:"required"`
	Y *float64 `json:"y" validate:"required"`
	Z *float64 `json:"z" validate:"required"`
}
