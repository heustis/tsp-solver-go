package model

import (
	"encoding/json"
	"fmt"
)

type HeapableCircuitMinClones struct {
	Vertices           []CircuitVertex
	deduplicator       func([]CircuitVertex) []CircuitVertex
	perimeterBuilder   PerimeterBuilder
	circuit            []CircuitVertex
	circuitEdges       []CircuitEdge
	closestEdges       *Heap
	length             float64
	unattachedVertices map[CircuitVertex]bool
	convexVertices     map[CircuitVertex]bool
}

func CreateHeapableCircuitMinClones(vertices []CircuitVertex, deduplicator func([]CircuitVertex) []CircuitVertex, perimeterBuilder PerimeterBuilder) *HeapableCircuitMinClones {
	return &HeapableCircuitMinClones{
		Vertices:         vertices,
		deduplicator:     deduplicator,
		perimeterBuilder: perimeterBuilder,
	}
}

func (c *HeapableCircuitMinClones) BuildPerimiter() {
	c.circuit, c.circuitEdges, c.unattachedVertices = c.perimeterBuilder.BuildPerimiter(c.Vertices)

	// Determine the initial length of the perimeter.
	c.length = 0.0
	for _, edge := range c.circuitEdges {
		c.length += edge.GetLength()
	}

	// Track the vertices in the initial circuit in convexVertices
	c.convexVertices = make(map[CircuitVertex]bool)
	for _, vertex := range c.circuit {
		c.convexVertices[vertex] = true
	}

	// Find the closest edge for all interior points, based on distance increase (rather than perpendicular distance)
	// total vertices = attached + unattached
	// complexity  = attached * unattached  = attached * (total - attached)  = total*attached - attached^2
	initialCandidates := []interface{}{}
	for _, edge := range c.circuitEdges {
		for v := range c.unattachedVertices {
			initialCandidates = append(initialCandidates, &DistanceToEdge{
				Vertex:   v,
				Edge:     edge,
				Distance: edge.DistanceIncrease(v),
			})
		}
	}
	c.closestEdges = NewHeap(GetDistanceToEdgeForHeap)
	c.closestEdges.PushAll(initialCandidates...)
}

func (c *HeapableCircuitMinClones) CloneAndUpdate() HeapableCircuit {
	// 1. Remove 'next closest' from heap - complexity O(log n)
	next, okay := c.closestEdges.PopHeap().(*DistanceToEdge)

	if next == nil || !okay {
		return nil
	} else if c.unattachedVertices[next.Vertex] || !c.closestEdges.AnyMatch(next.HasVertex) {
		// 2a. If there are no more items in the heap with the vertex in 'next closest', this is the last edge that the vertex can attach to (i.e. all other possibilities have been tried).
		//     Similarly, if this is the first time we are encountering a vertex, attach it at the specified location without cloning.
		//     So do not clone the circuit, rather attach 'next closest' to this circuit.
		// AnyMatch is O(n), attachVertex is O(n)
		c.AttachVertex(next)
		return nil
	} else {
		// 2b. If the 'next closest' vertex is already attached, clone the circuit with the 'next closest' vertex attached to the 'next closest' edge.
		// O(n)
		clone := &HeapableCircuitMinClones{
			Vertices:           make([]CircuitVertex, len(c.Vertices)),
			circuit:            make([]CircuitVertex, len(c.circuit)),
			circuitEdges:       make([]CircuitEdge, len(c.circuitEdges)),
			closestEdges:       c.closestEdges.Clone(),
			length:             c.length,
			unattachedVertices: make(map[CircuitVertex]bool),
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
		clone.closestEdges.DeleteAll(func(x interface{}) bool {
			current := x.(*DistanceToEdge)
			return current.Vertex.Equals(next.Vertex)
		})

		// Move the vertex from the previous location to the new location in the clone.
		//O(n)
		clone.DetachVertex(next.Vertex)
		clone.AttachVertex(next)

		return clone
	}
}

func (c *HeapableCircuitMinClones) Delete() {
	for k := range c.unattachedVertices {
		delete(c.unattachedVertices, k)
	}
	for k := range c.convexVertices {
		delete(c.convexVertices, k)
	}
	if c.closestEdges != nil {
		c.closestEdges.Delete()
	}
	c.Vertices = nil
	c.circuit = nil
	c.circuitEdges = nil
	c.closestEdges = nil
}

func (c *HeapableCircuitMinClones) GetAttachedVertices() []CircuitVertex {
	return c.circuit
}

func (c *HeapableCircuitMinClones) GetAttachedEdges() []CircuitEdge {
	return c.circuitEdges
}

func (c *HeapableCircuitMinClones) GetConvexVertices() map[CircuitVertex]bool {
	return c.convexVertices
}

func (c *HeapableCircuitMinClones) GetClosestEdges() *Heap {
	return c.closestEdges
}

func (c *HeapableCircuitMinClones) GetLength() float64 {
	return c.length
}

func (c *HeapableCircuitMinClones) GetLengthWithNext() float64 {
	if next := c.closestEdges.Peek(); next != nil {
		nextDistToEdge := next.(*DistanceToEdge)
		nextDist := nextDistToEdge.Distance
		if len(c.unattachedVertices) == 0 && nextDist > 0 {
			return c.length // If the circuit is complete and the next vertex to attach increases the perimeter length, the circuit is optimal.
		} else {
			return c.length + nextDist
		}
	} else {
		return c.length
	}
}

func (c *HeapableCircuitMinClones) GetUnattachedVertices() map[CircuitVertex]bool {
	return c.unattachedVertices
}

func (c *HeapableCircuitMinClones) Prepare() {
	c.Vertices = c.deduplicator(c.Vertices)
	c.circuit = []CircuitVertex{}
	c.circuitEdges = []CircuitEdge{}
	c.closestEdges = NewHeap(GetDistanceToEdgeForHeap)
	c.length = 0.0
	c.unattachedVertices = make(map[CircuitVertex]bool)
	c.convexVertices = make(map[CircuitVertex]bool)
}

func (c *HeapableCircuitMinClones) AttachVertex(toAttach *DistanceToEdge) {
	// 1. Update the circuitEdges and retrieve the newly created edges
	var edgeIndex int
	c.circuitEdges, edgeIndex = SplitEdge2(c.circuitEdges, toAttach.Edge, toAttach.Vertex)
	if edgeIndex < 0 {
		expectedEdgeJson, _ := json.Marshal(toAttach.Edge)
		actualCircuitJson, _ := json.Marshal(c.circuitEdges)
		initialVertices, _ := json.Marshal(c.Vertices)
		panic(fmt.Errorf("edge not found in circuit=%p, expected=%s, \ncircuit=%s \nvertices=%s", c, string(expectedEdgeJson), string(actualCircuitJson), string(initialVertices)))
	}
	edgeA, edgeB := c.circuitEdges[edgeIndex], c.circuitEdges[edgeIndex+1]

	// 2. Update the circuit
	// Complexity is O(n), due to shifting array elements to insert new vertex.
	c.circuit = InsertVertex(c.circuit, edgeIndex+1, toAttach.Vertex)
	delete(c.unattachedVertices, toAttach.Vertex)

	// 3. Update the circuit length and the distances increases as a result of the attached vertex.
	c.length += edgeA.GetLength() + edgeB.GetLength() - toAttach.Edge.GetLength()

	// 4. Replace any references to the merged edge with two entries for the newly created edges..
	//    Complexity is O(n)
	distanceIncreases := c.getDistanceIncreases()
	c.closestEdges.ReplaceAll(func(x interface{}) []interface{} {
		current := x.(*DistanceToEdge)
		existingIncrease, okay := distanceIncreases[current.Vertex]
		if !okay {
			existingIncrease = 0.0
		}
		if current.Edge.Equals(toAttach.Edge) {
			if current.Vertex.Equals(toAttach.Vertex) {
				return []interface{}{}
			}
			return []interface{}{
				&DistanceToEdge{
					Vertex:   current.Vertex,
					Edge:     edgeA,
					Distance: edgeA.DistanceIncrease(current.Vertex) - existingIncrease,
				},
				&DistanceToEdge{
					Vertex:   current.Vertex,
					Edge:     edgeB,
					Distance: edgeB.DistanceIncrease(current.Vertex) - existingIncrease,
				},
			}
		} else if current.Vertex.Equals(toAttach.Vertex) || current.Vertex.Equals(toAttach.Edge.GetStart()) || current.Vertex.Equals(toAttach.Edge.GetEnd()) {
			return []interface{}{
				&DistanceToEdge{
					Vertex:   current.Vertex,
					Edge:     current.Edge,
					Distance: current.Edge.DistanceIncrease(current.Vertex) - existingIncrease,
				},
			}
		} else {
			return []interface{}{x}
		}
	})
}

func (c *HeapableCircuitMinClones) DetachVertex(toDetach CircuitVertex) {
	// 1. Remove the vertex from the circuit
	c.circuit = DeleteVertex2(c.circuit, toDetach)

	// 2. Remove the edge with the vertex from the circuitEdges
	var detachedEdgeA, detachedEdgeB, mergedEdge CircuitEdge
	c.circuitEdges, detachedEdgeA, detachedEdgeB, mergedEdge = MergeEdges2(c.circuitEdges, toDetach)

	// 4. Update the circuit distance by removing the old distance increase and adding in the new distance increase.
	c.unattachedVertices[toDetach] = true
	c.length += mergedEdge.GetLength() - detachedEdgeA.GetLength() - detachedEdgeB.GetLength()

	// 5. Replace any references to the merged edges in the heap with a single entry for the merged edge.
	//    Complexity is O(n)
	distanceIncreases := c.getDistanceIncreases()
	replacedVertices := make(map[CircuitVertex]bool)
	c.closestEdges.ReplaceAll(func(x interface{}) []interface{} {
		current := x.(*DistanceToEdge)
		existingIncrease, okay := distanceIncreases[current.Vertex]
		if !okay {
			existingIncrease = 0.0
		}
		if current.Edge.Equals(detachedEdgeA) || current.Edge.Equals(detachedEdgeB) {
			if replacedVertices[current.Vertex] {
				// Only create one entry for the merged edge, even if the source heap has two entries.
				return []interface{}{}
			}
			replacedVertices[current.Vertex] = true
			// Do not allow an entry to be created for either of the vertices of the merged edge.
			// For example, if point B is detached from between points A & C, point C could have an existing entry for A-B, which would be replaced by A-C.
			// The way that this scenario happens is that B and C are both internal points, B is attached first, C is attached between B and D, leaving an entry for A-B for vertex C.
			if current.Vertex.Equals(mergedEdge.GetStart()) || current.Vertex.Equals(mergedEdge.GetEnd()) {
				return []interface{}{}
			}
			return []interface{}{&DistanceToEdge{
				Vertex:   current.Vertex,
				Edge:     mergedEdge,
				Distance: mergedEdge.DistanceIncrease(current.Vertex) - existingIncrease,
			}}
		} else if current.Vertex.Equals(toDetach) || current.Vertex.Equals(mergedEdge.GetStart()) || current.Vertex.Equals(mergedEdge.GetEnd()) {
			return []interface{}{
				&DistanceToEdge{
					Vertex:   current.Vertex,
					Edge:     current.Edge,
					Distance: current.Edge.DistanceIncrease(current.Vertex) - existingIncrease,
				},
			}
		} else {
			return []interface{}{x}
		}
	})
}

func (c *HeapableCircuitMinClones) getDistanceIncreases() map[CircuitVertex]float64 {
	distanceIncreases := make(map[CircuitVertex]float64)
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

var _ HeapableCircuit = (*HeapableCircuitMinClones)(nil)
