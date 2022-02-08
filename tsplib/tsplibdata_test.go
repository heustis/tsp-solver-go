package tsplib_test

import (
	"fmt"
	"testing"

	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/fealos/lee-tsp-go/solver"
	"github.com/fealos/lee-tsp-go/tsplib"
	"github.com/stretchr/testify/assert"
)

var filenames []string = []string{
	"a280",
	"ali535",
	"att48",
	"att532",
	"berlin52",
	"bier127",
	"brd14051",
	"burma14",
	"ch130",
	"ch150",
	"d198",
	"d493",
	"d657",
	"d1291",
	"d1655",
	"d2103",
	"d18512",
	"dsj1000",
	"eil51",
	"eil76",
	"eil101",
	"fl417",
	"fl1400",
	"fl1577",
	"fl3795",
	"fnl4461",
	"gil262",
	"gr96",
	"gr137",
	"gr202",
	"gr666",
	"kroA100",
	"kroC100",
	"kroD100",
	"lin105",
	"pcb442",
	"pr76",
	"pr1002",
	"pr2392",
	"rd100",
	"st70",
	"tsp225",
	"ulysses16",
	"ulysses22",
	"usa13509",
}

func TestNewData(t *testing.T) {
	assert := assert.New(t)

	data, err := tsplib.NewData("../test-data/tsplib/not_a_file.tsp")
	assert.NotNil(err)
	assert.Nil(data)

	data, err = tsplib.NewData("../test-data/tsplib/a280.tsp.malformedX")
	assert.NotNil(err)
	assert.Nil(data)

	data, err = tsplib.NewData("../test-data/tsplib/a280.tsp.malformedY")
	assert.NotNil(err)
	assert.Nil(data)

	data, err = tsplib.NewData("../test-data/tsplib/a280.tsp")
	assert.Nil(err)
	assert.NotNil(data)

	assert.Equal(`a280`, data.GetName())
	assert.Equal(`drilling problem (Ludwig)`, data.GetComment())
	assert.Equal(280, data.GetNumPoints())

	vertices := data.GetVertices()
	assert.Len(vertices, 280)

	assert.Equal(model2d.NewVertex2D(288, 149), vertices[0])
	assert.Equal(model2d.NewVertex2D(288, 129), vertices[1])
	assert.Equal(model2d.NewVertex2D(32, 129), vertices[55])
	assert.Equal(model2d.NewVertex2D(280, 133), vertices[279])

	bestRoute := data.GetBestRoute()
	assert.Len(bestRoute, 280)
	assert.Equal(model2d.NewVertex2D(288, 149), bestRoute[0])
	assert.Equal(model2d.NewVertex2D(288, 129), bestRoute[1])
	assert.Equal(model2d.NewVertex2D(288, 109), bestRoute[2])
	assert.Equal(model2d.NewVertex2D(270, 133), bestRoute[278])
	assert.Equal(model2d.NewVertex2D(280, 133), bestRoute[279])

	assert.InDelta(2586.7696475631606, data.GetBestRouteLength(), model.Threshold)
}

func TestNewData_ShouldSupportExponents(t *testing.T) {
	assert := assert.New(t)

	data, err := tsplib.NewData("../test-data/tsplib/d1291.tsp")
	assert.Nil(err)
	assert.NotNil(data)

	vertices := data.GetVertices()
	assert.Len(vertices, 1291)
	assert.Equal(model2d.NewVertex2D(0.0, 0.0), vertices[0])
	assert.Equal(model2d.NewVertex2D(837, 958.3), vertices[1])
}

func TestSolveAndCompare(t *testing.T) {
	for _, filename := range filenames {
		t.Run(filename, func(t *testing.T) {
			fmt.Println(filename)
			assert := assert.New(t)
			data, err := tsplib.NewData(fmt.Sprintf("../test-data/tsplib/%s.tsp", filename))
			assert.Nil(err)
			assert.NotNil(data)

			err = data.SolveAndCompare("concon", func(cv []model.CircuitVertex) model.Circuit {
				c := circuit.NewConvexConcave(model2d.DeduplicateVertices(data.GetVertices()), model2d.BuildPerimiter, false)
				solver.FindShortestPathCircuit(c)
				return c
			})
			assert.Nil(err, filename)

			// err = data.SolveAndCompare("concon_by_edge", func(cv []model.CircuitVertex) model.Circuit {
			// 	c := circuit.NewConvexConcaveByEdge(model2d.DeduplicateVertices(data.GetVertices()), model2d.BuildPerimiter, false)
			// 	solver.FindShortestPathCircuit(c)
			// 	return c
			// })
			// assert.Nil(err)

			// err = data.SolveAndCompare("concon_updates", func(cv []model.CircuitVertex) model.Circuit {
			// 	c := circuit.NewConvexConcave(model2d.DeduplicateVertices(data.GetVertices()), model2d.BuildPerimiter, true)
			// 	solver.FindShortestPathCircuit(c)
			// 	return c
			// })
			// assert.Nil(err)
		})
	}
}

func TestSolveAndCompareConfidence(t *testing.T) {
	t.Skip("Confidence circuits are significantly less performant than greed Concave Convex and Disparity algorithms, but only occasionally provides a couple percentage points of improved accuracy")
	filename := "eil76"
	fmt.Println(filename)
	assert := assert.New(t)
	data, err := tsplib.NewData(fmt.Sprintf("../test-data/tsplib/%s.tsp", filename))
	assert.Nil(err)
	assert.NotNil(data)

	for _, zscore := range []float64{
		0.5,
		0.75,
		1.0,
		1.5,
	} {
		err = data.SolveAndCompare(fmt.Sprintf("confidence_%g_", zscore), func(cv []model.CircuitVertex) model.Circuit {
			c := circuit.NewConvexConcaveConfidence(model2d.DeduplicateVertices(data.GetVertices()), model2d.BuildPerimiter)
			c.(*circuit.ConvexConcaveConfidence).Significance = zscore
			solver.FindShortestPathCircuit(c)
			return c
		})
		assert.Nil(err, filename)
	}
}
