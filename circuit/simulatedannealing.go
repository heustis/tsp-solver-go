package circuit

import (
	"math"
	"math/rand"
	"time"

	"github.com/fealos/lee-tsp-go/model"
)

// SimulatedAnnealing - TODO
type SimulatedAnnealing struct {
	circuit              []model.CircuitVertex
	farthestDistance     float64
	maxIterations        float64
	numIterations        float64
	preferCloseNeighbors bool
	random               *rand.Rand
	temperatureFunction  func(currentIteration float64, maxIterations float64) float64
}

func NewSimulatedAnnealing(circuit []model.CircuitVertex, maxIterations int, preferCloseNeighbors bool) model.Circuit {
	return &SimulatedAnnealing{
		circuit:              circuit,
		farthestDistance:     computeFarthestDistance(circuit),
		maxIterations:        float64(maxIterations),
		numIterations:        0.0,
		preferCloseNeighbors: preferCloseNeighbors,
		random:               rand.New(rand.NewSource(time.Now().UnixNano())),
		temperatureFunction:  CalculateTemperatureLinear,
	}
}

func NewSimulatedAnnealingFromCircuit(circuit model.Circuit, maxIterations int, preferCloseNeighbors bool) model.Circuit {
	for nextVertex, nextEdge := circuit.FindNextVertexAndEdge(); nextVertex != nil; nextVertex, nextEdge = circuit.FindNextVertexAndEdge() {
		circuit.Update(nextVertex, nextEdge)
	}

	initCircuit := circuit.GetAttachedVertices()
	return &SimulatedAnnealing{
		circuit:              initCircuit,
		farthestDistance:     computeFarthestDistance(initCircuit),
		maxIterations:        float64(maxIterations),
		numIterations:        0.0,
		preferCloseNeighbors: preferCloseNeighbors,
		random:               rand.New(rand.NewSource(time.Now().UnixNano())),
		temperatureFunction:  CalculateTemperatureLinear,
	}
}

func (s *SimulatedAnnealing) FindNextVertexAndEdge() (model.CircuitVertex, model.CircuitEdge) {
	// If we have reached the number of iterations we are done, so return (nil,nil)
	if s.numIterations >= s.maxIterations {
		return nil, nil
	}
	// We will determine the next vertex in Update(), so just return the first vertex in the circuit since it will be ignored by Update().
	return s.circuit[0], nil
}

func (s *SimulatedAnnealing) GetAttachedVertices() []model.CircuitVertex {
	return s.circuit
}

func (s *SimulatedAnnealing) GetLength() float64 {
	return model.Length(s.circuit)
}

func (s *SimulatedAnnealing) GetUnattachedVertices() map[model.CircuitVertex]bool {
	return make(map[model.CircuitVertex]bool)
}

// SetSeed sets the seed used by the SimulatedAnnealing for random number generation.
// This is to facilitate consistent unit tests.
func (s *SimulatedAnnealing) SetSeed(seed int64) {
	s.random = rand.New(rand.NewSource(seed))
}

// SetTemperatureFunction updates the function used in each iteration of Update() to calculate the temperature.
// By default SimulatedAnnealing uses a linear temperature function, but this package also provides a geometric temperature function, and enables custom temperature functions.
func (s *SimulatedAnnealing) SetTemperatureFunction(temperatureFunction func(currentIteration float64, maxIterations float64) float64) {
	s.temperatureFunction = temperatureFunction
}

func (s *SimulatedAnnealing) Update(vertexToAdd model.CircuitVertex, edgeToSplit model.CircuitEdge) {
	if s.numIterations >= s.maxIterations {
		return
	}

	s.numIterations++

	// This section could be included in FindNextVertexAndEdge, but it is more performant to have the indices here.
	// Select two random vertices to check if they should be swapped.
	numVertices := len(s.circuit)
	indexA := s.random.Intn(numVertices)
	indexAPrev := (indexA + numVertices - 1) % numVertices
	indexANext := (indexA + 1) % numVertices

	var indexB int
	if s.preferCloseNeighbors {
		indexB = s.getRandomNeighbor(indexA)
	} else {
		// Select a random vertex for B, but don't allow it to be the same as A.
		indexB = s.random.Intn(numVertices)
		for indexA == indexB {
			indexB = s.random.Intn(numVertices)
		}
	}
	indexBPrev := (indexB + numVertices - 1) % numVertices
	indexBNext := (indexB + 1) % numVertices

	// Calculate the effect swapping the two vertices will have on the length of the circuit.
	lengthACurrent := s.circuit[indexAPrev].DistanceTo(s.circuit[indexA]) + s.circuit[indexA].DistanceTo(s.circuit[indexANext])
	lengthANew := s.circuit[indexAPrev].DistanceTo(s.circuit[indexB]) + s.circuit[indexB].DistanceTo(s.circuit[indexANext])

	lengthBCurrent := s.circuit[indexBPrev].DistanceTo(s.circuit[indexB]) + s.circuit[indexB].DistanceTo(s.circuit[indexBNext])
	lengthBNew := s.circuit[indexBPrev].DistanceTo(s.circuit[indexA]) + s.circuit[indexA].DistanceTo(s.circuit[indexBNext])

	// If the two vertices are adjacent, need to add the length of the edge A->B to each new length.
	if indexA == indexBPrev || indexA == indexBNext {
		distAToB := s.circuit[indexA].DistanceTo(s.circuit[indexB])
		lengthANew += distAToB
		lengthBNew += distAToB
	}

	edgeADelta := lengthANew - lengthACurrent
	edgeBDelta := lengthBNew - lengthBCurrent

	// Scale delta so that it has a meaningful value in the acceptance function, since cooridinates from -100 to +100 will produce different deltas than coordinates from -10000 to +10000.
	// The temperature is always between 0 and 1, decreasing from near 1 to near 0 as annealing progresses.
	// The delta could be limited between 0 and 1 as well, so that all posibilities are feasable at a temperature of 1.
	// However, we know that any intersecting edges are not optimal, so we can optimize this by allowing the delta to exceed 1 in bad use cases.
	// The worst case delta is is if B and A are the farthest vertices from each other and both go from their closest vertices to their farthest vertices, and the best case is the reverse.
	// This worst case is guaranteed to be less than 4*|B-A|, but we will use |B-A| since it is okay if we ignore the possibilities that are close to the worst case.
	deltaIncrease := (edgeADelta + edgeBDelta) / s.farthestDistance

	temperature := s.temperatureFunction(s.numIterations, s.maxIterations)

	// Swap the two vertices if it would decrease the size of the circuit, or if the increase is within the acceptable bounds defined by the acceptance function.
	if testValue, acceptanceThreshold := s.random.Float64(), math.Exp(-deltaIncrease/temperature); deltaIncrease <= 0.0 || testValue < acceptanceThreshold {
		s.circuit[indexA], s.circuit[indexB] = s.circuit[indexB], s.circuit[indexA]
	}
}

// getRandomNeighbor weighs vertices based on their distance from the vertex at the supplied index, then randomly selects a vertex based on the weights.
func (s *SimulatedAnnealing) getRandomNeighbor(index int) (neighborIndex int) {
	// len-1 to ignore the vertex at the supplied index.
	weights := make([]*weightedVertex, len(s.circuit)-1)
	totalWeight := 0.0
	for i, weightIndex := 0, 0; i < len(s.circuit); i++ {
		if i != index {
			// Invert the distance between the points, so that closer points have larger weights than farther points (e.g. 1/5 > 1/500).
			weight := 1.0 / s.circuit[index].DistanceTo(s.circuit[i])
			weights[weightIndex] = &weightedVertex{
				weight:      weight,
				vertexIndex: i,
			}
			totalWeight += weight
			weightIndex++
		}
	}

	// Select a random index by weight:
	// 1) Select a random value between [0,totalWeight)
	// 2) Iterate through the weighted values, subtracting their weight from the random weight
	// 3) Once the random weight has a value of 0 or less, select the vertex with the weight that caused it to transition to 0 or negative.
	randomWeight := s.random.Float64() * totalWeight
	for _, w := range weights {
		randomWeight -= w.weight
		if randomWeight <= 0 {
			return w.vertexIndex
		}
	}
	// This should never be reached, since the random weight should never be greater than the total weight.
	return weights[len(weights)-1].vertexIndex
}

type weightedVertex struct {
	vertexIndex int
	weight      float64
}

// CalculateTemperatureGeometric calculates temperature according to the equation t'=t*X, so that it decreases geometricaly as the model iterates.
// For this implementation we are using 5.0 since it produces selective acceptance without making the later iterations useless:
// * 1.0 -> 0.99^99=0.3697, 0.999^999=0.3681,
// * 5.0 -> 0.95^99=0.0062, 0.995^999=0.0066,
// * 10.0 -> 0.90^99=0.000059, 0.990^999=0.000044,
func CalculateTemperatureGeometric(currentIteration float64, maxIterations float64) float64 {
	return math.Pow(1.0-5.0/maxIterations, currentIteration)
}

// CalculateTemperatureLinear calculates temperature according to the function t'=t-X, so that it decreases linearly as the model iterates.
func CalculateTemperatureLinear(currentIteration float64, maxIterations float64) float64 {
	return 1.0 - currentIteration/maxIterations
}

func computeFarthestDistance(circuit []model.CircuitVertex) float64 {
	farthestDistance := 0.0
	// Find the distance between the two farthest vertices, for scaling the delta
	for _, v := range circuit {
		farthestFromV := model.FindFarthestPoint(v, circuit)
		if testDistance := farthestFromV.DistanceTo(v); testDistance > farthestDistance {
			farthestDistance = testDistance
		}
	}
	return farthestDistance
}

var _ model.Circuit = (*SimulatedAnnealing)(nil)
