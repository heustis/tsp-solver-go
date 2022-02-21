package circuit_test

import (
	"testing"

	"github.com/heustis/tsp-solver-go/circuit"
	"github.com/heustis/tsp-solver-go/model"
	"github.com/heustis/tsp-solver-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestNewSimulatedAnnealing(t *testing.T) {
	assert := assert.New(t)

	initVertices := []model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(-7, 6),
	}

	c := circuit.NewSimulatedAnnealing(initVertices, 100, false)

	assert.NotNil(c)
	assert.Equal(initVertices, c.GetAttachedVertices())
	assert.InDelta(123.95617933216532, c.GetLength(), model.Threshold)
	assert.Len(c.GetUnattachedVertices(), 0)
	nextVertex, nextEdge := c.FindNextVertexAndEdge()
	assert.NotNil(nextVertex)
	assert.Nil(nextEdge)
}

func TestNewSimulatedAnnealingFromCircuit(t *testing.T) {
	assert := assert.New(t)

	initVertices := model2d.DeduplicateVertices([]model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(-7, 6),
	})

	circuitVertices := make([]model.CircuitVertex, len(initVertices))
	copy(circuitVertices, initVertices)
	c := circuit.NewSimulatedAnnealingFromCircuit(circuit.NewConvexConcave(initVertices, model2d.BuildPerimiter, false), 100, false)

	assert.NotNil(c)

	expectedVertices := []model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(-7, 6),
	}

	assert.Equal(expectedVertices, c.GetAttachedVertices())
	assert.InDelta(106.59678993710583, c.GetLength(), model.Threshold)
	assert.Len(c.GetUnattachedVertices(), 0)
	nextVertex, nextEdge := c.FindNextVertexAndEdge()
	assert.NotNil(nextVertex)
	assert.Nil(nextEdge)
}

func TestUpdate_SimulatedAnnealing(t *testing.T) {
	assert := assert.New(t)

	initVertices := []model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(-7, 6),
	}

	c := circuit.NewSimulatedAnnealing(initVertices, 100, false)
	c.SetSeed(1)

	assert.NotNil(c)
	assert.Equal(initVertices, c.GetAttachedVertices())
	assert.InDelta(123.95617933216532, c.GetLength(), model.Threshold)
	assert.InDelta(model.Length(c.GetAttachedVertices()), c.GetLength(), model.Threshold)

	c.Update(c.FindNextVertexAndEdge())
	assert.InDelta(127.97344237445196, c.GetLength(), model.Threshold)
	assert.InDelta(model.Length(c.GetAttachedVertices()), c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(0, 0),
	}, c.GetAttachedVertices())

	for i := 0; i < 98; i++ {
		c.Update(c.FindNextVertexAndEdge())
	}

	assert.InDelta(109.55350205245301, c.GetLength(), model.Threshold)
	assert.InDelta(model.Length(c.GetAttachedVertices()), c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
	}, c.GetAttachedVertices())

	nextVertex, nextEdge := c.FindNextVertexAndEdge()
	assert.NotNil(nextVertex)
	assert.Nil(nextEdge)
	c.Update(nextVertex, nextEdge)
	assert.InDelta(109.55350205245301, c.GetLength(), model.Threshold)
	assert.InDelta(model.Length(c.GetAttachedVertices()), c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
	}, c.GetAttachedVertices())

	nextVertex, nextEdge = c.FindNextVertexAndEdge()
	assert.Nil(nextVertex)
	assert.Nil(nextEdge)
	c.Update(model2d.NewVertex2D(3, 0), nil)
	assert.InDelta(109.55350205245301, c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
	}, c.GetAttachedVertices())

	c = circuit.NewSimulatedAnnealing(initVertices, 1000, false)
	c.SetSeed(1)
	for i := 0; i < 1000; i++ {
		c.Update(c.FindNextVertexAndEdge())
	}
	assert.InDelta(109.55350205245301, c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
	}, c.GetAttachedVertices())
}

func TestUpdate_SimulatedAnnealing_PreferCloseNeighbors(t *testing.T) {
	assert := assert.New(t)

	initVertices := []model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(-7, 6),
	}

	c := circuit.NewSimulatedAnnealing(initVertices, 100, true)
	c.SetSeed(1)

	assert.NotNil(c)
	assert.Equal(initVertices, c.GetAttachedVertices())
	assert.InDelta(123.95617933216532, c.GetLength(), model.Threshold)

	c.Update(c.FindNextVertexAndEdge())
	assert.InDelta(127.97344237445196, c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(0, 0),
	}, c.GetAttachedVertices())

	for i := 0; i < 98; i++ {
		c.Update(c.FindNextVertexAndEdge())
	}

	assert.InDelta(122.25667258059363, c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(15, -15),
	}, c.GetAttachedVertices())

	nextVertex, nextEdge := c.FindNextVertexAndEdge()
	assert.NotNil(nextVertex)
	assert.Nil(nextEdge)
	c.Update(nextVertex, nextEdge)
	assert.InDelta(122.25667258059363, c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(15, -15),
	}, c.GetAttachedVertices())

	nextVertex, nextEdge = c.FindNextVertexAndEdge()
	assert.Nil(nextVertex)
	assert.Nil(nextEdge)
	c.Update(model2d.NewVertex2D(3, 0), nil)
	assert.InDelta(122.25667258059363, c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(15, -15),
	}, c.GetAttachedVertices())
}

func TestUpdate_SimulatedAnnealingFromCircuit(t *testing.T) {
	assert := assert.New(t)

	initVertices := model2d.DeduplicateVertices([]model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(-7, 6),
	})

	c := circuit.NewSimulatedAnnealingFromCircuit(circuit.NewConvexConcave(initVertices, model2d.BuildPerimiter, false), 100, false)
	c.SetSeed(1)

	assert.NotNil(c)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(-7, 6),
	}, c.GetAttachedVertices())
	assert.InDelta(106.59678993710583, c.GetLength(), model.Threshold)

	c.Update(c.FindNextVertexAndEdge())
	assert.InDelta(112.8596677665242, c.GetLength(), model.Threshold)
	assert.InDelta(model.Length(c.GetAttachedVertices()), c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(15, -15),
	}, c.GetAttachedVertices())

	for i := 0; i < 98; i++ {
		c.Update(c.FindNextVertexAndEdge())
	}

	assert.InDelta(106.59678993710581, c.GetLength(), model.Threshold)
	assert.InDelta(model.Length(c.GetAttachedVertices()), c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(0, 0),
	}, c.GetAttachedVertices())

	nextVertex, nextEdge := c.FindNextVertexAndEdge()
	assert.NotNil(nextVertex)
	assert.Nil(nextEdge)
	c.Update(nextVertex, nextEdge)
	assert.InDelta(106.59678993710581, c.GetLength(), model.Threshold)
	assert.InDelta(model.Length(c.GetAttachedVertices()), c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(0, 0),
	}, c.GetAttachedVertices())

	nextVertex, nextEdge = c.FindNextVertexAndEdge()
	assert.Nil(nextVertex)
	assert.Nil(nextEdge)
	c.Update(model2d.NewVertex2D(3, 0), nil)
	assert.InDelta(106.59678993710581, c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(0, 0),
	}, c.GetAttachedVertices())
}

func TestUpdate_SimulatedAnnealingFromCircuit_GeometricTemperature(t *testing.T) {
	assert := assert.New(t)

	initVertices := model2d.DeduplicateVertices([]model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(-7, 6),
	})

	c := circuit.NewSimulatedAnnealingFromCircuit(circuit.NewConvexConcave(initVertices, model2d.BuildPerimiter, false), 100, false)
	c.SetSeed(1)
	c.SetTemperatureFunction(circuit.CalculateTemperatureGeometric)

	assert.NotNil(c)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(-7, 6),
	}, c.GetAttachedVertices())
	assert.InDelta(106.59678993710583, c.GetLength(), model.Threshold)

	c.Update(c.FindNextVertexAndEdge())
	assert.InDelta(112.8596677665242, c.GetLength(), model.Threshold)
	assert.InDelta(model.Length(c.GetAttachedVertices()), c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(15, -15),
	}, c.GetAttachedVertices())

	for i := 0; i < 98; i++ {
		c.Update(c.FindNextVertexAndEdge())
	}

	assert.InDelta(106.59678993710581, c.GetLength(), model.Threshold)
	assert.InDelta(model.Length(c.GetAttachedVertices()), c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(3, 13),
	}, c.GetAttachedVertices())

	nextVertex, nextEdge := c.FindNextVertexAndEdge()
	assert.NotNil(nextVertex)
	assert.Nil(nextEdge)
	c.Update(nextVertex, nextEdge)
	assert.InDelta(106.59678993710581, c.GetLength(), model.Threshold)
	assert.InDelta(model.Length(c.GetAttachedVertices()), c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(3, 13),
	}, c.GetAttachedVertices())

	nextVertex, nextEdge = c.FindNextVertexAndEdge()
	assert.Nil(nextVertex)
	assert.Nil(nextEdge)
	c.Update(model2d.NewVertex2D(3, 0), nil)
	assert.InDelta(106.59678993710581, c.GetLength(), model.Threshold)
	assert.Equal([]model.CircuitVertex{
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(3, 13),
	}, c.GetAttachedVertices())
}

func TestCalculateTemperatureGeometric(t *testing.T) {
	assert := assert.New(t)

	assert.InDelta(0.52446047, circuit.CalculateTemperatureGeometric(10, 80), model.Threshold)
	assert.InDelta(0.59873693, circuit.CalculateTemperatureGeometric(10, 100), model.Threshold)

	assert.InDelta(0.07565733, circuit.CalculateTemperatureGeometric(40, 80), model.Threshold)
	assert.InDelta(0.07694497, circuit.CalculateTemperatureGeometric(50, 100), model.Threshold)

	assert.InDelta(0.02081021, circuit.CalculateTemperatureGeometric(60, 80), model.Threshold)
	assert.InDelta(0.02758369, circuit.CalculateTemperatureGeometric(70, 100), model.Threshold)

	assert.InDelta(0.00623213, circuit.CalculateTemperatureGeometric(99, 100), model.Threshold)
	assert.InDelta(0.00668740, circuit.CalculateTemperatureGeometric(999, 1000), model.Threshold)
}

func TestCalculateTemperatureLinear(t *testing.T) {
	assert := assert.New(t)

	assert.InDelta(.875, circuit.CalculateTemperatureLinear(10, 80), model.Threshold)
	assert.InDelta(.9, circuit.CalculateTemperatureLinear(10, 100), model.Threshold)

	assert.InDelta(.5, circuit.CalculateTemperatureLinear(40, 80), model.Threshold)
	assert.InDelta(.5, circuit.CalculateTemperatureLinear(50, 100), model.Threshold)

	assert.InDelta(.25, circuit.CalculateTemperatureLinear(60, 80), model.Threshold)
	assert.InDelta(.3, circuit.CalculateTemperatureLinear(70, 100), model.Threshold)

	assert.InDelta(.01, circuit.CalculateTemperatureLinear(99, 100), model.Threshold)
	assert.InDelta(.001, circuit.CalculateTemperatureLinear(999, 1000), model.Threshold)
}
