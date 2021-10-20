package old

// MaleablePoints tracks which points are not part of a Polygon's convex hull.
// These points can adjust which edge they are part of, as new Vertices are added to the Polygon.
type MaleablePoints struct {
	external           map[*Vertex]bool
	attachedVertices   map[*Vertex]*DistanceToEdge
	unattachedVertices map[*Vertex]*DistanceToEdge
	toUpdate           map[*Vertex]bool
}

// IsInteriorVertex returns true if the supplied vertex is not part of the Polygon's convex hull.
func (i MaleablePoints) IsInteriorVertex(v *Vertex) bool {
	_, exists := i.unattachedVertices[v]
	if exists {
		return exists
	}

	_, exists = i.attachedVertices[v]
	return exists
}

func (i MaleablePoints) UpdateAttachedPoints(edge *Edge) {
	for vertex, currentDistance := range i.attachedVertices {
		closer, distance := currentDistance.IsCloser(edge, vertex)
		if closer {
			i.attachedVertices[vertex] = distance
			i.toUpdate[vertex] = true
		}
	}
}

// UpdateClosestEdge checks if the distance to a
func (i MaleablePoints) UpdateClosestEdge(edge *Edge) {
	for vertex, currentDistance := range i.unattachedVertices {
		closer, distance := currentDistance.IsCloser(edge, vertex)
		if closer {
			i.unattachedVertices[vertex] = distance
		}
	}
}
