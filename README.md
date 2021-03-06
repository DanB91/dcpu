## DCPU

This repository contains [DCPU assembly](http://dcpu.com) utilities and code.
Mostly commonly used library bits and bobs, along with a comprehensive
assembler and unit testing framework.


Commandline tools:

* **dcpu-asm**: This is a commandline assembler with a wide range of options,
  including pre- and post-processors.
* **dcpu-test**: This program runs unit tests as defined in the `lib` 
  directory. We use this to verify newly written code does what we
  want it to do.
* **dcpu-prof**: This is an interactive commandline tool that can analyze
  profiling data generated by the `prof` package.
* **dcpu-data**: This is a small tool which generates DCPU assembly source
  from any input file. Useful if you want to embed binary data in your
  programs.
* **dcpu-fmt**: This tool formats DCPU source files according to some
  predefined styling rules.

Packages:

* **parser**: This holds a package that parses assembly source and turns it
  into an Abstract Syntax Tree.
* **parser/util**: This package contains some parser related utility
  bits and bobs.
* **asm**: This holds aan assembler. It turns an AST into a compiled
  program, ready to be passed to the CPU for execution.
* **cpu**: A CPU emulator implementation. It adds the necessary instructions
  to make unit tests behave properly. As such, it may not be ideal to use
  as a standalone emulator.
* **cpu/hw/**: List of hardware components that can be hooked into the CPU.
* **prof**: this package holds a profiler for DASM code. It maintains
  information like cycle costs about a currently executing program.
  An emulator can use it to generate a profile file which can then be examined
  by the `dcpu-prof` tool for performance bottlenecks.

DCPU asssembly:

* **lib/**: This directory holds often used assembly code and unit tests.

Misc stuff:

* **dcpu.lang**: This file is a DCPU syntax file for GtkSourceView
  compatible editors (like Gedit). It should be installed in the
  language-specs directory.
  For me this is at: `/usr/share/gtksourceview-3.0/language-specs/`. 
* **dcpu-run**: This shell script ties `dcpu-test` and `dcpu-prof` together
  to run unit tests on a given input source file. It generates profiling
  data and then displays a focused overview of this data.


Refer to the README of each individual tool for more info.


### Usage

To install all tools in one go, do the following:

    $ go get https://github.com/jteeuwen/dcpu/<pkg|cmd>

The commandline programs will be installed where ever your `$GOBIN` points to.
They are now ready for use.


### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

