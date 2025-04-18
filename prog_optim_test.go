package main

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

func TestReshuffle(t *testing.T) {
	initPeanoLibrary()
	prog := sampleProgram(newSampleParams())
	prog_mutated := reshuffle(prog)
	// prog_mutated := point_mutate(prog)
	// assert that sometimes prog != prog_mutated?
	tassert(len(prog) == len(prog_mutated), func() { t.Error("Mutated Programs not same length.") })

	s1 := NewSetFromSlice(getCallSyms(prog))
	s2 := NewSetFromSlice(getCallSyms(prog_mutated))
	tassert(s1.Difference(s2).Size() == 0, func() { t.Error("Reshuffling should't change (multi)set of FnCalls") })

	// p1, p2 := formatProgram(prog), formatProgram(prog_mutated)
	// log.Println("TestMutate", p1, p2)

	// for i := range prog {
	// 	log.Println(prog[i])
	// 	log.Println(prog_mutated[i])
	// }
}

// func Fuzz_point_mutate(f *testing.F) {
// 	for i := range 50 {
// 		f.Add(i)
// 	}
// 	f.Fuzz(t_point_mutate)
// }

func initPeanoLibrary() {
	fn_library = make(map[Sym]Fun)
	addPeanoLib()
}

func TestSampleProgramLong(t *testing.T) {
	initPeanoLibrary()
	for range 1000 {
		p := sampleProgram(newSampleParams())
		validateOrFail(p, "sampled direct")
	}
}

func TestReshuffleLong(t *testing.T) {
	initPeanoLibrary()
	for range 1000 {
		p := sampleProgram(newSampleParams())
		p = reshuffle(p)
		validateOrFail(p, "reshuffled")
	}
}

func TestPointMutateLong(t *testing.T) {
	initPeanoLibrary()
	for range 1000 {
		p := sampleProgram(newSampleParams())
		p, _ = pointMutate(p)
		validateOrFail(p, "mutated")
	}
}

func TestBasicgenRewireLong(t *testing.T) {
	initPeanoLibrary()
	hadSuccess, hadFailure := false, false // tracks whether we achieved True and False returns (sometimes each)
	for range 1000 {
		p := sampleProgram(newSampleParams())
		p, isSuccess := rewireBase(p)
		hadSuccess = hadSuccess || isSuccess
		hadFailure = hadFailure || !isSuccess
		validateOrFail(p, "rewired")
	}
	if !hadSuccess {
		t.Error("never succeeded in rewiring")
	}
	if !hadFailure {
		t.Error("never failed in rewiring")
	}
}

func TestBasicgenComboLong(t *testing.T) {
	initPeanoLibrary()
	p := sampleProgram(newSampleParams())
	for range 1000 {
		x := rand.Float32()
		switch {
		case x < 0.3:
			p = reshuffle(p)
			panicIfInvalid(p)
		case x < 0.6:
			p, _ = pointMutate(p)
			panicIfInvalid(p)
		case x < 0.9:
			p, _ = rewireBase(p)
			panicIfInvalid(p)
			// p = rewire(p)
		default:
		}
		panicIfInvalid(p)
	}
}

func TestDeltaD(t *testing.T) {
	input := make([]int, 2500)
	input[0] = 1
	input[1000] = 2
	input[2000] = 3
	input[2499] = 5
	test := func(s []int) bool {
		b1 := slices.Contains(s, 1)
		b2 := slices.Contains(s, 2)
		b3 := slices.Contains(s, 3)
		b4 := slices.Contains(s, 5)
		return b1 && b2 && b3 && b4
	}
	target := []int{1, 2, 3, 5}
	result := deltaD(input, test)
	if !slices.Equal(result, target) {
		t.Error("DeltaD failed on target: ", target)
	} else {
		// fmt.Println("We found a minimal sequence and it is....")
		// fmt.Println(result)
	}
}

func TestDeltaDProgram(t *testing.T) {
	sp := newSampleParams()
	sp.Program_length = 100
	input := sampleProgram(sp)
	// why []Statement and not Program is loadbearing?
	test := func(p []Statement) bool {
		vm, _ := evalProgram(p)
		b1 := false
		for _, v := range vm {
			if v.value.(int) == 73 {
				b1 = true
			}
		}
		return b1
	}
	for !test(input) {
		input = sampleProgram(sp)
	}
	vm, _ := evalProgram(input)
	fmt.Println("The starter program:")
	printProgramAndValues(input, vm)
	result := deltaD(input, test)
	fmt.Println("The resulting program:")
	vm, _ = evalProgram(result)
	printProgramAndValues(result, vm)

	// if !slices.Equal(result, target) {
	// 	t.Error("DeltaD failed on target: ", target)
	// }
}
