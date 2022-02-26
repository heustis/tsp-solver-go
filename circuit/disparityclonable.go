package circuit

import (
	"sort"

	"github.com/heustis/tsp-solver-go/model"
	"github.com/heustis/tsp-solver-go/stats"
)

const minimumSignificance = 1.0
const maxClones uint16 = 1000

// DisparityClonable relies on the following priniciples to approximate the smallest concave circuit:
// 1. That the minimum convex hull of a 2D circuit must be traversed in that order to be optimal (i.e. swapping the order of any 2 vertices in the hull will result in edges intersecting.)
//    1a. This means that each convex hull vertex may have any number of interior points between and the next convex hull vertex, but that adding the interior vertices to the circuit cannot reorder these vertices.
// 2. Interior points are either near an edge, near a corner, or near the middle of a set of edges (i.e. similarly close to several edges, possibly all edges).
//    2a. A point that is close to a single edge will have a significant disparity between the distance increase of its closest edge, and the distance increase of all other edges.
//    2b. A point that is close to a corner of two edges will have a significant disparity between the distance increase of those two corner edges, and the distance increase of all other edges.
//    2c. A point that is near the middle of a group of edges may or may not have a significant disparity between its distance increase
// 3. As interior points are connected to the circuit, other points will move from '2c' to '2a' or '2b' (or become exterior points).
//    3a. This is because the new concave edges will be closer to the other interior points than the previous convex edges were.
//    3b. If a point becomes exterior, ignore edges that would intersect a closer edge if the point attached to the farther edge.
//        In other words, if the exterior point is close to a concave corner, it could attach to either edge without intersecting the other.
//        However, if it is near a convex corner, the farther edge would have to cross the closer edge to attach to the point.
//    3c. If all points are in 2c, clone the circuit once per edge and attach that edge to its closest edge, then solve each of those clones in parallel.
type DisparityClonable struct {
	significance     float64
	maxClones        uint16
	perimeterBuilder model.PerimeterBuilder
	circuits         []*disparityClonableCircuit
}

func NewDisparityClonable(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) *DisparityClonable {
	circuitEdges, unattachedVertices := perimeterBuilder(vertices)

	initCircuit := &disparityClonableCircuit{
		edges:     circuitEdges,
		distances: make(map[model.CircuitVertex]*stats.DistanceGaps),
		length:    0.0,
	}
	circuits := []*disparityClonableCircuit{initCircuit}

	// Find the closest edge for all interior points, based on distance increase; store them in a heap for retrieval from closest to farthest.
	for vertex := range unattachedVertices {
		initCircuit.distances[vertex] = stats.NewDistanceGaps(vertex, circuitEdges)
	}

	for _, edge := range circuitEdges {
		initCircuit.length += edge.GetLength()
	}

	return &DisparityClonable{
		significance:     minimumSignificance,
		maxClones:        maxClones,
		circuits:         circuits,
		perimeterBuilder: perimeterBuilder,
	}
}

func (c *DisparityClonable) FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge) {
	// Since updating may need to clone the circuits, and each circuit may need to updated with a different vertex, we just need to return any unattached point and edge.
	if len(c.circuits) > 0 {
		for k, v := range c.circuits[0].distances {
			return k, v.ClosestEdges[0].Edge
		}
	}
	return nil, nil
}

func (c *DisparityClonable) GetAttachedEdges() []model.CircuitEdge {
	if len(c.circuits) > 0 {
		return c.circuits[0].edges
	}
	return []model.CircuitEdge{}
}

func (c *DisparityClonable) GetAttachedVertices() []model.CircuitVertex {
	if len(c.circuits) > 0 && len(c.circuits[0].edges) > 0 {
		vertices := make([]model.CircuitVertex, len(c.circuits[0].edges))
		for i, edge := range c.circuits[0].edges {
			vertices[i] = edge.GetStart()
		}
		return vertices
	}
	return []model.CircuitVertex{}
}

func (c *DisparityClonable) GetLength() float64 {
	if len(c.circuits) > 0 {
		return c.circuits[0].length
	}
	return 0.0
}

func (c *DisparityClonable) GetUnattachedVertices() map[model.CircuitVertex]bool {
	unattachedVertices := make(map[model.CircuitVertex]bool)
	if len(c.circuits) > 0 {
		for k := range c.circuits[0].distances {
			unattachedVertices[k] = true
		}
	}
	return unattachedVertices
}

func (c *DisparityClonable) SetMaxClones(max uint16) {
	c.maxClones = max
}

func (c *DisparityClonable) SetSignificance(minSignificance float64) {
	c.significance = minSignificance
}

func (c *DisparityClonable) Update(ignoredVertex model.CircuitVertex, ignoredEdge model.CircuitEdge) {
	// Don't update if the perimeter has not been built, nor if the shortest circuit is completed.
	if len(c.circuits) == 0 || len(c.circuits[0].distances) == 0 {
		return
	}
	// Note: track updated and cloned circuits in a separate array once we need to clone at least one circuit.
	// Do not mutate the 'circuits' array while we are iterating over it, replace it with the updated/cloned array afterward (if appropriate).
	var updatedCircuits []*disparityClonableCircuit
	useUpdated := false
	for i, circuit := range c.circuits {
		if clones := circuit.update(c.significance); len(clones) > 0 || useUpdated {
			// Add all previously processed circuits the first time the clone array is constructed.
			if !useUpdated {
				updatedCircuits = make([]*disparityClonableCircuit, 0, len(c.circuits)+len(clones))
				updatedCircuits = append(updatedCircuits, c.circuits[0:i]...)
				useUpdated = true
			}
			updatedCircuits = append(updatedCircuits, circuit)
			updatedCircuits = append(updatedCircuits, clones...)
		}
	}
	if useUpdated {
		c.circuits = updatedCircuits
	}
	// Sort the updated slice from smallest to largest, preferring circuits that are close to completion.
	sort.Slice(c.circuits, func(i, j int) bool {
		return c.circuits[i].getLengthPerVertex() < c.circuits[j].getLengthPerVertex()
	})

	if len(c.circuits) > int(c.maxClones) {
		c.circuits = c.circuits[0:c.maxClones]
	}
}

var _ model.Circuit = (*DisparityClonable)(nil)
