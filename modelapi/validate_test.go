package modelapi_test

import (
	"math"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/heustis/tsp-solver-go/modelapi"
	"github.com/stretchr/testify/assert"
)

func TestValidateTspRequest(t *testing.T) {
	assert := assert.New(t)
	validate := validator.New()
	// At least one array must not be present
	assert.EqualError(validate.Struct(modelapi.TspRequest{}),
		"Key: 'TspRequest.Points2D' Error:Field validation for 'Points2D' failed on the 'required_without_all' tag\n"+
			"Key: 'TspRequest.Points3D' Error:Field validation for 'Points3D' failed on the 'required_without_all' tag\n"+
			"Key: 'TspRequest.PointsGraph' Error:Field validation for 'PointsGraph' failed on the 'required_without_all' tag")

	// Each array must not be empty
	assert.EqualError(validate.Struct(modelapi.TspRequest{Points2D: []*modelapi.Point2D{}}), `Key: 'TspRequest.Points2D' Error:Field validation for 'Points2D' failed on the 'isdefault|min=3' tag`)
	assert.EqualError(validate.Struct(modelapi.TspRequest{Points3D: []*modelapi.Point3D{}}), `Key: 'TspRequest.Points3D' Error:Field validation for 'Points3D' failed on the 'isdefault|min=3' tag`)
	assert.EqualError(validate.Struct(modelapi.TspRequest{PointsGraph: []*modelapi.PointGraph{}}), `Key: 'TspRequest.PointsGraph' Error:Field validation for 'PointsGraph' failed on the 'isdefault|min=3' tag`)

	// Each array must have at least 3 entries.
	assert.EqualError(validate.Struct(modelapi.TspRequest{Points2D: []*modelapi.Point2D{{X: float64Pointer(1), Y: float64Pointer(2)}, {X: float64Pointer(3), Y: float64Pointer(4)}}}), `Key: 'TspRequest.Points2D' Error:Field validation for 'Points2D' failed on the 'isdefault|min=3' tag`)
	assert.EqualError(validate.Struct(modelapi.TspRequest{Points3D: []*modelapi.Point3D{{X: float64Pointer(1), Y: float64Pointer(2), Z: float64Pointer(3)}, {X: float64Pointer(3), Y: float64Pointer(4), Z: float64Pointer(5)}}}), `Key: 'TspRequest.Points3D' Error:Field validation for 'Points3D' failed on the 'isdefault|min=3' tag`)
	assert.EqualError(validate.Struct(modelapi.TspRequest{PointsGraph: []*modelapi.PointGraph{{Id: "a", Neighbors: []modelapi.PointGraphNeighbor{{Id: "b", Distance: *float64Pointer(1)}}}, {Id: "a", Neighbors: []modelapi.PointGraphNeighbor{{Id: "b", Distance: *float64Pointer(1)}}}}}), `Key: 'TspRequest.PointsGraph' Error:Field validation for 'PointsGraph' failed on the 'isdefault|min=3' tag`)

	// Each entry in the array must be non-nil.
	assert.EqualError(validate.Struct(modelapi.TspRequest{Points2D: []*modelapi.Point2D{{X: float64Pointer(1), Y: float64Pointer(2)}, {X: float64Pointer(3), Y: float64Pointer(4)}, nil}}), `Key: 'TspRequest.Points2D[2]' Error:Field validation for 'Points2D[2]' failed on the 'required' tag`)
	assert.EqualError(validate.Struct(modelapi.TspRequest{Points3D: []*modelapi.Point3D{{X: float64Pointer(1), Y: float64Pointer(2), Z: float64Pointer(3)}, nil, {X: float64Pointer(3), Y: float64Pointer(4), Z: float64Pointer(5)}}}), `Key: 'TspRequest.Points3D[1]' Error:Field validation for 'Points3D[1]' failed on the 'required' tag`)
	assert.EqualError(validate.Struct(modelapi.TspRequest{PointsGraph: []*modelapi.PointGraph{nil, {Id: "a", Neighbors: []modelapi.PointGraphNeighbor{{Id: "b", Distance: *float64Pointer(1)}}}, {Id: "b", Neighbors: []modelapi.PointGraphNeighbor{{Id: "c", Distance: *float64Pointer(2)}}}}}), `Key: 'TspRequest.PointsGraph[0]' Error:Field validation for 'PointsGraph[0]' failed on the 'required' tag`)

	// Valid data
	assert.Nil(validate.Struct(modelapi.TspRequest{Points2D: []*modelapi.Point2D{{X: float64Pointer(1), Y: float64Pointer(2)}, {X: float64Pointer(3), Y: float64Pointer(4)}, {X: float64Pointer(5), Y: float64Pointer(6)}}}))
	assert.Nil(validate.Struct(modelapi.TspRequest{Points3D: []*modelapi.Point3D{{X: float64Pointer(1), Y: float64Pointer(2), Z: float64Pointer(3)}, {X: float64Pointer(3), Y: float64Pointer(4), Z: float64Pointer(5)}, {X: float64Pointer(6), Y: float64Pointer(7), Z: float64Pointer(8)}}}))
	assert.Nil(validate.Struct(modelapi.TspRequest{PointsGraph: []*modelapi.PointGraph{
		{Id: "a", Neighbors: []modelapi.PointGraphNeighbor{{Id: "b", Distance: *float64Pointer(1)}}},
		{Id: "b", Neighbors: []modelapi.PointGraphNeighbor{{Id: "c", Distance: *float64Pointer(2)}}},
		{Id: "c", Neighbors: []modelapi.PointGraphNeighbor{{Id: "a", Distance: *float64Pointer(2)}, {Id: "b", Distance: *float64Pointer(1)}}},
	}}))

	// Data in the arrays must be valid.
	assert.EqualError(validate.Struct(modelapi.TspRequest{Points2D: []*modelapi.Point2D{{X: float64Pointer(1), Y: float64Pointer(2)}, {Y: float64Pointer(4)}, {X: float64Pointer(5), Y: float64Pointer(6)}}}), `Key: 'TspRequest.Points2D[1].X' Error:Field validation for 'X' failed on the 'required' tag`)
	assert.EqualError(validate.Struct(modelapi.TspRequest{Points3D: []*modelapi.Point3D{{Y: float64Pointer(2), Z: float64Pointer(3)}, {X: float64Pointer(3), Y: float64Pointer(4), Z: float64Pointer(5)}, {X: float64Pointer(6), Y: float64Pointer(7), Z: float64Pointer(8)}}}), `Key: 'TspRequest.Points3D[0].X' Error:Field validation for 'X' failed on the 'required' tag`)
	assert.EqualError(validate.Struct(modelapi.TspRequest{PointsGraph: []*modelapi.PointGraph{
		{Id: "a", Neighbors: []modelapi.PointGraphNeighbor{{Id: "b", Distance: *float64Pointer(1)}}},
		{Id: "b", Neighbors: []modelapi.PointGraphNeighbor{{Id: "c", Distance: *float64Pointer(2)}}},
		{Id: "c", Neighbors: []modelapi.PointGraphNeighbor{}},
	}}), `Key: 'TspRequest.PointsGraph[2].Neighbors' Error:Field validation for 'Neighbors' failed on the 'min' tag`)

	// Object in algorithms cannot be nil
	assert.EqualError(validate.Struct(modelapi.TspRequest{Algorithms: []*modelapi.Algorithm{nil, {AlgorithmType: modelapi.ALG_CLOSEST_CLONE}}, Points2D: []*modelapi.Point2D{{X: float64Pointer(1), Y: float64Pointer(2)}, {X: float64Pointer(3), Y: float64Pointer(4)}, {X: float64Pointer(5), Y: float64Pointer(6)}}}),
		`Key: 'TspRequest.Algorithms[0]' Error:Field validation for 'Algorithms[0]' failed on the 'required' tag`)
	// Data in algorithms must be valid.
	assert.EqualError(validate.Struct(modelapi.TspRequest{Algorithms: []*modelapi.Algorithm{{AlgorithmType: "Other"}}, Points2D: []*modelapi.Point2D{{X: float64Pointer(1), Y: float64Pointer(2)}, {X: float64Pointer(3), Y: float64Pointer(4)}, {X: float64Pointer(5), Y: float64Pointer(6)}}}),
		`Key: 'TspRequest.Algorithms[0].AlgorithmType' Error:Field validation for 'AlgorithmType' failed on the 'oneof' tag`)

	// Only one of the ararys should be populated (cannot mix and match point types).
	assert.EqualError(validate.Struct(modelapi.TspRequest{
		Points2D: []*modelapi.Point2D{{X: float64Pointer(1), Y: float64Pointer(2)}, {X: float64Pointer(3), Y: float64Pointer(4)}, {X: float64Pointer(5), Y: float64Pointer(6)}},
		Points3D: []*modelapi.Point3D{{X: float64Pointer(1), Y: float64Pointer(2), Z: float64Pointer(3)}, {X: float64Pointer(3), Y: float64Pointer(4), Z: float64Pointer(5)}, {X: float64Pointer(6), Y: float64Pointer(7), Z: float64Pointer(8)}},
	}),
		"Key: 'TspRequest.Points2D' Error:Field validation for 'Points2D' failed on the 'excluded_with' tag\n"+
			"Key: 'TspRequest.Points3D' Error:Field validation for 'Points3D' failed on the 'excluded_with' tag")

	assert.EqualError(validate.Struct(modelapi.TspRequest{
		Points2D: []*modelapi.Point2D{{X: float64Pointer(1), Y: float64Pointer(2)}, {X: float64Pointer(3), Y: float64Pointer(4)}, {X: float64Pointer(5), Y: float64Pointer(6)}},
		PointsGraph: []*modelapi.PointGraph{
			{Id: "a", Neighbors: []modelapi.PointGraphNeighbor{{Id: "b", Distance: *float64Pointer(1)}}},
			{Id: "b", Neighbors: []modelapi.PointGraphNeighbor{{Id: "c", Distance: *float64Pointer(2)}}},
			{Id: "c", Neighbors: []modelapi.PointGraphNeighbor{{Id: "a", Distance: *float64Pointer(2)}, {Id: "b", Distance: *float64Pointer(1)}}},
		},
	}),
		"Key: 'TspRequest.Points2D' Error:Field validation for 'Points2D' failed on the 'excluded_with' tag\n"+
			"Key: 'TspRequest.PointsGraph' Error:Field validation for 'PointsGraph' failed on the 'excluded_with' tag")

	assert.EqualError(validate.Struct(modelapi.TspRequest{
		Points3D: []*modelapi.Point3D{{X: float64Pointer(1), Y: float64Pointer(2), Z: float64Pointer(3)}, {X: float64Pointer(3), Y: float64Pointer(4), Z: float64Pointer(5)}, {X: float64Pointer(6), Y: float64Pointer(7), Z: float64Pointer(8)}},
		PointsGraph: []*modelapi.PointGraph{
			{Id: "a", Neighbors: []modelapi.PointGraphNeighbor{{Id: "b", Distance: *float64Pointer(1)}}},
			{Id: "b", Neighbors: []modelapi.PointGraphNeighbor{{Id: "c", Distance: *float64Pointer(2)}}},
			{Id: "c", Neighbors: []modelapi.PointGraphNeighbor{{Id: "a", Distance: *float64Pointer(2)}, {Id: "b", Distance: *float64Pointer(1)}}},
		},
	}),
		"Key: 'TspRequest.Points3D' Error:Field validation for 'Points3D' failed on the 'excluded_with' tag\n"+
			"Key: 'TspRequest.PointsGraph' Error:Field validation for 'PointsGraph' failed on the 'excluded_with' tag")

	assert.EqualError(validate.Struct(modelapi.TspRequest{
		Points2D: []*modelapi.Point2D{{X: float64Pointer(1), Y: float64Pointer(2)}, {X: float64Pointer(3), Y: float64Pointer(4)}, {X: float64Pointer(5), Y: float64Pointer(6)}},

		Points3D: []*modelapi.Point3D{{X: float64Pointer(1), Y: float64Pointer(2), Z: float64Pointer(3)}, {X: float64Pointer(3), Y: float64Pointer(4), Z: float64Pointer(5)}, {X: float64Pointer(6), Y: float64Pointer(7), Z: float64Pointer(8)}},
		PointsGraph: []*modelapi.PointGraph{
			{Id: "a", Neighbors: []modelapi.PointGraphNeighbor{{Id: "b", Distance: *float64Pointer(1)}}},
			{Id: "b", Neighbors: []modelapi.PointGraphNeighbor{{Id: "c", Distance: *float64Pointer(2)}}},
			{Id: "c", Neighbors: []modelapi.PointGraphNeighbor{{Id: "a", Distance: *float64Pointer(2)}, {Id: "b", Distance: *float64Pointer(1)}}},
		},
	}),
		"Key: 'TspRequest.Points2D' Error:Field validation for 'Points2D' failed on the 'excluded_with' tag\n"+
			"Key: 'TspRequest.Points3D' Error:Field validation for 'Points3D' failed on the 'excluded_with' tag\n"+
			"Key: 'TspRequest.PointsGraph' Error:Field validation for 'PointsGraph' failed on the 'excluded_with' tag")
}

func TestValidatePoint2D(t *testing.T) {
	assert := assert.New(t)
	validate := validator.New()
	assert.EqualError(validate.Struct(modelapi.Point2D{}), "Key: 'Point2D.X' Error:Field validation for 'X' failed on the 'required' tag\nKey: 'Point2D.Y' Error:Field validation for 'Y' failed on the 'required' tag")
	assert.EqualError(validate.Struct(modelapi.Point2D{X: float64Pointer(1.23)}), `Key: 'Point2D.Y' Error:Field validation for 'Y' failed on the 'required' tag`)
	assert.EqualError(validate.Struct(modelapi.Point2D{Y: float64Pointer(2.34)}), `Key: 'Point2D.X' Error:Field validation for 'X' failed on the 'required' tag`)
	assert.Nil(validate.Struct(modelapi.Point2D{X: float64Pointer(1.23), Y: float64Pointer(2.34)}))
	assert.Nil(validate.Struct(modelapi.Point2D{X: float64Pointer(0), Y: float64Pointer(0.00)}))
}

func TestValidatePoint3D(t *testing.T) {
	assert := assert.New(t)
	validate := validator.New()
	assert.EqualError(validate.Struct(modelapi.Point3D{}), "Key: 'Point3D.X' Error:Field validation for 'X' failed on the 'required' tag\nKey: 'Point3D.Y' Error:Field validation for 'Y' failed on the 'required' tag\nKey: 'Point3D.Z' Error:Field validation for 'Z' failed on the 'required' tag")
	assert.EqualError(validate.Struct(modelapi.Point3D{X: float64Pointer(1.23)}), "Key: 'Point3D.Y' Error:Field validation for 'Y' failed on the 'required' tag\nKey: 'Point3D.Z' Error:Field validation for 'Z' failed on the 'required' tag")
	assert.EqualError(validate.Struct(modelapi.Point3D{Y: float64Pointer(2.34)}), "Key: 'Point3D.X' Error:Field validation for 'X' failed on the 'required' tag\nKey: 'Point3D.Z' Error:Field validation for 'Z' failed on the 'required' tag")
	assert.EqualError(validate.Struct(modelapi.Point3D{Z: float64Pointer(2.34)}), "Key: 'Point3D.X' Error:Field validation for 'X' failed on the 'required' tag\nKey: 'Point3D.Y' Error:Field validation for 'Y' failed on the 'required' tag")
	assert.EqualError(validate.Struct(modelapi.Point3D{X: float64Pointer(1.23), Y: float64Pointer(2.34)}), `Key: 'Point3D.Z' Error:Field validation for 'Z' failed on the 'required' tag`)
	assert.EqualError(validate.Struct(modelapi.Point3D{X: float64Pointer(1.23), Z: float64Pointer(2.34)}), `Key: 'Point3D.Y' Error:Field validation for 'Y' failed on the 'required' tag`)
	assert.EqualError(validate.Struct(modelapi.Point3D{Y: float64Pointer(1.23), Z: float64Pointer(2.34)}), `Key: 'Point3D.X' Error:Field validation for 'X' failed on the 'required' tag`)
	assert.Nil(validate.Struct(modelapi.Point3D{X: float64Pointer(1), Y: float64Pointer(2.2), Z: float64Pointer(3.45)}))
	assert.Nil(validate.Struct(modelapi.Point3D{X: float64Pointer(0), Y: float64Pointer(0.00), Z: float64Pointer(0.0)}))
}

func TestValidatePointGraph(t *testing.T) {
	assert := assert.New(t)
	validate := validator.New()
	assert.EqualError(validate.Struct(modelapi.PointGraph{}), "Key: 'PointGraph.Id' Error:Field validation for 'Id' failed on the 'required' tag\nKey: 'PointGraph.Neighbors' Error:Field validation for 'Neighbors' failed on the 'required' tag")
	assert.EqualError(validate.Struct(modelapi.PointGraph{Id: "a"}), `Key: 'PointGraph.Neighbors' Error:Field validation for 'Neighbors' failed on the 'required' tag`)
	assert.EqualError(validate.Struct(modelapi.PointGraph{Id: "b", Neighbors: []modelapi.PointGraphNeighbor{}}), `Key: 'PointGraph.Neighbors' Error:Field validation for 'Neighbors' failed on the 'min' tag`)
	assert.EqualError(validate.Struct(modelapi.PointGraph{Id: "c", Neighbors: []modelapi.PointGraphNeighbor{{Id: "b"}}}), `Key: 'PointGraph.Neighbors[0].Distance' Error:Field validation for 'Distance' failed on the 'required' tag`)
	assert.EqualError(validate.Struct(modelapi.PointGraph{Id: "d", Neighbors: []modelapi.PointGraphNeighbor{{Id: ""}}}), "Key: 'PointGraph.Neighbors[0].Id' Error:Field validation for 'Id' failed on the 'required' tag\nKey: 'PointGraph.Neighbors[0].Distance' Error:Field validation for 'Distance' failed on the 'required' tag")
	assert.EqualError(validate.Struct(modelapi.PointGraph{Id: "e", Neighbors: []modelapi.PointGraphNeighbor{{Id: "", Distance: 1.23}}}), `Key: 'PointGraph.Neighbors[0].Id' Error:Field validation for 'Id' failed on the 'required' tag`)
	assert.EqualError(validate.Struct(modelapi.PointGraph{Id: "f", Neighbors: []modelapi.PointGraphNeighbor{{Distance: 1.23}}}), `Key: 'PointGraph.Neighbors[0].Id' Error:Field validation for 'Id' failed on the 'required' tag`)
	assert.EqualError(validate.Struct(modelapi.PointGraph{Id: "g", Neighbors: []modelapi.PointGraphNeighbor{{Id: "b", Distance: -1.23}}}), `Key: 'PointGraph.Neighbors[0].Distance' Error:Field validation for 'Distance' failed on the 'min' tag`)
	assert.EqualError(validate.Struct(modelapi.PointGraph{Id: "h", Neighbors: []modelapi.PointGraphNeighbor{{Id: "b", Distance: -math.SmallestNonzeroFloat64}}}), `Key: 'PointGraph.Neighbors[0].Distance' Error:Field validation for 'Distance' failed on the 'min' tag`)
	assert.EqualError(validate.Struct(modelapi.PointGraph{Id: "i", Neighbors: []modelapi.PointGraphNeighbor{{Id: "b", Distance: 0.12}, {Id: "b", Distance: 1.23}}}), `Key: 'PointGraph.Neighbors' Error:Field validation for 'Neighbors' failed on the 'unique' tag`)
	assert.EqualError(validate.Struct(modelapi.PointGraph{Id: "j", Neighbors: []modelapi.PointGraphNeighbor{{Id: "b", Distance: 0.12}, {Id: "c", Distance: -1.23}}}), `Key: 'PointGraph.Neighbors[1].Distance' Error:Field validation for 'Distance' failed on the 'min' tag`)
	// Validator/v10 does not support `unique` with nil values in the array, so the array does not use pointers.
	// Once that is supported
	// assert.EqualError(validate.Struct(modelapi.PointGraph{Id: "k", Neighbors: []modelapi.PointGraphNeighbor{{Id: "b", Distance: 0.12}, {Id: "c", Distance: 1.23}, nil}}), `Key: 'PointGraph.Neighbors' Error:Field validation for 'Neighbors' failed on the 'unique' tag`)
	assert.Nil(validate.Struct(modelapi.PointGraph{Id: "l", Neighbors: []modelapi.PointGraphNeighbor{{Id: "b", Distance: 0.12}, {Id: "c", Distance: 1.23}}}))
}

func float64Pointer(f float64) *float64 {
	return &f
}
