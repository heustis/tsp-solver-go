package tspmodel

import (
	"container/list"
	"math"
)

// DeduplicateVerticesNoSorting is an O(n*n) algorithm for deduplicating an array of vertices, and returning a copy of the array containing only the unique vertices.
// This function does not modify the supplied array of vertices.
// Note 1: the order that vertices are encountered in an array matters for deduplication, for example:
//  - Given vertices A, B, and C.
//  - A and B are within tspmodel.Threshold of each other, as are B and C, but A and C are not.
//  - If A or C is encountered first, both A and C will be included in the output.
//  - However, if B is encountered first, only B will be in the output.
// Note 2: This should be used in situations where a hash map (O(n)) or sorting the array (O(n*log(n))) would be insufficient, for example:
//  - Sorting 3D vertices can result in duplicate entries in the output, due to having to sort by Y or Z first.
//  - If we sort by Y first, it is possible to have a point that matches in X and Y and is significantly different in Z in between points that match in X, Y, and Z.
//  - If vertices A and C are effectively equal, and vertex B has the same X and Y, but a significantly different Z, it could be sorted between A and C, preventing A and C from deduplicating.
func DeduplicateVerticesNoSorting(vertices []CircuitVertex) []CircuitVertex {
	// Note: we aren't using a set for deduplication due to using the Threshold for equality checks
	uniqueVertices := make([]CircuitVertex, 0, len(vertices))

	// Check each already added point to see if it is a duplicate of the current point.
	for _, v := range vertices {
		shouldAdd := true
		for _, added := range uniqueVertices {
			if v.Equals(added) {
				shouldAdd = false
				break
			}
		}

		if shouldAdd {
			uniqueVertices = append(uniqueVertices, v)
		}
	}
	return uniqueVertices
}

// DeleteVertex removes the vertex at the specified index in the supplied array, and returns the updated array.
// This may update the supplied array, so it should be updated with the returned array.
// This version of delete vertex should not be used if the array may get cloned or is a clone.
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

// DeleteVertexCopy returns a copy of the supplied array with the specified vertex removed.
// This does not modify the supplied array, so it is safe to use with algorithms that clone the vertex array.
func DeleteVertexCopy(vertices []CircuitVertex, toDelete CircuitVertex) []CircuitVertex {
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

// FindClosestEdge finds, and returns, the edge that is the closest to this vertex.
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

// FindClosestEdgeList finds, and returns, the edge that is the closest to this vertex in the supplied linked list.
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

// FindFarthestPoint finds the vertex in the array that is farthest from the supplied target vertex.
func FindFarthestPoint(target CircuitVertex, points []CircuitVertex) CircuitVertex {
	var farthestPoint CircuitVertex
	farthestDistance := 0.0

	for _, point := range points {
		if distance := point.DistanceTo(target); distance > farthestDistance {
			farthestDistance = distance
			farthestPoint = point
		}
	}

	return farthestPoint
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

// IsEdgeCloser returns true if the candidate edge is closer than the current closest edge.
func IsEdgeCloser(v CircuitVertex, candidateEdge CircuitEdge, currentEdge CircuitEdge) bool {
	return candidateEdge.DistanceIncrease(v) < currentEdge.DistanceIncrease(v)
}
