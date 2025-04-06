This little language ended up pretty simple.
Not by design but by dumb iteration and trying to do first, dumbest thing that came to mind.


Instead of have FnCall have a name symbol that points to an interned true function we actually
copy the function pointer into every op of the program. The `value` has type `any` and stores
the func ptr. The name is a unique within the Library `map[name]->FnCall`.

    type FnCall struct {
    	value  any
    	name   string
    	ptypes []Type
    	rtype  Type
    }

Instead of having Programs and Pieces we just have two different types Program and UncheckedProgram
with the same structure, but Programs must maintain an invariant that they aren't missing any syms
and can be passed to `evalProgram()` directly. Unchecked programs _may_ have this property but we
don't know it. There are validators and rectifiers that will repair / fill-in the blanks thus either
converting or confirming that we have a true Program.

A Program is a list of Statements.
A Statement maps input syms to a single output sym via a FnCall with appropriately typed args.

A Type is a string (for now).
A Sym is a sting (and always will be).

There is no scoping so Syms don't have to refer to their scope level.
There are no LitNums or LitStr, but because `FnCall.value` has type any we can construct thunks that
close over arbitrary values and store them there while in the middle of program construction.

There are no `FnDef`s, but maybe there are other ways of getting the same effect, either through
brute force copy-pasting chunks of programs, or by defining FnCalls that are actually built out
of more complex funcs internally. 

evalProgram() returns a ValueMap : map[Sym]Value. 
A Value is basically just an any that remembers what Type it thinks it is.
We can initialize the ValueMap during evalProgram with Values that
we want to make available to the program, but didn't place in the program body.

No Target interface or Campaign interface. No eval context. No modules. One big namespace.
No global Logger or threaded, stateful logger. No JSON.

CAN in principle write programs as text, parse and execute them... 


ProGen              ; TacticsLang 
======              ; =========== 
Value, Fun          ; Func, Call, LitNum, LitStr 
Statements          ; Ops, OpGeneral
Program (Unchecked) ; Block, Piece, "fragment"?
Sym = String        ; Syms have Scope
Type = String       ; CatalogType, TypedSym

Program Sampler     ; Piece Sampler + Program Sampler
Program Mutation    ; Piece Mutation
Genetic Optimizer   ; PieceSampler
i64                 ; f64 with casting
fn_lib, val_lib     ; Catalog, RootCatalog
???                 ; EvalContext
???                 ; PieceBuilder
No state yet        ; Target trait. Target/env has state. Tree fuzzing makes sense.
???                 ; Dispatcher, Logger

In process          ; Interpreter + Shared Heap


No type enforcement between the function defined in a Fun and the stated arg/p types.
TL DOES have this! It's probably necessary to provide to customers.
Go progen could provide a function that uses reflection to confirm that the types are good. 


