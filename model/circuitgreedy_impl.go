package model

import (
	"encoding/json"
	"fmt"
)

type CircuitGreedyImpl struct {
	Vertices           []CircuitVertex
	deduplicator       func([]CircuitVertex) []CircuitVertex
	perimeterBuilder   PerimeterBuilder
	circuitEdges       []CircuitEdge
	closestEdges       *Heap
	unattachedVertices map[CircuitVertex]bool
}

func NewCircuitGreedyImpl(vertices []CircuitVertex, deduplicator func([]CircuitVertex) []CircuitVertex, perimeterBuilder PerimeterBuilder) Circuit {
	return &CircuitGreedyImpl{
		Vertices:         vertices,
		deduplicator:     deduplicator,
		perimeterBuilder: perimeterBuilder,
	}
}

func (c *CircuitGreedyImpl) BuildPerimiter() {
	_, c.circuitEdges, c.unattachedVertices = c.perimeterBuilder.BuildPerimiter(c.Vertices)

	// Find the closest edge for all interior points, based on distance increase; store them in a heap for retrieval from closest to farthest.
	for vertex := range c.unattachedVertices {
		closest := vertex.FindClosestEdge(c.circuitEdges)
		c.closestEdges.PushHeap(&DistanceToEdge{
			Vertex:   vertex,
			Edge:     closest,
			Distance: closest.DistanceIncrease(vertex),
		})
	}
}

func (c *CircuitGreedyImpl) FindNextVertexAndEdge() (CircuitVertex, CircuitEdge) {
	if next, okay := c.closestEdges.PopHeap().(*DistanceToEdge); okay {
		return next.Vertex, next.Edge
	} else {
		return nil, nil
	}
}

func (c *CircuitGreedyImpl) GetAttachedEdges() []CircuitEdge {
	return c.circuitEdges
}

func (c *CircuitGreedyImpl) GetAttachedVertices() []CircuitVertex {
	vertices := make([]CircuitVertex, len(c.circuitEdges))
	for i, edge := range c.circuitEdges {
		vertices[i] = edge.GetStart()
	}
	return vertices
}

func (c *CircuitGreedyImpl) GetClosestEdges() *Heap {
	return c.closestEdges
}

func (c *CircuitGreedyImpl) GetLength() float64 {
	length := 0.0
	for _, edge := range c.circuitEdges {
		length += edge.GetLength()
	}
	return length
}

func (c *CircuitGreedyImpl) GetUnattachedVertices() map[CircuitVertex]bool {
	return c.unattachedVertices
}

func (c *CircuitGreedyImpl) Prepare() {
	c.unattachedVertices = make(map[CircuitVertex]bool)
	c.closestEdges = NewHeap(GetDistanceToEdgeForHeap)
	c.circuitEdges = []CircuitEdge{}

	c.Vertices = c.deduplicator(c.Vertices)
}

func (c *CircuitGreedyImpl) Update(vertexToAdd CircuitVertex, edgeToSplit CircuitEdge) {
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

func (c *CircuitGreedyImpl) updateInteriorPoints(removedEdge CircuitEdge, edgeA CircuitEdge, edgeB CircuitEdge) {
	for _, x := range c.closestEdges.GetValues() {
		previous := x.(*DistanceToEdge)
		distA := edgeA.DistanceIncrease(previous.Vertex)
		distB := edgeB.DistanceIncrease(previous.Vertex)
		if distA < previous.Distance && distA <= distB {
			previous.Edge = edgeA
			previous.Distance = distA
		} else if distB < previous.Distance {
			previous.Edge = edgeB
			previous.Distance = distB
		} else if previous.Edge == removedEdge {
			previous.Edge = previous.Vertex.FindClosestEdge(c.circuitEdges)
			previous.Distance = previous.Edge.DistanceIncrease(previous.Vertex)
		}
	}
	c.closestEdges.Heapify()
}

var _ Circuit = (*CircuitGreedyImpl)(nil)
