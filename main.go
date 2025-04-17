package main

import (
	"flag"
	"fmt"
	"os"
)

var dbname *string

func main() {
	dbname = flag.String("d", "", "database to connect")
	gob := flag.Int("gob", -1, "Can we build it? Yes we can!")
	lib := NewLib()
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println("dbname = ", *dbname)
		flag.Usage()
		return
	}

	if *gob != -1 {
		fn_library = make(map[Sym]Fun)
		value_library = make(map[Sym]Value)
		lib.addBasicMathLib()
		// delete(fn_library, "one")
		gobTheBuilder(*gob)
		os.Exit(0)
	}

	p := sample2lvl()
	vm, _ := evalProgram(p)
	printProgramAndValues(p, vm)

	// deltaDebug()
	// benchmarkSampleProgram()
	// runPow2()
	// runWire()
	// runGenetic()
}
