package model

// DeleteVertex removes the vertex at the specified index in the supplied array, and returns the updated array.
// This may update the supplied array, so it should be updated with the returned array.
func DeleteVertex(vertices []CircuitVertex, index int) []CircuitVertex {
	if lastIndex := len(vertices) - 1; lastIndex < 0 {
		return vertices
	} else if index <= 0 {
		return vertices[1:]
	} else if index >= lastIndex {
		return vertices[:lastIndex]
	} else {
		return append(vertices[:index], vertices[index+1:]...)
	}
}

// IndexOfEdge returns the index (starting at 0) of the edge in the array. If the edge is not in the array, -1 will be returned.
func IndexOfEdge(edges []CircuitEdge, edge CircuitEdge) int {
	for index, e := range edges {
		// Compare pointers first, for performance, but then check start and end points, in case the same edge is created multiple times.
		if e == edge ||
			((e.GetStart() == edge.GetStart() || e.GetStart().DistanceTo(edge.GetStart()) < Threshold) &&
				(e.GetEnd() == edge.GetEnd() || e.GetEnd().DistanceTo(edge.GetEnd()) < Threshold)) {
			return index
		}
	}
	return -1
}

// IndexOfVertex returns the index (starting at 0) of the vertex in the array. If the vertex is not in the array, -1 will be returned.
func IndexOfVertex(vertices []CircuitVertex, vertex CircuitVertex) int {
	for index, v := range vertices {
		// Compare pointers first, for performance, but then check coordinates, in case the same vertex is created multiple times.
		if v == vertex || v.DistanceTo(vertex) <= Threshold {
			return index
		}
	}
	return -1
}

// MergeEdges combines the edges so that the attached vertex at the specified index is no longer used in the edges.
// The vertexIndex should be the index of the vertex in the GetAttachedVertices array.
// The array of CircuitEdge must be ordered so that the 0-th edge starts with the 0-th vertex in the GetAttachedVertices array.
// This may update the supplied array, so it should be updated with the returned array.
// In addition to returning the updated array, this also returns the two detached edges.
func MergeEdges(edges []CircuitEdge, vertexIndex int) ([]CircuitEdge, CircuitEdge, CircuitEdge) {
	var detachedEdgeA CircuitEdge
	var detachedEdgeB CircuitEdge

	// There must be at least 2 edges to merge edges.
	if lastIndex := len(edges) - 1; lastIndex <= 0 {
		return []CircuitEdge{}, nil, nil
	} else if vertexIndex <= 0 {
		detachedEdgeA = edges[lastIndex]
		detachedEdgeB = edges[0]
		edges = edges[1:]
		// Need additional -1 since array has one fewer element in it.
		edges[lastIndex-1] = detachedEdgeA.Merge(detachedEdgeB)
	} else {
		if vertexIndex >= lastIndex {
			vertexIndex = lastIndex
		}

		detachedEdgeA = edges[vertexIndex-1]
		detachedEdgeB = edges[vertexIndex]
		edges = append(edges[:vertexIndex-1], edges[vertexIndex:]...)
		edges[vertexIndex-1] = detachedEdgeA.Merge(detachedEdgeB)
	}

	return edges, detachedEdgeA, detachedEdgeB
}

// SplitEdge replaces the supplied edge with the two edges that are created by adding the supplied vertex to the edge.
// This requires the supplied edge to exist in the array of circuit edges.
// If it does not exist, the array will be returned unchanged, along with an index of -1.
// If it does exist, this will return the updated array and the index of the replaced edge (which becomes the index of the first new edge, the index of the second new edge is always that index+1).
func SplitEdge(edges []CircuitEdge, edgeToSplit CircuitEdge, vertexToAdd CircuitVertex) ([]CircuitEdge, int) {
	edgeIndex := IndexOfEdge(edges, edgeToSplit)
	if edgeIndex == -1 {
		return edges, -1
	}

	edgeA, edgeB := edgeToSplit.Split(vertexToAdd)
	// copy all elements starting at the index one to the right to create a duplicate record at index and index+1.
	edges = append(edges[:edgeIndex+1], edges[edgeIndex:]...)
	// replace both duplicated edges so that the previous edge is no longer in the circuit and the two supplied edges replace it.
	edges[edgeIndex] = edgeA
	edges[edgeIndex+1] = edgeB

	return edges, edgeIndex
}
