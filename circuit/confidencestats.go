package circuit

import (
	"fmt"
	"math"

	"github.com/fealos/lee-tsp-go/model"
)

type vertexStatistics struct {
	ClosestEdges              []*model.DistanceToEdge
	DistanceGaps              []float64
	DistanceAverage           float64
	DistanceStandardDeviation float64
	GapAverage                float64
	GapStandardDeviation      float64
}

func (stats *vertexStatistics) processStats() {
	numEdges := len(stats.ClosestEdges)
	numGaps := float64(numEdges - 1)

	stats.DistanceAverage = 0
	stats.GapAverage = 0
	stats.DistanceGaps = make([]float64, numEdges-1)

	// Compute averages
	for current, next := 0, 1; current < numEdges; current, next = current+1, next+1 {
		currentDistance := stats.ClosestEdges[current].Distance
		stats.DistanceAverage += currentDistance / float64(numEdges)
		if next < numEdges {
			distanceGap := stats.ClosestEdges[next].Distance - currentDistance
			stats.DistanceGaps[current] = distanceGap
			stats.GapAverage += distanceGap / numGaps
		}
	}

	stats.DistanceStandardDeviation = 0
	stats.GapStandardDeviation = 0
	// Compute standard deviations
	for current, next := 0, 1; current < numEdges; current, next = current+1, next+1 {
		currentDeviation := stats.ClosestEdges[current].Distance - stats.DistanceAverage
		stats.DistanceStandardDeviation += currentDeviation * currentDeviation / float64(numEdges)
		if next < numEdges {
			currentGapDeviation := stats.DistanceGaps[current] - stats.GapAverage
			stats.GapStandardDeviation += currentGapDeviation * currentGapDeviation / numGaps
		}
	}

	stats.DistanceStandardDeviation = math.Sqrt(stats.DistanceStandardDeviation)
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

	s := fmt.Sprintf("{\r\n\t\"vertex\":%d,\r\n\t\"distanceAverage\":%g,\r\n\t\"distanceStdDev\":%g,\r\n\t\"gapAverage\":%g,\r\n\t\"gapStdDev\":%g,\r\n\t\"closestEdges\":[",
		vertexIndexLookup[stats.ClosestEdges[0].Vertex], stats.DistanceAverage, stats.DistanceStandardDeviation, stats.GapAverage, stats.GapStandardDeviation)

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
		if i == lastIndex {
			s += fmt.Sprintf("{\"gap\":%g,\"byGapStdDev\":%g,\"byDistStdDev\":%g}", gap, gap/stats.GapStandardDeviation, gap/stats.DistanceStandardDeviation)
		} else {
			s += fmt.Sprintf("{\"gap\":%g,\"byGapStdDev\":%g,\"byDistStdDev\":%g},", gap, gap/stats.GapStandardDeviation, gap/stats.DistanceStandardDeviation)
		}
	}

	s += "]}"

	return s
}
