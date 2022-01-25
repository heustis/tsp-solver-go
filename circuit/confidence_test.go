package circuit_test

import (
	"testing"

	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/stretchr/testify/assert"
)

func TestBuildPerimeter_ConvexConcaveConfidence(t *testing.T) {
	assert := assert.New(t)
	c := circuit.NewConvexConcaveConfidence([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	},
		model2d.DeduplicateVertices,
		model2d.BuildPerimiter,
	).(*circuit.ConvexConcaveConfidence)

	c.Prepare()
	c.BuildPerimiter()

	assert.Len(c.Vertices, 8)

	vertices := c.GetAttachedVertices()
	assert.Len(vertices, 5)
	assert.Equal(model2d.NewVertex2D(-15, -15), vertices[0])
	assert.Equal(model2d.NewVertex2D(15, -15), vertices[1])
	assert.Equal(model2d.NewVertex2D(9, 6), vertices[2])
	assert.Equal(model2d.NewVertex2D(3, 13), vertices[3])
	assert.Equal(model2d.NewVertex2D(-7, 6), vertices[4])

	assert.InDelta(95.73863479511238, c.GetLength(), model.Threshold)

	edges := c.GetAttachedEdges()
	assert.Len(edges, 5)
	assert.True(c.Vertices[0].EdgeTo(c.Vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(c.Vertices[7].EdgeTo(c.Vertices[6]).Equals(c.GetAttachedEdges()[1]))
	assert.True(c.Vertices[6].EdgeTo(c.Vertices[4]).Equals(c.GetAttachedEdges()[2]))
	assert.True(c.Vertices[4].EdgeTo(c.Vertices[1]).Equals(c.GetAttachedEdges()[3]))
	assert.True(c.Vertices[1].EdgeTo(c.Vertices[0]).Equals(c.GetAttachedEdges()[4]))

	unattached := c.GetUnattachedVertices()
	assert.Len(unattached, 3)
	assert.True(unattached[c.Vertices[2]])
	assert.True(unattached[c.Vertices[3]])
	assert.True(unattached[c.Vertices[5]])
}

func TestPrepare_ConvexConcaveConfidence(t *testing.T) {
	assert := assert.New(t)
	c := circuit.NewConvexConcaveConfidence([]model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(-15-model.Threshold/3.0, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(-15+model.Threshold/3.0, -15-model.Threshold/3.0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(7, 6),
		model2d.NewVertex2D(-7, 6),
	},
		model2d.DeduplicateVertices,
		model2d.BuildPerimiter,
	).(*circuit.ConvexConcaveConfidence)

	c.Prepare()

	assert.NotNil(c.Vertices)
	assert.Len(c.Vertices, 7)
	assert.ElementsMatch(c.Vertices, []model.CircuitVertex{
		model2d.NewVertex2D(-15+model.Threshold/3.0, -15-model.Threshold/3.0),
		model2d.NewVertex2D(-7, 6),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(7, 6),
		model2d.NewVertex2D(15, -15),
	})

	assert.NotNil(c.GetUnattachedVertices())
	assert.Len(c.GetUnattachedVertices(), 0)

	assert.NotNil(c.GetAttachedVertices())
	assert.Len(c.GetAttachedVertices(), 0)

	assert.NotNil(c.GetAttachedEdges())
	assert.Len(c.GetAttachedEdges(), 0)

	assert.Equal(0.0, c.GetLength())

	v, e := c.FindNextVertexAndEdge()
	assert.Nil(v)
	assert.Nil(e)
}

func TestUpdate_ConvexConcaveConfidence(t *testing.T) {
	assert := assert.New(t)
	c := circuit.NewConvexConcaveConfidence([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	},
		model2d.DeduplicateVertices,
		model2d.BuildPerimiter,
	).(*circuit.ConvexConcaveConfidence)

	c.Prepare()
	c.BuildPerimiter()

	assert.Len(c.Vertices, 8)
	assert.Len(c.GetAttachedVertices(), 5)
	assert.Len(c.GetAttachedEdges(), 5)
	assert.Len(c.GetUnattachedVertices(), 3)

	c.Update(c.FindNextVertexAndEdge())
	edges := c.GetAttachedEdges()
	vertices := c.GetAttachedVertices()

	assert.Len(c.Vertices, 8)
	assert.Len(vertices, 6)
	assert.Len(edges, 6)
	assert.Len(c.GetUnattachedVertices(), 2)

	assert.True(c.Vertices[0].EdgeTo(c.Vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(c.Vertices[7].EdgeTo(c.Vertices[3]).Equals(c.GetAttachedEdges()[1]))
	assert.True(c.Vertices[3].EdgeTo(c.Vertices[6]).Equals(c.GetAttachedEdges()[2]))
	assert.True(c.Vertices[6].EdgeTo(c.Vertices[4]).Equals(c.GetAttachedEdges()[3]))
	assert.True(c.Vertices[4].EdgeTo(c.Vertices[1]).Equals(c.GetAttachedEdges()[4]))
	assert.True(c.Vertices[1].EdgeTo(c.Vertices[0]).Equals(c.GetAttachedEdges()[5]))

	assert.Equal(model2d.NewVertex2D(15, -15), vertices[1])
	assert.Equal(model2d.NewVertex2D(3, 0), vertices[2])
	assert.Equal(model2d.NewVertex2D(9, 6), vertices[3])

	c.Update(c.FindNextVertexAndEdge())
	edges = c.GetAttachedEdges()
	vertices = c.GetAttachedVertices()

	assert.Len(c.Vertices, 8)
	assert.Len(vertices, 7)
	assert.Len(edges, 7)
	assert.Len(c.GetUnattachedVertices(), 1)

	assert.True(c.Vertices[0].EdgeTo(c.Vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(c.Vertices[7].EdgeTo(c.Vertices[3]).Equals(c.GetAttachedEdges()[1]))
	assert.True(c.Vertices[3].EdgeTo(c.Vertices[5]).Equals(c.GetAttachedEdges()[2]))
	assert.True(c.Vertices[5].EdgeTo(c.Vertices[6]).Equals(c.GetAttachedEdges()[3]))
	assert.True(c.Vertices[6].EdgeTo(c.Vertices[4]).Equals(c.GetAttachedEdges()[4]))
	assert.True(c.Vertices[4].EdgeTo(c.Vertices[1]).Equals(c.GetAttachedEdges()[5]))
	assert.True(c.Vertices[1].EdgeTo(c.Vertices[0]).Equals(c.GetAttachedEdges()[6]))

	assert.Equal(model2d.NewVertex2D(3, 0), vertices[2])
	assert.Equal(model2d.NewVertex2D(8, 5), vertices[3])
	assert.Equal(model2d.NewVertex2D(9, 6), vertices[4])

	c.Update(c.FindNextVertexAndEdge())
	edges = c.GetAttachedEdges()
	vertices = c.GetAttachedVertices()

	assert.Len(c.Vertices, 8)
	assert.Len(vertices, 8)
	assert.Len(edges, 8)
	assert.Len(c.GetUnattachedVertices(), 0)

	assert.True(c.Vertices[0].EdgeTo(c.Vertices[7]).Equals(c.GetAttachedEdges()[0]))
	assert.True(c.Vertices[7].EdgeTo(c.Vertices[2]).Equals(c.GetAttachedEdges()[1]))
	assert.True(c.Vertices[2].EdgeTo(c.Vertices[3]).Equals(c.GetAttachedEdges()[2]))
	assert.True(c.Vertices[3].EdgeTo(c.Vertices[5]).Equals(c.GetAttachedEdges()[3]))
	assert.True(c.Vertices[5].EdgeTo(c.Vertices[6]).Equals(c.GetAttachedEdges()[4]))
	assert.True(c.Vertices[6].EdgeTo(c.Vertices[4]).Equals(c.GetAttachedEdges()[5]))
	assert.True(c.Vertices[4].EdgeTo(c.Vertices[1]).Equals(c.GetAttachedEdges()[6]))
	assert.True(c.Vertices[1].EdgeTo(c.Vertices[0]).Equals(c.GetAttachedEdges()[7]))

	assert.Equal(model2d.NewVertex2D(0, 0), vertices[2])
	assert.Equal(model2d.NewVertex2D(3, 0), vertices[3])

	c.Update(c.FindNextVertexAndEdge())

	assert.Len(c.Vertices, 8)
	assert.Len(c.GetAttachedVertices(), 8)
	assert.Len(c.GetAttachedEdges(), 8)
	assert.Len(c.GetUnattachedVertices(), 0)
}

func TestUpdate_ShouldNotRemoveAttachedInteriorPointFromPerimeterIfNewEdgeIsCloserThanPreviousEdge_ConvexConcaveConfidence(t *testing.T) {
	assert := assert.New(t)
	c := circuit.NewConvexConcaveConfidence([]model.CircuitVertex{
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(4.7, 2.0),
		model2d.NewVertex2D(5.0, 2.25),
		model2d.NewVertex2D(5, 5),
		model2d.NewVertex2D(6.0, 2.5),
		model2d.NewVertex2D(10, 0),
	},
		model2d.DeduplicateVertices,
		model2d.BuildPerimiter,
	).(*circuit.ConvexConcaveConfidence)

	c.Prepare()
	c.BuildPerimiter()

	c.Update(c.FindNextVertexAndEdge())
	vertices := c.GetAttachedVertices()
	assert.Len(vertices, 4)
	assert.Equal(c.Vertices[0], vertices[0])
	assert.Equal(c.Vertices[2], vertices[1])
	assert.Equal(c.Vertices[5], vertices[2])
	assert.Equal(c.Vertices[3], vertices[3])

	c.Update(c.FindNextVertexAndEdge())
	vertices = c.GetAttachedVertices()
	assert.Len(vertices, 5)
	assert.Equal(c.Vertices[0], vertices[0])
	assert.Equal(c.Vertices[2], vertices[1])
	assert.Equal(c.Vertices[4], vertices[2])
	assert.Equal(c.Vertices[5], vertices[3])
	assert.Equal(c.Vertices[3], vertices[4])

	c.Update(c.FindNextVertexAndEdge())
	vertices = c.GetAttachedVertices()
	assert.Len(vertices, 6)
	assert.Equal(c.Vertices[0], vertices[0])
	assert.Equal(c.Vertices[1], vertices[1])
	assert.Equal(c.Vertices[2], vertices[2])
	assert.Equal(c.Vertices[4], vertices[3])
	assert.Equal(c.Vertices[5], vertices[4])
	assert.Equal(c.Vertices[3], vertices[5])

	c.Update(c.FindNextVertexAndEdge())
	vertices = c.GetAttachedVertices()
	assert.Len(vertices, 6)
	assert.Equal(c.Vertices[0], vertices[0])
	assert.Equal(c.Vertices[1], vertices[1])
	assert.Equal(c.Vertices[2], vertices[2])
	assert.Equal(c.Vertices[4], vertices[3])
	assert.Equal(c.Vertices[5], vertices[4])
	assert.Equal(c.Vertices[3], vertices[5])

	v, _ := c.FindNextVertexAndEdge()
	assert.Nil(v)
}

func TestString_ConvexConcaveConfidence(t *testing.T) {
	assert := assert.New(t)
	c := circuit.NewConvexConcaveConfidence([]model.CircuitVertex{
		// Note: the circuit is sorted by Prepare(), so the indices will change as specified below.
		model2d.NewVertex2D(-15, -15), // Index 0 after sorting
		model2d.NewVertex2D(0, 0),     // Index 2 after sorting
		model2d.NewVertex2D(15, -15),  // Index 7 after sorting
		model2d.NewVertex2D(3, 0),     // Index 3 after sorting
		model2d.NewVertex2D(3, 13),    // Index 4 after sorting
		model2d.NewVertex2D(8, 5),     // Index 5 after sorting
		model2d.NewVertex2D(9, 6),     // Index 6 after sorting
		model2d.NewVertex2D(-7, 6),    // Index 1 after sorting
	},
		model2d.DeduplicateVertices,
		model2d.BuildPerimiter,
	).(*circuit.ConvexConcaveConfidence)

	assert.Equal("{\r\n\t\"vertices\":[{\"x\":-15,\"y\":-15},{\"x\":0,\"y\":0},{\"x\":15,\"y\":-15},{\"x\":3,\"y\":0},{\"x\":3,\"y\":13},{\"x\":8,\"y\":5},{\"x\":9,\"y\":6},{\"x\":-7,\"y\":6}],\r\n\t\"edges\":[],\r\n\t\"edgeDistances\":[]}", c.String())

	c.Prepare()
	c.BuildPerimiter()

	s := c.String()
	assert.Contains(s, "{\r\n\t\"vertices\":[{\"x\":-15,\"y\":-15},{\"x\":-7,\"y\":6},{\"x\":0,\"y\":0},{\"x\":3,\"y\":0},{\"x\":3,\"y\":13},{\"x\":8,\"y\":5},{\"x\":9,\"y\":6},{\"x\":15,\"y\":-15}],")
	assert.Contains(s, "\r\n\t\"edges\":[{\"start\":0,\"end\":7,\"distance\":30},{\"start\":7,\"end\":6,\"distance\":21.840329667841555},{\"start\":6,\"end\":4,\"distance\":9.219544457292887},{\"start\":4,\"end\":1,\"distance\":12.206555615733702},{\"start\":1,\"end\":0,\"distance\":22.47220505424423}],")
	assert.Contains(s, "\r\n\t\"edgeDistances\":[")
	assert.Contains(s, "{\r\n\t\"vertex\":2,\r\n\t\"gapAverage\":1.7445576486450838,\r\n\t\"gapStdDev\":0.925454494640393,\r\n\t\"closestEdges\":[{\"edge\":4,\"distance\":7.9605428386450825},{\"edge\":1,\"distance\":10.189527594146842},{\"edge\":3,\"distance\":10.35465290568552},{\"edge\":0,\"distance\":12.426406871192853},{\"edge\":2,\"distance\":14.938773433225418}],\r\n\t\"gaps\":[{\"gap\":2.2289847555017595,\"gapZScore\":0.5234477866412126},{\"gap\":0.16512531153867727,\"gapZScore\":-1.7066558607186104},{\"gap\":2.071753965507334,\"gapZScore\":0.35355203173916167},{\"gap\":2.5123665620325646,\"gapZScore\":0.8296560423382361}]}")
	assert.Contains(s, "{\r\n\t\"vertex\":5,\r\n\t\"gapAverage\":5.569272159359096,\r\n\t\"gapStdDev\":4.475521504357779,\r\n\t\"closestEdges\":[{\"edge\":1,\"distance\":0.763503994948632},{\"edge\":2,\"distance\":1.628650237136812},{\"edge\":3,\"distance\":12.26072189469581},{\"edge\":0,\"distance\":21.669121408673433},{\"edge\":4,\"distance\":23.040592632385017}],\r\n\t\"gaps\":[{\"gap\":0.8651462421881799,\"gapZScore\":-1.0510788323976428},{\"gap\":10.632071657558997,\"gapZScore\":1.1312200138621373},{\"gap\":9.408399513977624,\"gapZScore\":0.8578055877690233},{\"gap\":1.371471223711584,\"gapZScore\":-0.9379467692335178}]}")
	assert.Contains(s, "{\r\n\t\"vertex\":3,\r\n\t\"gapAverage\":1.6964493303307373,\r\n\t\"gapStdDev\":2.7229600745194156,\r\n\t\"closestEdges\":[{\"edge\":1,\"distance\":5.854324418695558},{\"edge\":2,\"distance\":12.26573691694568},{\"edge\":3,\"distance\":12.4553481739569},{\"edge\":4,\"distance\":12.620447763166336},{\"edge\":0,\"distance\":12.640121740018508}],\r\n\t\"gaps\":[{\"gap\":6.411412498250122,\"gapZScore\":1.7315579512312698},{\"gap\":0.18961125701121873,\"gapZScore\":-0.5533823604025724},{\"gap\":0.16509958920943646,\"gapZScore\":-0.562384206603387},{\"gap\":0.01967397685217165,\"gapZScore\":-0.6157913842253105}]}")

}
