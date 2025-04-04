package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
)

// isValid checks whether evaluateProgram(p) will err. whether a program
// is internally consistent, meaning is each argsym in each FnCall a reference
// to a sym that exists. The Zero value is hacked in a way that doesn't allow
// for it to live in the Library as a Value...
func isValid[T UncheckedProgram | Program](program T) bool {
	present := NewSet[Sym]()
	if len(program) == 0 {
		ErrorLog.Printf("fail isValid: program len = 0")
		return false
	}
	if program[0].fn.value == nil {
		ErrorLog.Printf("fail isValid: beings with nil")
		return false
	}
	for i, stmt := range program {
		for j, arg := range stmt.argsyms {
			desired_type := stmt.fn.ptypes[j]
			if !present.Contains(arg) {
				ErrorLog.Printf("missing arg: %v with type: %v on line: %v \n", arg, desired_type, i)
				ErrorLog.Printf("the catalog is:\n %v \n", present)
				return false
			}
		}
		present.Add(stmt.outsym)
	}
	return true
}

func panic_if_invalid[T UncheckedProgram | Program](prog T) {
	validateOrFail(prog, "The program is invalid.")
}

func validateOrFail[T UncheckedProgram | Program](prog T, msg string) {
	if !isValid(prog) {
		ErrorLog.Print(msg)
		printProgram(prog, Errr)
		panic(msg)
	}
}

func getSymSet(program UncheckedProgram) SymSet {
	symset := NewSet[Sym]()
	for _, stmt := range program {
		symset.Add(stmt.outsym)
		for _, arg := range stmt.argsyms {
			symset.Add(arg)
		}
	}
	return symset
}

func (program UncheckedProgram) rename_syms(renames map[Sym]Sym) {
	for i, stmt := range program {
		newname, ok := renames[stmt.outsym]
		if ok {
			stmt.outsym = newname
		}
		for i, arg := range stmt.argsyms {
			newname, ok = renames[arg]
			if ok {
				stmt.argsyms[i] = newname
			}
		}
		program[i] = stmt
	}
}

// Change the sym names used in `tochange` to be disjoint from names used in `fixed` (either argsyms or outsyms).
func (tochange UncheckedProgram) uniquify_syms(fixed UncheckedProgram) {
	symset_p1 := getSymSet(fixed)
	symset_p2 := getSymSet(tochange)

	gensym := GenSym{idx: 0}
	renames := make(map[Sym]Sym)
	for _, item := range symset_p2.Elements() {
		renames[item] = gensym.genUnique(symset_p1)
	}

	tochange.rename_syms(renames)
}

type InterleaveConfig struct {
}

// Merge a and b into a new Program by interleaving their statements.
// Proceed from top to bottom.
// program = empty
// while either a or b isn't empty
//
//	choose a or b randomly and pull the top statement
//	if the arguments can be fulfilled by existing syms in the catalog
//	   append statement to program
//	   delete statement from a/b source
//	   add outsym,type to catalog
func interleave(a, b Program, cfg InterleaveConfig) Program {
	validateOrFail(a, "Fail on a.")
	validateOrFail(b, "Fail on b.")
	l, m := UncheckedProgram(CopyProgram(a)), UncheckedProgram(CopyProgram(b))
	// TODO: rename syms from both programs to be simple v0, v1, ... vN
	m.uniquify_syms(l)
	validateOrFail(l, "Fail on l.")
	validateOrFail(m, "Fail on m.")
	catalog := make(map[Type][]Sym)
	prog := make(Program, 0)
	var stmt Statement
	var p *UncheckedProgram
	i := 0
outer:
	for len(prog) < len(l) {
		i += 1
		if rand.Float32() < 0.5 && len(l) != 0 || len(m) == 0 {
			p = &l
		} else {
			p = &m
		}
		stmt = (*p)[0]
		for j, argtype := range stmt.fn.ptypes {
			arg_options := catalog[argtype]
			if len(arg_options) == 0 {
				continue outer
			}
			stmt.argsyms[j] = arg_options[rand.Intn(len(arg_options))]
		}
		// pop statement from program head
		*p = (*p)[1:]
		catalog[stmt.fn.rtype] = append(catalog[stmt.fn.rtype], stmt.outsym)
		prog = append(prog, stmt)
	}
	validateOrFail(prog, "Fail on newly minted prog.")
	return prog
}

// Walk from top to bottom, updating map of type->[]sym
// Each argument can be resampled from the map
// func rewireOrFail(prog UncheckedProgram) (Program, error) {
// 	// orig := make(Program, len(prog))
// 	orig_ := deepcopy.Copy(prog).(UncheckedProgram)
// 	orig, errInvalid := validate(orig_)
// 	catalog := make(map[Type][]Sym)
// 	for _, stmt := range prog {
// 		for j, argtype := range stmt.fn.ptypes {
// 			arg_options := catalog[argtype]
// 			if len(arg_options) == 0 {
// 				if errInvalid != nil {
// 					return Program{}, fmt.Errorf("There isn't a valid rewiring of the program %v \n.", prog)
// 				} else {
// 					// TODO: clean up. Don't use errors to communicate like this.
// 					return orig, fmt.Errorf("There isn't a valid rewiring of the program %v but orig is valid.\n.", prog)
// 				}
// 			}
// 			n := rand.Intn(len(arg_options))
// 			stmt.argsyms[j] = arg_options[n]
// 		}
// 		catalog[stmt.fn.rtype] = append(catalog[stmt.fn.rtype], stmt.outsym)
// 	}
// 	return Program(prog), nil
// }

func CopyProgram[T Program | UncheckedProgram](prog T) T {
	panic_if_invalid(prog)
	prog_copy := make(T, len(prog))
	for i := range prog_copy {
		argsyms := make([]Sym, len(prog[i].argsyms))
		for j := range prog[i].argsyms {
			argsyms[j] = prog[i].argsyms[j]
		}
		prog_copy[i] = Statement{
			fn:      prog[i].fn,
			outsym:  prog[i].outsym,
			argsyms: argsyms,
		}
	}
	panic_if_invalid(prog_copy)
	return prog_copy
}

func rewire_base[T Program | UncheckedProgram](prog T) (Program, bool) {
	catalog := make(map[Type][]Sym)
	panic_if_invalid(prog)
	prog_copy := CopyProgram(prog)
	panic_if_invalid(prog_copy)

	for i, stmt := range prog_copy {
		for j, argtype := range stmt.fn.ptypes {
			arg_options := catalog[argtype]
			if len(arg_options) == 0 {
				panic_if_invalid(prog_copy)
				return Program(prog_copy), false // no new wiring found
			}
			n := rand.Intn(len(arg_options))
			stmt.argsyms[j] = arg_options[n]
		}
		catalog[stmt.fn.rtype] = append(catalog[stmt.fn.rtype], stmt.outsym)
		prog_copy[i] = stmt
	}
	panic_if_invalid(prog_copy)
	return Program(prog_copy), true
}

type ProgramHistoryRow struct {
	Prog   string
	reward float64
	time   int
}

type GPParams struct {
	N_rounds   int
	N_programs int
	Ltype      Mutate
}

type Mutate = int

const (
	NoMut Mutate = iota
	Mut
	ReGen
)

var (
	program_history []ProgramHistoryRow
	cheating        Cheating
)

type Cheating = int

const (
	NoCheating Cheating = iota
	ZeroValue
	ZeroOnlyOnce
)

// We could actually run analysis on values AFTER generating all the programs and decouple
// these two things. Then have different kinds of analyses on []ValueMap (incl histograms).
// But this forces us to realize the entire thing in momory! Better to have run_simple_prog
// just generate programs and eval them continuously, then put them on a chan?
// func run_simple_program_gen(ch chan ValueMap, nprog int, sp SampleParams) {
// 	for i := range nprog {
// 		prog := sampleProgram(sp)
// 		validateOrFail(prog, fmt.Sprintf("invalid program at i=%v\n", i))
// 		// print_program(prog, Info)
// 		values, _ := evalProgram(prog)
// 		ch <- values
// 	}
// }

func Run_genetic_program_optimization(p GPParams) {

	program_history = make([]ProgramHistoryRow, p.N_programs*p.N_rounds)
	programs := make([]Program, p.N_programs)
	new_programs := make([]Program, p.N_programs)
	scores_current := make([]float64, p.N_programs)

	// initialize programs
	for i := range p.N_programs {
		prog := sampleProgram(newSampleParams())
		validateOrFail(prog, fmt.Sprintf("Failed creating sample program %v .\n", i))
		programs[i] = prog
	}

	for i := range p.N_rounds {

		// eval scores
		for k, prog := range programs {
			validateOrFail(prog, fmt.Sprintf("Failed evaluating sample program %v .\n", k))
			printProgram(prog, Info)
			_, reward := evalProgram(prog)
			scores_current[k] = reward
			program_history[i*p.N_programs+k] = ProgramHistoryRow{formatProgram(prog), reward, i}
		}

		// filter out the bottom half of programs
		sort.Slice(programs, func(k, j int) bool {
			return scores_current[k] > scores_current[j]
		})

		// NOTE: this is allowed to alias programs, becuase new_programs contains copies.
		// we've got N programs in the population
		// we want to take the top `n_keep` and explore them
		// or maybe we want to apply some ostu-like threshold
		// then we take these keepers and mutate them in various ways?
		// rewire, shuffle, mutate, insert, delete, prune, grow and interleave.
		// and maybe add a new fresh program from the generator for fun.
		n_keep := p.N_programs / 2
		best_programs := programs[:n_keep]

		// every new program maps (via mod) to an old, top-tier program (and potentially a 2nd)
		for k := range new_programs {
			n := k % len(best_programs)
			switch p.Ltype {
			case 0:
				new_programs[k] = programs[k]
			case 1:
				new_programs[k] = mutate_program(best_programs[n], best_programs)
			case 2:
				new_programs[k] = sampleProgram(newSampleParams())
			}
		}

		// programs = new_programs
		n_cop := copy(programs, new_programs) // TODO: just swap buffer pointers
		assert(n_cop == len(programs), "wrong number of programs!")
		assert(n_cop == len(new_programs), "wrong number of new_programs!")

		reward := mean(scores_current[:])

		fmt.Printf(UPSTART+"Finished genetic optimization round %v score %.3e \n", i, reward)

		global_time += 1
	}
	fmt.Println()

	InfoLog.Println("Final Programs after Genetic Algs are:")
	for _, prog := range programs {
		printProgram(prog, Info)
	}

}

// copies oprog
func reshuffle(oprog Program) Program {
	count := 0
	if !isValid(oprog) {
		panic("we fucked up")
	}
	// validate_or_fail(oprog, "we fucked up somehow badly")
	prog := CopyProgram(oprog)
	for {
		rand.Shuffle(len(prog), func(i, j int) {
			prog[i], prog[j] = prog[j], prog[i]
		})
		if isValid(prog) {
			// fmt.Println("We found a valid reshuffle!")
			return prog
		}
		if count == 1000 {
			// fmt.Println("Revert to original program.")
			// validate_or_fail(oprog, "we fucked up somehow")
			return oprog
		}
		count += 1
	}
}

func mutate_program(p Program, best_programs []Program) Program {
	r := rand.Float32()
	switch {
	case r < 0.1:
		// fmt.Println("sampleProgram_fromFragmentLib")
		return sampleProgram(newSampleParams())
	case r < 0.2:
		// fmt.Println("reshuffle")
		return reshuffle(p)
	case r < 0.3:
		// fmt.Println("point_mutate")
		pnew, _ := PointMutate(p)
		return pnew
	case r < 0.4:
		// fmt.Println("rewire")
		p, _ := rewire_base(p)
		return p
	case r < 0.5:
		// fmt.Println("prune")
		return prune(p)
	case r < 0.7:
		// fmt.Println("grow")
		return grow(p)
	case r < 1.0:
		// fmt.Println("interleave. len(p)=", len(p))
		n := rand.Intn(len(best_programs))
		return interleave(p, best_programs[n], InterleaveConfig{})
	default:
		return p
	}
}

// Replace a FnCall in Program with one from Library with matching rtype and ptypes.
func PointMutate(prog_ Program) (Program, bool) {
	libinv := buildLibraryInverse()
	replaceables := make(map[int]*Set[Sym])
	prog := CopyProgram(prog_) // TODO: does this reassign the pointer? I think yes.
	valid_fnset := NewSet[Sym]()
	var i int
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("The panic'd program is: ", i, " program \n ", formatProgram(prog))
			fmt.Println("Len = ", len(prog), "raw \n", prog)
			fmt.Println("The libinv.provides : ", libinv.provides)
			panic("Continue panicing")
		}
	}()
	for i = range prog {
		// require same rtype
		valid_fnset.Clear()
		// tmp := libinv.provides[prog[i].fn.rtype] // WARN: This doesn't copy the underlying map! Modifications affect libinv!

		// require same ptypes
		// fmt.Println("rtype: ", prog[i].fn.rtype)
		for fnsym := range libinv.provides[prog[i].fn.rtype].emap {
			isEqPTypes := reflect.DeepEqual(prog[i].fn.ptypes, fn_library[fnsym].ptypes)
			if isEqPTypes {
				valid_fnset.Add(fnsym)
			}
		}

		valid_fnset.Remove(Sym(prog[i].fn.name))
		InfoLog.Printf("The fn %+v on line %v is replaceable by: %+v \n", prog[i].fn.name, i, valid_fnset)

		if valid_fnset.Size() > 0 {
			replaceables[i] = NewSetFromMapKeys(valid_fnset.emap)
			InfoLog.Println("Found replaceable line ", i, " fn set ", valid_fnset)
		}
	}
	if len(replaceables) == 0 {
		InfoLog.Println("No replaceable lines!")
		return prog, false
	}
	lines := sortedKeys(replaceables)
	line := lines[rand.Intn(len(lines))] // select random line with replaceable fn
	options, exists := replaceables[line]
	assert(exists, fmt.Sprintf("Must exist by construction. line=%v , options=%+v \n", line, options.Elements()))
	assert(options.Size() > 0, fmt.Sprintf("Must have multiple options by construction! line=%v and options=%+v", line, options.Elements()))
	newfn := options.Elements()[rand.Intn(options.Size())]
	InfoLog.Printf("Replacing %v with %v at line %v.\n", prog[line].fn.name, newfn, line)
	newprog := make(Program, len(prog))
	copy(newprog, prog)
	newprog[line].fn = fn_library[newfn]

	validateOrFail(newprog, "Mutation produced invalid program.")

	return newprog, true
}

func prune(prog Program) Program {
	return prog
}

func grow(prog Program) Program {
	return prog
}

func getCallSyms(prog Program) []Sym {
	s := make([]Sym, len(prog))
	for i := range prog {
		s[i] = Sym(prog[i].fn.name)
	}
	return s
}

// We're trying to cut a program in half, keeping a valid base.
// Half could mean 1.
// TODO: finish me
func divide(prog Program) Program {
	new := CopyProgram(prog)
	// In a dataflow lang it's always safe to simply pick a line and remove everything afterwards.
	// This is not true if we have openResource() and dropResource() far apart in the program.
	// If dropResource() is implicit then there is no problem.

	return new
}

// divide is just one step in `minimize(prog, f)`
// where f : Values -> bool
// In order to minimize we have to run f on the values produced by evalProgram(prog)
// and maintain the state of the boolean.
// We can do a much better job of minimizing if we understand the full dataflow
//  in addition to just treating the whole program naively as if every line depended
//  on all previous lines.
//  We already have code which ensures that a program is valid by analyzing deps.
//  We need to be able to answer questions like:
//   - Does line 15 depend on line 12?
//   - What set of lines does 16 depend on? What set depends on 16?
//  Then the full alg would be:
