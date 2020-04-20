package tsp

import (
	"math"
)

// Vertex represents a 2-dimensional point
type Vertex struct {
	x float64
	y float64
}

// GetX returns the X coordinate of the vertex
func (v Vertex) GetX() float64 {
	return v.x
}

// GetY returns the Y coordinate of the vertex
func (v Vertex) GetY() float64 {
	return v.y
}

// DistanceTo returns the distance between the two vertices
func (v Vertex) DistanceTo(other *Vertex) float64 {
	return math.Sqrt(v.DistanceToSquared(other))
}

// DistanceToSquared returns the square of the distance between the two vertices
func (v Vertex) DistanceToSquared(other *Vertex) float64 {
	xDiff := other.x - v.x
	yDiff := other.y - v.y
	return xDiff*xDiff + yDiff*yDiff
}

// NewVertex creates a vertex
func NewVertex(x float64, y float64) *Vertex {
	return &Vertex{x: x, y: y}
}
