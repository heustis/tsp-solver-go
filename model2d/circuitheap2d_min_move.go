package model2d

import (
	"container/heap"
	"encoding/json"
	"fmt"

	"github.com/fealos/lee-tsp-go/model"
)

type HeapableCircuit2DMinClonesMove struct {
	Vertices           []*Vertex2D
	circuit            []model.CircuitVertex
	circuitEdges       []model.CircuitEdge
	closestEdges       *model.Heap
	length             float64
	midpoint           *Vertex2D
	unattachedVertices map[model.CircuitVertex]bool
	distanceIncreases  map[model.CircuitVertex]float64
}

func (c *HeapableCircuit2DMinClonesMove) BuildPerimiter() {
	// 1. Find point farthest from midpoint
	// Restricts problem-space to a circle around the midpoint, with radius equal to the distance to the point.
	farthestFromMid := findFarthestPoint(c.midpoint, c.Vertices)
	delete(c.unattachedVertices, farthestFromMid)
	c.circuit = append(c.circuit, farthestFromMid)

	// 2. Find point farthest from point in step 1.
	// Restricts problem-space to intersection of step 1 circle,
	// and a circle centered on the point from step 1 with a radius equal to the distance between the points found in step 1 and 2.
	farthestFromFarthest := findFarthestPoint(farthestFromMid, c.Vertices)
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

func (c *HeapableCircuit2DMinClonesMove) CloneAndUpdate() model.HeapableCircuit {
	// 1. Remove 'next closest' from heap
	next, okay := c.closestEdges.PopHeap().(*heapDistanceToEdge)

	if next == nil || !okay {
		return nil
	} else if c.unattachedVertices[next.vertex] || !c.closestEdges.AnyMatch(next.hasVertex) {
		// 2a. If there are no more items in the heap with the vertex in 'next closest', this is the last edge that the vertex can attach to (i.e. all other possibilities have been tried).
		//     Similarly, if this is the first time we are encountering a vertex, attach it at the specified location without cloning.
		//     So do not clone the circuit, rather attach 'next closest' to this circuit.
		c.attachVertex(next)
		return nil
	} else {
		// 2b. If the 'next closest' vertex is already attached, clone the circuit with the 'next closest' vertex attached to the 'next closest' edge.
		clone := &HeapableCircuit2DMinClonesMove{
			Vertices:           append(make([]*Vertex2D, 0, len(c.Vertices)), c.Vertices...),
			circuit:            append(make([]model.CircuitVertex, 0, len(c.circuit)), c.circuit...),
			circuitEdges:       append(make([]model.CircuitEdge, 0, len(c.circuitEdges)), c.circuitEdges...),
			closestEdges:       c.closestEdges.Clone(),
			length:             c.length,
			midpoint:           c.midpoint,
			unattachedVertices: make(map[model.CircuitVertex]bool),
			distanceIncreases:  make(map[model.CircuitVertex]float64),
		}
		for k, v := range c.unattachedVertices {
			clone.unattachedVertices[k] = v
		}
		for k, v := range c.distanceIncreases {
			clone.distanceIncreases[k] = v
		}

		clone.MoveVertex(next)
		return clone
	}
}

func (c *HeapableCircuit2DMinClonesMove) Delete() {
	for k := range c.unattachedVertices {
		delete(c.unattachedVertices, k)
	}
	for k := range c.distanceIncreases {
		delete(c.distanceIncreases, k)
	}
	if c.closestEdges != nil {
		c.closestEdges.DeleteAll(func(x interface{}) bool { return true })
	}
	c.Vertices = nil
	c.circuit = nil
	c.circuitEdges = nil
	c.midpoint = nil
	c.closestEdges = nil
}

func (c *HeapableCircuit2DMinClonesMove) GetAttachedVertices() []model.CircuitVertex {
	return c.circuit
}

func (c *HeapableCircuit2DMinClonesMove) GetLength() float64 {
	return c.length
}

func (c *HeapableCircuit2DMinClonesMove) GetLengthWithNext() float64 {
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

func (c *HeapableCircuit2DMinClonesMove) GetUnattachedVertices() map[model.CircuitVertex]bool {
	return c.unattachedVertices
}

func (c *HeapableCircuit2DMinClonesMove) Prepare() {
	c.unattachedVertices = make(map[model.CircuitVertex]bool)
	c.distanceIncreases = make(map[model.CircuitVertex]float64)
	c.closestEdges = model.NewHeap(getDistanceToEdgeForHeap)
	c.circuit = []model.CircuitVertex{}
	c.circuitEdges = []model.CircuitEdge{}
	c.length = 0.0

	c.Vertices = deduplicateVertices(c.Vertices)

	numVertices := float64(len(c.Vertices))
	c.midpoint = &Vertex2D{0.0, 0.0}

	for _, v := range c.Vertices {
		c.unattachedVertices[v] = true
		c.distanceIncreases[v] = 0.0
		c.midpoint.X += v.X / numVertices
		c.midpoint.Y += v.Y / numVertices
	}
}

func (c *HeapableCircuit2DMinClonesMove) attachVertex(toAttach *heapDistanceToEdge) {
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
	c.insertVertex(edgeIndex+1, toAttach.vertex)
	delete(c.unattachedVertices, toAttach.vertex)

	// 3. Update the circuit length and the distances increases as a result of the attached vertex.
	c.length += toAttach.distance
	c.distanceIncreases[toAttach.vertex] += toAttach.distance
	startDistanceDelta := c.updateDistanceIncrease(edgeIndex)
	endDistanceDelta := c.updateDistanceIncrease(edgeIndex + 2)

	// 4. Replace any items in the heap with the edge that was split, with one entry for each of the newly created edges.
	//    Also, update the distance for each heapDistanceToEdge with a vertex affected by the attachment.
	//    Note: need to create a copy of the heapDistanceToEdge, when updatign distance, so that other versions of the circuit are unaffected.
	c.closestEdges = c.closestEdges.ReplaceAll(func(x interface{}) []interface{} {
		current := x.(*heapDistanceToEdge)
		if current.edge.GetStart() == toAttach.edge.GetStart() && current.edge.GetEnd() == toAttach.edge.GetEnd() {
			return []interface{}{
				&heapDistanceToEdge{
					vertex:   current.vertex,
					edge:     edgeA,
					distance: edgeA.DistanceIncrease(current.vertex) - c.distanceIncreases[current.vertex],
				},
				&heapDistanceToEdge{
					vertex:   current.vertex,
					edge:     edgeB,
					distance: edgeB.DistanceIncrease(current.vertex) - c.distanceIncreases[current.vertex],
				},
			}
		} else if current.vertex == toAttach.vertex {
			return []interface{}{current.cloneAndAdjustDistance(-toAttach.distance)}
		} else if current.vertex == toAttach.edge.GetStart() {
			return []interface{}{current.cloneAndAdjustDistance(-startDistanceDelta)}
		} else if current.vertex == toAttach.edge.GetEnd() {
			return []interface{}{current.cloneAndAdjustDistance(-endDistanceDelta)}
		} else {
			return []interface{}{current}
		}
	})
}

func (c *HeapableCircuit2DMinClonesMove) MoveVertex(newLocation *heapDistanceToEdge) {
	// 1. Remove the vertex from the circuit
	detachIndex := model.IndexOfVertex(c.circuit, newLocation.vertex)
	c.circuit = model.DeleteVertex(c.circuit, detachIndex)

	// 2. Remove the edge with the vertex from the circuitEdges
	var detachedEdgeA model.CircuitEdge
	var detachedEdgeB model.CircuitEdge
	c.circuitEdges, detachedEdgeA, detachedEdgeB = model.MergeEdges(c.circuitEdges, detachIndex)

	// 3. Retrieve the merged edge from the circuit
	edgeLen := len(c.circuitEdges)
	mergedEdge := c.circuitEdges[(detachIndex-1+edgeLen)%edgeLen]

	// 4. Add the new vertex and edges to the circuit
	var edgeIndex int
	c.circuitEdges, edgeIndex = model.SplitEdge(c.circuitEdges, newLocation.edge, newLocation.vertex)
	c.insertVertex(edgeIndex+1, newLocation.vertex)
	edgeLen++

	// 5. Retrieve the edges created by the split.
	splitEdgeA, splitEdgeB := c.circuitEdges[edgeIndex], c.circuitEdges[edgeIndex+1%edgeLen]

	// 6. Update the circuit distance by adding in the new distance increase.
	// Note: the new location's distance already accounts for the previous location's distance, so the old distance doesn't need to be subtracted.
	c.distanceIncreases[newLocation.vertex] += newLocation.distance
	c.length += newLocation.distance

	// 7. Determine the new distance increase/decrease for start and end vertices of the merged and split edges.
	// Note: The merged start and end may need to be offset if the new location is earlier in the circuit array than the previous location.
	splitStartDistanceDelta := c.updateDistanceIncrease(edgeIndex)
	splitEndDistanceDelta := c.updateDistanceIncrease(edgeIndex + 2)
	if edgeIndex+1 < detachIndex {
		detachIndex++
	}
	mergedStartDistanceDelta := c.updateDistanceIncrease(detachIndex - 1)
	mergedEndDistanceDelta := c.updateDistanceIncrease(detachIndex)

	// 8. Update the heap entries
	c.closestEdges = c.closestEdges.ReplaceAll(func(x interface{}) []interface{} {
		current := x.(*heapDistanceToEdge)
		if current.edge.GetStart() == detachedEdgeA.GetStart() && current.edge.GetEnd() == detachedEdgeA.GetEnd() {
			// 8a. For the merged edge, since we are going from 2 edges to 1 edge, only one of the two entries needs to be replace in the heap, so keep the entry.
			return []interface{}{
				&heapDistanceToEdge{
					vertex:   current.vertex,
					edge:     mergedEdge,
					distance: mergedEdge.DistanceIncrease(current.vertex) - c.distanceIncreases[current.vertex],
				},
			}
		} else if current.edge.GetStart() == detachedEdgeB.GetStart() && current.edge.GetEnd() == detachedEdgeB.GetEnd() {
			// 8b. For the merged edge, since we are going from 2 edges to 1 edge, only one of the two entries needs to be replace in the heap, so drop this entry.
			return []interface{}{}
		} else if current.edge.GetStart() == newLocation.edge.GetStart() && current.edge.GetEnd() == newLocation.edge.GetEnd() {
			// 8c. Replace any items in the heap with the edge that was split, with one entry for each of the newly created edges.
			return []interface{}{
				&heapDistanceToEdge{
					vertex:   current.vertex,
					edge:     splitEdgeA,
					distance: splitEdgeA.DistanceIncrease(current.vertex) - c.distanceIncreases[current.vertex],
				},
				&heapDistanceToEdge{
					vertex:   current.vertex,
					edge:     splitEdgeB,
					distance: splitEdgeB.DistanceIncrease(current.vertex) - c.distanceIncreases[current.vertex],
				},
			}
		} else if current.vertex == newLocation.vertex {
			// 8d. For any closest edges involving the moved vertex, adjust them so that they account for removing the vertex from its new location in the perimeter.
			return []interface{}{current.cloneAndAdjustDistance(-newLocation.distance)}
		} else if current.vertex == splitEdgeA.GetStart() {
			// 8g. For any closest edges involving the split edge's start vertex, adjust their distances so that they account for the change in perimeter distance when removing the start vertex.
			return []interface{}{current.cloneAndAdjustDistance(-splitStartDistanceDelta)}
		} else if current.vertex == splitEdgeB.GetEnd() {
			// 8h. For any closest edges involving the split edge's end vertex, adjust their distances so that they account for the change in perimeter distance when removing the end vertex.
			return []interface{}{current.cloneAndAdjustDistance(-splitEndDistanceDelta)}
		} else if current.vertex == mergedEdge.GetStart() {
			// 8e. For any closest edges involving the merged edge's start vertex, adjust their distances so that they account for the change in perimeter distance when removing the start vertex.
			return []interface{}{current.cloneAndAdjustDistance(-mergedStartDistanceDelta)}
		} else if current.vertex == mergedEdge.GetEnd() {
			// 8f. For any closest edges involving the merged edge's end vertex, adjust their distances so that they account for the change in perimeter distance when removing the end vertex.
			return []interface{}{current.cloneAndAdjustDistance(-mergedEndDistanceDelta)}
		} else {
			// 8i. Ignore any items in the heap that do not contain the merged edge, the split edge, the moved vertex, or any of the vertices adjacent to the moved vertex (before or after moving).
			return []interface{}{current}
		}
	})
}

func (c *HeapableCircuit2DMinClonesMove) insertVertex(index int, vertex model.CircuitVertex) {
	if index >= len(c.circuit) {
		c.circuit = append(c.circuit, vertex)
	} else {
		// copy all elements starting at the index one to the right to create a duplicate record at index and index+1.
		c.circuit = append(c.circuit[:index+1], c.circuit[index:]...)
		// update only the vertex at the index, so that there are no duplicates and the vertex is at the index.
		c.circuit[index] = vertex
	}
}

func (d *heapDistanceToEdge) cloneAndAdjustDistance(distanceAdjustment float64) *heapDistanceToEdge {
	return &heapDistanceToEdge{
		vertex:   d.vertex,
		edge:     d.edge,
		distance: d.distance + distanceAdjustment,
	}
}

func (c *HeapableCircuit2DMinClonesMove) updateDistanceIncrease(vertexIndex int) float64 {
	circuitLen := len(c.circuit)
	vertex := c.circuit[(vertexIndex+circuitLen)%circuitLen]
	prev := c.circuit[(vertexIndex-1+circuitLen)%circuitLen]
	next := c.circuit[(vertexIndex+1)%circuitLen]
	distanceIncrease := vertex.DistanceTo(prev) + vertex.DistanceTo(next) - prev.DistanceTo(next)
	distanceDelta := distanceIncrease - c.distanceIncreases[vertex]
	c.distanceIncreases[vertex] = distanceIncrease
	return distanceDelta
}

var _ model.HeapableCircuit = (*HeapableCircuit2DMinClonesMove)(nil)
