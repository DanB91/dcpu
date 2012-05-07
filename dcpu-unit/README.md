## DCPU Unit Tester

This app is a DCPU assembler and runtime. It is specifically tailored
to run assembly unit tests as defined below.

The idea is to slap this somewhere into your existing tool chain.
The tool will look for all `*_test.dasm` files in the given path
and run them. Any code that generates errors causes the tool to stop,
yielding appropriate error context.

These tests can be written using the routines defined in `lib/test/`.
The assertion functions perform various comparisons on input
values and panic when these fail. This uses the custom `PANIC` instruction.
It simply prints a sipplied error string and exits the tool.

### *_test.dasm

These contain the actual test code that should be executed. Just like
the dcpu-pp pre-processor, this tool automatically resolves references to
external labels. In principle, one should use these test files to test
the behaviour of one and only one function.

For example: `lib/string/memchr_test.dasm` runs various tests to
probe the behaviour of the `memchr` function. It pushes in a set of
values through CPU registers, calls `memchr` and then performs the unit test.

Any additional data that is required by the tests, can be defined
at the end of the source file. With the exception of the `test` instruction,
the entire `*.test` file is a valid DASM source file.

### Test functions

Example code for a single test unit may look like this:

	 set a, data
	 set b, 3
	 set c, 0
	 jsr memchr
	 
	 set b, data
	 jsr asserteq

	 exit

	:data
	 dat 1, 2, 3, 4, 5

This defines some inputs, the calls `memchr` and compares the value in the
A register with something we expect it to be. `asserteq` panic if this
is not the case.

If all tests pass successfully, the tool exits cleanly.

### Runtime tracing

The `-t` flag will print runtime trace output for each instruction
as it is executed. This allows fine grained insight into what is happening.
This covers the current PC, opcode and operands, all register contents
and the original source file and line that created this instruction.

Here is an example of trace output for a test program.

	./dcpu-unit -i ../lib/ -t ../lib/
	0000: 0001 001c 001f | 0000 0000 0000 0000 0000 0000 0000 0000 | ffff 0000 0000 | memchr.test:1 | set pc, main
	0007: 0001 0000 0023 | 0000 0000 0000 0000 0000 0000 0000 0000 | ffff 0000 0000 | memchr.test:7 | set a, data
	0008: 0001 0001 0024 | 0002 0000 0000 0000 0000 0000 0000 0000 | ffff 0000 0000 | memchr.test:8 | set b, 3
	0009: 0001 0002 0021 | 0002 0003 0000 0000 0000 0000 0000 0000 | ffff 0000 0000 | memchr.test:9 | set c, 0
	000a: 0000 0001 001f | 0002 0003 0000 0000 0000 0000 0000 0000 | ffff 0000 0000 | memchr.test:10 | jsr memchr
	000e: 0012 0001 0003 | 0002 0003 0000 0000 0000 0000 0000 0000 | fffe 0000 0000 | memchr.dasm:15 | ife 0, c ; num is zero -- No compare needed.
	0010: 0012 0008 0001 | 0002 0003 0000 0000 0000 0000 0000 0000 | fffe 0000 0000 | memchr.dasm:19 | ife [a], b
	0012: 0002 0000 0022 | 0002 0003 0000 0000 0000 0000 0000 0000 | fffe 0000 0000 | memchr.dasm:21 | add a, 1
	0013: 0003 0002 0022 | 0003 0003 0000 0000 0000 0000 0000 0000 | fffe 0000 0000 | memchr.dasm:22 | sub c, 1
	0014: 0014 0002 0021 | 0003 0003 ffff 0000 0000 0000 0000 0000 | fffe ffff 0000 | memchr.dasm:23 | ifg c, 0
	0015: 0001 001c 0031 | 0003 0003 ffff 0000 0000 0000 0000 0000 | fffe ffff 0000 | memchr.dasm:24 | set pc, memchr_loop
	0010: 0012 0008 0001 | 0003 0003 ffff 0000 0000 0000 0000 0000 | fffe ffff 0000 | memchr.dasm:19 | ife [a], b
	0012: 0002 0000 0022 | 0003 0003 ffff 0000 0000 0000 0000 0000 | fffe ffff 0000 | memchr.dasm:21 | add a, 1
	0013: 0003 0002 0022 | 0004 0003 ffff 0000 0000 0000 0000 0000 | fffe 0000 0000 | memchr.dasm:22 | sub c, 1
	0014: 0014 0002 0021 | 0004 0003 fffe 0000 0000 0000 0000 0000 | fffe 0000 0000 | memchr.dasm:23 | ifg c, 0
	0015: 0001 001c 0031 | 0004 0003 fffe 0000 0000 0000 0000 0000 | fffe 0000 0000 | memchr.dasm:24 | set pc, memchr_loop
	0010: 0012 0008 0001 | 0004 0003 fffe 0000 0000 0000 0000 0000 | fffe 0000 0000 | memchr.dasm:19 | ife [a], b
	0011: 0001 001c 0018 | 0004 0003 fffe 0000 0000 0000 0000 0000 | ffff 0000 0000 | memchr.dasm:20 | set pc, pop
	000c: 0000 001e 0000 | 0004 0003 fffe 0000 0000 0000 0000 0000 | ffff 0000 0000 | memchr.test:11 | test


### Clock speed

The `-c N` flag defines the cpu's clock speed in nanoseconds.
Set this to a higher value to slow the CPU down. Combined with `-t`, this
can be a powerful debugging tool.


### Dependencies

    $ go get github.com/jteeuwen/dcpu/asm
    $ go get github.com/jteeuwen/dcpu/cpu
    $ go get github.com/jteeuwen/dcpu/parser

### Usage

    $ go get github.com/jteeuwen/dcpu/dcpu-unit

### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.
