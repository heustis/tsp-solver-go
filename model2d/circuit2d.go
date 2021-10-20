package model2d

import (
	"math"

	"github.com/fealos/lee-tsp-go/model"
)

type Circuit2D struct {
	Vertices           []*Vertex2D
	midpoint           *Vertex2D
	circuit            []model.CircuitVertex
	circuitEdges       []model.CircuitEdge
	closestEdges       map[model.CircuitVertex]*distanceToEdge
	interiorVertices   map[model.CircuitVertex]bool
	unattachedVertices map[model.CircuitVertex]bool
}

type distanceToEdge struct {
	edge     model.CircuitEdge
	distance float64
}

func (c *Circuit2D) BuildPerimiter() {
	// 1. Find point farthest from midpoint
	// Restricts problem-space to a circle around the midpoint, with radius equal to the distance to the point.
	farthestFromMid := findFarthestPoint(c.midpoint, c.Vertices)
	delete(c.unattachedVertices, farthestFromMid)
	c.circuit = append(c.circuit, farthestFromMid)

	// 2. Find point farthest from point in step 1.
	// Restricts problem-space to intersection of step 1 circle,
	// and a circle centered on the point from step 1 with a radius equal to the distance between the points found in step 1 and 2.
	farthestFromFarthest := findFarthestPoint(farthestFromMid, c.Vertices)
	delete(c.unattachedVertices, farthestFromFarthest)
	c.circuit = append(c.circuit, farthestFromFarthest)

	// 3. Created edges 1 -> 2 and 2 -> 1
	c.circuitEdges = append(c.circuitEdges, NewEdge2D(farthestFromMid, farthestFromFarthest))
	c.circuitEdges = append(c.circuitEdges, NewEdge2D(farthestFromFarthest, farthestFromMid))

	// 4. Find the exterior point farthest from its closest edge; 3rd point only.
	// For the third point, we can simplify this since both edges are the same (but flipped).
	// Also, the third point determines whether our vertices are ordered clockwise or counter-clockwise.
	// For this algorithm we will use counter-clockwise ordering, meaning the exterior points will be to the right of their closest edge (while the perimeter is convex).

	var farthestFromEdges *Vertex2D
	farthestFromEdgesDistance := 0.0

	for vertex := range c.unattachedVertices {
		v2d := vertex.(*Vertex2D)
		e2d := c.circuitEdges[0].(*Edge2D)
		if v2d.isLeftOf(e2d) {
			e2d = c.circuitEdges[1].(*Edge2D)
		}

		dist := v2d.distanceToEdge(e2d)
		if dist > farthestFromEdgesDistance {
			farthestFromEdgesDistance = dist
			farthestFromEdges = v2d
		}

		c.closestEdges[vertex] = &distanceToEdge{
			edge:     e2d,
			distance: dist,
		}
	}

	edgeToSplit := c.closestEdges[farthestFromEdges].edge
	delete(c.closestEdges, farthestFromEdges)
	c.addVertexToEdge(farthestFromEdges, edgeToSplit, c.updateExteriorPoints)

	// 5. Find the exterior point farthest from its closest edge; 4th point and beyond.
	// Split the closest edge by adding the point to it, and consequently to the perimeter.
	// Check all remaining exterior points to see if they are now interior points, and update the model as appropriate.
	// Repeat until all points are interior or perimeter points.
	for len(c.closestEdges) > 0 {
		var pointFarthestFromClosestEdge model.CircuitVertex
		farthestDistanceFromClosestEdge := 0.0

		for v, dist := range c.closestEdges {
			if dist.distance > farthestDistanceFromClosestEdge {
				pointFarthestFromClosestEdge = v
				farthestDistanceFromClosestEdge = dist.distance
			}
		}

		e := c.closestEdges[pointFarthestFromClosestEdge].edge
		delete(c.closestEdges, pointFarthestFromClosestEdge)
		c.addVertexToEdge(pointFarthestFromClosestEdge, e, c.updateExteriorPoints)
	}

	// 6. Find the closest edge for all interior points, based on distance increase (rather than perpendicular distance)
	for vertex := range c.unattachedVertices {
		closest := vertex.FindClosestEdge(c.circuitEdges)
		c.closestEdges[vertex] = &distanceToEdge{
			edge:     closest,
			distance: closest.DistanceIncrease(vertex),
		}
	}
}

func (c *Circuit2D) FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge) {
	var nextVertex model.CircuitVertex = nil
	nextEdge := &distanceToEdge{distance: math.MaxFloat64}

	for vertex := range c.unattachedVertices {
		edge := c.closestEdges[vertex]
		if edge.distance < nextEdge.distance {
			nextVertex = vertex
			nextEdge = edge
		}
	}
	return nextVertex, nextEdge.edge
}

func (c *Circuit2D) GetAttachedVertices() []model.CircuitVertex {
	return c.circuit
}

func (c *Circuit2D) GetInteriorVertices() map[model.CircuitVertex]bool {
	return c.interiorVertices
}

func (c *Circuit2D) GetLength() float64 {
	length := 0.0
	for _, edge := range c.circuitEdges {
		length += edge.GetLength()
	}
	return length
}

func (c *Circuit2D) GetUnattachedVertices() map[model.CircuitVertex]bool {
	return c.unattachedVertices
}

func (c *Circuit2D) Prepare() {
	c.interiorVertices = make(map[model.CircuitVertex]bool)
	c.unattachedVertices = make(map[model.CircuitVertex]bool)
	c.closestEdges = make(map[model.CircuitVertex]*distanceToEdge)
	c.circuit = []model.CircuitVertex{}
	c.circuitEdges = []model.CircuitEdge{}

	c.Vertices = deduplicateVertices(c.Vertices)

	numVertices := float64(len(c.Vertices))
	c.midpoint = &Vertex2D{0.0, 0.0}

	for _, v := range c.Vertices {
		c.unattachedVertices[v] = true
		c.midpoint.X += v.X / numVertices
		c.midpoint.Y += v.Y / numVertices
	}
}

func (c *Circuit2D) Update(vertexToAdd model.CircuitVertex, edgeToSplit model.CircuitEdge) {
	if vertexToAdd != nil {
		c.addVertexToEdge(vertexToAdd, edgeToSplit, c.updateInteriorPoints)
	}
}

func (c *Circuit2D) addVertexToEdge(vertexToAdd model.CircuitVertex, edgeToSplit model.CircuitEdge, updateFunc func(previousEdge model.CircuitEdge, edgeA model.CircuitEdge, edgeB model.CircuitEdge)) {
	var edgeIndex int
	c.circuitEdges, edgeIndex = model.SplitEdge(c.circuitEdges, edgeToSplit, vertexToAdd)
	if edgeIndex >= 0 {
		c.insertVertex(edgeIndex+1, vertexToAdd)
		delete(c.unattachedVertices, vertexToAdd)
		updateFunc(edgeToSplit, c.circuitEdges[edgeIndex], c.circuitEdges[edgeIndex+1])
	}
}

func (c *Circuit2D) indexOfAttachedVertex(vertex model.CircuitVertex) int {
	for index, v := range c.circuit {
		if v == vertex {
			return index
		}
	}
	return -1
}

func (c *Circuit2D) insertVertex(index int, vertex model.CircuitVertex) {
	if index >= len(c.circuit) {
		c.circuit = append(c.circuit, vertex)
	} else {
		// copy all elements starting at the index one to the right to create a duplicate record at index and index+1.
		c.circuit = append(c.circuit[:index+1], c.circuit[index:]...)
		// update only the vertex at the index, so that there are no duplicates and the vertex is at the index.
		c.circuit[index] = vertex
	}
}

func (c *Circuit2D) updateExteriorPoints(previousEdge model.CircuitEdge, edgeA model.CircuitEdge, edgeB model.CircuitEdge) {
	for v := range c.unattachedVertices {
		// If the vertex was previously an exterior point and the edge closest to it was split, it could now be an interior point.
		if closest, wasExterior := c.closestEdges[v]; wasExterior && closest.edge == previousEdge {
			var newClosest *distanceToEdge
			if distA, distB := edgeA.DistanceIncrease(v), edgeB.DistanceIncrease(v); distA < distB {
				newClosest = &distanceToEdge{
					edge:     edgeA,
					distance: distA,
				}
			} else {
				newClosest = &distanceToEdge{
					edge:     edgeB,
					distance: distB,
				}
			}

			// If the vertex is now interior, stop tracking its closest edge (until the convex perimeter is fully constructed) and add it to the interior edge list.
			// Otherwise, it is still exterior, so update its closest edge.
			if v.(*Vertex2D).isLeftOf(newClosest.edge.(*Edge2D)) {
				delete(c.closestEdges, v)
				c.interiorVertices[v] = true
			} else {
				c.closestEdges[v] = newClosest
			}
		}
	}
}

func (c *Circuit2D) updateInteriorPoints(previousEdge model.CircuitEdge, edgeA model.CircuitEdge, edgeB model.CircuitEdge) {
	for vertex := range c.interiorVertices {
		if edgeA.GetEnd() == vertex {
			// Don't update the vertex that was just added to the edge.
			continue
		} else if edgeA.GetStart() == vertex || edgeB.GetEnd() == vertex {
			// Update the closest edge of the start and end vertices, if they are interior, so that they account for the newly added vertex.
			vertexIndex := c.indexOfAttachedVertex(vertex)
			interiorLen := len(c.circuit)
			previousVertex := c.circuit[(vertexIndex+interiorLen-1)%interiorLen]
			nextVertex := c.circuit[(vertexIndex+1)%interiorLen]
			closest := NewEdge2D(previousVertex.(*Vertex2D), nextVertex.(*Vertex2D))
			c.closestEdges[vertex] = &distanceToEdge{
				edge:     closest,
				distance: closest.DistanceIncrease(vertex),
			}
			continue
		}

		// Update the closest edge for all interior vertices, including attached vertices.
		wasUpdated := false
		previousClosest := c.closestEdges[vertex]
		distA := edgeA.DistanceIncrease(vertex.(*Vertex2D))
		distB := edgeB.DistanceIncrease(vertex.(*Vertex2D))
		if distA < previousClosest.distance && distA <= distB {
			c.closestEdges[vertex] = &distanceToEdge{
				edge:     edgeA,
				distance: distA,
			}
			wasUpdated = true
		} else if distB < previousClosest.distance {
			c.closestEdges[vertex] = &distanceToEdge{
				edge:     edgeB,
				distance: distB,
			}
			wasUpdated = true
		}

		// If the point was already attached to the model, and one of the new edges is closer:
		// need to undo its current attachement and re-add it to the unattached vertices so the it can be attached to the new closer edge.
		if _, isUnattached := c.unattachedVertices[vertex]; wasUpdated && !isUnattached {
			c.unattachedVertices[vertex] = true

			// detach vertex - replace both previous edges with a single edge
			index := c.indexOfAttachedVertex(vertex)
			c.circuit = model.DeleteVertex(c.circuit, index)

			var detachedEdgeA model.CircuitEdge
			var detachedEdgeB model.CircuitEdge
			c.circuitEdges, detachedEdgeA, detachedEdgeB = model.MergeEdges(c.circuitEdges, index)

			// Update the closest edge of any vertices with either of the detached edges as their previous closest edge.
			for v := range c.interiorVertices {
				e := c.closestEdges[v]
				if e.edge == detachedEdgeA || e.edge == detachedEdgeB {
					newClosest := v.FindClosestEdge(c.circuitEdges)
					c.closestEdges[v] = &distanceToEdge{
						edge:     newClosest,
						distance: newClosest.DistanceIncrease(v),
					}
				}
			}
		}
	}
}

var _ model.Circuit = (*Circuit2D)(nil)
