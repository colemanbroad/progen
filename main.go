package main

import (
	"flag"
	"fmt"
)

var dbname *string

func main() {
	dbname = flag.String("d", "", "Database to connect")
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Println("dbname = ", *dbname)
		flag.Usage()
		return
	}

	// deltaDebug()
	testDeltaD()

	// benchmarkSampleProgram()
	// runPow2()
	// runWire()
	// runGenetic()
}
