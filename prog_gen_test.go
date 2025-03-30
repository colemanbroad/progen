package main

import (
	"math"
	"math/rand"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	// Open the log file
	file, err := os.OpenFile(logdir+"test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		ErrorLog.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	// Set log output to the file
	InfoLog.SetOutput(file)
	ErrorLog.SetOutput(file)

	// Run tests
	code := m.Run()

	// Exit with the test code
	os.Exit(code)
}

func Test_ShortestPath(t *testing.T) {
	for fn_name := range Library {
		bestpath := shortestFuncPath(make([]string, 0), string(fn_name))
		InfoLog.Printf("The best path to %v is: %v", fn_name, bestpath)
	}
}

func Test_TypeGraph_and_LibraryInverse(t *testing.T) {
	buildLibraryInverse()
	type_set := NewSet[Type]()
	type_set.Add(Type("bool"))
	buildTypeGraph(type_set)
}

func Test_BasicPrgramGen(t *testing.T) {
	for range 100 {
		sp := newSampleParams()
		if rand.Float32() < 0.5 {
			sp.Wire_nearby = !sp.Wire_nearby
		}
		p1 := UncheckedProgram(sampleProgram_fromFragmentLib(sp))
		program := Program(p1)
		valmap, r_delta := evalProgram(program)
		_, _ = valmap, r_delta
	}
}

func Test_TruncatedExponential(t *testing.T) {
	// n := time.Now().UnixNano()
	n := int64(42)
	ran := rand.New(rand.NewSource(n))
	ran.Float64()

	// Example parameters
	lambda := 1.0 // Rate parameter
	trunc := 5.0  // Truncation point

	// TODO: Actually compare results of distribution with analytical probabilities.
	xs := [1000]float64{}
	for range len(xs) * 100 {
		x := TruncatedExponentialSampler(ran, lambda, trunc)
		bucket := int(math.Floor(x / trunc * 1000.0))
		xs[bucket] += 1
	}
}

func Test_BasicProgramEval(t *testing.T) {
	program := Program{
		Statement{
			fn:      Library["one"],
			outsym:  "v1",
			argsyms: []Sym{},
		},
		Statement{
			fn:      Library["one"],
			outsym:  "v2",
			argsyms: []Sym{},
		},
		Statement{
			fn:      Library["add"],
			outsym:  "v3",
			argsyms: []Sym{"v1", "v2"},
		},
	}
	// fmt.Println(program)
	values, _ := evalProgram(program)
	v, exists := values["v3"]
	if !exists {
		t.Fail()
	}
	if !(reflect.ValueOf(v.value).Int() == 2) {
		t.Fail()
	}
}
