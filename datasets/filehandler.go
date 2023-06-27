package dataset

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/franciscobonand/aoc-max-clique/aco"
)

type Dataset struct {
    Input aco.Graph
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
func Write(fname string, data [][]float64) error {
    content := "gen,evals,repeated,bestfit,worstfit,meanfit,maxsize,minsize,meansize,betterCxChild,worseCxChild\n"
    for i, line := range data {
        content += fmt.Sprintf("%d,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f\n",
            i,
            line[0],
            line[1],
            line[2],
            line[3],
            line[4],
            line[5],
            line[6],
            line[7],
            line[8],
            line[9],
        )
    }
    bcontent := []byte(content)
    f := fmt.Sprintf("analysis/%s", fname)
    return os.WriteFile(f, bcontent, fs.ModePerm)
}
