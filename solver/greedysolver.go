package solver

import "github.com/fealos/lee-tsp-go/model"

func FindShortestPathGreedy(circuit model.Circuit) {
	circuit.Prepare()
	circuit.BuildPerimiter()

	for nextVertex, nextEdge := circuit.FindNextVertexAndEdge(); nextVertex != nil; nextVertex, nextEdge = circuit.FindNextVertexAndEdge() {
		circuit.Update(nextVertex, nextEdge)
	}
}
