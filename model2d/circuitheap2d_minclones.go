package model2d

import (
	"encoding/json"
	"fmt"

	"github.com/fealos/lee-tsp-go/model"
)

type HeapableCircuit2DMinClones struct {
	Vertices           []model.CircuitVertex
	circuit            []model.CircuitVertex
	circuitEdges       []model.CircuitEdge
	closestEdges       *model.Heap
	length             float64
	unattachedVertices map[model.CircuitVertex]bool
	convexVertices     map[model.CircuitVertex]bool
	// distanceIncreases  map[model.CircuitVertex]float64
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
		// AnyMatch is O(n), attachVertex is O(n)
		c.attachVertex(next)
		return nil
	} else {
		// 2b. If the 'next closest' vertex is already attached, clone the circuit with the 'next closest' vertex attached to the 'next closest' edge.
		// O(n)
		clone := &HeapableCircuit2DMinClones{
			Vertices:           make([]model.CircuitVertex, len(c.Vertices)),
			circuit:            make([]model.CircuitVertex, len(c.circuit)),
			circuitEdges:       make([]model.CircuitEdge, len(c.circuitEdges)),
			closestEdges:       c.closestEdges.Clone(),
			length:             c.length,
			unattachedVertices: make(map[model.CircuitVertex]bool),
			// distanceIncreases:  make(map[model.CircuitVertex]float64),
		}

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
		// for k, v := range c.distanceIncreases {
		// 	clone.distanceIncreases[k] = v
		// }

		// Need to update the heap after setting the distance increases, for the value function to properly heapify the data.
		// clone.closestEdges = model.NewHeap(getDistanceToEdgeForHeap)
		// clone.closestEdges.PushAllFrom(c.closestEdges)

		// Update one of the circuits' heaps to no longer have entries for this vertex.
		// The circuit chosen is the one which has the smaller distance increase for this vertex,
		// since the larger distance increase will make the vertex's other heapDistanceToEdges closer to the heap's root.
		// O(n)
		// heapToUpdate := c.closestEdges
		// if next.distance < c.distanceIncreases[next.vertex] {
		// 	heapToUpdate = clone.closestEdges
		// }
		clone.closestEdges.DeleteAll(func(x interface{}) bool {
			current := x.(*heapDistanceToEdge)
			return current.vertex.Equals(next.vertex)
		})

		// Move the vertex from the previous location to the new location in the clone.
		//O(n)
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
		delete(c.convexVertices, k)
	}
	// for k := range c.distanceIncreases {
	// 	delete(c.distanceIncreases, k)
	// }
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
		nextDist := nextDistToEdge.distance //- c.distanceIncreases[nextDistToEdge.vertex]
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
	// c.distanceIncreases = make(map[model.CircuitVertex]float64)
	c.closestEdges = model.NewHeap(getDistanceToEdgeForHeap)
	c.circuit = []model.CircuitVertex{}
	c.circuitEdges = []model.CircuitEdge{}
	c.length = 0.0

	c.Vertices = DeduplicateVertices(c.Vertices)

	// for _, v := range c.Vertices {
	// 	c.distanceIncreases[v] = 0.0
	// }
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
	c.length += edgeA.GetLength() + edgeB.GetLength() - toAttach.edge.GetLength()
	// Complexity is O(n), due to updateDistanceIncrease looking through the circuit for the index of the vertex.
	// c.distanceIncreases[toAttach.vertex] = toAttach.distance
	// c.updateDistanceIncrease(toAttach.edge.GetStart())
	// c.updateDistanceIncrease(toAttach.edge.GetEnd())

	// 4. Replace any references to the merged edge with two entries for the newly created edges..
	//    Complexity is O(n)
	distanceIncreases := c.getDistanceIncreases()
	c.closestEdges.ReplaceAll(func(x interface{}) []interface{} {
		current := x.(*heapDistanceToEdge)
		existingIncrease, okay := distanceIncreases[current.vertex]
		if !okay {
			existingIncrease = 0.0
		}
		if current.edge.Equals(toAttach.edge) {
			if current.vertex.Equals(toAttach.vertex) {
				return []interface{}{}
			}
			return []interface{}{
				&heapDistanceToEdge{
					vertex:   current.vertex,
					edge:     edgeA,
					distance: edgeA.DistanceIncrease(current.vertex) - existingIncrease,
				},
				&heapDistanceToEdge{
					vertex:   current.vertex,
					edge:     edgeB,
					distance: edgeB.DistanceIncrease(current.vertex) - existingIncrease,
				},
			}
		} else if current.vertex.Equals(toAttach.vertex) || current.vertex.Equals(toAttach.edge.GetStart()) || current.vertex.Equals(toAttach.edge.GetEnd()) {
			return []interface{}{
				&heapDistanceToEdge{
					vertex:   current.vertex,
					edge:     current.edge,
					distance: current.edge.DistanceIncrease(current.vertex) - existingIncrease,
				},
			}
		} else {
			return []interface{}{x}
		}
	})
}

func (c *HeapableCircuit2DMinClones) detachVertex(toDetach model.CircuitVertex) {
	// 1. Remove the vertex from the circuit
	c.circuit = model.DeleteVertex2(c.circuit, toDetach)

	// 2. Remove the edge with the vertex from the circuitEdges
	var detachedEdgeA, detachedEdgeB, mergedEdge model.CircuitEdge
	c.circuitEdges, detachedEdgeA, detachedEdgeB, mergedEdge = model.MergeEdges2(c.circuitEdges, toDetach)

	// 4. Update the circuit distance by removing the old distance increase and adding in the new distance increase.
	c.unattachedVertices[toDetach] = true
	c.length += mergedEdge.GetLength() - detachedEdgeA.GetLength() - detachedEdgeB.GetLength()
	// c.distanceIncreases[toDetach]
	// c.distanceIncreases[toDetach] = 0
	// c.updateDistanceIncrease(mergedEdge.GetStart())
	// c.updateDistanceIncrease(mergedEdge.GetEnd())

	// 5. Replace any references to the merged edges in the heap with a single entry for the merged edge.
	//    Complexity is O(n)
	distanceIncreases := c.getDistanceIncreases()
	replacedVertices := make(map[*Vertex2D]bool)
	c.closestEdges.ReplaceAll(func(x interface{}) []interface{} {
		current := x.(*heapDistanceToEdge)
		existingIncrease, okay := distanceIncreases[current.vertex]
		if !okay {
			existingIncrease = 0.0
		}
		if current.edge.Equals(detachedEdgeA) || current.edge.Equals(detachedEdgeB) {
			if replacedVertices[current.vertex] {
				// Only create one entry for the merged edge, even if the source heap has two entries.
				return []interface{}{}
			}
			replacedVertices[current.vertex] = true
			// Do not allow an entry to be created for either of the vertices of the merged edge.
			// For example, if point B is detached from between points A & C, point C could have an existing entry for A-B, which would be replaced by A-C.
			// The way that this scenario happens is that B and C are both internal points, B is attached first, C is attached between B and D, leaving an entry for A-B for vertex C.
			if current.vertex.Equals(mergedEdge.GetStart()) || current.vertex.Equals(mergedEdge.GetEnd()) {
				return []interface{}{}
			}
			return []interface{}{&heapDistanceToEdge{
				vertex:   current.vertex,
				edge:     mergedEdge,
				distance: mergedEdge.DistanceIncrease(current.vertex) - existingIncrease,
			}}
		} else if current.vertex.Equals(toDetach) || current.vertex.Equals(mergedEdge.GetStart()) || current.vertex.Equals(mergedEdge.GetEnd()) {
			return []interface{}{
				&heapDistanceToEdge{
					vertex:   current.vertex,
					edge:     current.edge,
					distance: current.edge.DistanceIncrease(current.vertex) - existingIncrease,
				},
			}
		} else {
			return []interface{}{x}
		}
	})
}

// func (c *HeapableCircuit2DMinClones) heapValueFunction(a interface{}) float64 {
// 	e := a.(*heapDistanceToEdge)
// 	return e.distance - c.distanceIncreases[e.vertex]
// }

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

// func (c *HeapableCircuit2DMinClones) updateDistanceIncrease(vertex model.CircuitVertex) {
// 	// No need to update convex vertices, since they will never be moved.
// 	if c.convexVertices[vertex] {
// 		return
// 	}
// 	index := model.IndexOfVertex(c.circuit, vertex)
// 	circuitLen := len(c.circuit)
// 	prev := c.circuit[(index-1+circuitLen)%circuitLen]
// 	next := c.circuit[(index+1)%circuitLen]
// 	c.distanceIncreases[vertex] = vertex.DistanceTo(prev) + vertex.DistanceTo(next) - prev.DistanceTo(next)
// }

func (c *HeapableCircuit2DMinClones) getDistanceIncreases() map[model.CircuitVertex]float64 {
	distanceIncreases := make(map[model.CircuitVertex]float64)
	circuitLen := len(c.circuit)
	if circuitLen >= 3 {
		prevIndex := circuitLen - 1
		nextIndex := 1
		for i, v := range c.circuit {
			// No need to include convex vertices, since they will never be moved.
			if !c.convexVertices[v] {
				prev := c.circuit[prevIndex]
				next := c.circuit[nextIndex]
				distanceIncreases[v] = v.DistanceTo(prev) + v.DistanceTo(next) - prev.DistanceTo(next)
			}
			prevIndex = i
			nextIndex = (nextIndex + 1) % circuitLen
		}
	}

	return distanceIncreases
}

var _ model.HeapableCircuit = (*HeapableCircuit2DMinClones)(nil)
