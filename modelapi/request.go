package modelapi

type TspRequest struct {
	Algorithms []string `json:"algorithms"`
	// "excluded_with" is not included in the validator docs, but it is in their source code,
	// see https://github.com/go-playground/validator/blob/v10.10.0/baked_in.go#L78
	Points2D    []*Vertex2D    `json:"points2d" validate:"required_without_all=Points3D PointsGraph,excluded_with=Points3D PointsGraph,isdefault|min=3,dive,required"`
	Points3D    []*Vertex3D    `json:"points3d" validate:"required_without_all=Points2D PointsGraph,excluded_with=Points2D PointsGraph,isdefault|min=3,dive,required"`
	PointsGraph []*VertexGraph `json:"pointsGraph" validate:"required_without_all=Points2D Points3D,excluded_with=Points2D Points3D,isdefault|min=3,dive,required"`
}
