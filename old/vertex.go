package old

import (
	"math"
)

// Vertex represents a 2-dimensional point
type Vertex struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// DistanceTo returns the distance between the two vertices
func (v Vertex) DistanceTo(other *Vertex) float64 {
	return math.Sqrt(v.DistanceToSquared(other))
}

// DistanceToSquared returns the square of the distance between the two vertices
func (v Vertex) DistanceToSquared(other *Vertex) float64 {
	xDiff := other.X - v.X
	yDiff := other.Y - v.Y
	return xDiff*xDiff + yDiff*yDiff
}

// NewVertex creates a vertex
func NewVertex(x float64, y float64) *Vertex {
	return &Vertex{X: x, Y: y}
}
