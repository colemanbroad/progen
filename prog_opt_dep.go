package main

import "fmt"

// Find missing syms.
// Determine if they can be attached to existing syms.
// For those that can't create a new sym by sampling a Function from the Library and adding a Statement
// TODO
func fillInMissing(program UncheckedProgram) Program {
	analysis := analyzeProgram(program)
	if len(analysis.missing.syms) == 0 {
		return Program(program)
	}

	// TODO: FIXME: This doesn't actually fill in missing syms.
	fmt.Println("Uh Oh there's gonna be trouble.")
	return Program(program)
}

type AnalyzedProgram struct {
	present, missing Catalog
}

// Determine which Syms are defined and which are missing.
func analyzeProgram(program UncheckedProgram) AnalyzedProgram {
	present := NewCatalog()
	missing := NewCatalog()

	// loop over existing program syms in random order
	for i, stmt := range program {
		// fmt.Println("add() #1 in analyzeProgram()")
		present.add(stmt.outsym, stmt.fn.rtype, uint16(i))
		for _, arg := range stmt.argsyms {
			symline := SymLine{sym: arg, line: uint16(i)}
			typ, ok := present.syms[symline]
			if !ok {
				// fmt.Println("add() #2 in analyzeProgram()")
				ErrorLog.Printf("missing arg: %v with type: %v on line: %v \n", arg, typ, i)
				missing.add(arg, typ, uint16(i))
			}
		}
	}

	return AnalyzedProgram{present, missing}
}

// Check that a program is valid i.e. all Syms are defined top-to-bottom.
func validate[T UncheckedProgram | Program](program T) (Program, error) {
	// present := make(SymSet)
	present := NewSet[Sym]()

	for i, stmt := range program {
		for j, arg := range stmt.argsyms {
			desired_type := stmt.fn.ptypes[j]
			if !present.Contains(arg) {
				ErrorLog.Printf("missing arg: %v with type: %v on line: %v \n", arg, desired_type, i)
				ErrorLog.Printf("the catalog is:\n %v \n", present)
				return Program{}, fmt.Errorf("Program is invalid:\n %v", formatProgram(program))
			}
		}
		present.Add(stmt.outsym)
	}

	return Program(program), nil

}

// A random subset of Missing Syms in p2 are attached to Syms created in p1.
// Sym names in p2 that overlap with those in p1 are changed.
// TODO:
func insert_concatPrograms(p1, p2 UncheckedProgram) UncheckedProgram {
	p2.uniquify_syms(p1)
	return append(p1, p2...)
}

func recombineUnchecked(a, b UncheckedProgram) (a2, b2 UncheckedProgram) {
	// randomly mix up the statements from p1 and p2 but return valid Programs.
	// a, b := UncheckedProgram(p1), UncheckedProgram(p2)
	b.uniquify_syms(a)
	n1, n2 := len(a), len(b)
	a2 = append(a[:n1/2], b[n2/2:]...)
	b2 = append(a[n1/2:], b[:n2/2]...)

	return a2, b2
}

func apply_all_changes(prog UncheckedProgram) UncheckedProgram {
	// shuffle(prog) // permute line order (may violate syms!) - maybe only allow valid permutations?
	// rewire(prog)  // maintain line order but change sym connections (maintain validity)
	// mutate(prog)  // replace a fn call with a new func?
	// replace_
	return UncheckedProgram{}
}

func _paramset(sym Sym) *Set[Type] {
	paramset := NewSet[Type]()
	for _, paramtype := range Library[sym].ptypes {
		paramset.Add(paramtype)
	}
	return paramset
}
