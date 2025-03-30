package main

import "testing"

func tassert(b bool, fn func()) {
	if !b {
		fn()
	}
}

func t_point_mutate(t *testing.T, seed int) {
	prog := sampleProgram_fromFragmentLib(newSampleParams())
	prog_mutated, isNew := point_mutate(prog)
	tassert(len(prog) == len(prog_mutated), func() { t.Error("Mutated Programs not same length.") })
	s1 := getCallSyms(prog)
	s2 := getCallSyms(prog_mutated)

	count_changes := 0
	for i := range s1 {
		if s1[i] != s2[i] {
			count_changes += 1
		}
	}
	if isNew {
		tassert(count_changes == 1, func() { t.Error("Incorrect mutation.") })
	} else {
		tassert(count_changes == 0, func() { t.Error("Incorrect mutation.") })
	}

	// p1, p2 := formatProgram(prog), formatProgram(prog_mutated)
	// log.Println("Test_mutate", p1, p2)

	// for i := range prog {
	// 	log.Println(prog[i])
	// 	log.Println(prog_mutated[i])
	// }
}
