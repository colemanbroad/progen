Random programs composed of small ints, +, and * are probably building numbers that
are somewhere in between normally distributed (from the `+`) and log-normal (from the `*`).

# Todo

-[x] try [zerolog](https://github.com/rs/zerolog)
-[x] try [slog](https://pkg.go.dev/log/slog@master#Handler)
    Send the json output directly into sqlite ?
    To duckdb ?
The answer is none of the above, because it's all too slow.
Just keep the data in memory and serialize it in bulk at the end.
For a long running process, just serialize intermittantly in bulk.

-[ ] What is the actual distribution of numbers sampled and sent to power_of_two? Is it an interesting distribution?
-[ ] add optimization
-[x] add pitsworld

Apache e-charts has tons on nice effects but is a charting embarrassment

-[ ] Allow for generic passing of stream of data with mapping from column names to plot dimensions.
-[ ] Multiple Streams are allowed.
-[ ] Allow for 2nd entrypoint where you add your own `run(Context)` function and have access to arbitrary canvas drawing tools.

-[ ] parallize rollouts? will this actually help or will we be screwed by reward array contention?
-[ ] fuzz a single pitsworld map instead of generating new ones randomly?


# Tips

profile you running program with:
`go tool pprof -http=":8080" http://localhost:6060/debug/pprof/profile?seconds=10`
to profile the CPU for 10s and then view the results in a web view.


# Todo

-[ ] optimizer
-[ ] wordle

# Program Gen

Programs are sequences of Statements.
Statements are tuples of (outsym, fn, argsyms).
Symbols map to Values which are created during Evaluation.

Evaluation proceeds top to bottom.
~Evaluation happens in phases.~

Q: How are we able to avoid the ditinction between Op and Call/Fn/LitNum/LitStr ?
A: LitNum/Str enable passing in arbitrary data directly to programs. We can't do that.

Q: Can we combine Programs together?
A: Yes, see e.g. rewire, mutate, shuffle, insert, delete, interleave.

Maybe we can put Programs though an analysis() step that returns an AnalyzedProgram
 which tells us which symbols and FuncDefns are missing definitions? 
Then we can manipulate the body of the Program without worrying that temporarily invalid Piece.inputs and returns.
Only an AnalyzedProgram without any missing Syms or Fns is ready to be Evaluated. 

This would be a mistake if we find ourselves re-anlyzing more frequently to the point where it becomes
a real cost.


## Program Gen Constraints 

What invariants are implicit in the data and what boundaries must maintain them?

1. The Program sent over the wire must be Valid (nothing Missing) (This is our responsibility).
2. The Values sent back over the wire for the Evaluated Program must be Valid.

OK, but what about internal boundaries? 
Is the program allowed to be arbitrarily invalid at every internal boundary?
No we shouldn't allow that, because constant checking would be inefficient!

Is there a way to make Validity true by construction? In it's current form this constraint must be checked. It is implicit.
Perhaps, but it would probably require a highly non-intuitive encoding to enforce that invariant.

So the first thing I want to try is putting an assertValid(program) at the end of any function that returns a Program.
This way whenever we call a func that returns a program we know that it'll be valid... 
NOTE: we can't hold the invariant that all Program params are valid, because this precludes defining funcs like checkValid(program) ! 

We could define two separate types ProgramChecked and ProgramUnchecked that tracks our knowledge about the state of the Valid invariant?
This way we can avoid checking unnecessarily, and we won't forget to check when it IS necessary!

and even if we take a program parameter
we can set another invariant that all programs passed 

# Functions as Values

Q: How can we incorporate functions-as-values into our Sym, Value, FnCall, Statement scheme?
A: Two options:
    1. Add FnDefn that, upon execution, adds new Func to the Library. FnDefn holds a subprogram / piece, calls Eval recursively, and returns a value?
    2. Skip FnDefn. If we want to reuse a piece of a program then we have to copy/paste it around at compile time.

Ideas:
A function can be returned from another function and associated with a symbol.
A function can be passed as an argument to a function through it's arg symbol.
If we want to use a function `f` that we've just defined then we pass it into `Call(f, args)`.
Usually the Function (in this case Call) defines the types of input and output.
But in this case the output type of the function would depend on the type of f (an arg type).

Q: How can we do type checking if we allow for things like `Call`? 

Call is a generic function, so we'd have to make the Type Checker smarter.

Call.rtype depends on Call.argtypes.
So we get Call.argtypes from lookup their syms in local catalog.

This will require a way of determining the argument types from their Syms.
Usually we get the types of Syms from the Statement.fn.ptypes and rtype.
If the first argument to Call is a FuncDefn then we're OK... But actually it should really be a Sym which is the 
outsym of a Statement with the fn:FuncDefn called "Define" that defines functions dynamically.

Define can't be implemented as a function.
It has to be implemented as something that surrounds some existing piece of code, figures out it's type, picks a return value...
If this is going to be done dynamically then we're going to have to enable it in the interpreter... Which we probably haven't done yet?

function and the remaining arguments are passed into the function then 
another lookup on the types of the argsyms

# Branching

SSA form defines `phi nodes` that merge different branches of control flow together, thus we can define e.g. loops 
which update a value in place instead of 

e.g. 

x1 = a
x2 = b
x3 = if c(x1,x2) x2 else x1
This would be written

IF(fn1, fn2, fn3) {
    if fn1() then return fn2() else return fn3()
}


# Rewire, Mutate, Shuffle

Imagine I have an API for sequences that changes the order of elements.
I could write a function `permute(i,j)` and that would be enough to sample all permutations. 
A dist over program length would determine a dist over permutations.
But I could also add `permute_triple_cyclic(i,j,k)` which doesn't give us access to any additional permutations, but does change the distribution!
In fact we could add dozens of these functions, and it might be even less obvious that we're duplicating work.
E.g. `type Permuter struct {current_index: u32}` and `Permuter.set(n)` and `Permuter.swap(n)` creates a new, different distribution over perms.

Think about the Algebra of +, *, ^ and iterate(x0, f, n) : T, (T -> T), Int -> T which applies f(f(f(... n times total ...(f(x0))))) n times in a row.

--------

Fundamentally there are three kinds of changes we want to make to programs: rewiring, mutating and shuffling.

1. (mut) Pointwise mutations of values - keeping DAG structure the same but replacing Values with new values of the same type. Alt: replacing value in LitNum or LitStr.
2. (mut) Pointwise mutations of functions - replace Funcs with new Funcs of the same func type
3. (rewiring) Keep the same funcs and outsyms but change the argsyms in a type-consistent way 
3. (mut + rewire) It's possible but unlikely to randomly mutate and rewire in a way that is equivalent to shuffling.
4. (shuffling) Statements while optionally preserving argsyms and funcs
5. (shuffling + rewire) Statements while preserving funcs but allowing rewiring argsyms (shuffling + rewire - mut)
6. shuffle + mut + rewire
7. insert / delete

Shuffling matters for stateful system. It determines order in which ops get executed.
Rewiring matters for pure functional systems. It determines dataflow.
Mutation + Insertion + Deletion determines the starting program. It determines how we explore program space.

---------------------------------------------------------------------------------------------------

NOTE: We can't shuffle arbitrarily without allowing some insertions: we can't put +(a, b) at the top without inserting e.g. a = 1, b = 9;

`rewire` syms keep statement order and funcs
`mutate` funcs (this includes LitNum and LitStr funcs. In general funcs are things that map symbols to values or transformations.)
`shuffle` statement order

Note that these three classes of operations are not orthogonal.
It's possible to perform a mutation + shuffle and end up with the same program. 
Also, the exact order of statements is often totally irrelevant to program behaviour. 
    It only matterss for funcs with side effects, although code that looks functional may in fact have deep, hidden side effects e.g. on cache that can determine later program values (e.g. timings)

What about `insert` which adds sub-programs to the body of a program and then wires them together.
At it's simplest this is just replacing f(a,b) with c = 2*a; f(c,b).

`insert` is one way of recombining two programs.
`interleave` is a kind of homologous recombination. 
In general when combining two programs we only need to make one kind of decision which is how to interleave.
Maybe we wanted to start off with some shuffling?
Then the questions is how to connect the two programs?

1. shuffle
2. interleave (insert/delete)
3. mutate
4. rewire

Can rewire conditional on the current wiring (e.g. given f(a,b) rewire to f(c,b) s.t. c is the first outsym with type(a) _before_ a appears.)

Rewiring that only pays attention to the current order of funcs in the program gives us total freedom to experiment, however not all 
func orderings are valid. 

NOTE: we can move all arg-free funcs to the start of the program WLOG ? No this is only true if side-effect-free (pure).
Knowing which funcs are pure would allow us to cut down on the number of unique programs, giving us an equivalence class of programs.

shuffle:
    for each line identify the line-no of each argsym.
    pick a random line and move it to a random location with line-no > max across argsyms.
    repeat

mutate:
    get the functype of a random line
    replace that func with a random func from Library of same functype

Combo Operations

"pruning":
mutate + rewire + delete
    choose a random func in program with rtype T
    cut all wires leading into func and replace with a random library func of type () -> T
    upstream funcs can now be deleted
    
"growing":
sample + mutate + insert
    sample a new, small program P2
    insert at start of P
    choose a random func f in P with rtype T 
    cut all wires into f 
    prune a random location in P1

      
# Conditional Generation

Imagine learning a map[func -> func]f32 that learned to associate a strength of connection between funcs. 
    funcs f1 -> f2 are connected if the outsym of f1 is an argsym of f2
    this connection is a wire
    a program contains many wires and gets an f32 reward, then we update the weight associated with f1->f2 based on the reward
    when constructing programs we can sample conditional on these connection strengths.

Hierarchical program construction

have to decide on 1) the shuffling and 2) the wiring
Homologous recombination maintains the per-program order, and interleaves roughly evenly.
Dumb concatenation is the simplest kind of interleaving. 

If we keep the program size fixed and the library fixed then there are only a finite number of programs constructable across all rewirings, mutations and shufflings.

It changes program length!

Note there are also operations like `uniqueifySymbols` which have NO effect on program behaviour.
In general any permutation of Syms will. 

So let's create functions based on these three types of things.

# Problems with Golang

Numeric types: `int` is not a generic int type, but rather a system-specific one. There is no system-specific float type, so we must specify float32 and float64.
We can't specify arbitrary int/uint types like i17 or u24.
~Can't define a generic function var like `var f func(T)int`~ Not true.
I have two equivalent funcs `func f(x []any)` and `func g[T any](x []T)`. For `x : []f32` calling `g(x)` works but `f(x)` doesn't compile.
Can't send signals to goroutines unless we pass around an explicit Context.

But also, [Go is my hammer](https://news.ycombinator.com/item?id=41223902).

# Relevant Work

Program Synthesis

- https://cs.nyu.edu/~davise/
- https://cs.nyu.edu/~davise/rck/intro.pdf
- https://deoxyribose.github.io/No-Shortcuts-to-Knowledge/#learning-as-probabilistic-program-synthesis
- https://evanthebouncy.github.io/program-synthesis-minimal/
- https://www.reddit.com/r/MachineLearning/comments/y378kk/p_a_minimalist_guide_to_program_synthesis/

- [Go interpreter for codegen eval](https://github.com/traefik/yaegi)

# Go testing and sometimes assertions

I wrote a testing function that had a branch. 
I wanted to know if both branches were hit and how often, but the _test.go files themselves are NOT under coverage. 
It turns out to be possible to move the function in question inside a file that _was_ under coverage, but it feels like the wrong place for test code.
I would like to have sometimes assertions, or coverage that keeps track of the percentage of my assertions that were hit and how often.
I can probably parse coverage.out, look specifically for `assert_sometimes(name, bool, message)` and save the results in a local db whenever we run testing, then e.g. plot them over time.
Code coverage (line level) is already quite good at showing you which lines you never hit.
But do we really want an `assert_sometimes()` in the code? 
So we want to know not just "did the assertion get called?" But also "what was the result? did it pass or fail?"
I think we can write to a DB at the end of the tests in test.Main?
We can parse the text to build a mapping from source lineno -> assertion name, then we can 

# Minimization

"Smaller" P have fewer lines (and therefore build smaller datastructures).
They should be "simpler" programs.

|P'| <= 2^P.len()

We want to search through contractions of the DAG from the bottom upwards.
The root nodes (definitions) are fixed, but we can prune the tree from other places.
Every node (line of the program) has a depth in the DAG (distance to nearest root).
A cut of this DAG is a valid program if all nodes' dependencies are satisfied.
We want to search through cuts.
Typically minimization is done from the "datastructure" perspective.
First consider a datastructure like a List. When we pass this list into f() it throws/returns a non-nil Error/nonzero exit code.
We want to know _why_ it threw the error, and we're suspicious that it has something to do with the List itself,
not because of our underlying environment is buggy and random.
Let's assume that if we try the same one again it will Error every time.
So what kind of list should we try next?
Maybe cut it in half? Why do our instincts tell us this is a good idea?
Maybe first we try it with the empty list or a list of size one.
That's because our programmer instincts are telling us "It doesn't work at all." Is the most likely thing.
Then we try it with different shapes.
Then with different numbers in different spots (let's assume we're working with a strongly typed language and f : [Num] -> error
Ok we try cutting it in half and playing with the numbers and eventually find that lists which have a [..., 0, 0, ... ] in them are bad.
Great. We solved by writing a minimizer for Lists.
What about the next g : Tree[Num] -> error? Do we need another custom minimizer?
What about f : SceneGraph -> error ?

What about finding a generic form / Piece describing the entire manifold of reward?
We can trim a program DAG down to the smallest subDAG that maintains T, but the result represents just a single point in program space.
The path we took to get there from P0 -> P1 -> ... -> PN may have strayed off the manifold of T at times,
 but it probably contains a subsequence of programs P'0 .. P'N that maintain T.

----  ----  ----  ----  ---- 

Imagine we are exploring Programs and we find one with a certain property T.

Let

P := The set of all programs.
P' := The subset of Programs that have T.
P'' := The subset of P' that have minimal length.

We may believe that P' have T because P'' have T, i.e. that all the extra lines in P' that differ from P'' are superfluous.
This may not be the case.
E.g. 

db : DB = connect_to_db()
s1 : string = new_write_query()
resp : HTTPResponse = submit_write(db, s1)
check_response(resp)

*How can we identify P'' in the fewest rollouts?*

Prune the DAG.
A program is valid if all arguments are defined.
A cut of edges that separates p into two components have still have <= ONE valid program.
A disinterleaving of p into two programs that doesn't  

-[ ] Does the udp_echo_test work? 
    With the bbfuzz campaign? 
    Did it introduce faults?
-[ ] Can we compare bbfuzz campaigns with different settings on the real curriculum? 
    Is it easy to add new things to the comparison? 
    Did we find any bugs or introduce any faults?

1. Query that returns sequence from tip to root. sort them. do some sequence analysis.

# Todo

-[ ] Add information to the Catalog that distinguishes between StdLib and Customer Fragments.

I have to use cusomer-research (production-environment?) for Curriculum, 
but I have to use customer-antithesis 

--- --- --- --- --- --- ---  --- --- --- --- --- --- ---  --- --- --- --- --- --- ---  

I want to know the statistics of connections. 
What is the avg distance between outsym and it's use as argsym?
What is the distribution?

I want to see how the wiring mode and learning affect Reward trajectories.

I want to be able to run a large number of trials, looping over hyperparams.
This will require full control of init() time.


# Peano Experiment 001

Plot x axis is the wiring depth to root of a symbol using a library that consists only
of `zero` and `succ(x)`. The y axis shows the count of symbols with that depth summed
over 1000 generated programs. Color shows the length of a program in lines (values are
10, 20 .. 100). Each panel shows a different value for the exponential decay rate used
in sampling an argument to  the `succ(x)` function. An increase of 1.0 in the rate
means probability decreases by `e`. The bottom panel (`zero`) shows flat sampling across
possible arguments.

The figure has two variants: when "cheating" we turn off the `zero` value after the first
Op is sampled, leaving only `succ(x)`. This is able to dramatically extend the depth of
the wiring distribution esp for decay rates 0.5 and 1.0.

# Program Minimization

Is like Binary Search in that it is O(log n) but works on DAGs instead of Sequences.
Source Code describing Dataflow has this structure.
Version Control History (VCH) has this structure.
You can search through subsets of the DAG of a program for the minimal program that
satisfies some bug/property/behaviour. And you can do the same to VCH, looking for the
minimal set of changes to the repo. The thing is... this algorithm is only efficient if we
make certain assumptions about the structure of these minimal programs / changesets. It
only works if one of the following conditions are met:
1. The change is a single line / changeset / node in the graph.
2. The change is a compact set (connected component) of adjacet lines / changesets / nodes.
3. The behaviour is not over-determined. There may be multiple nodes that cause the
behaviour, but we won't detect them until we remove the last one. Then we'll have to
go back and try the full program - the final line. We'll do this and find that the bug
remains, and the search process begins again. The result of this search process should be
a set of DAGs with a single node colored to represent that removing that node toggles the
condition. E.g. consider the program testing condition c. Good test

Non monotonic

    [1] x := 0
    [2] x += 1
    [3] x += 1
    [4] x += 1
    [5] x += 1
    [6] x += 1
    [7] c := x%2 == 0

The condition changes if we remove ANY of lines [2..6], but not if we remove any PAIR of those lines. But yes if we remove any TRIPLE of lines, etc.

Overdetermined

    [1] s := init_server()
    [2] kill_server(s)
    [3] kill_server(s)
    [4] kill_server(s)
    [5] kill_server(s)
    [6] c := is_server_alive?()

Good test

    [1] x := 1
    [2] x *= 2
    [3] x *= 3
    [4] x *= 5
    [5] x *= 7
    [6] x *= 11
    [7] c := x%3 == 0

The right thing for a minimizer to report in these three cases are:

Non monotonic:

  Non monotonic tests can either (a) stop greedily as soon as `c` changes. (b) continue
  on until `c` changes BACK, then report that the system is non monotonic. (c) continue
  on until we have identified the smallest possible program + delta that makes the initial
  change, in this case there is one possible base (lines 1 and 7) and five red lines
  (lines 2..6).

          x := 0
    [RED] x += 1 
          c := x%2 == 0

Overdetermined

    [1] s := init_server()
    [2] kill_server(s)
    [3] kill_server(s)
    [4] kill_server(s)
    [5] kill_server(s)
    [6] c := is_server_alive?()


I don't see how looking for the minimal DAG is possible in general.
We believe the structure is a DAG... but actually it's just the nodes that matter and the
edges can be rewired. If two nodes produce the same value (e.g. the number 3.0) it doesn't
matter which of them feeds that value as an argument to a downstream function. If we make
a change to the repo README it's not actually a hard requirement of the behaviour of the
program after that changeset.

The funny thing about the DAG is that it might not actually show HARD dependence, i.e. it
may be possible to cut out a node (src line or changeset), maybe tie the graph back
together and have things work fine.

!!! I found [this post by Regher](https://blog.regehr.org/archives/527) on generalized delta debugging.
That led to [bugpoint](https://llvm.org/docs/Bugpoint.html) mentioned by Chris Lattner.

# Genetic Experiments

These ran 20 campaigns a GeneticOptimizer of 1000 rounds of 20 programs with different mutations strategies and length ... 
But now I want to do the standard sampler with depth + power-of-two experiment. 20 campaigns 1000 rounds 20 programs of length 20 (default params).

Initial results suggest decay has little to no effect.


-[ ] Add full program text parsing + eval (no js syntax. boring syntax).

# Depthmap

I can create the map during program construction and afterwards.
What should I do with it?
- Use it to prefer deeper syms
- Use it to study how depth relates to diverse values

    If I have a depthmap AND valuemap I can aggregate them. What is the number of distinct
    values available as a function of program depth?  

Ok the results are in... The fraction of unique value increases monotonically with depth. The count of total unique values increases to a point, then decreases again.
The total depth distribution decays exponentially in each case: Flat sampling, exponential rate 0.1, exponential rate 1.0.
1000 programs length 100, and the depth 100. The tail of the exponential is greatly extended by changing the decay rate! 
Why does it work here but not for the Peano system?
I think because the Peano system had *half* of it's fragments (Zero / succ) would reset the depth to zero. If we build programs naively and we're interested in depth then we're going to be constantly running into the issue that the depth is reset to zero on every line with 50% probability. We can fix this by:

- removing Zero from Library and 
    - either 
        - adding it as the first line of the program [p(zero) = 1/n_lines]
        - adding it as a Value

- We can have a set of Values in the same way that we currently have a Library.
- Values live in the ValueLib map[Sym]Value and the Sym keys can be used as _arguments_ to FnCalls directly, without first being added to the program body as an outsym. 

---

What are we observing/measuring?
What are we controlling?

    - variables controlling program generation
    - variables controlling campaigns

What is the fragement set?
(What is the workload?)
How are we generating programs?

# Basic Math + Measure Depth vs Diversity

# Network model

Nodes send messages (values) at a time (float64) to a destination.

Nodes are state machines.
All comms are to a single NetworkManager which keeps messages around, maybe delays them. Maybe sends them on.
All nodes and the network manager have internal monotonically increasing clocks.
All message send and arrival times are recorded. Time increases in random intervals on each Node/NetwormManager.

The NetworkManager keeps the state of a bunch of Links.
Links can be Up, Delay, Hold or Drop. When a link is Delay and it receives a message we decide immediately on the
new time of arrival of the message and the Link's clock increases by delta t. When a link is Hold it keeps the
Messages in a buffer and increments time and changes state. Only later do we decide if/when to send the messages
in the buffer. 

Faults are just ways of changing link state either as individual Links or in groups.
Faults may accumulate?

# Todo

-[ ] generics?
-[ ] funcs as args and return values?
-[ ] closures (syntax for defining a func -- fn as return value -- that
    references symbols in the current scope and function in the lib
    without passing them into the arguments of another fn.)
-[ ] loops
-[ ] branching
-[ ] update in place (overwite sym) 

# Experiments Run So Far

Fragment Sets
- Add,Mul,One
- Add,Mul,LShift,One
- Succ, (Zero as value)
- Succ, Zero
- Add,Mul,PowerOfTwo
- Add,Mul,LShift,PowerOfTwo

- Pitsworld?
- Network Controller
- Network Workloads (echoserver, basic Raft?) 

Analyses (Margins) (Projections) 
- Depth (min distance to root value)
- Value

- wire.sql
    We used different variants of the Peano fragments.
    We played with wiring methodology (decay) and program length.
    We measured `depth` and `count`.
    We wanted to know if we could control the depth distribution.
    We determined that escaping exponential distributions is hard unless you use Zero as a Value.
- p2.sql
    We used the basic math + power-of-two fragment sets.
    We played with wiring methodology (decay).
    We measured `value`, `reward`, and `time`.
    We wanted to know if
    - we could make broader depth distributions.
    - broader distributions would find more powers of two.
    We determined that we CAN extend the exponential depth distribution with wiring decay.
    And this DOES have an effect on the number of unique powers of two (+4) and the cumulative reward (+60%).
- genetic.sql
    We used the basic math + power-of-two fragment sets.
    We played with ???
    We measured power-of-two (value, reward, time) AND wiring (depth, count).

- The p2 experiments can be rerun without touching the DB, just printing tables.

What experiments _should_ we run?
Our major, high-level questions are about efficiently finding bugs across all software.
But we're starting off with distributed systems.
And the most common dist sys bugs (and checkers) center around the orderings of events.
Formal methods try to prove that communicating state machines maintain invariants across
all possible event orderings.
Why are multiple orderings possible?
1. Messages travel across a network and experience different random delays / disappear.
2. Each machine proceeds at an indepent pace.

Which of these is the greater source of variation?
Each machine is processing concurrently while messages are being sent.
The system as a whole may have a lifetime in days, months or years before a full-system reset.
A single packet can take (best case) ~200ms to travel to the other side of the world.
A photon traveling through vacuum/glass takes 134/200ms to traverse the circumference. 
But on faulty networks it can be arbitrarily long. [source for a real distribution here?]

Each computer is very precisely a state machine, although the granularity of states
consisdered by dist-sys checking tools and the way "state" is used in dist-sys algorithms
e.g. Raft is much more coarse. So while the machine may change states at a rate 1e9/s in
one sense, it may maintain a single dist-sys state e.g. (leader, follower, voter) for arbitrarily long times!

--- --- --- --- --- --- --- --- --- --- --- --- --- --- --- 

# See delta.md for program minimization

[delta.md]

# Hierarchical models

I think the only way that I can think of at the moment to achieve what I believe Dave wants, and what will give us interesting programs across an arbitrary set of margins is the following:

- We enumerate or sample a sequence of meta-programs
- For each metaprogram we sample a program []Stmt in a particular way, e.g. using a subset of fragments, with weird wiring, or using a certain pattern. Each sub-program is concatenated to form the full program.

This makes a lot possible, but it makes everything more complicated.
First, we have to think of and implement a set of meta-fragments or patterns.
Then we have to implement a mata-program sampler which enumerates meta programs and samples them.

How to do credit assignment?

What kind of grammar for metaprograms? Are all pieces independent? Is there a state machine that tells me which pieces are allowed? Or is there a stack-like context that tells me which pieces are allowed? Can the pieces set/unset arbitrary state?
Does each piece return a Program, or are there other types of pieces that *only* set context? E.g. imagine taking the product
between Pieces that return state, pieces that control the set of available fragments, and pieces that control wiring.
Are these all independent? 


# Annoyances with work setup

nixos: The default version of Delve is TOO OLD and doesn't work right with my Go version. Maybe I can just `go install` it?
