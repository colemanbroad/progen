package main

import (
	"math/rand/v2"
	"testing"
)

func Test_reshuffle(t *testing.T) {
	prog := sampleProgram_fromFragmentLib(newSampleParams())
	prog_mutated := reshuffle(prog)
	// prog_mutated := point_mutate(prog)
	// assert that sometimes prog != prog_mutated?
	tassert(len(prog) == len(prog_mutated), func() { t.Error("Mutated Programs not same length.") })

	s1 := NewSetFromSlice(getCallSyms(prog))
	s2 := NewSetFromSlice(getCallSyms(prog_mutated))
	tassert(s1.Difference(s2).Size() == 0, func() { t.Error("Reshuffling should't change (multi)set of FnCalls") })

	// p1, p2 := formatProgram(prog), formatProgram(prog_mutated)
	// log.Println("Test_mutate", p1, p2)

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

func Test_basicgen(t *testing.T) {
	for range 1000 {
		p := sampleProgram_fromFragmentLib(newSampleParams())
		validateOrFail(p, "sampled direct")
	}
}

func Test_basicgen_shuffle(t *testing.T) {
	for range 10 {
		p := sampleProgram_fromFragmentLib(newSampleParams())
		p = reshuffle(p)
		validateOrFail(p, "reshuffled")
	}
}

func Test_basicgen_pointmut(t *testing.T) {
	for range 1000 {
		p := sampleProgram_fromFragmentLib(newSampleParams())
		p, _ = point_mutate(p)
		validateOrFail(p, "mutated")
	}
}

func Test_basicgen_rewire(t *testing.T) {
	hadSuccess, hadFailure := false, false // tracks whether we achieved True and False returns (sometimes each)
	for range 1000 {
		p := sampleProgram_fromFragmentLib(newSampleParams())
		p, isSuccess := rewire_base(p)
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

func Test_basicgen_combo(t *testing.T) {
	p := sampleProgram_fromFragmentLib(newSampleParams())
	for range 1000 {
		x := rand.Float32()
		switch {
		case x < 0.3:
			p = reshuffle(p)
			panic_if_invalid(p)
		case x < 0.6:
			p, _ = point_mutate(p)
			panic_if_invalid(p)
		case x < 0.9:
			p, _ = rewire_base(p)
			panic_if_invalid(p)
			// p = rewire(p)
		default:
		}
		panic_if_invalid(p)
	}
}
