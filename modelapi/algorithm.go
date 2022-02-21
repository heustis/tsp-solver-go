package modelapi

import (
	"math"

	"github.com/heustis/tsp-solver-go/circuit"
	"github.com/heustis/tsp-solver-go/model"
)

type AlgorithmType string

const (
	ALG_ANNEALING        AlgorithmType = "ANNEALING"
	ALG_CLOSEST_CLONE    AlgorithmType = "CLOSEST_CLONE"
	ALG_CLOSEST_GREEDY   AlgorithmType = "CLOSEST_GREEDY"
	ALG_DISPARITY_CLONE  AlgorithmType = "DISPARITY_CLONE"
	ALG_DISPARITY_GREEDY AlgorithmType = "DISPARITY_GREEDY"
	ALG_GENETIC          AlgorithmType = "GENETIC"
)

type TemperatureFunctionType string

const (
	TEMP_DEFAULT   TemperatureFunctionType = ""
	TEMP_GEOMETRIC TemperatureFunctionType = "GEOMETRIC"
	TEMP_LINEAR    TemperatureFunctionType = "LINEAR"
)

// Algorithm represents a union of the possible configuration data used by different types of circuits, so that the API can appear to be polymorphic.
type Algorithm struct {
	AlgorithmType         AlgorithmType           `json:"algorithmType" validate:"required,oneof=ANNEALING CLOSEST_CLONE CLOSEST_GREEDY DISPARITY_CLONE DISPARITY_GREEDY GENETIC"`
	CloneByInitEdges      *bool                   `json:"cloneByInitEdges,omitempty"`
	CloneOnFirstAttach    *bool                   `json:"cloneOnFirstAttach,omitempty"`
	MaxClones             *int64                  `json:"maxClones,omitempty"`
	MaxCrossovers         int                     `json:"maxCrossovers,omitempty" validate:"isdefault|min=1"`
	MaxIterations         int                     `json:"maxIterations,omitempty" validate:"isdefault|min=1,required_if=AlgorithmType GENETIC,required_if=AlgorithmType ANNEALING"`
	MinSignificance       *float64                `json:"minSignificance,omitempty" validate:"omitempty,min=0"`
	MutationRate          *float64                `json:"mutationRate,omitempty" validate:"omitempty,min=0,max=1"`
	NumChildren           int                     `json:"numChildren,omitempty" validate:"isdefault|min=1,required_if=AlgorithmType GENETIC"`
	NumParents            int                     `json:"numParents,omitempty" validate:"isdefault|min=1,required_if=AlgorithmType GENETIC"`
	PrecursorAlgorithm    *Algorithm              `json:"precursorAlgorithm,omitempty" validate:"omitempty,dive"`
	PreferCloseNeighbors  *bool                   `json:"preferCloseNeighbors,omitempty"`
	Seed                  *int64                  `json:"seed,omitempty"`
	ShouldBuildConvexHull *bool                   `json:"shouldBuildConvexHull,omitempty"`
	TemperatureFunction   TemperatureFunctionType `json:"temperatureFunction,omitempty" validate:"omitempty,oneof=GEOMETRIC LINEAR"`
	UpdateInteriorPoints  *bool                   `json:"updateInteriorPoints,omitempty"`
	UseRelativeDisparity  *bool                   `json:"useRelativeDisparity,omitempty"`
}

func (alg *Algorithm) GetCircuitFunction() func(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) model.Circuit {
	switch alg.AlgorithmType {
	case ALG_ANNEALING:
		return alg.CreateSimulatedAnnealing
	case ALG_CLOSEST_CLONE:
		return alg.CreateClosestClone
	case ALG_DISPARITY_CLONE:
		return alg.CreateDisparityClone
	case ALG_DISPARITY_GREEDY:
		return alg.CreateDisparityGreedy
	case ALG_GENETIC:
		return alg.CreateGenetic
	default:
		return alg.CreateClosestGreedy
	}
}

func (alg *Algorithm) CreateClosestClone(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) model.Circuit {
	c := circuit.NewClonableCircuitImpl(vertices, perimeterBuilder)
	c.CloneOnFirstAttach = isTrue(alg.CloneOnFirstAttach)
	solver := circuit.NewClonableCircuitSolver(c)
	if alg.MaxClones != nil {
		solver.MaxClones = int(*alg.MaxClones)
	}
	return solver
}

func (alg *Algorithm) CreateClosestGreedy(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) model.Circuit {
	if isTrue(alg.CloneByInitEdges) {
		return circuit.NewConvexConcaveByEdge(vertices, perimeterBuilder, isTrue(alg.UpdateInteriorPoints))
	} else {
		return circuit.NewConvexConcave(vertices, perimeterBuilder, isTrue(alg.UpdateInteriorPoints))
	}
}

func (alg *Algorithm) CreateDisparityClone(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) model.Circuit {
	c := circuit.NewConvexConcaveConfidence(vertices, perimeterBuilder)
	if alg.MaxClones != nil {
		if *alg.MaxClones < 1 || *alg.MaxClones > math.MaxInt16 {
			c.MaxClones = math.MaxUint16
		} else {
			c.MaxClones = uint16(*alg.MaxClones)
		}
	}
	if alg.MinSignificance != nil {
		c.Significance = *alg.MinSignificance
	}
	return c
}

func (alg *Algorithm) CreateDisparityGreedy(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) model.Circuit {
	return circuit.NewConvexConcaveDisparity(vertices, perimeterBuilder, isTrue(alg.UseRelativeDisparity))
}

func (alg *Algorithm) CreateGenetic(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) model.Circuit {
	var c *circuit.GeneticAlgorithm
	if isTrue(alg.ShouldBuildConvexHull) {
		c = circuit.NewGeneticAlgorithmWithPerimeterBuilder(vertices, perimeterBuilder, alg.NumParents, alg.NumChildren, alg.MaxIterations)
	} else {
		c = circuit.NewGeneticAlgorithm(vertices, alg.NumParents, alg.NumChildren, alg.MaxIterations)
	}
	if alg.MaxCrossovers > 0 {
		c.SetMaxCrossovers(alg.MaxCrossovers)
	}
	if alg.MutationRate != nil {
		c.SetMutationRate(*alg.MutationRate)
	}
	if alg.Seed != nil {
		c.SetSeed(*alg.Seed)
	}
	return c
}

func (alg *Algorithm) CreateSimulatedAnnealing(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) model.Circuit {
	var c *circuit.SimulatedAnnealing
	if alg.PrecursorAlgorithm != nil {
		precursorCircuit := alg.PrecursorAlgorithm.GetCircuitFunction()(vertices, perimeterBuilder)
		c = circuit.NewSimulatedAnnealingFromCircuit(precursorCircuit, alg.MaxIterations, isTrue(alg.PreferCloseNeighbors))
	} else {
		c = circuit.NewSimulatedAnnealing(vertices, alg.MaxIterations, isTrue(alg.PreferCloseNeighbors))
	}
	if alg.Seed != nil {
		c.SetSeed(*alg.Seed)
	}
	// The default temperature function is linear, so don't need to update it unless it is different.
	if alg.TemperatureFunction == TEMP_GEOMETRIC {
		c.SetTemperatureFunction(circuit.CalculateTemperatureGeometric)
	}
	return c
}

func isTrue(b *bool) bool {
	return b != nil && *b
}
