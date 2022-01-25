package solver

import (
	"math"

	"github.com/fealos/lee-tsp-go/graph"
	"github.com/fealos/lee-tsp-go/model"
)

// FindShortestPathNPNoChecks checks all possible combinations of paths to find the shortest path.
// It accepts an unordered set of vertices, and returns the ordered list of vertices.
// To minimize memory usage (N^2) it uses a tree and depth-first search, though its computation complexity is still O(n!).
func FindShortestPathNPNoChecks(vertices []model.CircuitVertex) ([]model.CircuitVertex, float64) {
	numVertices := len(vertices)
	if numVertices == 0 {
		return nil, 0.0
	}
	// Step 1: Prepare root of tree
	treeRoot := createNode(0, 0, nil)
	treeRoot.createChildren(numVertices)

	// Step 2: Navigate down tree, creating children as we go
	for current := treeRoot; treeRoot.shortestPath == nil; {
		// If at leaf node, compute length of path from root to leaf, then return up one level
		if current.depth == numVertices-1 {
			current.computeLeafPath(vertices)
			current = current.parent
		} else if nextChild := current.findUnprocessedChild(); nextChild != nil {
			// If uncomputed children nodes exist, process the next child node
			current = nextChild
			current.createChildren(numVertices)
		} else {
			// If no uncomputed children, find shortest path among them, update current node, and delete children
			for _, child := range current.children {
				if child.shortestPathLength < current.shortestPathLength {
					current.shortestPath = child.shortestPath
					current.shortestPathLength = child.shortestPathLength
				}
				child.deleteNode()
			}
			current.children = nil
			current = current.parent
		}
	}

	return treeRoot.shortestPath, treeRoot.shortestPathLength
}

// FindShortestPathNPWithChecks checks all possible combinations of paths to find the shortest path.
// It accepts an unordered set of vertices, and returns the ordered list of vertices.
// To minimize memory usage (N^2) it uses a tree and depth-first search, though its computation complexity is still O(n!).
// This will perform the following checks to reduce unnecessary computation:
// - Check if path to child intersects any ancestor paths (other than parent path, since they share a vertex), skip child if true.
// - (not implemented) Track found paths to minimize re-computation; O(n^4) memory usage, so may not be worth it.
// - Track best circuit length thus far, if child node would exceed that length, skip child if true.
func FindShortestPathNPWithChecks(vertices []model.CircuitVertex) ([]model.CircuitVertex, float64) {
	numVertices := len(vertices)
	if numVertices == 0 {
		return nil, 0.0
	}
	// Step 1: Prepare root of tree
	treeRoot := createNode(0, 0, nil)
	treeRoot.createChildren(numVertices)

	shortestPathLen := math.MaxFloat64
	_, isGraph := vertices[0].(*graph.GraphVertex)

	// Step 2: Navigate down tree, creating children as we go
	for current := treeRoot; treeRoot.shortestPath == nil; {
		// If at leaf node, compute length of path from root to leaf, then return up one level
		if current.depth == numVertices-1 {
			current.computeLeafPath(vertices)
			current = current.parent
		} else if nextChild := current.findUnprocessedChild(); nextChild != nil {
			// If uncomputed children nodes exist, process the next child node
			current = nextChild
			// Ignore paths that self intersect, since those are guaranteed to be non-optimal (for non-graphs), or are longer than the optimal path so far.
			// Do not perform this check for graphs, since those may need to be self-intersecting to produce the optimal result.
			if !isGraph && current.intersectsAnyAncestorsOrIsLongerThanBestLength(vertices, shortestPathLen) {
				// Use empty array and max length to ignore a branch in the tree.
				current.shortestPath = []model.CircuitVertex{}
				current = current.parent
			} else {
				current.createChildren(numVertices)
			}
		} else {
			// If no uncomputed children, find shortest path among them, update current node, and delete children
			current.shortestPath = []model.CircuitVertex{}
			for _, child := range current.children {
				if child.shortestPathLength < current.shortestPathLength {
					current.shortestPath = child.shortestPath
					current.shortestPathLength = child.shortestPathLength
					if child.shortestPathLength < shortestPathLen {
						shortestPathLen = child.shortestPathLength
					}
				}
				child.deleteNode()
			}
			current.children = nil
			current = current.parent
		}
	}

	return treeRoot.shortestPath, treeRoot.shortestPathLength
}

type treeNodeTSP struct {
	parent             *treeNodeTSP
	children           []*treeNodeTSP
	index              int
	depth              int
	shortestPathLength float64
	shortestPath       []model.CircuitVertex
}

func createNode(index int, depth int, parent *treeNodeTSP) *treeNodeTSP {
	return &treeNodeTSP{
		parent:             parent,
		children:           []*treeNodeTSP{},
		index:              index,
		depth:              depth,
		shortestPathLength: math.MaxFloat64,
		shortestPath:       nil,
	}
}

// createChildren creates one treeNodeTSP for each index that is not already and ancester of the current node nor is the current node.
func (t *treeNodeTSP) createChildren(numVertices int) {
	existingIndices := t.computeExistingIndices()
	for index := 0; index < numVertices; index++ {
		if _, exists := existingIndices[index]; !exists {
			child := createNode(index, t.depth+1, t)
			t.children = append(t.children, child)
		}
	}
}

// computeExistingIndices returns a map of all indices that are along the path from the root of the tree through the current node.
func (t *treeNodeTSP) computeExistingIndices() map[int]bool {
	existingIndices := make(map[int]bool)

	for current := t; current != nil; current = current.parent {
		existingIndices[current.index] = true
	}

	return existingIndices
}

// computeLeafPath stores the path from the root node to this leaf in this leaf, along with the length of the path.
func (t *treeNodeTSP) computeLeafPath(vertices []model.CircuitVertex) {
	// Only allow for leaf nodes to compute the path
	if t.depth != len(vertices)-1 {
		return
	}

	t.shortestPath = make([]model.CircuitVertex, len(vertices))

	var previousVertex model.CircuitVertex
	currentNode := t
	pathLength := 0.0
	// Since we are navigating the tree from leaf to root, need to add vertices to end of the slice first so that the root is the start of the slice
	for i := len(vertices) - 1; i >= 0; i-- {
		currentVertex := vertices[currentNode.index]
		t.shortestPath[i] = currentVertex

		// Compute the path length as we add vertices to the path, so that we don't have to iterate through the vertices multiple times
		if previousVertex != nil {
			pathLength += currentVertex.DistanceTo(previousVertex)
		}
		previousVertex = currentVertex
		currentNode = currentNode.parent
	}
	// Add the distance from the first vertex to the last vertex
	pathLength += previousVertex.DistanceTo(vertices[t.index])

	t.shortestPathLength = pathLength
}

// deleteNode cleans up the current node and any of the node's descendents.
func (t *treeNodeTSP) deleteNode() {
	t.parent = nil
	t.shortestPath = nil
	for _, c := range t.children {
		c.deleteNode()
	}
	t.children = nil
}

// findUnprocessedChild returns the first child in this node that does not have an optimal path.
// If all children in this node have been computed, nil is returned.
func (t *treeNodeTSP) findUnprocessedChild() *treeNodeTSP {
	for _, child := range t.children {
		if child.shortestPath == nil {
			return child
		}
	}
	return nil
}

// intersectsAnyAncestorsOrIsLongerThanBestLength returns true if the path to the current node is self intersecting or is shorter than the shortest completed circuit.
func (t *treeNodeTSP) intersectsAnyAncestorsOrIsLongerThanBestLength(vertices []model.CircuitVertex, bestLength float64) bool {
	// Ignore the first 3 vertices, since either there are no edges to intersect, or there is only one other edge (which ends at the parent to this vertex)
	if t.depth < 3 {
		return false
	}
	edgeT := vertices[t.parent.index].EdgeTo(vertices[t.index])
	pathLength := edgeT.GetLength()

	// Go to the parent two levels up, to ensure that there are no shared vertices between this edge and the current edge.
	for current := t.parent.parent; current.parent != nil; current = current.parent {
		edgeCurrent := vertices[current.parent.index].EdgeTo(vertices[current.index])
		pathLength += edgeCurrent.GetLength()
		if pathLength > bestLength {
			return true
		} else if edgeT.Intersects(edgeCurrent) {
			return true
		}
	}

	return false
}
