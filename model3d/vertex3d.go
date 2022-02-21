package model3d

import (
	"fmt"
	"math"
	"strconv"

	"github.com/heustis/tsp-solver-go/model"
)

// Vertex3D represents a 3-dimensional point
type Vertex3D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Add returns a new Vertex3D that is the sum of this vertex and the supplied vertex.
func (v *Vertex3D) Add(other *Vertex3D) *Vertex3D {
	return &Vertex3D{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}

// DistanceTo returns the distance between the two vertices.
func (v *Vertex3D) DistanceTo(other model.CircuitVertex) float64 {
	o := other.(*Vertex3D)
	return math.Sqrt(v.DistanceToSquared(o))
}

// DistanceToEdge returns the shortest distance between this point and the supplied edge.
func (v *Vertex3D) DistanceToEdge(e *Edge3D) float64 {
	// return v.DistanceTo(v.ProjectToEdge(e))
	x := (v.X - e.Start.X)
	y := (v.Y - e.Start.Y)
	z := (v.Z - e.Start.Z)
	vector := e.GetVector()
	dot := x*vector.X + y*vector.Y + z*vector.Z

	xDiff := (e.Start.X + (vector.X * dot)) - v.X
	yDiff := (e.Start.Y + (vector.Y * dot)) - v.Y
	zDiff := (e.Start.Z + (vector.Z * dot)) - v.Z

	return math.Sqrt(xDiff*xDiff + yDiff*yDiff + zDiff*zDiff)
}

// DistanceToSquared returns the square of the distance between the two vertices.
func (v *Vertex3D) DistanceToSquared(other *Vertex3D) float64 {
	xDiff := other.X - v.X
	yDiff := other.Y - v.Y
	zDiff := other.Z - v.Z
	return xDiff*xDiff + yDiff*yDiff + zDiff*zDiff
}

// DotProduct returns the dot product of this vertex and the supplied vertex.
func (v *Vertex3D) DotProduct(other *Vertex3D) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// EdgeTo returns a new Edge3D from this vertex (the start) to the supplied vertex (the end).
func (v *Vertex3D) EdgeTo(end model.CircuitVertex) model.CircuitEdge {
	return NewEdge3D(v, end.(*Vertex3D))
}

// Equals checks if the two vertices are equal.
// It compares pointers first, for performance, but then checks X, Y, and Z, in case the same vertex is created multiple times.
func (v *Vertex3D) Equals(other interface{}) bool {
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

// Multiply returns a new Vertex3D with its coordinates multiplied by the suppied value.
func (v *Vertex3D) Multiply(scalar float64) *Vertex3D {
	return &Vertex3D{X: v.X * scalar, Y: v.Y * scalar, Z: v.Z * scalar}
}

// ProjectToEdge returns a new Vertex3D which is the closest point on the supplied edge to this vertex.
func (v *Vertex3D) ProjectToEdge(e *Edge3D) *Vertex3D {
	// e.Start.add(e.vector.Multiply(v.subtract(e.Start).dotProduct(e.vector)))
	x := (v.X - e.Start.X)
	y := (v.Y - e.Start.Y)
	z := (v.Z - e.Start.Z)
	vector := e.GetVector()
	dot := x*vector.X + y*vector.Y + z*vector.Z

	return &Vertex3D{
		X: e.Start.X + (vector.X * dot),
		Y: e.Start.Y + (vector.Y * dot),
		Z: e.Start.Z + (vector.Z * dot),
	}
}

// String prints the vertex as a string.
func (v *Vertex3D) String() string {
	xString := strconv.FormatFloat(v.X, 'f', -1, 64)
	yString := strconv.FormatFloat(v.Y, 'f', -1, 64)
	zString := strconv.FormatFloat(v.Z, 'f', -1, 64)

	return fmt.Sprintf(`{"x":%s,"y":%s,"z":%s}`, xString, yString, zString)
}

// Subtract returns a new Vertex3D that is the difference between this vertex and the supplied vertex.
func (v *Vertex3D) Subtract(other *Vertex3D) *Vertex3D {
	return &Vertex3D{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z}
}

// NewVertex3D creates a new 3-dimensional vertex.
func NewVertex3D(x float64, y float64, z float64) *Vertex3D {
	return &Vertex3D{X: x, Y: y, Z: z}
}

var _ model.CircuitVertex = (*Vertex3D)(nil)
