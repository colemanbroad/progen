package main

// This is my crappy attempt to write dictproduct. It's not even generic, but first
// specific to the WiringExperiment.
type ParamSets struct {
	Wire_nearby    []bool
	Program_length []int
	WireDecayLen   []float64
	n_iter         []int
	cheat          []bool
}

// How can I write dictproduct??
// First lets write it for the above: sp, n_iter, cheating
func iterate_wirings(params ParamSets) {
	params = ParamSets{
		Wire_nearby:    []bool{true, false},
		Program_length: []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
		WireDecayLen:   []float64{1.0, 0.5, 0.1},
		n_iter:         []int{1000},
		cheat:          []bool{true, false},
	}
}
