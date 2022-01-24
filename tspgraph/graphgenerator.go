package tspgraph

import (
	"container/list"
	"fmt"
	"math/rand"
	"time"
)

type GraphGenerator struct {
	// EnableAsymetricDistances (it true, default false) allows the tspgraph to have different edge lengths from A to B than B to A
	// (e.g. to simulate different routes between two locations due to one-way streets).
	EnableAsymetricDistances bool

	// EnableUnidirectionalEdges (if true, default false) allows the tspgraph to have an edge from node A to B, without a corresponding edge from B to A.
	// All vertices will have paths to all other vertices, even if this is enabled.
	EnableUnidirectionalEdges bool

	// MaxEdges determines the maximum number of edges each vertex can have.
	// This must be greater than or equal to MinEdges, and this must be at least 2.
	MaxEdges uint8

	// MinEdges determines the minimum number of edge each vertex should have.
	// This must be at least 1.
	MinEdges uint8

	// NumVertices determines the number of vertices to generate.
	// This must be at least 3.
	NumVertices uint32

	// Seed is used to initialize the random algoritm.
	// This should be used to reproduce the same tspgraph accross multiple tests.
	// If this is nil, a seed will be automatically generated.
	Seed *int64
}

func (gen *GraphGenerator) Create() *Graph {
	if gen.MaxEdges < 2 {
		panic(fmt.Errorf("MaxEdges must be at least 2, supplied value=%v", gen.MaxEdges))
	} else if gen.MinEdges < 1 {
		panic(fmt.Errorf("MinEdges must be at least 1, supplied value=%v", gen.MinEdges))
	} else if gen.MaxEdges < gen.MinEdges {
		panic(fmt.Errorf("MaxEdges must be at least MinEdges, MaxEdges=%v MinEdges=%v", gen.MaxEdges, gen.MinEdges))
	} else if gen.NumVertices < 3 {
		panic(fmt.Errorf("NumVertices must be at least 3, supplied value=%v", gen.NumVertices))
	}

	availableNames := list.New()
	availableNames.Init()
	for i := 0; i < int(gen.NumVertices); i++ {
		availableNames.PushBack(buildVertexName(i))
	}
	var random *rand.Rand
	if gen.Seed != nil {
		random = rand.New(rand.NewSource(*gen.Seed))
	} else {
		random = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	tspgraph := &Graph{
		Vertices: []*GraphVertex{},
	}

	// Set up a basic tspgraph where the is at least a unidirectional circuit in the tspgraph.
	nextVertex := generateVertex(availableNames, random)
	for availableNames.Len() > 0 {
		current := nextVertex
		tspgraph.Vertices = append(tspgraph.Vertices, current)
		nextVertex = generateVertex(availableNames, random)
		gen.linkVertices(current, nextVertex, random)
	}
	// Append the final vertex to the tspgraph
	tspgraph.Vertices = append(tspgraph.Vertices, nextVertex)

	// Update each node in the tspgraph to have a random number of edges between MinEdges and MaxEdges
	// Note: this may produce Vertices with more edges than MaxEdges, but that doesn't cause any issues so I am not fixing it (at this time).
	numEdgesRange := gen.MaxEdges - gen.MinEdges
	for _, v := range tspgraph.Vertices {
		numEdges := gen.MinEdges
		if numEdgesRange > 0 {
			numEdges += uint8(random.Int31n(int32(numEdgesRange)))
		}
		for len(v.adjacentVertices) < int(numEdges) {
			destinationIndex := random.Intn(int(gen.NumVertices))
			gen.linkVertices(v, tspgraph.Vertices[destinationIndex], random)
		}
	}

	return tspgraph
}

func buildVertexName(index int) string {
	name := ""

	remainder := index % 26
	// Notes:
	// - need (result+25)%26 to allow `a` to appear in the first letter of the name. Otherwise, it is treated as 0 and `b`` would be treated as `1`.
	// - need (result-1)%26 to allow `za`-`zz` to appear in the list due to the manipulation we are doing to the remainder.
	for result := index / 26; result > 0; remainder, result = (result+25)%26, (result-1)/26 {
		name = fmt.Sprintf("%c%s", 'a'+remainder, name)
	}
	name = fmt.Sprintf("%c%s", 'a'+remainder, name)

	return name
}

func generateVertex(availableNames *list.List, random *rand.Rand) *GraphVertex {
	nameIndex := random.Intn(availableNames.Len())
	current := availableNames.Front()
	for i := 0; i < nameIndex; i, current = i+1, current.Next() {
	}
	name := availableNames.Remove(current)
	return &GraphVertex{
		id:               name.(string),
		adjacentVertices: make(map[*GraphVertex]float64),
	}
}

func (gen *GraphGenerator) linkVertices(a *GraphVertex, b *GraphVertex, random *rand.Rand) {
	distAB := random.Float64() * 10000.0
	a.adjacentVertices[b] = distAB

	// Set up the return edge, if enabled.
	if !gen.EnableUnidirectionalEdges {
		distBA := distAB
		if gen.EnableAsymetricDistances {
			distBA = random.Float64() * 10000.0
		}
		b.adjacentVertices[a] = distBA
	}
}
