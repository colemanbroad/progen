package main

import (
	"fmt"
	"maps"
	"math/rand/v2"
)

func cloneLib(lib Library) Library {
	l2 := NewLib()
	l2.fns = maps.Clone(lib.fns)
	l2.vals = maps.Clone(lib.vals)
	return l2
}

// Create a program from two calls to lib.sampleProgram() with
// different params and libraries.
func sample2lvl() Program {
	value_library = make(map[Sym]Value)
	sp := newSampleParams()

	lib := NewLib()
	lib.addBasicMathLib()
	sp.Program_length = 3
	p0 := lib.sampleProgram(sp)

	sp.Program_length = 50
	// lib2 := cloneLib(lib)
	delete(lib.fns, "one")
	sp.Prefix = p0
	// fmt.Printf("lib2 %+v", lib2)
	p1 := lib.sampleProgram(sp)
	return p1
}

func iterate() {
	N := 12 // legnth of array
	M := 10 // [0..M) domain of each array element

	digits := make([]int, N)
	for i := range N {
		digits[i] = 0
	}

	// find all primes coprime with M
	primes := PrimesUpTo(200)
	j := 0
	for i, p := range primes {
		if p%M == 0 || M%p == 0 {
			continue
		}
		primes[j] = primes[i]
		j += 1
	}
	primes = primes[:j]
	fmt.Println(primes)

	index := make([]int, N)
	for i := range N {
		index[i] = 0
	}
	for t := range 50 {
		for j := range N {
			// digits[j] = (digits[j] + primes[j]) % M // doesn't work. every 10th is 0,0,0,0,... then repeats.
			digits[j] = (digits[j]*primes[j] + 101) % M // doesn't work. each digit is on it's own cycle hitting a subset of 0..9 in a cycle.
			// We want every digit to see every value in some order (a permutation). But we want them all to use a
			// different permutation for different digits. And we don't want the permutations to repeat! We want
			// each digit to use a different permutation at t0 and to visit all possible permutations in a different
			// order! Different permutations of the sequence of possible permutations!
		}
		fmt.Println(t, digits)
	}
}

// PrimesUpTo returns a slice of all prime numbers ≤ n
func PrimesUpTo(n int) []int {
	if n < 2 {
		return []int{}
	}
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}
	for p := 2; p*p <= n; p++ {
		if isPrime[p] {
			for multiple := p * p; multiple <= n; multiple += p {
				isPrime[multiple] = false
			}
		}
	}
	primes := []int{}
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

// RandomPermutation returns a random permutation of 0 through n-1
func RandomPermutation(n int) []int {
	perm := make([]int, n)
	for i := range n {
		perm[i] = i
	}
	// rand.Seed(time.Now().UnixNano())
	for i := range n {
		j := rand.IntN(n-i) + i
		perm[i], perm[j] = perm[j], perm[i]
	}
	return perm
}
