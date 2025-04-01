package main

import (
	"flag"
	"fmt"
)

var dbname *string

func main() {

	dbname = flag.String("d", "", "Database to connect")
	flag.Parse()
	// wireexp := flag.Bool("w", false, "Run the wiring experiment")
	// create_tables := flag.Bool("c", false, "Create new tables in db.")
	// learn := flag.Int("l", 0, "Should we learn over time? 0: no, 1: yes, 2: yes + interleave")
	// n_iter := flag.Int("n", 0, "Run n rounds of genetic optimization. Otherwise run sqlpeek.")

	if flag.NFlag() == 0 {
		fmt.Println("dbname = ", *dbname)
		flag.Usage()
		return
	}

	runP2()
	// runWire()
	// runGeneticExperiment()

}
