package model

// DeleteVertex removes the vertex at the specified index in the supplied array, and returns the updated array.
// This may or may not update the supplied array, so it should be updated with the returned array.
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

// DeleteVertex2 removes the specified vertex from the supplied array, and returns the updated array.
// This may or may not update the supplied array, so it should be updated with the returned array.
func DeleteVertex2(vertices []CircuitVertex, toDelete CircuitVertex) []CircuitVertex {
	updatedLen := len(vertices) - 1
	if updatedLen < 0 {
		return []CircuitVertex{}
	}
	updated := make([]CircuitVertex, updatedLen)
	i := 0
	for _, v := range vertices {
		if !v.Equals(toDelete) {
			if i >= updatedLen {
				return vertices
			}
			updated[i] = v
			i++
		}
	}
	return updated
}

// IndexOfEdge returns the index (starting at 0) of the edge in the array. If the edge is not in the array, -1 will be returned.
func IndexOfEdge(edges []CircuitEdge, edge CircuitEdge) int {
	for index, e := range edges {
		if e.Equals(edge) {
			return index
		}
	}
	return -1
}

// IndexOfVertex returns the index (starting at 0) of the vertex in the array. If the vertex is not in the array, -1 will be returned.
func IndexOfVertex(vertices []CircuitVertex, vertex CircuitVertex) int {
	for index, v := range vertices {
		if v.Equals(vertex) {
			return index
		}
	}
	return -1
}

// InsertVertex inserts the supplied vertex at the specified index, 0-based.
// If the index is greater than the last index in the array, the vertex will be appended to the end of the array.
// This may modify the supplied array, so it should be updated with the returned array.
func InsertVertex(vertices []CircuitVertex, index int, vertex CircuitVertex) []CircuitVertex {
	if index >= len(vertices) {
		return append(vertices, vertex)
	} else {
		// copy all elements starting at the index one to the right to create a duplicate record at index and index+1.
		vertices = append(vertices[:index+1], vertices[index:]...)
		// update only the vertex at the index, so that there are no duplicates and the vertex is at the index.
		vertices[index] = vertex
		return vertices
	}
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

// MergeEdges2 combines the edges so that the supplied vertex is no longer used in the edges.
// The array of CircuitEdge must be ordered so that the 0-th edge starts with the 0-th vertex in the GetAttachedVertices array.
// This may or may not update the supplied array, so it should be updated with the returned array.
// In addition to returning the updated array, this also returns the two detached edges and the merged edge.
func MergeEdges2(edges []CircuitEdge, vertex CircuitVertex) ([]CircuitEdge, CircuitEdge, CircuitEdge, CircuitEdge) {
	var detachedEdgeA CircuitEdge
	var detachedEdgeB CircuitEdge
	var mergedEdge CircuitEdge

	lastIndex := len(edges) - 1

	// There must be at least 2 edges to merge edges.
	if lastIndex <= 0 {
		return []CircuitEdge{}, nil, nil, nil
	}

	updated := make([]CircuitEdge, lastIndex)
	updatedIndex := 0
	for i, e := range edges {
		if e.GetEnd().Equals(vertex) {
			detachedEdgeA = e
			detachedEdgeB = edges[(i+1)%len(edges)]
			mergedEdge = detachedEdgeA.Merge(detachedEdgeB)
			updated[updatedIndex] = mergedEdge
			updatedIndex++
		} else if !e.GetStart().Equals(vertex) {
			updated[updatedIndex] = e
			updatedIndex++
		}
	}

	return updated, detachedEdgeA, detachedEdgeB, mergedEdge
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

// SplitEdge2 replaces the supplied edge with the two edges that are created by adding the supplied vertex to the edge.
// This requires the supplied edge to exist in the array of circuit edges.
// This will return the updated array and the index of the first new edge.
func SplitEdge2(edges []CircuitEdge, edgeToSplit CircuitEdge, vertexToAdd CircuitVertex) ([]CircuitEdge, int) {
	updated := make([]CircuitEdge, len(edges)+1)
	edgeIndex := -1

	updatedIndex := 0
	for _, e := range edges {
		if e.Equals(edgeToSplit) {
			edgeA, edgeB := edgeToSplit.Split(vertexToAdd)
			edgeIndex = updatedIndex
			updated[updatedIndex] = edgeA
			updatedIndex++
			updated[updatedIndex] = edgeB
		} else {
			updated[updatedIndex] = e
		}
		updatedIndex++
	}

	if edgeIndex == -1 {
		return edges, edgeIndex
	} else {
		return updated, edgeIndex
	}
}
