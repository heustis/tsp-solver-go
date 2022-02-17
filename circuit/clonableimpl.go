package circuit

import (
	"encoding/json"
	"fmt"

	"github.com/heustis/lee-tsp-go/model"
)

type ClonableCircuitImpl struct {
	CloneOnFirstAttach bool
	Vertices           []model.CircuitVertex
	circuitEdges       []model.CircuitEdge
	closestEdges       *model.Heap
	length             float64
	vertexMetadata     map[model.CircuitVertex]*vertexStatus
	// Optional Feature: max closest edges per vertex
}

type vertexStatus struct {
	isUnattached     bool
	isConcave        bool
	distanceIncrease float64
}

func NewClonableCircuitImpl(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) *ClonableCircuitImpl {
	circuitEdges, unattachedVertices := perimeterBuilder(vertices)

	// Determine the initial length of the perimeter.
	length := 0.0
	vertexMetadata := make(map[model.CircuitVertex]*vertexStatus)
	for _, edge := range circuitEdges {
		length += edge.GetLength()
		vertexMetadata[edge.GetStart()] = &vertexStatus{
			isUnattached:     false,
			isConcave:        true,
			distanceIncrease: 0.0,
		}
	}

	// Find the closest edge for all interior points, based on distance increase (rather than perpendicular distance)
	// O(v*e) due to storing the entries, one per edge per unattached vertex, in an array then heapifing the array.
	// Note: it would be O(v*e * log(v*e)) to use PushHeap() on each entry independently.
	closestEdges := model.NewHeap(model.GetDistanceToEdgeForHeap)
	for v := range unattachedVertices {
		vertexMetadata[v] = &vertexStatus{
			isUnattached:     true,
			isConcave:        false,
			distanceIncrease: 0.0,
		}
		for _, edge := range circuitEdges {
			// Note: Using Push, not PushHeap to just append elements for now, will heapify after all elements are pushed.
			closestEdges.Push(&model.DistanceToEdge{
				Vertex:   v,
				Edge:     edge,
				Distance: edge.DistanceIncrease(v),
			})
		}
	}
	closestEdges.Heapify()

	return &ClonableCircuitImpl{
		CloneOnFirstAttach: false,
		Vertices:           vertices,
		circuitEdges:       circuitEdges,
		closestEdges:       closestEdges,
		length:             length,
		vertexMetadata:     vertexMetadata,
	}
}

func (c *ClonableCircuitImpl) CloneAndUpdate() ClonableCircuit {
	// 1. Remove 'next closest' from heap - complexity O(log n)
	next, okay := c.closestEdges.PopHeap().(*model.DistanceToEdge)

	if next == nil || !okay {
		return nil
	} else if (!c.CloneOnFirstAttach || c.hasOneUnattachedVertex()) && c.vertexMetadata[next.Vertex].isUnattached {
		// 2a. If the current vertex is unattached, attach it at the next location, provided that either:
		//     it is the last remaining vertex, or the circuit is configured not to clone on first attachment.
		c.AttachVertex(next)
		return nil
	} else if !c.closestEdges.AnyMatch(next.HasVertex) {
		// 2c. If there are no more items in the heap with the vertex in 'next closest', this is the last edge that the vertex can attach to (i.e. all other possibilities have been tried).
		//     So do not clone the circuit, rather attach 'next closest' to this circuit.
		// AnyMatch is O(n), attachVertex is O(n)
		c.AttachVertex(next)
		return nil
	} else {
		// 2c. If the 'next closest' vertex is already attached, or if this is configured to clone on the first attachment,
		// clone the circuit with the 'next closest' vertex attached to the 'next closest' edge.
		// O(n) due to copying circuits and interior vertices
		clone := &ClonableCircuitImpl{
			CloneOnFirstAttach: c.CloneOnFirstAttach,
			Vertices:           c.Vertices,
			circuitEdges:       make([]model.CircuitEdge, len(c.circuitEdges)),
			closestEdges:       c.closestEdges.Clone(),
			length:             c.length,
			vertexMetadata:     make(map[model.CircuitVertex]*vertexStatus),
		}
		copy(clone.circuitEdges, c.circuitEdges)

		for k, v := range c.vertexMetadata {
			// We clone on write, so we don't need to create copies at this time.
			clone.vertexMetadata[k] = v
		}

		if c.CloneOnFirstAttach {
			clone.AttachVertex(next)
		} else {
			// Move the vertex from the previous location to the new location in the clone, if the point was already attached.
			clone.MoveVertex(next)
		}

		return clone
	}
}

func (c *ClonableCircuitImpl) Delete() {
	for k := range c.vertexMetadata {
		delete(c.vertexMetadata, k)
	}
	if c.closestEdges != nil {
		c.closestEdges.Delete()
	}
	c.Vertices = nil
	c.circuitEdges = nil
	c.closestEdges = nil
}

func (c *ClonableCircuitImpl) FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge) {
	if next, okay := c.closestEdges.Peek().(*model.DistanceToEdge); okay && next != nil {
		return next.Vertex, next.Edge
	}
	return nil, nil
}

func (c *ClonableCircuitImpl) GetAttachedVertices() []model.CircuitVertex {
	vertices := make([]model.CircuitVertex, len(c.circuitEdges))
	for i, edge := range c.circuitEdges {
		vertices[i] = edge.GetStart()
	}
	return vertices
}

func (c *ClonableCircuitImpl) GetAttachedEdges() []model.CircuitEdge {
	return c.circuitEdges
}

func (c *ClonableCircuitImpl) GetClosestEdges() *model.Heap {
	return c.closestEdges
}

func (c *ClonableCircuitImpl) GetLength() float64 {
	return c.length
}

func (c *ClonableCircuitImpl) GetLengthWithNext() float64 {
	if next := c.closestEdges.Peek(); next != nil {
		nextDistToEdge := next.(*model.DistanceToEdge)
		if len(c.circuitEdges) == len(c.Vertices) && nextDistToEdge.Distance > 0 {
			return c.length // If the circuit is complete and the next vertex to attach increases the perimeter length, the circuit is optimal.
		} else {
			return c.length + nextDistToEdge.Distance
		}
	} else {
		return c.length
	}
}

func (c *ClonableCircuitImpl) GetUnattachedVertices() map[model.CircuitVertex]bool {
	unattachedVertices := make(map[model.CircuitVertex]bool)
	for k, v := range c.vertexMetadata {
		if v.isUnattached {
			unattachedVertices[k] = true
		}
	}
	return unattachedVertices
}

func (c *ClonableCircuitImpl) AttachVertex(toAttach *model.DistanceToEdge) {
	// 1. Update the circuitEdges and retrieve the newly created edges
	var edgeIndex int
	c.circuitEdges, edgeIndex = model.SplitEdgeCopy(c.circuitEdges, toAttach.Edge, toAttach.Vertex)
	if edgeIndex < 0 {
		expectedEdgeJson, _ := json.Marshal(toAttach.Edge)
		actualCircuitJson, _ := json.Marshal(c.circuitEdges)
		initialVertices, _ := json.Marshal(c.Vertices)
		panic(fmt.Errorf("edge not found in circuit=%p, expected=%s, \ncircuit=%s \nvertices=%s", c, string(expectedEdgeJson), string(actualCircuitJson), string(initialVertices)))
	}
	edgeA, edgeB := c.circuitEdges[edgeIndex], c.circuitEdges[edgeIndex+1]

	// 2. Update the circuit length and the distances increases as a result of the attached vertex.
	//    Note - the model.DistanceToEdge already accounts for both the existing edge and the new edge.
	c.length += toAttach.Distance

	updatedVertices := make(map[model.CircuitVertex]bool)
	updatedVertices[toAttach.Vertex] = true
	updatedVertices[toAttach.Edge.GetStart()] = true
	updatedVertices[toAttach.Edge.GetEnd()] = true
	c.updateDistanceIncreases(updatedVertices)

	// 3. Replace any references to the merged edge with two entries for the newly created edges.
	//    Complexity is O(v*e) due having one entry in the heap per edge per initially unattached vertex.
	c.closestEdges.ReplaceAll(func(x interface{}) interface{} {
		current := x.(*model.DistanceToEdge)
		existingIncrease := c.vertexMetadata[current.Vertex].distanceIncrease
		if current.Edge.GetStart() == toAttach.Edge.GetStart() && current.Edge.GetEnd() == toAttach.Edge.GetEnd() {
			if current.Vertex == toAttach.Vertex {
				// 3a. If the current DistanceToEdge equals toAttach, remove it, since its vertex was just attached to its edge.
				return nil
			}
			// 3b. Replace any DistanceToEdges that reference the split edge with one entry for each of the new edges.
			return []interface{}{
				&model.DistanceToEdge{
					Vertex:   current.Vertex,
					Edge:     edgeA,
					Distance: edgeA.DistanceIncrease(current.Vertex) - existingIncrease,
				},
				&model.DistanceToEdge{
					Vertex:   current.Vertex,
					Edge:     edgeB,
					Distance: edgeB.DistanceIncrease(current.Vertex) - existingIncrease,
				},
			}
		} else if c.CloneOnFirstAttach && current.Vertex == toAttach.Vertex {
			// 3c. If the circuit is cloned on every attachment, remove any items in the heap with toAttach's vertex.
			return nil
		} else if updatedVertices[current.Vertex] {
			// 3d. If the circuit is only cloned when moving a vertex, update the distances in the heap to account for the vertex's new position in the circuit.
			//     This updates any entries for the attached vertex, as well as the start and end points of the split edge.
			return &model.DistanceToEdge{
				Vertex:   current.Vertex,
				Edge:     current.Edge,
				Distance: current.Edge.DistanceIncrease(current.Vertex) - existingIncrease,
			}
		} else {
			// 3e. If the DistanceToEdge does not reference the attached vertex, nor the split edge, it is unmodified by this AttachVertex.
			return x
		}
	})
}

func (c *ClonableCircuitImpl) MoveVertex(toMove *model.DistanceToEdge) {
	// 1. Remove the edge with the vertex from the circuitEdges
	var mergedEdge, splitEdgeA, splitEdgeB model.CircuitEdge
	c.circuitEdges, mergedEdge, splitEdgeA, splitEdgeB = model.MoveVertex(c.circuitEdges, toMove.Vertex, toMove.Edge)
	if mergedEdge == nil {
		toMoveJson, _ := json.Marshal(toMove.Vertex)
		targetEdgeJson, _ := json.Marshal(toMove.Edge)
		actualCircuitJson, _ := json.Marshal(c.circuitEdges)
		panic(fmt.Errorf("cannot move vertex to edge circuit=%p, vertex=%s, edge=%s circuit=%s", c, string(toMoveJson), string(targetEdgeJson), string(actualCircuitJson)))
	}

	// 2. Update the circuit distance and the distances increases as a result of the attached vertex.
	//    Note - the model.DistanceToEdge already accounts for both the existing edge and the new edge.
	c.length += toMove.Distance

	updatedVertices := make(map[model.CircuitVertex]bool)
	updatedVertices[toMove.Vertex] = true
	updatedVertices[toMove.Edge.GetStart()] = true
	updatedVertices[toMove.Edge.GetEnd()] = true
	updatedVertices[mergedEdge.GetStart()] = true
	updatedVertices[mergedEdge.GetEnd()] = true
	c.updateDistanceIncreases(updatedVertices)

	// 3. Replace any references to the merged edges in the heap with a single entry for the merged edge.
	//    Replace any reference to the split edge with two entries, one for each new edge.
	//    Update the distance increases for the modified vertices.
	//    Complexity is O(v*e) due having one entry in the heap per edge per initially unattached vertex.
	replacedVertices := make(map[model.CircuitVertex]bool)
	c.closestEdges.ReplaceAll(func(x interface{}) interface{} {
		current := x.(*model.DistanceToEdge)
		existingIncrease := c.vertexMetadata[current.Vertex].distanceIncrease
		if current.Vertex == toMove.Vertex {
			// 3a. Remove any entries pertatining to the moved vertex, since the source circuit this is cloned from already has those entries.
			//     This introduces some inaccuracies in the optimum circuit computation, since the new location may have a larger distance increase than the previous location,
			//     which influences the order that DistanceToEdges appear in the heap.
			return nil
		} else if current.Edge.GetStart() == toMove.Vertex || current.Edge.GetEnd() == toMove.Vertex {
			// 3b. Only create one entry for the merged edge, even if the source heap has two entries.
			if replacedVertices[current.Vertex] {
				return nil
			}
			replacedVertices[current.Vertex] = true

			// 3c. Do not allow an entry to be created for either of the vertices of the merged edge.
			//     For example, if point B is detached from between points A & C, point C could have an existing entry for A-B, which would be replaced by A-C.
			//     The way that this scenario happens is that B and C are both internal points, B is attached first, C is attached between B and D, leaving an entry for A-B for vertex C.
			if current.Vertex == mergedEdge.GetStart() || current.Vertex == mergedEdge.GetEnd() {
				return nil
			}
			// 3d. Create a new the entry for the merged edge, in place of the edge that reference the moved vertex.
			return &model.DistanceToEdge{
				Vertex:   current.Vertex,
				Edge:     mergedEdge,
				Distance: mergedEdge.DistanceIncrease(current.Vertex) - existingIncrease,
			}
		} else if current.Edge.GetStart() == toMove.Edge.GetStart() && current.Edge.GetEnd() == toMove.Edge.GetEnd() {
			// 3e. Replace any DistanceToEdges that reference the split edge with one entry for each of the new edges.
			return []interface{}{
				&model.DistanceToEdge{
					Vertex:   current.Vertex,
					Edge:     splitEdgeA,
					Distance: splitEdgeA.DistanceIncrease(current.Vertex) - existingIncrease,
				},
				&model.DistanceToEdge{
					Vertex:   current.Vertex,
					Edge:     splitEdgeB,
					Distance: splitEdgeB.DistanceIncrease(current.Vertex) - existingIncrease,
				},
			}
		} else if updatedVertices[current.Vertex] {
			// 3d. Update the distances in the heap to account for the vertex's new position in the circuit.
			//     This updates any entries for the attached vertex, the start and end points of the split edge, and the start and end points for the merged edge.
			return &model.DistanceToEdge{
				Vertex:   current.Vertex,
				Edge:     current.Edge,
				Distance: current.Edge.DistanceIncrease(current.Vertex) - existingIncrease,
			}
		} else {
			return x
		}
	})
}

func (c *ClonableCircuitImpl) hasOneUnattachedVertex() bool {
	numUnattached := 0
	for _, v := range c.vertexMetadata {
		if v.isUnattached {
			numUnattached++
			if numUnattached > 1 {
				return false
			}
		}
	}
	return numUnattached == 1
}

func (c *ClonableCircuitImpl) updateDistanceIncreases(updatedVertices map[model.CircuitVertex]bool) {
	circuitLen := len(c.circuitEdges)
	if circuitLen >= 3 {
		prev := c.circuitEdges[circuitLen-1]
		for _, edge := range c.circuitEdges {
			v := edge.GetStart()
			if currentMetadata := c.vertexMetadata[v]; updatedVertices[v] && !currentMetadata.isConcave {
				c.vertexMetadata[v] = &vertexStatus{
					isUnattached:     false,
					isConcave:        false,
					distanceIncrease: prev.GetLength() + edge.GetLength() - prev.GetStart().DistanceTo(edge.GetEnd()),
				}
			}
			prev = edge
		}
	}
}

var _ ClonableCircuit = (*ClonableCircuitImpl)(nil)
