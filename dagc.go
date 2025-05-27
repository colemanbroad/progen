package main

import (
	"fmt"
)

// We want to count the number of ways of making any type available in the catalog.
// This is a step along the way to counting the number of ways to make a program length n,
// which is a step along the way to iterating over programs in a nice even way.

// First, we iterate over all the fragments and add one to the type and track the size of
// the dag required to build it.
// Then we iterate again and new types have nonzero possible ways to build,
// but the number of ways is the size of the input space i.e. tuple space,
// but the size of the subgraphs may differ for each of those ways!
// 		We could keep track of distributions???
// Then we

type MyType uint32

const (
	A MyType = iota
	B
	C
	D
)

type FnT struct {
	atypes []MyType
	rtype  MyType
}

var (
	f = FnT{[]MyType{}, A}
	g = FnT{[]MyType{}, B}
	h = FnT{[]MyType{A, B}, C}
	i = FnT{[]MyType{}, A}
	j = FnT{[]MyType{A}, A}
	k = FnT{[]MyType{A, C}, A}
	l = FnT{[]MyType{A, C}, D}
)

type Cata map[string]FnT

// A smarter way might be to just map between names.
type CataGraph2 struct {
	in  map[MyType][]string // funcs producing T
	out map[MyType][]string // funcs requiring T
}

func buildTypeCatalog(catalog Cata) CataGraph2 {
	cg := CataGraph2{
		in:  map[MyType][]string{},
		out: map[MyType][]string{},
	}
	for name, f := range catalog {
		funcs, _ := cg.in[f.rtype]
		funcs = append(funcs, name)
		cg.in[f.rtype] = funcs
		for _, a := range f.atypes {
			funcs, _ := cg.out[a]
			funcs = append(funcs, name)
			cg.out[a] = funcs
		}
	}
	return cg
}

func main() {
	catalog := Cata{
		"f": f, "g": g, "h": h,
		"i": i, "j": j, "k": k,
		"l": l,
	}
	fmt.Println(f, g, h)
	fmt.Println(catalog)

	cata2 := buildTypeCatalog(catalog)
	fmt.Printf("cg := %+v \n", cata2)

	determineLevel(catalog, cata2)

}

// We can use the TF-Graph to build an index of types and their trasitive requirements.
//
// We can assign a level to every F and T, which begins at zero for F : () -> T foreach T,
// and zero foreach T produced by each F. Then it increments to one for F : T -> T2 that require only
// T : lvl=0 and then T produced by those F *and not already a lower level*.
func determineLevel(catalog Cata, cat2 CataGraph2) {
	lvlF := map[string]int{}
	lvlT := map[MyType]int{}

	tset := NewSet[MyType]()
	fset := NewSet[string]()
	level := 0

	for fset.Size() < len(catalog) {
		addedOne := false
		starting_set := NewSetFromSlice(tset.Elements())
		// TODO: we shouldn't need to iterate over previously added Fns
		for name, f := range catalog {
			inputset := NewSetFromSlice(f.atypes)
			if inputset.Difference(starting_set).Size() != 0 {
				continue
			}
			// buildable. but does it already exist?
			_, ok := lvlF[name]
			if ok {
				continue
			}
			// it's buildable and hasn't been added yet
			addedOne = true
			lvlF[name] = level
			fset.Add(name)
			// has the rtype been added ?
			_, ok = lvlT[f.rtype]
			if ok {
				continue
			}
			// the rtype hasn't been added yet
			lvlT[f.rtype] = level
			tset.Add(f.rtype)
		}
		if addedOne == false && fset.Size() < len(catalog) {
			panic("Our catalog isn't buildable!")
		}
		level += 1
	}

	fmt.Printf("%+v \n", lvlF)
	fmt.Printf("%+v \n", lvlT)
}

// OK, now we can determine the level for any F or T, as well if the Catalog is buildable.
// But we want to *count* the ways that an F or T could be built!  An F can be run in a
// number of ways equal to the size of the input space (number of possible input tuple
// values). A T can be built in a number of ways equal to the sum of this number across
// input funcs.

// When counting it's very easy to construct a catalog with types that can be built in
// infinite ways. Any time we have an F : (..., T) -> T we have inf ways of building T.
// And this inf pollutes all downstream types i.e. F2 : T -> T2 has inf ways of building
// T2 because it inherits from T. This inf occurrs whenever we have a loop in the graph,
// even if it spans multiple F, e.g. T1 -> F1 -> T2 -> F2 -> T1.
//
// How should we deal with these inifinities?
// Doesn't this remind you of Feynman Diagrams?
//
// Also, I think this particular kind of factor graph should be called an FT-Graph.
//
// One way to deal with these infs is to limit the integral to DAGs of a certain size.
//
// We can always generate programs of infinite size simply by using multiple disconnected
// components, or by repeating
//
// Counting the product space of inputs to F is not enough! This models every input as independent, but we can use the same input in multiple arguments of the same type! Is this sufficient to
// count all the dags? What about correlations between even more distant objects, things that aren't args to the same F ? Have we undercounted by ignoring the DAGs with lots of shared reuse of Ts?
// What if, instead, we counted DAGs with each node labeled F, and then removed the dags that violated our type constraints?
//
// What about counting higher-kinded-types and generics in these graphs?
// Maybe it's easier to start with a Partial Order / Relation
//
// The basic technique of building programs directly by sampling random fragments with random connections tends to dramatically oversample DAGs that have a lot of F where the order doesn't matter. E.g. the program `A(); B(); C(); A(); C();` has one DAG but 5! programs that could be sampled. While `A(B(C(A(C))))` has one DAG and one program.

// t0 counts the number of ways we can make a type (ignoring program size).
// Many types have infinite way of being constructed!
// This happens whenever there is a loop in the TypeGraph upstream of the desired type T.
// How do you even *have* a loop in the TypeGraph? The TypeGraph must have factor graph
// structure, s.t. there is bipartite separation between types and funcs. So a loop exists
// when F : (... , T) -> T, which forms a loop T -> F -> T in the TypeGraph.
// So the DAG is a factor-graph-DAG. And we want to count/enumerate these factor-graph-dags.
//
// We should make this factor graph dag explicit?
