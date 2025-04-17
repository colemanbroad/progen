package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

func BenchmarkSampleProgram(b *testing.B) {
	if false {
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	lib := NewLib()
	// fn_library = make(map[Sym]Fun)
	lib.addBasicMathLib()
	// addPeanoLib()
	t0 := time.Now()
	sp := newSampleParams()
	for i := range 14 {
		sp.Program_length = 1 << i
		fmt.Println("program length: ", sp.Program_length)
		prog := lib.sampleProgram(sp)
		fmt.Println("time Gen: ", time.Now().Sub(t0))
		t0 = time.Now()
		evalProgram(prog)
		fmt.Println("time Eval: ", time.Now().Sub(t0))
		t0 = time.Now()
	}
}
