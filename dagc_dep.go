


func mapit(catalog Cata) {
	t0 := map[MyType]int{}
	for _, fn := range catalog {
		if len(fn.atypes) == 0 {
			fmt.Println(fn)
			cnt, ok := t0[fn.rtype]
			if !ok {
				cnt = 0
			}
			t0[fn.rtype] = cnt + 1
		}
	}

	ts := NewSetFromMapKeys(t0)
	for _, fn := range catalog {
		ltypes := NewSetFromSlice(fn.atypes)
		if ltypes.Difference(ts).Size() == 0 {
			fmt.Println(fn)
			cnt, ok := t0[fn.rtype]
			if !ok {
				cnt = 0
			}
			t0[fn.rtype] = cnt + 1
		}
	}
	fmt.Println(ts)

}

type TNode struct {
	name string
	T    MyType
	out  []*FNode
	in   []*FNode
}

type FNode struct {
	name string
	fn   FnT
	out  *TNode
	in   []*TNode
}

// The actual T,F Nodes live here, and they contain pointers to each other.
type CataGraph struct {
	ts map[MyType]TNode
	fs map[string]FNode // name of function
}

// func newCataGraph(catalog []FnT) CataGraph {
// 	cg := CataGraph{}
// 	names := []string{"f", "g", "h"}
// 	for i, fn := range catalog {
// 		fnode := FNode{
// 			name: names[i],
// 			fn:   f,
// 			out:  &TNode{},
// 			in:   []*TNode{},
// 		}
// 		for _, a := range fn.atypes {
// 			tnode, ok := cg.ts[a]
// 			if !ok {
// 				tnode = TNode{
// 					name: string(a),
// 					T:    a,
// 					out:  []FNode{},
// 					in:   []FNode{},
// 				}
// 			}
// 		}
// 		cg.ts
// 	}
// }
