
Testing distributed systems (DS) is a difficult but valuable niche within software testing with traditionally only very little overlap with the larger practice.
In DS the environment may produce faults that are not under your control, causing communication between machines to fail.

Communication by message passing, but it's very difficult to send a message *exactly one time*[^1]
Thus the importance of *idempotent* designs, retries, acks, and all kinds of techniques and tricks that appear much more often in DS.

[^1]: It's easy to send a message at least once or at most once, but difficult to send it *exactly* once.

We cannot stop to observe the state of different machines. 
The state of the system is the combined state of each machine + the state of the environment.
Common failure modes are documented in Liskov's [Generalized Isolation Level Definitions][liskov00].
[Jepsen] is a DS testing LLC that applies fuzzing-esqe techniques to distributed systems to check for these failure modes. 

*How does Jepsen find bugs?*

Q: Does Jepsen use a fault injector? 
A: Yes! They test against real systems, running in the wild but simulate faults (at the control node?). 

How does Jepsen find minimal repros? Is it manual? 
They use [elle] to find conflicts in traces. Elle should be a part of default antithesis?
Or maybe "default antithesis - distsys"? It will require logging in a format elle understands.
Does it require recording traces from the perspective of a "control node"?

*What can Antithesis learn from Jepsen?*

Q: Does Antithesis use something like [elle]?
Q: Should PS use elle or should it be deeply integrated? 

Q: What does *determinism* buy us?
Perfect reproducibility. Tree Fuzzing and Minimization? 

We've got a *fault injector* (FI) that simulates things like dropped and delayed messages and network partitions.
The hope is that by simulating these behaviors in a controlled environment we 

*How does FDB find bugs?*

FoundationDB turned the ad-hoc and bug ridden authoring of DS  

[liskov00]: https://pmg.csail.mit.edu/papers/icde00.pdf
[Jepsen]: https://jepsen.io/
[elle]: https://github.com/jepsen-io/elle

Communicating Sequential Processes.
