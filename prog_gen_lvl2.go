package main

import "maps"

func cloneLib(lib Library) Library {
	l2 := NewLib()
	l2.fns = maps.Clone(lib.fns)
	l2.vals = maps.Clone(lib.vals)
	return l2
}

func sample2lvl() Program {
	value_library = make(map[Sym]Value)
	var p Program
	p0 := Program{}
	sp := newSampleParams()

	// fn_library = make(map[Sym]Fun)
	lib := NewLib()
	lib.addBasicMathLib()
	sp.Program_length = 3
	p = lib.sampleProgram(sp)
	p0 = append(p0, p...)

	sp.Program_length = 10
	lib2 := cloneLib(lib)
	delete(lib2.fns, "one")
	p = lib2.sampleProgram(sp)
	p0 = append(p0, p...)

	return p
}
