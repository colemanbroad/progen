package main

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
)

type Value struct {
	value any
	name  string
	vtype Type
}

type FnCall struct {
	value  any
	name   string
	ptypes []Type
	rtype  Type
}

func (fn FnCall) String() string {
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
	fn      FnCall
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

var Library map[Sym]FnCall

// func init_library() {

// 	// functions
// 	Library = make(map[Sym]FnCall)
// 	addPeanoLib()

// 	addPowerOfTwo()
// 	// add_pitsworld()
// 	// basiMathLib()

// 	// f = func() string { return fmt.Sprint(rand.Intn(10)) }
// 	// addFn(f, "rand_s", []Type{}, "string")

// 	// f = func(a string) string { return a + a }
// 	// addFn(f, "repeat", []Type{"string"}, "string")

// 	// f = func(a, b string) string { return a + b }
// 	// addFn(f, "cat", []Type{"string", "string"}, "string")

// }

func addBasiMathLib() {
	var f any

	f = func() int { return rand.Intn(10) }
	addFuncToLibrary(f, "rand", []Type{}, "int")

	f = func() int { return 1 }
	addFuncToLibrary(f, "one", []Type{}, "int")

	f = func(a, b int) int { return a + b }
	addFuncToLibrary(f, "add", []Type{"int", "int"}, "int")

	f = func(a, b int) int { return a * b }
	addFuncToLibrary(f, "mul", []Type{"int", "int"}, "int")

	f = func(a, b int) int { return a << b }
	addFuncToLibrary(f, "<<", []Type{"int", "int"}, "int")

	f = func(isPos bool) float64 {
		x := rand.Float64()
		if !isPos {
			x *= -1
		}
		return x
	}
	addFuncToLibrary(f, "samplePosOrNeg", []Type{"bool"}, "f64")
}

func addPeanoLib() {
	var f any
	f = func(a int) int { return a + 1 }
	addFuncToLibrary(f, "succ", []Type{"int"}, "int")
	f = func() int { return 0 }
	addFuncToLibrary(f, "zero", []Type{}, "int")
}

func evalStatement(stmt Statement, locals ValueMap) {
	var r any
	var f any
	// type erase the function we're calling
	f = stmt.fn.value
	g := reflect.ValueOf(f)

	// then erase the args, (but haven't they arleady been erased?)
	// TODO: avoid alloc: make this fix-size array?
	args := make([]reflect.Value, 0)
	for _, a := range stmt.argsyms {
		args = append(args, reflect.ValueOf(locals[a].value))
	}

	// fmt.Println("f, g, args", f, g, args)
	// printProgram(Program{stmt}, Info)
	// return val cast to `any` type
	r = g.Call(args)[0].Interface()

	// Check that the types are correct
	// Here's where we can introduce the logic of "any" types?
	rtype := Type(reflect.TypeOf(r).Name())
	// if rtype != stmt.fn.rtype {
	// 	ErrorLog.Fatalf("Val type doesn't match fn rtype: %v %v", rtype, stmt.fn.rtype)
	// }

	val := Value{value: r, name: fmt.Sprintf("%v", r), vtype: rtype}
	locals[stmt.outsym] = val
}

func evalProgram(program Program) (values ValueMap, reward float64) {
	r0 := Reward_total
	locals := make(ValueMap)
	// TODO: How are we going to allow for ZeroValue in Programs with the same semantics as Values
	// in Rust GenTactics?
	if cheating == ZeroValue {
		locals["Zero"] = Value{
			value: 0,
			name:  "Zero",
			vtype: "int",
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
	fdef := FnCall{
		value:  f,
		name:   name,
		ptypes: ptypes,
		rtype:  rtype,
	}
	Library[Sym(fdef.name)] = fdef
}

func sampleFuncFromLibrary() FnCall {
	keys := make([]Sym, 0)
	for k := range Library {
		keys = append(keys, k)
	}
	k := keys[rand.Intn(len(keys))]
	return Library[k]
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
func sampleProgram_fromFragmentLib(sp SampleParams) Program {
	gensym := GenSym{idx: 0}
	program := make(Program, 0)

	// TODO: use the Catalog with lineno info to have more control over initial wiring
	local_catalog := NewCatalog()
	// idx := 0
	// n := rand.Intn(20)
	n := sp.Program_length
	line_no := uint16(0)

stmtLoop:
	for len(program) < n {
		// fmt.Println("Prog = ")
		// print_program(program, Info)
		var f FnCall
		switch cheating {
		case Normal:
			f = sampleFuncFromLibrary()
		case ZeroValue:
			m := len(program)
			if rand.Float32() < 1.0/float32(m) {
				stmt := Statement{
					fn:      Library["succ"],
					outsym:  gensym.gen(),
					argsyms: []Sym{"Zero"},
				}
				local_catalog.add(stmt.outsym, stmt.fn.rtype, line_no)
				line_no += 1
				program = append(program, stmt)
				continue
			} else {
				f = sampleFuncFromLibrary()
			}
		case ZeroOnlyOnce:
			if len(program) > 0 {
				f = Library["succ"]
			} else {
				f = sampleFuncFromLibrary()
			}
		}
		stmt := Statement{
			fn:      f,
			argsyms: make([]Sym, 0),
			outsym:  "void_sym",
		}
		// sample a random sym of the appropriate type from the Program above for each argument
		if len(f.ptypes) != 0 {
			for _, ptype := range f.ptypes {
				symtypes, exist := local_catalog.syms_inv[ptype]
				if !exist {
					continue stmtLoop
				}
				var n int
				if sp.Wire_nearby {
					// Exponentially distributed sampling.
					m := len(symtypes)
					n = m - int(TruncatedExponentialSampler(nil, sp.WireDecayLen, float64(m))) - 1
				} else {
					// n = (idx * 103823) % len(symtypes)
					// idx = idx + 1 // TODO: change the math here to control the wiring
					n = rand.Intn(len(symtypes))
				}
				stmt.argsyms = append(stmt.argsyms, symtypes[n].sym)
			}
		}
		stmt.outsym = gensym.gen()
		// fmt.Printf("add() in sampleProgram() sym = %v \n", stmt.outsym)
		// fmt.Printf("gensym = %v", gensym)
		local_catalog.add(stmt.outsym, stmt.fn.rtype, line_no)
		line_no += 1
		program = append(program, stmt)
	}
	return program
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

	for _, fn := range Library {
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
		fn := Library[Sym(fn_name)]
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
