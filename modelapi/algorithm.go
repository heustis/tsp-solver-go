package modelapi

import (
	"math"

	"github.com/heustis/lee-tsp-go/circuit"
	"github.com/heustis/lee-tsp-go/model"
)

type AlgorithmType string

const (
	ALG_ANNEALING      AlgorithmType = "ANNEALING"
	ALG_CLONABLE       AlgorithmType = "CLONABLE"
	ALG_CONCAVE_CONVEX AlgorithmType = "CONCAVE_CONVEX"
	ALG_CONFIDENCE     AlgorithmType = "CONFIDENCE"
	ALG_DISPARITY      AlgorithmType = "DISPARITY"
	ALG_GENETIC        AlgorithmType = "GENETIC"
)

type TemperatureFunctionType string

const (
	TEMP_DEFAULT   TemperatureFunctionType = ""
	TEMP_GEOMETRIC TemperatureFunctionType = "GEOMETRIC"
	TEMP_LINEAR    TemperatureFunctionType = "LINEAR"
)

// Algorithm represents a union of the possible configuration data used by different types of circuits, so that the API can appear to be polymorphic.
type Algorithm struct {
	AlgorithmType         AlgorithmType           `json:"algorithmType" validate:"required,oneof=ANNEALING CLONABLE CONCAVE_CONVEX CONFIDENCE DISPARITY GENETIC"`
	CloneByInitEdges      *bool                   `json:"cloneByInitEdges,omitempty"`
	CloneOnFirstAttach    *bool                   `json:"cloneOnFirstAttach,omitempty"`
	MaxClones             *int64                  `json:"maxClones,omitempty"`
	MaxIterations         int                     `json:"maxIterations,omitempty" validate:"isdefault|min=1,required_if=AlgorithmType GENETIC,required_if=AlgorithmType ANNEALING"`
	MinSignificance       *float64                `json:"minSignificance,omitempty" validate:"omitempty,min=0"`
	MutationRate          *float64                `json:"mutationRate,omitempty" validate:"omitempty,min=0,max=1"`
	NumChildren           int                     `json:"numChildren,omitempty" validate:"isdefault|min=1,required_if=AlgorithmType GENETIC"`
	NumParents            int                     `json:"numParents,omitempty" validate:"isdefault|min=1,required_if=AlgorithmType GENETIC"`
	PrecursorAlgorithm    *Algorithm              `json:"PrecursorAlgorithm,omitempty" validate:"omitempty,dive"`
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
	case ALG_CLONABLE:
		return alg.CreateClonable
	case ALG_CONFIDENCE:
		return alg.CreateConfidence
	case ALG_DISPARITY:
		return alg.CreateDisparity
	case ALG_GENETIC:
		return alg.CreateGenetic
	default:
		return alg.CreateConcaveConvex
	}
}

func (alg *Algorithm) CreateClonable(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) model.Circuit {
	c := circuit.NewClonableCircuitImpl(vertices, perimeterBuilder)
	if alg.CloneOnFirstAttach != nil {
		c.CloneOnFirstAttach = *alg.CloneOnFirstAttach
	}
	solver := circuit.NewClonableCircuitSolver(c)
	if alg.MaxClones != nil {
		solver.MaxClones = int(*alg.MaxClones)
	}
	return solver
}

func (alg *Algorithm) CreateConcaveConvex(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) model.Circuit {
	if getBool(alg.CloneByInitEdges) {
		return circuit.NewConvexConcaveByEdge(vertices, perimeterBuilder, getBool(alg.UpdateInteriorPoints))
	} else {
		return circuit.NewConvexConcave(vertices, perimeterBuilder, getBool(alg.UpdateInteriorPoints))
	}
}

func (alg *Algorithm) CreateConfidence(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) model.Circuit {
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

func (alg *Algorithm) CreateDisparity(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) model.Circuit {
	return circuit.NewConvexConcaveDisparity(vertices, perimeterBuilder, getBool(alg.UseRelativeDisparity))
}

func (alg *Algorithm) CreateGenetic(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) model.Circuit {
	var c *circuit.GeneticAlgorithm
	if getBool(alg.ShouldBuildConvexHull) {
		c = circuit.NewGeneticAlgorithmWithPerimeterBuilder(vertices, perimeterBuilder, alg.NumParents, alg.NumChildren, alg.MaxIterations)
	} else {
		c = circuit.NewGeneticAlgorithm(vertices, alg.NumParents, alg.NumChildren, alg.MaxIterations)
	}
	if alg.MutationRate != nil {
		c.SetMutationRate(*alg.MutationRate)
	}
	return c
}

func (alg *Algorithm) CreateSimulatedAnnealing(vertices []model.CircuitVertex, perimeterBuilder model.PerimeterBuilder) model.Circuit {
	var c *circuit.SimulatedAnnealing
	if alg.PrecursorAlgorithm != nil {
		precursorCircuit := alg.PrecursorAlgorithm.GetCircuitFunction()(vertices, perimeterBuilder)
		c = circuit.NewSimulatedAnnealingFromCircuit(precursorCircuit, alg.MaxIterations, getBool(alg.PreferCloseNeighbors))
	} else {
		c = circuit.NewSimulatedAnnealing(vertices, alg.MaxIterations, getBool(alg.PreferCloseNeighbors))
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

func getBool(b *bool) bool {
	return b != nil && *b
}
