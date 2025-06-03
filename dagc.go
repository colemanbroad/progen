package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
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

type MyType = Type

const (
	A MyType = "A"
	B MyType = "B"
	C MyType = "C"
	D MyType = "D"
)

func FnT(name string, ptypes []MyType, rtype MyType) Fun {
	return Fun{
		value:  nil,
		name:   name,
		ptypes: ptypes,
		rtype:  rtype,
	}
}

var (
	f = FnT("f", []MyType{}, A)
	g = FnT("g", []MyType{}, B)
	h = FnT("h", []MyType{A, B}, C)
	i = FnT("i", []MyType{}, A)
	j = FnT("j", []MyType{A}, A)
	k = FnT("k", []MyType{A, C}, A)
	l = FnT("l", []MyType{A, C}, D)
)

type Cata map[string]Fun

// A smarter way might be to just map between names.
type TypeCatalog struct {
	in  map[MyType][]string // funcs producing T
	out map[MyType][]string // funcs requiring T
}

// Inverts a catalog to find fns which provide/require a given type.
func buildTypeCatalog(catalog Cata) TypeCatalog {
	cg := TypeCatalog{
		in:  map[MyType][]string{},
		out: map[MyType][]string{},
	}
	for name, f := range catalog {
		funcs, _ := cg.in[f.rtype]
		funcs = append(funcs, name)
		cg.in[f.rtype] = funcs
		for _, a := range f.ptypes {
			funcs, _ := cg.out[a]
			funcs = append(funcs, name)
			cg.out[a] = funcs
		}
	}
	return cg
}

func test_dagc() {
	catalog := Cata{
		"f": f, "g": g, "h": h,
		"i": i, "j": j, "k": k,
		"l": l,
	}
	fmt.Println(f, g, h)
	fmt.Println(catalog)

	cata2 := buildTypeCatalog(catalog)
	fmt.Printf("cata2 (TypeCatalog) := %+v \n", cata2)

	determineLevel(catalog, cata2)
	type_index := buildTypeIndex(catalog, cata2)
	fmt.Printf("type_index (map[u32][]string) = %+v \n", type_index)
}

// We can use the TF-Graph to build an index of types and their trasitive requirements.
//
// We can assign a level to every F and T, which begins at zero for F : () -> T foreach T,
// and zero foreach T produced by each F. Then it increments to one for F : T -> T2 that require only
// T : lvl=0 and then T produced by those F *and not already a lower level*.
func determineLevel(catalog Cata, cat2 TypeCatalog) {
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
			inputset := NewSetFromSlice(f.ptypes)
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

	fmt.Printf("lvlF = %+v \n", lvlF)
	fmt.Printf("lvlT = %+v \n", lvlT)
}

func (m *Fun) ToBytes2() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := binary.Write(buf, binary.LittleEndian, m.rtype); err != nil {
		return nil, err
	}
	length := uint32(len(m.ptypes))
	if err := binary.Write(buf, binary.LittleEndian, length); err != nil {
		return nil, err
	}
	for _, val := range m.ptypes {
		if err := binary.Write(buf, binary.LittleEndian, val); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (m *Fun) ToBytes() []byte {
	if m == nil {
		panic("nope")
	}
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, m.rtype)
	binary.Write(buf, binary.LittleEndian, uint32(len(m.ptypes)))
	for _, val := range m.ptypes {
		binary.Write(buf, binary.LittleEndian, val)
	}
	return buf.Bytes()
}

func (m *Fun) hashOf() uint32 {
	by := m.ToBytes()
	return crc32.ChecksumIEEE(by)
}

// Build index from hash(Type) to []Fragment.name.
func buildTypeIndex(catalog Cata, cat2 TypeCatalog) map[uint32][]string {
	typecat := map[uint32][]string{}
	for name, fn := range catalog {
		l, ok := typecat[fn.hashOf()]
		if !ok {
			l = make([]string, 0, 10)
		}
		l = append(l, name)
		typecat[fn.hashOf()] = l
	}
	return typecat
}

// type Node struct {
// 	f Fun
// }
// type DataFlow struct {
// 	nodes []Statement
// }

func sampleDataflow() {

	zero := FnT("zero", []MyType{}, "int")
	zero.value = func() int { return 0 }
	one := FnT("one", []MyType{}, "int")
	one.value = func() int { return 1 }
	plus := FnT("plus", []MyType{"int", "int"}, "int")
	plus.value = func(a, b int) int { return a + b }

	catalog := Cata{
		"00": zero,
		"11": one,
		"++": plus,
	}

	tc := buildTypeCatalog(catalog)
	fmt.Printf("tc = %#v \n", tc)

	// The plan is to eventually
	// 0. pick n
	// 1. pick bold(h) = h0, h1, ..., hn
	// 2. pick m_i = h^*_i choose h_i  foreach i in [n]
	// 3. randomly linearize

	// Termset is insufficient, because we need to know WHO to attach to!
	// Not just Fun, But something like a full blown Statement!

	type Termset = *Set[*Fun]
	terms := map[uint16]map[Type]Termset{}

	// level 0 first

	catalog_used := NewSet[string]()
	available_types := NewSet[Type]()

	// This works for zero-arguments types as well!
	buildable := func(ptypes []Type) bool {
		for _, t := range ptypes {
			if !available_types.Contains(t) {
				return false
			}
		}
		return true
	}

	level := uint16(0)
	for catalog_used.Size() < len(catalog) {
		newtypes := NewSet[Type]()
		terms[level] = map[Type]*Set[*Fun]{}
		fmt.Printf("terms = %+v \n", terms)
		fmt.Println("size comparison : ", catalog_used.Size(), len(catalog))
		terms[level] = map[Type]*Set[*Fun]{}
		for name, f := range catalog {
			b0 := catalog_used.Contains(name)
			// b1 := len(f.ptypes) == 0
			b2 := buildable(f.ptypes)
			if b0 || !b2 {
				fmt.Println("f = ", f)
				fmt.Println("b0, b2 = ", b0, b2)
				continue
			}
			termset, ok := terms[level][f.rtype]
			if !ok {
				termset = NewSet[*Fun]()
				termset.Add(&f)
				terms[level][f.rtype] = termset
				fmt.Printf("terms = %+v \n", terms)
			} else {
				termset.Add(&f)
			}
			catalog_used.Add(name)
			newtypes.Add(f.rtype)
		}
		level += 1
		for _, t := range newtypes.Elements() {
			available_types.Add(t)
		}
	}

	fmt.Printf("terms = %+v \n", terms)

	// terms : (lvl:int, t:type) -> Set<Term>
	// How do we know what values are available in terms[i+1] given terms[0:i] ?
	// We do:
	//   terms[i+1] = union l[i+1, t] forall t in types
	// where
	//   l[i+1, t] = union values(l[:i], fn, i) forall fn in Catalog where fn.rtype = t
	// where
	//   values(terms, fn, i) = "all possible combinations of arguments to fn taken"
	//
	// OK, wait... We need to index the catalog by hash(fn.type) when doing type-first flow gen,
	// but when building the flow we need to index `terms` by fn.rtype.
	// Do we need to index `terms` by hash(fn.type) ?
	//

}

// Now that we have the index from Type -> Fragments we can start with building
// the more complex structure: the dataflow. How are we going to build it?
// Well we start by sampling a _type_ from the keys to typecat. Then we check all
// the valid syms of the correct types as args, and each time we check the FULL
// program DAG's hash value to see if it already exists (has already been created, even during
// an intermediate stage earlier in program gen). If so, then we skip it and move on. Even
// if the final structure had a totally different hash function.

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
