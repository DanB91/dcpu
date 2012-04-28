## DCPU Pre-processor

This app is a DCPU Assembly pre-processor. It parses the given input file(s).
When run in a project dir, it considers all `.dasm` files in that
directory to be part of the same project.

### External references

Any references to undefined labels are assumed to be defined
in external files names `<labelname>.dasm`. The 'assembler' looks in a
predefined path to resolve these files. It then simply includes the file
contents into the output source file. This allows us to keep source files
small and manageable. The filenames are expected to be all lower-case.

All this presupposes that each included file contains at least one label
with the same name as the reference. The pre-processor checks the
included code to make sure this is the case.

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

For example, the `-scramble` switch below controls the Scramble processor.

    $ dcpu-pp -h 
    Usage of ./dcpu-pp:
      -a=false: Dump the source code parse tree to stdout.
      -h=false: Display this help.
      -i="": Colon-separated list of additional include paths.
      -o="": Name of destination file. Defaults to stdout.
      -scramble=false: Obfuscate label names and label references.
      -v=false: Display version information.

### Usage

    go get github.com/jteeuwen/dcpu/dcpu-pp

### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.
