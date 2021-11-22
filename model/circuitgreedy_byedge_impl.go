package model

import (
	"math"
)

type CircuitGreedyByEdgeImpl struct {
	Vertices         []CircuitVertex
	deduplicator     func([]CircuitVertex) []CircuitVertex
	perimeterBuilder PerimeterBuilder
	circuits         []Circuit
}

func NewCircuitGreedyByEdgeImpl(vertices []CircuitVertex, deduplicator func([]CircuitVertex) []CircuitVertex, perimeterBuilder PerimeterBuilder) Circuit {
	return &CircuitGreedyByEdgeImpl{
		Vertices:         vertices,
		deduplicator:     deduplicator,
		perimeterBuilder: perimeterBuilder,
	}
}

func (c *CircuitGreedyByEdgeImpl) BuildPerimiter() {
	_, circuitEdges, unattachedVertices := c.perimeterBuilder.BuildPerimiter(c.Vertices)

	closestEdges := make(map[CircuitVertex]*DistanceToEdge)
	toAttach := make(map[*CircuitGreedyImpl]*DistanceToEdge)

	initLength := 0.0
	for _, edge := range circuitEdges {
		initLength += edge.GetLength()
	}

	// Create a greedy circuit for each edge, with each circuit attaching that edge to its closest point.
	// This allows the greedy algorithm to detect scenarios where the points are individually closer to various edges, but are collectively closer to a different edge.
	// This increases the complexity of this circuit implementation to O(n^3), the unsmiplified form being O(e*(n-e)*(n-e)), since the greedy implementation is O(n^2) or O((n-e)^2).
	for _, e := range circuitEdges {
		circuit := &CircuitGreedyImpl{
			circuitEdges:       make([]CircuitEdge, len(circuitEdges)),
			Vertices:           c.Vertices,
			closestEdges:       NewHeap(GetDistanceToEdgeForHeap),
			unattachedVertices: make(map[CircuitVertex]bool),
			length:             initLength,
		}

		copy(circuit.circuitEdges, circuitEdges)
		for k, v := range unattachedVertices {
			circuit.unattachedVertices[k] = v
		}

		vertexClosestToEdge := &DistanceToEdge{
			Distance: math.MaxFloat64,
		}
		for v := range unattachedVertices {
			d := e.DistanceIncrease(v)
			if d < vertexClosestToEdge.Distance {
				vertexClosestToEdge = &DistanceToEdge{
					Vertex:   v,
					Edge:     e,
					Distance: d,
				}
			}

			if prevClosest, okay := closestEdges[v]; !okay || d < prevClosest.Distance {
				closestEdges[v] = &DistanceToEdge{
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
				// Need to create a new DistanceToEdge for each circuit, due to how greedy circuits update DistanceToEdges
				circuit.closestEdges.Push(&DistanceToEdge{
					Vertex:   dist.Vertex,
					Edge:     dist.Edge,
					Distance: dist.Distance,
				})
			}
		}
		circuit.Update(closestToEdge.Vertex, closestToEdge.Edge)
	}
}

func (c *CircuitGreedyByEdgeImpl) FindNextVertexAndEdge() (CircuitVertex, CircuitEdge) {
	if shortest := c.getShortestCircuit(); shortest != nil && len(shortest.GetUnattachedVertices()) > 0 {
		next := shortest.(*CircuitGreedyImpl).closestEdges.Peek().(*DistanceToEdge)
		return next.Vertex, next.Edge
	} else {
		return nil, nil
	}
}

func (c *CircuitGreedyByEdgeImpl) GetAttachedVertices() []CircuitVertex {
	if shortest := c.getShortestCircuit(); shortest != nil {
		return shortest.GetAttachedVertices()
	}
	return []CircuitVertex{}
}

func (c *CircuitGreedyByEdgeImpl) GetLength() float64 {
	if shortest := c.getShortestCircuit(); shortest != nil {
		return shortest.GetLength()
	}
	return 0.0
}

func (c *CircuitGreedyByEdgeImpl) GetUnattachedVertices() map[CircuitVertex]bool {
	if shortest := c.getShortestCircuit(); shortest != nil {
		return shortest.GetUnattachedVertices()
	}
	return make(map[CircuitVertex]bool)
}

func (c *CircuitGreedyByEdgeImpl) Prepare() {
	c.circuits = []Circuit{}
	c.Vertices = c.deduplicator(c.Vertices)
}

func (c *CircuitGreedyByEdgeImpl) Update(ignoredVertex CircuitVertex, ignoredEdge CircuitEdge) {
	for _, circuit := range c.circuits {
		circuit.Update(circuit.FindNextVertexAndEdge())
	}
}

func (c *CircuitGreedyByEdgeImpl) getShortestCircuit() Circuit {
	shortestLen := math.MaxFloat64
	var shortest Circuit
	for _, circuit := range c.circuits {
		if l := circuit.GetLength(); l < shortestLen {
			shortest = circuit
			shortestLen = l
		}
	}
	return shortest
}

var _ Circuit = (*CircuitGreedyImpl)(nil)
