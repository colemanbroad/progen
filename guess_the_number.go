package main

import "fmt"

func gobTheBuilder(n int) {
	sp := newSampleParams()
	sp.Program_length = 1000
	input := sampleProgram(sp)
	// why []Statement and not Program is loadbearing?
	test := func(p []Statement) bool {
		vm, _ := evalProgram(p)
		b1 := false
		for _, v := range vm {
			if v.value.(int) == n {
				b1 = true
			}
		}
		return b1
	}
	count := 0
	for !test(input) {
		input = sampleProgram(sp)
		count += 1
		if count%1000 == 0 {
			fmt.Println(count)
		}
	}
	vm, _ := evalProgram(input)
	fmt.Println("The starter program:")
	printProgramAndValues(input, vm)
	result := deltaD(input, test)
	fmt.Println("The resulting program:")
	vm, _ = evalProgram(result)
	printProgramAndValues(result, vm)
}
