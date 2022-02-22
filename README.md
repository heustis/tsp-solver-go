# TSP Solver - GO Package

[![Go Report Card](https://goreportcard.com/badge/github.com/heustis/tsp-solver-go)](https://goreportcard.com/report/github.com/heustis/tsp-solver-go)

## Table of Contents
1. [How To Use](#how-to-use)
2. [Vertex Types](#vertex-types)
3. [Algorithms](#algorithms)

## How To Use

### Using the package in a service

1. Import the package into your project.  
  ```sh
  go get github.com/heustis/tsp-solver-go
  ```
2. Convert your data into one of the supported vertex types (see below), or create your own.  
4. Construct one or more algorithms (from the `circuit` subpackage) for approximating the optimum the circuit.
5. Use the solver package to update the circuit until it is complete.
6. Get the best computed circuit from the circuit.

For example:
```go
func convert(myData []MyModel) []model.CircuitVertex {
  tspModel := make([]model.CircuitVertex, len(myData))
  for i, myVertex := range myData {
    tspModel[i] = model2d.NewVertex2D(myVertex.X, myVertex.Y)
  }
  return tspModel
}

func ProcessMyData(myData []MyModel) []model.CircuitVertex {
  tspModel := convert(myData)
  c := circuit.NewConvexConcave(tspModel, model2d.BuildPerimiter, false)
  solver.FindShortestPathCircuit(c)
  return c.GetAttachedVertices()
}
```

### Using the package to back a JSON API

1. Read through the [OpenApi document](https://github.com/heustis/tsp-solver-go/blob/master/openapi.yaml) to understand the prebuilt API.
2. Import the package into your project.  
  ```sh
  go get github.com/heustis/tsp-solver-go
  ```
3. In your project, either:
    1. use the prebuilt API models (in the `modelapi` sub-package),
    2. compose or nest the prebuilt models into structs that match your API requirements, or
    3. create your own models and translation layer.
4. If using the prebuilt API, use
  ```go
    var request *modelapi.TspRequest
    if err := json.Unmarshal([]byte(requestJson), &request); err != nil {
		  // Handle error
	  }
    response := solver.FindShortestPathApi(request)
  ```

### Contributing to the package

1. Fork this repository
2. Create a new branch with the prefix `hotfix/` or `feature/` and a suffix describing why the branch exists (e.g. the ID of an issue, or brief description of the feature being added)
3. Make your changes.
4. Run the tests
    ```sh
     go test ./... -cover
    ```
    1. Make sure all tests still pass
    2. Ensure that the code coverage has not decreased as a result of your changes (TODO - upload file with coverage for comparison)
6. Add or update any relevant documentation:
    1. README
    2. OpenAPI
    3. Performance Comparisons
7. Push your changes, on your branch, to your fork
8. Create a pull request (in github.com/heustis/tsp-solver-go) from your branch to master
    1. An authorized user will review your changes.
    2. If they leave any comments, please respond to the comment or make the requested change in a timely manner.
    3. Once approved, an authorized user will merge your changes and release a new version of the project (in accordance with the [Go docs](https://go.dev/doc/modules/version-numbers))

## Vertex Types

### 2-Dimensional

#### Go
```go
vertices := []model.CircuitVertex{
  model2d.NewVertex2D(-15, -15),
  model2d.NewVertex2D(0, 0),
  model2d.NewVertex2D(15, -15),
  model2d.NewVertex2D(3, 0),
  model2d.NewVertex2D(15, -15+model.Threshold/2.0),
  model2d.NewVertex2D(3, 13),
  model2d.NewVertex2D(8, 5),
  model2d.NewVertex2D(9, 6),
  model2d.NewVertex2D(-7, 6),
}
vertices = model2d.DeduplicateVertices(vertices)
```

#### JSON
```json
{
  "x":0.0000012345,
  "y":123.45
}
```

### 3-Dimensional

#### Go
```go
vertices := []model.CircuitVertex{
  model3d.NewVertex3D(-15, -15, -5.0),
  model3d.NewVertex3D(0, 0, model.Threshold/9.0),
  model3d.NewVertex3D(15, -15, -5.0),
  model3d.NewVertex3D(-15-model.Threshold/3.0, -15, -5),
  model3d.NewVertex3D(0, 0, 0.0),
  model3d.NewVertex3D(0, model.Threshold*2, 0.0),
  model3d.NewVertex3D(-15+model.Threshold/3.0, -15-model.Threshold/3.0, -5+model.Threshold/4),
  model3d.NewVertex3D(3, 0, 3),
  model3d.NewVertex3D(3, 13, 4),
  model3d.NewVertex3D(7, 6, 5),
  model3d.NewVertex3D(-7, 6, 5),
}
vertices = model.DeduplicateVerticesNoSorting(vertices)
```

#### JSON
```json
{
  "x":1.23,
  "y":2.34,
  "z":4.56
}
```

### Graph (graph.GraphVertex)

#### Go
```go
vA := graph.NewGraphVertex("a")
vB := graph.NewGraphVertex("b")
vC := graph.NewGraphVertex("c")
vA.AddAdjacentVertex(vB, 123.4)
vB.AddAdjacentVertex(vA, 123.4)
vB.AddAdjacentVertex(vC, 23.45)
vC.AddAdjacentVertex(vB, 34.56) // the distances between neighbors can be asymmetric
vC.AddAdjacentVertex(vA, 100)   // the neighbors can be asymmetric

// Graphs can have circular references, and need to be cleaned up.
// The easiest way to handle this is to:
// 1. create a NewGraph wrapping the vertices,
// 2. defer deleting the Graph,
// 3. apply the algorithms within this scope
// 4. convert the best computed circuit into an non-circular reference format (e.g. modelapi.ToApiFromGraph(g))
g := graph.NewGraph([]*GraphVertex{vA, vB, vC})
defer g.Delete()
```

#### JSON
```json
[
  {
    "name": "a",
    "neighbors": [
      { "name": "b", "distance": 123.4 }
    ]
  },
  {
    "name": "b",
    "neighbors": [
      { "name": "a", "distance": 123.4 },
      { "name": "c", "distance": 23.45 }
    ]
  },
  {
    "name": "c",
    "neighbors": [
      { "name": "b", "distance": 34.56 },
      { "name": "a", "distance": 100 }
    ]
  }
]
```

## Algorithms

This sections defines the various apporoaches for approximating or solving the TSP that are supported by this package.

Most algorithms in tsp-solver-go are located in the `circuit` sub-package, including all algorithms supported by the HTTP API.
These algorithms all implement the `model.Circuit` interface.

There are a couple of NP-complete algorithms in the `solver` package, which do not implement this interface, but those are for testing purposes.

### Convex Concave Variants

These set of algorithms first build the minimum convex hull around the set of points, so that all points are either vertices on the hull or interior to the hull.
Once the convex hull has been created, the algorithms take different approaches for determining the 

The reasong to create the convex hull first, is that points in the convex hull must be traversed in that order regardless of where the internal points attach to that hull. If any of the points in the hull were to be visted in a different order, it would result in the circuit self-intersecting which is less efficient than a non-self-intersecting circuit.

The algorithm this package uses to compute the convex hull is:
1. Compute the midpoint of the points. 
    * This may be an existing point (e.g. in a graph) or a new temporary point which is discarded after the hull is created (e.g. in 2D and 3D).
2. Find the point farthest from the midpoint.  
    * This constrains the possible locations of the remaining points to a circle or sphere with a radius equal to the distance of this point from the midpoint.
3. Find the point farthest from the point in 2.  
    * This further constrains all points, by intersecting the previous circle/sphere with a circle/sphere centered around the point from 2, with a radius to the point in 3.
4. Creates initial edges 1b->1c and 1c->1b _(note: all other points are exterior at this time)_  
5. Finds the exterior point farthest from its closest edge and attach it to the circuit by splitting its closest edge.  
    * For computing this distance, in 2D and 3D this package uses the perpendicular distance of the point to the edge, in graphs it uses the distance increase that results from attaching the point to the edge.
    * By selecting the farthest point by perpendicular distance, the angle through this point is guaranteed to not exceed 180 degrees (i.e. it must be convex). Having equally far points to either side of this point would produce a straight line through this point, and no point can be farther off or it would have been selected.
6. Find any points that were external to the circuit and are now internal to the circuit, and stop considering them for future iterations.  
7. Repeat 5 and 6 until all points are attached to the circuit or internal to the circuit.

#### Smallest Increase

#### Smallest Increase With Cloning

#### Disparity

#### Disparity With Cloning

### Simulated Annealing

### Genetic Algorithm
