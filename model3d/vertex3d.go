package model3d

import (
	"fmt"
	"math"

	"github.com/fealos/lee-tsp-go/model"
)

// Vertex3D represents a 3-dimensional point
type Vertex3D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// DistanceTo returns the distance between the two vertices
func (v *Vertex3D) DistanceTo(other model.CircuitVertex) float64 {
	o := other.(*Vertex3D)
	return math.Sqrt(v.DistanceToSquared(o))
}

// DistanceToSquared returns the square of the distance between the two vertices
func (v *Vertex3D) DistanceToSquared(other *Vertex3D) float64 {
	xDiff := other.X - v.X
	yDiff := other.Y - v.Y
	zDiff := other.Z - v.Z
	return xDiff*xDiff + yDiff*yDiff + zDiff*zDiff
}

func (v *Vertex3D) Equals(other interface{}) bool {
	// Compare pointers first, for performance, but then check X and Y, in case the same vertex is created multiple times.
	if v == other {
		return true
	} else if other == (*Vertex3D)(nil) || other == nil {
		return v == (*Vertex3D)(nil)
	} else if otherVertex, okay := other.(*Vertex3D); okay && v != (*Vertex3D)(nil) {
		return math.Abs(v.X-otherVertex.X) < model.Threshold && math.Abs(v.Y-otherVertex.Y) < model.Threshold && math.Abs(v.Z-otherVertex.Z) < model.Threshold
	} else {
		return false
	}
}

// FindClosestEdge finds, and returns, the edge that is the closest to the vertex.
func (v *Vertex3D) FindClosestEdge(currentCircuit []model.CircuitEdge) model.CircuitEdge {
	var closest model.CircuitEdge = nil
	closestDistanceIncrease := math.MaxFloat64
	for _, candidate := range currentCircuit {
		candidateDistanceIncrease := candidate.DistanceIncrease(v)
		if candidateDistanceIncrease < closestDistanceIncrease {
			closest = candidate
			closestDistanceIncrease = candidateDistanceIncrease
		}
	}
	return closest
}

// IsEdgeCloser checks if the supplied edge is closer than the current closest edge.
func (v *Vertex3D) IsEdgeCloser(candidateEdge model.CircuitEdge, currentEdge model.CircuitEdge) bool {
	return candidateEdge.DistanceIncrease(v) < currentEdge.DistanceIncrease(v)
}

// ToString prints the vertex as a string.
func (v *Vertex3D) ToString() string {
	return fmt.Sprintf(`{"x":%v,"y":%v,"z":%v}`, v.X, v.Y, v.Z)
}

func (v *Vertex3D) add(other *Vertex3D) *Vertex3D {
	return &Vertex3D{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}

func (v *Vertex3D) distanceToEdge(e *Edge3D) float64 {
	return v.DistanceTo(v.projectToEdge(e))
}

func (v *Vertex3D) dotProduct(other *Vertex3D) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v *Vertex3D) multiply(scalar float64) *Vertex3D {
	return &Vertex3D{X: v.X * scalar, Y: v.Y * scalar, Z: v.Z * scalar}
}

func (v *Vertex3D) projectToEdge(e *Edge3D) *Vertex3D {
	return e.Start.add(e.vector.multiply(v.subtract(e.Start).dotProduct(e.vector)))
}

func (v *Vertex3D) subtract(other *Vertex3D) *Vertex3D {
	return &Vertex3D{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z}
}

// NewVertex3D creates a vertex
func NewVertex3D(x float64, y float64, z float64) *Vertex3D {
	return &Vertex3D{X: x, Y: y, Z: z}
}

var _ model.CircuitVertex = (*Vertex3D)(nil)
