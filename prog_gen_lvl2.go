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
	sp := newSampleParams()

	lib := NewLib()
	lib.addBasicMathLib()
	sp.Program_length = 3
	p0 := lib.sampleProgram(sp)

	sp.Program_length = 50
	lib2 := cloneLib(lib)
	delete(lib2.fns, "one")
	sp.Prefix = p0
	// fmt.Printf("lib2 %+v", lib2)
	p1 := lib2.sampleProgram(sp)
	return p1
}
