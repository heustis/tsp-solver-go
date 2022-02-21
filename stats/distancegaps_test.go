package stats_test

import (
	"testing"

	"github.com/heustis/tsp-solver-go/model"
	"github.com/heustis/tsp-solver-go/model2d"
	"github.com/heustis/tsp-solver-go/stats"
	"github.com/stretchr/testify/assert"
)

func TestCloneDistanceGaps(t *testing.T) {
	assert := assert.New(t)

	statsA := &stats.DistanceGaps{
		ClosestEdges: []*model.DistanceToEdge{
			{
				Vertex:   model2d.NewVertex2D(2, 3),
				Edge:     model2d.NewEdge2D(model2d.NewVertex2D(1, 1), model2d.NewVertex2D(3, 4)),
				Distance: 3.45,
			},
			{
				Vertex:   model2d.NewVertex2D(-6, 7),
				Edge:     model2d.NewEdge2D(model2d.NewVertex2D(1, -7), model2d.NewVertex2D(-8, 9)),
				Distance: 4.56,
			},
			{
				Vertex:   model2d.NewVertex2D(0, 0),
				Edge:     model2d.NewEdge2D(model2d.NewVertex2D(1, -7), model2d.NewVertex2D(-8, 9)),
				Distance: 7.89,
			},
		},
		Gaps:                 []float64{1.11, 3.33},
		GapAverage:           2.22,
		GapStandardDeviation: .555,
	}

	statsB := statsA.Clone()
	assert.NotNil(statsB)
	assert.Equal(statsA, statsB)

	// Validate that the arrays of closest edges can be updated independently.
	assert.Len(statsA.ClosestEdges, 3)
	assert.Len(statsB.ClosestEdges, 3)

	statsA.ClosestEdges = append(statsA.ClosestEdges, &model.DistanceToEdge{
		Vertex:   model2d.NewVertex2D(10, 10),
		Edge:     model2d.NewEdge2D(model2d.NewVertex2D(1, -7), model2d.NewVertex2D(-8, 9)),
		Distance: 17.89,
	})

	assert.Len(statsA.ClosestEdges, 4)
	assert.Len(statsB.ClosestEdges, 3)

	statsB.ClosestEdges = append(statsB.ClosestEdges, &model.DistanceToEdge{
		Vertex:   model2d.NewVertex2D(15, 15),
		Edge:     model2d.NewEdge2D(model2d.NewVertex2D(3, 3), model2d.NewVertex2D(40, 40)),
		Distance: 0.0,
	})

	assert.Len(statsA.ClosestEdges, 4)
	assert.Len(statsB.ClosestEdges, 4)

	// Validate that the arrays of gaps can be updated independently.
	assert.Len(statsA.Gaps, 2)
	assert.Len(statsB.Gaps, 2)

	statsA.Gaps = append(statsA.Gaps, 23.45)

	assert.Len(statsA.Gaps, 3)
	assert.Len(statsB.Gaps, 2)

	statsB.Gaps = append(statsB.Gaps, 55.55)

	assert.Len(statsA.Gaps, 3)
	assert.Len(statsB.Gaps, 3)
}

func TestNewDistanceGaps(t *testing.T) {
	assert := assert.New(t)

	vertices := []model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(-7, 6),
	}

	edges := []model.CircuitEdge{
		vertices[0].EdgeTo(vertices[2]),
		vertices[2].EdgeTo(vertices[6]),
		vertices[6].EdgeTo(vertices[4]),
		vertices[4].EdgeTo(vertices[7]),
		vertices[7].EdgeTo(vertices[0]),
	}

	stats1 := stats.NewDistanceGaps(vertices[1], edges)
	assert.NotNil(stats1)
	assert.Len(stats1.ClosestEdges, 5)
	assert.Len(stats1.Gaps, 4)

	assert.Equal(stats1.ClosestEdges[0], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[4],
		Distance: 7.9605428386450825,
	})
	assert.Equal(stats1.ClosestEdges[1], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[1],
		Distance: 10.189527594146842,
	})
	assert.Equal(stats1.ClosestEdges[2], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[3],
		Distance: 10.35465290568552,
	})
	assert.Equal(stats1.ClosestEdges[3], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[0],
		Distance: 12.426406871192853,
	})
	assert.Equal(stats1.ClosestEdges[4], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[2],
		Distance: 14.938773433225418,
	})

	assert.Equal([]float64{2.2289847555017595, 0.16512531153867727, 2.071753965507334, 2.5123665620325646}, stats1.Gaps)
	assert.Equal(1.7445576486450838, stats1.GapAverage)
	assert.Equal(0.925454494640393, stats1.GapStandardDeviation)

	stats3 := stats.NewDistanceGaps(vertices[3], edges)
	assert.NotNil(stats3)
	assert.Len(stats3.ClosestEdges, 5)
	assert.Len(stats3.Gaps, 4)

	assert.Equal(stats3.ClosestEdges[0], &model.DistanceToEdge{
		Vertex:   vertices[3],
		Edge:     edges[1],
		Distance: 5.854324418695558,
	})
	assert.Equal(stats3.ClosestEdges[1], &model.DistanceToEdge{
		Vertex:   vertices[3],
		Edge:     edges[2],
		Distance: 12.26573691694568,
	})
	assert.Equal(stats3.ClosestEdges[2], &model.DistanceToEdge{
		Vertex:   vertices[3],
		Edge:     edges[3],
		Distance: 12.4553481739569,
	})
	assert.Equal(stats3.ClosestEdges[3], &model.DistanceToEdge{
		Vertex:   vertices[3],
		Edge:     edges[4],
		Distance: 12.620447763166336,
	})
	assert.Equal(stats3.ClosestEdges[4], &model.DistanceToEdge{
		Vertex:   vertices[3],
		Edge:     edges[0],
		Distance: 12.640121740018508,
	})

	assert.Equal([]float64{6.411412498250122, 0.18961125701121873, 0.16509958920943646, 0.01967397685217165}, stats3.Gaps)
	assert.Equal(1.6964493303307373, stats3.GapAverage)
	assert.Equal(2.7229600745194156, stats3.GapStandardDeviation)
}

func TestStringDistanceGaps(t *testing.T) {
	assert := assert.New(t)

	vertices := []model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(-7, 6),
	}
	vertexLookup := make(map[model.CircuitVertex]int)
	for i, v := range vertices {
		vertexLookup[v] = i
	}

	edges := []model.CircuitEdge{
		vertices[0].EdgeTo(vertices[2]),
		vertices[2].EdgeTo(vertices[6]),
		vertices[6].EdgeTo(vertices[4]),
		vertices[4].EdgeTo(vertices[7]),
		vertices[7].EdgeTo(vertices[0]),
	}
	edgeLookup := make(map[model.CircuitEdge]int)
	for i, e := range edges {
		edgeLookup[e] = i
	}

	assert.Equal("{}", (&stats.DistanceGaps{}).String(vertexLookup, edgeLookup))

	stats1 := stats.NewDistanceGaps(vertices[1], edges)
	assert.Equal("{\r\n\t\"vertex\":1,\r\n\t\"gapAverage\":1.7445576486450838,\r\n\t\"gapStdDev\":0.925454494640393,\r\n\t\"closestEdges\":[{\"edge\":4,\"distance\":7.9605428386450825},{\"edge\":1,\"distance\":10.189527594146842},{\"edge\":3,\"distance\":10.35465290568552},{\"edge\":0,\"distance\":12.426406871192853},{\"edge\":2,\"distance\":14.938773433225418}],\r\n\t\"gaps\":[{\"gap\":2.2289847555017595,\"gapZScore\":0.5234477866412126},{\"gap\":0.16512531153867727,\"gapZScore\":-1.7066558607186104},{\"gap\":2.071753965507334,\"gapZScore\":0.35355203173916167},{\"gap\":2.5123665620325646,\"gapZScore\":0.8296560423382361}]}",
		stats1.String(vertexLookup, edgeLookup))

	stats3 := stats.NewDistanceGaps(vertices[3], edges)
	assert.Equal("{\r\n\t\"vertex\":3,\r\n\t\"gapAverage\":1.6964493303307373,\r\n\t\"gapStdDev\":2.7229600745194156,\r\n\t\"closestEdges\":[{\"edge\":1,\"distance\":5.854324418695558},{\"edge\":2,\"distance\":12.26573691694568},{\"edge\":3,\"distance\":12.4553481739569},{\"edge\":4,\"distance\":12.620447763166336},{\"edge\":0,\"distance\":12.640121740018508}],\r\n\t\"gaps\":[{\"gap\":6.411412498250122,\"gapZScore\":1.7315579512312698},{\"gap\":0.18961125701121873,\"gapZScore\":-0.5533823604025724},{\"gap\":0.16509958920943646,\"gapZScore\":-0.562384206603387},{\"gap\":0.01967397685217165,\"gapZScore\":-0.6157913842253105}]}",
		stats3.String(vertexLookup, edgeLookup))
}

func TestUpdateStatsShouldUpdateStartOfList(t *testing.T) {
	assert := assert.New(t)

	vertices := []model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(-7, 6),
	}

	edges := []model.CircuitEdge{
		vertices[0].EdgeTo(vertices[2]),
		vertices[2].EdgeTo(vertices[6]),
		vertices[6].EdgeTo(vertices[4]),
		vertices[4].EdgeTo(vertices[7]),
		vertices[7].EdgeTo(vertices[0]),
	}

	stats1 := stats.NewDistanceGaps(vertices[1], edges)
	edgeA, edgeB := edges[1].Split(vertices[3])
	statsClone := stats1.Clone()
	stats1.UpdateStats(edges[1], edgeA, edgeB)
	statsClone.UpdateStats(edges[1], edgeB, edgeA)

	assert.Len(stats1.ClosestEdges, 6)
	assert.Len(stats1.Gaps, 5)
	assert.Equal(stats1, statsClone)

	assert.Equal(stats1.ClosestEdges[0], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edgeA,
		Distance: 5.003830723297881,
	})
	assert.Equal(stats1.ClosestEdges[1], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edgeB,
		Distance: 5.331372452153399,
	})
	assert.Equal(stats1.ClosestEdges[2], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[4],
		Distance: 7.9605428386450825,
	})
	assert.Equal(stats1.ClosestEdges[3], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[3],
		Distance: 10.35465290568552,
	})
	assert.Equal(stats1.ClosestEdges[4], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[0],
		Distance: 12.426406871192853,
	})
	assert.Equal(stats1.ClosestEdges[5], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[2],
		Distance: 14.938773433225418,
	})

	assert.Equal([]float64{0.32754172885551824, 2.6291703864916833, 2.3941100670404367, 2.071753965507334, 2.5123665620325646}, stats1.Gaps)
	assert.Equal(1.9869885419855073, stats1.GapAverage)
	assert.Equal(0.8503077588916949, stats1.GapStandardDeviation)
}

func TestUpdateStats_ShouldUpdateMiddleOfList(t *testing.T) {
	assert := assert.New(t)

	vertices := []model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(-7, 6),
	}

	edges := []model.CircuitEdge{
		vertices[0].EdgeTo(vertices[2]),
		vertices[2].EdgeTo(vertices[6]),
		vertices[6].EdgeTo(vertices[4]),
		vertices[4].EdgeTo(vertices[7]),
		vertices[7].EdgeTo(vertices[0]),
	}

	stats1 := stats.NewDistanceGaps(vertices[1], edges)
	stats1.ClosestEdges[0].Distance *= 0.1
	stats1.ClosestEdges[1].Distance *= 0.1
	stats1.ClosestEdges[2].Distance *= 0.1

	edgeA, edgeB := edges[1].Split(vertices[3])
	statsClone := stats1.Clone()
	stats1.UpdateStats(edges[1], edgeA, edgeB)
	statsClone.UpdateStats(edges[1], edgeB, edgeA)

	assert.Len(stats1.ClosestEdges, 6)
	assert.Len(stats1.Gaps, 5)
	assert.Equal(stats1, statsClone)

	assert.Equal(stats1.ClosestEdges[0], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[4],
		Distance: 0.7960542838645083,
	})
	assert.Equal(stats1.ClosestEdges[1], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[3],
		Distance: 1.035465290568552,
	})
	assert.Equal(stats1.ClosestEdges[2], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edgeA,
		Distance: 5.003830723297881,
	})
	assert.Equal(stats1.ClosestEdges[3], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edgeB,
		Distance: 5.331372452153399,
	})
	assert.Equal(stats1.ClosestEdges[4], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[0],
		Distance: 12.426406871192853,
	})
	assert.Equal(stats1.ClosestEdges[5], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[2],
		Distance: 14.938773433225418,
	})

	assert.Equal([]float64{0.2394110067040437, 3.968365432729329, 0.32754172885551824, 7.095034419039454, 2.5123665620325646}, stats1.Gaps)
	assert.Equal(2.828543829872182, stats1.GapAverage)
	assert.Equal(2.5518904202095998, stats1.GapStandardDeviation)
}

func TestUpdateStats_ShouldUpdateEndOfList(t *testing.T) {
	assert := assert.New(t)

	vertices := []model.CircuitVertex{
		model2d.NewVertex2D(-15, -15),
		model2d.NewVertex2D(0, 0),
		model2d.NewVertex2D(15, -15),
		model2d.NewVertex2D(3, 0),
		model2d.NewVertex2D(3, 13),
		model2d.NewVertex2D(8, 5),
		model2d.NewVertex2D(9, 6),
		model2d.NewVertex2D(-7, 6),
	}

	edges := []model.CircuitEdge{
		vertices[0].EdgeTo(vertices[2]),
		vertices[2].EdgeTo(vertices[6]),
		vertices[6].EdgeTo(vertices[4]),
		vertices[4].EdgeTo(vertices[7]),
		vertices[7].EdgeTo(vertices[0]),
	}

	stats1 := stats.NewDistanceGaps(vertices[1], edges)
	stats1.ClosestEdges[0].Distance *= 0.1
	stats1.ClosestEdges[1].Distance *= 0.1
	stats1.ClosestEdges[2].Distance *= 0.1
	stats1.ClosestEdges[3].Distance *= 0.1
	stats1.ClosestEdges[4].Distance *= 0.1

	edgeA, edgeB := edges[1].Split(vertices[3])
	statsClone := stats1.Clone()
	stats1.UpdateStats(edges[1], edgeA, edgeB)
	statsClone.UpdateStats(edges[1], edgeB, edgeA)

	assert.Len(stats1.ClosestEdges, 6)
	assert.Len(stats1.Gaps, 5)
	assert.Equal(stats1, statsClone)

	assert.Equal(stats1.ClosestEdges[0], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[4],
		Distance: 0.7960542838645083,
	})
	assert.Equal(stats1.ClosestEdges[1], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[3],
		Distance: 1.035465290568552,
	})
	assert.Equal(stats1.ClosestEdges[2], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[0],
		Distance: 1.2426406871192854,
	})
	assert.Equal(stats1.ClosestEdges[3], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edges[2],
		Distance: 1.4938773433225419,
	})
	assert.Equal(stats1.ClosestEdges[4], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edgeA,
		Distance: 5.003830723297881,
	})
	assert.Equal(stats1.ClosestEdges[5], &model.DistanceToEdge{
		Vertex:   vertices[1],
		Edge:     edgeB,
		Distance: 5.331372452153399,
	})

	assert.Equal([]float64{0.2394110067040437, 0.20717539655073347, 0.2512366562032564, 3.509953379975339, 0.32754172885551824}, stats1.Gaps)
	assert.Equal(0.9070636336577782, stats1.GapAverage)
	assert.Equal(1.302044029110145, stats1.GapStandardDeviation)
}
