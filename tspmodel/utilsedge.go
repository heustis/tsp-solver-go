package tspmodel

import (
	"container/list"
)

// IndexOfEdge returns the index (starting at 0) of the edge in the array. If the edge is not in the array, -1 will be returned.
func IndexOfEdge(edges []CircuitEdge, edge CircuitEdge) int {
	for index, e := range edges {
		if e.Equals(edge) {
			return index
		}
	}
	return -1
}

// MergeEdgesByIndex combines the edges so that the attached vertex at the specified index is no longer used in the edges.
// The vertexIndex is the index of the edge starting with the vertex to detach in the edges array.
// After merging:
// - the replacement edge is stored in vertexIndex-1 (or the last entry if vertexIndex is 0),
// - the edge at vertexIndex is removed,
//  - the length of updatedEdges is one less than edges.
// This may update the supplied array, so it should be updated with the returned array.
// In addition to returning the updated array, this also returns the two detached edges.
func MergeEdgesByIndex(edges []CircuitEdge, vertexIndex int) (updatedEdges []CircuitEdge, detachedEdgeA CircuitEdge, detachedEdgeB CircuitEdge) {
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

// MergeEdgesByVertex combines the edges so that the supplied vertex is no longer used in the edges.
// The array of edges must be ordered so that the 0th edge starts with the 0th vertex in the GetAttachedVertices array, the 1st edge starts with the 1st vertex, and so on.
// If the supplied vertex is not in the array of edges, or there are too few edges (less than 2), the array will be returned unmodified (with nil for the other variables).
// After successfully merging:
// - the merged edge replaces the edge ending with the supplied vertex in the array,
// - the edge starting with the supplied vertex is removed from the array,
//  - the length of updatedEdges is one less than edges.
// This may update the supplied array, so it should be updated with the returned array.
// In addition to returning the updated array, this also returns the two detached edges and the merged edge.
func MergeEdgesByVertex(edges []CircuitEdge, vertex CircuitVertex) (updatedEdges []CircuitEdge, detachedEdgeA CircuitEdge, detachedEdgeB CircuitEdge, mergedEdge CircuitEdge) {
	vertexIndex := -1
	for i, e := range edges {
		if e.GetStart() == vertex {
			vertexIndex = i
			break
		}
	}

	if vertexIndex < 0 || len(edges) < 2 {
		return edges, nil, nil, nil
	}

	updated, detachedEdgeA, detachedEdgeB := MergeEdgesByIndex(edges, vertexIndex)
	updatedLen := len(updated)
	return updated, detachedEdgeA, detachedEdgeB, updated[(vertexIndex-1+updatedLen)%updatedLen]
}

// MergeEdgesCopy combines the edges so that the supplied vertex is no longer used in the edges.
// After successfully merging:
// - the merged edge replaces the edge ending with the supplied vertex in the array,
// - the edge starting with the supplied vertex is removed from the array,
//  - the length of updatedEdges is one less than edges.
// This does not modify the supplied array, so it is safe to use with algorithms that clone the edges array into multiple circuits.
// In addition to returning the updated array, this also returns the two detached edges and the merged edge.
func MergeEdgesCopy(edges []CircuitEdge, vertex CircuitVertex) (updatedEdges []CircuitEdge, detachedEdgeA CircuitEdge, detachedEdgeB CircuitEdge, mergedEdge CircuitEdge) {
	if len(edges) < 2 {
		return edges, nil, nil, nil
	}

	for i, e := range edges {
		if e.GetStart() == vertex {
			lenEdges := len(edges)
			updatedEdges = make([]CircuitEdge, 0, lenEdges-1)

			if i == 0 {
				detachedEdgeA = edges[lenEdges-1]
				detachedEdgeB = edges[i]
				mergedEdge = detachedEdgeA.Merge(detachedEdgeB)
				updatedEdges = append(updatedEdges, edges[1:lenEdges-1]...)
				updatedEdges = append(updatedEdges, mergedEdge)
				return updatedEdges, detachedEdgeA, detachedEdgeB, mergedEdge
			} else {
				detachedEdgeA = edges[i-1]
				detachedEdgeB = edges[i]
				mergedEdge = detachedEdgeA.Merge(detachedEdgeB)
				updatedEdges = append(updatedEdges, edges[:i-1]...)
				updatedEdges = append(updatedEdges, mergedEdge)
				if i < lenEdges-1 {
					updatedEdges = append(updatedEdges, edges[i+1:]...)
				}
				return updatedEdges, detachedEdgeA, detachedEdgeB, mergedEdge
			}
		}
	}

	return edges, nil, nil, nil

}

// MergeEdgesList combines the edges so that the supplied vertex is no longer used in the linked list of edges.
// In addition to updating the linked list, this also returns the two detached edges and the linked list element for the merged edge.
// If the vertex is not presed in the list of edges, it will be unmodified and nil will be returned.
func MergeEdgesList(edges *list.List, vertex CircuitVertex) (detachedEdgeA CircuitEdge, detachedEdgeB CircuitEdge, mergedLink *list.Element) {
	// There must be at least 2 edges to merge edges.
	if edges.Len() < 2 {
		return nil, nil, nil
	}

	for i, link := 0, edges.Front(); i < edges.Len(); i, link = i+1, link.Next() {
		detachedEdgeB = link.Value.(CircuitEdge)
		if detachedEdgeB.GetStart() == vertex {
			mergedLink = link.Prev()
			if link.Prev() == nil {
				mergedLink = edges.Back()
			}
			detachedEdgeA = mergedLink.Value.(CircuitEdge)
			mergedLink.Value = detachedEdgeA.Merge(detachedEdgeB)
			edges.Remove(link)
			return detachedEdgeA, detachedEdgeB, mergedLink
		}
	}
	return nil, nil, nil
}

// MoveVertex removes an attached vertex from its current location and moves it so that it splits the supplied edge.
// The vertices adjacent to the vertex's original location will be merged into a new edge.
// This may update the supplied array, so it should be updated with the returned array.
// In addition to returning the updated array, this also returns the merged edge and the two edges at the vertex's new location.
// If the vertex or edge is not in the circuit, this will return the original, unmodified, array and nil for the edges.
// Complexity: MoveVertex is O(N)
func MoveVertex(edges []CircuitEdge, vertex CircuitVertex, edge CircuitEdge) (updatedEdges []CircuitEdge, mergedEdge CircuitEdge, splitEdgeA CircuitEdge, splitEdgeB CircuitEdge) {
	// There must be at least 3 edges to move edges (2 that are initially attached to the vertex, and 1 other edge to attach it to).
	numEdges := len(edges)
	if numEdges < 3 {
		return edges, nil, nil, nil
	}

	// To avoid creating an extra array, this algorithm bubbles the second edge from the moved vertex's original location so that it is adjacent to the destination location.
	// These indices are used to enable that bubbling, and to track where to put the merged and split edges.
	mergedIndex, fromIndex, toIndex := -1, -1, -1

	prevIndex := len(edges) - 1
	for i, e := range edges {
		if e.GetStart() == vertex {
			mergedIndex = prevIndex
			fromIndex = i
		} else if e.GetStart() == edge.GetStart() && e.GetEnd() == edge.GetEnd() {
			toIndex = i
		}
		prevIndex = i
	}

	// If either index is less than zero, then either the vertex to move or the destination edge do not exist in the array.
	if fromIndex < 0 || toIndex < 0 {
		return edges, nil, nil, nil
	}

	// Merge the source edges, and split the destination edge.
	splitEdgeA, splitEdgeB = edge.Split(vertex)
	mergedEdge = edges[mergedIndex].Merge(edges[fromIndex])
	edges[mergedIndex] = mergedEdge

	// Determine whether the second edge (from the original location) needs to be bubbled up or down the array.
	var delta int
	if fromIndex > toIndex {
		delta = -1
		edges[toIndex] = splitEdgeA
		edges[fromIndex] = splitEdgeB
	} else {
		delta = 1
		edges[toIndex] = splitEdgeB
		edges[fromIndex] = splitEdgeA
	}

	for next := fromIndex + delta; next != toIndex; fromIndex, next = next, next+delta {
		edges[next], edges[fromIndex] = edges[fromIndex], edges[next]
	}

	return edges, mergedEdge, splitEdgeA, splitEdgeB
}

// SplitEdge replaces the supplied edge with the two edges that are created by adding the supplied vertex to the edge.
// This may update the supplied array, so it should be updated with the returned array.
// If the supplied edge does not exist, the array will be returned unchanged, along with an index of -1.
// If the supplied edge does exist, this will return the updated array and the index of the replaced edge (which becomes the index of the first new edge, the index of the second new edge is always that index+1).
func SplitEdge(edges []CircuitEdge, edgeToSplit CircuitEdge, vertexToAdd CircuitVertex) (updatedEdges []CircuitEdge, edgeIndex int) {
	edgeIndex = IndexOfEdge(edges, edgeToSplit)
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

// SplitEdgeCopy replaces the supplied edge with the two edges that are created by adding the supplied vertex to the edge.
// This does not modify the supplied array, so it is safe to use with algorithms that clone the edges array into multiple circuits.
// If the supplied edge does not exist, the array will be returned unchanged, along with an index of -1.
// If the supplied edge does exist, this will return the updated array and the index of the replaced edge (which becomes the index of the first new edge, the index of the second new edge is always that index+1).
func SplitEdgeCopy(edges []CircuitEdge, edgeToSplit CircuitEdge, vertexToAdd CircuitVertex) (updated []CircuitEdge, edgeIndex int) {
	updated = make([]CircuitEdge, len(edges)+1)
	edgeIndex = -1

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

// SplitEdgeList replaces the supplied edge with the two edges that are created by adding the supplied vertex to the edge.
// This requires the supplied edge to exist in the linked list of circuit edges.
// If it does not exist, the linked list will remain unchanged, and nil will be returned.
// If it does exist, this will update the linked list, and return the newly added element in the linked list.
func SplitEdgeList(edges *list.List, edgeToSplit CircuitEdge, vertexToAdd CircuitVertex) *list.Element {
	for i, link := 0, edges.Front(); i < edges.Len(); i, link = i+1, link.Next() {
		if edge := link.Value.(CircuitEdge); edge.Equals(edgeToSplit) {
			edgeA, edgeB := edge.Split(vertexToAdd)
			link.Value = edgeA
			return edges.InsertAfter(edgeB, link)
		}
	}
	return nil
}
