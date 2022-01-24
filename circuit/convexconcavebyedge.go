package circuit

import (
	"math"

	"github.com/fealos/lee-tsp-go/tspmodel"
)

type ConvexConcaveByEdge struct {
	Vertices              []tspmodel.CircuitVertex
	deduplicator          func([]tspmodel.CircuitVertex) []tspmodel.CircuitVertex
	perimeterBuilder      tspmodel.PerimeterBuilder
	circuits              []tspmodel.Circuit
	enableInteriorUpdates bool
}

func NewConvexConcaveByEdge(vertices []tspmodel.CircuitVertex, deduplicator tspmodel.Deduplicator, perimeterBuilder tspmodel.PerimeterBuilder, enableInteriorUpdates bool) tspmodel.Circuit {
	return &ConvexConcaveByEdge{
		Vertices:              vertices,
		deduplicator:          deduplicator,
		perimeterBuilder:      perimeterBuilder,
		enableInteriorUpdates: enableInteriorUpdates,
	}
}

func (c *ConvexConcaveByEdge) BuildPerimiter() {
	circuitEdges, unattachedVertices := c.perimeterBuilder(c.Vertices)

	closestEdges := make(map[tspmodel.CircuitVertex]*tspmodel.DistanceToEdge)
	toAttach := make(map[*ConvexConcave]*tspmodel.DistanceToEdge)

	initLength := 0.0
	for _, edge := range circuitEdges {
		initLength += edge.GetLength()
	}

	// Create a greedy circuit for each edge, with each circuit attaching that edge to its closest point.
	// This allows the greedy algorithm to detect scenarios where the points are individually closer to various edges, but are collectively closer to a different edge.
	// This increases the complexity of this circuit implementation to O(n^3), the unsmiplified form being O(e*(n-e)*(n-e)), since the greedy implementation is O(n^2) or O((n-e)^2).
	for _, e := range circuitEdges {
		circuit := &ConvexConcave{
			circuitEdges:          make([]tspmodel.CircuitEdge, len(circuitEdges)),
			Vertices:              c.Vertices,
			closestEdges:          tspmodel.NewHeap(tspmodel.GetDistanceToEdgeForHeap),
			unattachedVertices:    make(map[tspmodel.CircuitVertex]bool),
			length:                initLength,
			enableInteriorUpdates: c.enableInteriorUpdates,
		}

		copy(circuit.circuitEdges, circuitEdges)
		for k, v := range unattachedVertices {
			circuit.unattachedVertices[k] = v
		}

		vertexClosestToEdge := &tspmodel.DistanceToEdge{
			Distance: math.MaxFloat64,
		}
		for v := range unattachedVertices {
			d := e.DistanceIncrease(v)
			if d < vertexClosestToEdge.Distance {
				vertexClosestToEdge = &tspmodel.DistanceToEdge{
					Vertex:   v,
					Edge:     e,
					Distance: d,
				}
			}

			if prevClosest, okay := closestEdges[v]; !okay || d < prevClosest.Distance {
				closestEdges[v] = &tspmodel.DistanceToEdge{
					Vertex:   v,
					Edge:     e,
					Distance: d,
				}
			}
		}
		toAttach[circuit] = vertexClosestToEdge
		c.circuits = append(c.circuits, circuit)
	}

	for circuit, closestToEdge := range toAttach {
		for _, dist := range closestEdges {
			if dist.Vertex != closestToEdge.Vertex {
				// Need to create a new tspmodel.DistanceToEdge for each circuit, due to how greedy circuits update DistanceToEdges
				circuit.closestEdges.Push(&tspmodel.DistanceToEdge{
					Vertex:   dist.Vertex,
					Edge:     dist.Edge,
					Distance: dist.Distance,
				})
			}
		}
		circuit.Update(closestToEdge.Vertex, closestToEdge.Edge)
	}
}

func (c *ConvexConcaveByEdge) FindNextVertexAndEdge() (tspmodel.CircuitVertex, tspmodel.CircuitEdge) {
	if shortest := c.getShortestCircuit(); shortest != nil && len(shortest.GetUnattachedVertices()) > 0 {
		next := shortest.(*ConvexConcave).closestEdges.Peek().(*tspmodel.DistanceToEdge)
		return next.Vertex, next.Edge
	} else {
		return nil, nil
	}
}

func (c *ConvexConcaveByEdge) GetAttachedVertices() []tspmodel.CircuitVertex {
	if shortest := c.getShortestCircuit(); shortest != nil {
		return shortest.GetAttachedVertices()
	}
	return []tspmodel.CircuitVertex{}
}

func (c *ConvexConcaveByEdge) GetLength() float64 {
	if shortest := c.getShortestCircuit(); shortest != nil {
		return shortest.GetLength()
	}
	return 0.0
}

func (c *ConvexConcaveByEdge) GetUnattachedVertices() map[tspmodel.CircuitVertex]bool {
	if shortest := c.getShortestCircuit(); shortest != nil {
		return shortest.GetUnattachedVertices()
	}
	return make(map[tspmodel.CircuitVertex]bool)
}

func (c *ConvexConcaveByEdge) Prepare() {
	c.circuits = []tspmodel.Circuit{}
	c.Vertices = c.deduplicator(c.Vertices)
}

func (c *ConvexConcaveByEdge) Update(ignoredVertex tspmodel.CircuitVertex, ignoredEdge tspmodel.CircuitEdge) {
	for _, circuit := range c.circuits {
		circuit.Update(circuit.FindNextVertexAndEdge())
	}
}

func (c *ConvexConcaveByEdge) getShortestCircuit() tspmodel.Circuit {
	shortestLen := math.MaxFloat64
	var shortest tspmodel.Circuit
	for _, circuit := range c.circuits {
		if l := circuit.GetLength(); l < shortestLen {
			shortest = circuit
			shortestLen = l
		}
	}
	return shortest
}

var _ tspmodel.Circuit = (*ConvexConcaveByEdge)(nil)
