## DCPU

This repository contains [DCPU assembly][1] utilities and code.
Mostly commonly used library bits and bobs, along with a comprehensive
assembler and unit testing framework.

[1]: http://dcpu.com

* **dcpu.lang**: This file is a DCPU syntax file for GtkSourceView
  compatible editors (like Gedit). It should be installed in the
  language-specs directory.
  For me this is at: `/usr/share/gtksourceview-3.0/language-specs/`. 
* **parser**: This holds a package that parses assembly source and turns it
  into an Abstract Syntax Tree.
* **asm**: This holds aan assembler. It turns an AST into a compiled
  program, ready to be passed to the CPU for execution.
* **cpu**: A CPU emulator implementation. It adds the necessary instructions
  to make unit tests behave properly. As such, it may not be ideal to use
  as a standalone emulator.
* **cpu/hw/**_: List of hardware components that can be hooked into the CPU.
* **dcpu-asm**: This is a commandline assembler with a wide range of options,
  including pre- and post-processors. Refer to its README for more info.
* **dcpu-test**: This program runs unit tests as defined in the `lib` 
  directory. We use this to verify newly written code does what we
  want it to do. Refer to its README for more info.
* **dcpu-data**: This is a small tool which generates DCPU assembly source
  from any input file. Useful if oyu want to embed binary data in your
  programs.


### Usage

To install all tools in one go, do the following:

    $ git clone https://github.com/jteeuwen/dcpu.git
    $ cd dcpu
    $ go install -ldflags "-X main.AppVersionRev `date -u +%s` -s" ./...

The `dcpu-asm`, `dcpu-test` and `dcpu-data` programs will be installed
where ever your `$GOBIN` points to. They are now ready for use.

The `-X main.AppVersionRev ...` bit in the `-ldflags` section of
`go install` automatically sets the `AppVersionRev` variable in version.go
to the current unix timestamp. This allows us to do auto-incremented
versioning every time we build.

### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

