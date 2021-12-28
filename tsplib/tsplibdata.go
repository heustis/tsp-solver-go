package tsplib

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fealos/lee-tsp-go/model"
	"github.com/fealos/lee-tsp-go/model2d"
)

type TspLibData struct {
	name            string
	comment         string
	numPoints       int
	vertices        []*model2d.Vertex2D
	bestRoute       []*model2d.Vertex2D
	bestRouteLength float64
}

func (data *TspLibData) GetBestRoute() []*model2d.Vertex2D {
	bestRouteCopy := make([]*model2d.Vertex2D, len(data.bestRoute))
	copy(bestRouteCopy, data.bestRoute)
	return bestRouteCopy
}

func (data *TspLibData) GetBestRouteLength() float64 {
	return data.bestRouteLength
}

func (data *TspLibData) GetComment() string {
	return data.comment
}

func (data *TspLibData) GetName() string {
	return data.name
}

func (data *TspLibData) GetNumPoints() int {
	return data.numPoints
}

func (data *TspLibData) GetVertices() []model.CircuitVertex {
	verticesCopy := make([]model.CircuitVertex, len(data.vertices))
	for i, v := range data.vertices {
		verticesCopy[i] = v
	}
	return verticesCopy
}

func (data *TspLibData) SolveAndCompare(solverName string, solver func([]model.CircuitVertex) model.Circuit) error {
	verticesCopy := make([]model.CircuitVertex, len(data.vertices))
	for i, v := range data.vertices {
		verticesCopy[i] = v
	}
	result := solver(verticesCopy)

	f, err := os.OpenFile(fmt.Sprintf(`../test-data/tsplib/output/%s.tsp.%s.tour`, data.name, solverName), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "NAME : %s\n", data.name)
	fmt.Fprintf(f, "TYPE : TOUR\n")
	fmt.Fprintf(f, "DIMENSION : %d\n", data.numPoints)
	fmt.Fprintf(f, "BEST KNOWN LENGTH : %f\n", data.bestRouteLength)
	fmt.Fprintf(f, "COMPUTED LENGTH : %f\n", result.GetLength())
	fmt.Fprintf(f, "TOUR_SECTION\n")
	for _, v := range result.GetAttachedVertices() {
		index := model.IndexOfVertex(verticesCopy, v)
		fmt.Fprintf(f, "%d\n", index)
	}
	fmt.Fprintf(f, "-1\n")

	return nil
}

func NewData(filePath string) (*TspLibData, error) {
	sourceFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer sourceFile.Close()

	sourceScanner := bufio.NewScanner(sourceFile)

	data := &TspLibData{}

	regexComment := regexp.MustCompile(`^COMMENT\s*:\s*(.+)$`)
	regexDimension := regexp.MustCompile(`^DIMENSION\s*:\s*([0-9]+)(.*)$`)
	regexName := regexp.MustCompile(`^NAME\s*:\s*(.+)$`)
	regexCoordinate := regexp.MustCompile(`\s*([0-9]+)\s*(-?[0-9.]+(?:e\+[0-9]+)?)\s*(-?[0-9.]+(?:e\+[0-9]+)?)$`)

	for inCoordinateSection := false; sourceScanner.Scan(); {
		line := strings.TrimSpace(sourceScanner.Text())
		if inCoordinateSection {
			if r := regexCoordinate.FindStringSubmatch(line); r != nil {
				var x, y float64
				x, err = strconv.ParseFloat(r[2], 64)
				if err != nil {
					return nil, fmt.Errorf(`failed to parse X coordinate from file=%s line=%s error=%v`, filePath, line, err)
				}
				y, err = strconv.ParseFloat(r[3], 64)
				if err != nil {
					return nil, fmt.Errorf(`failed to parse Y coordinate from file=%s line=%s error=%v`, filePath, line, err)
				}
				data.vertices = append(data.vertices, model2d.NewVertex2D(x, y))
			} else {
				break
			}
		} else if strings.Compare(line, `NODE_COORD_SECTION`) == 0 {
			inCoordinateSection = true
		} else if r := regexName.FindStringSubmatch(line); r != nil {
			data.name = strings.TrimSpace(r[1])
		} else if r := regexComment.FindStringSubmatch(line); r != nil {
			data.comment = strings.TrimSpace(r[1])
		} else if r := regexDimension.FindStringSubmatch(line); r != nil {
			data.numPoints, _ = strconv.Atoi(strings.TrimSpace(r[1]))
			data.vertices = make([]*model2d.Vertex2D, 0, data.numPoints)
			data.bestRoute = make([]*model2d.Vertex2D, 0, data.numPoints)
			data.bestRouteLength = 0.0
		}
	}

	if sourceScanner.Err() != nil {
		return nil, sourceScanner.Err()
	}

	// Once the data has been loaded, also load the best known path for comparison.
	// Not all test data will have a best known solution, ignore those cases.
	solutionFile, err := os.Open(filePath + ".opt.tour")
	if err != nil {
		return data, nil
	}
	defer solutionFile.Close()

	solutionScanner := bufio.NewScanner(solutionFile)

	for numParsed, inTourSection := 0, false; solutionScanner.Scan(); {
		line := strings.TrimSpace(solutionScanner.Text())
		if inTourSection {
			if parsedIndex, err := strconv.Atoi(line); err == nil && parsedIndex > 0 {
				v := data.vertices[parsedIndex-1] // Need -1 to make index 0-based.
				data.bestRoute = append(data.bestRoute, v)
				numParsed++
				if numParsed > 1 {
					data.bestRouteLength += data.bestRoute[numParsed-2].DistanceTo(v)
				}
			} else {
				break
			}
		} else if strings.Compare(line, `TOUR_SECTION`) == 0 {
			inTourSection = true
		}
	}

	if l := len(data.bestRoute); l > 0 {
		data.bestRouteLength += data.bestRoute[l-1].DistanceTo(data.bestRoute[0])
	}

	return data, nil
}
