package model2d

import (
	"fmt"

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

// GetVector returns the normalized (length=1.0) vector from the edge's start to the edges end
func (e *Edge2D) GetVector() *Vertex2D {
	return e.vector
}

// GetLength returns the length of the edge
func (e *Edge2D) GetLength() float64 {
	return e.length
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
