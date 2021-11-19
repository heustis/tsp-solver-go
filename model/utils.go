package model

import (
	"container/list"
	"encoding/json"
	"fmt"
	"strings"
)

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

	updated, detachedEdgeA, detachedEdgeB := MergeEdges(edges, vertexIndex)
	updatedLen := len(updated)
	return updated, detachedEdgeA, detachedEdgeB, updated[(vertexIndex-1+updatedLen)%updatedLen]
}

// MergeEdges2 combines the edges so that the supplied vertex is no longer used in the linked list of edges.
// In addition to updating the linked list, this also returns the two detached edges and the linked list element for the merged edge.
// If the vertex is not presed in the list of edges, it will be unmodified and nil will be returned.
func MergeEdgesList(edges *list.List, vertex CircuitVertex) (CircuitEdge, CircuitEdge, *list.Element) {
	for i, link := 0, edges.Front(); i < edges.Len(); i, link = i+1, link.Next() {
		edge := link.Value.(CircuitEdge)
		if edge.GetStart() == vertex {
			prev := link.Prev().Value.(CircuitEdge)
			link.Value = prev.Merge(edge)
			edges.Remove(link.Prev())
			return prev, edge, link
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
func MoveVertex(edges []CircuitEdge, vertex CircuitVertex, edge CircuitEdge) ([]CircuitEdge, CircuitEdge, CircuitEdge, CircuitEdge) {
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
	splitEdgeA, splitEdgeB := edge.Split(vertex)
	mergedEdge := edges[mergedIndex].Merge(edges[fromIndex])
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

// ToString converts an object to a string.
func ToString(value interface{}) string {
	if p, okay := value.(Printable); okay {
		return p.ToString()
	} else if jsonBytes, err := json.Marshal(value); err == nil && strings.Compare("null", string(jsonBytes)) != 0 {
		return string(jsonBytes)
	} else {
		return fmt.Sprintf("%v", value)
	}
}
