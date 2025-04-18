package main

import (
	"maps"
)

func cloneLib(lib Library) Library {
	l2 := NewLib()
	l2.fns = maps.Clone(lib.fns)
	l2.vals = maps.Clone(lib.vals)
	return l2
}

func sample2lvl() Program {
	value_library = make(map[Sym]Value)
	// var p Program
	// prog := Program{}
	sp := newSampleParams()

	lib := NewLib()
	lib.addBasicMathLib()
	sp.Program_length = 3
	p0 := lib.sampleProgram(sp)
	// prog = append(prog, p...)

	sp.Program_length = 10
	lib2 := cloneLib(lib)
	delete(lib2.fns, "one")
	sp.Prefix = p0
	// fmt.Printf("lib %+v", lib2)
	p1 := lib2.sampleProgram(sp)
	// prog = p
	// printProgram(p0, Fmt)
	// printProgram(p1, Fmt)
	// fmt.Println(p0)
	// fmt.Println(p1)
	// p0 = append(p0, p1[len(p0):]...)

	return p1
}
