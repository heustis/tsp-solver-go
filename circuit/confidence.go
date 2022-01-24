package circuit

import (
	"fmt"
	"sort"

	"github.com/fealos/lee-tsp-go/tspmodel"
	"github.com/fealos/lee-tsp-go/tspstats"
)

const minimumSignificance = 1.0
const maxClones uint16 = 1000

// ConvexConcaveConfidence relies on the following priniciples to approximate the smallest concave circuit:
// 1. That the minimum convex hull of a 2D circuit must be traversed in that order to be optimal (i.e. swapping the order of any 2 vertices in the hull will result in edges intersecting.)
//    1a. This means that each convex hull vertex may have any number of interior points between and the next convex hull vertex, but that adding the interior vertices to the circuit cannot reorder these vertices.
// 2. Interior points are either near an edge, near a corner, or near the middle of a set of edges (i.e. similarly close to several edges, possibly all edges).
//    2a. A point that is close to a single edge will have a significant disparity between the distance increase of its closest edge, and the distance increase of all other edges.
//    2b. A point that is close to a corner of two edges will have a significant disparity between the distance increase of those two corner edges, and the distance increase of all other edges.
//    2c. A point that is near the middle of a group of edges may or may not have a significant disparity between its distance increase
// 3. As interior points are connected to the circuit, other points will move from '2c' to '2a' or '2b' (or become exterior points).
//    2a. This is because the new concave edges will be closer to the other interior points than the previous convex edges were.
//    2b. If a point becomes exterior, ignore edges that would intersect a closer edge if the point attached to the farther edge.
//        In other words, if the exterior point is close to a concave corner, it could attach to either edge without intersecting the other.
//        However, if it is near a convex corner, the farther edge would have to cross the closer edge to attach to the point.
//    2c. If all points are in 2c, clone the circuit once per edge and attach that edge to its closest edge, then solve each of those clones in parallel.
type ConvexConcaveConfidence struct {
	Vertices         []tspmodel.CircuitVertex
	Significance     float64
	MaxClones        uint16
	deduplicator     func([]tspmodel.CircuitVertex) []tspmodel.CircuitVertex
	perimeterBuilder tspmodel.PerimeterBuilder
	circuits         []*confidenceCircuit
}

func NewConvexConcaveConfidence(vertices []tspmodel.CircuitVertex, deduplicator tspmodel.Deduplicator, perimeterBuilder tspmodel.PerimeterBuilder) tspmodel.Circuit {
	return &ConvexConcaveConfidence{
		Vertices:         vertices,
		Significance:     minimumSignificance,
		MaxClones:        maxClones,
		deduplicator:     deduplicator,
		perimeterBuilder: perimeterBuilder,
	}
}

func (c *ConvexConcaveConfidence) BuildPerimiter() {
	circuitEdges, unattachedVertices := c.perimeterBuilder(c.Vertices)

	initCircuit := &confidenceCircuit{
		edges:     circuitEdges,
		distances: make(map[tspmodel.CircuitVertex]*tspstats.DistanceGaps),
		length:    0.0,
	}
	c.circuits = []*confidenceCircuit{initCircuit}

	// Find the closest edge for all interior points, based on distance increase; store them in a heap for retrieval from closest to farthest.
	for vertex := range unattachedVertices {
		initCircuit.distances[vertex] = tspstats.NewDistanceGaps(vertex, circuitEdges)
	}

	for _, edge := range circuitEdges {
		initCircuit.length += edge.GetLength()
	}
}

func (c *ConvexConcaveConfidence) FindNextVertexAndEdge() (tspmodel.CircuitVertex, tspmodel.CircuitEdge) {
	// Since updating may need to clone the circuits, and each circuit may need to updated with a different vertex, we just need to return any unattached point and edge.
	if len(c.circuits) > 0 {
		for k, v := range c.circuits[0].distances {
			return k, v.ClosestEdges[0].Edge
		}
	}
	return nil, nil
}

func (c *ConvexConcaveConfidence) GetAttachedEdges() []tspmodel.CircuitEdge {
	if len(c.circuits) > 0 {
		return c.circuits[0].edges
	}
	return []tspmodel.CircuitEdge{}
}

func (c *ConvexConcaveConfidence) GetAttachedVertices() []tspmodel.CircuitVertex {
	if len(c.circuits) > 0 && len(c.circuits[0].edges) > 0 {
		vertices := make([]tspmodel.CircuitVertex, len(c.circuits[0].edges))
		for i, edge := range c.circuits[0].edges {
			vertices[i] = edge.GetStart()
		}
		return vertices
	}
	return []tspmodel.CircuitVertex{}
}

func (c *ConvexConcaveConfidence) GetLength() float64 {
	if len(c.circuits) > 0 {
		return c.circuits[0].length
	}
	return 0.0
}

func (c *ConvexConcaveConfidence) GetUnattachedVertices() map[tspmodel.CircuitVertex]bool {
	unattachedVertices := make(map[tspmodel.CircuitVertex]bool)
	if len(c.circuits) > 0 {
		for k := range c.circuits[0].distances {
			unattachedVertices[k] = true
		}
	}
	return unattachedVertices
}

func (c *ConvexConcaveConfidence) Prepare() {
	c.Vertices = c.deduplicator(c.Vertices)
}

func (c *ConvexConcaveConfidence) Update(ignoredVertex tspmodel.CircuitVertex, ignoredEdge tspmodel.CircuitEdge) {
	// Don't update if the perimeter has not been built, nor if the shortest circuit is completed.
	if len(c.circuits) == 0 || len(c.circuits[0].distances) == 0 {
		return
	}
	// Note: track updated and cloned circuits in a separate array once we need to clone at least one circuit.
	// Do not mutate the 'circuits' array while we are iterating over it, replace it with the updated/cloned array afterward (if appropriate).
	var updatedCircuits []*confidenceCircuit
	useUpdated := false
	for i, circuit := range c.circuits {
		if clones := circuit.update(c.Significance); len(clones) > 0 || useUpdated {
			// Add all previously processed circuits the first time the clone array is constructed.
			if !useUpdated {
				updatedCircuits = make([]*confidenceCircuit, 0, len(c.circuits)+len(clones))
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
	// Sort the updated slice from smallest to largest.
	sort.Slice(c.circuits, func(i, j int) bool {
		return c.circuits[i].length < c.circuits[j].length
	})

	if len(c.circuits) > int(c.MaxClones) {
		c.circuits = c.circuits[0:c.MaxClones]
	}
}

func (c *ConvexConcaveConfidence) String() string {
	s := "{\r\n\t\"vertices\":["

	vertexIndexLookup := make(map[tspmodel.CircuitVertex]int)
	edgeIndexLookup := make(map[tspmodel.CircuitEdge]int)

	lastIndex := len(c.Vertices) - 1
	for i, v := range c.Vertices {
		vertexIndexLookup[v] = i
		s += v.String()
		if i != lastIndex {
			s += ","
		}
	}

	if len(c.circuits) == 0 {
		s += "],\r\n\t\"edges\":[],\r\n\t\"edgeDistances\":[]}"
		return s
	}

	s += "],\r\n\t\"edges\":["
	lastIndex = len(c.circuits[0].edges) - 1
	for i, e := range c.circuits[0].edges {
		edgeIndexLookup[e] = i
		s += fmt.Sprintf(`{"start":%d,"end":%d,"distance":%g}`, vertexIndexLookup[e.GetStart()], vertexIndexLookup[e.GetEnd()], e.GetLength())
		if i != lastIndex {
			s += ","
		}
	}

	s += "],\r\n\t\"edgeDistances\":["
	lastIndex = len(c.circuits[0].distances) - 1
	i := 0
	for _, dist := range c.circuits[0].distances {
		s += dist.String(vertexIndexLookup, edgeIndexLookup)
		if i != lastIndex {
			s += ","
		}
		i++
	}

	s += "]}"

	return s
}

var _ tspmodel.Circuit = (*ConvexConcaveConfidence)(nil)
