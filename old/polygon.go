package old

// Polygon is a representation of a 2-dimensional shape,
// with some additional fields that are useful for solving the Traveling Salesperson problem.
// Polygon objects are mutable!
type Polygon struct {
	vertices []*Vertex
	edges    []*Edge
	edgeMap  map[*Edge]bool
}

// GetVertices returns the vertices in the order that they form the perimeter.
func (poly *Polygon) GetVertices() []*Vertex {
	return poly.vertices
}

// GetEdges returns the edges that form the perimeter of the polygon.
// The first edge is from the 0th Vertex to the 1st Vertex.
// The second edge is from the 1st Vertex to the 2nd Vertex.
// ....
// The last edge is from the last (n-th) Vertex to the 0th Vertex.
func (poly *Polygon) GetEdges() []*Vertex {
	return poly.vertices
}

// MovePoint updates the polygon by moving the supplied vertex from its current location
// to between the two points in the supplied edge.
func (poly *Polygon) MovePoint(toMove *Vertex, edge *Edge) {
	var updatedPerimiter []*Vertex
	var updatedEdges []*Edge

	numCurrentVertices := len(poly.vertices)

	for i, current := range poly.vertices {
		next := poly.vertices[(i+1)%numCurrentVertices]

		if current == toMove {
			delete(poly.edgeMap, poly.edges[i])
		} else if next == toMove {
			updatedPerimiter = append(updatedPerimiter, current)

			edgeWithoutMoved := NewEdge(current, poly.vertices[(i+2)%numCurrentVertices])
			updatedEdges = append(updatedEdges, edgeWithoutMoved)
			poly.edgeMap[edgeWithoutMoved] = true

			delete(poly.edgeMap, poly.edges[i])
		} else if current == edge.GetStart() {
			updatedPerimiter = append(updatedPerimiter, toMove)

			currentToMoved := NewEdge(current, toMove)
			updatedEdges = append(updatedEdges, currentToMoved)
			poly.edgeMap[currentToMoved] = true

			movedToNext := NewEdge(toMove, next)
			updatedEdges = append(updatedEdges, movedToNext)
			poly.edgeMap[movedToNext] = true

			delete(poly.edgeMap, poly.edges[i])
		} else {
			updatedPerimiter = append(updatedPerimiter, current)
			updatedEdges = append(updatedEdges, poly.edges[i])
		}
	}

	poly.vertices = updatedPerimiter
	poly.edges = updatedEdges
}

// AddPoint updates the polygon with the supplied vertex added between the two points in the supplied edge.
func (poly *Polygon) AddPoint(toAdd *Vertex, edge *Edge) {
	var updatedPerimiter []*Vertex
	var updatedEdges []*Edge

	numCurrentVertices := len(poly.vertices)

	for i, current := range poly.vertices {
		updatedPerimiter = append(updatedPerimiter, current)

		next := poly.vertices[(i+1)%numCurrentVertices]

		if current == edge.GetStart() {
			updatedPerimiter = append(updatedPerimiter, toAdd)
			updatedEdges = append(updatedEdges, NewEdge(current, toAdd))
			updatedEdges = append(updatedEdges, NewEdge(toAdd, next))
		} else {
			updatedEdges = append(updatedEdges, poly.edges[i])
		}
	}

	poly.vertices = updatedPerimiter
	poly.edges = updatedEdges
}

// FindClosestEdge determines which of the polygon's edges is closest to the supplied vertex.
// Returns DistanceToEdge representing the closest edge and the distance to the vertex.
func (poly *Polygon) FindClosestEdge(vertex *Vertex) *DistanceToEdge {
	var closestEdge *DistanceToEdge
	for i, edge := range poly.edges {
		if i == 0 {
			closestEdge = NewDistanceToEdge(edge, edge.DistanceIncrease(vertex))
		} else {
			_, closestEdge = closestEdge.IsCloser(edge, vertex)
		}
	}
	return closestEdge
}

// NewPolygon creates a new polygon object
func NewPolygon(vertices []*Vertex) *Polygon {
	var edges []*Edge
	edgeMap := make(map[*Edge]bool)

	numVertices := len(vertices)
	for i := 0; i < numVertices; i++ {
		edge := NewEdge(vertices[i], vertices[(i+1)%numVertices])
		edges = append(edges, edge)
		edgeMap[edge] = true
	}
	return &Polygon{
		edges:    edges,
		edgeMap:  edgeMap,
		vertices: vertices,
	}
}
