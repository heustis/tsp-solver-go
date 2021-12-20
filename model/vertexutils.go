package model

import (
	"container/list"
	"math"
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

// FindClosestEdge finds, and returns, the edge that is the closest to the vertex.
func FindClosestEdge(vertex CircuitVertex, currentCircuit []CircuitEdge) CircuitEdge {
	var closest CircuitEdge = nil
	closestDistanceIncrease := math.MaxFloat64
	for _, candidate := range currentCircuit {
		if candidateDistanceIncrease := candidate.DistanceIncrease(vertex); candidateDistanceIncrease < closestDistanceIncrease {
			closest = candidate
			closestDistanceIncrease = candidateDistanceIncrease
		}
	}
	return closest
}

// FindClosestEdgeList finds, and returns, the edge that is the closest to the vertex in the supplied linked list.
func FindClosestEdgeList(vertex CircuitVertex, currentCircuit *list.List) CircuitEdge {
	var closest CircuitEdge = nil
	closestDistanceIncrease := math.MaxFloat64
	for i, link := 0, currentCircuit.Front(); i < currentCircuit.Len(); i, link = i+1, link.Next() {
		candidate := link.Value.(CircuitEdge)
		// Ignore edges already containing the vertex.
		if candidate.GetEnd() == vertex || candidate.GetStart() == vertex {
			continue
		}
		candidateDistanceIncrease := candidate.DistanceIncrease(vertex)
		if candidateDistanceIncrease < closestDistanceIncrease {
			closest = candidate
			closestDistanceIncrease = candidateDistanceIncrease
		}
	}
	return closest
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

// IsEdgeCloser checks if the supplied edge is closer than the current closest edge.
func IsEdgeCloser(v CircuitVertex, candidateEdge CircuitEdge, currentEdge CircuitEdge) bool {
	return candidateEdge.DistanceIncrease(v) < currentEdge.DistanceIncrease(v)
}
