package experimental

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"

	"github.com/fealos/lee-tsp-go/model"
)

// ConvexConcaveWeightedEdges is significantly worse than the other greedy algorithms, see `results_2d_comp_greedy_3.tsv`.
// I tested it with 8, 4, and 1 points in the weighting array (see below).
// With 1 point this produced the same results as `circuitgreedy_impl.go`, which was expected since weighing only one point is the same as only considering which point is closest to its closest edge.
type ConvexConcaveWeightedEdges struct {
	Vertices           []model.CircuitVertex
	deduplicator       func([]model.CircuitVertex) []model.CircuitVertex
	perimeterBuilder   model.PerimeterBuilder
	circuitEdges       []model.CircuitEdge
	closestVertices    map[model.CircuitEdge]*weightedEdge
	unattachedVertices map[model.CircuitVertex]bool
	length             float64
}

func NewConvexConcaveWeightedEdges(vertices []model.CircuitVertex, deduplicator func([]model.CircuitVertex) []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) model.Circuit {
	return &ConvexConcaveWeightedEdges{
		Vertices:         vertices,
		deduplicator:     deduplicator,
		perimeterBuilder: perimeterBuilder,
	}
}

func (c *ConvexConcaveWeightedEdges) BuildPerimiter() {
	c.circuitEdges, c.unattachedVertices = c.perimeterBuilder.BuildPerimiter(c.Vertices)

	// Find the closest edge for all interior points, based on distance increase; store them in a heap for retrieval from closest to farthest.
	c.length = 0.0
	for _, edge := range c.circuitEdges {
		c.closestVertices[edge] = newWeightedEdge(edge, c.unattachedVertices)
		c.length += edge.GetLength()
	}
}

func (c *ConvexConcaveWeightedEdges) FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge) {
	var closestEdge model.CircuitEdge
	closestVertices := &weightedEdge{
		weightedDistance: math.MaxFloat64,
	}
	for e, w := range c.closestVertices {
		if w.weightedDistance < closestVertices.weightedDistance {
			closestVertices = w
			closestEdge = e
		}
	}
	if len(closestVertices.closestVertices) == 0 {
		return nil, nil
	} else {
		return closestVertices.closestVertices[0].vertex, closestEdge
	}
}

func (c *ConvexConcaveWeightedEdges) GetAttachedEdges() []model.CircuitEdge {
	return c.circuitEdges
}

func (c *ConvexConcaveWeightedEdges) GetAttachedVertices() []model.CircuitVertex {
	vertices := make([]model.CircuitVertex, len(c.circuitEdges))
	for i, edge := range c.circuitEdges {
		vertices[i] = edge.GetStart()
	}
	return vertices
}

func (c *ConvexConcaveWeightedEdges) GetClosestVertices() map[model.CircuitEdge]*weightedEdge {
	return c.closestVertices
}

func (c *ConvexConcaveWeightedEdges) GetLength() float64 {
	return c.length
}

func (c *ConvexConcaveWeightedEdges) GetUnattachedVertices() map[model.CircuitVertex]bool {
	return c.unattachedVertices
}

func (c *ConvexConcaveWeightedEdges) Prepare() {
	c.unattachedVertices = make(map[model.CircuitVertex]bool)
	c.closestVertices = make(map[model.CircuitEdge]*weightedEdge)
	c.circuitEdges = []model.CircuitEdge{}
	c.length = 0.0

	c.Vertices = c.deduplicator(c.Vertices)
}

func (c *ConvexConcaveWeightedEdges) Update(vertexToAdd model.CircuitVertex, edgeToSplit model.CircuitEdge) {
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
		delete(c.closestVertices, edgeToSplit)

		for e, w := range c.closestVertices {
			w.removeVertex(vertexToAdd, e, c.unattachedVertices)
		}
		edgeA, edgeB := c.circuitEdges[edgeIndex], c.circuitEdges[edgeIndex+1]
		c.closestVertices[edgeA] = newWeightedEdge(edgeA, c.unattachedVertices)
		c.closestVertices[edgeB] = newWeightedEdge(edgeB, c.unattachedVertices)

		c.length += edgeA.GetLength() + edgeB.GetLength() - edgeToSplit.GetLength()
	}
}

var weights = [8]float64{0.5, 0.25, 0.125, 1.0 / 16.0, 1.0 / 32.0, 1.0 / 64.0, 1.0 / 128.0, 1.0 / 256.0}

// var weights = [4]float64{0.9, 0.09, 0.009, .0009}
// var weights = [1]float64{1.0}

type weightedVertex struct {
	distance float64
	vertex   model.CircuitVertex
}

type weightedEdge struct {
	weightedDistance float64
	closestVertices  []*weightedVertex
}

func (w *weightedEdge) GetClosestVertices() []*weightedVertex {
	return w.closestVertices
}

func (w *weightedEdge) GetDistance() float64 {
	return w.weightedDistance
}

func newWeightedEdge(edge model.CircuitEdge, unattachedVertices map[model.CircuitVertex]bool) *weightedEdge {
	lenClosest := int(math.Min(float64(len(weights)), float64(len(unattachedVertices))))
	lastIndex := lenClosest - 1
	w := &weightedEdge{
		weightedDistance: 0.0,
		closestVertices:  make([]*weightedVertex, lenClosest),
	}
	nextIndex := 0
	for v := range unattachedVertices {
		vertexDistance := edge.DistanceIncrease(v)
		// The first n vertices can be added directly to the array, then sorted once all are added.
		if nextIndex < lenClosest {
			w.closestVertices[nextIndex] = &weightedVertex{
				vertex:   v,
				distance: vertexDistance,
			}
			if nextIndex == lastIndex {
				sort.Slice(w.closestVertices, func(i, j int) bool {
					return w.closestVertices[i].distance < w.closestVertices[j].distance
				})
			}
			nextIndex++
		} else if vertexDistance < w.closestVertices[lastIndex].distance {
			w.closestVertices[lastIndex] = &weightedVertex{
				vertex:   v,
				distance: vertexDistance,
			}
			// Bubbling is normally too inefficient for sorting, but this array has a max of 8 entries so it isn't too impactful.
			for i, j := lastIndex, lastIndex-1; j >= 0; i, j = j, j-1 {
				// Stop bubbling once this vertex is farther than the next vertex in the array ("next" meaning closer to index 0).
				if vertexDistance > w.closestVertices[j].distance {
					break
				}
				w.closestVertices[i], w.closestVertices[j] = w.closestVertices[j], w.closestVertices[i]
			}
		}
	}
	w.updateDistance()
	return w
}

func (w *weightedEdge) removeVertex(vertex model.CircuitVertex, edge model.CircuitEdge, unattachedVertices map[model.CircuitVertex]bool) {
	numVertices := len(w.closestVertices)
	vertexIndex := numVertices - 1
	for ; vertexIndex >= 0 && w.closestVertices[vertexIndex].vertex != vertex; vertexIndex-- {
	}

	if vertexIndex < 0 {
		return
	}

	// If there are unattached vertices that are not included in the weighted average, add the next closest vertex into the average.
	if len(unattachedVertices) >= numVertices {
		nextClosest := &weightedVertex{
			distance: math.MaxFloat64,
		}
		for v := range unattachedVertices {
			if dist := edge.DistanceIncrease(v); dist < nextClosest.distance {
				nextClosest.vertex = v
				nextClosest.distance = dist
			}
		}

		w.closestVertices[vertexIndex] = nextClosest
		// Bubble the newly added vertex to the last position in the array, because
		// a newly added vertex will be farther away than any vertex already in the list (due to newWeightedEdge selecting the 8 closest vertices).
		for i, j := vertexIndex, vertexIndex+1; j < numVertices; i, j = j, j+1 {
			w.closestVertices[i], w.closestVertices[j] = w.closestVertices[j], w.closestVertices[i]
		}
	} else if vertexIndex == 0 {
		w.closestVertices = w.closestVertices[1:]
	} else if vertexIndex == numVertices-1 {
		w.closestVertices = w.closestVertices[:vertexIndex]
	} else {
		w.closestVertices = append(w.closestVertices[:vertexIndex], w.closestVertices[vertexIndex+1:]...)
	}
	w.updateDistance()
}

func (w *weightedEdge) updateDistance() {
	w.weightedDistance = 0.0
	for i, v := range w.closestVertices {
		w.weightedDistance += v.distance * weights[i]
	}
}

var _ model.Circuit = (*ConvexConcaveWeightedEdges)(nil)
