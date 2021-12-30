package circuit

import (
	"encoding/json"
	"fmt"

	"github.com/fealos/lee-tsp-go/model"
)

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
// TODO - test what counts as a significant gap: standard deviation of distance increases, standard deviation of gaps, etc.
type ConvexConcaveConfidence struct {
	Vertices         []model.CircuitVertex
	deduplicator     func([]model.CircuitVertex) []model.CircuitVertex
	perimeterBuilder model.PerimeterBuilder
	circuitEdges     []model.CircuitEdge
	edgeDistances    map[model.CircuitVertex]*vertexStatistics
	length           float64
}

type vertexStatistics struct {
	ClosestEdges                       []*model.DistanceToEdge
	DistanceAverage                    float64
	DistanceStandardDeviation          float64
	DistanceDisparityAverage           float64
	DistanceDisparityStandardDeviation float64
}

func NewConvexConcaveConfidence(vertices []model.CircuitVertex, deduplicator func([]model.CircuitVertex) []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder, enableInteriorUpdates bool) model.Circuit {
	return &ConvexConcaveConfidence{
		Vertices:         vertices,
		deduplicator:     deduplicator,
		perimeterBuilder: perimeterBuilder,
	}
}

func (c *ConvexConcaveConfidence) BuildPerimiter() {
	var unattachedVertices map[model.CircuitVertex]bool
	c.circuitEdges, unattachedVertices = c.perimeterBuilder.BuildPerimiter(c.Vertices)

	// Find the closest edge for all interior points, based on distance increase; store them in a heap for retrieval from closest to farthest.
	for vertex := range unattachedVertices {
		model.FindClosestEdge(vertex, c.circuitEdges) //TODO
	}

	c.length = 0.0
	for _, edge := range c.circuitEdges {
		c.length += edge.GetLength()
	}
}

func (c *ConvexConcaveConfidence) FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge) {
	return nil, nil //TODO
}

func (c *ConvexConcaveConfidence) GetAttachedEdges() []model.CircuitEdge {
	return c.circuitEdges
}

func (c *ConvexConcaveConfidence) GetAttachedVertices() []model.CircuitVertex {
	vertices := make([]model.CircuitVertex, len(c.circuitEdges))
	for i, edge := range c.circuitEdges {
		vertices[i] = edge.GetStart()
	}
	return vertices
}

func (c *ConvexConcaveConfidence) GetLength() float64 {
	return c.length
}

func (c *ConvexConcaveConfidence) GetUnattachedVertices() map[model.CircuitVertex]bool {
	unattachedVertices := make(map[model.CircuitVertex]bool)
	for k := range c.edgeDistances {
		unattachedVertices[k] = true
	}
	return unattachedVertices
}

func (c *ConvexConcaveConfidence) Prepare() {
	c.edgeDistances = make(map[model.CircuitVertex]*vertexStatistics)
	c.circuitEdges = []model.CircuitEdge{}
	c.length = 0.0

	c.Vertices = c.deduplicator(c.Vertices)
}

func (c *ConvexConcaveConfidence) Update(vertexToAdd model.CircuitVertex, edgeToSplit model.CircuitEdge) {
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
		//TODO
	}
}

var _ model.Circuit = (*ConvexConcaveConfidence)(nil)
