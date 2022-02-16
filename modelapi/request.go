package modelapi

type TspRequest struct {
	Algorithms []*Algorithm `json:"algorithms,omitempty" validate:"dive,required"`
	// "excluded_with" is not fully documented in the validator docs, but it is in their source code,
	// see https://github.com/go-playground/validator/blob/v10.10.0/baked_in.go#L78
	Points2D    []*Point2D    `json:"points2d,omitempty" validate:"required_without_all=Points3D PointsGraph,excluded_with=Points3D PointsGraph,isdefault|min=3,dive,required"`
	Points3D    []*Point3D    `json:"points3d,omitempty" validate:"required_without_all=Points2D PointsGraph,excluded_with=Points2D PointsGraph,isdefault|min=3,dive,required"`
	PointsGraph []*PointGraph `json:"pointsGraph,omitempty" validate:"required_without_all=Points2D Points3D,excluded_with=Points2D Points3D,isdefault|min=3,dive,required"`
}
