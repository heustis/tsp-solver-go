package tspstats

import (
	"fmt"
	"math"
	"sort"

	"github.com/fealos/lee-tsp-go/tspmodel"
)

// DistanceGaps tracks the distance between one vertex and each edge in the current circuit.
// From those distances, it analyzes the gaps between successive distances, sorted from smallest to largest,
// to facilitate finding statistically significant gaps, which indicates that the vertex will likely be attached to an edge prior to that gap.
type DistanceGaps struct {
	ClosestEdges         []*tspmodel.DistanceToEdge
	Gaps                 []float64
	GapAverage           float64
	GapStandardDeviation float64
}

// NewDistanceGaps accepts a vertex and the edges of the current circuit, and returnes a populated DistanceGaps.
func NewDistanceGaps(vertex tspmodel.CircuitVertex, edges []tspmodel.CircuitEdge) *DistanceGaps {
	stats := &DistanceGaps{
		ClosestEdges: make([]*tspmodel.DistanceToEdge, len(edges)),
	}
	for i, e := range edges {
		stats.ClosestEdges[i] = &tspmodel.DistanceToEdge{
			Vertex:   vertex,
			Edge:     e,
			Distance: e.DistanceIncrease(vertex),
		}
	}
	sort.Slice(stats.ClosestEdges, func(i, j int) bool {
		return stats.ClosestEdges[i].Distance < stats.ClosestEdges[j].Distance
	})
	stats.processStats()
	return stats
}

// Clone creates a deep copy of a DistanceGaps, if it needs to be used in a circuit that tracks multiple copies of the circuit (in various stages of completion).
func (stats *DistanceGaps) Clone() *DistanceGaps {
	clone := &DistanceGaps{
		ClosestEdges:         make([]*tspmodel.DistanceToEdge, len(stats.ClosestEdges)),
		Gaps:                 make([]float64, len(stats.Gaps)),
		GapAverage:           stats.GapAverage,
		GapStandardDeviation: stats.GapStandardDeviation,
	}
	copy(clone.ClosestEdges, stats.ClosestEdges)
	copy(clone.Gaps, stats.Gaps)
	return clone
}

// processStats calculates the gaps between successive distances in ClosestEdge, the average gap, and the standard deviation of the gaps.
func (stats *DistanceGaps) processStats() {
	numGaps := len(stats.ClosestEdges) - 1
	numGapsFloat := float64(numGaps)

	stats.GapAverage = 0
	stats.Gaps = make([]float64, numGaps)

	// Compute averages
	for current, next := 0, 1; current < numGaps; current, next = current+1, next+1 {
		distanceGap := stats.ClosestEdges[next].Distance - stats.ClosestEdges[current].Distance
		stats.Gaps[current] = distanceGap
		stats.GapAverage += distanceGap / numGapsFloat
	}

	stats.GapStandardDeviation = 0
	// Compute standard deviations
	for current, next := 0, 1; current < numGaps; current, next = current+1, next+1 {
		currentGapDeviation := stats.Gaps[current] - stats.GapAverage
		stats.GapStandardDeviation += currentGapDeviation * currentGapDeviation / numGapsFloat
	}

	stats.GapStandardDeviation = math.Sqrt(stats.GapStandardDeviation)
}

// UpdateStats replaces the removed edge with the two edges that result from its split, then updates the statistics for this vertex.
func (stats *DistanceGaps) UpdateStats(removedEdge tspmodel.CircuitEdge, edgeA tspmodel.CircuitEdge, edgeB tspmodel.CircuitEdge) {
	prevNumEdges := len(stats.ClosestEdges)
	numEdges := prevNumEdges + 1

	vertex := stats.ClosestEdges[0].Vertex

	closer := &tspmodel.DistanceToEdge{
		Vertex:   vertex,
		Edge:     edgeA,
		Distance: edgeA.DistanceIncrease(vertex),
	}

	farther := &tspmodel.DistanceToEdge{
		Vertex:   vertex,
		Edge:     edgeB,
		Distance: edgeB.DistanceIncrease(vertex),
	}

	if farther.Distance < closer.Distance {
		closer, farther = farther, closer
	}

	// Update the closest edges list - note: the list is already sorted
	updatedEdges := make([]*tspmodel.DistanceToEdge, numEdges)
	for src, dest, isCloserInList, isFartherInList := 0, 0, false, false; dest < numEdges; dest++ {
		if src >= prevNumEdges {
			if !isCloserInList {
				updatedEdges[dest] = closer
				isCloserInList = true
			} else {
				updatedEdges[dest] = farther
				isFartherInList = true
			}
		} else {
			srcEdge := stats.ClosestEdges[src]
			if !isCloserInList && closer.Distance < srcEdge.Distance {
				updatedEdges[dest] = closer
				isCloserInList = true
			} else if !isFartherInList && farther.Distance < srcEdge.Distance {
				updatedEdges[dest] = farther
				isFartherInList = true
			} else if srcEdge.Edge == removedEdge {
				src++
				dest-- // Need to keep the destination at the same position for the next iteration, since nothing was copied this iteration.
			} else {
				updatedEdges[dest] = srcEdge
				src++
			}
		}
	}
	stats.ClosestEdges = updatedEdges
	stats.processStats()
}

// String converts the DistanceGaps to a string, with the edges printed as their index in the circuit, and vertices as their index in the initial request.
func (stats *DistanceGaps) String(vertexIndexLookup map[tspmodel.CircuitVertex]int, edgeIndexLookup map[tspmodel.CircuitEdge]int) string {
	if len(stats.ClosestEdges) <= 0 {
		return `{}`
	}

	s := fmt.Sprintf("{\r\n\t\"vertex\":%d,\r\n\t\"gapAverage\":%g,\r\n\t\"gapStdDev\":%g,\r\n\t\"closestEdges\":[",
		vertexIndexLookup[stats.ClosestEdges[0].Vertex], stats.GapAverage, stats.GapStandardDeviation)

	lastIndex := len(stats.ClosestEdges) - 1
	for i, e := range stats.ClosestEdges {
		if i == lastIndex {
			s += fmt.Sprintf("{\"edge\":%d,\"distance\":%g}", edgeIndexLookup[e.Edge], e.Distance)
		} else {
			s += fmt.Sprintf("{\"edge\":%d,\"distance\":%g},", edgeIndexLookup[e.Edge], e.Distance)
		}
	}

	s += "],\r\n\t\"gaps\":["
	lastIndex = len(stats.Gaps) - 1
	for i, gap := range stats.Gaps {
		s += fmt.Sprintf("{\"gap\":%g,\"gapZScore\":%g}", gap, (gap-stats.GapAverage)/stats.GapStandardDeviation)
		if i != lastIndex {
			s += ","
		}
	}

	s += "]}"

	return s
}
