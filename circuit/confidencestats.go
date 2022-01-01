package circuit

import (
	"fmt"
	"math"

	"github.com/fealos/lee-tsp-go/model"
)

type vertexStatistics struct {
	ClosestEdges         []*model.DistanceToEdge
	DistanceGaps         []float64
	GapAverage           float64
	GapStandardDeviation float64
}

func (stats *vertexStatistics) processStats() {
	numGaps := len(stats.ClosestEdges) - 1
	numGapsFloat := float64(numGaps)

	stats.GapAverage = 0
	stats.DistanceGaps = make([]float64, numGaps)

	// Compute averages
	for current, next := 0, 1; current < numGaps; current, next = current+1, next+1 {
		distanceGap := stats.ClosestEdges[next].Distance - stats.ClosestEdges[current].Distance
		stats.DistanceGaps[current] = distanceGap
		stats.GapAverage += distanceGap / numGapsFloat
	}

	stats.GapStandardDeviation = 0
	// Compute standard deviations
	for current, next := 0, 1; current < numGaps; current, next = current+1, next+1 {
		currentGapDeviation := stats.DistanceGaps[current] - stats.GapAverage
		stats.GapStandardDeviation += currentGapDeviation * currentGapDeviation / numGapsFloat
	}

	stats.GapStandardDeviation = math.Sqrt(stats.GapStandardDeviation)
}

func (stats *vertexStatistics) updateClosestEdges(removedEdge model.CircuitEdge, edgeA model.CircuitEdge, edgeB model.CircuitEdge) {
	prevNumEdges := len(stats.ClosestEdges)
	numEdges := prevNumEdges + 1

	vertex := stats.ClosestEdges[0].Vertex

	closer := &model.DistanceToEdge{
		Vertex:   vertex,
		Edge:     edgeA,
		Distance: edgeA.DistanceIncrease(vertex),
	}

	farther := &model.DistanceToEdge{
		Vertex:   vertex,
		Edge:     edgeB,
		Distance: edgeB.DistanceIncrease(vertex),
	}

	if farther.Distance < closer.Distance {
		closer, farther = farther, closer
	}

	// Update the closest edges list - note: the list is already sorted
	updatedEdges := make([]*model.DistanceToEdge, numEdges)
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
			} else {
				updatedEdges[dest] = srcEdge
				src++
			}
		}
	}
}

func (stats *vertexStatistics) UpdateStats(removedEdge model.CircuitEdge, edgeA model.CircuitEdge, edgeB model.CircuitEdge) {
	stats.updateClosestEdges(removedEdge, edgeA, edgeB)
	stats.processStats()
}

func (stats *vertexStatistics) ToString(vertexIndexLookup map[model.CircuitVertex]int, edgeIndexLookup map[model.CircuitEdge]int) string {
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
	lastIndex = len(stats.DistanceGaps) - 1
	for i, gap := range stats.DistanceGaps {
		s += fmt.Sprintf("{\"gap\":%g,\"byGapStdDev\":%g}", gap, gap/stats.GapStandardDeviation)
		if i != lastIndex {
			s += ","
		}
	}

	s += "]}"

	return s
}
