---
title: Compilers Brief & Brisk
type: blogs
---

# Compilers Brief & Brisk

Most compilers out there follow a particular architecture:

![](https://upload.wikimedia.org/wikipedia/commons/c/cc/Compiler_design.svg)

### Introduction
In this article I intend to dissect this architecture piece by piece in some detail. 

Consider this article a supplement to the plethora of resources out there on compilers. It exists as a self contained resource to get your toes wet in the world of programming language design and implementation.

The audience for this article is someone who has very limited knowledge as to how a compiler works, i.e. you know that they compile into assembly at most. Though I do presume that the reader has a good understanding of data structures & algorithms. 

It is by no means reflective of modern 'production' compilers with millions of lines of code! But rather a very brief/brisk 'compilers for dummies' resource to get an idea of what goes on in a compiler.

### Disclaimer!
I would like to disclaim that I am by no means an expert in compilers! I don't have a doctorate in Compilers, I did not study this at an academic level in any way - most of what I am sharing is what I have learned in my spare time for fun. In addition, I am not claiming what I write to be the de-facto approach for engineering a compiler, but rather introducing approaches that would be applicable for a small toy compiler.

## Frontend
Referring back to the previous diagram, the arrows on the left pointing into the box are programming languages we all know and love today - like C! The frontend looks something like this:

    Lexical Analysis -> Parser

### Lexical Analysis
When I was first learning about compilers and language design, this was described to me as a 'fancy way of saying tokenization'. So let's go with that. The 'lexer' typically takes input in a form of strings or a stream of characters, and recognizes patterns in those characters to cut them up into tokens. 

In the case of a compiler, the input string would be a program that the programmer would write. This would be read in from a file into a string, and then the lexer would tokenize the programs source code.

    enum TokenType {
        Identifier,
        Number,
    };

    struct Token {
        std::string Lexeme;
        TokenType type;

        // ...
        // It's also handy to store things in here
        // like the position of the token (start to end row:col)
    };

In the above snippet of some kind of C-like language, we have a Token structure which contains the aforementioned `lexeme`, as well as a TokenType to distinguish what kind of lexeme is stored.

Note: This post isn't necessarily a guide to creating a language with code samples, but I will add a few bits code snippets here and there to assist translating these ideas.

Lexers are usually the easiest part of the compiler to make, in fact the entire frontend is usually quite simple relative to the other pieces of the puzzle. Though this depends on how hard you can make it for yourself 😉

Take the following piece of C code:

    int main() {
        printf("Hello world!\n");
        return 0;
    }

If you were to read this from a file into a string, and then linearly scan through the string you could probably cut these into tokens. We naturally identify these tokens looking at the language ourselves, e.g. it's clear that `int` is a "word", and `0` in the return statement is a "number".
The lexer does the same kind of thing, and we can go into as much detail as necessary to ease the process later. For example, you could lex:

    0xdeadbeef
    1231234234
    3.1412
    55.5555
    0b0001

As "numbers", or you could categorize them further as:

    0xdeadbeef      HexNumber
    1231234234      WholeNumber
    3.1412          FloatingNumber
    55.5555         FloatingNumber
    0b0001          BinaryNumber

As for defining "words", it can be difficult. Most languages define a word as a grouping of letters and digits, the identifier typically must _start_ with a letter (or an underscore), e.g.

    123foobar := 3
    person-age := 5
    fmt.Println(123foobar)

Is **not** valid Go code as it is probably parsed into these tokens:

    Number(123), Identifier(foobar), Symbol(:=), Number(3) ...

Most identifiers that we encounter are of the form:

    foo_bar
    __uint8_t
    fooBar123

Lexers will also have to deal with other problems, e.g.:

* Whitespace,
* Comments - Multi Line and Single Line,
* Identifiers,
* Numbers, bases, number 'formatting', e.g. `1_000_000`,
* Input encoding, e.g. supporting UTF8 rather than ASCII

... and before you think about using regular expressions to do this, I would not recommend it! It's much easier to write a lexer from scratch, but I highly recommend reading this [blog post](https://commandcenter.blogspot.com/2011/08/regular-expressions-in-lexing-and.html) from our lord and saviour Rob Pike. 
Though, there are many articles on why Regex is not the right tool for the job so I think I'll skim over that segment for this article.
It's also a lot more fun to write a lexer than it is pulling your hair out over a long winded regular expression you have pasted into regex101.com at 5:24 in the morning.

My first 'programming language' I used the `split(str)` function to tokenize my input - I didn't get very far.

### Parsing
Parsing is a bit more of a complicated beast compared to Lexical Analysis. There are many kinds of parsers out there and parser generators. This is where things start to get a bit more serious.

In compilers, a parser will usually take an input of tokens, and produce a tree of some sort. This could be an 'Abstract Syntax Tree', or a 'Parse Tree'. Both of which are similar at the core, but do share differences.

You could think of these stages so far as functions:

    fn lex(string input) []Token {...}
    fn parse(tokens []Token) AST {...}

    let input = "int main() { return 0; }";
    let tokens = lex(input);
    let parse_tree = parse(tokens);
    // ....

Compilers are usually built up in lots of little components which take inputs and mutate them or convert them into different outputs. Which is partly why functional languages are very good for creating compilers! And also the sweet pattern matching and usually very extensive standard libraries.
Fun fact: the first implementation of the [Rust](https://en.wikipedia.org/wiki/Rust_(programming_language)) compiler was implemented in OCaml.
And a piece of advice is to keep these components as self contained and simple as possible, keeping everything modular simplifies the entire process. I feel as if this philosophy applies to many aspects of software engineering however.

### Trees!
#### Parse Tree
WTF is... - a parse tree? Sometimes referred to as a 'syntax tree', is a much more dense tree that represents the source program. They contain _all_ (or most) of the information of the input program, usually matching what is described in the grammar of your language. Each node in the tree would be a nonterminal or terminal in the grammar, e.g. a NumberConstant node or a StringConstant node.

#### Abstract Syntax Tree
The Abstract Syntax Tree (AST) is, as the name suggests, an 'abstract' syntax tree. The Parse Tree contains a lot of (possibly superfluous) information about your program. The point of the AST is that we don't need all of this information to do our job. It throws away a lot of the useless structural/grammatical information that doesn't contribute to the semantics of the program.
For example, perhaps you have an expression in your tree _Parse Tree_ like `((5 + 5) - 3) + 2`. You would store the parenthesis in the Parse Tree, and maybe that the values 5, 5, 3, and 2 are atoms, but once you can derive the associations, you can abstract away these details in the AST as we only need to know the values (the numbers) and their operators as well as the order of the operations.

Here's another free for re-use image I found that shows the AST for `a + b / c`. It's a bit large, sorry! I blame svbtle for that.

![a + b / c in the form of a tree](https://upload.wikimedia.org/wikipedia/commons/6/68/Parsing_Example.png)

An AST could be represented as such:

    interface Expression { ... };

    struct UnaryExpression {
        Expression value;
        char op;
    };

    struct BinaryExpression {
        Expression lhand, rhand;
        string op; // string because the op could be more than 1 char.
    };

    interface Node { ... };

    // or for something like a variable
    struct Variable : Node {
        Token identifier;
        Expression value;
    };

This is a very limited representation of an AST, but you could hopefully see how you would structure your nodes.

As for parsing them, you could have a procedure like:

    Node parseNode() {
        Token current = consume();
        switch (current.lexeme) {
        case "var":
            return parseVariableNode();
        // ...
        }
        panic("unrecognized input!");
    }

    Node n = parseNode();
    if (n != null) {
        // append to some list of top level nodes?
        // or append to a block of nodes!
    }

And hopefully you get the gist of how it would recursively parse other nodes from the top level constructs in the language. Though I'll leave that to you to learn about the specifics of implementing a recursive descent parser.

### Grammars
To parse from a set of tokens into an AST can be a tricky task. Usually you would start with some kind of grammar for your language.

Grammar is basically a definition of how a language is structured. There are a few languages for defining languages, which can be described (or bootstrapped) with themselves.

[Extended Backus-Naur Form](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form) (EBNF) is an example of a language for defining languages. It is based off [BNF](https://dlang.org/spec/grammar.html) which is a bit more angle bracket-y.

Here's an example of some EBNF taken from the Wikipedia Article:

    digit excluding zero = "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" ;
    digit                = "0" | digit excluding zero ;

Production rules are defined, which tell the reader what pattern of terminals make up what 'nonterminal'. Terminals are part of the grammars alphabet, e.g. the token "if", or in the example above "0" and "1" are terminals. Non-terminals are the opposite, they are on the left of the production rules, and can be considered variables or 'named references' to a grouping of terminals _and_ non terminals.

Many languages have specifications, which contain their grammars that you can read. Here is the spec for [Go](https://golang.org/ref/spec#Function_declarations), and [Rust](https://doc.rust-lang.org/reference/), as well as [D](https://dlang.org/spec/grammar.html).

#### Recursive Descent Parsing
The easiest approach is using a 'recursive descent parser'. There are many approaches to parsing, and this is one of them.

Recursive descent is a top down parser built from a set of recursive procedures. It's a much simpler to write a parser, given that your grammar has no [left recursion](https://en.wikipedia.org/wiki/Left_recursion). For most hobby/toy languages, it's a sufficient technique for parsing. GCC uses a hand-written recursive descent parser, though it has used YACC before.
Though there can be issues with parsing these languages, especially something like C, where:

    foo * bar

Could be interpreted as:

    int foo = 3;
    int bar = 4;
    foo * bar; // unused expression

Or it could be interpreted as:

    typedef struct {
        int b;
    } foo;
    foo* bar;
    bar.b = 3;

[Clang](https://clang.llvm.org/features.html) also uses a recursive-descent parser in its implementation:

> Because it is plain C++ code, recursive descent makes it very easy for new developers to understand the code, it easily supports ad-hoc rules and other strange hacks required by C/C++, and makes it straightforward to implement excellent diagnostics and error recovery.

A few alternative approaches to parsing that are worth reading into are:
* (top down) LL, Recursive descent
* (bottom up) LR, shift reduce, recursive ascent, ...

#### Parser Generators!
Parser generators are a very good approach to take too. There are trade offs though, as there are with any choice you make when creating a piece of software.

Parser generators are usually very fast, they are a lot easier than writing your own parser (and getting a performant result from it), though they are usually not very user friendly, and don't allow for great error messages. In addition, you then have to learn how to use your parser generator, as well as when it comes to bootstrapping your compiler, you probably have to bootstrap your parser generator too.

[ANTLR](https://www.antlr.org/) is an example of a parser generator, though there are plenty more out there.

I think they are a tool for when you don't want to focus on writing the frontend, but would rather get on with writing the middle and the backend parts of your compiler/interpreter, or dealing with whatever else you want to parse.

#### The Application of Parsing
If you haven't guessed yet. Just the frontend of a compiler (lex/parse) is _very_ applicable to other problems:

* Syntax highlighting;
* HTML/CSS parsing for a layout engine;
* Transpilers: TypeScript, CoffeeScript;
* Assemblers;
* REGEX;
* Screen scraping;
* URL parsing;
* Formatting tools like `gofmt`;
* SQL parsing; 

And much more.

## The Middle
Semantic analysis! Analyzing the semantics of the language is one of the harder parts of compiler design.
It involves ensuring that the input programs are _correct_. Without semantic analysis the programmer would have to be trusted to write code that is correct all the time - which is impossible!
Not only this, but it can be extremely difficult to compile a program if you can't analyze that the semantics are correct in this analysis phase of the compiler.

I remember reading a diagram a while ago that showed the percentages of the front, middle, and back ends of a compiler and how they were split up and a while ago it was something along the lines of...

    F: 20% M: 20%: B: 60%

But nowadays I feel it's more like this:

    F: 5% M: 60% B: 35%

Frontends are mostly the job of a generator or can be done very quickly to a language that is mostly context free and/or doesn't have any ambiguities in the grammar (throw a recursive descent parser at it!).

With technology like LLVM, most of the work of optimisation can be offloaded onto the framework, already providing a plethora of optimisations out of the box! 
So that leads us with semantic analysis which is a very integral part of the compilation phase. 
For example, a language like [Rust](https://en.wikipedia.org/wiki/Rust_(programming_language)) with its ownership memory model, the compiler mostly acts as a big fancy machine that performs all types of static analysis on the input forms. Part of the job is converting the input into a more manageable form to do this analysis.

Because of this, semantic analysis is definitely eating up a lot of the focus for compiler architecture since a lot of the tedious ground work like optimising generated assembly, or readng input into an AST is done for you. 

### Semantic Passes
Most compilers will perform a large amount of 'semantic passes' over the AST (or some abstract form to represent the code) during semantic analysis. [This](https://blogs.msdn.microsoft.com/ericlippert/2010/02/04/how-many-passes/) article goes into detail about most of the passes that are performed in the .NET C# compiler (from 2010).
Note: this wont be an exhaustive list of every pass!

### 'Top Level' Declaration Pass
The compiler will go over all of the 'top level' declarations in the modules and acknowledge that they exist. It's top level as it doesn't go any further into blocks, it will simply declare what structures, functions, etc. exist in what module.

### Name/Symbol Resolution
This pass will go through all of the blocks of code in the functions, etc. and resolve them, i.e. find the symbols that they resolve to. This is a common pass, and is usually where the error `No such symbol XYZ` comes from when you compile your Go code.

This can be a very tricky pass to do, especially if you have cyclic dependencies in your dependency graph. Some languages do not allow for cycles, e.g. Go will error if you have packages that form a cycle. Cyclic dependencies are generally considered to be a side effect of poor architecture in a program.

### Type Inference
This is a pass in which the compiler that will go through all variables and infer their types. Type inference can be done using a process called 'unification', or 'type unification'. Though you can have some very simple implementations for simpler type systems.

An example could be assigning expression nodes a type, e.g. an IntegerConstantNode would have the type IntegerType(64). And then you have some function e.g. unify(t1, t2) which picks the widest type that can be used for inferring the type of more complex expressions. Then it's a matter of assigning the left hand variable the right hand values inferred type.

### Mutability Pass
Languages like Rust will potentially have a pass to ensure that immutable values are not re-assigned:

    let x = 3;
    x = 4; // BAD!

    let mut y = 5;
    y = 6; // OK!

This pass in the compiler will run through all of the blocks and functions and ensure that they are 'const correct', i.e. we are not mutating anything we shouldn't, and that all values passed to certain functions are constant or mutable where need be.

This is done with symbol information that is collected from prior passes. A symbol table is built up in the semantic pass which contains information like the token name, and whether the variable is mutable or not. It could also contain other information, e.g. in the case of C++, if the symbol is `extern` or `static`.

### Symbol Tables
A symbol table or 'stab' is a lookup table for symbols that exist in your program. These exist per scope and contain all of the symbol information for the said scope.
The symbol information is contains properties like the name of the symbol, it could contain the type information too, as well as if the symbol is mutable or not, if it should be externally linked, is it in static memory, etc.

### Scope
Scope is an important concept in programming languages. Of course your language doesn't necessarily have to allow for nested scope, it could all be in one namespace!
Though representing Scope is an interesting problem in compiler design, scope behaves like (or is) a stack data structure in most c-like languages.
Usually, you would push and pop scope, and normally these would control names, e.g. allowing for shadowing of variables:

    { // push scope
        let x = 3;
        { // push scope
            let x = 4; // OK!
        } // pop scope
    } // pop scope

And this can be represented in a few ways:

    struct Scope {
        Scope* outer;
        SymbolTable symbols;
    }

Somewhat irrelevant, but interesting reading/knowledge: [Spaghetti stack](https://en.wikipedia.org/wiki/Parent_pointer_tree). This is a data structure that was used to store the scopes in their counterpart block AST nodes. It's often referred to as a spaghetti stack!

### Type Systems
Many of these headings could be their own articles, but I feel like this heading probably takes the cake for that.
There is a lot of information on type systems out there, and there are many kinds of type systems, and a lot of heated debate about everything. I wont gloss over this topic too much at all, but I will link to this excellent article by [Steve Klabnik](https://blog.steveklabnik.com/posts/2010-07-17-what-to-know-before-debating-type-systems).
Though the point of this header, is that the type system is something that is enforced and defined semantically in the middle of the compiler with aid of the compilers representations as well as the analysis of these representations. 

### Ownership
Ownership is a concept that is becoming more and more prevalent in the programming world. Ownership and move semantics are principles in the language [Rust](https://en.wikipedia.org/wiki/Rust_(programming_language)) &mdash; and hopefully more to come. There are many forms of static analysis performed on Rust code that checks that input conforms to a set of rules with regards to memory: who owns what memory, when the memory dies, and how many references (or borrows) to these values/memory there are.

The beauty of Rust is that this is all enforced at compile time, during the middle of the compiler, so there is no garbage collection or reference counting forced upon the programmer. These semantics are offset to the type system and can be enforced before the program even exists as a compiled binary.

I can't speak on the internals of how this all works, but I can tell you that it is the work of static analysis and some cool research by the folks at Mozilla and the people behind [Cyclone](https://en.wikipedia.org/wiki/Cyclone_%28programming_language%29).

### Control Flow Graphs
To represent a programs flow, we use Control Flow Graphs (CFG), which contains all the paths that may be traversed during the execution of a program. This is used in semantic analysis to handle dead code elimination, e.g. blocks that wont ever be reached, or functions, or even modules. It can also be used to determinte if a loop wont stop iterating, for example. Or unreachable code, e.g. you call a `panic` or return in a loop and the code outside of the loop doesn't get to execute. [Data Flow Analysis](https://en.wikipedia.org/wiki/Data-flow_analysis) plays a prominent role during the semantic phase of a compiler, and it's worth reading up on the types of analysis you can do, how it works, and what optimisations can come from it.

## Backend
![a barren wasteland](https://upload.wikimedia.org/wikipedia/commons/8/85/Hutong_Barren_Wasteland.jpg)

The final part in our architecture diagram.

This is where most of the work is done to produce our binary executable. There are a few ways to do this, which we will discuss in the later segments of this article.

The semantic analysis phase doesn't necessarily have to mutate a lot of the information on the tree, and it's probably a better idea not to in terms of avoiding some spaghetti mess.

### A note on transpilers
Transpilers are another form of compiler, in which the compiler transpiles into another 'source level' language, e.g. you could write something that compiles into C code. I think this is somewhat pointless though if your language doesn't have a lot to offer on top of the language its compiling to. It mostly seems to make sense for languages that are relatively high level or the language itself is limited.
However, compiling to C is a very established habit in the history of compilers, in fact the first C++ compiler 'Cfront' compiled into C code.

A good example for this is JavaScript. TypeScript (and many other languages) transpile into JavaScript to introduce more features to the language, and most importantly a sensible type system with various amounts of static analysis to catch bugs and errors before we encounter them at runtime.

This is one type of 'target' for a compiler, and it's usually the easiest as you don't have to think in more lower level concepts about assigning variables, or handling optimisations, etc. as you are mostly piggy backing on top of another language. Though the obvious downside is that you have a lot of overhead, and are usually confined within the language you are compiling to.

### LLVM
Many modern compilers will opt for using LLVM as their backend: Rust, Swift, C/C++ (clang), D, Haskell.
This can be considered the 'easier route', as most of the work is done for you in supporting a wide variety of architectures, and providing an insurmountable level of optimisation. In contrast to the aforementioned route of transpilation, you get quite a lot of control with LLVM too. More so than if you were to compile to C. For example, you can decide how large types will be, e.g. 1 bit, 4 bits, 8 bits, 16 bits - which isn't as easy in C, and sometimes not possible, or not even defined for certain platforms.

### Generating Assembly
Generating code directly for a specific architecture, i.e. machine code or assembly is technically the most common route with countless languages opting for this route.

Go is an example of a modern language that does not take advantage of the LLVM framework (as of writing this). It generates code for a few platforms/architectures.

There are lots of pros and cons to this, though with technology like LLVM available it's unwise to generate assembly code yourself as it is unlikely a toy compiler that has its own assembly backend would surpass LLVMs level of optimisation for one platform let alone multiple.

That being said, a considerable benefit of generating assembly yourself is that your compiler will likely be a lot faster than if you were to use a framework like LLVM where it has to build your IR, optimise it, etc. and then eventually write it out as assembly (or whatever target you pick).

Regardless, it's still enjoyable to attempt. And is especially interesting if you wanted to learn more about programming in assembly or the lower levels of how languages work. The easiest way to approach this is to walk the AST, or walk the generated IR (if you have one) and 'emit' assembly instructions to a file using `fprintf` or some file writer utility. This is how [8cc](https://github.com/rui314/8cc) works.

### Bytecode Generation
Another option is generating bytecode for some kind of virtual machine or a bytecode interpreter. Java is a prime example of this, in fact the [JVM](https://en.wikipedia.org/wiki/Java_virtual_machine) has spawned an entire family of languages that generate bytecode for it, e.g. Kotlin.

There are many benefits of generating bytecode, the main reason for Java was for portability. If you can have your virtual machine run anywhere, any code that executes on the virtual machine will run anywhere too. And it's a lot simple to make something like an abstract set of bytecode instructions run on machines than it is to target 50 bajillion computer architectures.
As far as I know the JVM will also JIT frequently run 'hot code' into a native function, and other such JIT tricks to squeeze out extra performance from code.

### Optimisations
Optimisations are integral to a compiler, no one wants slow code! They will usually be the larger part of the backend, and there is a large magnitude of research on squeezing out the extra bits of performance with code. 
If you ever compile some C code and run it with all optimisations on full flex, it can be amazing what kind of madness it can produce. [godbolt's compiler explorer](https://godbolt.org/) is a great tool to look into how existing compilers generate their code, what instructions relate to what source code, as well as you can specify certain levels of optimisations, targets, versions of compilers, etc.

A good start if you are ever writing a compiler is to write simple programs in C and turn off all optimisations, as well as strip the debug symbols, and have a look at what code GCC generates. It can be a handy reference if you ever get stuck!

The importance with optimisations is that you can trade off accuracy of the program for speed, and finding right balance can be difficult. Some optimisations are also very specific on their use case and can in some instances produce the wrong result. These optimisations usually don't find themselves in production compilers for obvious reasons!

An interesting comment taken from the [HN thread](https://news.ycombinator.com/item?id=19756087) on this article from user 'rwmj' is that you only need around 8 optimisation passes to get 80% of the best case performance from your compiler, all of which were catalogued in 1971! This was in a [compilers talk](http://venge.net/graydon/talks/CompilerTalk-2019.pdf) from Graydon Hoare, the mastermind behind Rust.

### IR
Having an intermediate representation (IR) is not required, but definitely beneficial. You could generate code from the AST, though it can become quite tedious and messy to do so, as well as it's quite difficult to optimise.

An IR can be thought of as a higher level representation of the code that you are generating for. It must be very accurate to what it represents, and must contain all of the information necessary to generate the code.

There are certain types of IR, or 'forms' you can make with the IR to allow for easier optimisations. One example of this is SSA, or Static Single Assignment, in which every variable is assigned exactly once.
Go builds an SSA based IR before it generates code. LLVM's IR is built upon the concept of SSA to provide its optimisations.

SSA provides a few optimisations by nature, for example constant propagation, dead code elimination, and (a big one) register allocation.

#### Register Allocation
Register allocation is not a necessity when it comes to generating code, but an optimisation. One abstraction we take for granted is that we can define as many variables as required for our programs. In assembly, however, we can either make use of the finite amount of registers [usually 16 to 32] available (and keep track of them in our heads), or we can spill to the stack.

Register allocation is an attempt to find what variables can go in what registers at what point of time (without overwriting other values). This is much more efficient than spilling to the stack, though can be quite expensive and impossible for a computer to calculate the perfect solution.
A few algorithms for register allocation are: graph colouring, which is as computationally hard problem (NP-complete). Or a linear scan which will scan the variables to determine their liveness ranges - as opposed to graph colouring which requires the code is in graph form to calculate the liveness of variables.

### Things to consider
There is a vast amount of information on compilers. So much to cover that it would not fit nicely into this article. That being said, I wanted to write up, or at least mention, a few bits and pieces that should be considered for any of your future endeavours.

#### [Name Mangling](https://en.wikipedia.org/wiki/Name_mangling)
If you are generating assembly where there isn't really any scope or namespace, you will have an issue with conflicting symbols in a lot of cases. Especially if your language supports function overloading or classes, etc.

    fn main() int {
        let x = 0;
        {
            let x = 0;
            {
                let x = 0;
            }
        }
        return 0;
    }

For instance, the previous example (if the variables don't get optimised out anyway 😉) you will have to mangle the names of those symbols so that they would not conflict in the generated assembly. Name mangling usually indicates type information too, or it could contain scope information, etc.

#### Debug Information
Tools like LLDB usually integrate with standards like [DWARF](https://en.wikipedia.org/wiki/DWARF). Another excellent feature of LLVM is that you get some relatively easy integration with the existing GNU debugger tool via. DWARF. Your language would probably need a debugger, it's always easiest to use someone else's, unless you were to roll your own.

#### Foreign Function Interface ([FFI](https://en.wikipedia.org/wiki/Libffi))
There is usually no escape from libc, you should probably read up on this and think about how you would incorporate this in your language. How will you hook into C code, or expose your languages code to C?

#### Linking
Writing your own linker is a task of its own. When your compiler generates code, does it generate a machine code of some sort (i.e. into a `.s`/`.asm` file)? Does it write the code directly to an object file? Jonathan Blow's programming language [Jai](https://www.youtube.com/watch?v=TH9VCN6UkyQ) supposedly writes all of the code into a single object file. There are many different approaches to this with varying trade offs.

## Further Reading

* [Jack Crenshaw](https://compilers.iecc.com/crenshaw/) - my personal gateway into the realm of programming language implementation
* [Crafting Interpreters](https://craftinginterpreters.com/)
* [An Introduction to LLVM (with Go)](https://blog.felixangell.com/an-introduction-to-llvm-in-go) - me!
* [PL/0](https://en.wikipedia.org/wiki/PL/0)
* The Dragon Book - a classic book that has _everything_
* [8cc](https://github.com/rui314/8cc)