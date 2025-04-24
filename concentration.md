*Concentration of Measure* is the generalization of the key problem facing fuzzing.
Roughly speaking, our programs stop finding bugs and making them bigger doesn't help.
More precisely, it is difficult to obtain broad coverage of the program state space
from a set of randomly generated sample programs.

We can't just supply independent, uniform random values to the functions, we can't just create programs
that sample functions from an API in this way either.

[Swarm testing][swarm12] is a fuzzing approach that deals with this problem by restricing the API randomly per input program.
But the problem also arises at the level of function inputs, and in 

Key researchers from the mathematical side are [ledoux] and [talgrand].
Key researchers from the fuzzing / CS side are [groce] and [regher].

[swarm12]: https://agroce.github.io/issta12.pdf
[ledoux]: https://scholar.google.com/citations?user=8IX7iCIAAAAJ&hl=en&oi=ao
[talgrand]: https://scholar.google.com/citations?user=mVbEaoIAAAAJ&hl=en&oi=sra
[groce]: https://scholar.google.com/citations?user=ewrrvq8AAAAJ&hl=en&scioq=swarm+testing&oi=sra
[regher]: https://scholar.google.com/citations?user=yCv02hMAAAAJ&hl=en

[measureNotes08]: https://web.math.princeton.edu/~naor/homepage%20files/Concentration%20of%20Measure.pdf

