package model3d

import (
	"fmt"

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

// GetVector returns the normalized (length=1.0) vector from the edge's start to the edges end
func (e *Edge3D) GetVector() *Vertex3D {
	return e.vector
}

// GetLength returns the length of the edge
func (e *Edge3D) GetLength() float64 {
	return e.length
}

// Merge creates a new edge starting from this edge's start vertex and ending at the supplied edge's end vertex.
func (e *Edge3D) Merge(other model.CircuitEdge) model.CircuitEdge {
	return NewEdge3D(e.Start, other.GetEnd().(*Vertex3D))
}

// Split creates two new edges "start-to-vertex" and "vertex-to-end" based on this edge and the supplied vertex.
func (e *Edge3D) Split(vertex model.CircuitVertex) (model.CircuitEdge, model.CircuitEdge) {
	return NewEdge3D(e.Start, vertex.(*Vertex3D)), NewEdge3D(vertex.(*Vertex3D), e.End)
}

// ToString prints the edge as a string.
func (e *Edge3D) ToString() string {
	return fmt.Sprintf(`{"start":%s,"end":%s}`, e.Start.ToString(), e.End.ToString())
}

// NewEdge3D creates a edge from the starting Vertex3D to the ending Vertex3D
func NewEdge3D(start *Vertex3D, end *Vertex3D) *Edge3D {
	length := start.DistanceTo(end)
	vector := NewVertex3D((end.X-start.X)/length, (end.Y-start.Y)/length, (end.Z-start.Z)/length)
	return &Edge3D{
		Start:  start,
		End:    end,
		vector: vector,
		length: length,
	}
}

var _ model.CircuitEdge = (*Edge3D)(nil)
