package modelapi_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/heustis/tsp-solver-go/circuit"
	"github.com/heustis/tsp-solver-go/model2d"
	"github.com/heustis/tsp-solver-go/modelapi"
	"github.com/heustis/tsp-solver-go/solver"
	"github.com/stretchr/testify/assert"
)

func TestValidateAlgorithm(t *testing.T) {
	assert := assert.New(t)
	validate := validator.New()
	assert.EqualError(validate.Struct(modelapi.Algorithm{}), "Key: 'Algorithm.AlgorithmType' Error:Field validation for 'AlgorithmType' failed on the 'required' tag")
	assert.EqualError(validate.Struct(modelapi.Algorithm{AlgorithmType: ""}), "Key: 'Algorithm.AlgorithmType' Error:Field validation for 'AlgorithmType' failed on the 'required' tag")
	assert.EqualError(validate.Struct(modelapi.Algorithm{AlgorithmType: "TEST"}), "Key: 'Algorithm.AlgorithmType' Error:Field validation for 'AlgorithmType' failed on the 'oneof' tag")

	assert.EqualError(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_ANNEALING}), "Key: 'Algorithm.MaxIterations' Error:Field validation for 'MaxIterations' failed on the 'required_if' tag")
	assert.EqualError(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_ANNEALING, MaxIterations: -5}), "Key: 'Algorithm.MaxIterations' Error:Field validation for 'MaxIterations' failed on the 'isdefault|min=1' tag")
	assert.EqualError(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_ANNEALING, MaxIterations: 5000000, TemperatureFunction: "OTHER"}), "Key: 'Algorithm.TemperatureFunction' Error:Field validation for 'TemperatureFunction' failed on the 'oneof' tag")
	assert.EqualError(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_ANNEALING, MaxIterations: 5000000, PrecursorAlgorithm: &modelapi.Algorithm{AlgorithmType: "UNKOWN"}}), "Key: 'Algorithm.PrecursorAlgorithm.AlgorithmType' Error:Field validation for 'AlgorithmType' failed on the 'oneof' tag")
	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_ANNEALING, MaxIterations: 5000000}))
	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_ANNEALING, MaxIterations: 5000000, TemperatureFunction: "LINEAR"}))
	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_ANNEALING, MaxIterations: 5000000, TemperatureFunction: "GEOMETRIC"}))

	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_CLOSEST_CLONE}))
	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_CLOSEST_CLONE, MaxClones: intPointer(15)}))
	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_CLOSEST_CLONE, CloneOnFirstAttach: boolPointer(false)}))
	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_CLOSEST_CLONE, MaxClones: intPointer(-1), CloneOnFirstAttach: boolPointer(false)}))
	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_CLOSEST_CLONE, MaxClones: intPointer(45), CloneOnFirstAttach: boolPointer(true)}))

	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_CLOSEST_GREEDY}))
	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_CLOSEST_GREEDY, CloneByInitEdges: boolPointer(false)}))
	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_CLOSEST_GREEDY, CloneByInitEdges: boolPointer(true)}))

	assert.EqualError(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_DISPARITY_CLONE, MinSignificance: float64Pointer(-.5)}), "Key: 'Algorithm.MinSignificance' Error:Field validation for 'MinSignificance' failed on the 'min' tag")
	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_DISPARITY_CLONE}))
	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_DISPARITY_CLONE, MinSignificance: float64Pointer(.25)}))

	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_DISPARITY_GREEDY}))
	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_DISPARITY_GREEDY, UseRelativeDisparity: boolPointer(false)}))
	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_DISPARITY_GREEDY, UseRelativeDisparity: boolPointer(true)}))

	assert.EqualError(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_GENETIC}),
		"Key: 'Algorithm.MaxIterations' Error:Field validation for 'MaxIterations' failed on the 'required_if' tag\n"+
			"Key: 'Algorithm.NumChildren' Error:Field validation for 'NumChildren' failed on the 'required_if' tag\n"+
			"Key: 'Algorithm.NumParents' Error:Field validation for 'NumParents' failed on the 'required_if' tag")

	assert.EqualError(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_GENETIC, MaxIterations: -1}),
		"Key: 'Algorithm.MaxIterations' Error:Field validation for 'MaxIterations' failed on the 'isdefault|min=1' tag\n"+
			"Key: 'Algorithm.NumChildren' Error:Field validation for 'NumChildren' failed on the 'required_if' tag\n"+
			"Key: 'Algorithm.NumParents' Error:Field validation for 'NumParents' failed on the 'required_if' tag")

	assert.EqualError(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_GENETIC, MaxIterations: 10, NumChildren: -5}),
		"Key: 'Algorithm.NumChildren' Error:Field validation for 'NumChildren' failed on the 'isdefault|min=1' tag\n"+
			"Key: 'Algorithm.NumParents' Error:Field validation for 'NumParents' failed on the 'required_if' tag")

	assert.EqualError(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_GENETIC, MaxIterations: 10, NumChildren: 5, NumParents: -15}),
		"Key: 'Algorithm.NumParents' Error:Field validation for 'NumParents' failed on the 'isdefault|min=1' tag")

	assert.EqualError(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_GENETIC, MaxIterations: 10, NumChildren: 5, NumParents: 15, MutationRate: float64Pointer(-.5)}),
		"Key: 'Algorithm.MutationRate' Error:Field validation for 'MutationRate' failed on the 'min' tag")

	assert.EqualError(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_GENETIC, MaxIterations: 10, NumChildren: 5, NumParents: 15, MutationRate: float64Pointer(5)}),
		"Key: 'Algorithm.MutationRate' Error:Field validation for 'MutationRate' failed on the 'max' tag")

	assert.EqualError(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_GENETIC, MaxIterations: 10, NumChildren: 5, NumParents: 15, MaxCrossovers: -3}),
		`Key: 'Algorithm.MaxCrossovers' Error:Field validation for 'MaxCrossovers' failed on the 'isdefault|min=1' tag`)

	assert.Nil(validate.Struct(modelapi.Algorithm{AlgorithmType: modelapi.ALG_GENETIC, MaxIterations: 10, NumChildren: 5, NumParents: 15, Seed: intPointer(12345), MaxCrossovers: 6}))
}

func TestGetProcessFunction(t *testing.T) {
	assert := assert.New(t)

	alg := &modelapi.Algorithm{}
	alg.AlgorithmType = modelapi.ALG_ANNEALING
	assert.True(reflect.ValueOf(alg.CreateSimulatedAnnealing).Pointer() == reflect.ValueOf(alg.GetCircuitFunction()).Pointer())

	alg.AlgorithmType = modelapi.ALG_CLOSEST_CLONE
	assert.True(reflect.ValueOf(alg.CreateClosestClone).Pointer() == reflect.ValueOf(alg.GetCircuitFunction()).Pointer())

	alg.AlgorithmType = modelapi.ALG_CLOSEST_GREEDY
	assert.True(reflect.ValueOf(alg.CreateClosestGreedy).Pointer() == reflect.ValueOf(alg.GetCircuitFunction()).Pointer())

	alg.AlgorithmType = modelapi.ALG_DISPARITY_CLONE
	assert.True(reflect.ValueOf(alg.CreateDisparityClone).Pointer() == reflect.ValueOf(alg.GetCircuitFunction()).Pointer())

	alg.AlgorithmType = modelapi.ALG_DISPARITY_GREEDY
	assert.True(reflect.ValueOf(alg.CreateDisparityGreedy).Pointer() == reflect.ValueOf(alg.GetCircuitFunction()).Pointer())

	alg.AlgorithmType = modelapi.ALG_GENETIC
	assert.True(reflect.ValueOf(alg.CreateGenetic).Pointer() == reflect.ValueOf(alg.GetCircuitFunction()).Pointer())
}

func TestCreateClosestClone(t *testing.T) {
	assert := assert.New(t)

	vertices := model2d.GenerateVertices(10)

	alg := &modelapi.Algorithm{AlgorithmType: modelapi.ALG_CLOSEST_CLONE}
	c := alg.CreateClosestClone(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.ClonableCircuitSolver{}, c)
	assert.Equal(-1, c.(*circuit.ClonableCircuitSolver).MaxClones)

	alg.MaxClones = intPointer(5)
	c = alg.CreateClosestClone(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.ClonableCircuitSolver{}, c)
	assert.Equal(5, c.(*circuit.ClonableCircuitSolver).MaxClones)

	alg.MaxClones = nil
	alg.CloneOnFirstAttach = boolPointer(false)
	c = alg.CreateClosestClone(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.ClonableCircuitSolver{}, c)
	assert.Equal(-1, c.(*circuit.ClonableCircuitSolver).MaxClones)

	alg.MaxClones = intPointer(10)
	c = alg.CreateClosestClone(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.ClonableCircuitSolver{}, c)
	assert.Equal(10, c.(*circuit.ClonableCircuitSolver).MaxClones)

	alg.MaxClones = nil
	alg.CloneOnFirstAttach = boolPointer(true)
	c = alg.CreateClosestClone(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.ClonableCircuitSolver{}, c)
	assert.Equal(-1, c.(*circuit.ClonableCircuitSolver).MaxClones)

	alg.MaxClones = intPointer(15)
	c = alg.CreateClosestClone(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.ClonableCircuitSolver{}, c)
	assert.Equal(15, c.(*circuit.ClonableCircuitSolver).MaxClones)
}

func TestCreateClosestGreedy(t *testing.T) {
	assert := assert.New(t)

	vertices := model2d.GenerateVertices(10)

	alg := &modelapi.Algorithm{AlgorithmType: modelapi.ALG_CLOSEST_GREEDY}
	assert.IsType(&circuit.ConvexConcave{}, alg.CreateClosestGreedy(vertices, model2d.BuildPerimiter))

	alg.CloneByInitEdges = boolPointer(false)
	assert.IsType(&circuit.ConvexConcave{}, alg.CreateClosestGreedy(vertices, model2d.BuildPerimiter))

	alg.CloneByInitEdges = boolPointer(true)
	assert.IsType(&circuit.ConvexConcaveByEdge{}, alg.CreateClosestGreedy(vertices, model2d.BuildPerimiter))
}

func TestCreateDisparityClone(t *testing.T) {
	assert := assert.New(t)

	vertices := model2d.GenerateVertices(10)

	alg := &modelapi.Algorithm{AlgorithmType: modelapi.ALG_DISPARITY_CLONE}
	c := alg.CreateDisparityClone(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.ConvexConcaveConfidence{}, c)
	assert.Equal(uint16(1000), c.(*circuit.ConvexConcaveConfidence).MaxClones)
	assert.Equal(1.0, c.(*circuit.ConvexConcaveConfidence).Significance)
	assert.Equal(vertices, c.(*circuit.ConvexConcaveConfidence).Vertices)

	for _, numClones := range []int64{10, 3, 21, 4400, math.MaxUint16, 15} {
		alg.MaxClones = intPointer(numClones)
		c = alg.CreateDisparityClone(vertices, model2d.BuildPerimiter)
		assert.IsType(&circuit.ConvexConcaveConfidence{}, c)
		assert.Equal(uint16(numClones), c.(*circuit.ConvexConcaveConfidence).MaxClones)
		assert.Equal(1.0, c.(*circuit.ConvexConcaveConfidence).Significance)
		assert.Equal(vertices, c.(*circuit.ConvexConcaveConfidence).Vertices)
	}

	for _, significance := range []float64{-1, 0.0, 12345.6789, 2.5} {
		alg.MinSignificance = float64Pointer(significance)
		c = alg.CreateDisparityClone(vertices, model2d.BuildPerimiter)
		assert.IsType(&circuit.ConvexConcaveConfidence{}, c)
		assert.Equal(uint16(15), c.(*circuit.ConvexConcaveConfidence).MaxClones)
		assert.Equal(significance, c.(*circuit.ConvexConcaveConfidence).Significance)
		assert.Equal(vertices, c.(*circuit.ConvexConcaveConfidence).Vertices)
	}

	for _, numClones := range []int64{1234567890, 0, -1, math.MinInt64, math.MaxInt64} {
		alg.MaxClones = intPointer(numClones)
		c = alg.CreateDisparityClone(vertices, model2d.BuildPerimiter)
		assert.IsType(&circuit.ConvexConcaveConfidence{}, c)
		assert.Equal(uint16(math.MaxUint16), c.(*circuit.ConvexConcaveConfidence).MaxClones)
		assert.Equal(2.5, c.(*circuit.ConvexConcaveConfidence).Significance)
		assert.Equal(vertices, c.(*circuit.ConvexConcaveConfidence).Vertices)
	}
}

func TestCreateDisparityGreedy(t *testing.T) {
	assert := assert.New(t)

	vertices := model2d.GenerateVertices(100)

	alg := &modelapi.Algorithm{AlgorithmType: modelapi.ALG_DISPARITY_GREEDY}
	c := alg.CreateDisparityGreedy(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.ConvexConcaveDisparity{}, c)

	alg.UseRelativeDisparity = boolPointer(false)
	c1 := alg.CreateDisparityGreedy(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.ConvexConcaveDisparity{}, c1)

	alg.UseRelativeDisparity = boolPointer(true)
	c2 := alg.CreateDisparityGreedy(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.ConvexConcaveDisparity{}, c2)

	solver.FindShortestPathCircuit(c)
	solver.FindShortestPathCircuit(c1)
	solver.FindShortestPathCircuit(c2)
	assert.Equal(c.GetAttachedVertices(), c1.GetAttachedVertices())
	assert.NotEqual(c1.GetAttachedVertices(), c2.GetAttachedVertices())
}

func TestCreateGenetic(t *testing.T) {
	assert := assert.New(t)

	vertices := model2d.GenerateVertices(10)

	alg := &modelapi.Algorithm{AlgorithmType: modelapi.ALG_GENETIC, NumChildren: 15, NumParents: 30, MaxIterations: 100}
	c := alg.CreateGenetic(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.GeneticAlgorithm{}, c)

	alg.ShouldBuildConvexHull = boolPointer(false)
	c = alg.CreateGenetic(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.GeneticAlgorithm{}, c)

	alg.ShouldBuildConvexHull = boolPointer(true)
	c = alg.CreateGenetic(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.GeneticAlgorithm{}, c)

	alg.MutationRate = float64Pointer(0.1234)
	c = alg.CreateGenetic(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.GeneticAlgorithm{}, c)

	alg.Seed = intPointer(7890)
	c = alg.CreateGenetic(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.GeneticAlgorithm{}, c)

	alg.MaxCrossovers = 5
	c = alg.CreateGenetic(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.GeneticAlgorithm{}, c)

	alg.MaxCrossovers = 50
	c = alg.CreateGenetic(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.GeneticAlgorithm{}, c)
}

func TestCreateSimulatedAnnealing(t *testing.T) {
	assert := assert.New(t)

	vertices := model2d.GenerateVertices(100)

	alg := &modelapi.Algorithm{AlgorithmType: modelapi.ALG_ANNEALING, MaxIterations: 25}
	c := alg.CreateSimulatedAnnealing(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.SimulatedAnnealing{}, c)

	alg.PreferCloseNeighbors = boolPointer(false)
	c = alg.CreateSimulatedAnnealing(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.SimulatedAnnealing{}, c)

	alg.PreferCloseNeighbors = boolPointer(true)
	c = alg.CreateSimulatedAnnealing(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.SimulatedAnnealing{}, c)

	alg.Seed = intPointer(1234)
	c = alg.CreateSimulatedAnnealing(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.SimulatedAnnealing{}, c)

	alg.TemperatureFunction = modelapi.TEMP_LINEAR
	c = alg.CreateSimulatedAnnealing(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.SimulatedAnnealing{}, c)

	alg.TemperatureFunction = modelapi.TEMP_GEOMETRIC
	c = alg.CreateSimulatedAnnealing(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.SimulatedAnnealing{}, c)

	alg.PrecursorAlgorithm = &modelapi.Algorithm{
		AlgorithmType: modelapi.ALG_CLOSEST_GREEDY,
	}
	c = alg.CreateSimulatedAnnealing(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.SimulatedAnnealing{}, c)

	alg.PrecursorAlgorithm = &modelapi.Algorithm{
		AlgorithmType: modelapi.ALG_DISPARITY_GREEDY,
	}
	c = alg.CreateSimulatedAnnealing(vertices, model2d.BuildPerimiter)
	assert.IsType(&circuit.SimulatedAnnealing{}, c)
}

func boolPointer(b bool) *bool {
	return &b
}

func intPointer(i int64) *int64 {
	return &i
}
