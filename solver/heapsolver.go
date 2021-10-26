package solver

import "github.com/fealos/lee-tsp-go/model"

func FindShortestPathHeap(circuit model.HeapableCircuit) (model.HeapableCircuit, int) {
	circuit.Prepare()
	circuit.BuildPerimiter()

	circuitHeap := model.NewHeap(func(a interface{}) float64 {
		return a.(model.HeapableCircuit).GetLengthWithNext()
	})
	circuitHeap.PushHeap(circuit)

	iterationCount := 0
	next := circuitHeap.PopHeap().(model.HeapableCircuit)
	for ; len(next.GetUnattachedVertices()) > 0 || next.GetLengthWithNext() < next.GetLength(); next = circuitHeap.PopHeap().(model.HeapableCircuit) {
		clone := next.CloneAndUpdate()
		circuitHeap.PushHeap(next)
		if clone != nil {
			circuitHeap.PushHeap(clone)
		}
		iterationCount++
	}

	// clean up the heap and each circuitHeap
	circuitHeap.Delete()

	return next, iterationCount
}
