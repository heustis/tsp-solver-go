package circuit_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/fealos/lee-tsp-go/circuit"
	"github.com/fealos/lee-tsp-go/model2d"
	"github.com/fealos/lee-tsp-go/tsplib"
	"github.com/stretchr/testify/assert"
)

func TestEdgeDistanceStats(t *testing.T) {
	assert := assert.New(t)

	for _, filename := range []string{
		"a280",
		"ali535",
		"att48",
		"att532",
		"berlin52",
		"bier127",
		// "brd14051",
		"burma14",
		"ch130",
		"ch150",
		"d198",
		"d493",
		"d657",
		"d1291",
		"d1655",
		"d2103",
		// "d18512",
		// "dsj1000",
		// "eil51",
		// "eil76",
		// "eil101",
		// "fl417",
		// "fl1400",
		// "fl1577",
		// "fl3795",
		// "fnl4461",
		// "gil262",
		// "gr96",
		// "gr137",
		// "gr202",
		// "gr666",
		// "kroA100",
		// "kroC100",
		// "kroD100",
		// "lin105",
		// "pcb442",
		// "pr76",
		// "pr1002",
		// "pr2392",
		// "rd100",
		// "st70",
		// "tsp225",
		// "ulysses16",
		// "ulysses22",
		// "usa13509",
	} {
		fmt.Println(filename)

		data, err := tsplib.NewData(fmt.Sprintf("../test-data/tsplib/%s.tsp", filename))
		assert.Nil(err)
		assert.NotNil(data)

		c := circuit.NewConvexConcaveConfidence(data.GetVertices(), model2d.DeduplicateVertices, &model2d.PerimeterBuilder2D{})
		c.Prepare()
		c.BuildPerimiter()

		f, err := os.OpenFile(fmt.Sprintf(`../results/gaps/%s.stats.json`, filename), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)
		assert.Nil(err)

		defer f.Close()

		f.WriteString(c.(*circuit.ConvexConcaveConfidence).ToString())
	}
}
