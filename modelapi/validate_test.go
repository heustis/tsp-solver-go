package modelapi_test

import (
	"math"
	"testing"

	"github.com/fealos/lee-tsp-go/modelapi"
	"github.com/go-playground/validator/v10"
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
	assert.EqualError(validate.Struct(modelapi.TspRequest{Points2D: []*modelapi.Vertex2D{}}), `Key: 'TspRequest.Points2D' Error:Field validation for 'Points2D' failed on the 'isdefault|min=3' tag`)
	assert.EqualError(validate.Struct(modelapi.TspRequest{Points3D: []*modelapi.Vertex3D{}}), `Key: 'TspRequest.Points3D' Error:Field validation for 'Points3D' failed on the 'isdefault|min=3' tag`)
	assert.EqualError(validate.Struct(modelapi.TspRequest{PointsGraph: []*modelapi.VertexGraph{}}), `Key: 'TspRequest.PointsGraph' Error:Field validation for 'PointsGraph' failed on the 'isdefault|min=3' tag`)

	// Each array must have at least 3 entries.
	assert.EqualError(validate.Struct(modelapi.TspRequest{Points2D: []*modelapi.Vertex2D{{X: float64Pointer(1), Y: float64Pointer(2)}, {X: float64Pointer(3), Y: float64Pointer(4)}}}), `Key: 'TspRequest.Points2D' Error:Field validation for 'Points2D' failed on the 'isdefault|min=3' tag`)
	assert.EqualError(validate.Struct(modelapi.TspRequest{Points3D: []*modelapi.Vertex3D{{X: float64Pointer(1), Y: float64Pointer(2), Z: float64Pointer(3)}, {X: float64Pointer(3), Y: float64Pointer(4), Z: float64Pointer(5)}}}), `Key: 'TspRequest.Points3D' Error:Field validation for 'Points3D' failed on the 'isdefault|min=3' tag`)
	assert.EqualError(validate.Struct(modelapi.TspRequest{PointsGraph: []*modelapi.VertexGraph{{Id: "a", Neighbors: []modelapi.VertexNeighbor{{Id: "b", Distance: *float64Pointer(1)}}}, {Id: "a", Neighbors: []modelapi.VertexNeighbor{{Id: "b", Distance: *float64Pointer(1)}}}}}), `Key: 'TspRequest.PointsGraph' Error:Field validation for 'PointsGraph' failed on the 'isdefault|min=3' tag`)

	// Each entry in the array must be non-nil.
	assert.EqualError(validate.Struct(modelapi.TspRequest{Points2D: []*modelapi.Vertex2D{{X: float64Pointer(1), Y: float64Pointer(2)}, {X: float64Pointer(3), Y: float64Pointer(4)}, nil}}), `Key: 'TspRequest.Points2D[2]' Error:Field validation for 'Points2D[2]' failed on the 'required' tag`)
	assert.EqualError(validate.Struct(modelapi.TspRequest{Points3D: []*modelapi.Vertex3D{{X: float64Pointer(1), Y: float64Pointer(2), Z: float64Pointer(3)}, nil, {X: float64Pointer(3), Y: float64Pointer(4), Z: float64Pointer(5)}}}), `Key: 'TspRequest.Points3D[1]' Error:Field validation for 'Points3D[1]' failed on the 'required' tag`)
	assert.EqualError(validate.Struct(modelapi.TspRequest{PointsGraph: []*modelapi.VertexGraph{nil, {Id: "a", Neighbors: []modelapi.VertexNeighbor{{Id: "b", Distance: *float64Pointer(1)}}}, {Id: "b", Neighbors: []modelapi.VertexNeighbor{{Id: "c", Distance: *float64Pointer(2)}}}}}), `Key: 'TspRequest.PointsGraph[0]' Error:Field validation for 'PointsGraph[0]' failed on the 'required' tag`)

	// Valid data
	assert.Nil(validate.Struct(modelapi.TspRequest{Points2D: []*modelapi.Vertex2D{{X: float64Pointer(1), Y: float64Pointer(2)}, {X: float64Pointer(3), Y: float64Pointer(4)}, {X: float64Pointer(5), Y: float64Pointer(6)}}}))
	assert.Nil(validate.Struct(modelapi.TspRequest{Points3D: []*modelapi.Vertex3D{{X: float64Pointer(1), Y: float64Pointer(2), Z: float64Pointer(3)}, {X: float64Pointer(3), Y: float64Pointer(4), Z: float64Pointer(5)}, {X: float64Pointer(6), Y: float64Pointer(7), Z: float64Pointer(8)}}}))
	assert.Nil(validate.Struct(modelapi.TspRequest{PointsGraph: []*modelapi.VertexGraph{
		{Id: "a", Neighbors: []modelapi.VertexNeighbor{{Id: "b", Distance: *float64Pointer(1)}}},
		{Id: "b", Neighbors: []modelapi.VertexNeighbor{{Id: "c", Distance: *float64Pointer(2)}}},
		{Id: "c", Neighbors: []modelapi.VertexNeighbor{{Id: "a", Distance: *float64Pointer(2)}, {Id: "b", Distance: *float64Pointer(1)}}},
	}}))

	// Only one of the ararys should be populated (cannot mix and match point types).
	assert.EqualError(validate.Struct(modelapi.TspRequest{
		Points2D: []*modelapi.Vertex2D{{X: float64Pointer(1), Y: float64Pointer(2)}, {X: float64Pointer(3), Y: float64Pointer(4)}, {X: float64Pointer(5), Y: float64Pointer(6)}},
		Points3D: []*modelapi.Vertex3D{{X: float64Pointer(1), Y: float64Pointer(2), Z: float64Pointer(3)}, {X: float64Pointer(3), Y: float64Pointer(4), Z: float64Pointer(5)}, {X: float64Pointer(6), Y: float64Pointer(7), Z: float64Pointer(8)}},
	}),
		"Key: 'TspRequest.Points2D' Error:Field validation for 'Points2D' failed on the 'excluded_with' tag\n"+
			"Key: 'TspRequest.Points3D' Error:Field validation for 'Points3D' failed on the 'excluded_with' tag")

	assert.EqualError(validate.Struct(modelapi.TspRequest{
		Points2D: []*modelapi.Vertex2D{{X: float64Pointer(1), Y: float64Pointer(2)}, {X: float64Pointer(3), Y: float64Pointer(4)}, {X: float64Pointer(5), Y: float64Pointer(6)}},
		PointsGraph: []*modelapi.VertexGraph{
			{Id: "a", Neighbors: []modelapi.VertexNeighbor{{Id: "b", Distance: *float64Pointer(1)}}},
			{Id: "b", Neighbors: []modelapi.VertexNeighbor{{Id: "c", Distance: *float64Pointer(2)}}},
			{Id: "c", Neighbors: []modelapi.VertexNeighbor{{Id: "a", Distance: *float64Pointer(2)}, {Id: "b", Distance: *float64Pointer(1)}}},
		},
	}),
		"Key: 'TspRequest.Points2D' Error:Field validation for 'Points2D' failed on the 'excluded_with' tag\n"+
			"Key: 'TspRequest.PointsGraph' Error:Field validation for 'PointsGraph' failed on the 'excluded_with' tag")

	assert.EqualError(validate.Struct(modelapi.TspRequest{
		Points3D: []*modelapi.Vertex3D{{X: float64Pointer(1), Y: float64Pointer(2), Z: float64Pointer(3)}, {X: float64Pointer(3), Y: float64Pointer(4), Z: float64Pointer(5)}, {X: float64Pointer(6), Y: float64Pointer(7), Z: float64Pointer(8)}},
		PointsGraph: []*modelapi.VertexGraph{
			{Id: "a", Neighbors: []modelapi.VertexNeighbor{{Id: "b", Distance: *float64Pointer(1)}}},
			{Id: "b", Neighbors: []modelapi.VertexNeighbor{{Id: "c", Distance: *float64Pointer(2)}}},
			{Id: "c", Neighbors: []modelapi.VertexNeighbor{{Id: "a", Distance: *float64Pointer(2)}, {Id: "b", Distance: *float64Pointer(1)}}},
		},
	}),
		"Key: 'TspRequest.Points3D' Error:Field validation for 'Points3D' failed on the 'excluded_with' tag\n"+
			"Key: 'TspRequest.PointsGraph' Error:Field validation for 'PointsGraph' failed on the 'excluded_with' tag")

	assert.EqualError(validate.Struct(modelapi.TspRequest{
		Points2D: []*modelapi.Vertex2D{{X: float64Pointer(1), Y: float64Pointer(2)}, {X: float64Pointer(3), Y: float64Pointer(4)}, {X: float64Pointer(5), Y: float64Pointer(6)}},

		Points3D: []*modelapi.Vertex3D{{X: float64Pointer(1), Y: float64Pointer(2), Z: float64Pointer(3)}, {X: float64Pointer(3), Y: float64Pointer(4), Z: float64Pointer(5)}, {X: float64Pointer(6), Y: float64Pointer(7), Z: float64Pointer(8)}},
		PointsGraph: []*modelapi.VertexGraph{
			{Id: "a", Neighbors: []modelapi.VertexNeighbor{{Id: "b", Distance: *float64Pointer(1)}}},
			{Id: "b", Neighbors: []modelapi.VertexNeighbor{{Id: "c", Distance: *float64Pointer(2)}}},
			{Id: "c", Neighbors: []modelapi.VertexNeighbor{{Id: "a", Distance: *float64Pointer(2)}, {Id: "b", Distance: *float64Pointer(1)}}},
		},
	}),
		"Key: 'TspRequest.Points2D' Error:Field validation for 'Points2D' failed on the 'excluded_with' tag\n"+
			"Key: 'TspRequest.Points3D' Error:Field validation for 'Points3D' failed on the 'excluded_with' tag\n"+
			"Key: 'TspRequest.PointsGraph' Error:Field validation for 'PointsGraph' failed on the 'excluded_with' tag")
}

func TestValidateVertex2d(t *testing.T) {
	assert := assert.New(t)
	validate := validator.New()
	assert.Error(validate.Struct(modelapi.Vertex2D{}))
	assert.Error(validate.Struct(modelapi.Vertex2D{X: float64Pointer(1.23)}))
	assert.Error(validate.Struct(modelapi.Vertex2D{Y: float64Pointer(2.34)}))
	assert.Nil(validate.Struct(modelapi.Vertex2D{X: float64Pointer(1.23), Y: float64Pointer(2.34)}))
	assert.Nil(validate.Struct(modelapi.Vertex2D{X: float64Pointer(0), Y: float64Pointer(0.00)}))
}

func TestValidateVertex3d(t *testing.T) {
	assert := assert.New(t)
	validate := validator.New()
	assert.Error(validate.Struct(modelapi.Vertex3D{}))
	assert.Error(validate.Struct(modelapi.Vertex3D{X: float64Pointer(1.23)}))
	assert.Error(validate.Struct(modelapi.Vertex3D{Y: float64Pointer(2.34)}))
	assert.Error(validate.Struct(modelapi.Vertex3D{Z: float64Pointer(2.34)}))
	assert.Error(validate.Struct(modelapi.Vertex3D{X: float64Pointer(1.23), Y: float64Pointer(2.34)}))
	assert.Error(validate.Struct(modelapi.Vertex3D{X: float64Pointer(1.23), Z: float64Pointer(2.34)}))
	assert.Error(validate.Struct(modelapi.Vertex3D{Y: float64Pointer(1.23), Z: float64Pointer(2.34)}))
	assert.Nil(validate.Struct(modelapi.Vertex3D{X: float64Pointer(1), Y: float64Pointer(2.2), Z: float64Pointer(3.45)}))
	assert.Nil(validate.Struct(modelapi.Vertex3D{X: float64Pointer(0), Y: float64Pointer(0.00), Z: float64Pointer(0.0)}))
}

func TestValidateVertexGraph(t *testing.T) {
	assert := assert.New(t)
	validate := validator.New()
	assert.Error(validate.Struct(modelapi.VertexGraph{}))
	assert.Error(validate.Struct(modelapi.VertexGraph{Id: "a"}))
	assert.Error(validate.Struct(modelapi.VertexGraph{Id: "b", Neighbors: []modelapi.VertexNeighbor{}}))
	assert.Error(validate.Struct(modelapi.VertexGraph{Id: "c", Neighbors: []modelapi.VertexNeighbor{{Id: "b"}}}))
	assert.Error(validate.Struct(modelapi.VertexGraph{Id: "d", Neighbors: []modelapi.VertexNeighbor{{Id: ""}}}))
	assert.Error(validate.Struct(modelapi.VertexGraph{Id: "e", Neighbors: []modelapi.VertexNeighbor{{Id: "", Distance: 1.23}}}))
	assert.Error(validate.Struct(modelapi.VertexGraph{Id: "f", Neighbors: []modelapi.VertexNeighbor{{Distance: 1.23}}}))
	assert.Error(validate.Struct(modelapi.VertexGraph{Id: "g", Neighbors: []modelapi.VertexNeighbor{{Id: "b", Distance: -1.23}}}))
	assert.Error(validate.Struct(modelapi.VertexGraph{Id: "h", Neighbors: []modelapi.VertexNeighbor{{Id: "b", Distance: -math.SmallestNonzeroFloat64}}}))
	assert.Error(validate.Struct(modelapi.VertexGraph{Id: "i", Neighbors: []modelapi.VertexNeighbor{{Id: "b", Distance: 0.12}, {Id: "b", Distance: 1.23}}}))
	assert.Error(validate.Struct(modelapi.VertexGraph{Id: "j", Neighbors: []modelapi.VertexNeighbor{{Id: "b", Distance: 0.12}, {Id: "c", Distance: -1.23}}}))
	assert.Error(validate.Struct(modelapi.VertexGraph{Id: "k", Neighbors: []modelapi.VertexNeighbor{{Id: "b", Distance: 0.12}, {Id: "c", Distance: -1.23}}}))
	// Validator/v10 does not support `unique` with nil values in the array, so the array does not use pointers.
	// Once that is supported
	// assert.Error(validate.Struct(modelapi.VertexGraph{Id: "k", Neighbors: []modelapi.VertexNeighbor{{Id: "b", Distance: 0.12}, {Id: "c", Distance: 1.23}, nil}}))
	assert.Nil(validate.Struct(modelapi.VertexGraph{Id: "l", Neighbors: []modelapi.VertexNeighbor{{Id: "b", Distance: 0.12}, {Id: "c", Distance: 1.23}}}))
}

func float64Pointer(f float64) *float64 {
	return &f
}
