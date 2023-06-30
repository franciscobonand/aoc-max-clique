package aco

import (
	"math"
	"math/rand"
	"sort"
	"strings"

	set "github.com/hashicorp/go-set"
	sts "github.com/montanaflynn/stats"
	"gonum.org/v1/gonum/floats"
)

type colony struct {
	ants            int
	generations     int
	minPheromone    float64
	maxPheromone    float64
	evaporationRate float64
	bestClique      []string
	graph           Graph
}

func NewColony(ants, iterations int, minPheromone, maxPheromone, evaporationRate float64, graph map[string]map[string]float64) *colony {
	return &colony{
		ants:            ants,
		generations:     iterations,
		minPheromone:    minPheromone,
		maxPheromone:    maxPheromone,
		evaporationRate: evaporationRate,
		graph:           graph,
	}
}

// Run returns the max clique and the stats for the run
func (c colony) Run() (int, [][]float64) {
	var mean, worst, repeated, sdev float64
	stats := make([][]float64, c.generations)
	for i := 0; i < c.generations; i++ {
		cliques := make([][]string, 0)
		for ant := 0; ant < c.ants; ant++ {
			clique := c.buildClique()
			sort.Strings((clique))
			cliques = append(cliques, clique)
		}
		c.updatePheromones(cliques)
		mean, worst, repeated, sdev = c.getStats(cliques)
		stats[i] = []float64{float64(len(c.bestClique)), worst, mean, repeated, sdev}
	}
	return len(c.bestClique), stats
}

func (c colony) getStats(cliques [][]string) (float64, float64, float64, float64) {
	worst := c.bestClique
	var cliqueStrs []string
	var sizes []float64
	for _, clique := range cliques {
		sizes = append(sizes, float64(len(clique)))
		if len(clique) < len(worst) {
			worst = clique
		}
		cliqueStrs = append(cliqueStrs, strings.Join(clique, ","))
	}
	unicliques := set.From[string](cliqueStrs)
	mean, _ := sts.Mean(sizes)
	sdev, _ := sts.StandardDeviation(sizes)
	return mean, float64(len(worst)), float64(c.ants - unicliques.Size()), sdev
}

// updatePheromones uses elitism (best solution is used to update the pheromones)
func (c *colony) updatePheromones(cliques [][]string) {
	bestClique := []string{}
	for _, clique := range cliques {
		if len(clique) > len(bestClique) {
			bestClique = clique
		}
	}

	if len(bestClique) > len(c.bestClique) {
		c.bestClique = bestClique
	}
	// Evaporate pheromone on all edges
	for k1, vertices := range c.graph {
		for k2, pheromone := range vertices {
			newVal := pheromone - (pheromone * c.evaporationRate)
			if newVal < c.minPheromone {
				newVal = c.minPheromone
			}
			c.graph[k1][k2] = newVal
		}
	}
	// Deposit pheromones for best ant
	bestDiff := float64(len(c.bestClique) - len(bestClique))
	for _, vertex1 := range bestClique {
		for _, vertex2 := range bestClique {
			if vertex1 == vertex2 {
				continue
			}
			newVal := c.graph[vertex1][vertex2] + (1 / (1 + bestDiff))
			if newVal > c.maxPheromone {
				newVal = c.maxPheromone
			}
			c.graph[vertex1][vertex2] = newVal
		}
	}
}

func (c colony) buildClique() []string {
	initialVertex := c.graph.GetRandomKey()
	clique := set.From([]string{initialVertex})
	candidates := set.From(c.graph.GetNeighbours(initialVertex))
	for {
		if candidates.Empty() {
			break
		}
		candStr := candidates.Slice()
		cliqueStr := clique.Slice()
		pheromoneFactors := c.getPheromoneFactors(candStr, cliqueStr)
		nextVertex := getRandomNode(candStr, pheromoneFactors)
		clique.Insert(nextVertex)
		nextNeighbours := set.From(c.graph.GetNeighbours(nextVertex))
		candidates = candidates.Intersect(nextNeighbours)
	}
	return clique.Slice()
}

func (c colony) getPheromoneFactors(vertices, clique []string) []float64 {
	factors := []float64{}
	for _, vertex := range vertices {
		sum := 0.0
		for _, cvert := range clique {
			sum += c.graph[vertex][cvert]
		}
		factors = append(factors, math.Pow(sum, 2))
	}
	return factors
}

func getChoiceProbability(factors []float64) (float64, []float64) {
	probabilities := []float64{}
	sum := 0.0
	for _, factor := range factors {
		sum += factor
	}
	for _, factor := range factors {
		probabilities = append(probabilities, factor/sum)
	}
	return sum, probabilities
}

func weightedChoice(nodes []string, weights []float64) string {
	cdf := make([]float64, len(weights))
	floats.CumSum(cdf, weights)
	val := rand.Float64() * cdf[len(cdf)-1]
	idx := sort.Search(len(cdf), func(i int) bool { return cdf[i] > val })
	return nodes[idx]
}

func getRandomNode(nodes []string, factors []float64) string {
	probSum, pheromProbs := getChoiceProbability(factors)
	val := rand.Float64() * probSum
	for idx := range nodes {
		val -= pheromProbs[idx]
		if val <= 0 {
			return nodes[idx]
		}
	}
	return nodes[len(nodes)-1]
}
