
# Basic delta debugging

Zeller[1] defines two algorithms `ddmin` which minimizes a failing test case and `dd`
which minimizes the failing test case while simultaneously maximizing the fraction of
the input that comes from an additional, *passing* test case. In this overview we will
focus only on the basic `ddmin`. The original `ddmin` is demonstrated on HTML, Mozilla
User Click data, and C source code. When targeting C they run ddmin on the sequence of
characters in the file, ignoring any information from a potential C tokenizer/parser.

Specifically [1] propose a coarse-to-fine approach that ignores any inherent structrue of
the input other than the raw sequence. The core algorithm partitions the input into `n=2`
subsequences by cutting it in the middle. Each partition is tested on it's own, then if
the condition is not reproduced the *complement* of each partition is tested. If at any
point we manage to repro the failure then we discard the unsued sequence and reset `n=2`,
otherwise at the of the `2n` test cases we double `n *= 2` and repeat.

The idea behind `ddmin` is so broad that it cries out for generalizations beyond strings
of characters to data with more structure. The best summary of such extensions is Regehr's
blog post[0]. Regehr says this about the original Zeller paper[1].

  > The enduring value of the paper is to popularize and assign a name
  to the idea of using a search algorithm to improve the quality of
  failure-inducing test cases.

Regehr's group tested four new techniques where each is a different generalization of the
original ddmin taking advantage of some aspect of the input structure. Interestingly they
find that

  > local minima can often by escaped by running the implementations one after the other.
 
# ddmin for structured data

How can `ddmin` best be applied to data with structure beyond that of a string of chars?

There is has been continuous research on this topic since [1] with perhaps the most
successful variant being Hierarchical Delta Debugging [HDD][2] (HDD), wihch is designed
for tree-structured inputs and requires knowledge of a (context-free) grammar describing
the input. [HDD][2] takes a coarse-to-fine ablation approach where first the top level
nodes in the parse tree are removed before moving down the tree. It can be applied to
program source by applying the parser and running HDD on the resulting AST. However this
approach breaks down when the input format allows for defining *references* that allow for
one location to point to a different span of the input. Obviously, this is pervasive in
programming languages and in practice means that a naive application of HDD will lead to a
large number of dangling references, potentially masking the true source of failure.

Additional techniques discussed in [0] include iterative delta debugging[4] and
Lithium[5]. Lithium in particular does two things I really like. First, instead
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

In the comments around [0] a user proposes a more flexible version of [2] that works
with DAGs intead of trees. The idea is that you have a linear thing (a source file) and
you have _guesses_ as to which chunks of the file are related in some way e.g. you make
guesses as to potential grammars. If we _knew_ the grammar of the file (which delimiters
are used to increase the scope) we would be able to form a parse tree (assuming the
grammar doesn't allow defining references). But if we don't know the delimiters we can
still guess and maybe our guesses overlap each other. Thus we've got potentially overlapping spans of source and these naturally form a partial order `A subset B`.

## An aside on spans and relations

Spans arise naturally in [2] as a tree of spans are naturally formed by a Context
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
    And here we test the " { " character (a far superior alternative to " } " ). 
                         │   │                                          │   │    
                         │   │                                          │   │    
                         └─C─┘                                          └─D─┘    

This produces the following DAG according to the `strict-subset-of` relation:

    D -> B
    A
    C

And the following Graph according to the `intersects` relation:

    A - B - D
    C - A

We could include the DAG produced by the `precedes` relation as well:

    C -> B,D

We can infer the graph of `intersects` (inter) from the `precedes` (<) DAG by noting that

    x<y => ! x inter y
    ! x<y AND ! y<x <=> x inter y   . 

Together


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

Common Properties of Relations:
- symmetric     ; aRb => bRa
- asymmetric    ; aRb => !bRa
- antisymmetric ; aRb AND bRa => a = b (this is a weird one. it references `equals` explicitly?!)
- transitive    ; aRb AND bRc => aRc
- intransitive  ; aRb AND bRc => !aRc
- reflexive     ; aRa forall a (this is also a werid one. it's not an implication. it's telling you literally what pairs are in the relation).
- irreflexive   ; !aRa forall a

Relations are equivalent to DiGraphs (maybe with 2-cycles).
Symmetric Relations (intersects) are equivalent to Graphs (maybe with 1-cycles!?)
Asymmetric Relations (child of) are equivalent to Digraphs WITHOUT 2-cycles.
Weak Partial Orders (<=) are reflexive, asymmetric and transitive, thus they contain all 1-cycles.
Strict Partial Orders (<, descendant of) are a special case of Asymmetric Relation without ANY cycles.
PO's can be represented as DAGs.
A Total Order is equivalent to a Chain.
There are many (most!) relations which are neither symmetric NOR asymmetric! E.g. "A likes B".
We can define multiple relations on a set of elements, resulting in multiple
coexisting (Di)Graphs.

The `intersects` relation is symmetric (representable by a graph).
The `subset` relation is partial order.

Relation        ; Properties
equal to        ; symmetric, reflexive, transitive
not equal to    ; asymmetric, irreflexive, ~transitive~
child of        ; asymmetric, irreflexive, intransitive
descendant of   ; 
related to      ; 
intersects      ; 
subset of       ; 
likes           ; 
correlated with ; ??

A `Span` has (at least) the following relations: `precedes`, `intersects`, `subset of`.
We only need ~one of~ `precedes` ~or `intersects`~, because we can derive `intersects`
from precedes (but not vice versa!). But we may prefer `intersects` because it could
save a little work.

    A precedes B => A !intersects B
    A !precedes B and B !precedes A <=> A intersects B




# linearizeability + DD

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


# What should WE do?

What we need to apply this to our programs is...
Know the structure of the program. Which nodes can be cut/replaced?
Which components are totally independent at the top level? (There may only be one
giant Connected Component.) When we make a cut we have to replace the argsyms
with something new. Or we have to replace the (Fun Args...) with a new Fun
and new Args (potentially none) of the right type. We're trying to cut
large swaths of the search space. Ideally we'd cut the program in half
each time!

    There may be a fundamental tension between creating programs that are DEEP
    (the goal of our depth distribution experiments) and creating programs that
    are easy to minimize!

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


# Biblio 

[1] Zeller and R. Hildebrandt. "Simplifying and isolating failure-inducing input."
[2] Ghassan Misherghi and Zhendong Su. "HDD: Hierarchical Delta Debugging"
[3] Chisel: Effective Program Debloating via Reinforcement Learning
[200] Search-Based Software Testing: Past, Present and Future
[201] A Survey on Software Fault Localization
[202] The Art, Science, and Engineering of Fuzzing: A Survey
[203] Examining Zero-Shot Vulnerability Repair with Large Language Models

[0]: https://blog.regehr.org/archives/527
[1]: https://ieeexplore.ieee.org/stamp/stamp.jsp?tp=&arnumber=988498
[2]: https://users.cs.northwestern.edu/~robby/courses/395-495-2009-fall/hdd.pdf
[3]: https://github.com/aspire-project/chisel
[4]: https://people.kth.se/~artho/papers/idd.pdf
[5]: https://www.squarefree.com/lithium/algorithm.html
[200]: https://ieeexplore.ieee.org/abstract/document/5954405
[201]: https://sci-hub.ru/10.1109/tse.2016.2521368
[202]: https://ieeexplore.ieee.org/abstract/document/8863940
[203]: https://ieeexplore.ieee.org/abstract/document/10179324
