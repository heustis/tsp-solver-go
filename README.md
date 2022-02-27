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
  validate := validator.New()
  if err := validate.Struct(request); err != nil {
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

#### About
These set of algorithms first build the minimum convex hull around the set of points, so that all points are either vertices on the hull or interior to the hull.
Once the convex hull has been created, these algorithms take different approaches for attaching the interior points to the hull.

The reasong to create the convex hull first, is that points in the convex hull must be traversed in that order regardless of where the internal points attach to that hull. If any of the points in the hull were to be visted in a different order, it would result in the circuit self-intersecting which is less efficient than a non-self-intersecting circuit.

#### Convex Hull Steps
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

#### Complexity
The complexity of creating the convex hull is O(n^2) because:
* One point is attached during each iteration (to a max of n iterations).
* Each iteration the exterior point farthest from its closest edge is selected in step 5, which is O(n).
* Each iteration the closest edge for each exterior points is updated and each point is checked to see if they are now interior, which is O(n).
    * Updating the closest edge for a single point is O(1) because only the two new edges and the previous edge are considered.
    * Checking for interior status for a single point is O(1) because only the point, closest edge, and midpoint are used.
* The selection and update are done independently each iteration, so the complexity of each iteration it is the maximum of their complexity, O(n).

### Convex Concave - Closest Greedy

#### About
This algorithm builds the convex hull, then greedily attaches points to the circuit by prioritizing points that have the smallest impact on the length of the circuit. In other words, it prefers the point, that when attached to its closest edge (by distance increase), increases the length of the circuit by the least.

* the distance increase is calculated as:
  ```go
  start.EdgeTo(point).GetLength() + point.EdgeTo(end).GetLength() - start.EdgeTo(end).GetLength()
  ```

#### Steps
This algorithm performs the following steps after generating the convex hull:
1. If `cloneByInitEdges` is set to true: 
    1. The convex hull is cloned once per edge in the convex hull.
    2. For each edge, it updates the corresponding clone by attaching the point closest to that edge to it.
    3. Each of these clones is processed according to the remaining steps simultaneously (as in, all clones will be updated by their N-th point prior any clone having their N+1-th point attached).
2. Calculates, and stores, the closest edge to each point based on the distance increase that results from inserting the point along that edge
    * Internally, this uses a heap to store the points + closest edges to avoid reevaluating all edges each iteration. Selecting the next point just involves popping the next item off the heap.
3. Selects the next point to attach to the circuit by finding the point with the closest edge.
4. Attaches the point to its closest edge.
5. Updates the closest edge of remaining unattached points, by comparing their previous closest edges to the newly created edges (after the split)
    1. If `updateInteriorPoints` is set to true, this will also check each attached interior point (other than the one that was just attached) to see if either newly created edge is closer to that point than the edge it was originally attached to. 
    2. If one of the new edges is closer, it will be detached and added to the heap of points + closest edges.
    3. Once all attached interior points have been checked, and detached if appropriate, the closest edges for all unattached interior points are updated.
6. Repeats 3-5 until all points are attached to the circuit.

#### Complexity:
* This algorithm is O(n^2) because it needs to attach each interior point to the circuit, and each time it attaches an interior point it needs to check if the newly created edges are closer to each remaining interior point than their current closest edge (so that subsequent updates are correct).
* If `updateInteriorPoints` is enabled, this becomes O(n^3) due to the check for closest edges that occurs whenever a point is detached.
* If `cloneByInitEdges` is enabled, this becomes O(n^3) due to updating each clone in each iteration.
* If both `updateInteriorPoints` and `cloneByInitEdges` are enabled this becomes O(n^4).

### Convex Concave - Closest with Cloning

#### About
This behaves similarly to Closest Greedy, in that it first builds a convex hull, then selects interior points to attach to the hull based on whichever point has the minimum distance increase. Unlike Closest Greedy, this clones the entire circuit either whenever a point would be attached to a location, or whenever an attached point would be reattached at a different location. This allows this algorithm to explore possibilities that would be missed by the greedy algorithm. 

#### Steps
1. Generates the convex hull (as described above).
2. Initializes metadata for each point:
    * whether the point is unattached,
    * whether the point is part of the initial convex hull,
    * the distance increase for the point, if it is attached and not part of the convex hull.
2. Calculates, and stores, the distance increase for each point and each edge.
    * Unlike Closest Greedy this tracks all points and edge combinations, to allow for points to attach to their 2nd, 3rd, ..., Nth edges.
    * Internally, this uses a heap to store the points + edge distances to avoid reevaluating all edges each iteration.
3.  Initializes the set of clones with an initial clone containing the convex hull, interior points, closest edges, and metadata.
    * Internally, this uses a heap to store the clones, sorted so that the shortest circuit is at the head of the heap.
4. Selects the clone to update by choosing the one with the shortest circuit.
5. Selects the next point and edge to attach to the clone's circuit by finding the point with the next smallest distance increase.
6. Attaches the point to the edge:
    * If `cloneOnFirstAttach` is `true`, this will:
        1. create a clone,
        2. attach the point, to the edge from 5, in the clone,
        3. update the metadata for the attached point and the two points in the split edge (attached, distance increase),
        4. remove any point+edge combinations with the attached point from the heap,
        5. replace any point+edge combinations containing the split edge in the heap, with two point+edge combinations for the two new edges, and
        6. update the heap distance increases for all of the affected points.
    * If `cloneOnFirstAttach` is `false`, and this is the first time the point is attached, this will:
        1. attach the point, to the edge from 5,
        2. update the metadata for the attached point and the two points in the split edge (attached, distance increase),
        3. replace any point+edge combinations containing the split edge in the heap, with two point+edge combinations for the two new edges, and
        4. update the heap distance increases for all of the affected points.
    * If `cloneOnFirstAttach` is `false`, and the point has been attached previously, this will:
        1. create a clone,
        2. detach the point from its current edge in the clone, by merging the two edges that the point was a part of,
        3. attach the point, to the edge from 5, in the clone,
        4. update the metadata for the attached point, the two points in the split edge, and the two points in the merged edge (attached, distance increase),
        5. remove any point+edge combinations with the attached point from the heap
        6. replace any point+edge combinations containing the split edge in the heap, with two point+edge combinations for the two new edges,
        7. replace any point+edge combinations containing either of the merged edges in the heap, with one entry for the merged edge, and
        8. update the heap distance increases for all of the affected points.
7. If `maxClones` is configured, and the number of clones exceeds the maximum, discard the clone with the worst length per attached point.

#### Complexity:
* This algorithm is O(n!).

### Convex Concave - Disparity Greedy

#### About
To understand why this algorithm works, consider the locations where an unattached point can be within a circuit: 
1. near a single edge,
    * This will have a significant disparity between the distance increase of its closest edge, and the distance increase of all other edges.
2. near a corner of two edges,
    * This will have an insignificant disparity between the distance increase of the two corner edges, and but a significant disparity between those two edges and the distance increase of all other edges.
3. in the middle of several edges,
    * The number of edges with an insignificant disparity is typically greater than two, but is more variable than the other locations.

This algorithm prioritizes category 1 points, since their closest edge is likely to be their optimum location, and defers processing category 2 points, which are harder to predict the optimum location.

As points are attached to the circuit, points in category 3 will move into categories 1 and 2 due to concave edges becoming closer to them than the initial convex edges were. Some points may become external points, since this doesn't prioritize the closest points, but that is okay as the closest edge to any of those points will be one of the new edges (so it won't create intersecting edges).

Eventually category 2 points need to be selected, but the earlier selections should improve the accuracy of these selections and reduce the impact of incorrect selections on the length of the circuit.

#### Steps
1. Generates the convex hull (as described above).
2. Calculates, and stores, the closest two edges to each point, based on the distance increase from inserting the point along that edge.
3. For each point, determines the disparity between the two closest edges:
    * If `useRelativeDisparity` is true, this calculates the disparity by dividing the larger distance increase by the smaller.
    * If `useRelativeDisparity` is false, this calculates the disparity by subtracting the smaller distance increase from the larger.
4. Selects the next point to attach to the circuit, by finding the point with the largest disparity between its two closest edges.
    * If two points have the same disparity, the point that is closer to its closest edge is chosen.
5. Attaches the selected point to its closest edge.
6. Updates the remaining unattached points, by comparing their previous closest two edges to the newly created edges (after the split), and updating their disparity if they are updated.
7. Repeats 3-5 until all points are attached to the circuit.

#### Complexity:
* This algorithm is O(n^2) because it needs to attach each interior point to the circuit, and each time it attaches an interior point it needs to check if the newly created edges are closer to each remaining interior point than their current closest edges, so that it can update their disparity and select the correct point + edge in subsequent iterations.

### Convex Concave - Disparity With Cloning

#### About
This behaves similarly to AlgorithmDisparityGreedy in that it first builds a convex hull and then it prioritizes points based on the disparity in distance increases from those points to the edges. However, unlike the greedy algorithm, this will clone the circuit if a point is close to multiple edges, and attach it to each of those edges.

To determine whether a gap is significant, this computes the following statistics for each point (in each clone):  
* the average gap in distance increases
* and standard deviation of the gaps

When selecting which point to attach next, this first finds the earliest significant gap. If multiple points have significant gaps at the same position, the most significant gap of those is chosen. For example:
* Point A has distance increases of 2, 3, 10, 12, ... (gaps 1, 7, 2, ...)
* Point B has distance increases of 1, 7, 8, 10, ... (gaps 6, 1, 2, ...)
* Point C has distance increases of 1, 2, 6, 7, ... (gaps 1, 4, 1, ...)
* Point B will be chosen first since its significant gap is at the 0th index
* Point A will be chosen second since it and C have significant gaps at the 1st index, but A's is more significant.  
* Since point A's gap is after two edges, the current circuit will be cloned into two circuits (incliding the current circuit), each with point A attaching a different one of those edges.  
* If point A were after three edges, three circuits would be created/updated (including the initial circuit).

Note: unlike AlgorithmClosestClone once a point is attached to a circuit, that point is not allowed to move.

#### Steps
1. Generates the convex hull (as described above).
2. Calculates, and stores, the distance increase to each edge for each point.
3. Sorts each point's distances increases from smallest to largest.
4. Calculates the following statistics for each point's distance increases:
    1. The gap between each distance increase (i.e. `gap[i] = distance[i+1] - distance[i]`).
    2. The average gap size.
    3. The standard deviation of the gaps.
5. Initializes the set of clones with the initial circuit data.
6. For each clone in the set (each iteration of this algorithm updates all the clones):
    1. Selects the point with earliest significant gap.
        * Earliest means smallest index in the gap array.
        * Significance defaults to 1 standard deviation, but is configurable
        * If multiple points have a significant gap at the same index, from those, the point with the closest edge is selected.
        * If no points have significant gaps, the point with the closest edge is selected.
    2. A clone of the current clone is created for each edge in the range `[1:n)` where `n` is the index of the significant gap.
    3. Each clone has the selected point attached to the corresponding edge (see the next step for what that entails).
    4. The current clone has the selected point attached to its closest edge:
        1. The point is inserted into the circuit by splitting the edge into two edges.
        2. The length of the circuit is updated.
        3. The stats of each point are updated by removing the split edge and adding the two new edges.  
           This resorts the distance increases and recomputes the gaps, average gap, and standard deviations.
7. Sort the clones based on circuit length per attached vertex.
8. If there are more clones that the configured maximum clones (default: 1000), trim the set of clones to only retain the configured amount.
9. Repeat steps 6-8 until the shortest clone is completed.

#### Complexity
* This algorithm is O(n^3 * log(maxClones * n^2) * maxClones) because each iteration:
  * one point is attached in each clone, so the total number of iterations is the number of points,
  * up to `maxClones` clones will have points attached,
  * each time a point is selected for attachment it is possible to create one clone per edge, and
  * the clones are sorted by length prior to trimming them to `maxClones`.

### Simulated Annealing

#### About
This implements [simulated annealing](https://en.wikipedia.org/wiki/Simulated_annealing) to stochastically approximate the optimum circuit through a set of points. 

Unlike the convex-concave algorithms this does not start from a convex hull and work towards a completed circuit. Rather this treats the supplied set of points as the initial circuit, or uses another algorithm to create an initial circuit, and mutates its to try to find a better sequencing of points for the circuit.

#### Steps
1. Randomly selects 2 points.
    * If `preferCloseNeighbors` is `true`, when selecting a second point it will prefer points that are close to the first selected point.
2. Determine how swapping the 2 points impacts the circuit length.
    * i.e. how much does swapping the points lengthen or shorten the circuit?
3. Scale this value based on the size of the coordinate space being used, so that it is meaningful regardless of if the coordinates are from -100 to +100 or -100000 to +100000
4. Use the configured temperature function to determine the acceptance value (based on the number of iterations, max iterations, and impact of the swap).
    * The temperature function is designed to reduce the probability of accepting a "bad" swap as the number of iterations approaches the maximum iterations. This allows early swaps to avoid local maxima, but later swaps focus on refining the current circuit towards its local maximum.
5. Generate a random number between `[0.0, 1.0)` as a test value.
6. Swap the points the test value it is less than the acceptance value, or if the swap would shorten the circuit.


#### Complexity
* This algorithm is O(maxIterations) because each iteration affects a maximum of 6 points (the 2 to swap and each of their adjacent points).
* If `preferCloseNeighbors` is `true`, the complexity is `O(n * maxIterations)` as selecting a neighboring point is O(n).
* If `precursorAlgorithm` is configured, the complexity is the maximum of O(precursorAlgorithm) and O(maxIterations).

### Genetic Algorithm

#### About
This implements a [genetic algorithm](https://en.wikipedia.org/wiki/Genetic_algorithm) to stochastically approximate the optimum circuit through a set of points. 

Unlike the convex-concave algorithms this does not start from a convex hull and work towards a completed circuit. Rather this starts with a randomly generated set of circuits, mutates them to create new circuits, and repeats this process a number of times. Returning the best circuit found during this process.

#### Steps
1. Initialization - a random set of parent circuits are created. By default these are random circuits, but users can optionally have the circuits based on the optimum convex hull.
2. Child circuits are created by:
      1. Randomly selecting two parent circuits.
      2. Using crossover to blend the parent circuits into a new circuit.
      3. Fixing any duplicate or missing points that result from the crossover.
      4. Mutating the circuit (as described in "mutationRate")
3. Selection - this uses an elitist slection algorithm, to ensure that the best solutions are not lost from one generation to the next. To do this the new children are combined with the previous generation of parents, and the top "numParents" (based on shortest circuit length) are retained for the next iteration.
4. Termination - this repeats steps 2 and 3 "maxIterations" times, then returns the best circuit found by this process.

#### Complexity
* This algorithm is the max of `O(numIterations * numChildren * n)` and `O(numIterations * (numParents+numChildren)*log(numParents+numChildren))` because each iteration:
  * starts with `numParents` circuits,
  * creates `numChildren` circuits,
  * performs up to `maxCrossovers` crossovers each time a child is created `(numChildren * n)`,
  * performs up to `n` mutations each time a child is created `(numChildren * n)`,
  * the parents and children are combined, selected, and trimmed to form the next generation of parents `O((numParents+numChildren)*log(numParents+numChildren))`
* If `shouldBuildConvexHull` is `true`, the complexity of the algorithm if the max of O(n^2) and the previous maximum.
