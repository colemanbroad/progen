package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Simplest race condition.
// main and func1 compete to read/write `m`
func f1() {
	m := "one"
	go func() {
		m = "two"
	}()
	fmt.Println(m)
}

// Use a chan to wait for func1 to be done writing m. Remove either
// of the channel send/recv lines to introduce a deadlock or a race
// condition.
func f2() {
	dchan := make(chan bool)
	m := "one"
	go func() {
		m = "two"
		dchan <- true // deadlock
	}()
	<-dchan // data race
	fmt.Println(m)
}

// The control flow may depend on buffer size! It may be hard to
// predict which thread the scheduler will choose first. 0 prints
// "two". 1 prints "one". This isn't technically a race condition,
// just a nondeterminism.
func f3() {
	dchan := make(chan bool, 1) // toggle 0 | 1
	m := "one"
	go func() {
		<-dchan
		m = "two"
		dchan <- true
	}()
	dchan <- true
	<-dchan
	fmt.Println(m)
}

// Change the timing and a different thread is scheduled first.
// This is _effectively_ a race condition. But it presents as a
// nondeterminism, and isn't picked up by -race flag.
func f4() {
	dchan := make(chan bool, 1)
	m := "one"
	go func() {
		<-dchan
		m = "two"
		dchan <- true
	}()
	dchan <- true
	time.Sleep(0 * time.Millisecond) // 0ms="one" ; 1ms="two"
	<-dchan
	fmt.Println(m)
}

// A nondeterminism isn't a race condition.
func f5() {
	dchan := make(chan bool)
	m := "one"
	f := func() {
		m = "two"
		dchan <- true
	}
	g := func() {
		m = "three"
		dchan <- true
	}
	if rand.Float32() < 0.5 {
		go f()
	} else {
		go g()
	}
	<-dchan
	fmt.Println(m)
}

func main2() {
	f1()
	f2()
	f3()
	f4()
	f5()
}
