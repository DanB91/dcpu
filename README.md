## DCPU

This repository contains [DCPU assembly][1] utilities and code.
Mostly commonly used library bits and bobs, along with a preprocessor
which should make the writing of larger programs a little less painful.

[1]: http://dcpu.com

* _dcpu.lang_: This file is a DCPU syntax file for GtkSourceView
  compatible editors (like Gedit). It should be installed in the
  language-specs directory.
  For me this is at: `/usr/share/gtksourceview-3.0/language-specs/`. 
* _parser_: This holds a package that parses assembly source and turns it
  into an Abstract Syntax Tree.
* _asm_: This holds aan assembler. It turns an AST into a compiled
  program, ready to be passed to the CPU for execution.
* _cpu_: A CPU emulator implementation. It adds the necessary instructions
  to make unit tests behave properly. As such, it may not be ideal to use
  as a standalone emulator.
* _cpu/hw/*_: List of hardware components that can be hooked into the CPU.
* _dcpu-pp_: This is a commandline tool that offers some pre-processing
  magic for `.dasm` source code. Refer to its README for more info.
* _dcpu-unit_: This program contains a custom assembler and emulator
  for the DCPU platform. It runs unit tests as defined in the `lib` 
  directory. We use this to verify newly written code does what we
  want it to do. Refer to its README for more info.
* _lib/_: This directory holds often used assembly code and unit tests.


### Usage

To install all tools in one go, do the following:

    $ git clone https://github.com/jteeuwen/dcpu.git
    $ cd dcpu
    $ go install -a -ldflags="-s" ./...

The `dcpu-pp` and `dcpu-unit` programs will not be installed whereever
your $GOBIN is set to and they are ready for use.

The `-a` switch ensures everything is freshly built (including the linked
Go core packages). The `-ldflags="-s"` switch causes binaries to be built
without all the debug symbols. This speeds them up a little and drastically
reduces their file size. The `./...` bit simply means: Build any and all
packages in this directory and all sub directories. The Go tool will
automatically figure out the dependency tree and build things in the right
order.

### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

