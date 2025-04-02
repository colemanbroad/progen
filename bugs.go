package main

import "fmt"

// The following code has a bug, which we found by generating the program!.
func bugOne() {
	// prog := `
	// 	a = 1
	// 	b = 1 << 63
	// 	c = 1 << b
	// `

	a := 1
	b := 1 << 32
	b = b << 31 // 1 << 63
	fmt.Println(b)
	b = sign(b) * b // Buggy impl of Abs(b). Fails only on INTMIN
	fmt.Println(b)
	x := a << b // FAIL: negative shift amount
	fmt.Println(x)
}
