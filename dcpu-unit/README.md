## DCPU Unit Tester

This app is a DCPU assembler and runtime.  It is specifically tailored
to run assembly unit tests as defined below.

The idea is to slap this somewhere into your existing tool chain.
The tool will perform all available unit tests and compares the
CPU's register states with some predefined correct values for each
test. If there is a mismatch, it will cough up an appropriate
warning message which points you to the relevant source files.

The way in which unit tests and expected outputs are defined,
borrows from The [TECS][1] tool chain, which uses a very similar
approach to unit test the behaviour of hardware components.

[1]: http://www1.idc.ac.il/tecs/

As can be seen in this repo's `lib` directory, unit tests can
be defined for as many components as you need them. Each test comes
with at least two files which are described below.

### *.test

These contain the actual test code that should be executed. Just like
the dcpu-pp pre-processor, this tool automatically resolves references to
external labels. In principle, one should use these test files to test
the behaviour of one and only one function.

For example: `lib/string/memchr.test` runs various tests to
probe the behaviour of the `memchr` function. It pushes in a set of
values through CPU registers, calls `memchr` and then performs the unit test.

Any additional data that is required by the tests, can be defined
at the end of the source file. With the exception of the `test` instruction,
the entire `*.test` file is a valid DASM source file.

### *.cmp

These are 'compare files'. They define the output that the unit tester
will generate from all tests as it /should/ be; one row per test.
So these are the values we expect the code to generate /if/ our implementation
is without bugs.

To make these comparisons meaningful, it is usually a good idea to write
the unit tests for common, as well as for edge cases.

When a function receives input it should not normally be receiving,
strange things may happen. These tests allow us to discover just how
strange these things may be and this allow us to harden our code and
avoid difficult to debug runtime errors.


### Test instruction

Example code for a single test unit may look like this:

	set a, data
	set b, 3
	set c, 5
	jsr memchr
	test
    
    :data dat 1, 2, 3, 4, 5

This defines some inputs, the calls `memchr` and issues the special `test`
instruction. It is this last instruction which denotes the end of a single
unit. We can have arbitrarily many units in a single test file.

It should be noted that each unit is considered a single, full program.
When a new test starts (`test` has just been executed), all state information
in the CPU is reset to its defaults. This includes registers and memory.
This ensures consistent behaviour.

When the `test` instruction is fired, the CPU performs the following steps:

* Fetch the current register states.
* Write a line into the output log, containing said states.
* Fetch the appropriate output line from the `*.cmp` file that accompanies
  this particular test.
* Compare the generated output with the one defined in the compare file.
* If these lines are identical, the CPU is reset and the next unit is run.
* If these lines are not identical, an appropriate error message is created
  and displayed. At which point, all unit testing stops and the tool exits.

If all tests pass successfully, the tool exits cleanly.

### Dependencies

    $ go get github.com/jteeuwen/dcpu/parser

### Usage

    $ go get github.com/jteeuwen/dcpu/dcpu-unit

### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.
