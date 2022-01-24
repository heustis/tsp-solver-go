package tspsolver

import "github.com/fealos/lee-tsp-go/tspmodel"

// FindShortestPathCircuit use a tspmodel.Circuit to approximate or solve the shortest path through a series of points.
// The complexity of this depends on the complexity of the algorithm used by the tspmodel.Circuit.
func FindShortestPathCircuit(circuit tspmodel.Circuit) {
	circuit.Prepare()
	circuit.BuildPerimiter()

	for nextVertex, nextEdge := circuit.FindNextVertexAndEdge(); nextVertex != nil; nextVertex, nextEdge = circuit.FindNextVertexAndEdge() {
		circuit.Update(nextVertex, nextEdge)
	}
}
