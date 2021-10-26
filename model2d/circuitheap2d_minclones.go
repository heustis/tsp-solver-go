package model2d

import (
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
	initialCandidates := []interface{}{}
	for _, edge := range c.circuitEdges {
		for v := range c.unattachedVertices {
			initialCandidates = append(initialCandidates, &heapDistanceToEdge{
				vertex:   v.(*Vertex2D),
				edge:     edge,
				distance: edge.DistanceIncrease(v),
			})
		}
		// c.closestEdges.PushAll(c.findCandidateVertices2(edge.(*Edge2D), c.unattachedVertices)...)
	}
	c.closestEdges.PushAll(initialCandidates...)
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
		// O(n)
		clone := &HeapableCircuit2DMinClones{
			Vertices:           make([]*Vertex2D, len(c.Vertices)),
			circuit:            make([]model.CircuitVertex, len(c.circuit)),
			circuitEdges:       make([]model.CircuitEdge, len(c.circuitEdges)),
			closestEdges:       nil,
			length:             c.length,
			unattachedVertices: make(map[model.CircuitVertex]bool),
			distanceIncreases:  make(map[model.CircuitVertex]float64),
		}
		clone.closestEdges = model.NewHeap(clone.heapValueFunction)
		clone.closestEdges.PushAllFrom(c.closestEdges)

		for i, v := range c.Vertices {
			clone.Vertices[i] = v
		}
		for i, c := range c.circuit {
			clone.circuit[i] = c
		}
		for i, e := range c.circuitEdges {
			clone.circuitEdges[i] = e
		}
		for k, v := range c.unattachedVertices {
			clone.unattachedVertices[k] = v
		}
		for k, v := range c.distanceIncreases {
			clone.distanceIncreases[k] = v
		}

		// Update one of the circuits' heaps to no longer have entries for this vertex.
		// The circuit chosen is the one which has the smaller distance increase for this vertex,
		// since the larger distance increase will make the vertex's other heapDistanceToEdges closer to the heap's root.
		// O(n)
		heapToUpdate := clone.closestEdges
		if c.distanceIncreases[next.vertex] < next.distance {
			heapToUpdate = c.closestEdges
		}
		heapToUpdate.DeleteAll(func(x interface{}) bool {
			current := x.(*heapDistanceToEdge)
			return current.vertex.Equals(next.vertex)
		})

		// Move the vertex from the previous location to the new location in the clone.
		//O(TODO)
		clone.detachVertex(next.vertex)
		clone.attachVertex(next)

		return clone
	}
}

func (c *HeapableCircuit2DMinClones) Delete() {
	for k := range c.unattachedVertices {
		delete(c.unattachedVertices, k)
	}
	for k := range c.convexVertices {
		delete(c.distanceIncreases, k)
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
		nextDistToEdge := next.(*heapDistanceToEdge)
		nextDist := nextDistToEdge.distance - c.distanceIncreases[nextDistToEdge.vertex]
		if len(c.unattachedVertices) == 0 && nextDist > 0 {
			return c.length // If the circuit is complete and the next vertex to attach increases the perimeter length, the circuit is optimal.
		} else {
			return c.length + nextDist
		}
	} else {
		return c.length
	}
}

func (c *HeapableCircuit2DMinClones) GetUnattachedVertices() map[model.CircuitVertex]bool {
	return c.unattachedVertices
}

func (c *HeapableCircuit2DMinClones) Prepare() {
	c.convexVertices = make(map[model.CircuitVertex]bool)
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
	c.circuitEdges, edgeIndex = model.SplitEdge2(c.circuitEdges, toAttach.edge, toAttach.vertex)
	if edgeIndex < 0 {
		expectedEdgeJson, _ := json.Marshal(toAttach.edge)
		actualCircuitJson, _ := json.Marshal(c.circuitEdges)
		initialVertices, _ := json.Marshal(c.Vertices)
		panic(fmt.Errorf("edge not found in circuit=%p, expected=%s, \ncircuit=%s \nvertices=%s", c, string(expectedEdgeJson), string(actualCircuitJson), string(initialVertices)))
	}
	edgeA, edgeB := c.circuitEdges[edgeIndex], c.circuitEdges[edgeIndex+1]

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

	// 4. Replace any references to the merged edge with two entries for the newly created edges..
	//    Complexity is O(n)
	c.closestEdges.ReplaceAll(func(x interface{}) []interface{} {
		current := x.(*heapDistanceToEdge)
		if current.edge.Equals(toAttach.edge) {
			if current.vertex.Equals(toAttach.vertex) {
				return []interface{}{}
			}
			return []interface{}{
				&heapDistanceToEdge{
					vertex:   current.vertex,
					edge:     edgeA,
					distance: edgeA.DistanceIncrease(current.vertex),
				},
				&heapDistanceToEdge{
					vertex:   current.vertex,
					edge:     edgeB,
					distance: edgeB.DistanceIncrease(current.vertex),
				},
			}
		} else {
			return []interface{}{x}
		}
	})

	// deleted := c.closestEdges.DeleteAll(func(x interface{}) bool {
	// 	current := x.(*heapDistanceToEdge)
	// 	return current.edge.Equals(toAttach.edge)
	// })

	// for _, x := range deleted {
	// 	current := x.(*heapDistanceToEdge)
	// 	c.closestEdges.PushHeap(&heapDistanceToEdge{
	// 		vertex:   current.vertex,
	// 		edge:     edgeA,
	// 		distance: edgeA.DistanceIncrease(current.vertex),
	// 	})
	// 	c.closestEdges.PushHeap(&heapDistanceToEdge{
	// 		vertex:   current.vertex,
	// 		edge:     edgeB,
	// 		distance: edgeB.DistanceIncrease(current.vertex),
	// 	})
	// }

	// 5. Add entries for each of the new edges - include any interior vertices, attached or unattached,
	//    that are not to the right of the edge nor would result in a different vertex becoming exterior if the vertex were attached to this edge.
	//    Complexity is O(n^2)
	// newEntries := append(c.findCandidateVertices(edgeA.(*Edge2D), deleted), c.findCandidateVertices(edgeB.(*Edge2D), deleted)...)
	// c.closestEdges.PushAll(newEntries...)
}

func (c *HeapableCircuit2DMinClones) detachVertex(toDetach model.CircuitVertex) {
	// 1. Remove the vertex from the circuit
	c.circuit = model.DeleteVertex2(c.circuit, toDetach)

	// 2. Remove the edge with the vertex from the circuitEdges
	var detachedEdgeA, detachedEdgeB, mergedEdge model.CircuitEdge
	c.circuitEdges, detachedEdgeA, detachedEdgeB, mergedEdge = model.MergeEdges2(c.circuitEdges, toDetach)

	// 4. Update the circuit distance by removing the old distance increase and adding in the new distance increase.
	c.unattachedVertices[toDetach] = true
	c.length -= c.distanceIncreases[toDetach]
	c.distanceIncreases[toDetach] = 0
	c.updateDistanceIncrease(mergedEdge.GetStart())
	c.updateDistanceIncrease(mergedEdge.GetEnd())

	// 5. Replace any references to the merged edges in the heap with a single entry for the merged edge.
	//    Complexity is O(n)
	c.closestEdges.ReplaceAll(func(x interface{}) []interface{} {
		current := x.(*heapDistanceToEdge)
		if current.edge.Equals(detachedEdgeA) {
			return []interface{}{&heapDistanceToEdge{
				vertex:   current.vertex,
				edge:     mergedEdge,
				distance: mergedEdge.DistanceIncrease(current.vertex),
			}}
		} else if current.edge.Equals(detachedEdgeB) {
			return []interface{}{}
		} else {
			return []interface{}{x}
		}
	})

	// deleted := c.closestEdges.DeleteAll(func(x interface{}) bool {
	// 	current := x.(*heapDistanceToEdge)
	// 	return current.edge.Equals(detachedEdgeA) || current.edge.Equals(detachedEdgeB)
	// })

	// addedVertices := make(map[*Vertex2D]bool)
	// for _, x := range deleted {
	// 	current := x.(*heapDistanceToEdge)
	// 	if !addedVertices[current.vertex] {
	// 		addedVertices[current.vertex] = true
	// 		c.closestEdges.PushHeap(&heapDistanceToEdge{
	// 			vertex:   current.vertex,
	// 			edge:     mergedEdge,
	// 			distance: mergedEdge.DistanceIncrease(current.vertex),
	// 		})
	// 	}
	// }

	// 6. Add entries for the new edge - include any interior vertices, attached or unattached,
	//    that are not to the right of the edge nor would result in a different vertex becoming exterior if the vertex were attached to this edge.
	//    Complexity is O(n^2)
	// newEntries := c.findCandidateVertices(mergedEdge.(*Edge2D), deleted)
	// c.closestEdges.PushAll(newEntries...)
}

func (c *HeapableCircuit2DMinClones) findCandidateVertices(edge *Edge2D, initialCandidates []interface{}) []interface{} {
	validCandidates := make(map[model.CircuitVertex]bool)
	for _, c := range initialCandidates {
		validCandidates[c.(*heapDistanceToEdge).vertex] = true
	}
	return c.findCandidateVertices2(edge, validCandidates)
}

func (c *HeapableCircuit2DMinClones) findCandidateVertices2(edge *Edge2D, validCandidates map[model.CircuitVertex]bool) []interface{} {
	candiates := []*heapDistanceToEdge{}

	// O(n) find interior vertices that can have this edge attached to them.
	for v := range validCandidates {
		// TODO - either only include validCandidates that are in the heap for the source edge, or do not filter the valid candidates (by leftOf, or the sorting + filtering below).
		// Otherwise, can attach a vertex, move the vertex, and prevent a different vertex from then attaching to the merged edge.

		// Ignore vertices that are part of the convex hull, or the edge itself.
		// Also ignore any vertices that are to the right of the edge:
		//   The circuit is counter-clockwise, so a point to the right of an edge will have another edge between this edge and that vertex.
		//   If that vertex were attached to this edge, it would automatically result in a longer circuit since it would result in crossed edges.
		v2d := v.(*Vertex2D)
		if !c.convexVertices[v] && v != edge.GetEnd() && v != edge.GetStart() && v2d.isLeftOf(edge) {
			candiates = append(candiates, &heapDistanceToEdge{
				vertex:   v2d,
				edge:     edge,
				distance: edge.DistanceIncrease(v),
			})
		}
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
