package main

import (
	crand "crypto/rand"
	"flag"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/franciscobonand/aoc-max-clique/aco"
	dataset "github.com/franciscobonand/aoc-max-clique/datasets"
)

var (
	ants, generations                           int
	file                                        string
	minPheromone, maxPheromone, evaporationRate float64
	seed                                        int64
	verbose, getstats                           bool
)

func main() {
	// go run . -ants 20 -gens 20 -minpheromone 0.1 -maxpheromone 10 -evaporation 0.1 -file "datasets/easy.col" -seed 4132
	initializeFlags()
	if evaporationRate <= 0.0 || evaporationRate >= 1.0 {
		log.Fatalln("evaporation must be between 0.0 and 1.0")
	}
	if minPheromone <= 0 || maxPheromone <= 0 || maxPheromone <= minPheromone {
		log.Fatalln("must be '0 < minpheromone < maxpheromone'")
	}

	ds, err := dataset.Read(file, maxPheromone)
	if err != nil {
		panic(err.Error())
	}

	fileaux := strings.Split(file, "/")
	file = fileaux[len(fileaux)-1]
	file = strings.Split(file, ".")[0]
	statsfile := strings.Join([]string{file, strconv.Itoa(ants), strconv.Itoa(generations), fmt.Sprintf("%.1f", evaporationRate)}, "-")

	fmt.Println(statsfile)
	// This will be used to run multiple experiments to generate stats
	var runqnt int64 = 1
	var run int64
	if getstats {
		runqnt = 30
	}

	allstats := make([][][]float64, runqnt)

	start := time.Now()

	for run = 0; run < runqnt; run++ {
		setSeed(seed + run)
		colony := aco.NewColony(ants, generations, minPheromone, maxPheromone, evaporationRate, ds.Input)
		maxClique, stats := colony.Run()
		if verbose {
			for i, stat := range stats {
				fmt.Printf("gen:%d,best:%f,worst:%f,mean:%f,rep:%f,sdev:%f\n", i+1, stat[0], stat[1], stat[2], stat[3], stat[4])
			}
		}
		fmt.Println("Run", run+1, "Max Clique:", maxClique)
		allstats[run] = stats
	}

	elapsed := time.Since(start)
	fmt.Printf("AOC took %s\n", elapsed)

	if getstats {
		stats := make([][]float64, generations)
		for gen := 0; gen < generations; gen++ {
			stats[gen] = make([]float64, 5)
			for stat := 0; stat < 5; stat++ {
				for run = 0; run < runqnt; run++ {
					stats[gen][stat] += allstats[run][gen][stat]
				}
				stats[gen][stat] /= float64(runqnt)
			}
		}
		err := dataset.Write(statsfile, stats, elapsed)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func initializeFlags() {
	flag.IntVar(&ants, "ants", 100, "number of ants")
	flag.IntVar(&generations, "gens", 10, "number of generations to run")
	flag.Float64Var(&minPheromone, "minpheromone", 0.1, "min value the pheromone of a vertex can have")
	flag.Float64Var(&maxPheromone, "maxpheromone", 5, "max value the pheromone of a vertex can have")
	flag.Float64Var(&evaporationRate, "evaporation", 0.2, "pheromone evaporation rate, must be (0.0, 1.0)")
	flag.StringVar(&file, "file", "easy", "file containing data to be processed")
	flag.Int64Var(&seed, "seed", 1, "seed for generating the initial population")
	flag.BoolVar(&verbose, "verbose", false, "Prints stats per generation")
	flag.BoolVar(&getstats, "getstats", false, "generate stats and saves into an outpu file on ./analysis/")
	flag.Parse()
}

// setSeed sets the given number as seed, or a random value if seed is <= 0
func setSeed(seed int64) int64 {
	if seed <= 0 {
		max := big.NewInt(2<<31 - 1)
		rseed, _ := crand.Int(crand.Reader, max)
		seed = rseed.Int64()
	}
	rand.Seed(seed)
	return seed
}
