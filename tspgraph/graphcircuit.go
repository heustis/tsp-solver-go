package tspgraph

import (
	"math"

	"github.com/fealos/lee-tsp-go/tspmodel"
)

type GraphCircuit struct {
	tspgraph           *Graph
	circuit            []tspmodel.CircuitEdge
	edges              map[tspmodel.CircuitVertex]map[tspmodel.CircuitVertex]tspmodel.CircuitEdge
	length             float64
	unattachedVertices map[tspmodel.CircuitVertex]bool
}

func NewGraphCircuit(g *Graph) *GraphCircuit {
	return &GraphCircuit{
		tspgraph: g,
	}
}

func (c *GraphCircuit) BuildPerimiter() {
	// Determine the "midpoint" by finding the node with the smallest average distance to all other nodes.
	var midpoint *GraphVertex
	midpointAvgDistance := math.MaxFloat64

	numVertices := float64(len(c.tspgraph.Vertices))
	for _, v := range c.tspgraph.Vertices {
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

	// Initialize the length - note that the edges could be asymetric in the tspgraph.
	c.length = c.edges[farthestFromMid][farthestFromFarthest].GetLength() + c.edges[farthestFromFarthest][farthestFromMid].GetLength()
	c.circuit = append(c.circuit, c.edges[farthestFromMid][farthestFromFarthest])
	c.circuit = append(c.circuit, c.edges[farthestFromFarthest][farthestFromMid])

	// Attach vertices to the circuit until all vertices are either interior vertices or attached to the circuit.
	exteriorVertices := make(map[tspmodel.CircuitVertex]bool)
	for k, v := range c.unattachedVertices {
		exteriorVertices[k] = v
	}
	// var interiorVertex tspmodel.CircuitVertex = nil
	for len(exteriorVertices) > 0 {
		// Since a tspgraph does not have to follow 2D and 3D geometric principles, we need to recompute the closest edges each time.
		// With 2D we know that the closest edge (or edges that are created by splitting it) will remain the closest edge of any external point since the perimeter is convex.
		// However, in a tspgraph A->B and B->C each may be farther away from an external point Z than A->C, causing a different edge (D->E) to become the closest edge to Z.
		closestEdges := c.findClosestEdges(exteriorVertices)
		farthestFromClosestEdge := &tspmodel.DistanceToEdge{
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
			if d, okay := edge.(tspmodel.Deletable); okay {
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
	// Note: Do not delete the tspgraph, since it was supplied to NewGraphCircuit rather than created by this GraphCircuit.
}

func (c *GraphCircuit) FindNextVertexAndEdge() (vertexToAdd tspmodel.CircuitVertex, edgeToSplit tspmodel.CircuitEdge) {
	closest := &tspmodel.DistanceToEdge{
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

func (c *GraphCircuit) GetAttachedVertices() []tspmodel.CircuitVertex {
	vertices := make([]tspmodel.CircuitVertex, len(c.circuit))
	for i, e := range c.circuit {
		vertices[i] = e.GetStart()
	}
	return vertices
}

func (c *GraphCircuit) GetEdgeFor(start tspmodel.CircuitVertex, end tspmodel.CircuitVertex) tspmodel.CircuitEdge {
	if e, okay := c.edges[start]; okay {
		return e[end]
	} else {
		return nil
	}
}

func (c *GraphCircuit) GetLength() float64 {
	return c.length
}

func (c *GraphCircuit) GetUnattachedVertices() map[tspmodel.CircuitVertex]bool {
	return c.unattachedVertices
}

func (c *GraphCircuit) Prepare() {
	c.circuit = []tspmodel.CircuitEdge{}
	c.edges = c.tspgraph.PathToAllFromAll()
	c.length = 0.0
	c.unattachedVertices = make(map[tspmodel.CircuitVertex]bool)
	for _, v := range c.tspgraph.Vertices {
		c.unattachedVertices[v] = true
	}
}

func (c *GraphCircuit) Update(vertexToAdd tspmodel.CircuitVertex, edgeToSplit tspmodel.CircuitEdge) {
	c.insertVertex(vertexToAdd, edgeToSplit)
}

func (c *GraphCircuit) distanceIncrease(v tspmodel.CircuitVertex, e tspmodel.CircuitEdge) float64 {
	return c.edges[e.GetStart()][v].GetLength() + c.edges[v][e.GetEnd()].GetLength() - e.GetLength()
}

func (c *GraphCircuit) farthestVertexFrom(v tspmodel.CircuitVertex) tspmodel.CircuitVertex {
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

func (c *GraphCircuit) findClosestEdges(vertices map[tspmodel.CircuitVertex]bool) map[tspmodel.CircuitVertex]*tspmodel.DistanceToEdge {
	closestEdges := make(map[tspmodel.CircuitVertex]*tspmodel.DistanceToEdge)
	for v := range vertices {
		closestEdges[v] = &tspmodel.DistanceToEdge{
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

func (c *GraphCircuit) insertVertex(v tspmodel.CircuitVertex, edgeToSplit tspmodel.CircuitEdge) (tspmodel.CircuitEdge, tspmodel.CircuitEdge) {
	edgeA, edgeB := c.edges[edgeToSplit.GetStart()][v], c.edges[v][edgeToSplit.GetEnd()]

	edgeIndex := tspmodel.IndexOfEdge(c.circuit, edgeToSplit)
	c.circuit = append(c.circuit[:edgeIndex+1], c.circuit[edgeIndex:]...)
	c.circuit[edgeIndex] = edgeA
	c.circuit[edgeIndex+1] = edgeB

	c.length += edgeA.GetLength() + edgeB.GetLength() - edgeToSplit.GetLength()
	delete(c.unattachedVertices, v)

	return edgeA, edgeB
}

var _ tspmodel.Circuit = (*GraphCircuit)(nil)
