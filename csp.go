package main

import (
	"math/rand/v2"
	"time"
)

const (
	n_global = 5
)

type Machine struct {
	state int
}

func (m *Machine) read() int {
	s := rand.Int64N(n_global)
	time.Sleep(time.Duration(s) * time.Millisecond)
	return m.state
}

func (m *Machine) write(new int) {
	s := rand.Int64N(n_global)
	time.Sleep(time.Duration(s) * time.Millisecond)
	m.state = new
	// s = rand.Int64N(n_global)
	// time.Sleep(time.Duration(s) * time.Millisecond)
}

// Define two state machines. We want them to perform a task where
// a is given a new state from a list, updates it's state, then a tells b
// to update to the same value... or we make
// Or we have two clients writing to the same db. The db holds state and
// the clients read it and write back the value + 1. The
func someTask() {
	a := Machine{1}
	b := Machine{2}
}
