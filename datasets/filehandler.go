package dataset

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/franciscobonand/aoc-max-clique/aco"
	sts "github.com/montanaflynn/stats"
)

type Dataset struct {
	Input  aco.Graph
	Output []float64
}

// Read reads a file resided in the given path.
// The path is relative to the directory the program is executed
func Read(fpath string, maxPheromone float64) (*Dataset, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ds := Dataset{}
	ds.Input = map[string]map[string]float64{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		content := strings.Split(line, " ")
		if content[0] == "e" {
			if content[1] == content[2] {
				continue
			}
			if _, ok := ds.Input[content[1]]; !ok {
				ds.Input[content[1]] = map[string]float64{}
			}
			if _, ok := ds.Input[content[2]]; !ok {
				ds.Input[content[2]] = map[string]float64{}
			}
			ds.Input[content[1]][content[2]] = maxPheromone
			ds.Input[content[2]][content[1]] = maxPheromone
		}
	}

	return &ds, nil
}

// TODO (currently no stats are being generated)
func Write(fname string, data [][]float64, elapsed time.Duration) error {
	content := "gen,best,worst,mean,rep,sdev\n"
	max := make([]float64, 0)
	for i, line := range data {
		content += fmt.Sprintf("%d,%f,%f,%f,%f,%f\n",
			i+1,
			line[0],
			line[1],
			line[2],
			line[3],
			line[4],
		)
		max = append(max, line[0])
	}
	maxSdev, _ := sts.StandardDeviation(max)
	content += fmt.Sprintf("BestSdev: %f, Time: %s", maxSdev, elapsed)
	bcontent := []byte(content)
	f := fmt.Sprintf("analysis/%s", fname)
	return os.WriteFile(f, bcontent, fs.ModePerm)
}
