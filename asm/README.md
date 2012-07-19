## ASM

This is a DASM assembler. It accepts an AST and turns it into a
compiled program you can run on the CPU.


## Assembly conventions

The assembler adheres to some code conventions.
It is advised to take these into account when you write code for this
assembler.


### Register states

* The A, B and C registers are used to push arguments into a function.
* The A register will hold the return value from a function if applicable.
* When calling a function, the A, B and C registers are considered free
  for all. Their contents are not preserved across calls. So if you are
  using them for anything important, you are required to save/restore them
  manually.
* The registers X, Y, Z, I and J are considered to be protected.
  This means they are guaranteed to be preserved across function calls.
  Any function that uses one of these registers must first push its value
  to the stack and pop it back once done.


### Custom function definitions

If the source code uses our custom function definition syntax, the
assembler will use this to automatically detect if any if the protected
registers is being used in the function. If so, it will inject the
necessary push/pop code for each of these at the beginning and end of
a function.

The following source uses the `I` register. The special `def ...` syntax
enables the assembler to do the boilerplate for us, Also note the use
of the `return` keyword here. It is another special instruction that our
assembler understands and will turn into valid code.

    def main
    	set i, 0
	
    :main_loop
    	ife i, 10
    		return
    	
    	add i, 1
    	set pc, main_loop
    end

When we run this through the assembler, it translates this into the
following code before assembly starts:

    :main
    	set push, i
    	set i, 0
	
    :main_loop
    	ife i, 10
    		set pc, $__main_epilog
    	
    	add i, 1
    	set pc, main_loop
	
    :$__main_epilog
		set i, pop
		set pc, pop


### Constants

We can define constant values using the `equ` instruction.
Constants can have any value and they are simply used to replace
any reference to them with said value at assembly time.

For example:

	equ FOO, 123
	equ BAR, 'A'
	
	set i, FOO
	set j, BAR

Translates into:
	
	set i, 123
	set j, 'A'

Constants can be local to a function. When the assembler parses the source
code, it will first resolve constants defined inside functions and then
the global ones.

For example:

	equ FOO, 123       ; global FOO constant -- valid everywhere.
	
	def main
		equ FOO, 321   ; local FOO constant -- Only valid inside `main`
		set a, FOO
	end

Translates into:
	
	:main
		set a, 321
	:$__main_epilog
		set pc, pop


### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.
