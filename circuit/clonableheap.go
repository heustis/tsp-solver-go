package circuit

import "github.com/fealos/lee-tsp-go/model"

type HeapableCircuit struct {
	Vertices           []model.CircuitVertex
	deduplicator       func([]model.CircuitVertex) []model.CircuitVertex
	perimeterBuilder   model.PerimeterBuilder
	circuitEdges       []model.CircuitEdge
	closestEdges       *model.Heap
	length             float64
	unattachedVertices map[model.CircuitVertex]bool
}

func NewHeapableCircuit(vertices []model.CircuitVertex, deduplicator model.Deduplicator, perimeterBuilder model.PerimeterBuilder) *HeapableCircuit {
	return &HeapableCircuit{
		Vertices:         vertices,
		deduplicator:     deduplicator,
		perimeterBuilder: perimeterBuilder,
	}
}

func (c *HeapableCircuit) BuildPerimiter() {
	c.circuitEdges, c.unattachedVertices = c.perimeterBuilder(c.Vertices)

	// Determine the initial length of the perimeter.
	c.length = 0.0
	for _, edge := range c.circuitEdges {
		c.length += edge.GetLength()
	}

	// Find the closest edge for all interior points, based on distance increase (rather than perpendicular distance)
	// O(v*e) due to storing the entries, one per edge per unattached vertex, in an array then heapifing the array.
	// Note: it would be O(v*e * log(v*e)) to push each entry independently.
	initialCandidates := []interface{}{}
	for v := range c.unattachedVertices {
		for _, edge := range c.circuitEdges {
			initialCandidates = append(initialCandidates, &model.DistanceToEdge{
				Edge:     edge,
				Distance: edge.DistanceIncrease(v),
				Vertex:   v,
			})
		}
	}
	c.closestEdges.PushAll(initialCandidates...)
}

func (c *HeapableCircuit) CloneAndUpdate() ClonableCircuit {
	// 1. Remove 'next closest' from heap
	next, okay := c.closestEdges.PopHeap().(*model.DistanceToEdge)

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
		clone := &HeapableCircuit{
			Vertices:           append(make([]model.CircuitVertex, 0, len(c.Vertices)), c.Vertices...),
			circuitEdges:       append(make([]model.CircuitEdge, 0, len(c.circuitEdges)), c.circuitEdges...),
			closestEdges:       c.closestEdges.Clone(),
			length:             c.length,
			unattachedVertices: make(map[model.CircuitVertex]bool),
		}
		for k, v := range c.unattachedVertices {
			clone.unattachedVertices[k] = v
		}

		clone.attachVertex(next)
		return clone
	}
}

func (c *HeapableCircuit) Delete() {
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

func (c *HeapableCircuit) FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge) {
	if next, okay := c.closestEdges.Peek().(*model.DistanceToEdge); okay && next != nil {
		return next.Vertex, next.Edge
	}
	return nil, nil
}

func (c *HeapableCircuit) GetAttachedVertices() []model.CircuitVertex {
	vertices := make([]model.CircuitVertex, len(c.circuitEdges))
	for i, edge := range c.circuitEdges {
		vertices[i] = edge.GetStart()
	}
	return vertices
}

func (c *HeapableCircuit) GetAttachedEdges() []model.CircuitEdge {
	return c.circuitEdges
}

func (c *HeapableCircuit) GetClosestEdges() *model.Heap {
	return c.closestEdges
}

func (c *HeapableCircuit) GetLength() float64 {
	return c.length
}

func (c *HeapableCircuit) GetLengthWithNext() float64 {
	if next := c.closestEdges.Peek(); next != nil {
		return c.length + next.(*model.DistanceToEdge).Distance
	} else {
		return c.length
	}
}

func (c *HeapableCircuit) GetUnattachedVertices() map[model.CircuitVertex]bool {
	return c.unattachedVertices
}

func (c *HeapableCircuit) Prepare() {
	c.Vertices = c.deduplicator(c.Vertices)
	c.circuitEdges = []model.CircuitEdge{}
	c.closestEdges = model.NewHeap(model.GetDistanceToEdgeForHeap)
	c.length = 0.0
	c.unattachedVertices = make(map[model.CircuitVertex]bool)
}

func (c *HeapableCircuit) attachVertex(toAttach *model.DistanceToEdge) {
	// 1. Update the circuitEdges
	var edgeIndex int
	c.circuitEdges, edgeIndex = model.SplitEdge(c.circuitEdges, toAttach.Edge, toAttach.Vertex)

	// 2. Update the circuit
	delete(c.unattachedVertices, toAttach.Vertex)
	c.length += toAttach.Distance

	// 3. Retrieve the newly created edges
	edgeA, edgeB := c.circuitEdges[edgeIndex], c.circuitEdges[edgeIndex+1%len(c.circuitEdges)]

	c.closestEdges.ReplaceAll(func(x interface{}) interface{} {
		current := x.(*model.DistanceToEdge)
		if current.Vertex == toAttach.Vertex {
			// 4a. Remove any items in the heap with the vertex that was attached to the circuit.
			return nil
		} else if current.Edge == toAttach.Edge {
			// 4b. Replace any items in the heap with the edge that was split, with one entry for each of the newly created edges.
			return []interface{}{
				&model.DistanceToEdge{
					Vertex:   current.Vertex,
					Edge:     edgeA,
					Distance: edgeA.DistanceIncrease(current.Vertex),
				},
				&model.DistanceToEdge{
					Vertex:   current.Vertex,
					Edge:     edgeB,
					Distance: edgeB.DistanceIncrease(current.Vertex),
				},
			}
		} else {
			// 4c. Ignore any items in the heap that do not contain the attached vertex nor the split edge.
			return current
		}
	})
}

var _ ClonableCircuit = (*HeapableCircuit)(nil)
