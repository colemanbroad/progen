package main

import "fmt"

type DepthStats struct {
	depth2values map[int]*Set[int]
	depthcount   map[int]int
}

func (ds DepthStats) update(prog Program, vals ValueMap) {
	dm := createDepthmap(prog)
	for s, val := range vals {
		depth := dm[s]
		d, ok := ds.depthcount[depth]
		if !ok {
			d = 0
		}
		ds.depthcount[depth] = d + 1
		valset, ok := ds.depth2values[depth]
		if !ok {
			valset = NewSet[int]()
			ds.depth2values[depth] = valset
		}
		intval, ok := val.value.(int)
		if ok {
			valset.Add(intval)
		}
	}
}

func (ds DepthStats) print() {
	keys := sortedKeys(ds.depth2values)
	fmt.Println("Depth\tUnique\tTotal\tRatio")
	for _, depth := range keys {
		depthset := ds.depth2values[depth]
		total := ds.depthcount[depth]
		fmt.Printf("%v\t%v\t%v\t%-8f \n", depth, depthset.Size(), total, float32(depthset.Size())/float32(total))
	}
}

func NewDepthStats() DepthStats {
	return DepthStats{
		depth2values: make(map[int]*Set[int]),
		depthcount:   make(map[int]int),
	}
}

func getDepth(depthmap map[Sym]int, args ...Sym) int {
	depth := 0
	for _, a := range args {
		d, ok := depthmap[a]
		if ok {
			depth = max(depth, d)
		}
	}
	return depth + 1
}

func createDepthmap(prog Program) (depthmap map[Sym]int) {
	depthmap = make(map[Sym]int)
	for _, line := range prog {
		depthmap[line.outsym] = getDepth(depthmap, line.argsyms...)
	}
	return depthmap
}
