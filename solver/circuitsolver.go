package solver

import "github.com/fealos/lee-tsp-go/model"

// FindShortestPathCircuit use a model.Circuit to approximate or solve the shortest path through a series of points.
// The complexity of this depends on the complexity of the algorithm used by the model.Circuit.
func FindShortestPathCircuit(circuit model.Circuit) {
	for nextVertex, nextEdge := circuit.FindNextVertexAndEdge(); nextVertex != nil; nextVertex, nextEdge = circuit.FindNextVertexAndEdge() {
		circuit.Update(nextVertex, nextEdge)
	}
}
