package model2d

import (
	"container/heap"
	"fmt"

	"github.com/fealos/lee-tsp-go/model"
)

type HeapableCircuit2D struct {
	Vertices           []model.CircuitVertex
	circuit            []model.CircuitVertex
	circuitEdges       []model.CircuitEdge
	closestEdges       *model.Heap
	length             float64
	midpoint           *Vertex2D
	unattachedVertices map[model.CircuitVertex]bool
}

type heapDistanceToEdge struct {
	vertex   *Vertex2D
	edge     model.CircuitEdge
	distance float64
}

func (h *heapDistanceToEdge) hasVertex(x interface{}) bool {
	dist := x.(*heapDistanceToEdge)
	return dist.vertex == h.vertex
}

func (h *heapDistanceToEdge) ToString() string {
	return fmt.Sprintf(`{"vertex":{"x":%v,"y":%v},"edge":{"start":{"x":%v,"y":%v},"end":{"x":%v,"y":%v}},"distance":%v}`,
		h.vertex.X, h.vertex.Y,
		h.edge.GetStart().(*Vertex2D).X, h.edge.GetStart().(*Vertex2D).Y,
		h.edge.GetEnd().(*Vertex2D).X, h.edge.GetEnd().(*Vertex2D).Y,
		h.distance)
}

func getDistanceToEdgeForHeap(a interface{}) float64 {
	return a.(*heapDistanceToEdge).distance
}

func (c *HeapableCircuit2D) BuildPerimiter() {
	// 1. Find point farthest from midpoint
	// Restricts problem-space to a circle around the midpoint, with radius equal to the distance to the point.
	farthestFromMid := findFarthestPoint(c.midpoint, c.Vertices).(*Vertex2D)
	delete(c.unattachedVertices, farthestFromMid)
	c.circuit = append(c.circuit, farthestFromMid)

	// 2. Find point farthest from point in step 1.
	// Restricts problem-space to intersection of step 1 circle,
	// and a circle centered on the point from step 1 with a radius equal to the distance between the points found in step 1 and 2.
	farthestFromFarthest := findFarthestPoint(farthestFromMid, c.Vertices).(*Vertex2D)
	delete(c.unattachedVertices, farthestFromFarthest)
	c.circuit = append(c.circuit, farthestFromFarthest)

	// 3. Created edges 1 -> 2 and 2 -> 1
	c.circuitEdges = append(c.circuitEdges, NewEdge2D(farthestFromMid, farthestFromFarthest))
	c.circuitEdges = append(c.circuitEdges, NewEdge2D(farthestFromFarthest, farthestFromMid))

	c.length = c.circuitEdges[0].GetLength() * 2

	// 4. Initialize the closestEdges map which will be used to find the exterior point farthest from its closest edge.
	// For the third point only, we can simplify this since both edges are the same (but flipped).
	// When the third point is inserted it will determine whether our vertices are ordered clockwise or counter-clockwise.
	// For this algorithm we will use counter-clockwise ordering, meaning the exterior points will be to the right of their closest edge (while the perimeter is convex).

	exteriorClosestEdges := make(map[model.CircuitVertex]*heapDistanceToEdge)

	for vertex := range c.unattachedVertices {
		v2d := vertex.(*Vertex2D)
		e2d := c.circuitEdges[0].(*Edge2D)
		if v2d.isLeftOf(e2d) {
			e2d = c.circuitEdges[1].(*Edge2D)
		}

		exteriorClosestEdges[vertex] = &heapDistanceToEdge{
			edge:     e2d,
			distance: v2d.distanceToEdge(e2d),
			vertex:   v2d,
		}
	}

	// 5. Find the exterior point farthest from its closest edge.
	// Split the closest edge by adding the point to it, and consequently to the perimeter.
	// Check all remaining exterior points to see if they are now interior points, and update the model as appropriate.
	// Repeat until all points are interior or perimeter points.
	// Complexity: This step in O(N^2) because it iterates once per vertex in the concave perimeter (N iterations) and for each of those iterations it:
	//             1. looks at each exterior point to find farthest from its closest point (O(N)); and then
	//             2. updates each exterior point that had the split edge as its closest edge (O(N)).
	for len(exteriorClosestEdges) > 0 {
		farthestFromClosestEdge := &heapDistanceToEdge{
			distance: 0.0,
		}
		for _, closest := range exteriorClosestEdges {
			if closest.distance > farthestFromClosestEdge.distance {
				farthestFromClosestEdge = closest
			}
		}

		var edgeIndex int
		c.circuitEdges, edgeIndex = model.SplitEdge(c.circuitEdges, farthestFromClosestEdge.edge, farthestFromClosestEdge.vertex)
		c.insertVertex(edgeIndex+1, farthestFromClosestEdge.vertex)
		delete(c.unattachedVertices, farthestFromClosestEdge.vertex)
		delete(exteriorClosestEdges, farthestFromClosestEdge.vertex)

		edgeA, edgeB := c.circuitEdges[edgeIndex], c.circuitEdges[edgeIndex+1]
		c.length += edgeA.GetLength() + edgeB.GetLength() - farthestFromClosestEdge.edge.GetLength()

		for v := range c.unattachedVertices {
			// If the vertex was previously an exterior point and the edge closest to it was split, it could now be an interior point.
			if closest, wasExterior := exteriorClosestEdges[v]; wasExterior && closest.edge == farthestFromClosestEdge.edge {
				var newClosest *heapDistanceToEdge
				if distA, distB := edgeA.DistanceIncrease(v), edgeB.DistanceIncrease(v); distA < distB {
					newClosest = &heapDistanceToEdge{
						edge:     edgeA,
						distance: distA,
						vertex:   v.(*Vertex2D),
					}
				} else {
					newClosest = &heapDistanceToEdge{
						edge:     edgeB,
						distance: distB,
						vertex:   v.(*Vertex2D),
					}
				}

				// If the vertex is now interior, stop tracking its closest edge (until the convex perimeter is fully constructed) and add it to the interior edge list.
				// Otherwise, it is still exterior, so update its closest edge.
				if v.(*Vertex2D).isLeftOf(newClosest.edge.(*Edge2D)) {
					delete(exteriorClosestEdges, v)
				} else {
					exteriorClosestEdges[v] = newClosest
				}
			}
		}
	}

	// 6. Find the closest edge for all interior points, based on distance increase (rather than perpendicular distance)
	for v := range c.unattachedVertices {
		for _, edge := range c.circuitEdges {
			heap.Push(c.closestEdges, &heapDistanceToEdge{
				edge:     edge,
				distance: edge.DistanceIncrease(v),
				vertex:   v.(*Vertex2D),
			})
		}
	}
}

func (c *HeapableCircuit2D) CloneAndUpdate() model.HeapableCircuit {
	// 1. Remove 'next closest' from heap
	next, okay := c.closestEdges.PopHeap().(*heapDistanceToEdge)

	if next == nil || !okay {
		return nil
	} else if len(c.unattachedVertices) == 1 || !c.closestEdges.AnyMatch(next.hasVertex) {
		// 2a. If there are no more items in the heap with the vertex in 'next closest', this is the last edge that the vertex can attach to (i.e. all other possibilities have been tried).
		//     Similarly, if there is only one unattached vertex left, it has to attach to 'next' since that is its most optimal placement in the circuit.
		//     So do not clone the circuit, rather attach 'next closest' to this circuit.
		c.attachVertex(next)
		return nil
	} else {
		// 2b. If there are more items in the heap with the vertex from 'next closest',
		//     clone the circuit so that in the clone the vertex will be attached to the 'next closest' edge,
		//     but in the original circuit that vertex can attach to a different edge.
		clone := &HeapableCircuit2D{
			Vertices:           append(make([]model.CircuitVertex, 0, len(c.Vertices)), c.Vertices...),
			circuit:            append(make([]model.CircuitVertex, 0, len(c.circuit)), c.circuit...),
			circuitEdges:       append(make([]model.CircuitEdge, 0, len(c.circuitEdges)), c.circuitEdges...),
			closestEdges:       c.closestEdges.Clone(),
			length:             c.length,
			midpoint:           c.midpoint,
			unattachedVertices: make(map[model.CircuitVertex]bool),
		}
		for k, v := range c.unattachedVertices {
			clone.unattachedVertices[k] = v
		}

		clone.attachVertex(next)
		return clone
	}
}

func (c *HeapableCircuit2D) Delete() {
	for k := range c.unattachedVertices {
		delete(c.unattachedVertices, k)
	}
	if c.closestEdges != nil {
		c.closestEdges.Delete()
	}
	c.Vertices = nil
	c.circuit = nil
	c.circuitEdges = nil
	c.midpoint = nil
	c.closestEdges = nil
}

func (c *HeapableCircuit2D) GetAttachedVertices() []model.CircuitVertex {
	return c.circuit
}

func (c *HeapableCircuit2D) GetLength() float64 {
	return c.length
}

func (c *HeapableCircuit2D) GetLengthWithNext() float64 {
	if next := c.closestEdges.Peek(); next != nil {
		return c.length + next.(*heapDistanceToEdge).distance
	} else {
		return c.length
	}
}

func (c *HeapableCircuit2D) GetUnattachedVertices() map[model.CircuitVertex]bool {
	return c.unattachedVertices
}

func (c *HeapableCircuit2D) Prepare() {
	c.unattachedVertices = make(map[model.CircuitVertex]bool)
	c.closestEdges = model.NewHeap(getDistanceToEdgeForHeap)
	c.circuit = []model.CircuitVertex{}
	c.circuitEdges = []model.CircuitEdge{}
	c.length = 0.0

	c.Vertices = DeduplicateVertices(c.Vertices)

	numVertices := float64(len(c.Vertices))
	c.midpoint = &Vertex2D{0.0, 0.0}

	for _, v := range c.Vertices {
		c.unattachedVertices[v] = true
		v2d := v.(*Vertex2D)
		c.midpoint.X += v2d.X / numVertices
		c.midpoint.Y += v2d.Y / numVertices
	}
}

func (c *HeapableCircuit2D) attachVertex(toAttach *heapDistanceToEdge) {
	// 1. Update the circuitEdges
	var edgeIndex int
	c.circuitEdges, edgeIndex = model.SplitEdge(c.circuitEdges, toAttach.edge, toAttach.vertex)

	// 2. Update the circuit
	c.insertVertex(edgeIndex+1, toAttach.vertex)
	delete(c.unattachedVertices, toAttach.vertex)
	c.length += toAttach.distance

	// 3. Retrieve the newly created edges
	edgeA, edgeB := c.circuitEdges[edgeIndex], c.circuitEdges[edgeIndex+1%len(c.circuitEdges)]

	c.closestEdges.ReplaceAll(func(x interface{}) []interface{} {
		current := x.(*heapDistanceToEdge)
		if current.vertex == toAttach.vertex {
			// 4a. Remove any items in the heap with the vertex that was attached to the circuit.
			return []interface{}{}
		} else if current.edge == toAttach.edge {
			// 4b. Replace any items in the heap with the edge that was split, with one entry for each of the newly created edges.
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
			// 4c. Ignore any items in the heap that do not contain the attached vertex nor the split edge.
			return []interface{}{current}
		}
	})
}

func (c *HeapableCircuit2D) insertVertex(index int, vertex model.CircuitVertex) {
	if index >= len(c.circuit) {
		c.circuit = append(c.circuit, vertex)
	} else {
		// copy all elements starting at the index one to the right to create a duplicate record at index and index+1.
		c.circuit = append(c.circuit[:index+1], c.circuit[index:]...)
		// update only the vertex at the index, so that there are no duplicates and the vertex is at the index.
		c.circuit[index] = vertex
	}
}

var _ model.HeapableCircuit = (*HeapableCircuit2D)(nil)
var _ model.Printable = (*heapDistanceToEdge)(nil)
