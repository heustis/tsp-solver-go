package circuit

import "github.com/heustis/tsp-solver-go/model"

// CompletedCircuit provides a no-op represenatation of a circuit, for use once an algorithm completes its computation.
// This allows for circuits with large memory requirements or circular references to be deleted without deleting the best computed circuit.
type CompletedCircuit struct {
	Circuit []model.CircuitVertex
	Length  float64
}

// NewCompletedCircuit returns a CompletedCircuit containing the result of the supplied Circuit.
// This will only account for vertices that are already attached to the circuit; any unattached vertices will be ignored.
func NewCompletedCircuit(c model.Circuit) *CompletedCircuit {
	return &CompletedCircuit{
		Circuit: c.GetAttachedVertices(),
		Length:  c.GetLength(),
	}
}

// CloneAndUpdate does nothing as the circuit is complete.
func (c *CompletedCircuit) CloneAndUpdate() ClonableCircuit {
	return nil
}

// Delete is implemented for compatibility with ClonableCircuit.
func (c *CompletedCircuit) Delete() {
	c.Circuit = nil
}

// FindNextVertexAndEdge determines the next vertex to add to the circuit, along with which edge it should be added to.
// This returns (nil,nil) because the circuit is complete.
func (c *CompletedCircuit) FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge) {
	return nil, nil
}

// GetAttachedVertices returns all vertices that have been added to the circuit.
// This returns them in the order they should be traversed as part of the circuit.
func (c *CompletedCircuit) GetAttachedVertices() []model.CircuitVertex {
	return c.Circuit
}

// GetLength returns the length of the circuit.
func (c *CompletedCircuit) GetLength() float64 {
	return c.Length
}

// GetLengthWithNext returns the length of the circuit, since it is complete.
func (c *CompletedCircuit) GetLengthWithNext() float64 {
	return c.Length
}

// GetUnattachedVertices returns an empty map, since the circuit is complete.
func (c *CompletedCircuit) GetUnattachedVertices() map[model.CircuitVertex]bool {
	return make(map[model.CircuitVertex]bool)
}

// Update does nothing as the circuit is complete.
func (c *CompletedCircuit) Update(vertexToAdd model.CircuitVertex, edgeToSplit model.CircuitEdge) {
}

var _ ClonableCircuit = (*CompletedCircuit)(nil)
var _ model.Circuit = (*CompletedCircuit)(nil)
