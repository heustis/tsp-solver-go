package circuit

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/heustis/tsp-solver-go/model"
)

// DisparityGreedy generates the convex hull, then performs the following steps:
// 1. Determines, and tracks, the closest two edges to each point, based on the distance increase from inserting the point along that edge.
// 2. For each point, determines the disparity between the two closest edges:
//     * If `useRelativeDisparity` is true, this calculates the disparity by dividing the larger distance increase by the smaller.
//     * If `useRelativeDisparity` is false, this calculates the disparity by subtracting the smaller distance increase from the larger.
// 3. Selects the next point to attach to the circuit, by finding the point with the largest disparity between its two closest edges.
//     * If two points have the same disparity, the point that is closer to its closest edge is chosen.
// 4. Attaches the selected point to its closest edge.
// 5. Updates the remaining unattached points, by comparing their previous closest two edges to the newly created edges (after the split), and updating their disparity if they are updated.
// 6. Repeats 3-5 until all points are attached to the circuit.
// This algorithm greedily attaches points to the convex hull by prioritizing points that have the smallest impact on the length of the circuit. In other words, it prefers the point, that when attached to its closest edge (by distance increase), increases the length of the circuit by the least.
//
// Complexity:
// * This algorithm is O(n^2) because it needs to attach each interior point to the circuit, and each time it attaches an interior point it needs to check if the newly created edges are closer to each remaining interior point than their current closest edges, so that it can update their disparity and select the correct point + edge in subsequent iterations.
type DisparityGreedy struct {
	circuitEdges  []model.CircuitEdge
	edgeDistances map[model.CircuitVertex]*vertexDisparity
	length        float64
}

// NewDisparityGreedy creates a new DisparityGreedy, builds the convex hull, and prepares the disparity metadata.
func NewDisparityGreedy(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder, useRelativeDisparity bool) *DisparityGreedy {
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

	return &DisparityGreedy{
		circuitEdges:  circuitEdges,
		edgeDistances: edgeDistances,
		length:        length,
	}
}

func (c *DisparityGreedy) FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge) {
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

func (c *DisparityGreedy) GetAttachedEdges() []model.CircuitEdge {
	return c.circuitEdges
}

func (c *DisparityGreedy) GetAttachedVertices() []model.CircuitVertex {
	vertices := make([]model.CircuitVertex, len(c.circuitEdges))
	for i, edge := range c.circuitEdges {
		vertices[i] = edge.GetStart()
	}
	return vertices
}

func (c *DisparityGreedy) GetLength() float64 {
	return c.length
}

func (c *DisparityGreedy) GetUnattachedVertices() map[model.CircuitVertex]bool {
	unattachedVertices := make(map[model.CircuitVertex]bool)
	for k := range c.edgeDistances {
		unattachedVertices[k] = true
	}
	return unattachedVertices
}

func (c *DisparityGreedy) Update(vertexToAdd model.CircuitVertex, edgeToSplit model.CircuitEdge) {
	if vertexToAdd != nil {
		var edgeIndex int
		c.circuitEdges, edgeIndex = model.SplitEdge(c.circuitEdges, edgeToSplit, vertexToAdd)
		if edgeIndex < 0 {
			expectedEdgeJson, _ := json.Marshal(edgeToSplit)
			actualCircuitJson, _ := json.Marshal(c.circuitEdges)
			panic(fmt.Errorf("edge not found in circuit=%p, expected=%s, \ncircuit=%s", c, string(expectedEdgeJson), string(actualCircuitJson)))
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

var _ model.Circuit = (*DisparityGreedy)(nil)
