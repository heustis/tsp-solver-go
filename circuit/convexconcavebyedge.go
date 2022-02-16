package circuit

import (
	"math"

	"github.com/fealos/lee-tsp-go/model"
)

type ConvexConcaveByEdge struct {
	Vertices              []model.CircuitVertex
	circuits              []model.Circuit
	enableInteriorUpdates bool
}

func NewConvexConcaveByEdge(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder, enableInteriorUpdates bool) *ConvexConcaveByEdge {
	circuitEdges, unattachedVertices := perimeterBuilder(vertices)

	closestEdges := make(map[model.CircuitVertex]*model.DistanceToEdge)
	toAttach := make(map[*ConvexConcave]*model.DistanceToEdge)

	initLength := 0.0
	for _, edge := range circuitEdges {
		initLength += edge.GetLength()
	}

	circuits := make([]model.Circuit, len(circuitEdges))
	// Create a greedy circuit for each edge, with each circuit attaching that edge to its closest point.
	// This allows the greedy algorithm to detect scenarios where the points are individually closer to various edges, but are collectively closer to a different edge.
	// This increases the complexity of this circuit implementation to O(n^3), the unsmiplified form being O(e*(n-e)*(n-e)), since the greedy implementation is O(n^2) or O((n-e)^2).
	for i, e := range circuitEdges {
		circuit := &ConvexConcave{
			circuitEdges:          make([]model.CircuitEdge, len(circuitEdges)),
			Vertices:              vertices,
			closestEdges:          model.NewHeap(model.GetDistanceToEdgeForHeap),
			unattachedVertices:    make(map[model.CircuitVertex]bool),
			length:                initLength,
			enableInteriorUpdates: enableInteriorUpdates,
		}

		copy(circuit.circuitEdges, circuitEdges)
		for k, v := range unattachedVertices {
			circuit.unattachedVertices[k] = v
		}

		vertexClosestToEdge := &model.DistanceToEdge{
			Distance: math.MaxFloat64,
		}
		for v := range unattachedVertices {
			d := e.DistanceIncrease(v)
			if d < vertexClosestToEdge.Distance {
				vertexClosestToEdge = &model.DistanceToEdge{
					Vertex:   v,
					Edge:     e,
					Distance: d,
				}
			}

			if prevClosest, okay := closestEdges[v]; !okay || d < prevClosest.Distance {
				closestEdges[v] = &model.DistanceToEdge{
					Vertex:   v,
					Edge:     e,
					Distance: d,
				}
			}
		}
		toAttach[circuit] = vertexClosestToEdge
		circuits[i] = circuit
	}

	for circuit, closestToEdge := range toAttach {
		for _, dist := range closestEdges {
			if dist.Vertex != closestToEdge.Vertex {
				// Need to create a new model.DistanceToEdge for each circuit, due to how greedy circuits update DistanceToEdges
				circuit.closestEdges.Push(&model.DistanceToEdge{
					Vertex:   dist.Vertex,
					Edge:     dist.Edge,
					Distance: dist.Distance,
				})
			}
		}
		circuit.Update(closestToEdge.Vertex, closestToEdge.Edge)
	}

	return &ConvexConcaveByEdge{
		Vertices:              vertices,
		enableInteriorUpdates: enableInteriorUpdates,
		circuits:              circuits,
	}
}

func (c *ConvexConcaveByEdge) FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge) {
	if shortest := c.getShortestCircuit(); shortest != nil && len(shortest.GetUnattachedVertices()) > 0 {
		next := shortest.(*ConvexConcave).closestEdges.Peek().(*model.DistanceToEdge)
		return next.Vertex, next.Edge
	} else {
		return nil, nil
	}
}

func (c *ConvexConcaveByEdge) GetAttachedVertices() []model.CircuitVertex {
	if shortest := c.getShortestCircuit(); shortest != nil {
		return shortest.GetAttachedVertices()
	}
	return []model.CircuitVertex{}
}

func (c *ConvexConcaveByEdge) GetLength() float64 {
	if shortest := c.getShortestCircuit(); shortest != nil {
		return shortest.GetLength()
	}
	return 0.0
}

func (c *ConvexConcaveByEdge) GetUnattachedVertices() map[model.CircuitVertex]bool {
	if shortest := c.getShortestCircuit(); shortest != nil {
		return shortest.GetUnattachedVertices()
	}
	return make(map[model.CircuitVertex]bool)
}

func (c *ConvexConcaveByEdge) Update(ignoredVertex model.CircuitVertex, ignoredEdge model.CircuitEdge) {
	for _, circuit := range c.circuits {
		circuit.Update(circuit.FindNextVertexAndEdge())
	}
}

func (c *ConvexConcaveByEdge) getShortestCircuit() model.Circuit {
	shortestLen := math.MaxFloat64
	var shortest model.Circuit
	for _, circuit := range c.circuits {
		if l := circuit.GetLength(); l < shortestLen {
			shortest = circuit
			shortestLen = l
		}
	}
	return shortest
}

var _ model.Circuit = (*ConvexConcaveByEdge)(nil)
