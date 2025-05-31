
#let horizontalrule = line(start: (25%, 0%), end: (75%, 0%))
#show figure.caption: it => text(it, fill: gray, size: 9pt)
#set heading(numbering: "1.")
#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node

#let question(bod) = [
  #v(1em)
  #align(center)[#bod]
  #v(1em)
]



= Program generation methods

- Forward
- Forward + tail weight (a generalization of Forward)
- Backward
- Dataflow-first
- Margin-guided (aka "tree-based")
- Two-level (this is really an idea for program combinators, not direct generation.)
- Trait-based (customers impl traits with known semantics)

_Each approach leads to a different distribution of programs._

== Forward sampling <forward>

Program generation is currently performed one line at a time ("forward") by

+ sample a random fragment from the catalog
+ sample a random tuple of existing, appropriately typed syms as arguments
+ reject the fragment if the arguments don't exist.

_Does this give a flat distribution across the set of possible programs?_

No, it doesn't. This process leads to a bias where fragments that are easier to insert

== Forward + tail weight

Fragments are sampled evenly from the catalog,
but wiring has a preference for syms defined more recently.

= Common data structures in program generation

Simplifications of `Program` and `Catalog` arranged from rich to poor.

== Program

A program introduces: `Sym`, `Statement`, and `Function`.

// == DAG

// A TacticsLang program can be represented as a DAG
// where each node corresponds to a statement defining a new sym.
// This representation is clean in TacticsLang as the language doesn't allow defining multiple syms per statement.

// #let nodes = (
//   `v1 = f()`,
//   `v2 = g()`,
//   `v3 = h(...)`,
//   `v4 = i(...)`,
// )
// #let mydia = diagram(node-stroke: 0.1em, node-corner-radius: 0.3em, {
//   let a = (1 / 2, 0)
//   let b = (0, 1 / 2)
//   let c = (1, 1)
//   let d = (1 / 2, 1.8)
//   node(a, nodes.at(0))
//   node(b, nodes.at(1))
//   node(c, nodes.at(2))
//   node(d, nodes.at(3))
//   edge(a, c, "-|>")
//   edge(b, c, "-|>")
//   edge(c, d, "-|>")
//   edge(b, d, "-|>")
// })

// #let fig = align(center + horizon)[
//   #stack(
//     dir: ltr,
//     ```python
//     v1 = f()
//     v2 = g()
//     v3 = h({v1, v2})
//     v4 = i({v2, v3})
//     ```,
//     text(mydia, 7pt),
//     spacing: 3em,
//   )]
// #let cap = [
//   The program on the left corresponds to the DAG on the right _evaluated in the order_: `v1 v2 v3 v4`,
//   but the order `v2 v1 v3 v4` would also be valid. The specific tuple of arguments for a node is also hidden in this view,
//   all we can say is that the function depends somehow on the set of inputs flowing into the node.
// ]
// #figure(fig, caption: cap) <dataflow-0>

== E-Graph

The E-graph is uniquely defined by the Catalog,
and has a node for every element of the Catalog and a box drawn around nodes with the same return type
and a numbered, directed edge out from each node pointing to a box of appropriate type corresponding to the arguments to that function.

== TF-Graphs

The TF-Graph is also uniquely defined by the Catalog.
The initial catalog is a `map<Name, Function>`,
where each fragment is a type with a single arrow
#footnote[Let's ignore higher-kinded and types (HKTs) and generics for now.]
mapping a tuple of input types to a single output type.

The TF-Graph for the catalog `{0, 1, +, True, False, =}` looks like this:

#let dia = diagram(node-stroke: 0.1em, node-shape: circle, {
  let a = node((0, 4 * 0.0), "0")
  let b = node((0, 4 * 0.2), "1")
  let c = node((0, 4 * 0.4), "+")
  let c2 = node((0, 4 * 0.6), $=$)

  let d = node((0, 4 * 0.8), "True")
  let e = node((0, 4 * 1.0), "False")
  // let f = node((0, 4 * 1.2), $=_b$)

  let g = node((1, 0.2), "Num")
  let h = node((1, 2.8), "Bool")

  a
  b
  c
  c2

  d
  e
  // f

  g
  h
  edge(a.value.pos.raw, g.value.pos.raw, "-|>")
  edge(b.value.pos.raw, g.value.pos.raw, "-|>")
  edge(c.value.pos.raw, g.value.pos.raw, "-|>")
  edge(c2.value.pos.raw, h.value.pos.raw, "-|>")

  edge(g.value.pos.raw, c.value.pos.raw, "-|>", bend: 40deg)
  edge(g.value.pos.raw, c2.value.pos.raw, "-|>", bend: 40deg)

  edge(d.value.pos.raw, h.value.pos.raw, "-|>")
  edge(e.value.pos.raw, h.value.pos.raw, "-|>")
  // edge(f.value.pos.raw, h.value.pos.raw, "-|>")
  // edge(h.value.pos.raw, f.value.pos.raw, "-|>", bend: 40deg)
})
#align(center, dia)

This graph is bipartite in $F$ and $T$, i.e. fragment nodes only point to type nodes,
and type nodes only point to fragments.
The TF-Graph is useful becuase it makes *reachability analysis* easy.
From the types' perspective, the arrows show which fragments provide and require them.

If there doesn't exist a path following the in-edges backwards to a source node,
then that type or fragment isn't constructable within the given catalog.

TF-Graph is defined as follows:

+ Add a node for every fragment in the catalog
+ Add a node for every distinct type referenced across all fragments
+ For each fragment $F$, for each type $T$ in the set of it's arguments' types: add an edge $T -> F$.
+ For each fragment $F$ add an edge pointing to it's return type $T$: $F -> T$.

// And here it is in pseudocode:

// ```go
// catalog := Catalog()
// g := EmptyTFGraph()
// for _, frag := range catalog {
//   g.add_node(frag)
//   for _, typ := range frag.argtypes {
//     g.add_node(typ)
//     g.add_directed_edge(typ, frag)
//   }
//   g.add_node(frag.rtype)
//   g.add_directed_edge(frag, frag.rtype)
// }
// ```


== Dataflow

A dataflow is a program where a function can be executed any time all it's inputs are ready.
If multiple functions are ready then the order of execution is ambiguous and the scheduler makes a decision.
Sometimes the order of exection is unambiguous (see @dataflow-1 left), but not usually.

The diagram below shows two dataflows. The one on the left corresponds to just a single order of execution,
while the flow on the right can be executed in $3! = 6$ different ways. This introduces a bias
in the dataflow distribution when we sample evenly across programs.

#let d1 = diagram(node-stroke: 0.1em, node-shape: circle, {
  let a = (1, 2 * 0 / 3)
  let b = (1, 2 * 1 / 3)
  let c = (1, 2 * 2 / 3)
  let d = (1, 2 * 3 / 3)
  node(a, "a")
  node(b, "b")
  node(c, "c")
  node(d, "d")
  edge(a, b, "-|>")
  edge(b, c, "-|>")
  edge(c, d, "-|>")
})
#let d2 = diagram(node-stroke: 0.1em, node-shape: circle, {
  let a = (1, 0)
  let b = (0, 1)
  let c = (1, 1)
  let d = (2, 1)
  node(a, "a")
  node(b, "b")
  node(c, "c")
  node(d, "d")
  edge(a, b, "-|>")
  edge(a, c, "-|>")
  edge(a, d, "-|>")
})
#figure(
  align(center)[ #stack(dir: ltr, d1, d2, spacing: 3em) ],
  caption: [
    There is only one way to linearize the dataflow on the left: the program `d(c(b(a())))`, but there are $3! = 6$ ways to linearize
    that on the right: `a(); b(); c(); d();` as well as `a(); d(); c(); b()`, etc...
  ],
)<dataflow-1>

// #text(fill: gray, size: 8pt)[
//   There is only one way to linearize the dataflow on the left: the program `d(c(b(a())))`, but there are $3! = 6$ ways to linearize
//   that on the right: `a(); b(); c(); d();` as well as `a(); d(); c(); b()`, etc...
// ]

// When naively sampling programs in a manner that seems unbiased it is easy to end up with a highly biased distribution of dataflows.
// For example: there are 4! programs dataflow-equivalent to `A(); B(); C(); D();`, but only one equivalent to `F(G(H(I())))`.
// Another way of stating this is "the type signatures of fragSo the type signatures of our functions affect the nu

order of execution of statements is only implicitly defined via the connections between symbols and may not be unique.
We go from `Program = List<Statement>` to `Dataflow = Set<Statement>`.
When our catalog consists only of pure functions then the dataflow determines everything about the program.

_But do we sample evenly from the space of programs? Doesn't the catalog affect this?_

Our fragment-first sampling technique (@forward)

=== Counting Dataflows

In order to remove the bias in the dataflow distribuion we either need to

- sample evenly across programs, but then compensate for the bias by rejecting some flows probabilistically
- sample dataflows directly.

#horizontalrule

Every statement has immediate and transitive dependencies.
If two statements perform the same operation and have the exact same transitive dependencies then they are *syntactically equivalent*.
This can be useful if we want to repeat the same calculation at different times,
or if the system has hidden state#footnote[it reads or writes some global mutable unreferenced by the program] that affects the results.

Let's begin by enumerating *minimal dataflows*, i.e. flows where no two terms are equivalent.
If we restrict the set of dataflows in this way then the space naturally factorizes in to layers,
where dependencies only flow in one direction.

#smallcaps[An example.]
Let the catalog contain three items `{0, 1, +}`.
Now let's build up layers containing all valid syntactic expressions,
where expressions in $l_i$ contain at least one reference to a term in $l_(i-1)$

// We put {`0, 1`} on layer $l_0$ for total size $c_0 = 2$.
#block(stroke: 0.05em, inset: 1em, width: 100%)[
  `l0 = 0, 1` \
  `l1 = 0+0, 0+1, 1+0, 1+1` \
  `      a    b    c    d              `
  #text(fill: gray)[definitions] \
  `l2 = a+a, a+b, b+a, b+b, a+c, + ... `
  #text(fill: gray)[all combinations of two $l_1$ arguments] \
  `     a+0, 0+a, a+1, + ...           `
  #text(fill: gray)[all combinations of one $l_1$ and one $l_0$.]
]

// v(1em)\ #align(center, body:body)\ #v(1em)

#question[_How many terms are there at each level?_]

Let's define $c_i = |l_i|$.
There are two terms in $l_0$ so $c_1 = 2 times 2 = 4$,
and then $c_2 = c_1^2 + 2c_1c_0 = 32$.

The the number of available terms at level $i$ is

#set math.equation(numbering: "[1]")

$
  c_(i+1) = underbrace(#h(1em) c_i^2 #h(1em), "two from" l_i) + underbrace(2c_i C_(i-1), "one from" l_i \ "and one from" l_(j<i))
$ <eq-1>

where $C_k = sum_(j=0)^k c_j$ is the total number of terms at level $k$ and below.

#question[_How many minimal dataflows have_ $n$ _terms?_]

A particular flow instance may not use every available term at each level,
but instead we pick a subset of terms $m_i subset.eq l_i$,
and our specific choices at level $i$ determine the set of possible choices at level $i+1$.

If we define $h_i = |m_i|$ as the number of terms selected from level $i$
then the number of possible terms to choose from at $i+1$, conditional on $h_(0..i) = [h_0, h_1, ..., h_i]$, is
$ h^*_(i+1) = h_i^2 + 2h_i H_(i-1) $
where $H_i = sum_(j=0)^i h_j$.
This means the number of possible flows $phi_n (bold(h))$ given our choices $bold(h) = [h_0, h_1, ..., h_n]$ is just the product of these choices $vec(h_i^*, h_i)$ at every level#footnote[we know that $h_i = 0$ for $i>n$.]
i.e.

$
  phi_n (bold(h)) = product_(i=0)^n vec(h_i^*, h_i)
$

where our choices for $bold(h)$ are bounded from above $h_i <= c_i$.
And the total number of flows $Phi_n$ is

$
  Phi_n = sum_bold(h) phi_n (bold(h)) bracket.double.l h_0 + h_1 + ... = n bracket.double.r
$

i.e. the sum across all $bold(h)$ that sum to $n$ total terms.

#question[What about catalogs beyond `{0, 1, +}` ?]

The main ideas still hold.
We can still segment the set of all possible terms into levels of finite size.
Now our counts $h_(i tau)$ are conditioned on level $i$ and type $tau in {T_1, T_2, ..., T_c}$.
And the recurrence changes

$ h_((i+1) tau) = $


=== Sampling Dataflows

Dataflow-first program generation splits the problem into two phases

+ sample a dataflow
+ sample a valid evaluation order of the terms.

When fuzzing pure, functional systems we can spend our effort sampling different dataflows,
while for impure-dominant systems we can more easily explore subtle differences in evaluation order and timing.
We sample a flow by,

+ choose $n$
+ choose $bold(h)$ consistent with $n$
+


We believe that exploring different DAGs is important for exploring software for a few reasons.
First, the DAG describes all the dataflow of the program.
So, to the extent that the system is pure, it's entire behaviour is determined by the DAG,
and all programs with identical DAGs are equivalent.
By implication, to the extent that the system is impure, it's entire behaviour is determined increasingly by the timing of fragment executions,
and not based on program dataflow.
We would like to control the dataflow as well as the order/timing of execution explicitly!
And we don't want either to be determined by the type signature of what happens to be in the catalog
(other than supplying basic constraints on the space).
_The catalog constrains the space of programs,
but it shouldn't bias the distribuion over that space!_

@forward[_Foward sampling_] is unable to sample evenly across DAGs,
because selecting fragments _first_ biases the distribution.
E.g. if we sample a fragment with no arguments then it will be inserted every time.
This adds new root nodes to the DAG and means the growth rate of root nodes is beyond our control and depends on the catalog composition.
A potential solution to this problem is to adjust the probability of sampling fragments depending upon their type.
E.g. we may want to make it less likely that fragments with simple types are inserted into the program,
or even place hard constraints on the number of fragments of a given type that are allowed.

#horizontalrule

One factorization of programs is

1. pick a DAG where nodes are equivalent up to type
2. pick a fragment for each node / type
3. pick a linearization for the DAG.

Another is

#horizontalrule

How can we count DAGs? Look at R. Stanley's _Enumerative Combinatorics_ for the answer.


// produces a value, and two values which are
// Let's begin by assuming that every statement in the dataflow produces a unique value.

