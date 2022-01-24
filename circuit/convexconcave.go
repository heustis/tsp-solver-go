package circuit

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/fealos/lee-tsp-go/tspmodel"
)

type ConvexConcave struct {
	Vertices              []tspmodel.CircuitVertex
	deduplicator          func([]tspmodel.CircuitVertex) []tspmodel.CircuitVertex
	perimeterBuilder      tspmodel.PerimeterBuilder
	circuitEdges          []tspmodel.CircuitEdge
	closestEdges          *tspmodel.Heap
	interiorVertices      map[tspmodel.CircuitVertex]bool
	unattachedVertices    map[tspmodel.CircuitVertex]bool
	length                float64
	enableInteriorUpdates bool
}

func NewConvexConcave(vertices []tspmodel.CircuitVertex, deduplicator tspmodel.Deduplicator, perimeterBuilder tspmodel.PerimeterBuilder, enableInteriorUpdates bool) tspmodel.Circuit {
	return &ConvexConcave{
		Vertices:              vertices,
		deduplicator:          deduplicator,
		perimeterBuilder:      perimeterBuilder,
		enableInteriorUpdates: enableInteriorUpdates,
	}
}

func (c *ConvexConcave) BuildPerimiter() {
	c.circuitEdges, c.unattachedVertices = c.perimeterBuilder(c.Vertices)

	// Find the closest edge for all interior points, based on distance increase; store them in a heap for retrieval from closest to farthest.
	for vertex := range c.unattachedVertices {
		if c.enableInteriorUpdates {
			c.interiorVertices[vertex] = true
		}
		closest := tspmodel.FindClosestEdge(vertex, c.circuitEdges)
		c.closestEdges.PushHeap(&tspmodel.DistanceToEdge{
			Vertex:   vertex,
			Edge:     closest,
			Distance: closest.DistanceIncrease(vertex),
		})
	}

	c.length = 0.0
	for _, edge := range c.circuitEdges {
		c.length += edge.GetLength()
	}
}

func (c *ConvexConcave) FindNextVertexAndEdge() (tspmodel.CircuitVertex, tspmodel.CircuitEdge) {
	if next, okay := c.closestEdges.PopHeap().(*tspmodel.DistanceToEdge); okay {
		return next.Vertex, next.Edge
	} else {
		return nil, nil
	}
}

func (c *ConvexConcave) GetAttachedEdges() []tspmodel.CircuitEdge {
	return c.circuitEdges
}

func (c *ConvexConcave) GetAttachedVertices() []tspmodel.CircuitVertex {
	vertices := make([]tspmodel.CircuitVertex, len(c.circuitEdges))
	for i, edge := range c.circuitEdges {
		vertices[i] = edge.GetStart()
	}
	return vertices
}

func (c *ConvexConcave) GetClosestEdges() *tspmodel.Heap {
	return c.closestEdges
}

func (c *ConvexConcave) GetInteriorVertices() map[tspmodel.CircuitVertex]bool {
	return c.interiorVertices
}

func (c *ConvexConcave) GetLength() float64 {
	return c.length
}

func (c *ConvexConcave) GetUnattachedVertices() map[tspmodel.CircuitVertex]bool {
	return c.unattachedVertices
}

func (c *ConvexConcave) Prepare() {
	c.interiorVertices = make(map[tspmodel.CircuitVertex]bool)
	c.unattachedVertices = make(map[tspmodel.CircuitVertex]bool)
	c.closestEdges = tspmodel.NewHeap(tspmodel.GetDistanceToEdgeForHeap)
	c.circuitEdges = []tspmodel.CircuitEdge{}
	c.length = 0.0

	c.Vertices = c.deduplicator(c.Vertices)
}

func (c *ConvexConcave) Update(vertexToAdd tspmodel.CircuitVertex, edgeToSplit tspmodel.CircuitEdge) {
	if vertexToAdd != nil {
		var edgeIndex int
		c.circuitEdges, edgeIndex = tspmodel.SplitEdge(c.circuitEdges, edgeToSplit, vertexToAdd)
		if edgeIndex < 0 {
			expectedEdgeJson, _ := json.Marshal(edgeToSplit)
			actualCircuitJson, _ := json.Marshal(c.circuitEdges)
			initialVertices, _ := json.Marshal(c.Vertices)
			panic(fmt.Errorf("edge not found in circuit=%p, expected=%s, \ncircuit=%s \nvertices=%s", c, string(expectedEdgeJson), string(actualCircuitJson), string(initialVertices)))
		}
		delete(c.unattachedVertices, vertexToAdd)
		if c.enableInteriorUpdates {
			c.updateInteriorPoints(edgeToSplit, c.circuitEdges[edgeIndex], c.circuitEdges[edgeIndex+1])
		} else {
			c.updateClosestEdges(edgeToSplit, c.circuitEdges[edgeIndex], c.circuitEdges[edgeIndex+1])
		}
	}
}

func (c *ConvexConcave) getClosestEdgeForAttachedPoint(vertex tspmodel.CircuitVertex) tspmodel.CircuitEdge {
	prev := c.circuitEdges[len(c.circuitEdges)-1]
	for _, edge := range c.circuitEdges {
		if edge.GetStart() == vertex {
			return prev.GetStart().EdgeTo(edge.GetEnd())
		}

		prev = edge
	}
	return nil
}

func (c *ConvexConcave) updateClosestEdges(removedEdge tspmodel.CircuitEdge, edgeA tspmodel.CircuitEdge, edgeB tspmodel.CircuitEdge) {
	c.length += edgeA.GetLength() + edgeB.GetLength() - removedEdge.GetLength()
	for _, x := range c.closestEdges.GetValues() {
		previous := x.(*tspmodel.DistanceToEdge)
		distA := edgeA.DistanceIncrease(previous.Vertex)
		distB := edgeB.DistanceIncrease(previous.Vertex)
		if distA < previous.Distance && distA <= distB {
			previous.Edge = edgeA
			previous.Distance = distA
		} else if distB < previous.Distance {
			previous.Edge = edgeB
			previous.Distance = distB
		} else if previous.Edge == removedEdge {
			previous.Edge = tspmodel.FindClosestEdge(previous.Vertex, c.circuitEdges)
			previous.Distance = previous.Edge.DistanceIncrease(previous.Vertex)
		}
	}
	c.closestEdges.Heapify()
}

func (c *ConvexConcave) updateInteriorPoints(removedEdge tspmodel.CircuitEdge, edgeA tspmodel.CircuitEdge, edgeB tspmodel.CircuitEdge) {
	c.length += edgeA.GetLength() + edgeB.GetLength() - removedEdge.GetLength()
	// Detach any interior, attached vertices that are now closer to either created edge than they are to their attached edge.
	for vertex := range c.interiorVertices {
		// Ignore unattached vertices and vertices attached to one of the newly created edges.
		if c.unattachedVertices[vertex] || edgeA.GetStart() == vertex || edgeA.GetEnd() == vertex || edgeB.GetEnd() == vertex {
			continue
		}
		closestAttached := c.getClosestEdgeForAttachedPoint(vertex)
		previousDistance := closestAttached.DistanceIncrease(vertex)
		if edgeA.DistanceIncrease(vertex) < previousDistance || edgeB.DistanceIncrease(vertex) < previousDistance {
			c.unattachedVertices[vertex] = true
			c.circuitEdges, _, _, _ = tspmodel.MergeEdgesByVertex(c.circuitEdges, vertex)
			// This will be  updated by ReplaceAll in the next step, so the edge value and distance are unimportant.
			c.closestEdges.PushHeap(&tspmodel.DistanceToEdge{
				Vertex:   vertex,
				Edge:     nil,
				Distance: math.MaxFloat64,
			})
		}
	}

	// Since multiple edges could have been replaced (due to both the newly attached point and any removed points) recalculate the closest edge for each unattached vertex.
	c.closestEdges.ReplaceAll(func(x interface{}) interface{} {
		previous := x.(*tspmodel.DistanceToEdge)
		previous.Edge = tspmodel.FindClosestEdge(previous.Vertex, c.circuitEdges)
		previous.Distance = previous.Edge.DistanceIncrease(previous.Vertex)
		return previous
	})
}

var _ tspmodel.Circuit = (*ConvexConcave)(nil)
