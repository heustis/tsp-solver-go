package circuit

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/fealos/lee-tsp-go/model"
)

type ConvexConcave struct {
	Vertices              []model.CircuitVertex
	circuitEdges          []model.CircuitEdge
	closestEdges          *model.Heap
	enableInteriorUpdates bool
	interiorVertices      map[model.CircuitVertex]bool
	length                float64
	unattachedVertices    map[model.CircuitVertex]bool
}

func NewConvexConcave(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder, enableInteriorUpdates bool) model.Circuit {
	circuitEdges, unattachedVertices := perimeterBuilder(vertices)

	// Find the closest edge for all interior points, based on distance increase; store them in a heap for retrieval from closest to farthest.
	interiorVertices := make(map[model.CircuitVertex]bool)
	closestEdges := model.NewHeap(model.GetDistanceToEdgeForHeap)
	for vertex := range unattachedVertices {
		if enableInteriorUpdates {
			interiorVertices[vertex] = true
		}
		closest := model.FindClosestEdge(vertex, circuitEdges)
		closestEdges.PushHeap(&model.DistanceToEdge{
			Vertex:   vertex,
			Edge:     closest,
			Distance: closest.DistanceIncrease(vertex),
		})
	}

	length := 0.0
	for _, edge := range circuitEdges {
		length += edge.GetLength()
	}

	return &ConvexConcave{
		Vertices:              vertices,
		circuitEdges:          circuitEdges,
		closestEdges:          closestEdges,
		enableInteriorUpdates: enableInteriorUpdates,
		interiorVertices:      interiorVertices,
		length:                length,
		unattachedVertices:    unattachedVertices,
	}
}

func (c *ConvexConcave) FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge) {
	if next, okay := c.closestEdges.PopHeap().(*model.DistanceToEdge); okay {
		return next.Vertex, next.Edge
	} else {
		return nil, nil
	}
}

func (c *ConvexConcave) GetAttachedEdges() []model.CircuitEdge {
	return c.circuitEdges
}

func (c *ConvexConcave) GetAttachedVertices() []model.CircuitVertex {
	vertices := make([]model.CircuitVertex, len(c.circuitEdges))
	for i, edge := range c.circuitEdges {
		vertices[i] = edge.GetStart()
	}
	return vertices
}

func (c *ConvexConcave) GetClosestEdges() *model.Heap {
	return c.closestEdges
}

func (c *ConvexConcave) GetInteriorVertices() map[model.CircuitVertex]bool {
	return c.interiorVertices
}

func (c *ConvexConcave) GetLength() float64 {
	return c.length
}

func (c *ConvexConcave) GetUnattachedVertices() map[model.CircuitVertex]bool {
	return c.unattachedVertices
}

func (c *ConvexConcave) Update(vertexToAdd model.CircuitVertex, edgeToSplit model.CircuitEdge) {
	if vertexToAdd != nil {
		var edgeIndex int
		c.circuitEdges, edgeIndex = model.SplitEdge(c.circuitEdges, edgeToSplit, vertexToAdd)
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

func (c *ConvexConcave) getClosestEdgeForAttachedPoint(vertex model.CircuitVertex) model.CircuitEdge {
	prev := c.circuitEdges[len(c.circuitEdges)-1]
	for _, edge := range c.circuitEdges {
		if edge.GetStart() == vertex {
			return prev.GetStart().EdgeTo(edge.GetEnd())
		}

		prev = edge
	}
	return nil
}

func (c *ConvexConcave) updateClosestEdges(removedEdge model.CircuitEdge, edgeA model.CircuitEdge, edgeB model.CircuitEdge) {
	c.length += edgeA.GetLength() + edgeB.GetLength() - removedEdge.GetLength()
	for _, x := range c.closestEdges.GetValues() {
		previous := x.(*model.DistanceToEdge)
		distA := edgeA.DistanceIncrease(previous.Vertex)
		distB := edgeB.DistanceIncrease(previous.Vertex)
		if distA < previous.Distance && distA <= distB {
			previous.Edge = edgeA
			previous.Distance = distA
		} else if distB < previous.Distance {
			previous.Edge = edgeB
			previous.Distance = distB
		} else if previous.Edge == removedEdge {
			previous.Edge = model.FindClosestEdge(previous.Vertex, c.circuitEdges)
			previous.Distance = previous.Edge.DistanceIncrease(previous.Vertex)
		}
	}
	c.closestEdges.Heapify()
}

func (c *ConvexConcave) updateInteriorPoints(removedEdge model.CircuitEdge, edgeA model.CircuitEdge, edgeB model.CircuitEdge) {
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
			c.circuitEdges, _, _, _ = model.MergeEdgesByVertex(c.circuitEdges, vertex)
			// This will be  updated by ReplaceAll in the next step, so the edge value and distance are unimportant.
			c.closestEdges.PushHeap(&model.DistanceToEdge{
				Vertex:   vertex,
				Edge:     nil,
				Distance: math.MaxFloat64,
			})
		}
	}

	// Since multiple edges could have been replaced (due to both the newly attached point and any removed points) recalculate the closest edge for each unattached vertex.
	c.closestEdges.ReplaceAll(func(x interface{}) interface{} {
		previous := x.(*model.DistanceToEdge)
		previous.Edge = model.FindClosestEdge(previous.Vertex, c.circuitEdges)
		previous.Distance = previous.Edge.DistanceIncrease(previous.Vertex)
		return previous
	})
}

var _ model.Circuit = (*ConvexConcave)(nil)
