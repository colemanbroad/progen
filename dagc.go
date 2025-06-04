package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"math/rand/v2"
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

	determineLevel(catalog)
	type_index := buildTypeIndex(catalog)
	fmt.Printf("type_index (map[u32][]string) = %+v \n", type_index)
}

// We can use the TF-Graph to build an index of types and their trasitive requirements.
//
// We can assign a level to every F and T, which begins at zero for F : () -> T foreach T,
// and zero foreach T produced by each F. Then it increments to one for F : T -> T2 that require only
// T : lvl=0 and then T produced by those F *and not already a lower level*.
func determineLevel(catalog Cata) {
	lvlF := map[string]int{}
	lvlT := map[MyType]int{}

	tset := NewSet[MyType]()
	fset := NewSet[string]()
	level := 0

	// c

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
func buildTypeIndex(catalog Cata) map[uint32][]string {
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

// Set of Nodes represented by index in DataFlow.nodes
// type Nodeset *Set[uint]
type Level = uint16
type Node struct {
	fn        Fun
	level     Level
	arg_nodes []uint
}

// DataFlow describes the complete flow of data in the program.
// Each node contains reference to previous nodes
type DataFlow struct {
	nodes     []Node
	max_level uint
}

func (d DataFlow) getTypeSetWhereLvl(min_level, max_level Level) *Set[Type] {
	s := NewSet[Type]()
	for _, n := range d.nodes {
		b0 := n.level >= min_level
		b1 := n.level <= max_level
		if b0 && b1 {
			s.Add(n.fn.rtype)
		}
	}
	return s
}

func (d DataFlow) findArgNodesForFn(fn Fun, lvl Level) []uint {
	s := []uint{}
	for _, a := range fn.ptypes {
		n1 := d.getNodesOfTypeWhereLevel(a, lvl-1, lvl-1)
		n2 := d.getNodesOfTypeWhereLevel(a, 0, lvl-1)
		n3 := n1.Union(n2)
		if n3.Size() == 0 {
			panic("unreachable")
		}
		ret, err := n3.Sample()
		if err != nil {
			panic("impossible")
		}
		s = append(s, uint(ret))
	}
	return s
}

func (d DataFlow) getNodesOfTypeWhereLevel(typ Type, min_level, max_level Level) *Set[int] {
	s := NewSet[int]()
	for i, n := range d.nodes {
		b0 := n.level >= min_level
		b1 := n.level <= max_level
		b2 := n.fn.rtype == typ
		if b0 && b1 && b2 {
			s.Add(i)
		}
	}
	return s
}

func (d DataFlow) getBuildableTypes(catalog Cata, reqd, optional *Set[Type]) *Set[uint32] {
	s := NewSet[uint32]()
	for _, fun := range catalog {
		a := NewSetFromSlice(fun.ptypes)
		if a.Intersection(reqd).Size() < reqd.Size() {
			continue
		}
		if a.Difference(optional).Size() > 0 {
			continue
		}
		s.Add(fun.hashOf())
	}
	return s
}

// We need to keep track of which nodes are at the fronteir of unevald nodes.
// This requires a simple set of fronteir nodes, and a way of figuing out which
// new nodes to add to the frontier after a frontier node has been added to the program.
//
// Initially, if a node has no deps we add it to the frontier.
// Then when we remove a node from the frontier we check it's deps
// against the set of nodes already in the program. If the set contains,
// then we add it to the frontier (greedily). The frontier set can dedup.
// This alg is O(program_length ^ 2), because we check the full flow after every new line.
//
// A faster version would just add all the level n nodes (in any order) before adding lvl n+1,
// but this is only a subset of linearizations.
func (d DataFlow) linearize() Program {
	prog := Program{}
	// gen := GenSym{
	// 	idx: 0,
	// }
	// first, let's just make a program in the node order.
	for i, n := range d.nodes {
		argsyms := []Sym{}
		for _, a := range n.arg_nodes {
			argsyms = append(argsyms, Sym(fmt.Sprintf("v%v", a)))
		}
		stmt := Statement{
			fn:      n.fn,
			outsym:  Sym(fmt.Sprintf("v%v", i)),
			argsyms: argsyms,
		}
		prog = append(prog, stmt)
		// for _, a := range n.arg_nodes {
		// 	args := []uint{}
		// 	// args = append(args, a, i)
		// 	// node_inv = append(node_inv)
		// }
	}
	return prog
}

// Add some terms to the current level.
// Terms are chosen from the set of possible ones: h_i^*
// This set is too large to build explicitly, so we're going to break it down and sample it instead.

// Once we know that a fragment is ACTUALLY, properly buildable,
// then we go about sampling args (node indices of appropriate type).
// Then we can add it!
// t0 := flow.getRTypeSet(min_level=0, max_level=level-1)
// t1 := flow.getRTypeSet(min_level=level-1, max_level=level-1)
// fnt2 := flow.getBuildableFnTypes(catalog, necessary=t1, optional=t1)
// t_chosen := multiSampleN(set=fnt2,n=5) // n may be greater than len(set)
// nodes := chooseWires(t_chosen, necessary=t0, optional=t1)
func NewDataFlow(catalog Cata) DataFlow {
	flow := DataFlow{
		nodes:     []Node{},
		max_level: 0,
	}
	type2fns := buildTypeIndex(catalog)

	maxlevel := Level(3)
	level := Level(0)
	for level <= maxlevel {
		n_terms := 3
		t0 := flow.getTypeSetWhereLvl(level-1, level-1)
		t1 := flow.getTypeSetWhereLvl(0, level-1)
		t2 := flow.getBuildableTypes(catalog, t0, t1)
		for range n_terms {
			typ, _ := t2.Sample()
			catFns := type2fns[typ]
			fn := catalog[catFns[rand.IntN(len(catFns))]]
			arg_nodes := flow.findArgNodesForFn(fn, level)
			node := Node{
				fn:        fn,
				level:     level,
				arg_nodes: arg_nodes,
			}
			flow.nodes = append(flow.nodes, node)
		}
		level += 1
	}

	return flow
}

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

	df := NewDataFlow(catalog)
	fmt.Printf("%+v\n", df)

	prog := df.linearize()
	printProgram(prog, Fmt)
	// fmt.Printf("%+v\n", prog)

}
