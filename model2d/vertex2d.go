package model2d

import (
	"container/list"
	"fmt"
	"math"

	"github.com/fealos/lee-tsp-go/model"
)

// Vertex2D represents a 2-dimensional point
type Vertex2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// DistanceTo returns the distance between the two vertices
func (v *Vertex2D) DistanceTo(other model.CircuitVertex) float64 {
	o := other.(*Vertex2D)
	return math.Sqrt(v.DistanceToSquared(o))
}

// DistanceToSquared returns the square of the distance between the two vertices
func (v *Vertex2D) DistanceToSquared(other *Vertex2D) float64 {
	xDiff := other.X - v.X
	yDiff := other.Y - v.Y
	return xDiff*xDiff + yDiff*yDiff
}

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

// FindClosestEdge finds, and returns, the edge that is the closest to the vertex.
func (v *Vertex2D) FindClosestEdge(currentCircuit []model.CircuitEdge) model.CircuitEdge {
	var closest model.CircuitEdge = nil
	closestDistanceIncrease := math.MaxFloat64
	for _, candidate := range currentCircuit {
		// Ignore edges already containing the vertex.
		if candidate.GetEnd() == v || candidate.GetStart() == v {
			continue
		}
		candidateDistanceIncrease := candidate.DistanceIncrease(v)
		if candidateDistanceIncrease < closestDistanceIncrease {
			closest = candidate
			closestDistanceIncrease = candidateDistanceIncrease
		}
	}
	return closest
}

func (v *Vertex2D) FindClosestEdgeList(currentCircuit *list.List) model.CircuitEdge {
	var closest model.CircuitEdge = nil
	closestDistanceIncrease := math.MaxFloat64
	for i, link := 0, currentCircuit.Front(); i < currentCircuit.Len(); i, link = i+1, link.Next() {
		candidate := link.Value.(*Edge2D)
		// Ignore edges already containing the vertex.
		if candidate.GetEnd() == v || candidate.GetStart() == v {
			continue
		}
		candidateDistanceIncrease := candidate.DistanceIncrease(v)
		if candidateDistanceIncrease < closestDistanceIncrease {
			closest = candidate
			closestDistanceIncrease = candidateDistanceIncrease
		}
	}
	return closest
}

// IsEdgeCloser checks if the supplied edge is closer than the current closest edge.
func (v *Vertex2D) IsEdgeCloser(candidateEdge model.CircuitEdge, currentEdge model.CircuitEdge) bool {
	return candidateEdge.DistanceIncrease(v) < currentEdge.DistanceIncrease(v)
}

// ToString prints the vertex as a string.
func (v *Vertex2D) ToString() string {
	return fmt.Sprintf(`{"x":%v,"y":%v}`, v.X, v.Y)
}

func (v *Vertex2D) add(other *Vertex2D) *Vertex2D {
	return &Vertex2D{X: v.X + other.X, Y: v.Y + other.Y}
}

func (v *Vertex2D) distanceToEdge(e *Edge2D) float64 {
	return v.DistanceTo(v.projectToEdge(e))
}

func (v *Vertex2D) dotProduct(other *Vertex2D) float64 {
	return v.X*other.X + v.Y*other.Y
}

func (v *Vertex2D) isLeftOf(e *Edge2D) bool {
	return e.vector.leftPerpendicular().dotProduct(v.subtract(e.Start)) > model.Threshold
}

func (v *Vertex2D) isRightOf(e *Edge2D) bool {
	return e.vector.rightPerpendicular().dotProduct(v.subtract(e.Start)) > model.Threshold
}

func (v *Vertex2D) leftPerpendicular() *Vertex2D {
	return &Vertex2D{X: -v.Y, Y: v.X}
}

func (v *Vertex2D) multiply(scalar float64) *Vertex2D {
	return &Vertex2D{X: v.X * scalar, Y: v.Y * scalar}
}

func (v *Vertex2D) projectToEdge(e *Edge2D) *Vertex2D {
	return e.Start.add(e.vector.multiply(v.subtract(e.Start).dotProduct(e.vector)))
}

func (v *Vertex2D) rightPerpendicular() *Vertex2D {
	return &Vertex2D{X: v.Y, Y: -v.X}
}

func (v *Vertex2D) subtract(other *Vertex2D) *Vertex2D {
	return &Vertex2D{X: v.X - other.X, Y: v.Y - other.Y}
}

// NewVertex2D creates a vertex
func NewVertex2D(x float64, y float64) *Vertex2D {
	return &Vertex2D{X: x, Y: y}
}

var _ model.CircuitVertex = (*Vertex2D)(nil)
