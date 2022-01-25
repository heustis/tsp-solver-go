package model3d

import (
	"fmt"
	"math"

	"github.com/fealos/lee-tsp-go/model"
)

// Edge3D represents the line segment between two points
type Edge3D struct {
	Start  *Vertex3D `json:"start"`
	End    *Vertex3D `json:"end"`
	vector *Vertex3D
	length float64
}

// DistanceIncrease returns the difference in length between the edge
// and the two edges formed by inserting the vertex between the edge's start and end.
// For example, if start->end has a length of 5, start->vertex has a length of 3,
//  and vertex->end has a length of 6, this will return 4 (i.e. 6 + 3 - 5)
func (e *Edge3D) DistanceIncrease(vertex model.CircuitVertex) float64 {
	return e.Start.DistanceTo(vertex) + e.End.DistanceTo(vertex) - e.length
}

func (e *Edge3D) Equals(other interface{}) bool {
	// Compare pointers first, for performance, but then check start and end points, in case the same edge is created multiple times.
	if e == other {
		return true
	} else if other == (*Edge3D)(nil) || other == nil {
		return e == (*Edge3D)(nil)
	} else if otherVertex, okay := other.(*Edge3D); okay && e != (*Edge3D)(nil) {
		return e.Start.Equals(otherVertex.Start) && e.End.Equals(otherVertex.End)
	} else {
		return false
	}
}

// GetStart returns the start vertex of the edge
func (e *Edge3D) GetStart() model.CircuitVertex {
	return e.Start
}

// GetEnd returns the end vertex of the edge
func (e *Edge3D) GetEnd() model.CircuitVertex {
	return e.End
}

// GetLength returns the length of the edge
func (e *Edge3D) GetLength() float64 {
	return e.length
}

// GetVector returns the normalized (length=1.0) vector from the edge's start to the edges end
func (e *Edge3D) GetVector() *Vertex3D {
	if e.vector == nil {
		e.vector = NewVertex3D((e.End.X-e.Start.X)/e.length, (e.End.Y-e.Start.Y)/e.length, (e.End.Z-e.Start.Z)/e.length)
	}
	return e.vector
}

// Intersects checks if the two edges go through at least one identical point.
func (e *Edge3D) Intersects(other model.CircuitEdge) bool {
	otherEdge3D := other.(*Edge3D)
	// See http://paulbourke.net/geometry/pointlineplane/
	// Note: due to point deduplication, we do not need to check for zero length edges.

	vec21 := e.End.Subtract(e.Start)
	vec43 := otherEdge3D.End.Subtract(otherEdge3D.Start)
	vec13 := e.Start.Subtract(otherEdge3D.Start)

	dot4321 := vec43.DotProduct(vec21)
	dot4343 := vec43.DotProduct(vec43)
	dot2121 := vec21.DotProduct(vec21)
	dot1321 := vec13.DotProduct(vec21)

	denominator := (dot2121 * dot4343) - (dot4321 * dot4321)
	if math.Abs(denominator) < model.Threshold {
		// Edges are parallel, check if they are colinear, then return true if they overlap.

		// For this we can do similar math to the denominator, using vec13 (the start-to-start vector) as the "other" edge for this check.
		dot1313 := vec13.DotProduct(vec13)
		startToStartDenominator := (dot2121 * dot1313) - (dot1321 * dot1321)

		return math.Abs(startToStartDenominator) < model.Threshold && (model.IsBetween(e.Start.X, otherEdge3D.Start.X, otherEdge3D.End.X) ||
			model.IsBetween(e.End.X, otherEdge3D.Start.X, otherEdge3D.End.X) ||
			model.IsBetween(otherEdge3D.Start.X, e.Start.X, e.End.X) ||
			model.IsBetween(otherEdge3D.End.X, e.Start.X, e.End.X))
	}

	dot1343 := vec13.DotProduct(vec43)

	numerator := (dot1343 * dot4321) - (dot1321 * dot4343)
	percentE := numerator / denominator

	// If the closest point is not within the the start and end points, then the line segments do not intersect, even if the infinite lines do.
	if percentE < -model.Threshold || percentE > 1.0+model.Threshold {
		return false
	}

	percentOther := (dot1343 + (dot4321 * percentE)) / dot4343

	// If the closest point is not within the the start and end points, then the line segments do not intersect, even if the infinite lines do.
	if percentOther < -model.Threshold || percentOther > 1.0+model.Threshold {
		return false
	}

	pointE := NewVertex3D(e.Start.X+(percentE*vec21.X), e.Start.Y+(percentE*vec21.Y), e.Start.Z+(percentE*vec21.Z))

	pointOther := NewVertex3D(otherEdge3D.Start.X+(percentOther*vec43.X), otherEdge3D.Start.Y+(percentOther*vec43.Y), otherEdge3D.Start.Z+(percentOther*vec43.Z))

	return pointE.Equals(pointOther)
}

// Merge creates a new edge starting from this edge's start vertex and ending at the supplied edge's end vertex.
func (e *Edge3D) Merge(other model.CircuitEdge) model.CircuitEdge {
	return NewEdge3D(e.Start, other.GetEnd().(*Vertex3D))
}

// Split creates two new edges "start-to-vertex" and "vertex-to-end" based on this edge and the supplied vertex.
func (e *Edge3D) Split(vertex model.CircuitVertex) (model.CircuitEdge, model.CircuitEdge) {
	return NewEdge3D(e.Start, vertex.(*Vertex3D)), NewEdge3D(vertex.(*Vertex3D), e.End)
}

// String prints the edge as a string.
func (e *Edge3D) String() string {
	return fmt.Sprintf(`{"start":%s,"end":%s}`, e.Start.String(), e.End.String())
}

// NewEdge3D creates a edge from the starting Vertex3D to the ending Vertex3D
func NewEdge3D(start *Vertex3D, end *Vertex3D) *Edge3D {
	length := start.DistanceTo(end)
	return &Edge3D{
		Start:  start,
		End:    end,
		vector: nil,
		length: length,
	}
}

var _ model.CircuitEdge = (*Edge3D)(nil)
