package circuit

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/stats"
)

type confidenceCircuit struct {
	edges     []model.CircuitEdge
	distances map[model.CircuitVertex]*stats.DistanceGaps
	length    float64
}

func (c *confidenceCircuit) attachVertex(distance *model.DistanceToEdge) {
	var edgeIndex int
	c.edges, edgeIndex = model.SplitEdgeCopy(c.edges, distance.Edge, distance.Vertex)
	if edgeIndex < 0 {
		expectedEdgeJson, _ := json.Marshal(distance.Edge)
		actualCircuitJson, _ := json.Marshal(c.edges)
		panic(fmt.Errorf("edge not found in circuit=%p, expected=%s, \ncircuit=%s", c, string(expectedEdgeJson), string(actualCircuitJson)))
	}
	c.length += distance.Distance
	edgeA, edgeB := c.edges[edgeIndex], c.edges[edgeIndex+1]
	for _, stats := range c.distances {
		stats.UpdateStats(distance.Edge, edgeA, edgeB)
	}
}

func (c *confidenceCircuit) clone() *confidenceCircuit {
	clone := &confidenceCircuit{
		edges:     make([]model.CircuitEdge, len(c.edges)),
		distances: make(map[model.CircuitVertex]*stats.DistanceGaps),
		length:    c.length,
	}
	copy(clone.edges, c.edges)

	for k, v := range c.distances {
		clone.distances[k] = v.Clone()
	}

	return clone
}

func (c *confidenceCircuit) findNext(significance float64) []*model.DistanceToEdge {
	// If there is only one vertex left to attach, attach it to its closest edge.
	if len(c.distances) == 1 {
		for _, stats := range c.distances {
			return stats.ClosestEdges[0:1]
		}
	}

	var vertexToUpdate model.CircuitVertex
	var closestVertex *model.DistanceToEdge

	// Find the most significant early gap to determine which vertex to attach to which edge (or edges).
	// Prioritize earlier significant gaps over later, but more significant, gaps (e.g. a gap with a Z-score of 3.5 at index 1 should be prioritized over a gap with a Z-score of 5 at index 2).
	gapIndex := math.MaxInt64
	gapSignificance := 0.0

	for v, stats := range c.distances {
		// Track the vertex closest to its nearest edge, in the event there are no significant gaps.
		if closestVertex == nil || stats.ClosestEdges[0].Distance < closestVertex.Distance {
			closestVertex = stats.ClosestEdges[0]
		}
		// Determine if the current vertex has a significant gap in its edge distances that is:
		// earlier than the current best, or more significant at the same index.
		for i, currentGap := range stats.Gaps {
			if i > gapIndex {
				break
			} else if currentSignificance := (currentGap - stats.GapAverage) / stats.GapStandardDeviation; currentSignificance < significance {
				// Note: do not use the absolute value for this computation, as we only want significantly large gaps, not significantly small gaps.
				continue
			} else if currentSignificance > gapSignificance || i < gapIndex {
				vertexToUpdate = v
				gapIndex = i
				gapSignificance = currentSignificance
			}
		}
	}

	// If all vertices lack significant gaps, select the vertex with the closest edge and clone the circuit once for each edge.
	if vertexToUpdate == nil {
		return c.distances[closestVertex.Vertex].ClosestEdges
	}

	return c.distances[vertexToUpdate].ClosestEdges[0 : gapIndex+1]
}

func (c *confidenceCircuit) update(significance float64) (clones []*confidenceCircuit) {
	next := c.findNext(significance)

	delete(c.distances, next[0].Vertex)

	if numClones := len(next) - 1; numClones > 0 {
		clones = make([]*confidenceCircuit, numClones)
		for i, cloneDistance := range next {
			if cloneIndex := i - 1; cloneIndex >= 0 {
				clones[cloneIndex] = c.clone()
				clones[cloneIndex].attachVertex(cloneDistance)
			}
		}
	} else {
		clones = nil
	}

	// In either case update this circuit with the first entry - this must happen after cloning,
	c.attachVertex(next[0])

	return clones
}
