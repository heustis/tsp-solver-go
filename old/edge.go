package old

// Edge represents the line segment between two points
type Edge struct {
	start  *Vertex
	end    *Vertex
	vector *Vertex
	length float64
}

// Contains returns true if the edge contains the supplied vertex
func (e Edge) Contains(v *Vertex) bool {
	return e.start == v || e.end == v
}

// GetStart returns the start vertex of the edge
func (e Edge) GetStart() *Vertex {
	return e.start
}

// GetEnd returns the end vertex of the edge
func (e Edge) GetEnd() *Vertex {
	return e.end
}

// GetVector returns the normalized (length=1.0) vector from the edge's start to the edges end
func (e Edge) GetVector() *Vertex {
	return e.vector
}

// GetLength returns the length of the edge
func (e Edge) GetLength() float64 {
	return e.length
}

// DistanceIncrease returns the difference in length between the edge
// and the two edges formed by inserting the vertex between the edge's start and end.
// For example, if start->end has a length of 5, start->vertex has a length of 3,
//  and vertex->end has a length of 6, this will return 4 (6 + 3 - 5)
func (e Edge) DistanceIncrease(vertex *Vertex) float64 {
	return e.start.DistanceTo(vertex) + e.end.DistanceTo(vertex) - e.length
}

// NewEdge creates a edge from the starting Vertex to the ending Vertex
func NewEdge(start *Vertex, end *Vertex) *Edge {
	length := start.DistanceTo(end)
	vector := NewVertex((end.X-start.X)/length, (end.Y-start.Y)/length)
	return &Edge{start: start, end: end, vector: vector, length: length}
}
