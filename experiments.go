package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"
	// "log"
)

var (
	program_history []ProgramHistoryRow
)

type Cheating = int

const (
	NoCheating Cheating = iota
	ZeroValue
	ZeroOnlyOnce
)

func initPeano(cheating Cheating) {
	fn_library = make(map[Sym]Fun)
	addPeanoLib()
	if cheating == ZeroOnlyOnce {
		// NOTE: explicit make() init not needed for program_prefix?
		program_prefix = Program{Statement{
			fn: Fun{
				value:  func() int { return 0 },
				name:   "LitZero",
				ptypes: []Type{},
				rtype:  "int",
			},
			outsym:  "symzero",
			argsyms: []Sym{},
		}}
		delete(fn_library, "zero")

		// if we remove Zero from fn_lib here then it won't be available
		//  when we need to eval the program later!
		// So... this method doesn't work. We could have a generic
		// "filter library syms" function that is applied during
		// sampleProgram... But this will not work. We could use a
		// Progen style Litnum here?! OK, let's try that!

	} else if cheating == ZeroValue {
		value_library = make(map[Sym]Value)
		// TODO: How are we going to allow for ZeroValue in Programs
		// with the same semantics as Values in Rust GenTactics?
		value_library["Zero"] = Value{
			value: 0,
			name:  "Zero",
			vtype: "int",
		}
		delete(fn_library, "zero")
	}
}

// Cheating: Do we remove Zero from the Library after line 1 ? Peano Specific.
// Decay: What is the decay rate in the exponential dist over lines when sampling an input arg.
// Proglen: Program length
func runPeano() {
	for _, c := range []Cheating{ZeroValue} {
		for _, decay := range []float64{0.1, 1.0, 0.0} {
			for _, proglen := range []int{10, 20, 50, 100} {
				initPeano(Cheating(c))
				sp := newSampleParams()
				sp.Wire_nearby = true
				if decay == 0.0 {
					sp.Wire_nearby = false
				}
				sp.WireDecayLen = decay
				sp.Program_length = proglen
				fmt.Println("Begin wiring: ", sp)
				nprog := 1
				vh := make(IntHistogram)
				for range nprog {
					prog := sampleProgram(sp)
					values, _ := evalProgram(prog)
					printProgramAndValues(prog, values)
					for _, v := range values {
						vh.add(v.value.(int))
					}
				}
				fmt.Println("ProgLen = ", sp.Program_length, " Value Hist: ", vh)
				savePeano(sp, nprog, vh, Cheating(c))
			}
		}
	}
}

func savePeano(sp SampleParams, nprog int, valuehist map[int]int, cheating Cheating) {
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
func runPow2() {
	fn_library = make(map[Sym]Fun)
	addBasicMathLib()
	// addPowerOfTwo()
	for _, proglen := range []int{100} {
		for _, decay := range []float64{0.0, 0.1, 1.0} {
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
			stats := NewDepthStats()
			global_time = 0
			for range 1 {
				// fmt.Println("i = ", i)
				prog := sampleProgram(sp)
				vals, _ := evalProgram(prog)
				// printProgramAndValues(prog, vals)
				stats.update(prog, vals)
				global_time += 1
			}
			stats.print()
		}
	}
}

func savePow2(sp SampleParams) {
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
func runGenetic() {
	fn_library = make(map[Sym]Fun)
	addBasicMathLib()
	addPowerOfTwo()
	for range 20 {
		init_history()
		// init_maphistory()
		init_reward()
		// Init_campaign()
		p := GPParams{N_rounds: 1000, N_programs: 20, Ltype: NoMut}
		RunGenetic(p)
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

func benchmarkSampleProgram() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	fn_library = make(map[Sym]Fun)
	addBasicMathLib()
	// addPeanoLib()
	sp := newSampleParams()
	sp.Wire_nearby = false
	var t0 time.Time
	for i := range 18 {
		sp.Program_length = 1 << i
		fmt.Println("program length: ", sp.Program_length)
		t0 = time.Now()
		prog := sampleProgram(sp)
		fmt.Println("time Gen: ", time.Now().Sub(t0))
		t0 = time.Now()
		evalProgram(prog)
		fmt.Println("time Eval: ", time.Now().Sub(t0))
	}
}
