package solver

import (
	"math"

	"github.com/fealos/lee-tsp-go/model"
)

// FindShortestPathNP checks all possible combinations of paths to find the shortest path.
// It accepts an unordered set of vertices, and returns the ordered list of vertices.
// To minimize memory usage (N^2) it uses a tree and depth-first search.
func FindShortestPathNP(vertices []model.CircuitVertex) ([]model.CircuitVertex, float64) {
	// Step 1: Prepare root of tree
	treeRoot := createNode(0, 0, nil)
	numVertices := len(vertices)
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

func (t *treeNodeTSP) createChildren(numVertices int) {
	existingIndices := t.computeExistingIndices()
	for index := 0; index < numVertices; index++ {
		if _, exists := existingIndices[index]; !exists {
			t.children = append(t.children, createNode(index, t.depth+1, t))
		}
	}
}

func (t *treeNodeTSP) computeExistingIndices() map[int]bool {
	existingIndices := make(map[int]bool)

	for current := t; current != nil; current = current.parent {
		existingIndices[current.index] = true
	}

	return existingIndices
}

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

func (t *treeNodeTSP) deleteNode() {
	t.parent = nil
	t.shortestPath = nil
	for _, c := range t.children {
		c.deleteNode()
	}
	t.children = nil
}

func (t *treeNodeTSP) findUnprocessedChild() *treeNodeTSP {
	for _, child := range t.children {
		if child.shortestPath == nil {
			return child
		}
	}
	return nil
}
