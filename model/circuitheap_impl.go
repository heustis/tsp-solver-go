package model

type HeapableCircuitImpl struct {
	Vertices           []CircuitVertex
	deduplicator       func([]CircuitVertex) []CircuitVertex
	perimeterBuilder   PerimeterBuilder
	circuitEdges       []CircuitEdge
	closestEdges       *Heap
	length             float64
	unattachedVertices map[CircuitVertex]bool
}

func NewHeapableCircuitImpl(vertices []CircuitVertex, deduplicator func([]CircuitVertex) []CircuitVertex, perimeterBuilder PerimeterBuilder) *HeapableCircuitImpl {
	return &HeapableCircuitImpl{
		Vertices:         vertices,
		deduplicator:     deduplicator,
		perimeterBuilder: perimeterBuilder,
	}
}

func (c *HeapableCircuitImpl) BuildPerimiter() {
	c.circuitEdges, c.unattachedVertices = c.perimeterBuilder.BuildPerimiter(c.Vertices)

	// Determine the initial length of the perimeter.
	c.length = 0.0
	for _, edge := range c.circuitEdges {
		c.length += edge.GetLength()
	}

	// 6. Find the closest edge for all interior points, based on distance increase (rather than perpendicular distance)
	// O(v*e) due to storing the entries, one per edge per unattached vertex, in an array then heapifing the array.
	// Note: it would be O(v*e * log(v*e)) to push each entry independently.
	initialCandidates := []interface{}{}
	for v := range c.unattachedVertices {
		for _, edge := range c.circuitEdges {
			initialCandidates = append(initialCandidates, &DistanceToEdge{
				Edge:     edge,
				Distance: edge.DistanceIncrease(v),
				Vertex:   v,
			})
		}
	}
	c.closestEdges.PushAll(initialCandidates...)
}

func (c *HeapableCircuitImpl) CloneAndUpdate() HeapableCircuit {
	// 1. Remove 'next closest' from heap
	next, okay := c.closestEdges.PopHeap().(*DistanceToEdge)

	if next == nil || !okay {
		return nil
	} else if len(c.unattachedVertices) == 1 || !c.closestEdges.AnyMatch(next.HasVertex) {
		// 2a. If there are no more items in the heap with the vertex in 'next closest', this is the last edge that the vertex can attach to (i.e. all other possibilities have been tried).
		//     Similarly, if there is only one unattached vertex left, it has to attach to 'next' since that is its most optimal placement in the circuit.
		//     So do not clone the circuit, rather attach 'next closest' to this circuit.
		c.attachVertex(next)
		return nil
	} else {
		// 2b. If there are more items in the heap with the vertex from 'next closest',
		//     clone the circuit so that in the clone the vertex will be attached to the 'next closest' edge,
		//     but in the original circuit that vertex can attach to a different edge.
		clone := &HeapableCircuitImpl{
			Vertices:           append(make([]CircuitVertex, 0, len(c.Vertices)), c.Vertices...),
			circuitEdges:       append(make([]CircuitEdge, 0, len(c.circuitEdges)), c.circuitEdges...),
			closestEdges:       c.closestEdges.Clone(),
			length:             c.length,
			unattachedVertices: make(map[CircuitVertex]bool),
		}
		for k, v := range c.unattachedVertices {
			clone.unattachedVertices[k] = v
		}

		clone.attachVertex(next)
		return clone
	}
}

func (c *HeapableCircuitImpl) Delete() {
	for k := range c.unattachedVertices {
		delete(c.unattachedVertices, k)
	}
	if c.closestEdges != nil {
		c.closestEdges.Delete()
	}
	c.Vertices = nil
	c.circuitEdges = nil
	c.closestEdges = nil
}

func (c *HeapableCircuitImpl) GetAttachedVertices() []CircuitVertex {
	vertices := make([]CircuitVertex, len(c.circuitEdges))
	for i, edge := range c.circuitEdges {
		vertices[i] = edge.GetStart()
	}
	return vertices
}

func (c *HeapableCircuitImpl) GetAttachedEdges() []CircuitEdge {
	return c.circuitEdges
}

func (c *HeapableCircuitImpl) GetClosestEdges() *Heap {
	return c.closestEdges
}

func (c *HeapableCircuitImpl) GetLength() float64 {
	return c.length
}

func (c *HeapableCircuitImpl) GetLengthWithNext() float64 {
	if next := c.closestEdges.Peek(); next != nil {
		return c.length + next.(*DistanceToEdge).Distance
	} else {
		return c.length
	}
}

func (c *HeapableCircuitImpl) GetUnattachedVertices() map[CircuitVertex]bool {
	return c.unattachedVertices
}

func (c *HeapableCircuitImpl) Prepare() {
	c.Vertices = c.deduplicator(c.Vertices)
	c.circuitEdges = []CircuitEdge{}
	c.closestEdges = NewHeap(GetDistanceToEdgeForHeap)
	c.length = 0.0
	c.unattachedVertices = make(map[CircuitVertex]bool)
}

func (c *HeapableCircuitImpl) attachVertex(toAttach *DistanceToEdge) {
	// 1. Update the circuitEdges
	var edgeIndex int
	c.circuitEdges, edgeIndex = SplitEdge(c.circuitEdges, toAttach.Edge, toAttach.Vertex)

	// 2. Update the circuit
	delete(c.unattachedVertices, toAttach.Vertex)
	c.length += toAttach.Distance

	// 3. Retrieve the newly created edges
	edgeA, edgeB := c.circuitEdges[edgeIndex], c.circuitEdges[edgeIndex+1%len(c.circuitEdges)]

	c.closestEdges.ReplaceAll(func(x interface{}) []interface{} {
		current := x.(*DistanceToEdge)
		if current.Vertex == toAttach.Vertex {
			// 4a. Remove any items in the heap with the vertex that was attached to the circuit.
			return []interface{}{}
		} else if current.Edge == toAttach.Edge {
			// 4b. Replace any items in the heap with the edge that was split, with one entry for each of the newly created edges.
			return []interface{}{
				&DistanceToEdge{
					Vertex:   current.Vertex,
					Edge:     edgeA,
					Distance: edgeA.DistanceIncrease(current.Vertex),
				},
				&DistanceToEdge{
					Vertex:   current.Vertex,
					Edge:     edgeB,
					Distance: edgeB.DistanceIncrease(current.Vertex),
				},
			}
		} else {
			// 4c. Ignore any items in the heap that do not contain the attached vertex nor the split edge.
			return []interface{}{current}
		}
	})
}

var _ HeapableCircuit = (*HeapableCircuitImpl)(nil)
var _ Printable = (*DistanceToEdge)(nil)
