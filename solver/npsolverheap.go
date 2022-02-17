package solver

import "github.com/heustis/lee-tsp-go/model"

// FindShortestPathNPHeap finds the shortest path by using a heap to grow the shortest circuit until in includes all the vertices.
// It accepts an unordered set of vertices, and returns the ordered list of vertices.
// This has a complexity of O(n!) and a memory usage of O(n!).
func FindShortestPathNPHeap(vertices []model.CircuitVertex) ([]model.CircuitVertex, float64) {
	// Prepare root of tree
	treeRoot := createHeapNode(0, 1, nil)
	numVertices := len(vertices)

	heap := model.NewHeap(func(a interface{}) float64 {
		return a.(*treeHeapNodeTSP).pathLength
	})

	// Add each child in the current node to the heap. The current node is the node in the heap with the shortest path so far.
	node := treeRoot
	for ; node.depth < numVertices; node = heap.PopHeap().(*treeHeapNodeTSP) {
		node.createChildren(numVertices)
		for _, c := range node.children {
			c.pathLength = c.computePathLen(vertices)
			heap.PushHeap(c)
		}
	}

	// Create a path from the root to the current node.
	pathLength := node.pathLength
	path := make([]model.CircuitVertex, numVertices)
	for n := node; n != nil; n = n.parent {
		path[n.depth-1] = vertices[n.index]
	}

	// Clean up the heap and tree.
	heap.Delete()
	treeRoot.deleteNode()

	return path, pathLength
}

type treeHeapNodeTSP struct {
	parent     *treeHeapNodeTSP
	children   []*treeHeapNodeTSP
	index      int
	depth      int
	pathLength float64
}

func createHeapNode(index int, depth int, parent *treeHeapNodeTSP) *treeHeapNodeTSP {
	return &treeHeapNodeTSP{
		parent:     parent,
		children:   []*treeHeapNodeTSP{},
		index:      index,
		depth:      depth,
		pathLength: 0.0,
	}
}

// createChildren creates one treeHeapNodeTSP for each index that is not already and ancester of the current node nor is the current node.
func (t *treeHeapNodeTSP) createChildren(numVertices int) {
	existingIndices := t.computeExistingIndices()
	for index := 0; index < numVertices; index++ {
		if _, exists := existingIndices[index]; !exists {
			t.children = append(t.children, createHeapNode(index, t.depth+1, t))
		}
	}
}

// computeExistingIndices returns a map of all indices that are along the path from the root of the tree through the current node.
func (t *treeHeapNodeTSP) computeExistingIndices() map[int]bool {
	existingIndices := make(map[int]bool)

	for current := t; current != nil; current = current.parent {
		existingIndices[current.index] = true
	}

	return existingIndices
}

// computePathLen determines the length of the path from the root to this node.
// If the current node is a leaf node, the length of the path back to the root is added to the circuit length.
func (t *treeHeapNodeTSP) computePathLen(vertices []model.CircuitVertex) float64 {
	if t.parent == nil {
		return 0.0
	} else if t.depth == len(vertices) {
		// The root node is always index 0, so we can determine the path length by utilizing the parent node's path length (rather than navigating through the tree).
		return t.parent.pathLength + vertices[t.parent.index].DistanceTo(vertices[t.index]) + vertices[0].DistanceTo(vertices[t.index])
	} else {
		return t.parent.pathLength + vertices[t.parent.index].DistanceTo(vertices[t.index])
	}
}

// deleteNode cleans up the current node and any of the node's descendents.
func (t *treeHeapNodeTSP) deleteNode() {
	t.parent = nil
	if t.children != nil {
		for _, c := range t.children {
			c.deleteNode()
		}
		t.children = nil
	}
}
