package model

import (
	"encoding/json"
	"fmt"
)

type HeapableCircuitMinClonesLimited struct {
	Vertices         []CircuitVertex
	deduplicator     func([]CircuitVertex) []CircuitVertex
	perimeterBuilder PerimeterBuilder
	circuitEdges     []CircuitEdge
	closestEdges     *Heap
	length           float64
	interiorVertices map[CircuitVertex]*vertexStatus
}

func NewHeapableCircuitMinClonesLimited(vertices []CircuitVertex, deduplicator func([]CircuitVertex) []CircuitVertex, perimeterBuilder PerimeterBuilder) *HeapableCircuitMinClonesLimited {
	return &HeapableCircuitMinClonesLimited{
		Vertices:         vertices,
		deduplicator:     deduplicator,
		perimeterBuilder: perimeterBuilder,
	}
}

func (c *HeapableCircuitMinClonesLimited) BuildPerimiter() {
	var unattachedVertices map[CircuitVertex]bool
	c.circuitEdges, unattachedVertices = c.perimeterBuilder.BuildPerimiter(c.Vertices)

	// Determine the initial length of the perimeter.
	c.length = 0.0
	for _, edge := range c.circuitEdges {
		c.length += edge.GetLength()
		c.interiorVertices[edge.GetStart()] = &vertexStatus{
			isUnattached:     false,
			isConcave:        true,
			distanceIncrease: 0.0,
		}
	}

	// Find the closest edge for all interior points, based on distance increase (rather than perpendicular distance)
	// total vertices = attached + unattached
	// complexity  = attached * unattached  = attached * (total - attached)  = total*attached - attached^2
	c.closestEdges = NewHeap(GetDistanceToEdgeForHeap)
	for v := range unattachedVertices {
		vertexHeap := NewHeap(GetDistanceToEdgeForHeap)
		c.interiorVertices[v] = &vertexStatus{
			isUnattached:     true,
			isConcave:        false,
			distanceIncrease: 0.0,
		}
		for _, edge := range c.circuitEdges {
			// Note: Using Push, not PushHeap to just append elements for now, will heapify after all elements are pushed.
			vertexHeap.Push(&DistanceToEdge{
				Vertex:   v,
				Edge:     edge,
				Distance: edge.DistanceIncrease(v),
			})
		}
		vertexHeap.Heapify()
		vertexHeap.TrimN(3)
		for v := vertexHeap.PopHeap(); v != nil; v = vertexHeap.PopHeap() {
			c.closestEdges.Push(v)
		}
	}
	c.closestEdges.Heapify()
}

func (c *HeapableCircuitMinClonesLimited) CloneAndUpdate() HeapableCircuit {
	// 1. Remove 'next closest' from heap - complexity O(log n)
	next, okay := c.closestEdges.PopHeap().(*DistanceToEdge)

	if next == nil || !okay {
		return nil
	} else if c.interiorVertices[next.Vertex].isUnattached || !c.closestEdges.AnyMatch(next.HasVertex) {
		// 2a. If there are no more items in the heap with the vertex in 'next closest', this is the last edge that the vertex can attach to (i.e. all other possibilities have been tried).
		//     Similarly, if this is the first time we are encountering a vertex, attach it at the specified location without cloning.
		//     So do not clone the circuit, rather attach 'next closest' to this circuit.
		// AnyMatch is O(n), attachVertex is O(n)
		c.AttachVertex(next)
		return nil
	} else {
		// 2b. If the 'next closest' vertex is already attached, clone the circuit with the 'next closest' vertex attached to the 'next closest' edge.
		// O(n)
		clone := &HeapableCircuitMinClonesLimited{
			Vertices:         c.Vertices,
			circuitEdges:     make([]CircuitEdge, len(c.circuitEdges)),
			closestEdges:     c.closestEdges.Clone(),
			length:           c.length,
			interiorVertices: make(map[CircuitVertex]*vertexStatus),
		}
		copy(clone.circuitEdges, c.circuitEdges)

		for k, v := range c.interiorVertices {
			// We clone on write, so we don't need to create copies at this time.
			clone.interiorVertices[k] = v
		}

		// Move the vertex from the previous location to the new location in the clone.
		//O(n)
		clone.MoveVertex(next)

		return clone
	}
}

func (c *HeapableCircuitMinClonesLimited) Delete() {
	for k := range c.interiorVertices {
		delete(c.interiorVertices, k)
	}
	if c.closestEdges != nil {
		c.closestEdges.Delete()
	}
	c.Vertices = nil
	c.circuitEdges = nil
	c.closestEdges = nil
}

func (c *HeapableCircuitMinClonesLimited) GetAttachedVertices() []CircuitVertex {
	vertices := make([]CircuitVertex, len(c.circuitEdges))
	for i, edge := range c.circuitEdges {
		vertices[i] = edge.GetStart()
	}
	return vertices
}

func (c *HeapableCircuitMinClonesLimited) GetAttachedEdges() []CircuitEdge {
	return c.circuitEdges
}

func (c *HeapableCircuitMinClonesLimited) GetClosestEdges() *Heap {
	return c.closestEdges
}

func (c *HeapableCircuitMinClonesLimited) GetLength() float64 {
	return c.length
}

func (c *HeapableCircuitMinClonesLimited) GetLengthWithNext() float64 {
	if next := c.closestEdges.Peek(); next != nil {
		nextDistToEdge := next.(*DistanceToEdge)
		if len(c.circuitEdges) == len(c.Vertices) && nextDistToEdge.Distance > 0 {
			return c.length // If the circuit is complete and the next vertex to attach increases the perimeter length, the circuit is optimal.
		} else {
			return c.length + nextDistToEdge.Distance
		}
	} else {
		return c.length
	}
}

func (c *HeapableCircuitMinClonesLimited) GetUnattachedVertices() map[CircuitVertex]bool {
	unattachedVertices := make(map[CircuitVertex]bool)
	for k, v := range c.interiorVertices {
		if v.isUnattached {
			unattachedVertices[k] = true
		}
	}
	return unattachedVertices
}

func (c *HeapableCircuitMinClonesLimited) Prepare() {
	c.Vertices = c.deduplicator(c.Vertices)
	c.circuitEdges = []CircuitEdge{}
	c.closestEdges = NewHeap(GetDistanceToEdgeForHeap)
	c.length = 0.0
	c.interiorVertices = make(map[CircuitVertex]*vertexStatus)
}

func (c *HeapableCircuitMinClonesLimited) AttachVertex(toAttach *DistanceToEdge) {
	// 1. Update the circuitEdges and retrieve the newly created edges
	var edgeIndex int
	//TODO - this can cause an index out of bounds exception, investigate prior to using this struct further.
	// Note: in preliminary tests this was already less accurate than the greedy algorithms.
	c.circuitEdges, edgeIndex = SplitEdge2(c.circuitEdges, toAttach.Edge, toAttach.Vertex)
	if edgeIndex < 0 {
		expectedEdgeJson, _ := json.Marshal(toAttach.Edge)
		actualCircuitJson, _ := json.Marshal(c.circuitEdges)
		initialVertices, _ := json.Marshal(c.Vertices)
		panic(fmt.Errorf("edge not found in circuit=%p, expected=%s, \ncircuit=%s \nvertices=%s", c, string(expectedEdgeJson), string(actualCircuitJson), string(initialVertices)))
	}
	edgeA, edgeB := c.circuitEdges[edgeIndex], c.circuitEdges[edgeIndex+1]

	// 2. Update the circuit length and the distances increases as a result of the attached vertex.
	//    Note - the DistanceToEdge already accounts for both the existing edge and the new edge.
	c.length += toAttach.Distance

	updatedVertices := make(map[CircuitVertex]bool)
	updatedVertices[toAttach.Vertex] = true
	updatedVertices[toAttach.Edge.GetStart()] = true
	updatedVertices[toAttach.Edge.GetEnd()] = true
	c.updateDistanceIncreases(updatedVertices)

	// 3. Replace any references to the merged edge with two entries for the newly created edges..
	//    Complexity is O(n)
	c.closestEdges.ReplaceAll2(func(x interface{}) interface{} {
		current := x.(*DistanceToEdge)
		existingIncrease := c.interiorVertices[current.Vertex].distanceIncrease
		if current.Edge.GetStart() == toAttach.Edge.GetStart() && current.Edge.GetEnd() == toAttach.Edge.GetEnd() {
			if current.Vertex == toAttach.Vertex {
				return nil
			} else if distA, distB := edgeA.DistanceIncrease(current.Vertex), edgeB.DistanceIncrease(current.Vertex); distA <= distB {
				return &DistanceToEdge{
					Vertex:   current.Vertex,
					Edge:     edgeA,
					Distance: distA - existingIncrease,
				}
			} else {
				return &DistanceToEdge{
					Vertex:   current.Vertex,
					Edge:     edgeB,
					Distance: distB - existingIncrease,
				}
			}
		} else if updatedVertices[current.Vertex] {
			return &DistanceToEdge{
				Vertex:   current.Vertex,
				Edge:     current.Edge,
				Distance: current.Edge.DistanceIncrease(current.Vertex) - existingIncrease,
			}
		} else {
			return x
		}
	})
}

func (c *HeapableCircuitMinClonesLimited) MoveVertex(toMove *DistanceToEdge) {
	// 1. Remove the edge with the vertex from the circuitEdges
	var mergedEdge, splitEdgeA, splitEdgeB CircuitEdge
	c.circuitEdges, mergedEdge, splitEdgeA, splitEdgeB = MoveVertex(c.circuitEdges, toMove.Vertex, toMove.Edge)
	if mergedEdge == nil {
		toMoveJson, _ := json.Marshal(toMove.Vertex)
		targetEdgeJson, _ := json.Marshal(toMove.Edge)
		actualCircuitJson, _ := json.Marshal(c.circuitEdges)
		panic(fmt.Errorf("cannot move vertex to edge circuit=%p, vertex=%s, edge=%s circuit=%s", c, string(toMoveJson), string(targetEdgeJson), string(actualCircuitJson)))
	}

	// 2. Update the circuit distance and the distances increases as a result of the attached vertex.
	//    Note - the DistanceToEdge already accounts for both the existing edge and the new edge.
	c.length += toMove.Distance

	updatedVertices := make(map[CircuitVertex]bool)
	updatedVertices[toMove.Vertex] = true
	updatedVertices[toMove.Edge.GetStart()] = true
	updatedVertices[toMove.Edge.GetEnd()] = true
	updatedVertices[mergedEdge.GetStart()] = true
	updatedVertices[mergedEdge.GetEnd()] = true
	c.updateDistanceIncreases(updatedVertices)

	// 3. Replace any references to the merged edges in the heap with a single entry for the merged edge.
	//    Complexity is O(n)
	replacedVertices := make(map[CircuitVertex]bool)
	c.closestEdges.ReplaceAll2(func(x interface{}) interface{} {
		current := x.(*DistanceToEdge)
		existingIncrease := c.interiorVertices[current.Vertex].distanceIncrease
		if current.Vertex == toMove.Vertex {
			return nil
		} else if current.Edge.GetStart() == toMove.Vertex || current.Edge.GetEnd() == toMove.Vertex {
			if replacedVertices[current.Vertex] {
				// Only create one entry for the merged edge, even if the source heap has two entries.
				return nil
			}
			replacedVertices[current.Vertex] = true
			// Do not allow an entry to be created for either of the vertices of the merged edge.
			// For example, if point B is detached from between points A & C, point C could have an existing entry for A-B, which would be replaced by A-C.
			// The way that this scenario happens is that B and C are both internal points, B is attached first, C is attached between B and D, leaving an entry for A-B for vertex C.
			if current.Vertex == mergedEdge.GetStart() || current.Vertex == mergedEdge.GetEnd() {
				return nil
			}
			return &DistanceToEdge{
				Vertex:   current.Vertex,
				Edge:     mergedEdge,
				Distance: mergedEdge.DistanceIncrease(current.Vertex) - existingIncrease,
			}
		} else if current.Edge.GetStart() == toMove.Edge.GetStart() && current.Edge.GetEnd() == toMove.Edge.GetEnd() {
			return []interface{}{
				&DistanceToEdge{
					Vertex:   current.Vertex,
					Edge:     splitEdgeA,
					Distance: splitEdgeA.DistanceIncrease(current.Vertex) - existingIncrease,
				},
				&DistanceToEdge{
					Vertex:   current.Vertex,
					Edge:     splitEdgeB,
					Distance: splitEdgeB.DistanceIncrease(current.Vertex) - existingIncrease,
				},
			}
		} else if updatedVertices[current.Vertex] {
			return &DistanceToEdge{
				Vertex:   current.Vertex,
				Edge:     current.Edge,
				Distance: current.Edge.DistanceIncrease(current.Vertex) - existingIncrease,
			}
		} else {
			return x
		}
	})
}

func (c *HeapableCircuitMinClonesLimited) updateDistanceIncreases(updatedVertices map[CircuitVertex]bool) {
	circuitLen := len(c.circuitEdges)
	if circuitLen >= 3 {
		prev := c.circuitEdges[circuitLen-1]
		for _, edge := range c.circuitEdges {
			v := edge.GetStart()
			if currentMetadata := c.interiorVertices[v]; updatedVertices[v] && !currentMetadata.isConcave {
				c.interiorVertices[v] = &vertexStatus{
					isUnattached:     false,
					isConcave:        currentMetadata.isConcave,
					distanceIncrease: prev.GetLength() + edge.GetLength() - edge.GetEnd().DistanceTo(prev.GetStart()),
				}
			}
			prev = edge
		}
	}
}

var _ HeapableCircuit = (*HeapableCircuitMinClonesLimited)(nil)
