package tspmodel2d

import (
	"fmt"
	"math"
	"strconv"

	"github.com/fealos/lee-tsp-go/tspmodel"
)

// Vertex2D represents a 2-dimensional point
type Vertex2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Add returns a new Vertex2D that is the sum of this vertex and the supplied vertex.
func (v *Vertex2D) Add(other *Vertex2D) *Vertex2D {
	return &Vertex2D{X: v.X + other.X, Y: v.Y + other.Y}
}

// DistanceTo returns the distance between this vertex and the supplied vertex.
func (v *Vertex2D) DistanceTo(other tspmodel.CircuitVertex) float64 {
	o := other.(*Vertex2D)
	return math.Sqrt(v.DistanceToSquared(o))
}

// DistanceToEdge returns the shortest distance between this point and the supplied edge.
func (v *Vertex2D) DistanceToEdge(e *Edge2D) float64 {
	return v.DistanceTo(v.ProjectToEdge(e))
}

// DistanceToSquared returns the square of the distance between the two vertices.
func (v *Vertex2D) DistanceToSquared(other *Vertex2D) float64 {
	xDiff := other.X - v.X
	yDiff := other.Y - v.Y
	return xDiff*xDiff + yDiff*yDiff
}

// DotProduct returns the dot product of this vertex and the supplied vertex.
func (v *Vertex2D) DotProduct(other *Vertex2D) float64 {
	return v.X*other.X + v.Y*other.Y
}

// EdgeTo returns a new Edge2D from this vertex (the start) to the supplied vertex (the end).
func (v *Vertex2D) EdgeTo(end tspmodel.CircuitVertex) tspmodel.CircuitEdge {
	return NewEdge2D(v, end.(*Vertex2D))
}

// Equals checks if the two vertices are equal.
// It compares pointers first, for performance, but then checks X and Y, in case the same vertex is created multiple times.
func (v *Vertex2D) Equals(other interface{}) bool {
	if v == other {
		return true
	} else if otherVertex, okay := other.(*Vertex2D); okay {
		return math.Abs(v.X-otherVertex.X) < tspmodel.Threshold && math.Abs(v.Y-otherVertex.Y) < tspmodel.Threshold
	} else {
		return false
	}
}

// IsLeftOf returns true if this vertex is to the left of the supplied edge.
func (v *Vertex2D) IsLeftOf(e *Edge2D) bool {
	// This math is the same as the following, but avoids creation of additional Vertex2Ds:
	// return e.vector.LeftPerpendicular().DotProduct(v.Subtract(e.Start)) > tspmodel.Threshold
	vector := e.GetVector()
	dot := -vector.Y*(v.X-e.Start.X) + vector.X*(v.Y-e.Start.Y)
	return dot > tspmodel.Threshold
}

// IsRightOf returns true if this vertex is to the right of the supplied edge.
func (v *Vertex2D) IsRightOf(e *Edge2D) bool {
	// This math is the same as the following, but avoids creation of additional Vertex2Ds:
	// return e.vector.RightPerpendicular().DotProduct(v.Subtract(e.Start)) > tspmodel.Threshold
	vector := e.GetVector()
	dot := vector.Y*(v.X-e.Start.X) - vector.X*(v.Y-e.Start.Y)
	return dot > tspmodel.Threshold
}

// LeftPerpendicular returns a new Vertex2D that is 90 degrees perpendicular of this vertex, to the left (counter-clockwise or anti-clockwise).
func (v *Vertex2D) LeftPerpendicular() *Vertex2D {
	return &Vertex2D{X: -v.Y, Y: v.X}
}

// Multiply returns a new Vertex2D with its coordinates multiplied by the suppied value.
func (v *Vertex2D) Multiply(scalar float64) *Vertex2D {
	return &Vertex2D{X: v.X * scalar, Y: v.Y * scalar}
}

// ProjectToEdge returns a new Vertex2D which is the closest point on the supplied edge to this vertex.
func (v *Vertex2D) ProjectToEdge(e *Edge2D) *Vertex2D {
	// This math is the same as the following, but avoids creation of additional Vertex2Ds:
	// return e.Start.Add(e.vector.Multiply(v.Subtract(e.Start).DotProduct(e.vector)))
	x := (v.X - e.Start.X)
	y := (v.Y - e.Start.Y)
	vector := e.GetVector()
	dot := x*vector.X + y*vector.Y

	return &Vertex2D{
		X: e.Start.X + (vector.X * dot),
		Y: e.Start.Y + (vector.Y * dot),
	}
}

// RightPerpendicular returns a new Vertex2D that is 90 degrees perpendicular of this vertex, to the right (clockwise).
func (v *Vertex2D) RightPerpendicular() *Vertex2D {
	return &Vertex2D{X: v.Y, Y: -v.X}
}

// Subtract returns a new Vertex2D that is the difference between this vertex and the supplied vertex.
func (v *Vertex2D) Subtract(other *Vertex2D) *Vertex2D {
	return &Vertex2D{X: v.X - other.X, Y: v.Y - other.Y}
}

// String prints the vertex as a string.
func (v *Vertex2D) String() string {
	xString := strconv.FormatFloat(v.X, 'f', -1, 64)
	yString := strconv.FormatFloat(v.Y, 'f', -1, 64)
	return fmt.Sprintf(`{"x":%s,"y":%s}`, xString, yString)
}

// NewVertex2D creates a new Vertex2D.
func NewVertex2D(x float64, y float64) *Vertex2D {
	return &Vertex2D{X: x, Y: y}
}

var _ tspmodel.CircuitVertex = (*Vertex2D)(nil)
