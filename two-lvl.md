
Hierarchical models allow us to chnage the way that programs are generated.
Instead of having a single sampleProgram() we'll have many functions with type `Library, Params, Program -> Program`.
The Library allows us to alter the set of Funs and Values that can be sampled.
Program is what we've built thus far (a prefix of lines 0..n). 

Examples include
- swarm : randomly remove a Fun from the Library.
- permute : sample each Fun one time in a particular order.
- threads : keep separate threads in the DAG, but interleave the lines e.g. `a = f(); b = f(); a2 = f(a); b2 = f(b); a3 = f(a, a2?); etc...`
- repeat : just pick a single op and repeat it `a1 = f(a0); a2 = f(a1); a3 = f(a2) ...`
- mutate : re-assign an existing sym `a = f(a)`
- repeat-mutate : sample a small program using inputs from above `loop 7 {a = f(a)}`
- repeat-mutate-block : sample a small program using inputs from above `loop N {b = f(a); c = g(b); a = h(c)}`
- put weights on the available params
- are all of these ops just different ways of putting weights on Fun, Value, and Sym for a duration?
- are these ideas subsumed by having fun with Sets of objects in teh same way as we played with Sets of Links/Nodes for the partition tests?


- `Library, Program, Params -> Program`
- `Library -> Library, LibraryWeights`


We can write a generator using these high level components, but can we generate them directly?
I think we can! Just make a new Lib consisting of these objects... And run sampleProgram()?
Be sure that you get one 

Example Problems
- password



---

I think the only way that I can think of at the moment to achieve what I believe Dave wants, and what will give us interesting programs across an arbitrary set of margins is the following:

- We enumerate or sample a sequence of meta-programs
- For each metaprogram we sample a program []Stmt in a particular way, e.g. using a subset of fragments, with weird wiring, or using a certain pattern. Each sub-program is concatenated to form the full program.

This makes a lot possible, but it makes everything more complicated.
First, we have to think of and implement a set of meta-fragments or patterns.
Then we have to implement a mata-program sampler which enumerates meta programs and samples them.

How to do credit assignment?

What kind of grammar for metaprograms?
Are all pieces independent?
Is there a state machine that tells me which pieces are allowed?
Or is there a stack-like context that tells me which pieces are allowed?
Can the pieces set/unset arbitrary state?
Does each piece return a Program, or are there other types of pieces that *only* set context?
E.g. imagine taking the product between Pieces that return state, pieces that control the set of available fragments, and pieces that control wiring.
Are these all independent? 

---

Explore programs not by random sampling at the top level, but by pure enumeration.
If you've got a space where some values are very unlikely (like a sequence of N ops that don't use some fragment) and we want to increase the likelihood of this kind of program we can create a hierarchical model that first picks a class of structure, many of which may be unlikely under the default model, and then samples from a new dist with support zero for programs outside of that structure.

Now we have to sample a program at both levels of the hierarchy.

What about an algorithm that did brute force enumeration of every possible combination of the following.

- set of functions used
- sequence of functions used
- possible wirings

hierarchical model that picks eight digits.
- independently randomly samples each digit.
- randomly samples first, then 2nd conditioned on first, then 3rd conditioned on 2nd
- just enumerates the numbers starting from 00000000
- enumerates the numbers in a random order?
- enumerates 00000000, 11111111, 22222222, etc, except instead of incrementing by 1 on each digit
    we increment the ith digit by d_i (mod 10) where d_i is the ith prime number (excluding 2 and five).

In general for a sequence of spaces if we want to cover all small pairs of margins quickly random sampling
will be good at covering margins we *aren't aware of* but is bad at getting to 100% coverage!
If we'd rather have 100% coverage of each of a small list of margins we should enumerate and iterate,
unless we have so many samples that coverage will be oversaturated!

For hierarchical models this makes a lot of sense at the higher levels.







