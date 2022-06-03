---
title: An Introduction to LLVM in Go
type: blogs
---

# An Introduction to LLVM in Go
LLVM is an infrastructure for creating compilers. It was initially created by Chris Lattner in 2000, and released in 2003. Since then it has evolved into an umbrella project that has a wide array of tools such as the LLVM Linker `lld`, LLVM Debugger `lldb`, and so on.

The banner feature of LLVM is its intermediate representation, commonly referred to as the `LLVM IR`. The idea of LLVM is that you can compile down to this IR, and then this IR can be JIT compiled, interpreted, or compiled into native assembly for the machine it's running on. The primary target of this IR is compilers, in fact there are many compilers out there that use LLVM: clang and clang++ for C and C++ respectively, `ldc2` for the D programming language, the Rust language, Swift, etc. There are projects like [emscripten](https://github.com/kripken/emscripten), which can compile LLVM BC (LLVM bitcode) into javascript to be executed in the browser. 

Typically in compiler design you would have to worry about register allocation, generating code for different architectures, and producing good code that is well optimized. The beauty of LLVM is that it does this for you. LLVM features a vast collection of optimisations, can target a variety of architectures, and has a nice API that makes code generation a whole lot more simpler.

## LLVM IR
Now lets take a swift look at the LLVM IR. If you take your average C program and run it through `clang` with the `-emit-llvm` and `-S` flag, it will produce an `.ll` file. This file extension means that it's LLVM IR.

Here's the C code I'm going to compile into LLVM IR:

```c
int main() {
	int a = 32;
	int b = 16;
	return a + b;
}
```

I've ran it through clang, specifying to not optimize anything: `clang test.c -S -emit-llvm -O0`:

	define i32 @main() #0 {
		%1 = alloca i32, align 4
		%a = alloca i32, align 4
		%b = alloca i32, align 4
		store i32 0, i32* %1
		store i32 32, i32* %a, align 4
		store i32 16, i32* %b, align 4
		%2 = load i32, i32* %a, align 4
		%3 = load i32, i32* %b, align 4
		%4 = add nsw i32 %2, %3
		ret i32 %4
	}

I've omitted a lot of the excess code for simplicity sake. If you have a look at the IR, it looks much like a more verbose, readable assembly. Something you may note is that the IR is strongly typed. There are type annotations everywhere, instructions, values, functions, etc.

Let's step through this IR and try get a grasp of what is going on. First off we have a function with a syntax much like a C-style function with the brace, type, name and the parenthesis for arguments.

In our function we have a bunch of values and instructions. In the IR we see 5 instructions, `alloca`, `store`, `load`, `add`, and `ret`.

Let's dissect the IR part-by-part to help understand how this works. Note that I've ignored a few things here, namely the alignment, and the `nsw` flags. You can learn more about those in the LLVM documentation, I'll be explaining the underlying semantics.

### Locals
Before we get onto the instructions, you should know what a local is. Locals are like variables. They are denoted with a percent symbol `%`. As the name suggests, they are local to the function they are defined in. This means that they cannot be modified/referenced outside of the function that declares them.

### `alloca`
This instruction will allocate memory in the stack frame. This memory is freed when the function is returned. The instruction returns a value, this is why we assign it to `%a`, etc. The value it returns is a pointer to the memory that is allocated. For example:

	%a = alloca i32, align 4

This instruction allocates space for a 32 bit signed integer on the stack. The pointer is stored in the local `a`.

### `store`
The store instruction will change the value at the given pointer to contain the given value. Here's an example to simplify the explanation:

	store i32 32, i32* %a, align 4

Here we tell LLVM to store the value 32 of type `i32` into the local `a` of the type `i32*` (a pointer to an i32). This instruction returns void, i.e. it returns nothing and cannot be assigned to a local.

### `load`
Finally, the `load` instruction. This instruction will return the value at the given memory address:

	%2 = load i32, i32* %a, align 4

In the example above, we load a value of type `i32` from the memory address `a` (which is a pointer to an `i32`). This value is stored into the local `2`. We have to load values because we can't dereference 

We now know what the instructions mean, so hopefully you should be able to read and understand more than half of the IR above. As for the rest of the instructions, they should be relatively straight-forward. `add` will perform addition on the given values and return the result. The `ret` instruction specifies the value to return from the function.

### LLVM API
LLVM provides an API for building this IR. The initial API is in C++, though there are bindings to a variety of languages from Lua, to OCaml, C, Go, and many more. 

In this article, I'll be using the Go bindings. Though before we start building some IR, we need to know and understand a few details:

#### Modules
A module is a group of definitions and declarations. This is the container, and we must make one. Typically modules are created per-file, so in our C example, that file was a module.

We create a module like so. We pass a string as the name of the module, we're going to call ours "main" as it's the main module:

```go
module := llvm.NewModule("main")
```

### Types
LLVM provides a wide variety of types, from primitive types like bytes, integers, floating point, to more complex types like Structures, Arrays, and Function Types.

There are some built in types, in the format of `TypeWidthType()`, so for example `Int16Type` is an integer with a width of 16 bits.

```go
foo := llvm.Int16Type()
bar := llvm.Int32Type()
```

We can specify arbitrary bit-widths:

```go
toast := llvm.IntType(32)
```

An array is like so:

```go
ages := llvm.ArrayType(llvm.Int32Type(), 16)
```

This is an array of 16 32-bit integers.

### Values
LLVM values can be returned from instructions, though they can be constants, functions, globals, ...

Below we create a constant integer of type `i32`, with the value `666`. The boolean parameter at the end is whether to sign extend.

```go
foo := llvm.ConstInt(llvm.Int32Type(), 666, false)
```

We can create floating point constants:

```go
bar := llvm.ConstFloat(llvm.FloatType(), 32.5)
```

And we can assign these values to variables, or pass them to functions, and so on. Here we create an add instruction that adds two constant values:

```go
a := llvm.ConstInt(llvm.Int32Type(), 12)
b := llvm.ConstInt(llvm.Int32Type(), 24)
c := llvm.ConstAdd(a, b)
```

### Basic Blocks
This is slightly different to how you may expect. In assembly, we use labels for functions, and control flow. LLVM is similar to this, though we have an explicit syntax for functions. How do we control the flow of our program? We use basic blocks. So the IR would look like this:

	define i32 @main() {
	entry:
		...
	0:
		...
	1:
		...
	}

We have our main function, and inside of this function we have three basic blocks. The entry block, and then the 0 and 1 block. You can have as many basic blocks as you want. They are used for things like jumping around, for instance looping, if statements, and so on.

In the Go bindings for LLVM, we define a basic block like so:

```go
llvm.AddBasicBlock(context, "entry")
```

Where the context is the function that we want to add the block to. This is _not_ a function type. We'll discuss this later, though.

## IR Builder
The IR Builder will create our IR for us. We feed it values, instructions, etc. and it will join them all together. The key part of the builder is that we can use it to reposition where we build, and append instructions in different places.

We can use this builder to append instructions to our module. Below we setup a builder, make a function and entry block, then append some simple instructions to store a constant:

```go
builder := llvm.NewBuilder()
// create a function "main"
// create a block "entry"

foo := builder.CreateAlloca(llvm.Int32Type(), "foo")
builder.CreateStore(foo, llvm.ConstInt(llvm.Int32Type(), 12, false))
```

This produces IR like so:

	define i32 @main() {
	entry:
		%foo = alloca i32
		store i32 12, i32* %foo
	}

### Functions
Functions are a type in LLVM. We need to specify a few things when we define this function type: the return type, parameter types, and if the function is variadic, i.e. if it takes a variable amount of arguments.

Here's our main function as we've seen it thus far:

```go
main := llvm.FunctionType(llvm.Int32Type(), []llvm.Type{}, false)
llvm.AddFunction(mod, "main", main)
```

The first parameter is the return type, so a 32 bit integer. Our function takes no parameters, so we pass an empty type array. And the function is not variadic, so we pass false in the for the last argument. Easy right?

The AddFunction will add the function to the given module as the given name. We can then reference this later (it's stored in a key/value map) like so:

```go
mainFunc := mod.NamedFunction("main")
```

This will lookup the function in the module.

Now we can piece together everything we've learned so far:

```go
// setup our builder and module
builder := llvm.NewBuilder()
mod := llvm.NewModule("my_module")

// create our function prologue
main := llvm.FunctionType(llvm.Int32Type(), []llvm.Type{}, false)
llvm.AddFunction(mod, "main", main)
block := llvm.AddBasicBlock(mod.NamedFunction("main"), "entry")

// note that we've set a function and need to tell
// the builder where to insert things to
builder.SetInsertPoint(block, block.FirstInstruction())

// int a = 32
a := builder.CreateAlloca(llvm.Int32Type(), "a")
builder.CreateStore(llvm.ConstInt(llvm.Int32Type(), 32, false), a)

// int b = 16
b := builder.CreateAlloca(llvm.Int32Type(), "b")
builder.CreateStore(llvm.ConstInt(llvm.Int32Type(), 16, false), b)
```

So far so good, though, because an `alloca` returns a pointer, we can't add them together. We have to generate a `load` to "dereference" our pointer.

```go
aVal := builder.CreateLoad(a, "a_val")
bVal := builder.CreateLoad(b, "b_val")
```

And then the arithmetic part. We'll be doing `a + b`, this is straight-forward as we need to create an add instruction:

```go
result := builder.CreateAdd(aVal, bVal, "ab_value")
```

Now we need to return this since our function returns an `i32`. 

```go
builder.CreateRet(result)
```

And that is it! But how do we execute this? We have a few choices, we can either:

* Use LLVM's JIT/execution engine
* Translate into IR -> BitCode -> Assembly -> Object -> Executable

Because the first option is a little more concise to fit into the executable, I'm going for that route. I'll leave it as an exercise to the reader to do the second route. If you create an executable, if you check the status code after running it, it should be the `48` which is the result. To do this in Bash, echo out `$?` environmental variable:

	$ ./a.out
	$ echo $?
	$ 48

If you want to print this to standard out, you will have to define the `printf` function, `putch` or some equivalent. Hopefully this tutorial provides you with enough to do this. If you get stuck (shameless plug), I'm working on a language called Ark that is built on top of LLVM and is written in Go. [You can check out our code generator here](https://github.com/ark-lang/ark/blob/master/src/codegen/LLVMCodegen/codegen.go).

There is documentation availible on the LLVM bindings [here](https://godoc.org/llvm.org/llvm/bindings/go/llvm). This has almost everything you need to know.

As well as the [LLVM specification](http://llvm.org/docs/LangRef.html), which covers everything in detail. This includes instructions, intrinsics, attributes, and everything else.

## Running our code
Enough rambling, lets get to it. Here's an overview of what this section involves:

* Verifying our module
* Initializing the execution engine
* Setting up our function call and executing it!

First let's verify our module is correct. 

```go
if ok := llvm.VerifyModule(mod, llvm.ReturnStatusAction); ok != nil {
	fmt.Println(ok.Error())
	// ideally you would dump and exit, but hey ho
}
mod.Dump()
```

This will print out the error if our module is invalid. An invalid module could be caused by a variety of things, though malformed IR is the most likely cause. The `mod.Dump()` call will dump the module IR to the standard out.

Now to initialize an execution engine:

```go
engine, err := llvm.NewExecutionEngine(mod)
if err != nil {
	fmt.Println(err.Error())
	// exit...
}
```

And finally, running our function and printing the result to stdout. We pass an empty array of GenericValues since our function takes no arguments:

```go
funcResult := engine.RunFunction(mod.NamedFunction("main"), []llvm.GenericValue{})
fmt.Printf("%d\n", funcResult.Int(false))
```

# Building
You need to have LLVM installed. Luckily for me this is as simple as:

```bash
$ pacman -S llvm
```

If you are on Windows, this may be trickier. On any other Linux distribution, search for llvm in your package manager. On Mac you can use Homebrew.

And then we install the go bindings. The release variable is 362, though if you are using say llvm 3.7.0, this should be 370, etc. The below will clone the LLVM repository into the GOPATH, then build and install the bindings.

```bash
$ release=RELEASE_362
$ svn co https://llvm.org/svn/llvm-project/llvm/tags/$release/final $GOPATH/src/llvm.org/llvm
$ cd $GOPATH/src/llvm.org/llvm/bindings/go
$ ./build.sh
$ go install llvm.org/llvm/bindings/go/llvm
```

Now in the go file, make sure you add the relevant imports, e.g. `import "llvm.org/llvm/bindings/go/llvm"`. Once this is done, you can run your go file and it should print out the result:

![](/assets/img/Screenshot-from-2016-03-13-17-42-03.png)

Done! Hopefully you can see how this can be used to build a compiler. The next step from here would be to check out the LLVM Kaleidoscope tutorial, or experiment and try implement your own thing.

## Full Code

```go
package main

import (
	"fmt"
	"llvm.org/llvm/bindings/go/llvm"
)

func main() {
	// setup our builder and module
	builder := llvm.NewBuilder()
	mod := llvm.NewModule("my_module")

	// create our function prologue
	main := llvm.FunctionType(llvm.Int32Type(), []llvm.Type{}, false)
	llvm.AddFunction(mod, "main", main)
	block := llvm.AddBasicBlock(mod.NamedFunction("main"), "entry")
	builder.SetInsertPoint(block, block.FirstInstruction())

	// int a = 32
	a := builder.CreateAlloca(llvm.Int32Type(), "a")
	builder.CreateStore(llvm.ConstInt(llvm.Int32Type(), 32, false), a)

	// int b = 16
	b := builder.CreateAlloca(llvm.Int32Type(), "b")
	builder.CreateStore(llvm.ConstInt(llvm.Int32Type(), 16, false), b)

	// return a + b
	bVal := builder.CreateLoad(b, "b_val")
	aVal := builder.CreateLoad(a, "a_val")
	result := builder.CreateAdd(aVal, bVal, "ab_val")
	builder.CreateRet(result)

	// verify it's all good
	if ok := llvm.VerifyModule(mod, llvm.ReturnStatusAction); ok != nil {
		fmt.Println(ok.Error())
	}
	mod.Dump()

	// create our exe engine
	engine, err := llvm.NewExecutionEngine(mod)
	if err != nil {
		fmt.Println(err.Error())
	}

	// run the function!
	funcResult := engine.RunFunction(mod.NamedFunction("main"), []llvm.GenericValue{})
	fmt.Printf("%d\n", funcResult.Int(false))
}
```