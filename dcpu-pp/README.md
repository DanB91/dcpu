## DCPU Pre-processor

This app is a DCPU Assembly pre-processor. It parses the given input file(s) and
spits out either a single assembly source file or an AST. The output can
optionally be transformed by a number of predefined processors which each
perform a specific function.

### External references

Any references to undefined labels are assumed to be defined
in external files named `$labelname.dasm`. The tool looks in predefined paths
to resolve these files. It then simply includes the file contents into the
output source file. This allows us to keep source files small and manageable.
The filenames are expected to be all lower-case.

For example, the following subroutine call refers to label `println`, which
is not defined in the current source file.

    ...
    set a, 1
    set b, 2
    jsr println
    ...

This tool will traverse the list of include path directories and tries to
find `$PATH/println.dasm`. If it can't find the file anywhere, the tool
throws an appropriate parse error with line/column information:

    test.dasm:40:7 Undefined reference: "println"

All this presupposes that `$PATH/println.dasm` contains at least one label
named `println`. The pre-processor checks the included code to make sure this
is the case. if there is no such label defined, it also means the reference
could not be resolved and the same parse error is thrown.

### pre-processor modes

The application operates on a list of pre-processor types which register
themselves at startup. They can be included in the parsing session through
their own commandline switches.

Invoke the program with the `-h` flag to see a list of options.
All boolean switches with a name `> 1` character can be considered
one of these registered processors and they each perform a specific
transformation on the complete AST. By supplying their respective switch
in the commandline invocation, you activate them. By default, they are
all disabled.

For example, the `-scramble` and `-strip` switches below control two processors.

    $ dcpu-pp -h 
    Usage: dcpu-pp [options] file
      -a=false: Dump the source code parse tree to stdout.
      -h=false: Display this help.
      -i="": Colon-separated list of additional include paths.
      -o="": Name of destination file. Defaults to stdout.
      -scramble=false: Obfuscate label names and label references.
      -strip=false: Remove all code comments.
      -v=false: Display version information.v=false: Display version information.

### Dependencies

    $ go get github.com/jteeuwen/dcpu/parser

### Usage

    $ go get github.com/jteeuwen/dcpu/dcpu-pp

Invocation for code in a single project directory is as follows.

    $ dcpu-pp main.dasm

To build the same code using the scramble and strip processors:

    $ dcpu-pp -scramble -strip main.dasm

### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.
