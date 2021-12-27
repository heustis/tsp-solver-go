package model2d

import (
	"fmt"
	"math"

	"github.com/fealos/lee-tsp-go/model"
)

// Edge2D represents the line segment between two points
type Edge2D struct {
	Start  *Vertex2D `json:"start"`
	End    *Vertex2D `json:"end"`
	vector *Vertex2D
	length float64
}

// DistanceIncrease returns the difference in length between the edge
// and the two edges formed by inserting the vertex between the edge's start and end.
// For example, if start->end has a length of 5, start->vertex has a length of 3,
//  and vertex->end has a length of 6, this will return 4 (i.e. 6 + 3 - 5)
func (e *Edge2D) DistanceIncrease(vertex model.CircuitVertex) float64 {
	return e.Start.DistanceTo(vertex) + e.End.DistanceTo(vertex) - e.length
}

func (v *Edge2D) Equals(other interface{}) bool {
	// Compare pointers first, for performance, but then check start and end points, in case the same edge is created multiple times.
	if v == other {
		return true
	} else if otherVertex, okay := other.(*Edge2D); okay {
		return v.Start.Equals(otherVertex.Start) && v.End.Equals(otherVertex.End)
	} else {
		return false
	}
}

// GetStart returns the start vertex of the edge
func (e *Edge2D) GetStart() model.CircuitVertex {
	return e.Start
}

// GetEnd returns the end vertex of the edge
func (e *Edge2D) GetEnd() model.CircuitVertex {
	return e.End
}

// GetLength returns the length of the edge
func (e *Edge2D) GetLength() float64 {
	return e.length
}

// GetVector returns the normalized (length=1.0) vector from the edge's start to the edges end
func (e *Edge2D) GetVector() *Vertex2D {
	return e.vector
}

// Intersects checks if the two edges go through at least one identical point.
func (e *Edge2D) Intersects(other model.CircuitEdge) bool {
	otherEdge2D := other.(*Edge2D)
	// See http://paulbourke.net/geometry/pointlineplane/
	eDeltaX := e.End.X - e.Start.X
	eDeltaY := e.End.Y - e.Start.Y
	otherDeltaX := otherEdge2D.End.X - otherEdge2D.Start.X
	otherDeltaY := otherEdge2D.End.Y - otherEdge2D.Start.Y
	denominator := otherDeltaY*eDeltaX - otherDeltaX*eDeltaY

	startToStartDeltaX := e.Start.X - otherEdge2D.Start.X
	startToStartDeltaY := e.Start.Y - otherEdge2D.Start.Y

	if math.Abs(denominator) < model.Threshold {
		// Edges are parallel, check if the edges are on the same line, then return true if lines overlap.

		// To do this, use the same math as the denominator, with the edges being this edge and an edge from this start to the other start.
		// If this is also parallel, the line segments are on the same infinite line.
		startToStartParallel := (otherEdge2D.Start.Y-e.Start.Y)*eDeltaX - (otherEdge2D.Start.X-e.Start.X)*eDeltaY

		return math.Abs(startToStartParallel) < model.Threshold && (model.IsBetween(e.Start.X, otherEdge2D.Start.X, otherEdge2D.End.X) ||
			model.IsBetween(e.End.X, otherEdge2D.Start.X, otherEdge2D.End.X) ||
			model.IsBetween(otherEdge2D.Start.X, e.Start.X, e.End.X) ||
			model.IsBetween(otherEdge2D.End.X, e.Start.X, e.End.X))
	}

	// Determine the percentage of this edge's length from the start to the intersecting point.
	// Needs to be between 0 and 1, negative indicates the intersection is before the start, greater than 1 indicates that it is after the end.
	numeratorE := otherDeltaX*startToStartDeltaY - otherDeltaY*startToStartDeltaX
	if intersectDistPercentE := numeratorE / denominator; intersectDistPercentE < -model.Threshold || intersectDistPercentE > 1.0+model.Threshold {
		return false
	}

	// Check that the intersecting point exists within the other edge's length based on its percentage.
	numeratorOther := eDeltaX*startToStartDeltaY - eDeltaY*startToStartDeltaX
	intersectDistPercentOther := numeratorOther / denominator
	return intersectDistPercentOther >= -model.Threshold && intersectDistPercentOther < 1.0+model.Threshold
}

// Merge creates a new edge starting from this edge's start vertex and ending at the supplied edge's end vertex.
func (e *Edge2D) Merge(other model.CircuitEdge) model.CircuitEdge {
	return NewEdge2D(e.Start, other.GetEnd().(*Vertex2D))
}

// Split creates two new edges "start-to-vertex" and "vertex-to-end" based on this edge and the supplied vertex.
func (e *Edge2D) Split(vertex model.CircuitVertex) (model.CircuitEdge, model.CircuitEdge) {
	return NewEdge2D(e.Start, vertex.(*Vertex2D)), NewEdge2D(vertex.(*Vertex2D), e.End)
}

// ToString prints the edge as a string.
func (e *Edge2D) ToString() string {
	return fmt.Sprintf(`{"start":%s,"end":%s}`, e.Start.ToString(), e.End.ToString())
}

// NewEdge2D creates a edge from the starting Vertex2D to the ending Vertex2D
func NewEdge2D(start *Vertex2D, end *Vertex2D) *Edge2D {
	length := start.DistanceTo(end)
	vector := NewVertex2D((end.X-start.X)/length, (end.Y-start.Y)/length)
	return &Edge2D{
		Start:  start,
		End:    end,
		vector: vector,
		length: length,
	}
}

var _ model.CircuitEdge = (*Edge2D)(nil)
