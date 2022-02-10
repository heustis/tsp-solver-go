package graph

type Graph struct {
	vertices []*GraphVertex
}

func NewGraph(vertices []*GraphVertex) *Graph {
	g := &Graph{vertices: vertices}
	for _, v := range vertices {
		v.InitializePaths()
	}
	return g
}

func (g *Graph) Delete() {
	for _, v := range g.vertices {
		v.Delete()
	}
	g.vertices = nil
}

func (g *Graph) GetVertices() []*GraphVertex {
	return g.vertices
}

func (g *Graph) String() string {
	s := `{"vertices":[`

	isFirst := true
	for _, v := range g.vertices {
		if !isFirst {
			s += ","
		}
		isFirst = false
		s += v.String()
	}

	s += `]}`
	return s
}
