package circuit_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestCompletedCircuit(t *testing.T) {
	assert := assert.New(t)

	expectedPath := []model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(-7, 6),
	}

	c := &circuit.CompletedCircuit{
		Circuit: []model.CircuitVertex{
			model2d.NewVertex2D(-15, -15),
			model2d.NewVertex2D(0, 0),
			model2d.NewVertex2D(15, -15),
			model2d.NewVertex2D(3, 0),
			model2d.NewVertex2D(3, 13),
			model2d.NewVertex2D(8, 5),
			model2d.NewVertex2D(9, 6),
			model2d.NewVertex2D(-7, 6),
		},

		Length: 12345.6789,
	}

	c.Prepare()
	assert.Equal(expectedPath, c.GetAttachedVertices())
	assert.Equal(12345.6789, c.GetLength())
	assert.Equal(12345.6789, c.GetLengthWithNext())
	assert.Len(c.GetUnattachedVertices(), 0)

	c.BuildPerimiter()
	assert.Equal(expectedPath, c.GetAttachedVertices())
	assert.Equal(12345.6789, c.GetLength())
	assert.Equal(12345.6789, c.GetLengthWithNext())
	assert.Len(c.GetUnattachedVertices(), 0)

	assert.Equal(c, c.CloneAndUpdate())
	assert.Len(c.GetUnattachedVertices(), 0)

	v, e := c.FindNextVertexAndEdge()
	assert.Nil(v)
	assert.Nil(e)

	c.Update(v, e)
	assert.Equal(expectedPath, c.GetAttachedVertices())
	assert.Equal(12345.6789, c.GetLength())
	assert.Equal(12345.6789, c.GetLengthWithNext())
	assert.Len(c.GetUnattachedVertices(), 0)

	c.Update(expectedPath[2], expectedPath[5].EdgeTo(expectedPath[6]))
	assert.Equal(expectedPath, c.GetAttachedVertices())
	assert.Equal(12345.6789, c.GetLength())
	assert.Equal(12345.6789, c.GetLengthWithNext())
	assert.Len(c.GetUnattachedVertices(), 0)

	c.Delete()
	assert.Nil(c.GetAttachedVertices())
}
