package graph_test

import (
	"fmt"
	"testing"

	"github.com/heustis/lee-tsp-go/graph"
	"github.com/stretchr/testify/assert"
)

func TestCreate_ShouldPanicIfMaxEdgesIsLessThan2(t *testing.T) {
	assert := assert.New(t)

	for i := 0; i < 2; i++ {
		assert.Panics(func() {
			g := &graph.GraphGenerator{
				MaxEdges:    uint8(i),
				MinEdges:    1,
				NumVertices: 10,
			}
			g.Create()
		}, i)
	}
}

func TestCreate_ShouldPanicIfMinEdgesIsLessThan1(t *testing.T) {
	assert := assert.New(t)

	assert.Panics(func() {
		g := &graph.GraphGenerator{
			MaxEdges:    2,
			MinEdges:    0,
			NumVertices: 10,
		}
		g.Create()
	})
}

func TestCreate_ShouldPanicIfMinEdgesIsGreaterThanMaxEdges(t *testing.T) {
	assert := assert.New(t)

	assert.Panics(func() {
		g := &graph.GraphGenerator{
			MaxEdges:    2,
			MinEdges:    4,
			NumVertices: 10,
		}
		g.Create()
	})
}

func TestCreate_ShouldPanicIfNumVerticesIsLessThan3(t *testing.T) {
	assert := assert.New(t)

	for i := 0; i < 3; i++ {
		assert.Panics(func() {
			g := &graph.GraphGenerator{
				MaxEdges:    4,
				MinEdges:    2,
				NumVertices: uint32(i),
			}
			g.Create()
		}, i)
	}
}

func TestCreate_ShouldProduceVerticesWithPathsToAllOtherVertices(t *testing.T) {
	assert := assert.New(t)

	gen := &graph.GraphGenerator{
		MaxEdges:    4,
		MinEdges:    2,
		NumVertices: uint32(10),
	}

	g := gen.Create()
	defer g.Delete()

	assert.NotNil(g)
	assert.Len(g.GetVertices(), 10)
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if j != i {
				e := g.GetVertices()[i].EdgeTo(g.GetVertices()[j])
				assert.NotNil(e, fmt.Sprintf(`cannot create path fromIndex=%d toIndex=%d graph=%s`, i, j, g.String()))
			}
		}
	}
}

func TestCreate_ShouldProduceVerticesWithoutDuplicateNames(t *testing.T) {
	assert := assert.New(t)

	gen := &graph.GraphGenerator{
		MaxEdges:    4,
		MinEdges:    2,
		NumVertices: uint32(30),
	}

	g := gen.Create()
	defer g.Delete()

	assert.NotNil(g)
	assert.Len(g.GetVertices(), 30)

	expectedNamesArray := []string{`a`, `b`, `c`, `d`, `e`, `f`, `g`, `h`, `i`, `j`, `k`, `l`, `m`, `n`, `o`, `p`, `q`, `r`, `s`, `t`, `u`, `v`, `w`, `x`, `y`, `z`, `aa`, `ab`, `ac`, `ad`}
	expectedNames := make(map[string]bool)
	for _, name := range expectedNamesArray {
		expectedNames[name] = true
	}

	for _, v := range g.GetVertices() {
		assert.True(expectedNames[v.GetId()], fmt.Sprintf(`unexpected or duplicate name for vertex, name=%s graph=%s`, v.GetId(), g.String()))
		delete(expectedNames, v.GetId())
	}

	assert.Len(expectedNames, 0)
}

func TestCreate_ShouldProduceVerticesWithCorrectNames(t *testing.T) {
	assert := assert.New(t)

	gen := &graph.GraphGenerator{
		MaxEdges:    4,
		MinEdges:    2,
		NumVertices: uint32(26 * 27),
	}

	g := gen.Create()
	defer g.Delete()

	assert.NotNil(g)

	alphabet := []string{`a`, `b`, `c`, `d`, `e`, `f`, `g`, `h`, `i`, `j`, `k`, `l`, `m`, `n`, `o`, `p`, `q`, `r`, `s`, `t`, `u`, `v`, `w`, `x`, `y`, `z`}
	expectedNames := make(map[string]bool)

	for _, firstLetter := range alphabet {
		expectedNames[firstLetter] = true
		for _, secondLetter := range alphabet {
			expectedNames[firstLetter+secondLetter] = true
		}
	}

	assert.Len(g.GetVertices(), len(expectedNames))

	for _, v := range g.GetVertices() {
		assert.True(expectedNames[v.GetId()], fmt.Sprintf(`unexpected or duplicate name for vertex, name=%s`, v.GetId()))
		delete(expectedNames, v.GetId())
	}

	assert.Len(expectedNames, 0)
}

func TestCreate_ShouldProduceSameGraphForSameSeed(t *testing.T) {
	assert := assert.New(t)

	var seed int64
	for seed = -5; seed < 5; seed++ {
		gen := &graph.GraphGenerator{
			MaxEdges:    4,
			MinEdges:    2,
			NumVertices: uint32(30),
			Seed:        &seed,
		}

		g := gen.Create()
		g2 := gen.Create()

		assert.NotNil(g)
		assert.NotNil(g2)
		assert.Len(g.GetVertices(), 30)
		assert.Len(g2.GetVertices(), 30)

		assert.Equal(g.String(), g2.String(), seed)

		g.Delete()
		g2.Delete()
	}
}

func TestCreate_ShouldProduceVerticesWithAppropriateNumberOfEdges(t *testing.T) {
	assert := assert.New(t)

	for min := 2; min < 10; min++ {
		for max := 0; max < 10; max++ {
			gen := &graph.GraphGenerator{
				MaxEdges:    uint8(min + max),
				MinEdges:    uint8(min),
				NumVertices: uint32(30),
			}

			g := gen.Create()

			assert.NotNil(g)
			assert.Len(g.GetVertices(), 30)
			for _, v := range g.GetVertices() {
				adj := v.GetAdjacentVertices()
				assert.GreaterOrEqual(len(adj), min, fmt.Sprintf(`vertex has too few adjacent vertices, name=%s graph=%s`, v.GetId(), g.String()))
				for other := range adj {
					symmetricDist, symmetricEdgeExists := other.GetAdjacentVertices()[v]
					assert.True(symmetricEdgeExists, fmt.Sprintf(`vertex has path without symmetric return path, sourceVertex=%s destinationVertex=%s graph=%s`, v.GetId(), other.GetId(), g.String()))
					assert.Equal(adj[other], symmetricDist)
				}
			}

			g.Delete()
		}
	}
}

func TestCreate_ShouldProduceAsymmetricEdgesIfEnabled(t *testing.T) {
	assert := assert.New(t)

	min := 5
	seed := int64(12456)
	gen := &graph.GraphGenerator{
		MaxEdges:                 uint8(10),
		MinEdges:                 uint8(min),
		NumVertices:              uint32(30),
		EnableAsymetricDistances: true,
		Seed:                     &seed,
	}

	g := gen.Create()

	assert.NotNil(g)
	assert.Len(g.GetVertices(), 30)

	hasAsymmetricEdge := false
	for _, v := range g.GetVertices() {
		adj := v.GetAdjacentVertices()
		assert.GreaterOrEqual(len(adj), min, fmt.Sprintf(`vertex has too few adjacent vertices, name=%s graph=%s`, v.GetId(), g.String()))
		for other := range adj {
			returnDist, returnEdgeExists := other.GetAdjacentVertices()[v]
			assert.True(returnEdgeExists, fmt.Sprintf(`vertex has path without return path, sourceVertex=%s destinationVertex=%s graph=%s`, v.GetId(), other.GetId(), g.String()))
			if !hasAsymmetricEdge {
				hasAsymmetricEdge = returnDist != adj[other]
			}
		}
	}
	assert.True(hasAsymmetricEdge, fmt.Sprintf(`graph did not generate asymetric edges though it was configured to %s`, g.String()))

	g.Delete()
}

func TestCreate_ShouldProduceUnidirectionalEdgesIfEnabled(t *testing.T) {
	assert := assert.New(t)

	min := 5
	seed := int64(12456)
	gen := &graph.GraphGenerator{
		MaxEdges:                  uint8(10),
		MinEdges:                  uint8(min),
		NumVertices:               uint32(30),
		EnableUnidirectionalEdges: true,
		Seed:                      &seed,
	}

	g := gen.Create()

	assert.NotNil(g)
	assert.Len(g.GetVertices(), 30)

	hasUnidirectionalEdge := false
	for _, v := range g.GetVertices() {
		adj := v.GetAdjacentVertices()
		assert.GreaterOrEqual(len(adj), min, fmt.Sprintf(`vertex has too few adjacent vertices, name=%s graph=%s`, v.GetId(), g.String()))
		for other := range adj {
			if !hasUnidirectionalEdge {
				_, returnEdgeExists := other.GetAdjacentVertices()[v]
				hasUnidirectionalEdge = !returnEdgeExists
			}
		}
	}
	assert.True(hasUnidirectionalEdge, fmt.Sprintf(`graph did not generate unidirectional edges though it was configured to %s`, g.String()))

	// Validate that even with unidirectional edges, we can create a path from any vertex to any other vertex.
	for i := 0; i < 30; i++ {
		for j := 0; j < 30; j++ {
			if j != i {
				e := g.GetVertices()[i].EdgeTo(g.GetVertices()[j])
				assert.NotNil(e, fmt.Sprintf(`cannot create path fromIndex=%d toIndex=%d graph=%s`, i, j, g.String()))
			}
		}
	}

	g.Delete()
}
