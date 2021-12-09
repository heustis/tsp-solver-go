package model

import (
	"encoding/json"
	"fmt"
	"math"
)

type CircuitGreedyWithUpdatesImpl struct {
	Vertices           []CircuitVertex
	deduplicator       func([]CircuitVertex) []CircuitVertex
	perimeterBuilder   PerimeterBuilder
	circuitEdges       []CircuitEdge
	closestEdges       *Heap
	interiorVertices   map[CircuitVertex]bool
	unattachedVertices map[CircuitVertex]bool
	length             float64
}

func NewCircuitGreedyWithUpdatesImpl(vertices []CircuitVertex, deduplicator func([]CircuitVertex) []CircuitVertex, perimeterBuilder PerimeterBuilder) Circuit {
	return &CircuitGreedyWithUpdatesImpl{
		Vertices:         vertices,
		deduplicator:     deduplicator,
		perimeterBuilder: perimeterBuilder,
	}
}

func (c *CircuitGreedyWithUpdatesImpl) BuildPerimiter() {
	c.circuitEdges, c.unattachedVertices = c.perimeterBuilder.BuildPerimiter(c.Vertices)

	// Find the closest edge for all interior points, based on distance increase (rather than perpendicular distance)
	for vertex := range c.unattachedVertices {
		c.interiorVertices[vertex] = true
		closest := FindClosestEdge(vertex, c.circuitEdges)
		c.closestEdges.PushHeap(&DistanceToEdge{
			Vertex:   vertex,
			Edge:     closest,
			Distance: closest.DistanceIncrease(vertex),
		})
	}

	c.length = 0.0
	for _, edge := range c.circuitEdges {
		c.length += edge.GetLength()
	}
}

func (c *CircuitGreedyWithUpdatesImpl) FindNextVertexAndEdge() (CircuitVertex, CircuitEdge) {
	if next, okay := c.closestEdges.PopHeap().(*DistanceToEdge); okay {
		return next.Vertex, next.Edge
	} else {
		return nil, nil
	}
}

func (c *CircuitGreedyWithUpdatesImpl) GetAttachedEdges() []CircuitEdge {
	return c.circuitEdges
}

func (c *CircuitGreedyWithUpdatesImpl) GetAttachedVertices() []CircuitVertex {
	vertices := make([]CircuitVertex, len(c.circuitEdges))
	for i, edge := range c.circuitEdges {
		vertices[i] = edge.GetStart()
	}
	return vertices
}

func (c *CircuitGreedyWithUpdatesImpl) GetClosestEdges() *Heap {
	return c.closestEdges
}

func (c *CircuitGreedyWithUpdatesImpl) GetInteriorVertices() map[CircuitVertex]bool {
	return c.interiorVertices
}

func (c *CircuitGreedyWithUpdatesImpl) GetLength() float64 {
	return c.length
}

func (c *CircuitGreedyWithUpdatesImpl) GetUnattachedVertices() map[CircuitVertex]bool {
	return c.unattachedVertices
}

func (c *CircuitGreedyWithUpdatesImpl) Prepare() {
	c.interiorVertices = make(map[CircuitVertex]bool)
	c.unattachedVertices = make(map[CircuitVertex]bool)
	c.closestEdges = NewHeap(GetDistanceToEdgeForHeap)
	c.circuitEdges = []CircuitEdge{}
	c.length = 0.0

	c.Vertices = c.deduplicator(c.Vertices)
}

func (c *CircuitGreedyWithUpdatesImpl) Update(vertexToAdd CircuitVertex, edgeToSplit CircuitEdge) {
	if vertexToAdd != nil {
		var edgeIndex int
		c.circuitEdges, edgeIndex = SplitEdge(c.circuitEdges, edgeToSplit, vertexToAdd)
		if edgeIndex < 0 {
			expectedEdgeJson, _ := json.Marshal(edgeToSplit)
			actualCircuitJson, _ := json.Marshal(c.circuitEdges)
			initialVertices, _ := json.Marshal(c.Vertices)
			panic(fmt.Errorf("edge not found in circuit=%p, expected=%s, \ncircuit=%s \nvertices=%s", c, string(expectedEdgeJson), string(actualCircuitJson), string(initialVertices)))
		}

		delete(c.unattachedVertices, vertexToAdd)
		c.updateInteriorPoints(edgeToSplit, c.circuitEdges[edgeIndex], c.circuitEdges[edgeIndex+1])
	}
}

func (c *CircuitGreedyWithUpdatesImpl) getClosestEdgeForAttachedPoint(vertex CircuitVertex) CircuitEdge {
	prev := c.circuitEdges[len(c.circuitEdges)-1]
	for _, edge := range c.circuitEdges {
		if edge.GetStart() == vertex {
			return prev.GetStart().EdgeTo(edge.GetEnd())
		}

		prev = edge
	}
	return nil
}

func (c *CircuitGreedyWithUpdatesImpl) updateInteriorPoints(removedEdge CircuitEdge, edgeA CircuitEdge, edgeB CircuitEdge) {
	c.length += edgeA.GetLength() + edgeB.GetLength() - removedEdge.GetLength()
	// Detach any interior, attached vertices that are now closer to either created edge than they are to their attached edge.
	for vertex := range c.interiorVertices {
		// Ignore unattached vertices and vertices attached to one of the newly created edges.
		if c.unattachedVertices[vertex] || edgeA.GetStart() == vertex || edgeA.GetEnd() == vertex || edgeB.GetEnd() == vertex {
			continue
		}
		closestAttached := c.getClosestEdgeForAttachedPoint(vertex)
		previousDistance := closestAttached.DistanceIncrease(vertex)
		if edgeA.DistanceIncrease(vertex) < previousDistance || edgeB.DistanceIncrease(vertex) < previousDistance {
			c.unattachedVertices[vertex] = true
			c.circuitEdges, _, _, _ = MergeEdges2(c.circuitEdges, vertex)
			// This will be  updated by ReplaceAll in the next step, so the edge value and distance are unimportant.
			c.closestEdges.PushHeap(&DistanceToEdge{
				Vertex:   vertex,
				Edge:     nil,
				Distance: math.MaxFloat64,
			})
		}
	}

	// Since multiple edges could have been replaced (due to both the newly attached point and any removed points) recalculate the closest edge for each unattached vertex.
	c.closestEdges.ReplaceAll2(func(x interface{}) interface{} {
		previous := x.(*DistanceToEdge)
		previous.Edge = FindClosestEdge(previous.Vertex, c.circuitEdges)
		previous.Distance = previous.Edge.DistanceIncrease(previous.Vertex)
		return previous
	})
}

var _ Circuit = (*CircuitGreedyWithUpdatesImpl)(nil)
