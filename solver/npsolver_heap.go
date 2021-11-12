package solver

import (
	"github.com/fealos/lee-tsp-go/model"
)

// FindShortestPathNPHeap finds the shortest path by using a heap to grow the shortest circuit until in includes all the vertices.
// It accepts an unordered set of vertices, and returns the ordered list of vertices.
func FindShortestPathNPHeap(vertices []model.CircuitVertex) ([]model.CircuitVertex, float64) {
	// Step 1: Prepare root of tree
	treeRoot := createHeapNode(0, 1, nil)
	numVertices := len(vertices)

	heap := model.NewHeap(func(a interface{}) float64 {
		return a.(*treeHeapNodeTSP).pathLength
	})

	node := treeRoot
	for ; node.depth < numVertices; node = heap.PopHeap().(*treeHeapNodeTSP) {
		node.createChildren(numVertices)
		for _, c := range node.children {
			c.pathLength = c.computePathLen(vertices)
			heap.PushHeap(c)
		}
	}

	pathLength := node.pathLength
	path := make([]model.CircuitVertex, numVertices)
	for n := node; n != nil; n = n.parent {
		path[n.depth-1] = vertices[n.index]
	}

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

func (t *treeHeapNodeTSP) createChildren(numVertices int) {
	existingIndices := t.computeExistingIndices()
	for index := 0; index < numVertices; index++ {
		if _, exists := existingIndices[index]; !exists {
			t.children = append(t.children, createHeapNode(index, t.depth+1, t))
		}
	}
}

func (t *treeHeapNodeTSP) computeExistingIndices() map[int]bool {
	existingIndices := make(map[int]bool)

	for current := t; current != nil; current = current.parent {
		existingIndices[current.index] = true
	}

	return existingIndices
}

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

func (t *treeHeapNodeTSP) deleteNode() {
	t.parent = nil
	if t.children != nil {
		for _, c := range t.children {
			c.deleteNode()
		}
		t.children = nil
	}
}
