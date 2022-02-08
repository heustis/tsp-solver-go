package circuit

import (
	"math/rand"
	"sort"
	"time"

	"github.com/fealos/lee-tsp-go/model"
)

// GeneticAlgorithm - TODO - add comments
type GeneticAlgorithm struct {
	currentGeneration []*geneticCircuit
	maxIterations     int
	mutationRate      float64
	numParents        int
	numChildren       int
	numIterations     int
	random            *rand.Rand
}

type geneticCircuit struct {
	circuit []model.CircuitVertex
	length  float64
}

func (g *geneticCircuit) difference(other *geneticCircuit) float64 {
	difference := 0.0
	startIndex := 0
	for ; startIndex < len(other.circuit) && other.circuit[startIndex] != g.circuit[0]; startIndex++ {
	}

	for i, j := 0, startIndex; i < len(g.circuit); i, j = i+1, (j+1)%len(other.circuit) {
		difference += g.circuit[i].DistanceTo(other.circuit[j])
	}
	return difference
}

func (g *geneticCircuit) setLength() {
	g.length = model.Length(g.circuit)
}

func NewGeneticAlgorithm(initCircuit []model.CircuitVertex, numParents int, numChildren int, maxIterations int) model.Circuit {
	circuitLen := len(initCircuit)
	initGeneration := make([]*geneticCircuit, numParents)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create an initial generation of random parents.
	for genIndex := 0; genIndex < numParents; genIndex++ {
		current := &geneticCircuit{
			circuit: make([]model.CircuitVertex, circuitLen),
		}
		initGeneration[genIndex] = current
		// Insert each vertex at a random index in the circuit, if there is already a value at that index, use the next 'nil' index.
		for _, v := range initCircuit {
			circuitIndex := random.Intn(circuitLen)
			for ; current.circuit[circuitIndex] != nil; circuitIndex = (circuitIndex + 1) % circuitLen {
			}
			current.circuit[circuitIndex] = v
		}
		current.setLength()
	}

	g := &GeneticAlgorithm{
		currentGeneration: initGeneration,
		maxIterations:     maxIterations,
		mutationRate:      0.1,
		numParents:        numParents,
		numChildren:       numChildren,
		numIterations:     0,
		random:            random,
	}
	g.sortGeneration()
	return g
}

func NewGeneticAlgorithmWithPerimeterBuilder(initCircuit []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder, numParents int, numChildren int, maxIterations int) model.Circuit {
	initEdges, interiorVertices := perimeterBuilder(initCircuit)
	circuitLen := len(interiorVertices)
	initGeneration := make([]*geneticCircuit, numParents)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create an initial generation of random parents.
	for genIndex := 0; genIndex < numParents; genIndex++ {
		current := &geneticCircuit{
			circuit: make([]model.CircuitVertex, 0, circuitLen),
		}
		initGeneration[genIndex] = current
		// Load the convex perimeter into the parent first.
		for _, e := range initEdges {
			current.circuit = append(current.circuit, e.GetStart())
		}

		// Insert each interior vertex at a random location along the perimeter.
		for v := range interiorVertices {
			vertexIndex := random.Intn(len(current.circuit))
			current.circuit = model.InsertVertex(current.circuit, vertexIndex, v)
		}
		current.setLength()
	}

	g := &GeneticAlgorithm{
		currentGeneration: initGeneration,
		maxIterations:     maxIterations,
		mutationRate:      0.1,
		numParents:        numParents,
		numChildren:       numChildren,
		numIterations:     0,
		random:            random,
	}
	g.sortGeneration()
	return g
}

func (g *GeneticAlgorithm) FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge) {
	// If we have reached the number of iterations we are done, so return (nil,nil)
	if g.numIterations >= g.maxIterations {
		return nil, nil
	}
	// This does not update circuits one vertex at a time, so just return the first vertex in the best circuit since it will be ignored by Update().
	return g.currentGeneration[0].circuit[0], nil
}

func (g *GeneticAlgorithm) GetAttachedVertices() []model.CircuitVertex {
	return g.currentGeneration[0].circuit
}

func (g *GeneticAlgorithm) GetLength() float64 {
	return g.currentGeneration[0].length
}

func (g *GeneticAlgorithm) GetUnattachedVertices() map[model.CircuitVertex]bool {
	return make(map[model.CircuitVertex]bool)
}

// SetMutationRate updates the GeneticAlgorithm's mutation rate, which determines how frequently child circuits will be mutated (after cross-over).
// The mutation rate should be between 0.0 (0% chance of mutation) and 1.0 (100% chance of mutation).
// Numbers greater than 1.0 will behave as though they were 1.0, and numbers less than 0.0 will behave as though they were 0.0.
// If unspecified, the mutation rate defaults to 0.1 (10%)
func (g *GeneticAlgorithm) SetMutationRate(mutationRate float64) {
	g.mutationRate = mutationRate
}

func (g *GeneticAlgorithm) Update(vertexToAdd model.CircuitVertex, edgeToSplit model.CircuitEdge) {
	if g.numIterations >= g.maxIterations {
		return
	}
	g.numIterations++

	// For i..numChildren, create a child circuit:
	//   a. Select random pairs of parents for cross-breeding.
	//   b. Select the first parent at random. // TODO (optional): prefer those that are smaller circuits
	//   c. For second parent, prefer those that are dissimilar to the first parent. // TODO (optional): also prefer those that are smaller circuits
	//   d. Create the child circuit by combining the parents circuits via crossover + mutation
	nextGeneration := make([]*geneticCircuit, g.numChildren)
	for childIndex := 0; childIndex < g.numChildren; childIndex++ {
		parentA := g.currentGeneration[g.random.Intn(g.numParents)]

		parentDifferences := make([]float64, g.numParents)
		totalDifferences := 0.0
		for i := 0; i < g.numParents; i++ {
			parentDiff := parentA.difference(g.currentGeneration[i])
			parentDifferences[i] = parentDiff
			totalDifferences += parentDiff
		}

		parentB := g.currentGeneration[g.numParents-1]
		for i, selector := 0, g.random.Float64()*totalDifferences; i < g.numParents; i++ {
			selector = selector - parentDifferences[i]
			if selector <= 0.0 {
				parentB = g.currentGeneration[i]
				break
			}
		}

		// Generate at least one crossover point at random, and create the child from the parents.
		crossoverIndices := []int{}
		for numCrossovers := 1 + g.random.Intn(len(parentA.circuit)-2); numCrossovers > 0; numCrossovers-- {
			crossoverIndices = append(crossoverIndices, 1+g.random.Intn(len(parentA.circuit)-2))
		}
		sort.Ints(crossoverIndices)
		childCircuit := crossover(parentA, parentB, crossoverIndices)

		// Check the child for duplicates and missing vertices, and fix them.
		g.fixMissingAndDuplicateVertices(childCircuit, parentA.circuit)

		//Mutate the circuit
		g.mutate(childCircuit)

		child := &geneticCircuit{
			circuit: childCircuit,
		}
		child.setLength()
		nextGeneration[childIndex] = child
	}

	// Perform the genetic algorithm's selection step - by combining the children with the parent generation and eliminate excess potential parents from the next generation based on circuit length.
	// Note 1: we do this at the end, rather than the start, to minimize the amount of data we store between updates.
	// Note 2: we do not discard parents unless they are less "fit", so that we do not lose good approximations between generations.
	g.currentGeneration = append(g.currentGeneration, nextGeneration...)
	g.sortGeneration()
	g.currentGeneration = g.currentGeneration[0:g.numParents]
}

func (g *GeneticAlgorithm) fixMissingAndDuplicateVertices(toFix []model.CircuitVertex, allVertices []model.CircuitVertex) (fixed []model.CircuitVertex) {
	// Track all vertices in a map
	missingVertices := make(map[model.CircuitVertex]bool)
	for _, v := range allVertices {
		missingVertices[v] = true
	}
	duplicateIndices := []int{}

	// Remove vertices from the map as we encounter them in the toFix array.
	// If the vertex was already removed, add its index to the duplicateIndices array.
	for i, v := range toFix {
		if missingVertices[v] {
			delete(missingVertices, v)
		} else {
			duplicateIndices = append(duplicateIndices, i)
		}
	}

	// Add each missing vertex to the array in place of a random duplicate.
	for missingVertex := range missingVertices {
		duplicateIndex := g.random.Intn(len(duplicateIndices))
		vertexIndex := duplicateIndices[duplicateIndex]
		toFix[vertexIndex] = missingVertex

		// Remove the replaced duplicate from the list of duplicates, to avoid reusing it.
		duplicateIndices = model.DeleteIndexInt(duplicateIndices, duplicateIndex)
	}

	return toFix
}

// mutate swaps random vertices in the child array, according to the mutation rate.
func (g *GeneticAlgorithm) mutate(child []model.CircuitVertex) {
	numVertices := len(child)
	for i := 0; i < numVertices; i++ {
		if g.random.Float64() < g.mutationRate {
			swapIndex := g.random.Intn(numVertices)
			child[i], child[swapIndex] = child[swapIndex], child[i]
		}
	}
}

//sortGeneration orders the current generation from shortest length to longest.
func (g *GeneticAlgorithm) sortGeneration() {
	sort.Slice(g.currentGeneration, func(i, j int) bool {
		return g.currentGeneration[i].length < g.currentGeneration[j].length
	})
}

func crossover(parentA *geneticCircuit, parentB *geneticCircuit, crossoverIndices []int) (child []model.CircuitVertex) {
	// TODO (Optional): Use same start vertex for both parents circuits during crossover.
	// startIndexB := 0
	// for ; startIndexB < len(parentB.circuit) && parentB.circuit[startIndexB] != parentA.circuit[0]; startIndexB++ {
	// }

	child = make([]model.CircuitVertex, 0, len(parentA.circuit))

	// Append everything prior to the first crossover, from parentA
	child = append(child, parentA.circuit[:crossoverIndices[0]]...)

	// Set the active parent to parentB, so that the next append uses parentB, even when there is only one crossover.
	activeParent := parentB
	lastCrossoverIndex := len(crossoverIndices) - 1
	for i, next := 0, 1; i < lastCrossoverIndex; i, next = i+1, next+1 {
		crossoverStart := crossoverIndices[i]
		crossoverEnd := crossoverIndices[next]
		// Ignore zero-length crossovers.
		if crossoverEnd == crossoverStart {
			continue
		}
		child = append(child, activeParent.circuit[crossoverStart:crossoverEnd]...)
		// Swap the active parent between each crossover, this needs to be after append so that the final append works correctly.
		if activeParent == parentA {
			activeParent = parentB
		} else {
			activeParent = parentA
		}
	}
	// Append everything after the last crossover, from the active parent.
	child = append(child, activeParent.circuit[crossoverIndices[lastCrossoverIndex]:]...)

	return child
}
