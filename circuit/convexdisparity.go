package circuit

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/heustis/tsp-solver-go/model"
)

type ConvexConcaveDisparity struct {
	Vertices      []model.CircuitVertex
	circuitEdges  []model.CircuitEdge
	edgeDistances map[model.CircuitVertex]*vertexDisparity
	length        float64
}

func NewConvexConcaveDisparity(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder, useRelativeDisparity bool) *ConvexConcaveDisparity {
	circuitEdges, unattachedVertices := perimeterBuilder(vertices)

	edgeDistances := make(map[model.CircuitVertex]*vertexDisparity)

	// Find the two closest edges for each unattached vertex, so that we can compute the disparity between those two edges' distance increases.
	// We will use that disparity to determine which vertex to attach during each iteration.
	// Since we are only looking at the first change in distance increases, we only need to store the closest two edges, and update them as the circuit changes.
	for vertex := range unattachedVertices {
		disparity := &vertexDisparity{
			closestEdge:       &model.DistanceToEdge{Vertex: vertex, Distance: math.MaxFloat64},
			secondClosestEdge: &model.DistanceToEdge{Vertex: vertex, Distance: math.MaxFloat64},
		}

		if useRelativeDisparity {
			disparity.disparityFunction = func(closer *model.DistanceToEdge, farther *model.DistanceToEdge) float64 {
				// Avoid divide by 0, also it should be impossible to have a negative distance increase at this stage of computation.
				if closer.Distance < model.Threshold {
					return math.MaxFloat64
				}
				return farther.Distance / closer.Distance
			}
		} else {
			disparity.disparityFunction = func(closer *model.DistanceToEdge, farther *model.DistanceToEdge) float64 {
				return farther.Distance - closer.Distance
			}
		}

		for _, e := range circuitEdges {
			disparity.update(e, e.DistanceIncrease(vertex))
		}

		disparity.amount = disparity.secondClosestEdge.Distance - disparity.closestEdge.Distance
		edgeDistances[vertex] = disparity
	}

	length := 0.0
	for _, edge := range circuitEdges {
		length += edge.GetLength()
	}

	return &ConvexConcaveDisparity{
		Vertices:      vertices,
		circuitEdges:  circuitEdges,
		edgeDistances: edgeDistances,
		length:        length,
	}
}

func (c *ConvexConcaveDisparity) FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge) {
	maxDisparityAmount := -1.0
	next := &model.DistanceToEdge{
		Distance: math.MaxFloat64,
	}
	// Find the vertex with the largest gap between its two closest edges.
	// If two vertices have approximately the same gap between their closest edges, select the vertex that is closer to its closest edge.
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
	disparityFunction func(closer *model.DistanceToEdge, farther *model.DistanceToEdge) float64
}

func (disparity *vertexDisparity) update(e model.CircuitEdge, distance float64) {
	if distance < disparity.secondClosestEdge.Distance {
		disparity.secondClosestEdge.Distance = distance
		disparity.secondClosestEdge.Edge = e
		// Swap the closest and second closest edges, if this edge is now the closest (since it is already stored in second closest)
		if distance < disparity.closestEdge.Distance {
			disparity.closestEdge, disparity.secondClosestEdge = disparity.secondClosestEdge, disparity.closestEdge
		}
		disparity.amount = disparity.disparityFunction(disparity.closestEdge, disparity.secondClosestEdge)
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
