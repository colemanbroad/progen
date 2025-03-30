package main

import (
	"testing"
)

func Test_create_set(t *testing.T) {
	// Create sets
	setA := NewSet[int]()
	setB := NewSet[int]()

	// Add elements
	setA.Add(1)
	setA.Add(2)
	setA.Add(3)

	setB.Add(3)
	setB.Add(4)
	setB.Add(5)

	// Perform operations
	InfoLog.Println("Set A:", setA.Elements())
	InfoLog.Println("Set B:", setB.Elements())

	union := setA.Union(setB)
	InfoLog.Println("Union:", union.Elements())

	intersection := setA.Intersection(setB)
	InfoLog.Println("Intersection:", intersection.Elements())

	difference := setA.Difference(setB)
	InfoLog.Println("Difference (A - B):", difference.Elements())
}
