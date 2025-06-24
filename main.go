package main

import (
	"flag"
)

var dbname *string

func init() {
	dbname = flag.String("d", "", "database to connect")
}

func main() {

	pick_number := flag.Int("n", -1, "Can we build it? Yes we can!")

	flag.Parse()
	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	if *pick_number != -1 {
		gobTheBuilder(*pick_number)
		return
	}

	// testDataflow()
	// test_dagc()

	// iterate()
	// os.Exit(0)

	// p := sample2lvl()
	// vm, _ := evalProgram(p)
	// printProgramAndValues(p, vm)

	// deltaDebug()
	// benchmarkSampleProgram()
	// runPow2()
	runPow2Dataflow()
	// runWire()
	// runGenetic()
}
