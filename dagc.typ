
#let horizontalrule = line(start: (25%, 0%), end: (75%, 0%))

= Program generation methods

- Forward
- Forward + tail weight
- DAG sampling
- Backward
- Tree-based (Controling margins)
- Two-level

_Each approach leads to a different distribution of programs._

#set heading(numbering: "1")

= Forward sampling <forward>

Program generation is currently performed one line at a time ("forward") by sampling

+ a random fragment from the catalog
+ a random tuple of existing, appropriately typed syms as arguments.

= Forward + tail weight

Fragments are sampled evenly from the catalog,
but wiring has a preference for syms defined more recently.

= DAG sampling

A TacticsLang program's dataflow can be represented as a DAG
where each node corresponds to a statement defining a new sym.
This representation is clean in TacticsLang as the language doesn't allow defining multiple syms per statement.

#import "@preview/fletcher:0.5.8" as fletcher: diagram, edge, node
#let nodes = (
  `v1 = f()`,
  `v2 = g()`,
  `v3 = h(v1, v2)`,
  `v4 = i(v3, v2)`,
)
#let mydia = diagram(node-stroke: 0.1em, node-corner-radius: 0.3em, {
  let a = (1 / 2, 0)
  let b = (0, 1 / 2)
  let c = (1, 1)
  let d = (1 / 2, 1.8)
  node(a, nodes.at(0))
  node(b, nodes.at(1))
  node(c, nodes.at(2))
  node(d, nodes.at(3))
  edge(a, c, "-|>")
  edge(b, c, "-|>")
  edge(c, d, "-|>")
  edge(b, d, "-|>")
})

#align(center + horizon)[
  #stack(
    dir: ltr,
    ```python
    v1 = f()
    v2 = g()
    v3 = h(v1, v2)
    v4 = i(v3, v2)
    ```,
    text(mydia, 7pt),
    spacing: 3em,
  )]
#text(fill: gray, size: 8pt)[
  The program on the left corresponds to the DAG on the right _evaluated in the order_: `v1 v2 v3 v4`,
  but the order `v2 v1 v3 v4` would also be valid.
]

The goal of DAG sampling is split program generation into two phases:

+ sample a random DAG
+ randomly select a (valid) evaluation order of the nodes.

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


= E-Graph

The E-graph is uniquely defined by the Catalog,
and has a node for every element of the Catalog and a box drawn around nodes with the same return type
and a numbered, directed edge out from each node pointing to a box of appropriate type corresponding to the arguments to that function.

= TF-Graphs

The TF-Graph is also uniquely defined by the Catalog.
The initial catalog is a `map[name]Fragment`,
where each fragment is a type with a single arrow
#footnote[Let's ignore higher-kinded types (HKTs) for now.]
mapping a tuple of input types to a single output type.
We can represent the fragments and types in a single structure, the TF-Graph, defined as follows:

+ Add a node for every fragment in the catalog
+ Add a node for every distinct type referenced across all fragments
+ For each fragment $F$, for each type $T$ in the set of it's arguments' types: add an edge $T -> F$.
+ For each fragment $F$ add an edge pointing to it's return type $T$: $F -> T$.

And here it is in pseudocode:

```go
catalog := Catalog()
g := EmptyTFGraph()
for _, frag := range catalog {
  g.add_node(frag)
  for _, typ := range frag.argtypes {
    g.add_node(typ)
    g.add_directed_edge(typ, frag)
  }
  g.add_node(frag.rtype)
  g.add_directed_edge(frag, frag.rtype)
}
```

This graph is bipartite in $F$ and $T$, i.e. fragment nodes only point to type nodes,
and type nodes only point to fragments.

= Dataflow

Let's try to make dataflows and see how easy it is.
We're going to just make an index type $->$ catalog.

