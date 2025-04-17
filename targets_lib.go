package main

// f = func() string { return fmt.Sprint(rand.Intn(10)) }
// addFn(f, "rand_s", []Type{}, "string")

// f = func(a string) string { return a + a }
// addFn(f, "repeat", []Type{"string"}, "string")

// f = func(a, b string) string { return a + b }
// addFn(f, "cat", []Type{"string", "string"}, "string")

func sign(x int) int {
	if x < 0 {
		return -1
	} else if x == 0 {
		return 0
	}
	return 1
}

func (lib Library) addBasicMathLib() {
	var f any

	// f = func() int { return rand.Intn(10) }
	// lib.addFuncToLibrary(f, "rand", []Type{}, "int")

	f = func() int { return 1 }
	lib.addFuncToLibrary(f, "one", []Type{}, "int")

	f = func(a, b int) int { return a + b }
	lib.addFuncToLibrary(f, "add", []Type{"int", "int"}, "int")

	f = func(a, b int) int { return a * b }
	lib.addFuncToLibrary(f, "mul", []Type{"int", "int"}, "int")

	f = func(a, b int) int {
		c := b % 64
		c = sign(c) * c
		r := a << c
		return r
	}
	lib.addFuncToLibrary(f, "<<", []Type{"int", "int"}, "int")

	// f = func(isPos bool) float64 {
	// 	x := rand.Float64()
	// 	if !isPos {
	// 		x *= -1
	// 	}
	// 	return x
	// }
	// lib.addFuncToLibrary(f, "samplePosOrNeg", []Type{"bool"}, "f64")
}
