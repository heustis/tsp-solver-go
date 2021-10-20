package model2d

import "github.com/fealos/lee-tsp-go/model"

type PerimeterBuilder2D struct {
}

func (builder *PerimeterBuilder2D) BuildPerimiter(verticesArg []*Vertex2D) ([]model.CircuitVertex, []model.CircuitEdge, map[model.CircuitVertex]bool) {
	numVertices := len(verticesArg)
	vertices := make([]*Vertex2D, numVertices)
	midpoint := NewVertex2D(0, 0)
	unattachedVertices := make(map[model.CircuitVertex]bool)
	circuit := make([]model.CircuitVertex, 0, numVertices)
	circuitEdges := make([]model.CircuitEdge, 0, numVertices)

	for i, v := range verticesArg {
		// v2d := v.(*Vertex2D)
		vertices[i] = v
		unattachedVertices[v] = true
		midpoint.X += v.X / float64(numVertices)
		midpoint.Y += v.Y / float64(numVertices)
	}

	// 1. Find point farthest from midpoint
	// Restricts problem-space to a circle around the midpoint, with radius equal to the distance to the point.
	farthestFromMid := findFarthestPoint(midpoint, vertices)
	delete(unattachedVertices, farthestFromMid)
	circuit = append(circuit, farthestFromMid)

	// 2. Find point farthest from point in step 1.
	// Restricts problem-space to intersection of step 1 circle,
	// and a circle centered on the point from step 1 with a radius equal to the distance between the points found in step 1 and 2.
	farthestFromFarthest := findFarthestPoint(farthestFromMid, vertices)
	delete(unattachedVertices, farthestFromFarthest)
	circuit = append(circuit, farthestFromFarthest)

	// 3. Created edges 1 -> 2 and 2 -> 1
	circuitEdges = append(circuitEdges, NewEdge2D(farthestFromMid, farthestFromFarthest))
	circuitEdges = append(circuitEdges, NewEdge2D(farthestFromFarthest, farthestFromMid))

	// 4. Initialize the closestEdges map which will be used to find the exterior point farthest from its closest edge.
	// For the third point only, we can simplify this since both edges are the same (but flipped).
	// When the third point is inserted it will determine whether our vertices are ordered clockwise or counter-clockwise.
	// For this algorithm we will use counter-clockwise ordering, meaning the exterior points will be to the right of their closest edge (while the perimeter is convex).

	exteriorClosestEdges := make(map[model.CircuitVertex]*heapDistanceToEdge)

	for vertex := range unattachedVertices {
		v2d := vertex.(*Vertex2D)
		e2d := circuitEdges[0].(*Edge2D)
		if v2d.isLeftOf(e2d) {
			e2d = circuitEdges[1].(*Edge2D)
		}

		exteriorClosestEdges[vertex] = &heapDistanceToEdge{
			edge:     e2d,
			distance: v2d.distanceToEdge(e2d),
			vertex:   v2d,
		}
	}

	// 5. Find the exterior point farthest from its closest edge.
	// Split the closest edge by adding the point to it, and consequently to the perimeter.
	// Check all remaining exterior points to see if they are now interior points, and update the model as appropriate.
	// Repeat until all points are interior or perimeter points.
	// Complexity: This step in O(N^2) because it iterates once per vertex in the concave perimeter (N iterations) and for each of those iterations it:
	//             1. looks at each exterior point to find farthest from its closest point (O(N)); and then
	//             2. updates each exterior point that had the split edge as its closest edge (O(N)).
	for len(exteriorClosestEdges) > 0 {
		farthestFromClosestEdge := &heapDistanceToEdge{
			distance: 0.0,
		}
		for _, closest := range exteriorClosestEdges {
			if closest.distance > farthestFromClosestEdge.distance {
				farthestFromClosestEdge = closest
			}
		}

		var edgeIndex int
		circuitEdges, edgeIndex = model.SplitEdge(circuitEdges, farthestFromClosestEdge.edge, farthestFromClosestEdge.vertex)
		insertVertex(circuit, edgeIndex+1, farthestFromClosestEdge.vertex)
		delete(unattachedVertices, farthestFromClosestEdge.vertex)
		delete(exteriorClosestEdges, farthestFromClosestEdge.vertex)

		edgeA, edgeB := circuitEdges[edgeIndex], circuitEdges[edgeIndex+1]

		for v := range unattachedVertices {
			// If the vertex was previously an exterior point and the edge closest to it was split, it could now be an interior point.
			if closest, wasExterior := exteriorClosestEdges[v]; wasExterior && closest.edge == farthestFromClosestEdge.edge {
				var newClosest *heapDistanceToEdge
				if distA, distB := edgeA.DistanceIncrease(v), edgeB.DistanceIncrease(v); distA < distB {
					newClosest = &heapDistanceToEdge{
						edge:     edgeA,
						distance: distA,
						vertex:   v.(*Vertex2D),
					}
				} else {
					newClosest = &heapDistanceToEdge{
						edge:     edgeB,
						distance: distB,
						vertex:   v.(*Vertex2D),
					}
				}

				// If the vertex is now interior, stop tracking its closest edge (until the convex perimeter is fully constructed) and add it to the interior edge list.
				// Otherwise, it is still exterior, so update its closest edge.
				if v.(*Vertex2D).isLeftOf(newClosest.edge.(*Edge2D)) {
					delete(exteriorClosestEdges, v)
				} else {
					exteriorClosestEdges[v] = newClosest
				}
			}
		}
	}

	return circuit, circuitEdges, unattachedVertices
}

func insertVertex(circuit []model.CircuitVertex, index int, vertex model.CircuitVertex) []model.CircuitVertex {
	if index >= len(circuit) {
		return append(circuit, vertex)
	} else {
		// copy all elements starting at the index one to the right to create a duplicate record at index and index+1.
		circuit = append(circuit[:index+1], circuit[index:]...)
		// update only the vertex at the index, so that there are no duplicates and the vertex is at the index.
		circuit[index] = vertex
		return circuit
	}
}
