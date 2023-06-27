package aco

import "math/rand"

type Graph map[string]map[string]float64

func (g Graph) GetNeighbours(key string) []string {
    neighbours := []string{}
    for k := range g[key] {
        neighbours = append(neighbours, k)
    }
    return neighbours
}

func (g Graph) GetRandomKey() string {
    keys := []string{}
    for k := range g {
        keys = append(keys, k)
    }
    val := rand.Intn(len(keys))
    return keys[val]
}
