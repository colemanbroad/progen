package main

import (
	"fmt"
	// "log"
)

// Cheating: Do we remove Zero from the Library after line 1 ? Peano Specific.
// Decay: What is the decay rate in the exponential dist over lines when sampling an input arg.
// Proglen: Program length
func runWire() {
	Library = make(map[Sym]FnCall)
	addPeanoLib()
	for _, c := range []Cheating{ZeroValue} {
		for _, decay := range []float64{0.1, 1.0, 0.0} {
			for _, proglen := range []int{10, 20, 50, 100} {
				cheating = Cheating(c)
				sp := newSampleParams()
				sp.Wire_nearby = true
				if decay == 0.0 {
					sp.Wire_nearby = false
				}
				sp.WireDecayLen = decay
				sp.Program_length = proglen
				fmt.Println("Begin wiring: ", sp)
				wire_inner(1000, sp)
			}
		}
	}
}

type ValueHistogram map[int]int

func (h ValueHistogram) add(val int) {
	c, exist := h[val]
	if !exist {
		c = 0
	}
	h[val] = c + 1
}

func wire_inner(nprog int, sp SampleParams) {
	vh := make(ValueHistogram)
	for range nprog {
		prog := sampleProgram(sp)
		values, _ := evalProgram(prog)
		for _, v := range values {
			vh.add(v.value.(int))
		}
	}
	fmt.Println("ProgLen = ", sp.Program_length, " Value Hist: ", vh)
	saveWire(sp, nprog, vh)
}

func saveWire(sp SampleParams, nprog int, valuehist map[int]int) {
	fmt.Println("saving: ", sp)
	db := ConnectSqlite(*dbname)
	defer db.Close()
	_, err := db.Exec("create table if not exists wiring(prog_l int, wr_decay real, wr_nearby bool, n_prog int, depth int, count int, cheating int)")
	check(err)
	// prog_l int, wr_decay real, wr_nearby bool, n_prog int, depth int, count int
	for k, v := range valuehist {
		_, err := db.Exec("insert into wiring values(?,?,?,?,?,?,?)",
			sp.Program_length, sp.WireDecayLen,
			sp.Wire_nearby, nprog, k, v, cheating)
		if err != nil {
			fmt.Println("Error saving wiring: ", sp, nprog)
			return
		}
	}
}

// How does deeper wiring affect the Powers of Two distribution?
func runP2() {
	Library = make(map[Sym]FnCall)
	addBasicMathLib()
	addPowerOfTwo()
	for _, proglen := range []int{100} {
		for _, decay := range []float64{0.1, 1.0, 0.0} {
			sp := newSampleParams()
			sp.Wire_nearby = true
			if decay == 0.0 {
				sp.Wire_nearby = false
			}
			sp.WireDecayLen = decay
			sp.Program_length = proglen
			fmt.Println("Begin wiring: ", sp)
			init_history()
			init_reward()
			global_time = 0
			for range 1000 {
				// fmt.Println("i = ", i)
				prog := sampleProgram(sp)
				_, _ = evalProgram(prog)
				global_time += 1
			}
			saveP2(sp)
		}
	}
}

func saveP2(sp SampleParams) {
	db := ConnectSqlite(*dbname)
	defer db.Close()
	var err error
	var s string
	campaign_id := generateRandomString(16)
	s = `create table if not exists wire_pow_of_two (value real, reward real, time int, campaign_id string, proglen int, decay float)`
	_, err = db.Exec(s)
	check(err)
	// s = `create table if not exists program_history (prog string, reward real, time int, campaign_id string)`
	// _, err = db.Exec(s)
	// check(err)
	// s = `create table if not exists campaigns (campaign_id string, dtime datetime, n_iter int, ltype int)`
	// _, err = db.Exec(s)
	// check(err)
	tx, err := db.Begin()
	check(err)
	s = "insert into wire_pow_of_two (value, reward, time, campaign_id, proglen, decay) values (?,?,?,?,?,?)"
	stmt, err := tx.Prepare(s)
	check(err)
	defer stmt.Close()
	for _, vr := range history_power_of_two {
		// fmt.Println("saving vr", vr)
		_, err = stmt.Exec(vr.Value, vr.Reward, vr.Time, campaign_id, sp.Program_length, sp.WireDecayLen)
		check(err)
	}
	err = tx.Commit()
	check(err)
}

// How does mutation affect PowerOfTwo?
func runGeneticExperiment() {
	Library = make(map[Sym]FnCall)
	addBasicMathLib()
	addPowerOfTwo()
	for range 20 {
		init_history()
		// init_maphistory()
		init_reward()
		// Init_campaign()
		p := GPParams{N_rounds: 1000, N_programs: 20, Ltype: NoMut}
		Run_genetic_program_optimization(p)
		saveGenetic(p)
	}
}

func saveGenetic(p GPParams) {
	db := ConnectSqlite(*dbname)
	defer db.Close()
	var err error
	var s string
	campaign_id := generateRandomString(16)
	// s = `create table if not exists history_power_of_two (value real, reward real, time int, campaign_id string, mut int)`
	// _, err = db.Exec(s)
	// check(err)
	// s = `create table if not exists program_history (prog string, reward real, time int, campaign_id string)`
	// _, err = db.Exec(s)
	// check(err)
	// s = `create table if not exists campaigns (campaign_id string, dtime datetime, n_iter int, ltype int)`
	// _, err = db.Exec(s)
	// check(err)
	tx, err := db.Begin()
	check(err)
	s = "insert into history_power_of_two (value, reward, time, campaign_id, mut) values (?,?,?,?,?)"
	stmt, err := tx.Prepare(s)
	check(err)
	defer stmt.Close()
	for _, vr := range history_power_of_two {
		_, err = stmt.Exec(vr.Value, vr.Reward, vr.Time, campaign_id, p.Ltype)
		check(err)
	}
	err = tx.Commit()
	check(err)
}

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
