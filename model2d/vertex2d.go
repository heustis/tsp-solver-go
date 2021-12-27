package model2d

import (
	"fmt"
	"math"

	"github.com/fealos/lee-tsp-go/model"
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

// DistanceTo returns the distance between the two vertices
func (v *Vertex2D) DistanceTo(other model.CircuitVertex) float64 {
	o := other.(*Vertex2D)
	return math.Sqrt(v.DistanceToSquared(o))
}

// DistanceToEdge returns the shortest distance between this point and the supplied edge.
func (v *Vertex2D) DistanceToEdge(e *Edge2D) float64 {
	return v.DistanceTo(v.ProjectToEdge(e))
}

// DistanceToSquared returns the square of the distance between the two vertices
func (v *Vertex2D) DistanceToSquared(other *Vertex2D) float64 {
	xDiff := other.X - v.X
	yDiff := other.Y - v.Y
	return xDiff*xDiff + yDiff*yDiff
}

// DotProduct returns the dot product of this vertex and the supplied vertex.
func (v *Vertex2D) DotProduct(other *Vertex2D) float64 {
	return v.X*other.X + v.Y*other.Y
}

// EdgeTo returns a new Edge2D from this vertex to the supplied vertex.
func (v *Vertex2D) EdgeTo(end model.CircuitVertex) model.CircuitEdge {
	return NewEdge2D(v, end.(*Vertex2D))
}

func (v *Vertex2D) Equals(other interface{}) bool {
	// Compare pointers first, for performance, but then check X and Y, in case the same vertex is created multiple times.
	if v == other {
		return true
	} else if otherVertex, okay := other.(*Vertex2D); okay {
		return math.Abs(v.X-otherVertex.X) < model.Threshold && math.Abs(v.Y-otherVertex.Y) < model.Threshold
	} else {
		return false
	}
}

// IsLeftOf returns true if this vertex is to the left of the supplied edge.
func (v *Vertex2D) IsLeftOf(e *Edge2D) bool {
	return e.vector.LeftPerpendicular().DotProduct(v.Subtract(e.Start)) > model.Threshold
}

// IsRightOf returns true if this vertex is to the right of the supplied edge.
func (v *Vertex2D) IsRightOf(e *Edge2D) bool {
	return e.vector.RightPerpendicular().DotProduct(v.Subtract(e.Start)) > model.Threshold
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
	return e.Start.Add(e.vector.Multiply(v.Subtract(e.Start).DotProduct(e.vector)))
}

// RightPerpendicular returns a new Vertex2D that is 90 degrees perpendicular of this vertex, to the right (clockwise).
func (v *Vertex2D) RightPerpendicular() *Vertex2D {
	return &Vertex2D{X: v.Y, Y: -v.X}
}

// Subtract returns a new Vertex2D that is the difference between this vertex and the supplied vertex.
func (v *Vertex2D) Subtract(other *Vertex2D) *Vertex2D {
	return &Vertex2D{X: v.X - other.X, Y: v.Y - other.Y}
}

// ToString prints the vertex as a string.
func (v *Vertex2D) ToString() string {
	return fmt.Sprintf(`{"x":%v,"y":%v}`, v.X, v.Y)
}

// NewVertex2D creates a new Vertex2D
func NewVertex2D(x float64, y float64) *Vertex2D {
	return &Vertex2D{X: x, Y: y}
}

var _ model.CircuitVertex = (*Vertex2D)(nil)
