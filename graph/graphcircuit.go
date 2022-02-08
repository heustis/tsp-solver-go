package graph

import (
	"math"

	"github.com/fealos/lee-tsp-go/model"
)

type GraphCircuit struct {
	graph              *Graph
	circuit            []model.CircuitEdge
	edges              map[model.CircuitVertex]map[model.CircuitVertex]model.CircuitEdge
	length             float64
	unattachedVertices map[model.CircuitVertex]bool
}

func NewGraphCircuit(g *Graph) *GraphCircuit {
	c := &GraphCircuit{
		graph:              g,
		circuit:            []model.CircuitEdge{},
		edges:              g.PathToAllFromAll(),
		length:             0.0,
		unattachedVertices: make(map[model.CircuitVertex]bool),
	}
	for _, v := range c.graph.Vertices {
		c.unattachedVertices[v] = true
	}
	c.BuildPerimiter()
	return c
}

func (c *GraphCircuit) BuildPerimiter() {
	// Determine the "midpoint" by finding the node with the smallest average distance to all other nodes.
	var midpoint *GraphVertex
	midpointAvgDistance := math.MaxFloat64

	numVertices := float64(len(c.graph.Vertices))
	for _, v := range c.graph.Vertices {
		avgDist := 0.0

		for _, edge := range c.edges[v] {
			// Divide prior to adding to avoid overflowing the distance
			avgDist += edge.GetLength() / numVertices
		}

		if avgDist < midpointAvgDistance {
			midpoint = v
			midpointAvgDistance = avgDist
		}
	}

	// The first point on the perimeter will be the point farthest from the midpoint.
	farthestFromMid := c.farthestVertexFrom(midpoint)
	delete(c.unattachedVertices, farthestFromMid)

	// The second point on the perimeter will be the point farthest from the previous point.
	farthestFromFarthest := c.farthestVertexFrom(farthestFromMid)
	delete(c.unattachedVertices, farthestFromFarthest)

	// Initialize the length - note that the edges could be asymetric in the graph.
	c.length = c.edges[farthestFromMid][farthestFromFarthest].GetLength() + c.edges[farthestFromFarthest][farthestFromMid].GetLength()
	c.circuit = append(c.circuit, c.edges[farthestFromMid][farthestFromFarthest])
	c.circuit = append(c.circuit, c.edges[farthestFromFarthest][farthestFromMid])

	// Attach vertices to the circuit until all vertices are either interior vertices or attached to the circuit.
	exteriorVertices := make(map[model.CircuitVertex]bool)
	for k, v := range c.unattachedVertices {
		exteriorVertices[k] = v
	}
	// var interiorVertex model.CircuitVertex = nil
	for len(exteriorVertices) > 0 {
		// Since a graph does not have to follow 2D and 3D geometric principles, we need to recompute the closest edges each time.
		// With 2D we know that the closest edge (or edges that are created by splitting it) will remain the closest edge of any external point since the perimeter is convex.
		// However, in a graph A->B and B->C each may be farther away from an external point Z than A->C, causing a different edge (D->E) to become the closest edge to Z.
		closestEdges := c.findClosestEdges(exteriorVertices)
		farthestFromClosestEdge := &model.DistanceToEdge{
			Distance: -1.0,
		}
		for _, closest := range closestEdges {
			if closest.Distance > farthestFromClosestEdge.Distance {
				farthestFromClosestEdge = closest
			}
		}

		c.insertVertex(farthestFromClosestEdge.Vertex, farthestFromClosestEdge.Edge)
		delete(exteriorVertices, farthestFromClosestEdge.Vertex)
		delete(closestEdges, farthestFromClosestEdge.Vertex)

		// Update exterior vertices.
		// A vertex is interior if it is closer to the perimeter vertices (other than those in its closest edge) than the closest vertex on its closest edge.
		// i.e. isExterior := dist(C, P) > dist(projection(E, C), P)
		for extVertex, extClosest := range closestEdges {
			isExterior := true

			var closestEdgeVertex *GraphVertex = nil
			for _, edgeVertex := range extClosest.Edge.(*GraphEdge).path {
				if closestEdgeVertex == nil || c.edges[extVertex][edgeVertex].GetLength() < c.edges[extVertex][closestEdgeVertex].GetLength() {
					closestEdgeVertex = edgeVertex
				}
			}

			for _, edge := range c.circuit {
				start := edge.GetStart()
				if start == extClosest.Edge.GetStart() || start == extClosest.Edge.GetEnd() {
					continue
				}
				if c.edges[extVertex][start].GetLength() < c.edges[closestEdgeVertex][start].GetLength() {
					isExterior = false
					break
				}
			}

			if !isExterior {
				delete(exteriorVertices, extVertex)
			}
		}

		// Delete the closestEdges map
		for k := range closestEdges {
			delete(closestEdges, k)
		}
	}
}

func (c *GraphCircuit) Delete() {
	for startVertex, edgeMap := range c.edges {
		for endVertex, edge := range edgeMap {
			if d, okay := edge.(model.Deletable); okay {
				d.Delete()
			}
			delete(edgeMap, endVertex)
		}
		delete(c.edges, startVertex)
	}
	for k := range c.unattachedVertices {
		delete(c.unattachedVertices, k)
	}
	c.circuit = nil
	// Note: Do not delete the graph, since it was supplied to NewGraphCircuit rather than created by this GraphCircuit.
}

func (c *GraphCircuit) FindNextVertexAndEdge() (vertexToAdd model.CircuitVertex, edgeToSplit model.CircuitEdge) {
	closest := &model.DistanceToEdge{
		Distance: math.MaxFloat64,
		Vertex:   nil,
		Edge:     nil,
	}
	for v := range c.unattachedVertices {
		for _, edge := range c.circuit {
			if dist := c.distanceIncrease(v, edge); dist < closest.Distance {
				closest.Vertex = v
				closest.Edge = edge
				closest.Distance = dist
			}
		}
	}

	return closest.Vertex, closest.Edge
}

func (c *GraphCircuit) GetAttachedVertices() []model.CircuitVertex {
	vertices := make([]model.CircuitVertex, len(c.circuit))
	for i, e := range c.circuit {
		vertices[i] = e.GetStart()
	}
	return vertices
}

func (c *GraphCircuit) GetEdgeFor(start model.CircuitVertex, end model.CircuitVertex) model.CircuitEdge {
	if e, okay := c.edges[start]; okay {
		return e[end]
	} else {
		return nil
	}
}

func (c *GraphCircuit) GetLength() float64 {
	return c.length
}

func (c *GraphCircuit) GetUnattachedVertices() map[model.CircuitVertex]bool {
	return c.unattachedVertices
}

func (c *GraphCircuit) Update(vertexToAdd model.CircuitVertex, edgeToSplit model.CircuitEdge) {
	c.insertVertex(vertexToAdd, edgeToSplit)
}

func (c *GraphCircuit) distanceIncrease(v model.CircuitVertex, e model.CircuitEdge) float64 {
	return c.edges[e.GetStart()][v].GetLength() + c.edges[v][e.GetEnd()].GetLength() - e.GetLength()
}

func (c *GraphCircuit) farthestVertexFrom(v model.CircuitVertex) model.CircuitVertex {
	farthestPoint := v
	farthestDist := 0.0

	for other, edge := range c.edges[v] {
		if edge.GetLength() > farthestDist {
			farthestDist = edge.GetLength()
			farthestPoint = other
		}
	}

	return farthestPoint
}

func (c *GraphCircuit) findClosestEdges(vertices map[model.CircuitVertex]bool) map[model.CircuitVertex]*model.DistanceToEdge {
	closestEdges := make(map[model.CircuitVertex]*model.DistanceToEdge)
	for v := range vertices {
		closestEdges[v] = &model.DistanceToEdge{
			Vertex:   v,
			Edge:     nil,
			Distance: math.MaxFloat64,
		}
		for _, edge := range c.circuit {
			if dist := c.distanceIncrease(v, edge); dist < closestEdges[v].Distance {
				closestEdges[v].Edge = edge
				closestEdges[v].Distance = dist
			}
		}
	}
	return closestEdges
}

func (c *GraphCircuit) insertVertex(v model.CircuitVertex, edgeToSplit model.CircuitEdge) (model.CircuitEdge, model.CircuitEdge) {
	edgeA, edgeB := c.edges[edgeToSplit.GetStart()][v], c.edges[v][edgeToSplit.GetEnd()]

	edgeIndex := model.IndexOfEdge(c.circuit, edgeToSplit)
	c.circuit = append(c.circuit[:edgeIndex+1], c.circuit[edgeIndex:]...)
	c.circuit[edgeIndex] = edgeA
	c.circuit[edgeIndex+1] = edgeB

	c.length += edgeA.GetLength() + edgeB.GetLength() - edgeToSplit.GetLength()
	delete(c.unattachedVertices, v)

	return edgeA, edgeB
}

var _ model.Circuit = (*GraphCircuit)(nil)
