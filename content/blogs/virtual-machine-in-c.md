---
title: Virtual Machine in C
type: blogs
---

# Virtual Machine in C

Here's the GitHub to show what we'll be making. You can also compare your code to this repository in case you have any errors: [GitHub Repository](http://www.github.com/felixangell/mac)

I felt like writing an article about building your very own virtual machine in the C programming language. I love working on lower level applications e.g. compilers, interpreters, parsers, virtual machines, etc. So I thought I'd write this article as learning how virtual machines work is a great way to introduce yourself into the general realm of lower level programming!

## Prerequisites & Notices

There are a few things that you need before we can continue:

- GCC/Clang/.. — I'm using clang, but you can use any modern compiler;
- Text Editor — I would suggest a text editor over an IDE (when writing C), I'll be using Emacs;
- Basic programming knowledge — Just the basics: variables, flow control, functions, structures, etc; and
- GNU Make — A build system so we aren't writing the same commands in the terminal over and over to compile our code

## Why should I write a Virtual Machine?

Here are some reasons why you should write a virtual machine:

- You want a deeper understanding of how computers work. This article will help you understand what your computer does at a lower level, a virtual machine provides a nice simpler layer of abstraction. And there's no better way to learn than build one, ey? 
- You just want to learn about virtual machines because it's fun.
- You want to learn more about how some programming languages work. For instance, various languages nowadays target virtual machines - usually written specifically for the language. Examples include the JVM, Lua's VM, Facebook's Hip-Hop VM (PHP/Hack), etc. There's also quite a large abstraction, from say a C++ program to assembly for your machine. When you think about it, we take a lot for granted when we write our programs in our fancy OOP paradigm language with garbage collection and all these nice features.

### Instruction Set

We'll be implementing our own instruction set, it will be relatively simple. I'll briefly cover instructions like moving values from registers, or jumping to other instructions, but hopefully you'll figure this out after you've read the article.

Our virtual machine will have a set of register: `A`, `B`, `C`, `D`, `E`, and `F`. These are general purpose registers, which means that they can be used for storing anything. This is as opposed to say special purpose registers, for example on x86, e.g: `ip`, `flag`, `ds`, ... 

A program will be a read-only sequence of instructions. The virtual machine is a stack-based virtual machine, which means that it has a stack we can push and pop values to, and a few registers we can use too. These are also a lot more simpler to implement than a register-based virtual machine.

Without further ado, here's an example of the instruction set we're going to be implementing in action. The semi-colons are comments on what the line will do.

```
PSH 5       ; pushes 5 to the stack
PSH 10      ; pushes 10 to the stack
ADD         ; pops two values on top of the stack, adds them pushes to stack
POP         ; pops the value on the stack, will also print it for debugging
SET A 0     ; sets register A to 0
HLT         ; stop the program
```

That's our instruction set, note that the POP instruction will print the instruction we popped, this is more of a debugging thing (ADD will push a result to the stack, so we can POP the value from the stack to verify it is there). 
I've also included a SET instruction, this is so you understand how registers can be accessed and written to. You can also try your hand at implementing instructions like `MOV A, B` (move the value A to B). `HLT` is the instruction to show we've finished executing the program.

## How does a Virtual Machine work?
Virtual Machines are more simple than you think, they follow a simple pattern, the "instruction cycle", which is: fetch; decode; and execute.
First we fetch the next instruction in the instruction list or code, we then decode the instruction and execute the decoded instruction.

## Project Structure
Before we start programming, we need to set-up our project. We need a folder where our project will be located, I like to keep my projects under `~/dev`. Here's how we would set up our project in the terminal. This is assuming you already have a `~/dev/` directory, but you can `cd` into anywhere you want your project to be.

```
$ cd ~/dev/
$ mkdir mac
$ cd mac
$ mkdir src
```

Above we `cd` into our `~/dev` directory, or wherever you want your project to be, we make a directory (I'm calling this VM "mac"). We then `cd` into that directory and make our `src` directory, which is where our code will be located.

## Makefile
Our makefile is relatively straight-forward, we won't need to separate anything into multiple files, and we won't be including anything so we just need to compile the file with some flags:

```makefile
SRC_FILES = main.c
CC_FLAGS = -Wall -Wextra -g -std=c11
CC = clang

all:
    ${CC} ${SRC_FILES} ${CC_FLAGS} -o mac
```

That should suffice for now, you can always improve it later on, but as long as it does the job, we should be fine.

## Program Instructions
Now for the Virtual Machines code. First we need to define the instructions for our program. For this, we can just use an `enum`, since our instructions are basically numbers from 0 - X.
In fact, an assembler program will take your assembly files and (sort of) translate all of the ops into their number counterparts. For example, if you wrote an assembler for mac, it would translate all `MOV` ops into the number `0`, and so on.

```c
typedef enum {
    PSH,
    ADD,
    POP,
    SET,
    HLT
} InstructionSet;
```

Now we can store a test program as an array. So for a test, we'll write a simple program to add the values `5` and `6`, then print them out.

Note: When I say print them out, really I'll just make it so that when we call "pop" our virtual machine will `printf` the value that we pop. In reality you wouldn't want to do this unless you're debugging or something. 

The instructions should be stored as an array, I'll define it somewhere at the top of the document; you could probably throw it in a header file. Here's our test program:

```c
const int program[] = {
    PSH, 5,
    PSH, 6,
    ADD,
    POP,
    HLT
};
```

The above program will push `5` and `6` to the stack, execute the add instruction which will pop the two values that are on the stack, add them together and push the result back on the stack. We then pop the result since our pop instruction will print the value (for debugging purposes). 

Finally, the `HLT` instruction means terminate the program. This is used so that if we had control flow we can terminate the program whenever. Our virtual machine will terminate itself naturally if we had no instructions left, though.

Now we have to implement the instruction cycle (fetch, decode, execute). Technically we don't really have to decode anything. This will make more sense later.

## Fetching the current instruction
Because we have stored our program as an array, it's simple to fetch the current instruction. A virtual machine has a counter, typically called a Program Counter, Instruction Pointer, ... these names are synonymous and your choice is personal preference. Usually they are shortened to PC or IP respectively.

If you remember before, I said that we would store the program counter as a register... we will do that, but later on. For now, we'll just create a variable at the top of our code called `ip`, and set that to `0`:

```c
int ip = 0;
```

This `ip` stands for instruction pointer. The program itself is stored as an array of integers. The `ip` variable serves as an index in the array as to which instruction is currently being executed.

```c
int ip = 0;

int main() {
    int instr = program[ip];
    return 0;
}
```

If we were to printf the `instr` variable, it should give us `PSH` (or `0`). We can write this as a fetch function like so:

```c
int fetch() {
    return program[ip];
}
```

This function will return the current instruction when called. So, what if we want the next instruction? We just increment the instruction pointer:

```c
int main() {
    int x = fetch(); // PSH
    ip++; // increment instruction pointer
    int y = fetch(); // 5
}
```

So how do we automate this? Well we know that a program runs until it is halted via the `HLT` instruction. So we just have an infinite loop that will keep looping until the current instruction is `HLT`.

```c
#include <stdbool.h> 

bool running = true;

int main() {
    while (running) {
       int x = fetch();
       if (x == HLT) running = false;
       ip++;
    }
}
```

This will work perfectly fine, but it's kind of messy. What we're doing is looping through each instruction, checking if the value of that instruction is `HLT`, if it is then stop the loop, otherwise eat the instruction and repeat.

## Evaluating an instruction
So this is the gist of our Virtual Machine, but we can do better. A virtual machine is so simple that you can write a huge switch statement. In fact, this is usually the best way to do it in terms of speed, as opposed to say a HashMap for all the instructions and some abstract class or interface with an `execute` method.

Each case in the switch statement would be an instruction that we defined in our enum. The eval function will take a single parameter, which is the instruction to evaluate. We won't do any of the instruction pointer increments in this function unless we're consuming operands.

```c
void eval(int instr) {
    switch (instr) {
    case HLT:
        running = false;
        break;
    }
}
```

Let's add this back into the main loop of the virtual machine:

```c
bool running = true;
int ip = 0;

// instruction enum
// eval function
// fetch function

int main() {
    while (running) {
        eval(fetch());
        ip++; // increment the ip every iteration
    }
}
```

## The stack!
Great, that should work perfectly. Now before we add the other instructions, we need a stack. The stack is a very simple data structure. We'll be using an array for this rather than a linked list. Because our stack is a fixed size, we don't have to worry about resizing/copying. And it's probably better in terms of cache efficiency to use an array rather than a linked list!

Similarly to how we have an `ip` that indexes the program array, we need a stack pointer (`sp`) to show where we are in the stack array.

Here's a little visualisation of our stack data structure:

```c
[] // empty

PSH 5 // put 5 on **top** of the stack
[5]

PSH 6 // 6 on top of the stack
[5, 6]

POP // pop the 6 off the top
[5]

POP // pop the 5
[] // empty

PSH 6 // push a 6...
[6]

PSH 5 // etc..
[6, 5]
```

Let's break down our program in terms of the stack:

```c
PSH, 5,
PSH, 6,
ADD,
POP,
HLT
```

Well first we push 5 to the stack

```c
[5]
```

Then we push 6:

```c
[5, 6]
```

Then the add instruction will basically pop these values and add them together and push the result on the stack:

```c
[5, 6]

// pop the top value, store it in a variable called a
a = pop; // a contains 6
[5] // stack contents

// pop the top value, store it in a variable called b
b = pop; // b contains 5
[] // stack contents

// now we add b and a. Note we do it backwards, in addition
// this doesn't matter, but in other potential instructions
// for instance divide 5 / 6 is not the same as 6 / 5
result = b + a;
push result // push the result to the stack
[11] // stack contents
```

Where does our stack pointer come into play? Well the stack pointer or `sp` is set to `-1`, this means it's empty. Arrays are zero-indexed in C, so if the sp was `0` it would be set to whatever random number the C compiler throws in there because the memory for an array is not zeroed out.

Now if we push 3 values, the sp would be 2. So here's an array with 3 values:

```
        -> sp -1
    psh -> sp 0
    psh -> sp 1
    psh -> sp 3

  sp points here (sp = 2)
       |
       V
[1, 5, 9]
 0  1  2 <- array indices or "addresses"
```

Now when we **pop** a value from the stack, we decrement the stack pointer, so that means that we're popping 9, and the top of the stack will be 5:

```
    sp points here (sp = 1)
        |
        V
    [1, 5]
     0  1 <- these are the array indices
```

When we want to see the top of the stack, we look at the value at the current `sp`. Okay, hopefully you should know how a stack works! If you're still confused, check out [this](https://en.wikipedia.org/wiki/Stack_(abstract_data_type)) article on wikipedia.

To implement a stack in C is relatively straight-forward. Along with our `ip` variable, we should define the `sp` variable and our array which will be the stack:

```c
int ip = 0;
int sp = -1;
int stack[256];

...
```

Now if we want to push a value to the stack, we increment the stack pointer **then** we set the value at the current sp (which we just incremented).

**The order here is very important!** If you set the value first, _then_ you increment the `sp` you will get some bad behaviour because we're writing to the memory at index `-1`, not good!

```c
// sp = -1
sp++; // sp = 0
stack[sp] = 5; // set value at stack[0] -> 5
// top of stack is now [5]
```

In our `eval` function, we can add the stack push like this:

```c
void eval(int instr) {
    switch (instr) {
        case HLT: {
            running = false;
            break;
        }
        case PSH: {
            sp++;
            stack[sp] = program[++ip];
            break;
        }
    }
}
```

There are a few differences between the previous eval function. Firstly, there are braces around the case clauses. If you aren't familiar with this trick, it gives the case a scope so you can define variables inside of the clause.

Secondly, the `program[++ip]` expression. Why are we doing that here? It's because our `PSH` instruction has an argument. `PSH, 5`. Immediately after the `PSH` op is the value that we want to push to the stack.
Here we increment the ip so that its pointing to the `5`, and then we access that value from the program array. 

    program = [ PSH, 5, PSH, 6, ]
                0    1  2    3

    when pushing:
    ip starts at 0 (PSH)
    ip++, so ip is now 1 (5)
    sp++, allocate some space on the stack
    stack[sp] = program[ip], put the value 5 on the stack

The `POP` instruction is as simple as decrementing the stack pointer. However, I wanted to make it so that the pop instruction will print the value it just popped. Because of this, we have to do a little bit more work:

```c
case POP: {
    // store the value at the stack in val_popped THEN decrement the stack ptr
    int val_popped = stack[sp--];

    // print it out!
    printf("popped %d\n", val_popped);
    break;
}
```

Finally, the `ADD` instruction. This one may be a little trickier to get your head around, and this is where we need our scope trick on the case clause because we're introducing some variables.

```c
case ADD: {
    // first we pop the stack and store it as 'a'
    int a = stack[sp--];

    // then we pop the top of the stack and store it as 'b'
    int b = stack[sp--];

    // we then add the result and push it to the stack
    int result = b + a;
    sp++; // increment stack pointer **before**
    stack[sp] = result; // set the value to the top of the stack

    // all done!
    break;
}
```

Note that the order here is important for certain operations! If you were implementing divide, you might have some troubles because:
    
    5 / 4 != 4 / 5

Stacks are LIFO (last in, first out). Meaning if we pushed 5 then 4, we would pop 4, then pop 5. If we did `pop() / pop()` this would give us the wrong expression, so it's crucial that you get the order correct.

## Registers
Registers are very easy to implement, I mentioned we would have the registers `A`, `B`, `C`, `D`, `E`, and `F`. We can use an enum for this like we did the instruction set.

```c
typedef enum {
   A, B, C, D, E, F,
   NUM_OF_REGISTERS
} Registers;
```

The last member in the enum `NUM_OF_REGISTERS` is a little trick so we can get the size of the registers even if you add more.

We'll store our registers in an array. Because we use an enum, A = 0, B = 1, C = 2, etc. So when we want to set the register A, it's as simple as saying `register[A] = some_value`.

```c
int registers[NUM_OF_REGISTERS];
```

Printing out the value in register `A`:

```c
printf("%d\n", registers[A]); // prints the value at the register A
```

### Instruction Pointer
**What about branching?** I'll leave that to you! Remember an instruction pointer points to the current instruction. Now because this is in the virtual machines source code, your best bet would be to have the instruction pointer as a register that you can read and manipulate from the virtual machines programs.

```c
typedef enum {
    A, B, C, D, E, F, PC, SP,
    NUM_OF_REGISTERS
} Registers;
```

Now we have to port our code to actually use these instruction and stack pointers. A quick and dirty method to do this with the existing code-base is to remove the `sp` and `ip` variables at the top and replace them with a define:

```c
#define sp (registers[SP])
#define ip (registers[IP])
```

That should be a decent fix so that you don't have to re-write a lot of your code, and it should function perfectly. However, this may not scale very well, and it is aliasing some code so I would suggest not using this method, but for a simple toy virtual machine it will suffice.

When it comes to branching in our code, I'll give you a hint. With our new `IP` register, we can branch by writing to this IP a different value. Try this sample below and see what it will do:

```
PSH 10
SET IP 0
```

Similar to the BASIC program that a lot of people are familiar with:

```
10 PRINT "Hello, World"
20 GOTO 10
```

However, since we are pushing values to the stack constantly, we will eventually get a stack overflow once we've pushed more than the amount of space we've defined.

Note that each 'word' is an instruction, so given the following program:

```
              ;  these are the instructions
PSH 10        ;  0 1
PSH 20        ;  2 3
SET IP 0      ;  4 5 6
```

If we wanted to jump to the second `set` instruction, we would set the `IP` register to `2` instead of `0`.

### Fin
And there you have it! If you run `make` in the projects root directory and you can execute the virtual machine: `./mac`.

You can check out the source code on the github [here](http://www.github.com/felixangell/mac). If you want to see a more developed version of the VM with the `MOV` and `SET` instructions then check out the `mac-improved` directory. 
The source code for the virtual machine we implemented in this article is in `mac.c`

### Further Reading
If you're interested in this topic and want to expand more, there is a lot of resources out there on the internet. Markus Persson (Notch) wrote a DCPU-16, which is basically a 16 bit virtual machine for the scrapped game 0x10c. 

There are a few implementations of it around GitHub you can check out. You could also look into emulating something like a simple CPU, e.g. a Zilog Z80. 
If you want to write an emulator for something like this, go check out the manual and see if you can implement the instruction set and the registers. There's a few implementations on GitHub if you need any help.

Thanks for reading!