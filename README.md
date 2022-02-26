# TSP Solver - GO Package

[![Go Report Card](https://goreportcard.com/badge/github.com/heustis/tsp-solver-go)](https://goreportcard.com/report/github.com/heustis/tsp-solver-go)
[![GoDoc](https://pkg.go.dev/badge/github.com/heustis/heustis/tsp-solver-go)](https://pkg.go.dev/github.com/heustis/tsp-solver-go)
[![License](https://img.shields.io/dub/l/vibe-d.svg)](https://github.com/heustis/tsp-solver-go/blob/master/LICENSE)

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
  c := circuit.NewClosestGreedy(tspModel, model2d.BuildPerimiter, false)
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
Once the convex hull has been created, these algorithms take different approaches for attaching the interior points to the hull.

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

### Convex Concave - Closest Greedy
This algorithm performs the following steps after generating the convex hull:
1. If `cloneByInitEdges` is set to true: 
    1. The convex hull is cloned once per edge in the convex hull.
    2. For each edge, it updates the corresponding clone by attaching the point closest to that edge to it.
    3. Each of these clones is processed according to the remaining steps simultaneously (as in, all clones will be updated by their N-th point prior any clone having their N+1-th point attached).
2. Determines, and tracks, the closest edge to each point based on the distance increase that results from inserting the point along that edge
    * the distance increase can be calculated as:
      ```go
      start.EdgeTo(point).GetLength() + point.EdgeTo(end).GetLength() - start.EdgeTo(end).GetLength()
      ```
    * Internally, this uses a heap to store the points + closest edges to avoid reevaluating all edges each iteration. Selecting the next point just involves popping the next item off the heap.
3. Selects the next point to attach to the circuit by finding the point with the closest edge.
4. Attaches the point to its closest edge.
5. Updates the closest edge of remaining unattached points, by comparing their previous closest edges to the newly created edges (after the split)
    1. If `updateInteriorPoints` is set to true, this will also check each attached interior point (other than the one that was just attached) to see if either newly created edge is closer to that point than the edge it was originally attached to. 
    2. If one of the new edges is closer, it will be detached and added to the heap of points + closest edges.
    3. Once all attached interior points have been checked, and detached if appropriate, the closest edges for all unattached interior points are updated.
6. Repeats 3-5 until all points are attached to the circuit.
This algorithm greedily attaches points to the convex hull by prioritizing points that have the smallest impact on the length of the circuit. In other words, it prefers the point, that when attached to its closest edge (by distance increase), increases the length of the circuit by the least.

Complexity:
* This algorithm is O(n^2) because it needs to attach each interior point to the circuit, and each time it attaches an interior point it needs to check if the newly created edges are closer to each remaining interior point than their current closest edge (so that subsequent updates are correct).
* If `updateInteriorPoints` is enabled, this becomes O(n^3) due to the check for closest edges that occurs whenever a point is detached.
* If `cloneByInitEdges` is enabled, this becomes O(n^3) due to updating each clone in each iteration.
* If both `updateInteriorPoints` and `cloneByInitEdges` are enabled this becomes O(n^4).

### Convex Concave - Closest with Cloning
This behaves similarly to Closest Greedy, in that it first builds a convex hull, then selects interior points to attach to the hull based on whichever point has the minimum distance increase. Unlike Closest Greedy, this clones the entire circuit either whenever a point would be attached to a location, or whenever an attached point would be reattached at a different location. This allows this algorithm to explore possibilities that would be missed by the greedy algorithm. 

This algorithm performs the following steps after generating the convex hull:
1. Initializes metadata for each point:
    * whether the point is unattached,
    * whether the point is part of the initial convex hull,
    * the distance increase for the point, if it is attached and not part of the convex hull.
2. Determines, and tracks, the closest edge to each point based on the distance increase that results from inserting the point along that edge
    * the distance increase can be calculated as:
      ```go
      start.EdgeTo(point).GetLength() + point.EdgeTo(end).GetLength() - start.EdgeTo(end).GetLength()
      ```
    * Internally, this uses a heap to store the points + closest edges to avoid reevaluating all edges each iteration.
3.  Initializes the set of clones with an initial clone containing the convex hull, interior points, closest edges, and metadata.
    * Internally, this uses a heap to store the clones, sorted so that the shortest circuit is at the head of the heap.
4. Selects the clone to update by choosing the one with the shortest circuit.



To enable this behavior, this algorithm tracks each interior point and the distance increase that would result from attaching the point to each edge in the circuit, exluding edges that the point has already been attached to. When comparing the the effect of attaching a point to the circuit, on the length of the circuit, both the distance increase of the new location and the distance decrease of removing the existing location are taken into account.

Complexity:
* This algorithm is O(N!).
* If `maxClones` is enabled 

### Convex Concave - Disparity Greedy

### Convex Concave - Disparity With Cloning

### Simulated Annealing

### Genetic Algorithm
