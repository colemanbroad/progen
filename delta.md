
Synonyms in the literature:

- reduction
- shrinking
- minimization
- debloating

# Sequences

Zeller[ZH02] defines two algorithms `ddmin` which minimizes a failing test case and `dd`
which minimizes the failing test case while simultaneously maximizing the fraction of
the input that comes from an additional, *passing* test case. In this overview we will
focus only on the basic `ddmin`. The original `ddmin` is demonstrated on HTML, Mozilla
User Click data, and C source code. When targeting C they run ddmin on the sequence of
characters in the file, ignoring any information from a potential C tokenizer/parser.

Specifically [ZH02] propose a coarse-to-fine approach that ignores any inherent structrue of
the input other than the raw sequence. The core algorithm partitions the input into `n=2`
subsequences by cutting it in the middle. Each partition is tested on it's own, then if
the condition is not reproduced the *complement* of each partition is tested. If at any
point we manage to repro the failure then we discard the unsued sequence and reset `n=2`,
otherwise at the of the `2n` test cases we double `n *= 2` and repeat.

The idea behind `ddmin` is so broad that it cries out for generalizations beyond strings
of characters to data with more structure. The best summary of the history of test-case
reduction and generalizations of delta debuggin can be found on John Regehr's blog posts
[0],[c-reduce]. Regehr says this about the original Zeller paper,

  > The enduring value of the paper is to popularize and assign a name
  to the idea of using a search algorithm to improve the quality of
  failure-inducing test cases.

Regehr's group tested four new techniques where each is a different generalization of the
original ddmin taking advantage of some aspect of the input structure. Interestingly they
find that

  > local minima can often by escaped by running the implementations one after the other.

This work culminated in [c-reduce] which differs in interesting ways from other approaches
to minimization. Perhaps most importantly, is it tries hard to escape local minima!

  > The C-Reduce core does not insist that transformations make the test case smaller, and
    in fact quite a few of its passes potentially increase the size of the test case, with
    the goal of eliminating sources of coupling within the test case, unblocking progress
    in other passes.

This "eliminating the source of coupling" idea is like our idea that upstream nodes
in the program DAG that prevent the DAG from being a Tree can be copied/duplicated
if it allows us to prune large branches of the resulting tree.

  > The sequence of transformation passes is carefully orchestrated such that passes that
    are likely to give the biggest wins -- such as those that remove entire functions --
    run first;

  > only a small proportion of the transformation passes is intended to be
    semantics-preserving [...] we only want to preserve enough semantics that we can
    probabilistically avoid breaking whatever property makes a test case interesting

However [c-reduce] still runs, fundamentally, on C/C++ source code, while *we are running
directly on program IR*.

# Trees

How can `ddmin` best be applied to data with structure beyond that of a string of chars?

There is has been continuous research on this topic since [ZH02] with perhaps the most
successful variant being Hierarchical Delta Debugging [HDD], wihch is designed
for tree-structured inputs and requires knowledge of a (context-free) grammar describing
the input. [HDD] takes a coarse-to-fine ablation approach where first the top level
nodes in the parse tree are removed before moving down the tree. It can be applied to
program source by applying the parser and running HDD on the resulting AST. However this
approach breaks down when the input format allows for defining *references* that allow for
one location to point to a different span of the input. Obviously, this is pervasive in
programming languages and in practice means that a naive application of HDD will lead to a
large number of dangling references, potentially masking the true source of failure.

Additional techniques discussed in [0] include iterative delta debugging[idd10] and
Lithium[lith09]. Lithium in particular does two things I really like. First, instead
of referring to test cases as "failing" or "passing" they refer to them as being
"interesting" or "uninteresting", which is a generalization that Antithesis has also
made. Second they have a nice definition of the "monotonicity" condition under which their
algorithm will be able to find SOME solution (if not the global optimum).

  > Monotonicity: no subsequence of an uninteresting file is ever interesting.

In the domain of digital images this is equivalent to increasing the threshold on an image
and then segmenting it. There may be many connected components, and those components may
divide and shrink as we increase the threshold. "Monotonic" means our interestingness
test is something like "does this sequence of statments S1 contain the subsequence S2" or
something like that. You can imagine asking the same question about sets of pixels.

# DAGs 

In the comments around [0] a user proposes a more flexible version of [HDD] that works
with DAGs intead of trees. The idea is that you have a linear thing (a source file) and
you have _guesses_ as to which chunks of the file are related in some way e.g. you make
guesses as to potential grammars. If we _knew_ the grammar of the file (which delimiters
are used to increase the scope) we would be able to form a parse tree (assuming the
grammar doesn't allow defining references). But if we don't know the delimiters we can
still guess and maybe our guesses overlap each other. Thus we've got potentially
overlapping spans of source and these naturally form a partial order `A subset B`.

Spans arise naturally in [HDD] as a tree of spans are naturally formed by a Context
Free grammar describing a given source file. They ALSO arise in the context of
distributed systems from "event traces" 
We've seen partial orders formed by spans before! They arise naturally when
thinking about concurrent process and spans of time. We want to know "does span A
happen before or after span B?" but we only have a definitive answer if the spans
don't overlap. 
                                      
                                         ┌────────────────( B )───────────────┐  
                                         │                                    │  
                           ┌─────────────────────{ A }────────────────────┐      
                           │                                              │      
    And here we test the " { " character (a far *cooler* alternative to " } " ). 
                         │   │                  │      │                │   │    
                         │   │                  │      │                │   │    
                         └─C─┘                  └*  E *┘                └─D─┘    
                                                                                 
This produces the following DAG (we only consider an irreflexive version of each relation)
according to the `subset-of` relation:

    A
    B
    C
    D -> B
    E -> A,B

And the following Graph according to the `intersects` relation:

    A — B,C,D,E
    B — D,E

    and by symmetry:

    C — A
    D — A,B
    E — A,B
    
We could include the DAG produced by the `precedes` relation as well:

    A
    B
    C -> B,D,E
    D
    E -> D

We can infer the graph of `intersects` (inter) from the `precedes` (<) DAG by noting that

    x<y => ! x inter y
    ! x<y AND ! y<x <=> x inter y   . 

Together [WIP]

See how these delimeters form spans which form partial order? Only the D span is
completely contained within B; all other pairs have no direct `subset` relation.

But there IS a difference between `A intersects B` and `A !intersects B`. 
This difference may actually be very important for a parser, because we expect
to be able to excise a span without disrupting any of it's non-intersecting spans.
So actually we have both the assymmetric `A <= B` and the symmetric `A ⋂ B` (intersects).
And it's possible that A⊂B, B⊂A, A⋂B, B⋂A can all be false! 
There is one implication:

    A !inter B => A !<= B 

Thus this should tell us how to go about searching for smaller input files.

1. Propose pairs of delimiters.
2. Determine all possible delimiter spans in the file.
3. Determine the `intersects` and `subset` relation for these spans.
4. Optimize using the resulting DAG.

The DAG is formed only by the `subset` relation.
The addition of the `intersects` relation complicates the situation, because it
is a symmetric relation.

### Deep in the weeds on Relations

Common Properties of Relations:


Property      | Definition
--------      | ----------
symmetric     | aRb => bRa
asymmetric    | aRb => !bRa
antisymmetric | aRb AND bRa => a = b (this is a weird one. it references `equals` explicitly?!)
transitive    | aRb AND bRc => aRc
intransitive  | aRb AND bRc => !aRc
reflexive     | aRa forall a[^fn2]
irreflexive   | !aRa forall a

[^fn2]: This is also a werid one. It's not an implication. It's telling you literally what
pairs are in the relation.

- Relations are equivalent to DiGraphs (maybe with 2-cycles).
- Symmetric Relations (e.g. intersects) are equivalent to Graphs (maybe with 1-cycles!?)
- Asymmetric Relations (e.g. child of) are equivalent to Digraphs WITHOUT 2-cycles.
- Weak Partial Orders (e.g. <=) are reflexive, asymmetric and transitive, thus they contain all 1-cycles.
- Strict Partial Orders (e.g. <, descendant of) are a special case of Asymmetric Relation without ANY cycles.
- PO's can be represented as DAGs.
- A Total Order is equivalent to a Chain.
- There are many (most!) relations which are neither symmetric NOR asymmetric! E.g. "A likes B".
- We can define multiple relations on a set of elements, resulting in multiple coexisting (Di)Graphs.

Relation        | Properties
--------        | ----------
equal to        | symmetric, reflexive, transitive
not equal to    | asymmetric, irreflexive, ~~transitive~~
child of        | asymmetric, irreflexive, intransitive
descendant of   | 
related to      | 
intersects      | 
subset of       | 
likes           | 
correlated with | ??


# Program Source

Programs have tons of special structure that we can take advantage of, or ignore and risk
wasting lots of effort. The current SotA for program minimization is a tool [Chisel][3], [3a]
from UPenn. They also refer to this task as program "debloating". The current benchmark
reducer is John Regehr's [c-reduce], which is designed for C/C++ code but apparently also
works well on other (C-like?) languages which can take advantage of the initial reduction
phases that don't rely on Clang's C/C++-specific analysis passes.

# Python function arguments (Property based testing)

The best example is [Hypothesis][hyp13], which is probably what used to be called [pydelta]? How
does Hypothesis do reduction? They have a few extra fancy tricks up their sleeves like in [invariant
parameters] produced during the `explain` phase? 

[invariant parameters]: https://hypothesis.readthedocs.io/en/latest/reference/api.html#controlling-what-runs

  > After shrinking to a minimal failing example, Hypothesis will try to find parts of
    the example – e.g. separate args to @given() – which can vary freely without changing
    the result of that minimal failing example. If the automated experiments run without
    finding a passing variation, we leave a comment in the final report:

    test_x_divided_by_y(
        x=0,  # or any other generated value
        y=0,
    )

, and allowing failing cases to be stored in a database of your choosing. This is
an interesting choice and I don't really get it...

  > Hypothesis takes a philosophical stance that property-based testing libraries, not
    the user, should be responsible for selecting the distribution. As an intentional
    design choice, Hypothesis therefore lets you control the domain of inputs to your
    test, but not the distribution.


# What should WE do?

What we need to apply this to our programs is...
Know the structure of the program. Which nodes can be cut/replaced?
Which components are totally independent at the top level? (There may only be one
giant Connected Component.) When we make a cut we have to replace the argsyms
with something new. Or we have to replace the (Fun Args...) with a new Fun
and new Args (potentially none) of the right type. We're trying to cut
large swaths of the search space. Ideally we'd cut the program in half
each time!

*There may be a fundamental tension between creating programs that are DEEP
(the goal of our depth distribution experiments) and creating programs that
are easy to minimize!*

First, we can cut out all the program that isn't upstream of the condition/
failure. Then, we can go about pruning the DAG. Wherever we use a Fun that
takes varargs we can try cutting a subset of the inputs. Wherever we have
a Fun that takes a typed arg, we can replace it with a minimal construct
of the same type. Some types can be created directly (we have an int value "one"
always available), but others may require small (minimal) programs to build.

Our condition/failure may not live in the program directly, but may come
from some external oracle with unknown dependencies on the different lines
of the Program. This means instead of immediately being able to prune all
the top level code that isn't upstream of the condition we must prune it
by running standard Delta Debugging on the set of top-level connected
components.

Delta Debugging is for _set structured_ input.
Hierarchical Delta Debugging is for _tree structured_ input. 
We need something for _DAG structured_ input.

1. Delta Debugging for distinct connected components.
2. 

A `Span` has (at least) the following relations: `precedes`, `intersects`, `subset of`.
We only need ~one of~ `precedes` ~or `intersects`~, because we can derive `intersects`
from precedes (but not vice versa!). But we may prefer `intersects` because it could
save a little work.

    A precedes B => A !intersects B
    A !precedes B and B !precedes A <=> A intersects B

# Distributed System Traces (observed, not controlled!)

Are there examples of DD for finding linearizeablility issues with CSP?
This would look like the following. I have a program which creates N state
machines. It sends and receives messages from them adding in delays and message
re-ordering. All the while it's checking the program invariants. 
Let's NOT take this path, as it's too far removed from the ideas we've already
implemented: the state machines (containers) under our control don't send messages
only to us, but they speak to each other via a NetworkManager intermediary. We
control each machine via it's interface as well as the NetworkManager. Our job
is to find a Program which generates inconsistent states by looking for logic
that _assumes_ that the order of messages sent is the order of messages received,
and that the order of messages received is the order in which they are processed.
And that the order in which they are processed is the order in which they will
respond, and that the order in which they respond is the order in which you observe
the responses.
These assumptions can be implicitly encoded in a variety of ways.

[WIP]

# Fuzzing

## C-Smith

With [c-smith] Regehr generates C programs and [guided tree search][gts21] tries to
balance the "embedded decision tree" corresponding to that program.

  > For certain well-structured special cases, such as generating strings from a regular
  grammar, algorithms exist to generate strings with uniform probability, but it is easy
  to see that in the general case, no such algorithm can exist. The proof of this is
  handwavy, but observe that an arbitrary-sized subtree can lurk past every unexplored
  edge in the decision tree. Without prior knowledge of the tree shape, there is simply no
  way to know which unexplored branches lead to one or two leaves and which lead to a vast
  number of leaves.

## Z3 model code 

[TypeFuzz] is interesting work on *program mutations* for finding bugs from 2021.

  > Generative Type-Aware Mutation is a hybrid of mutation-based and grammar-based
  fuzzing and features an infinite mutation space overcoming a major limitation of
  OpFuzz, the state-of-the-art fuzzer for SMT solvers. We have realized Generative
  Type-Aware Mutation in a practical SMT solver bug hunting tool, TypeFuzz. During our
  testing period with TypeFuzz, we reported over 237 bugs in the state-of-the-art SMT
  solvers Z3 and CVC4.

  > The reports shown are reduced bug triggers after bug reduction with pydelta and
  C-Reduce.

References [pydelta],[C-Reduce]. 

## SQL queries

[SQL98] is some of the first work using test-case reduction according to [c-reduce],
even if this wasn't it's primary purpose. This [fig](pics/Screenshot 2025-04-11 at
7.16.35 PM.png) is a simple diagram of the coverage problem.

There is ongoing work on fuzzing databases [204] [205], although I guess it mostly focuses
on query plan coverage and not on dist-sys correctness.

# Biblio 

[ZH02]: https://ieeexplore.ieee.org/stamp/stamp.jsp?tp=&arnumber=988498
Zeller and R. Hildebrandt. "Simplifying and isolating failure-inducing input."  

[HDD]: https://users.cs.northwestern.edu/~robby/courses/395-495-2009-fall/hdd.pdf
Ghassan Misherghi and Zhendong Su. "HDD: Hierarchical Delta Debugging"  

[3]: https://github.com/aspire-project/chisel
Chisel: Effective Program Debloating via Reinforcement Learning

[204]: https://www.semanticscholar.org/reader/ec682d9c7d68149dcd8932acd01a751f2f8b5611
Testing Database Engines via Query Plan Guidance  

[0]: https://blog.regehr.org/archives/527
[3a]: https://pardisp.github.io/_papers/chisel-poster.pdf  
[idd10]: https://people.kth.se/~artho/papers/artho-idd-10.pdf
[lith09]: https://www.squarefree.com/lithium/algorithm.html
[205]: https://scholar.google.com/scholar?as_ylo=2021&q=Massive+Stochastic+Testing+of+SQL&hl=en&as_sdt=0,9
[pydelta]: missing!?
[c-reduce]: https://blog.regehr.org/archives/1678
[SQL98]: https://www.semanticscholar.org/paper/Massive-Stochastic-Testing-of-SQL-Slutz/74b2c1bce3963fbb1300dc0995b9e275f3393cb9?p2df
[hyp13]: https://hypothesis.readthedocs.io/en/latest/
[WIP]: WorkInProgress
[TypeFuzz]: https://dl.acm.org/doi/pdf/10.1145/3485529
[gts21]: https://github.com/regehr/guided-tree-search
[c-smith]: https://github.com/csmith-project/csmith

# Wip

[PDD21]: https://xiongyingfei.github.io/papers/FSE21a.pdf

[RLDD18]: https://chisel.cis.upenn.edu/papers/ccs18.pdf

[200]: https://ieeexplore.ieee.org/abstract/document/5954405
Search-Based Software Testing: Past, Present and Future  

[201]: https://sci-hub.ru/10.1109/tse.2016.2521368
A Survey on Software Fault Localization  

[202]: https://ieeexplore.ieee.org/abstract/document/8863940
The Art, Science, and Engineering of Fuzzing: A Survey  

[203]: https://ieeexplore.ieee.org/abstract/document/10179324
Examining Zero-Shot Vulnerability Repair with Large Language Models  

[flakey21]: https://dl.acm.org/doi/fullHtml/10.1145/3476105#Bib0057
A survey of flakey tests
