package model

type HeapableCircuitComplete struct {
	Circuit []CircuitVertex
	Length  float64
}

func (c *HeapableCircuitComplete) BuildPerimiter() {
}

func (c *HeapableCircuitComplete) CloneAndUpdate() HeapableCircuit {
	return c
}

func (c *HeapableCircuitComplete) Delete() {
	c.Circuit = nil
}

func (c *HeapableCircuitComplete) FindNextVertexAndEdge() (CircuitVertex, CircuitEdge) {
	return nil, nil
}

func (c *HeapableCircuitComplete) GetAttachedVertices() []CircuitVertex {
	return c.Circuit
}

func (c *HeapableCircuitComplete) GetLength() float64 {
	return c.Length
}

func (c *HeapableCircuitComplete) GetLengthWithNext() float64 {
	return c.Length
}

func (c *HeapableCircuitComplete) GetUnattachedVertices() map[CircuitVertex]bool {
	return make(map[CircuitVertex]bool)
}

func (c *HeapableCircuitComplete) Prepare() {
}

func (c *HeapableCircuitComplete) Update(vertexToAdd CircuitVertex, edgeToSplit CircuitEdge) {
}

var _ HeapableCircuit = (*HeapableCircuitComplete)(nil)
var _ Circuit = (*HeapableCircuitComplete)(nil)
