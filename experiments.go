package main

import (
	"fmt"
	// "log"
)

func runWiringExperiment() {
	Library = make(map[Sym]FnCall)
	addPeanoLib()
	for _, c := range []Cheating{Normal, ZeroOnlyOnce} {
		for _, decay := range []float64{0.1, 1.0, 0.0} {
			for _, i := range []int{1, 2, 5, 10} {
				cheating = Cheating(c)
				sp := newSampleParams()
				sp.Wire_nearby = true
				if decay == 0.0 {
					sp.Wire_nearby = false
				}
				sp.WireDecayLen = decay
				sp.Program_length = 10 * i
				fmt.Println("Begin wiring: ", sp)
				run_basic_program_gen(1000, sp)
			}
		}
	}
}

func saveWiring(sp SampleParams, nprog int, valuehist map[int]int) {
	db := ConnectSqlite("wiring.sqlite")
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

func runGeneticExperiment() {
	Library = make(map[Sym]FnCall)
	addBasiMathLib()
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
	var err error
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
	stmt, err := tx.Prepare("insert into history_power_of_two (value, reward, time, campaign_id, mut) values (?,?,?,?,?)")
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
