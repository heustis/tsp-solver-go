package circuit

import "github.com/fealos/lee-tsp-go/model"

// CompletedCircuit provides a no-op represenatation of a circuit, for use once an algorithm completes its computation.
// This allows for circuits with large memory requirements or circular references to be deleted without deleting the best computed circuit.
type CompletedCircuit struct {
	Circuit []model.CircuitVertex
	Length  float64
}

func (c *CompletedCircuit) CloneAndUpdate() ClonableCircuit {
	return c
}

func (c *CompletedCircuit) Delete() {
	c.Circuit = nil
}

func (c *CompletedCircuit) FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge) {
	return nil, nil
}

func (c *CompletedCircuit) GetAttachedVertices() []model.CircuitVertex {
	return c.Circuit
}

func (c *CompletedCircuit) GetLength() float64 {
	return c.Length
}

func (c *CompletedCircuit) GetLengthWithNext() float64 {
	return c.Length
}

func (c *CompletedCircuit) GetUnattachedVertices() map[model.CircuitVertex]bool {
	return make(map[model.CircuitVertex]bool)
}

func (c *CompletedCircuit) Update(vertexToAdd model.CircuitVertex, edgeToSplit model.CircuitEdge) {
}

var _ ClonableCircuit = (*CompletedCircuit)(nil)
var _ model.Circuit = (*CompletedCircuit)(nil)
