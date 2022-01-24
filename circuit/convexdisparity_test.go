package circuit_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/tspmodel"
	"github.com/fealos/lee-tsp-go/tspmodel2d"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerimeter_ConvexConcaveDisparity(t *testing.T) {
	assert := assert.New(t)
	c := circuit.NewConvexConcaveDisparity([]tspmodel.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		tspmodel2d.NewVertex2D(-15, -15), // Index 0 after sorting
		tspmodel2d.NewVertex2D(0, 0),     // Index 2 after sorting
		tspmodel2d.NewVertex2D(15, -15),  // Index 7 after sorting
		tspmodel2d.NewVertex2D(3, 0),     // Index 3 after sorting
		tspmodel2d.NewVertex2D(3, 13),    // Index 4 after sorting
		tspmodel2d.NewVertex2D(8, 5),     // Index 5 after sorting
		tspmodel2d.NewVertex2D(9, 6),     // Index 6 after sorting
		tspmodel2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	},
		tspmodel2d.DeduplicateVertices,
		tspmodel2d.BuildPerimiter,
		false,
	).(*circuit.ConvexConcaveDisparity)

	c.Prepare()
	c.BuildPerimiter()

	assert.Len(c.Vertices, 8)

	assert.Len(c.GetAttachedVertices(), 5)
	assert.Equal(tspmodel2d.NewVertex2D(-15, -15), c.GetAttachedVertices()[0])
	assert.Equal(tspmodel2d.NewVertex2D(15, -15), c.GetAttachedVertices()[1])
	assert.Equal(tspmodel2d.NewVertex2D(9, 6), c.GetAttachedVertices()[2])
	assert.Equal(tspmodel2d.NewVertex2D(3, 13), c.GetAttachedVertices()[3])
	assert.Equal(tspmodel2d.NewVertex2D(-7, 6), c.GetAttachedVertices()[4])

	assert.InDelta(95.73863479511238, c.GetLength(), tspmodel.Threshold)

	assert.Len(c.GetAttachedEdges(), 5)
	assert.True(c.Vertices[0].EdgeTo(c.Vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(c.Vertices[7].EdgeTo(c.Vertices[6]).Equals(c.GetAttachedEdges()[1]))
	assert.True(c.Vertices[6].EdgeTo(c.Vertices[4]).Equals(c.GetAttachedEdges()[2]))
	assert.True(c.Vertices[4].EdgeTo(c.Vertices[1]).Equals(c.GetAttachedEdges()[3]))
	assert.True(c.Vertices[1].EdgeTo(c.Vertices[0]).Equals(c.GetAttachedEdges()[4]))

	assert.Len(c.GetUnattachedVertices(), 3)
	assert.True(c.GetUnattachedVertices()[c.Vertices[2]])
	assert.True(c.GetUnattachedVertices()[c.Vertices[3]])
	assert.True(c.GetUnattachedVertices()[c.Vertices[5]])
}

func TestPrepare_ConvexConcaveDisparity(t *testing.T) {
	assert := assert.New(t)
	c := circuit.NewConvexConcaveDisparity([]tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(-15, -15),
		tspmodel2d.NewVertex2D(0, 0),
		tspmodel2d.NewVertex2D(15, -15),
		tspmodel2d.NewVertex2D(-15-tspmodel.Threshold/3.0, -15),
		tspmodel2d.NewVertex2D(0, 0),
		tspmodel2d.NewVertex2D(-15+tspmodel.Threshold/3.0, -15-tspmodel.Threshold/3.0),
		tspmodel2d.NewVertex2D(3, 0),
		tspmodel2d.NewVertex2D(3, 13),
		tspmodel2d.NewVertex2D(7, 6),
		tspmodel2d.NewVertex2D(-7, 6),
	},
		tspmodel2d.DeduplicateVertices,
		tspmodel2d.BuildPerimiter,
		false,
	).(*circuit.ConvexConcaveDisparity)

	c.Prepare()

	assert.NotNil(c.Vertices)
	assert.Len(c.Vertices, 7)
	assert.ElementsMatch(c.Vertices, []tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(-15+tspmodel.Threshold/3.0, -15-tspmodel.Threshold/3.0),
		tspmodel2d.NewVertex2D(-7, 6),
		tspmodel2d.NewVertex2D(0, 0),
		tspmodel2d.NewVertex2D(3, 0),
		tspmodel2d.NewVertex2D(3, 13),
		tspmodel2d.NewVertex2D(7, 6),
		tspmodel2d.NewVertex2D(15, -15),
	})

	assert.NotNil(c.GetUnattachedVertices())
	assert.Len(c.GetUnattachedVertices(), 0)

	assert.NotNil(c.GetAttachedVertices())
	assert.Len(c.GetAttachedVertices(), 0)

	assert.NotNil(c.GetAttachedEdges())
	assert.Len(c.GetAttachedEdges(), 0)
}

func TestUpdate_ConvexConcaveDisparity(t *testing.T) {
	assert := assert.New(t)
	c := circuit.NewConvexConcaveDisparity([]tspmodel.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		tspmodel2d.NewVertex2D(-15, -15), // Index 0 after sorting
		tspmodel2d.NewVertex2D(0, 0),     // Index 2 after sorting
		tspmodel2d.NewVertex2D(15, -15),  // Index 7 after sorting
		tspmodel2d.NewVertex2D(3, 0),     // Index 3 after sorting
		tspmodel2d.NewVertex2D(3, 13),    // Index 4 after sorting
		tspmodel2d.NewVertex2D(8, 5),     // Index 5 after sorting
		tspmodel2d.NewVertex2D(9, 6),     // Index 6 after sorting
		tspmodel2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	},
		tspmodel2d.DeduplicateVertices,
		tspmodel2d.BuildPerimiter,
		false,
	).(*circuit.ConvexConcaveDisparity)

	c.Prepare()
	c.BuildPerimiter()

	assert.Len(c.Vertices, 8)
	assert.Len(c.GetAttachedVertices(), 5)
	assert.Len(c.GetAttachedEdges(), 5)
	assert.Len(c.GetUnattachedVertices(), 3)

	// Init circuit := (-15, -15), (15, -15), (9,6), (3,13), (-7,6)

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.Vertices, 8)
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Len(c.GetAttachedEdges(), 6)
	assert.Len(c.GetUnattachedVertices(), 2)

	assert.True(c.Vertices[0].EdgeTo(c.Vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(c.Vertices[7].EdgeTo(c.Vertices[3]).Equals(c.GetAttachedEdges()[1]))
	assert.True(c.Vertices[3].EdgeTo(c.Vertices[6]).Equals(c.GetAttachedEdges()[2]))
	assert.True(c.Vertices[6].EdgeTo(c.Vertices[4]).Equals(c.GetAttachedEdges()[3]))
	assert.True(c.Vertices[4].EdgeTo(c.Vertices[1]).Equals(c.GetAttachedEdges()[4]))
	assert.True(c.Vertices[1].EdgeTo(c.Vertices[0]).Equals(c.GetAttachedEdges()[5]))

	assert.Equal(tspmodel2d.NewVertex2D(3, 0), c.GetAttachedVertices()[2])
	assert.Equal(tspmodel2d.NewVertex2D(9, 6), c.GetAttachedVertices()[3])

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.Vertices, 8)
	assert.Len(c.GetAttachedVertices(), 7)
	assert.Len(c.GetAttachedEdges(), 7)
	assert.Len(c.GetUnattachedVertices(), 1)

	assert.True(c.Vertices[0].EdgeTo(c.Vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(c.Vertices[7].EdgeTo(c.Vertices[3]).Equals(c.GetAttachedEdges()[1]))
	assert.True(c.Vertices[3].EdgeTo(c.Vertices[5]).Equals(c.GetAttachedEdges()[2]))
	assert.True(c.Vertices[5].EdgeTo(c.Vertices[6]).Equals(c.GetAttachedEdges()[3]))
	assert.True(c.Vertices[6].EdgeTo(c.Vertices[4]).Equals(c.GetAttachedEdges()[4]))
	assert.True(c.Vertices[4].EdgeTo(c.Vertices[1]).Equals(c.GetAttachedEdges()[5]))
	assert.True(c.Vertices[1].EdgeTo(c.Vertices[0]).Equals(c.GetAttachedEdges()[6]))

	assert.Equal(tspmodel2d.NewVertex2D(3, 0), c.GetAttachedVertices()[2])
	assert.Equal(tspmodel2d.NewVertex2D(8, 5), c.GetAttachedVertices()[3])
	assert.Equal(tspmodel2d.NewVertex2D(9, 6), c.GetAttachedVertices()[4])

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.Vertices, 8)
	assert.Len(c.GetAttachedVertices(), 8)
	assert.Len(c.GetAttachedEdges(), 8)
	assert.Len(c.GetUnattachedVertices(), 0)

	assert.True(c.Vertices[0].EdgeTo(c.Vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(c.Vertices[7].EdgeTo(c.Vertices[2]).Equals(c.GetAttachedEdges()[1]))
	assert.True(c.Vertices[2].EdgeTo(c.Vertices[3]).Equals(c.GetAttachedEdges()[2]))
	assert.True(c.Vertices[3].EdgeTo(c.Vertices[5]).Equals(c.GetAttachedEdges()[3]))
	assert.True(c.Vertices[5].EdgeTo(c.Vertices[6]).Equals(c.GetAttachedEdges()[4]))
	assert.True(c.Vertices[6].EdgeTo(c.Vertices[4]).Equals(c.GetAttachedEdges()[5]))
	assert.True(c.Vertices[4].EdgeTo(c.Vertices[1]).Equals(c.GetAttachedEdges()[6]))
	assert.True(c.Vertices[1].EdgeTo(c.Vertices[0]).Equals(c.GetAttachedEdges()[7]))

	assert.Equal(tspmodel2d.NewVertex2D(0, 0), c.GetAttachedVertices()[2])
	assert.Equal(tspmodel2d.NewVertex2D(3, 0), c.GetAttachedVertices()[3])
	assert.Equal(tspmodel2d.NewVertex2D(8, 5), c.GetAttachedVertices()[4])
	assert.Equal(tspmodel2d.NewVertex2D(9, 6), c.GetAttachedVertices()[5])

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.Vertices, 8)
	assert.Len(c.GetAttachedVertices(), 8)
	assert.Len(c.GetAttachedEdges(), 8)
	assert.Len(c.GetUnattachedVertices(), 0)
}

func TestUpdate_ConvexConcaveDisparityRelative(t *testing.T) {
	assert := assert.New(t)
	c := circuit.NewConvexConcaveDisparity([]tspmodel.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		tspmodel2d.NewVertex2D(-15, -15), // Index 0 after sorting
		tspmodel2d.NewVertex2D(0, 0),     // Index 2 after sorting
		tspmodel2d.NewVertex2D(15, -15),  // Index 7 after sorting
		tspmodel2d.NewVertex2D(3, 0),     // Index 3 after sorting
		tspmodel2d.NewVertex2D(3, 13),    // Index 4 after sorting
		tspmodel2d.NewVertex2D(8, 5),     // Index 5 after sorting
		tspmodel2d.NewVertex2D(9, 6),     // Index 6 after sorting
		tspmodel2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	},
		tspmodel2d.DeduplicateVertices,
		tspmodel2d.BuildPerimiter,
		true,
	).(*circuit.ConvexConcaveDisparity)

	c.Prepare()
	c.BuildPerimiter()

	assert.Len(c.Vertices, 8)
	assert.Len(c.GetAttachedVertices(), 5)
	assert.Len(c.GetAttachedEdges(), 5)
	assert.Len(c.GetUnattachedVertices(), 3)

	// Init circuit := (-15, -15), (15, -15), (9,6), (3,13), (-7,6)

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.Vertices, 8)
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Len(c.GetAttachedEdges(), 6)
	assert.Len(c.GetUnattachedVertices(), 2)

	assert.True(c.Vertices[0].EdgeTo(c.Vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(c.Vertices[7].EdgeTo(c.Vertices[3]).Equals(c.GetAttachedEdges()[1]))
	assert.True(c.Vertices[3].EdgeTo(c.Vertices[6]).Equals(c.GetAttachedEdges()[2]))
	assert.True(c.Vertices[6].EdgeTo(c.Vertices[4]).Equals(c.GetAttachedEdges()[3]))
	assert.True(c.Vertices[4].EdgeTo(c.Vertices[1]).Equals(c.GetAttachedEdges()[4]))
	assert.True(c.Vertices[1].EdgeTo(c.Vertices[0]).Equals(c.GetAttachedEdges()[5]))

	assert.Equal(tspmodel2d.NewVertex2D(3, 0), c.GetAttachedVertices()[2])
	assert.Equal(tspmodel2d.NewVertex2D(9, 6), c.GetAttachedVertices()[3])

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.Vertices, 8)
	assert.Len(c.GetAttachedVertices(), 7)
	assert.Len(c.GetAttachedEdges(), 7)
	assert.Len(c.GetUnattachedVertices(), 1)

	assert.True(c.Vertices[0].EdgeTo(c.Vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(c.Vertices[7].EdgeTo(c.Vertices[3]).Equals(c.GetAttachedEdges()[1]))
	assert.True(c.Vertices[3].EdgeTo(c.Vertices[5]).Equals(c.GetAttachedEdges()[2]))
	assert.True(c.Vertices[5].EdgeTo(c.Vertices[6]).Equals(c.GetAttachedEdges()[3]))
	assert.True(c.Vertices[6].EdgeTo(c.Vertices[4]).Equals(c.GetAttachedEdges()[4]))
	assert.True(c.Vertices[4].EdgeTo(c.Vertices[1]).Equals(c.GetAttachedEdges()[5]))
	assert.True(c.Vertices[1].EdgeTo(c.Vertices[0]).Equals(c.GetAttachedEdges()[6]))

	assert.Equal(tspmodel2d.NewVertex2D(3, 0), c.GetAttachedVertices()[2])
	assert.Equal(tspmodel2d.NewVertex2D(8, 5), c.GetAttachedVertices()[3])
	assert.Equal(tspmodel2d.NewVertex2D(9, 6), c.GetAttachedVertices()[4])

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.Vertices, 8)
	assert.Len(c.GetAttachedVertices(), 8)
	assert.Len(c.GetAttachedEdges(), 8)
	assert.Len(c.GetUnattachedVertices(), 0)

	assert.True(c.Vertices[0].EdgeTo(c.Vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(c.Vertices[7].EdgeTo(c.Vertices[2]).Equals(c.GetAttachedEdges()[1]))
	assert.True(c.Vertices[2].EdgeTo(c.Vertices[3]).Equals(c.GetAttachedEdges()[2]))
	assert.True(c.Vertices[3].EdgeTo(c.Vertices[5]).Equals(c.GetAttachedEdges()[3]))
	assert.True(c.Vertices[5].EdgeTo(c.Vertices[6]).Equals(c.GetAttachedEdges()[4]))
	assert.True(c.Vertices[6].EdgeTo(c.Vertices[4]).Equals(c.GetAttachedEdges()[5]))
	assert.True(c.Vertices[4].EdgeTo(c.Vertices[1]).Equals(c.GetAttachedEdges()[6]))
	assert.True(c.Vertices[1].EdgeTo(c.Vertices[0]).Equals(c.GetAttachedEdges()[7]))

	assert.Equal(tspmodel2d.NewVertex2D(0, 0), c.GetAttachedVertices()[2])
	assert.Equal(tspmodel2d.NewVertex2D(3, 0), c.GetAttachedVertices()[3])
	assert.Equal(tspmodel2d.NewVertex2D(8, 5), c.GetAttachedVertices()[4])
	assert.Equal(tspmodel2d.NewVertex2D(9, 6), c.GetAttachedVertices()[5])

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.Vertices, 8)
	assert.Len(c.GetAttachedVertices(), 8)
	assert.Len(c.GetAttachedEdges(), 8)
	assert.Len(c.GetUnattachedVertices(), 0)
}

func TestUpdate_ShouldNotRemoveAttachedInteriorPointFromPerimeterIfNewEdgeIsCloserThanPreviousEdge_ConvexConcaveDisparity(t *testing.T) {
	assert := assert.New(t)
	c := circuit.NewConvexConcaveDisparity([]tspmodel.CircuitVertex{
		tspmodel2d.NewVertex2D(0, 0),
		tspmodel2d.NewVertex2D(4.7, 2.0),
		tspmodel2d.NewVertex2D(5.0, 2.25),
		tspmodel2d.NewVertex2D(5, 5),
		tspmodel2d.NewVertex2D(6.0, 2.5),
		tspmodel2d.NewVertex2D(10, 0),
	},
		tspmodel2d.DeduplicateVertices,
		tspmodel2d.BuildPerimiter,
		false,
	).(*circuit.ConvexConcaveDisparity)

	c.Prepare()
	c.BuildPerimiter()

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 4)
	assert.Equal(c.Vertices[0], c.GetAttachedVertices()[0])
	assert.Equal(c.Vertices[5], c.GetAttachedVertices()[1])
	assert.Equal(c.Vertices[4], c.GetAttachedVertices()[2])
	assert.Equal(c.Vertices[3], c.GetAttachedVertices()[3])

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 5)
	assert.Equal(c.Vertices[0], c.GetAttachedVertices()[0])
	assert.Equal(c.Vertices[1], c.GetAttachedVertices()[1])
	assert.Equal(c.Vertices[5], c.GetAttachedVertices()[2])
	assert.Equal(c.Vertices[4], c.GetAttachedVertices()[3])
	assert.Equal(c.Vertices[3], c.GetAttachedVertices()[4])

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Equal(c.Vertices[0], c.GetAttachedVertices()[0])
	assert.Equal(c.Vertices[1], c.GetAttachedVertices()[1])
	assert.Equal(c.Vertices[2], c.GetAttachedVertices()[2])
	assert.Equal(c.Vertices[5], c.GetAttachedVertices()[3])
	assert.Equal(c.Vertices[4], c.GetAttachedVertices()[4])
	assert.Equal(c.Vertices[3], c.GetAttachedVertices()[5])

	c.Update(c.FindNextVertexAndEdge())
	assert.Len(c.GetAttachedVertices(), 6)
	assert.Equal(c.Vertices[0], c.GetAttachedVertices()[0])
	assert.Equal(c.Vertices[1], c.GetAttachedVertices()[1])
	assert.Equal(c.Vertices[2], c.GetAttachedVertices()[2])
	assert.Equal(c.Vertices[5], c.GetAttachedVertices()[3])
	assert.Equal(c.Vertices[4], c.GetAttachedVertices()[4])
	assert.Equal(c.Vertices[3], c.GetAttachedVertices()[5])

	v, _ := c.FindNextVertexAndEdge()
	assert.Nil(v)
}

func TestUpdate_ConvexConcaveDisparity_ShouldPanicIfEdgeIsNotInCircuit(t *testing.T) {
	assert := assert.New(t)
	c := circuit.NewConvexConcaveDisparity([]tspmodel.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		tspmodel2d.NewVertex2D(-15, -15), // Index 0 after sorting
		tspmodel2d.NewVertex2D(0, 0),     // Index 2 after sorting
		tspmodel2d.NewVertex2D(15, -15),  // Index 7 after sorting
		tspmodel2d.NewVertex2D(3, 0),     // Index 3 after sorting
		tspmodel2d.NewVertex2D(3, 13),    // Index 4 after sorting
		tspmodel2d.NewVertex2D(8, 5),     // Index 5 after sorting
		tspmodel2d.NewVertex2D(9, 6),     // Index 6 after sorting
		tspmodel2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	},
		tspmodel2d.DeduplicateVertices,
		tspmodel2d.BuildPerimiter,
		false,
	).(*circuit.ConvexConcaveDisparity)

	c.Prepare()
	c.BuildPerimiter()

	assert.Len(c.Vertices, 8)
	assert.Len(c.GetAttachedVertices(), 5)
	assert.Len(c.GetAttachedEdges(), 5)
	assert.Len(c.GetUnattachedVertices(), 3)

	assert.Panics(func() { c.Update(c.Vertices[2], c.Vertices[0].EdgeTo(c.Vertices[4])) })
}
