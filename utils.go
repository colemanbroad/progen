package main

import (
	"sort"

	"golang.org/x/exp/constraints"
)

const (
	CLEAR   = "\033[K"
	UPSTART = "\033[F"
	UP      = "\033[A"
	START   = "\r"
)

func assert(b bool, msg string) {
	if !b {
		ErrorLog.Fatal(msg)
	}
}

func sortedKeys[K constraints.Ordered, V any](m map[K]V) []K {
	keys := make([]K, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func mean(m []float64) float64 {
	n, sum := 0, 0.0
	for _, mi := range m {
		n += 1
		sum += mi
	}
	return sum / float64(n)
}
