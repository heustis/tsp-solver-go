package circuit

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/fealos/lee-tsp-go/model"
)

type ConvexConcaveDisparity struct {
	Vertices         []model.CircuitVertex
	deduplicator     func([]model.CircuitVertex) []model.CircuitVertex
	perimeterBuilder model.PerimeterBuilder
	circuitEdges     []model.CircuitEdge
	edgeDistances    map[model.CircuitVertex]*vertexDisparity
	length           float64
}

func NewConvexConcaveDisparity(vertices []model.CircuitVertex, deduplicator func([]model.CircuitVertex) []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) model.Circuit {
	return &ConvexConcaveDisparity{
		Vertices:         vertices,
		deduplicator:     deduplicator,
		perimeterBuilder: perimeterBuilder,
	}
}

func (c *ConvexConcaveDisparity) BuildPerimiter() {
	var unattachedVertices map[model.CircuitVertex]bool
	c.circuitEdges, unattachedVertices = c.perimeterBuilder.BuildPerimiter(c.Vertices)

	// Find the closest edge for all interior points, based on distance increase; store them in a heap for retrieval from closest to farthest.
	for vertex := range unattachedVertices {
		disparity := &vertexDisparity{
			closestEdge:       &model.DistanceToEdge{Vertex: vertex, Distance: math.MaxFloat64},
			secondClosestEdge: &model.DistanceToEdge{Vertex: vertex, Distance: math.MaxFloat64},
		}

		for _, e := range c.circuitEdges {
			disparity.update(e, e.DistanceIncrease(vertex))
		}

		disparity.amount = disparity.secondClosestEdge.Distance - disparity.closestEdge.Distance
		c.edgeDistances[vertex] = disparity
	}

	c.length = 0.0
	for _, edge := range c.circuitEdges {
		c.length += edge.GetLength()
	}
}

func (c *ConvexConcaveDisparity) FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge) {
	maxDisparityAmount := -1.0
	next := &model.DistanceToEdge{
		Distance: math.MaxFloat64,
	}
	for _, v := range c.edgeDistances {
		if v.amount > maxDisparityAmount+model.Threshold || (v.amount > maxDisparityAmount-model.Threshold && v.closestEdge.Distance < next.Distance) {
			next = v.closestEdge
			maxDisparityAmount = v.amount
		}
	}
	return next.Vertex, next.Edge
}

func (c *ConvexConcaveDisparity) GetAttachedEdges() []model.CircuitEdge {
	return c.circuitEdges
}

func (c *ConvexConcaveDisparity) GetAttachedVertices() []model.CircuitVertex {
	vertices := make([]model.CircuitVertex, len(c.circuitEdges))
	for i, edge := range c.circuitEdges {
		vertices[i] = edge.GetStart()
	}
	return vertices
}

func (c *ConvexConcaveDisparity) GetLength() float64 {
	return c.length
}

func (c *ConvexConcaveDisparity) GetUnattachedVertices() map[model.CircuitVertex]bool {
	unattachedVertices := make(map[model.CircuitVertex]bool)
	for k := range c.edgeDistances {
		unattachedVertices[k] = true
	}
	return unattachedVertices
}

func (c *ConvexConcaveDisparity) Prepare() {
	c.edgeDistances = make(map[model.CircuitVertex]*vertexDisparity)
	c.circuitEdges = []model.CircuitEdge{}
	c.length = 0.0

	c.Vertices = c.deduplicator(c.Vertices)
}

func (c *ConvexConcaveDisparity) Update(vertexToAdd model.CircuitVertex, edgeToSplit model.CircuitEdge) {
	if vertexToAdd != nil {
		var edgeIndex int
		c.circuitEdges, edgeIndex = model.SplitEdge(c.circuitEdges, edgeToSplit, vertexToAdd)
		if edgeIndex < 0 {
			expectedEdgeJson, _ := json.Marshal(edgeToSplit)
			actualCircuitJson, _ := json.Marshal(c.circuitEdges)
			initialVertices, _ := json.Marshal(c.Vertices)
			panic(fmt.Errorf("edge not found in circuit=%p, expected=%s, \ncircuit=%s \nvertices=%s", c, string(expectedEdgeJson), string(actualCircuitJson), string(initialVertices)))
		}
		delete(c.edgeDistances, vertexToAdd)
		edgeA, edgeB := c.circuitEdges[edgeIndex], c.circuitEdges[edgeIndex+1]
		c.length += edgeA.GetLength() + edgeB.GetLength() - edgeToSplit.GetLength()

		for vertex, disparity := range c.edgeDistances {
			disparity.remove(edgeToSplit)
			disparity.update(edgeA, edgeA.DistanceIncrease(vertex))
			disparity.update(edgeB, edgeB.DistanceIncrease(vertex))
		}
	}
}

type vertexDisparity struct {
	closestEdge       *model.DistanceToEdge
	secondClosestEdge *model.DistanceToEdge
	amount            float64
}

func (disparity *vertexDisparity) update(e model.CircuitEdge, distance float64) {
	if distance < disparity.secondClosestEdge.Distance {
		disparity.secondClosestEdge.Distance = distance
		disparity.secondClosestEdge.Edge = e
		// Swap the closest and second closest edges, if this edge is now the closest (since it is already stored in second closest)
		if distance < disparity.closestEdge.Distance {
			disparity.closestEdge, disparity.secondClosestEdge = disparity.secondClosestEdge, disparity.closestEdge
		}
	}
}

func (disparity *vertexDisparity) remove(e model.CircuitEdge) bool {
	if e == disparity.closestEdge.Edge {
		disparity.closestEdge.Distance = math.MaxFloat64
		disparity.closestEdge.Edge = nil
		disparity.closestEdge, disparity.secondClosestEdge = disparity.secondClosestEdge, disparity.closestEdge
		return true
	} else if e == disparity.secondClosestEdge.Edge {
		disparity.secondClosestEdge.Distance = math.MaxFloat64
		disparity.secondClosestEdge.Edge = nil
		return true
	}
	return false
}

var _ model.Circuit = (*ConvexConcaveDisparity)(nil)
