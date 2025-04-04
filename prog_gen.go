package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"reflect"
)

// TODO: Is the Value type actually necessary? We can just store
// ValueMap as map[Sym]any... But Value let's us remember what (r)type
// the value _thinks_ it is! In case of future conflict.
type Value struct {
	value any
	name  string
	vtype Type
}

type Fun struct {
	value  any
	name   string
	ptypes []Type
	rtype  Type
}

func (fn Fun) String() string {
	return fmt.Sprintf("Fn {name:%v, ptypes:%v, rtype:%v}", fn.name, fn.ptypes, fn.rtype)
}

type Sym string
type Type string
type SymLine struct {
	sym  Sym
	line uint16
}
type SymSet = *Set[Sym]
type Statement struct {
	fn      Fun
	outsym  Sym
	argsyms []Sym
}
type Program []Statement
type UncheckedProgram []Statement
type ValueMap map[Sym]Value

type LogType uint8

const (
	Info LogType = iota
	Warn
	Errr
)

var fn_library map[Sym]Fun
var value_library map[Sym]Value
var program_prefix Program // NOTE: when construcing the prefix we must uphold the constraints inherent in Program (vs UncheckedProgram)

func evalStatement(stmt Statement, locals ValueMap) {
	var r any
	// var f any
	// f = stmt.fn.value
	// type erase the function we're calling
	g := reflect.ValueOf(stmt.fn.value)

	// then erase the args, (Q: but haven't they arleady been erased?)
	// TODO: avoid alloc: make this fix-size array?
	args := make([]reflect.Value, 0)
	for _, a := range stmt.argsyms {
		args = append(args, reflect.ValueOf(locals[a].value))
	}

	// fmt.Println("f, g, args", f, g, args)
	// printProgram(Program{stmt}, Info)

	// return val cast to `any` type aka Interface()
	r = g.Call(args)[0].Interface()

	// Here's where we can introduce the logic of "any" types? generic types?
	//
	rtype := Type(reflect.TypeOf(r).Name())

	// Check that the types are correct
	// if rtype != stmt.fn.rtype {
	// 	ErrorLog.Fatalf("Val type doesn't match fn rtype: %v %v", rtype, stmt.fn.rtype)
	// }

	val := Value{value: r, name: fmt.Sprintf("%v", r), vtype: rtype}
	locals[stmt.outsym] = val
}

func NewValueMap() map[Sym]Value {
	locals := make(ValueMap)
	return locals
}

func evalProgram(program Program) (values ValueMap, reward float64) {
	r0 := Reward_total
	locals := NewValueMap()
	if value_library != nil {
		for sym, val := range value_library {
			locals[sym] = val
		}
	}
	for _, stmt := range program {
		evalStatement(stmt, locals)
	}
	delta := Reward_total - r0
	// if len(history) >= 1 && history[len(history)-1].Value > 1<<20 {
	// 	fmt.Print("The value produced is too big.\n\n")
	// 	fmt.Print("The global time is: ", global_time, "\n\n")
	// printProgramAndValues(program, locals)
	// 	os.Exit(0)
	// }
	return locals, delta
}

func addFuncToLibrary(f any, name string, ptypes []Type, rtype Type) {
	fdef := Fun{
		value:  f,
		name:   name,
		ptypes: ptypes,
		rtype:  rtype,
	}
	fn_library[Sym(fdef.name)] = fdef
}

func sampleFuncFromLibrary() Fun {
	keys := make([]Sym, 0)
	// fmt.Println("fn_library = ", fn_library)
	for k := range fn_library {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		log.Fatalln("error: did you forget to add fragments to the library?")
	}
	k := keys[rand.Intn(len(keys))]
	return fn_library[k]
}

type SampleParams struct {
	Wire_nearby    bool
	Program_length int
	WireDecayLen   float64
}

func newSampleParams() SampleParams {
	return SampleParams{
		Wire_nearby:    true,
		Program_length: 20,
		WireDecayLen:   1.0,
	}
}

// Then generating a program is easier if we keep track of the
// values we will have at our disposal at every point.
// i.e. add lines conditional on what's already in the program.
// Depends on globals: fn_library, value_library and program_prefix, but only fn_library must be non-nil.
func sampleProgram(sp SampleParams) Program {
	gensym := GenSym{idx: 0}
	program := make(Program, 0, sp.Program_length)
	// TODO: use the Catalog with lineno info to have more control over initial wiring
	local_catalog := NewCatalog()
	depthmap := make(map[Sym]int)

	if program_prefix != nil {
		prefix := CopyProgram(program_prefix)
		for line_no, stmt := range prefix {
			depthmap[stmt.outsym] = getDepth(depthmap, stmt.argsyms...)
			local_catalog.add(stmt.outsym, stmt.fn.rtype, uint16(line_no))
			// gensym.add(stmt.outsym) // TODO: impl this so we can be sure to avoid Sym collisions!
			program = append(program, stmt)
		}
	}

	global_catalog := NewCatalog()
	for sym, val := range value_library {
		global_catalog.add(sym, val.vtype, 0) // FIXME! line is wrong
	}

stmtLoop:
	for len(program) < sp.Program_length {
		var f Fun
		f = sampleFuncFromLibrary()
		stmt := Statement{
			fn:      f,
			argsyms: make([]Sym, len(f.ptypes)),
			outsym:  "void_sym",
		}
		// sample a random sym of the appropriate type from the Program above for each argument
		if len(f.ptypes) != 0 {
			for i, ptype := range f.ptypes {
				all_syms, _ := local_catalog.syms_inv[ptype] // TODO: revert SymLine idea. or add Values as SymLines with Line = -1 ?
				all_syms = append(global_catalog.syms_inv[ptype], all_syms...)
				// WARN: This is necessary because of exponential sampling. We have to make
				// an independent decision about what probability to assign to syms from value_library vs the program body.
				if len(all_syms) == 0 {
					continue stmtLoop
				}
				var n int
				if sp.Wire_nearby {
					// Exponentially distributed sampling.
					m := len(all_syms)
					n = m - int(TruncatedExponentialSampler(nil, sp.WireDecayLen, float64(m))) - 1
				} else {
					// n = (idx * 103823) % len(symtypes)
					// idx = idx + 1 // TODO: change the math here to control the wiring
					n = rand.Intn(len(all_syms))
				}
				arg := all_syms[n].sym
				if cheating == ZeroValue && rand.Float64() < 1.0/(float64(len(program)+1)) {
					arg = Sym("Zero")
				}
				stmt.argsyms[i] = arg //= append(stmt.argsyms, arg)
			}
		}
		stmt.outsym = gensym.gen()
		// fmt.Printf("add() in sampleProgram() sym = %v \n", stmt.outsym)
		// fmt.Printf("gensym = %v", gensym)
		depthmap[stmt.outsym] = getDepth(depthmap, stmt.argsyms...)
		local_catalog.add(stmt.outsym, stmt.fn.rtype, uint16(len(program))) // WARN: will break when len(prog) >= 2^16
		program = append(program, stmt)
	}
	// fmt.Print(depthmap)
	return program
}

func createDepthmap(prog Program) (depthmap map[Sym]int) {
	depthmap = make(map[Sym]int)
	for _, line := range prog {
		depthmap[line.outsym] = getDepth(depthmap, line.argsyms...)
	}
	return depthmap
}

func getDepth(depthmap map[Sym]int, args ...Sym) int {
	depth := 0
	for _, a := range args {
		d, ok := depthmap[a]
		if ok {
			depth = max(depth, d)
		}
	}
	return depth + 1
}

// TruncatedExponentialSampler draws a f64 value from an exponential distribution
// with rate λ, truncated to [0, trunc). It does not allocate mass at trunc.
func TruncatedExponentialSampler(rng *rand.Rand, lambda, trunc float64) float64 {
	if lambda <= 0 {
		panic("TruncatedExponentialSampler: lambda must be > 0")
	}
	if trunc <= 0 {
		panic("TruncatedExponentialSampler: trunc must be > 0")
	}

	// if rng == nil {
	// 	return rand.ExpFloat64() / lambda
	// }
	var u float64
	if rng == nil {
		u = rand.Float64()
	} else {
		u = rng.Float64()
	}

	// u := rng.Float64() // U ~ Uniform(0, 1)
	cutoff := 1 - math.Exp(-lambda*trunc)
	// Solve for y in:
	//   1 - e^{-lambda y} = u * cutoff
	// => y = -1/lambda * ln(1 - u * cutoff)
	return -math.Log(1-u*cutoff) / lambda
}

// TruncatedExponentialSampler samples from a truncated exponential distribution
// with rate lambda (λ) and truncation point trunc, without putting extra probability mass at trunc.
func TruncatedExponentialSampler2(lambda, trunc float64, ran *rand.Rand) float64 {
	// Ensure trunc is positive
	if trunc <= 0 {
		panic("truncation point must be positive")
	}

	// Generate a random value between 0 and 1
	var u float64
	if ran == nil {
		u = rand.Float64()
	} else {
		u = ran.Float64()
	}

	// t' ~ c e^-t/tau = p(t')
	// cdf(t') = int_t=0,inf p(t') = -tau c e^-t/tau | t=inf - t=0 = 0 - (- tau c) = tau c

	// CDF at the truncation point
	cdfTrunc := 1 - math.Exp(-lambda*trunc)

	// Apply the inverse CDF for the truncated exponential distribution
	sample := -math.Log(1-u*cdfTrunc) / lambda
	return sample
}

func formatProgram[T Program | UncheckedProgram](program T) string {
	str := ""
	for i, s := range program {
		c := s.fn.name
		for _, sym := range s.argsyms {
			c += " " + string(sym)
		}
		str += fmt.Sprintf("%v: %v = %14v # \n", i, s.outsym, c)
	}
	return str
}

func printProgramAndValues(program Program, vm ValueMap) {
	str := "\n"
	for i, s := range program {
		c := s.fn.name
		for _, sym := range s.argsyms {
			c += " " + string(sym)
		}
		r := vm[s.outsym].value
		str += fmt.Sprintf("%v: %v = %14v # %v \n", i, s.outsym, c, r)
	}
	fmt.Println(str)
	// InfoLog.Println(str)
}

func printProgram[T Program | UncheckedProgram](program T, loggg LogType) {
	if loggg == Info {
		InfoLog.Print(formatProgram(program))
	} else if loggg == Errr {
		panic(formatProgram(program))
	}
}

type GenSym struct {
	idx int16
}

func (g *GenSym) gen() Sym {
	s := Sym(fmt.Sprintf("v%v", g.idx))
	g.idx += 1
	// fmt.Printf("g.idx = %v \n", g.idx)
	return s
}

func (g *GenSym) genUnique(existing SymSet) Sym {
	for {
		s := g.gen()
		ok := existing.Contains(s)
		if !ok {
			return s
		}
	}
}

type Catalog struct {
	syms     map[SymLine]Type
	syms_inv map[Type][]SymLine
}

func NewCatalog() Catalog {
	return Catalog{
		syms:     make(map[SymLine]Type),
		syms_inv: make(map[Type][]SymLine),
	}
}

func (cat *Catalog) add(s_new Sym, t_new Type, line uint16) {
	symline := SymLine{sym: s_new, line: line}
	t, ok := cat.syms[symline]
	if ok {
		assert(t == t_new, fmt.Sprintf("The sym %v has conflicting types %v and %v.", s_new, t, t_new))
		InfoLog.Printf("The Sym %v is being overwritten.\n", s_new)
	} else {
		cat.syms[symline] = t_new
		a := cat.syms_inv[t_new]
		a = append(a, symline)
		cat.syms_inv[t_new] = a // Must re-assign to map in case `a` has been moved
	}
}

func buildTypeGraphOneStep(typeset_start *Set[Type]) (typeset *Set[Type], funcset *Set[string]) {
	typeset = NewSet[Type]()
	funcset = NewSet[string]()

	for _, fn := range fn_library {
		paramset := NewSet[Type]()
		for _, ptype := range fn.ptypes {
			paramset.Add(ptype)
		}
		if paramset.Difference(typeset_start).Size() == 0 {
			funcset.Add(fn.name)
			typeset.Add(fn.rtype)
		}
	}
	return typeset, funcset
}

// Return a minimal set of funcs which must be called in order to make the target func callable.
func shortestFuncPath(path []string, target_func string) []string {

	// assume input path is valid
	// get types buildable from funcs in path
	// forall funcs f in Library callable with these types
	// true_p = [], true_len = -1
	// if f == targetfunc: return path
	//   else
	//   p = shortestFuncPath(path + f)
	//   if true_p = [] OR len(p) < len(true_p): true_p = p
	// return true_p

	fns_path := NewSet[string]()

	typeset := NewSet[Type]()
	for _, fn_name := range path {
		fns_path.Add(fn_name)
		fn := fn_library[Sym(fn_name)]
		typeset.Add(fn.rtype)
	}

	_, fns_callable := buildTypeGraphOneStep(typeset)
	fns_callable = fns_callable.Difference(fns_path)
	bestpath := make([]string, 0)
	for fn_name := range fns_callable.emap {
		if fn_name == target_func {
			return path
		}
		newpath := append(path, fn_name)
		newpath = shortestFuncPath(newpath, target_func)
		if len(bestpath) == 0 || len(newpath) < len(bestpath) {
			bestpath = newpath
		}
	}
	return bestpath
}
