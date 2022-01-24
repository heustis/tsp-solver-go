package circuit

import "github.com/fealos/lee-tsp-go/tspmodel"

type CompletedCircuit struct {
	Circuit []tspmodel.CircuitVertex
	Length  float64
}

func (c *CompletedCircuit) BuildPerimiter() {
}

func (c *CompletedCircuit) CloneAndUpdate() ClonableCircuit {
	return c
}

func (c *CompletedCircuit) Delete() {
	c.Circuit = nil
}

func (c *CompletedCircuit) FindNextVertexAndEdge() (tspmodel.CircuitVertex, tspmodel.CircuitEdge) {
	return nil, nil
}

func (c *CompletedCircuit) GetAttachedVertices() []tspmodel.CircuitVertex {
	return c.Circuit
}

func (c *CompletedCircuit) GetLength() float64 {
	return c.Length
}

func (c *CompletedCircuit) GetLengthWithNext() float64 {
	return c.Length
}

func (c *CompletedCircuit) GetUnattachedVertices() map[tspmodel.CircuitVertex]bool {
	return make(map[tspmodel.CircuitVertex]bool)
}

func (c *CompletedCircuit) Prepare() {
}

func (c *CompletedCircuit) Update(vertexToAdd tspmodel.CircuitVertex, edgeToSplit tspmodel.CircuitEdge) {
}

var _ ClonableCircuit = (*CompletedCircuit)(nil)
var _ tspmodel.Circuit = (*CompletedCircuit)(nil)
