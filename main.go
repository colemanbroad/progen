package main

import (
	"flag"
	"fmt"
	"os"
)

var dbname *string

func main() {
	dbname = flag.String("d", "", "Database to connect")
	gob := flag.Int("gob", 73, "Can we build it? Yes we can!")
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println("dbname = ", *dbname)
		flag.Usage()
		return
	}

	if gob != nil {
		fn_library = make(map[Sym]Fun)
		value_library = make(map[Sym]Value)
		addBasicMathLib()
		// delete(fn_library, "one")
		gobTheBuilder(*gob)
		os.Exit(0)
	}
	// deltaDebug()

	// benchmarkSampleProgram()
	runPow2()
	// runWire()
	// runGenetic()
}
