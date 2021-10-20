package model2d

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/fealos/lee-tsp-go/model"
)

type HeapableCircuit2DMinClones struct {
	Vertices           []*Vertex2D
	circuit            []model.CircuitVertex
	circuitEdges       []model.CircuitEdge
	closestEdges       *model.Heap
	length             float64
	unattachedVertices map[model.CircuitVertex]bool
	convexVertices     map[model.CircuitVertex]bool
	distanceIncreases  map[model.CircuitVertex]float64
}

func (c *HeapableCircuit2DMinClones) BuildPerimiter() {
	c.circuit, c.circuitEdges, c.unattachedVertices = (&PerimeterBuilder2D{}).BuildPerimiter(c.Vertices)

	// Determine the initial length of the perimeter.
	for _, edge := range c.circuitEdges {
		c.length += edge.GetLength()
	}

	// Track the vertices in the initial circuit in convexVertices
	for _, vertex := range c.circuit {
		c.convexVertices[vertex] = true
	}

	// Find the closest edge for all interior points, based on distance increase (rather than perpendicular distance)
	// total vertices = attached + unattached
	// complexity  = attached * unattached  = attached * (total - attached)  = total*attached - attached^2
	for _, edge := range c.circuitEdges {
		for v := range c.unattachedVertices {
			heap.Push(c.closestEdges, &heapDistanceToEdge{
				edge:     edge,
				distance: edge.DistanceIncrease(v),
				vertex:   v.(*Vertex2D),
			})
		}
	}
}

func (c *HeapableCircuit2DMinClones) CloneAndUpdate() model.HeapableCircuit {
	// 1. Remove 'next closest' from heap - complexity O(log n)
	next, okay := c.closestEdges.PopHeap().(*heapDistanceToEdge)

	if next == nil || !okay {
		return nil
	} else if c.unattachedVertices[next.vertex] || !c.closestEdges.AnyMatch(next.hasVertex) {
		// 2a. If there are no more items in the heap with the vertex in 'next closest', this is the last edge that the vertex can attach to (i.e. all other possibilities have been tried).
		//     Similarly, if this is the first time we are encountering a vertex, attach it at the specified location without cloning.
		//     So do not clone the circuit, rather attach 'next closest' to this circuit.
		// AnyMatch is O(n), attachVertex is O(TODO)
		c.attachVertex(next)
		return nil
	} else {
		// 2b. If the 'next closest' vertex is already attached, clone the circuit with the 'next closest' vertex attached to the 'next closest' edge.
		// Complexity of cloning is O(n)
		clone := &HeapableCircuit2DMinClones{
			Vertices:           append(make([]*Vertex2D, 0, len(c.Vertices)), c.Vertices...),
			circuit:            append(make([]model.CircuitVertex, 0, len(c.circuit)), c.circuit...),
			circuitEdges:       append(make([]model.CircuitEdge, 0, len(c.circuitEdges)), c.circuitEdges...),
			closestEdges:       c.closestEdges.Clone(),
			length:             c.length - c.distanceIncreases[next.vertex],
			unattachedVertices: make(map[model.CircuitVertex]bool),
			distanceIncreases:  make(map[model.CircuitVertex]float64),
		}
		for k, v := range c.unattachedVertices {
			clone.unattachedVertices[k] = v
		}
		for k, v := range c.distanceIncreases {
			clone.distanceIncreases[k] = v
		}
		clone.distanceIncreases[next.vertex] = next.distance

		clone.closestEdges.DeleteAll(func(x interface{}) bool {
			return false //TODO
		})

		clone.detachVertex(next.vertex)
		clone.attachVertex(next)

		// Update the current circuit's heap to no longer have entries for this vertex, since the clone will have closer entries for each of them (due to adjusting the distance)
		c.closestEdges.ReplaceAll(func(x interface{}) []interface{} {
			current := x.(*heapDistanceToEdge)
			if current.vertex == next.vertex {
				return []interface{}{}
			}
			return []interface{}{current}
		})

		return clone
	}
}

func (c *HeapableCircuit2DMinClones) Delete() {
	for k := range c.unattachedVertices {
		delete(c.unattachedVertices, k)
	}
	for k := range c.distanceIncreases {
		delete(c.distanceIncreases, k)
	}
	if c.closestEdges != nil {
		c.closestEdges.Delete()
	}
	c.Vertices = nil
	c.circuit = nil
	c.circuitEdges = nil
	c.closestEdges = nil
}

func (c *HeapableCircuit2DMinClones) GetAttachedVertices() []model.CircuitVertex {
	return c.circuit
}

func (c *HeapableCircuit2DMinClones) GetLength() float64 {
	return c.length
}

func (c *HeapableCircuit2DMinClones) GetLengthWithNext() float64 {
	if next := c.closestEdges.Peek(); next != nil {
		nextDist := next.(*heapDistanceToEdge)
		if len(c.unattachedVertices) == 0 && nextDist.distance > 0 {
			return c.length // If the circuit is complete and the next vertex to attach increases the perimeter length, the circuit is optimal.
		} else {
			return c.length + nextDist.distance
		}
	} else {
		return c.length
	}
}

func (c *HeapableCircuit2DMinClones) GetUnattachedVertices() map[model.CircuitVertex]bool {
	return c.unattachedVertices
}

func (c *HeapableCircuit2DMinClones) Prepare() {
	c.unattachedVertices = make(map[model.CircuitVertex]bool)
	c.distanceIncreases = make(map[model.CircuitVertex]float64)
	c.closestEdges = model.NewHeap(c.heapValueFunction)
	c.circuit = []model.CircuitVertex{}
	c.circuitEdges = []model.CircuitEdge{}
	c.length = 0.0

	c.Vertices = deduplicateVertices(c.Vertices)

	for _, v := range c.Vertices {
		c.distanceIncreases[v] = 0.0
	}
}

func (c *HeapableCircuit2DMinClones) attachVertex(toAttach *heapDistanceToEdge) {
	// 1. Update the circuitEdges and retrieve the newly created edges
	var edgeIndex int
	c.circuitEdges, edgeIndex = model.SplitEdge(c.circuitEdges, toAttach.edge, toAttach.vertex)
	if edgeIndex < 0 {
		expectedEdgeJson, _ := json.Marshal(toAttach.edge)
		actualCircuitJson, _ := json.Marshal(c.circuitEdges)
		panic(fmt.Errorf("edge not found in circuit, expected=%s, circuit=%s", string(expectedEdgeJson), string(actualCircuitJson)))
	}
	edgeLen := len(c.circuitEdges)
	edgeA, edgeB := c.circuitEdges[edgeIndex], c.circuitEdges[edgeIndex+1%edgeLen]

	// 2. Update the circuit
	// Complexity is O(n), due to shifting array elements to insert new vertex.
	c.insertVertex(edgeIndex+1, toAttach.vertex)
	delete(c.unattachedVertices, toAttach.vertex)

	// 3. Update the circuit length and the distances increases as a result of the attached vertex.
	// Complexity is O(n), due to updateDistanceIncrease looking through the circuit for the index of the vertex.
	c.length += toAttach.distance
	c.distanceIncreases[toAttach.vertex] = toAttach.distance
	c.updateDistanceIncrease(toAttach.edge.GetStart())
	c.updateDistanceIncrease(toAttach.edge.GetEnd())

	// 4. Remove any references to the removed edge from the heap.
	//    Complexity is O(n)
	c.closestEdges.DeleteAll(func(x interface{}) bool {
		current := x.(*heapDistanceToEdge)
		return current.edge == toAttach.edge
	})

	// 5. Add entries for each of the new edges - include any interior vertices, attached or unattached,
	//    that are not to the right of the edge nor would result in a different vertex becoming exterior if the vertex were attached to this edge.
	//    Complexity is O(n^2)
	newEntries := append(c.findCandidateVertices(edgeA.(*Edge2D)), c.findCandidateVertices(edgeB.(*Edge2D))...)
	c.closestEdges.PushAll(newEntries...)
}

func (c *HeapableCircuit2DMinClones) detachVertex(toDetach model.CircuitVertex) {
	// 1. Remove the vertex from the circuit
	vertexIndex := model.IndexOfVertex(c.circuit, toDetach)
	c.circuit = model.DeleteVertex(c.circuit, vertexIndex)

	// 2. Remove the edge with the vertex from the circuitEdges
	var detachedEdgeA model.CircuitEdge
	var detachedEdgeB model.CircuitEdge
	c.circuitEdges, detachedEdgeA, detachedEdgeB = model.MergeEdges(c.circuitEdges, vertexIndex)

	// 3. Retrieve the merged edge from the circuit
	edgeLen := len(c.circuitEdges)
	mergedEdge := c.circuitEdges[(vertexIndex-1+edgeLen)%edgeLen]

	// 4. Update the circuit distance by removing the old distance increase and adding in the new distance increase.
	c.unattachedVertices[toDetach] = true
	c.length -= c.distanceIncreases[toDetach]
	c.distanceIncreases[toDetach] = 0
	c.updateDistanceIncrease(mergedEdge.GetStart())
	c.updateDistanceIncrease(mergedEdge.GetEnd())

	// 5. Remove any references to the removed edges from the heap.
	//    Complexity is O(n)
	c.closestEdges.DeleteAll(func(x interface{}) bool {
		current := x.(*heapDistanceToEdge)
		return current.edge == detachedEdgeA || current.edge == detachedEdgeB
	})

	// 6. Add entries for the new edge - include any interior vertices, attached or unattached,
	//    that are not to the right of the edge nor would result in a different vertex becoming exterior if the vertex were attached to this edge.
	newEntries := c.findCandidateVertices(mergedEdge.(*Edge2D))
	c.closestEdges.PushAll(newEntries...)
}

func (c *HeapableCircuit2DMinClones) findCandidateVertices(edge *Edge2D) []interface{} {
	candiates := []*heapDistanceToEdge{}

	// O(n) find interior vertices that can have this edge attached to them.
	for _, v := range c.Vertices {
		// Ignore vertices that are part of the convex hull, or the edge itself.
		// Also ignore any vertices that are to the right of the edge.
		// The circuit is counter-clockwise, so a point to the right of an edge will have another edge between this edge and that vertex.
		// If that vertex were attached to this edge, it would automatically result in a longer circuit since it would result in crossed edges.
		if c.convexVertices[v] || v == edge.GetEnd() || v == edge.GetStart() || v.isRightOf(edge) {
			continue
		}
		candiates = append(candiates, &heapDistanceToEdge{
			vertex:   v,
			edge:     edge,
			distance: edge.DistanceIncrease(v),
		})
	}

	// O(n * log n) sort the candidate vertices to make the next filtering faster.
	sort.Slice(candiates, func(i, j int) bool {
		return candiates[i].distance < candiates[j].distance
	})

	// O(n^2) check each candidate to see if it would cause any other vertex (in the candidate list) to shift from interior to exterior.
	// If so, ignore the vertex, as the other vertex should be chosen first.
	// Note: only need to check vertices that are closer to the edge than the candidate vertex
	filteredCandiates := []interface{}{}
	for i, c := range candiates {
		splitEdgeA, splitEdgeB := edge.Split(c.vertex)

		include := true
		for j := i - 1; j >= 0 && include; j-- {
			if vertexJ := candiates[j].vertex; vertexJ.isRightOf(splitEdgeA.(*Edge2D)) && vertexJ.isRightOf(splitEdgeB.(*Edge2D)) {
				include = false
			}
		}

		if include {
			filteredCandiates = append(filteredCandiates, c)
		}
	}

	return filteredCandiates
}

func (c *HeapableCircuit2DMinClones) heapValueFunction(a interface{}) float64 {
	e := a.(*heapDistanceToEdge)
	return e.distance - c.distanceIncreases[e.vertex]
}

func (c *HeapableCircuit2DMinClones) insertVertex(index int, vertex model.CircuitVertex) {
	if index >= len(c.circuit) {
		c.circuit = append(c.circuit, vertex)
	} else {
		// copy all elements starting at the index one to the right to create a duplicate record at index and index+1.
		c.circuit = append(c.circuit[:index+1], c.circuit[index:]...)
		// update only the vertex at the index, so that there are no duplicates and the vertex is at the index.
		c.circuit[index] = vertex
	}
}

// func (c *HeapableCircuit2DMinClones) updateDistances() {
// 	circuitLen := len(c.circuit)
// 	c.length = 0.0
// 	for prev, i, next := circuitLen-1, 0, 1; i < circuitLen; prev, i, next = i, i+1, (next+1)%circuitLen {
// 		prevVertex, currentVertex, nextVertex := c.circuit[prev], c.circuit[i], c.circuit[next]
// 		prevToCurrent := prevVertex.DistanceTo(currentVertex)
// 		c.distanceIncreases[c.circuit[i]] = (prevToCurrent + currentVertex.DistanceTo(nextVertex)) - prevVertex.DistanceTo(nextVertex)
// 		c.length += prevToCurrent
// 	}
// 	for v := range c.unattachedVertices {
// 		c.distanceIncreases[v] = 0.0
// 	}
// }

func (c *HeapableCircuit2DMinClones) updateDistanceIncrease(vertex model.CircuitVertex) {
	// No need to update convex vertices, since they will never be moved.
	if c.convexVertices[vertex] {
		return
	}
	index := model.IndexOfVertex(c.circuit, vertex)
	circuitLen := len(c.circuit)
	prev := c.circuit[(index-1+circuitLen)%circuitLen]
	next := c.circuit[(index+1)%circuitLen]
	c.distanceIncreases[vertex] = vertex.DistanceTo(prev) + vertex.DistanceTo(next) - prev.DistanceTo(next)
}

var _ model.HeapableCircuit = (*HeapableCircuit2DMinClones)(nil)
