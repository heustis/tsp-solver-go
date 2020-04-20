package tsp

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		panic("Usage: " + args[0] + " <number_of_vertices>")
	}

	numVertices, err := strconv.Atoi(args[1])
	if err != nil || numVertices < 3 {
		panic("number_of_vertices must be an integer greater than 2")
	}

	vertices := generateVertices(numVertices)

	midpoint := findMidpoint(vertices)

	farthestVertex := findFarthestVertex(midpoint, vertices)
	farthestFromFarthestVertex := findFarthestVertex(farthestVertex, vertices)

	polygon := NewPolygon([]*Vertex{farthestVertex, farthestFromFarthestVertex})

	polygonJSON, err := json.Marshal(polygon)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(polygonJSON))
}

func generateVertices(size int) []*Vertex {
	var vertices []*Vertex
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		vertices = append(vertices, NewVertex(r.NormFloat64(), r.NormFloat64()))
	}
	return vertices
}

func findMidpoint(vertices []*Vertex) *Vertex {
	count := float64(len(vertices))
	midpoint := NewVertex(0, 0)
	for _, vertex := range vertices {
		midpoint = NewVertex(midpoint.GetX()+(vertex.GetX()/count), midpoint.GetY()+(vertex.GetY()/count))
	}
	return midpoint
}

func findFarthestVertex(origin *Vertex, vertices []*Vertex) *Vertex {
	farthestDistance := 0.0
	var farthestVertex *Vertex
	for _, current := range vertices {
		distanceToCurrent := origin.DistanceToSquared(current)
		if distanceToCurrent > farthestDistance {
			farthestDistance = distanceToCurrent
			farthestVertex = current
		}
	}
	return farthestVertex
}

func isExternal(vertex *Vertex, cloesestEdge *Edge, poly *Polygon) bool {

}
