package circuit_test

import (
	"testing"

	"github.com/heustis/lee-tsp-go/circuit"
	"github.com/heustis/lee-tsp-go/model"
	"github.com/heustis/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestNewGeneticAlgorithm(t *testing.T) {
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

	c := circuit.NewGeneticAlgorithm(initVertices, 10, 50, 10)
	assert.NotNil(c)

	initCircuit := c.GetAttachedVertices()
	assert.Len(initCircuit, len(initVertices))
	assert.Len(c.GetUnattachedVertices(), 0)
	nextVertex, nextEdge := c.FindNextVertexAndEdge()
	assert.NotNil(nextVertex)
	assert.Equal(initCircuit[0], nextVertex)
	assert.Nil(nextEdge)
}

func TestNewGeneticAlgorithmWithPerimeterBuilder(t *testing.T) {
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

	c := circuit.NewGeneticAlgorithmWithPerimeterBuilder(initVertices, model2d.BuildPerimiter, 10, 50, 10)
	assert.NotNil(c)

	initCircuit := c.GetAttachedVertices()
	assert.Len(initCircuit, len(initVertices))
	assert.Len(c.GetUnattachedVertices(), 0)
	nextVertex, nextEdge := c.FindNextVertexAndEdge()
	assert.NotNil(nextVertex)
	assert.Equal(initCircuit[0], nextVertex)
	assert.Nil(nextEdge)
}

func TestUpdate_GeneticAlgorithm(t *testing.T) {
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

	c := circuit.NewGeneticAlgorithm(initVertices, 10, 50, 10)
	assert.NotNil(c)
	previousLen := c.GetLength()
	previousCircuit := c.GetAttachedVertices()

	for i := 0; i < 10; i++ {
		c.Update(c.FindNextVertexAndEdge())
		currentLength := c.GetLength()
		assert.LessOrEqual(currentLength, previousLen)
		previousLen = currentLength
		currentCircuit := c.GetAttachedVertices()
		for _, v := range initVertices {
			assert.Contains(currentCircuit, v)
		}
		previousCircuit = currentCircuit
	}

	c.Update(c.FindNextVertexAndEdge())
	assert.Equal(previousLen, c.GetLength())
	assert.Equal(previousCircuit, c.GetAttachedVertices())
}

func TestUpdate_GeneticAlgorithmWithPerimeterBuilder(t *testing.T) {
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

	c := circuit.NewGeneticAlgorithmWithPerimeterBuilder(initVertices, model2d.BuildPerimiter, 10, 50, 10)
	assert.NotNil(c)
	previousLen := c.GetLength()
	previousCircuit := c.GetAttachedVertices()

	for i := 0; i < 10; i++ {
		c.Update(c.FindNextVertexAndEdge())
		currentLength := c.GetLength()
		assert.LessOrEqual(currentLength, previousLen)
		previousLen = currentLength
		currentCircuit := c.GetAttachedVertices()
		for _, v := range initVertices {
			assert.Contains(currentCircuit, v)
		}
		previousCircuit = currentCircuit
	}

	c.Update(c.FindNextVertexAndEdge())
	assert.Equal(previousLen, c.GetLength())
	assert.Equal(previousCircuit, c.GetAttachedVertices())
}
