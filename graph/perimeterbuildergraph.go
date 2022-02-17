package graph

import (
	"math"

	"github.com/heustis/lee-tsp-go/model"
)

// BuildPerimiter produces the smallest convex perimeter that can encompass all the vertices in the supplied array.
// This returns both the edges comprising the convex perimeter and the set of unattached (interior) vertices.
// This will panic if any of the vertices in the array are not of type GraphVertex.
func BuildPerimiter(vertices []model.CircuitVertex) (circuitEdges []model.CircuitEdge, unattachedVertices map[model.CircuitVertex]bool) {
	unattachedVertices = make(map[model.CircuitVertex]bool)
	circuitEdges = []model.CircuitEdge{}

	// Determine the "midpoint" by finding the node with the smallest average distance to all other nodes.
	var midpoint *GraphVertex
	midpointAvgDistance := math.MaxFloat64

	numVertices := float64(len(vertices))
	for _, v := range vertices {
		unattachedVertices[v] = true

		avgDist := 0.0

		for _, edge := range v.(*GraphVertex).paths {
			// Divide prior to adding to avoid overflowing the distance
			avgDist += edge.GetLength() / numVertices
		}

		if avgDist < midpointAvgDistance {
			midpoint = v.(*GraphVertex)
			midpointAvgDistance = avgDist
		}
	}

	// The first point on the perimeter will be the point farthest from the midpoint.
	farthestFromMid := farthestVertexFrom(midpoint)
	delete(unattachedVertices, farthestFromMid)

	// The second point on the perimeter will be the point farthest from the previous point.
	farthestFromFarthest := farthestVertexFrom(farthestFromMid)
	delete(unattachedVertices, farthestFromFarthest)

	// Create the first two edges - note that the edges could be asymetric in the graph.
	circuitEdges = append(circuitEdges, farthestFromMid.EdgeTo(farthestFromFarthest))
	circuitEdges = append(circuitEdges, farthestFromFarthest.EdgeTo(farthestFromMid))

	// Attach vertices to the circuit until all vertices are either interior vertices or attached to the circuit.
	exteriorVertices := make(map[model.CircuitVertex]bool)
	for k, v := range unattachedVertices {
		exteriorVertices[k] = v
	}
	// var interiorVertex model.CircuitVertex = nil
	for len(exteriorVertices) > 0 {
		// Since a graph does not have to follow 2D and 3D geometric principles, we need to recompute the closest edges each time.
		// With 2D we know that the closest edge (or edges that are created by splitting it) will remain the closest edge of any external point since the perimeter is convex.
		// However, in a graph A->B and B->C each may be farther away from an external point Z than A->C, causing a different edge (D->E) to become the closest edge to Z.
		closestEdges := findClosestEdges(circuitEdges, exteriorVertices)
		farthestFromClosestEdge := &model.DistanceToEdge{
			Distance: -1.0,
		}
		for _, closest := range closestEdges {
			if closest.Distance > farthestFromClosestEdge.Distance {
				farthestFromClosestEdge = closest
			}
		}

		circuitEdges = insertVertex(circuitEdges, farthestFromClosestEdge.Vertex, farthestFromClosestEdge.Edge)
		delete(unattachedVertices, farthestFromClosestEdge.Vertex)
		delete(exteriorVertices, farthestFromClosestEdge.Vertex)
		delete(closestEdges, farthestFromClosestEdge.Vertex)

		// Update exterior vertices.
		// A vertex is interior if it is closer to the perimeter vertices (other than those in its closest edge) than the closest vertex on its closest edge.
		// i.e. isExterior := dist(C, P) > dist(projection(E, C), P)
		for extVertex, extClosest := range closestEdges {
			isExterior := true

			var closestEdgeVertex *GraphVertex = nil
			for _, edgeVertex := range extClosest.Edge.(*GraphEdge).path {
				if closestEdgeVertex == nil || extVertex.DistanceTo(edgeVertex) < extVertex.DistanceTo(closestEdgeVertex) {
					closestEdgeVertex = edgeVertex
				}
			}

			for _, edge := range circuitEdges {
				start := edge.GetStart()
				if start == extClosest.Edge.GetStart() || start == extClosest.Edge.GetEnd() {
					continue
				}
				if extVertex.DistanceTo(start) < closestEdgeVertex.DistanceTo(start) {
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

	return circuitEdges, unattachedVertices
}

func farthestVertexFrom(v model.CircuitVertex) model.CircuitVertex {
	farthestPoint := v
	farthestDist := 0.0

	for _, edge := range v.(*GraphVertex).GetPaths() {
		if edge.GetLength() > farthestDist {
			farthestDist = edge.GetLength()
			farthestPoint = edge.GetEnd()
		}
	}

	return farthestPoint
}

func findClosestEdges(circuitEdges []model.CircuitEdge, vertices map[model.CircuitVertex]bool) map[model.CircuitVertex]*model.DistanceToEdge {
	closestEdges := make(map[model.CircuitVertex]*model.DistanceToEdge)
	for v := range vertices {
		closestEdges[v] = &model.DistanceToEdge{
			Vertex:   v,
			Edge:     nil,
			Distance: math.MaxFloat64,
		}
		for _, edge := range circuitEdges {
			if dist := edge.DistanceIncrease(v); dist < closestEdges[v].Distance {
				closestEdges[v].Edge = edge
				closestEdges[v].Distance = dist
			}
		}
	}
	return closestEdges
}

func insertVertex(circuitEdges []model.CircuitEdge, v model.CircuitVertex, edgeToSplit model.CircuitEdge) []model.CircuitEdge {
	edgeA, edgeB := edgeToSplit.GetStart().EdgeTo(v), v.EdgeTo(edgeToSplit.GetEnd())

	edgeIndex := model.IndexOfEdge(circuitEdges, edgeToSplit)
	circuitEdges = append(circuitEdges[:edgeIndex+1], circuitEdges[edgeIndex:]...)
	circuitEdges[edgeIndex] = edgeA
	circuitEdges[edgeIndex+1] = edgeB

	return circuitEdges
}
