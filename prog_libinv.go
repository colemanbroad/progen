package main

import "fmt"

type FuncSet = *Set[Sym]

type LibraryInverse struct {
	provides map[Type]FuncSet
	requires map[Type]FuncSet
}

func (lib LibraryInverse) String() string {
	str := "provides:\n"
	for name, fnset := range lib.provides {
		str += fmt.Sprintf("type: %v, fn_set:\n%v\n", name, fnset)
	}
	str += "requires:\n"
	for name, fnset := range lib.requires {
		str += fmt.Sprintf("type: %v, fn_set:\n%v\n", name, fnset)
	}
	return str
}

// implement bredth first search on the Type Graph to find the shortest path between t1 and t2.
// construct layered DAG of types and funcs
func buildTypeGraph(type_set *Set[Type]) {

	// type_set := NewSet[Type]()

	all_funcs := NewSet[string]()
	all_types := NewSet[Type]()
	lvl_funcs := make(map[int]Set[string], 0)
	lvl_types := make(map[int]Set[Type], 0)

	lvl := 0

	all_types = all_types.Union(type_set)
	lvl_types[lvl] = *type_set

	// for every type in the set figure out what funcs require it and then if those funcs have all
	// for every func in the library figure out if it's set of argtypes << set of available types.
	// if yes, then add it to _producible_
	// if it's rtype isn't in typeset add it's rtype to typeset and mark the typeset as modified this iteration.
	// if the typeset isn't modified after a full iteration through the Lib then we're DONE
	for {
		type_set = NewSet[Type]()
		func_set := NewSet[string]()

		for _, fn := range Library {

			if all_funcs.Contains(fn.name) {
				continue
			}

			paramset := NewSet[Type]()
			for _, ptype := range fn.ptypes {
				paramset.Add(ptype)
			}

			if paramset.Difference(all_types).Size() == 0 {
				func_set.Add(fn.name)
				type_set.Add(fn.rtype)
			}
		}

		// fmt.Println("The size is ", func_set.Size())
		// fmt.Println("all_funcs is: ", all_funcs)

		if func_set.Size() == 0 {
			break
		}

		all_funcs = all_funcs.Union(func_set)
		l := lvl_funcs[lvl]
		l = *l.Union(func_set)
		lvl_funcs[lvl] = l

		lvl += 1

		all_types = all_types.Union(type_set)
		r := lvl_types[lvl]
		r = *r.Union(type_set)
		lvl_types[lvl] = r
	}

	InfoLog.Println("The TypeGraph is\n", lvl_types)
	InfoLog.Println("The FuncGraph is\n", lvl_funcs)
}

// TODO: Build TypeGraph data structure that makes determining TypePath easy
// Find a short program that includes the target function.

func buildLibraryInverse() LibraryInverse {
	// In order to build `fn` we have to first build it's input types
	// This requires finding funcs that produce that type.
	// For each func we try adding it and hope we don't get stuck, but when we do get stuck we backtrack...
	// It's like a pathfinding algorithms. It's a sequence of ops that have to be strung together in the right way.
	// Let's assume we don't care about wiring for now. The only thing we need is to have a program that includes `fn` and
	// also has an empty set of missing syms.

	// First let's compute the graph of Types and how they depend on each other.
	// A TypeGraph is a DirectedGraph of Types and Functions.
	// We have a Library which is indexed by function and tells us the argtypes and rtype
	// We want to have something indexed by Type and tells us which functions require that type or provide that type.
	// We can then make a Path through

	lib := LibraryInverse{
		provides: make(map[Type]FuncSet),
		requires: make(map[Type]FuncSet),
	}

	for _, fn := range Library {
		fset, ok := lib.provides[fn.rtype]
		if !ok {
			fset = NewSet[Sym]()
			// fset = make(FuncSet3, 0)
			// fset = NewSet[*FuncDefn]()
		}

		fset.Add(Sym(fn.name))
		// fset[&fn] = struct{}{}
		lib.provides[fn.rtype] = fset
	}

	for _, fn := range Library {
		for _, ptype := range fn.ptypes {
			fset, ok := lib.requires[ptype]
			if !ok {
				// fset = make(FuncSet3, 0)
				// fset = FuncSet(*NewSet[*FnCall]())
				fset = NewSet[Sym]()
			}
			// fset[&fn] = struct{}{}
			fset.Add(Sym(fn.name))
			lib.requires[ptype] = fset
		}
	}

	return lib
}
