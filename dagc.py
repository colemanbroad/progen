
from collections import namedtuple, Counter
import pdb

Fn = namedtuple("Fn", ["atypes", "rtype"])

frags = {
    "f1": Fn([None], "a"),
    "f2": Fn(["a"], "b"),
    "f3": Fn(["a", "b"], "b"),
    "f4": Fn(["a", "b"], "c"),
    "f5": Fn(["a", "b"], "c"),
    "f6": Fn(["c", "c"], "d"),
}

def count():
    print(frags)
    # first we want to iterate through the fragments and find the Types that can be built
    # from nothing, and count how many ways we can build them.
    t0 = Counter([fn.rtype for f,fn in frags.items() if fn.atypes[0] is None])
    print(t0)
    # Now what about lvl 2? How many ways can we build using only things from lvl1 ?
    t1 = Counter([fn.rtype for f,fn in frags.items()
                 if set(fn.atypes) <= set(t0) and fn.rtype not in set(t0)])
    print(t1)
    tmpset = set(t0) | set(t1)
    t2 = Counter([fn.rtype for f,fn in frags.items()
                 if set(fn.atypes) <= tmpset and fn.rtype not in tmpset])
    print(t2)

count()
