package model3d

import (
	"math"

	"github.com/fealos/lee-tsp-go/model"
)

type PerimeterBuilder3D struct {
}

func (builder *PerimeterBuilder3D) BuildPerimiter(verticesArg []model.CircuitVertex) ([]model.CircuitEdge, map[model.CircuitVertex]bool) {
	numVertices := len(verticesArg)
	vertices := make([]*Vertex3D, numVertices)
	midpoint := NewVertex3D(0, 0, 0)
	unattachedVertices := make(map[model.CircuitVertex]bool)
	distanceToMidpoint := make(map[model.CircuitVertex]float64)
	exteriorClosestEdges := make(map[model.CircuitVertex]*model.DistanceToEdge)
	circuitEdges := make([]model.CircuitEdge, 0, numVertices)

	for i, v := range verticesArg {
		v3d := v.(*Vertex3D)
		vertices[i] = v3d
		unattachedVertices[v] = true
		midpoint.X += v3d.X / float64(numVertices)
		midpoint.Y += v3d.Y / float64(numVertices)
		midpoint.Z += v3d.Z / float64(numVertices)
	}

	// 1. Find point farthest from midpoint
	// Restricts problem-space to a sphere around the midpoint, with radius equal to the distance to the point.
	farthestFromMid := findFarthestPoint3D(midpoint, vertices)
	delete(unattachedVertices, farthestFromMid)
	distanceToMidpoint[farthestFromMid] = farthestFromMid.DistanceTo(midpoint)

	// 2. Find point farthest from point in step 1.
	// Restricts problem-space to intersection of step 1 sphere,
	// and a sphere centered on the point from step 1 with a radius equal to the distance between the points found in step 1 and 2.
	farthestFromFarthest := findFarthestPoint3D(farthestFromMid, vertices)
	delete(unattachedVertices, farthestFromFarthest)
	distanceToMidpoint[farthestFromFarthest] = farthestFromFarthest.DistanceTo(midpoint)

	// 3. Created edges 1 -> 2 and 2 -> 1
	circuitEdges = append(circuitEdges, NewEdge3D(farthestFromMid, farthestFromFarthest))
	circuitEdges = append(circuitEdges, NewEdge3D(farthestFromFarthest, farthestFromMid))

	// 4. Initialize the closestEdges map which will be used to find the exterior point farthest from its closest edge.
	// For the third point only, we can simplify this since both edges are the same (but flipped).
	// When the third point is inserted it will determine whether our vertices are ordered clockwise or counter-clockwise.
	// For the 3D version of this algorithm we will not deliberately select CW vs CCW.
	for vertex := range unattachedVertices {
		v3d := vertex.(*Vertex3D)
		e3d := circuitEdges[0].(*Edge3D)

		exteriorClosestEdges[vertex] = &model.DistanceToEdge{
			Edge:     e3d,
			Distance: v3d.distanceToEdge(e3d),
			Vertex:   v3d,
		}

		distanceToMidpoint[v3d] = v3d.DistanceTo(midpoint)
	}

	// 5. Find the exterior point farthest from its closest edge.
	// Split the closest edge by adding the point to it, and consequently to the perimeter.
	// Check all remaining exterior points to see if they are now interior points, and update the model as appropriate.
	// Repeat until all points are interior or perimeter points.
	// Complexity: This step in O(N^2) because it iterates once per vertex in the concave perimeter (N iterations) and for each of those iterations it:
	//             1. looks at each exterior point to find farthest from its closest point (O(N)); and then
	//             2. updates each exterior point that had the split edge as its closest edge (O(N)).
	for len(exteriorClosestEdges) > 0 {
		farthestFromClosestEdge := &model.DistanceToEdge{
			Distance: 0.0,
		}
		for _, closest := range exteriorClosestEdges {
			if closest.Distance > farthestFromClosestEdge.Distance {
				farthestFromClosestEdge = closest
			}
		}

		var edgeIndex int
		circuitEdges, edgeIndex = model.SplitEdge(circuitEdges, farthestFromClosestEdge.Edge, farthestFromClosestEdge.Vertex)
		delete(unattachedVertices, farthestFromClosestEdge.Vertex)
		delete(exteriorClosestEdges, farthestFromClosestEdge.Vertex)

		// Check if either adjacent point on the perimeter is now interior, and detach it if it is. This is required for 3D but not 2D.
		if circuitLen := len(circuitEdges); circuitLen > 3 {
			//TODO - Debug and fix
			start, added, end := farthestFromClosestEdge.Edge.GetStart(), farthestFromClosestEdge.Vertex, farthestFromClosestEdge.Edge.GetEnd()
			beforeStart, afterEnd := circuitEdges[(edgeIndex-1+circuitLen)%circuitLen].GetStart(), circuitEdges[(edgeIndex+2)%circuitLen].GetEnd()
			if isInterior(start, beforeStart.EdgeTo(added), midpoint, distanceToMidpoint[start]) {
				circuitEdges, _, _, _ = model.MergeEdges2(circuitEdges, start)
				unattachedVertices[start] = true
			}
			if isInterior(end, added.EdgeTo(afterEnd), midpoint, distanceToMidpoint[end]) {
				circuitEdges, _, _, _ = model.MergeEdges2(circuitEdges, end)
				unattachedVertices[end] = true
			}
		}

		for v := range unattachedVertices {
			// If the vertex was previously an exterior point and the edge closest to it was split, it could now be an interior point.
			if closest, wasExterior := exteriorClosestEdges[v]; wasExterior {
				// Update the closest edge to this point. For 3 dimensions, need to check all edges.
				v3d := v.(*Vertex3D)
				newClosest := &model.DistanceToEdge{
					Vertex:   v3d,
					Distance: math.MaxFloat64,
				}
				for _, edge := range circuitEdges {
					e3d := edge.(*Edge3D)
					dist := v3d.distanceToEdge(e3d)
					if dist < newClosest.Distance {
						newClosest.Distance = dist
						newClosest.Edge = e3d
					}
				}

				// If the vertex is now interior, stop tracking its closest edge (until the convex perimeter is fully constructed) and add it to the interior edge list.
				// Otherwise, it is still exterior, so update its closest edge.
				if newClosest.Edge != closest.Edge && isInterior(v, newClosest.Edge, midpoint, distanceToMidpoint[v]) {
					delete(exteriorClosestEdges, v)
				} else {
					exteriorClosestEdges[v] = newClosest
				}
			}
		}
	}

	return circuitEdges, unattachedVertices
}

func findFarthestPoint3D(target *Vertex3D, points []*Vertex3D) *Vertex3D {
	var farthestPoint *Vertex3D
	farthestDistance := 0.0

	for _, point := range points {
		if distance := point.DistanceTo(target); distance > farthestDistance {
			farthestDistance = distance
			farthestPoint = point
		}
	}

	return farthestPoint
}

func isInterior(v model.CircuitVertex, closestEdge model.CircuitEdge, midpoint *Vertex3D, distanceToMidpoint float64) bool {
	projected := v.(*Vertex3D).projectToEdge(closestEdge.(*Edge3D))
	projectedDist := projected.DistanceTo(midpoint)
	return projectedDist > distanceToMidpoint
}

var _ model.PerimeterBuilder = (*PerimeterBuilder3D)(nil)
