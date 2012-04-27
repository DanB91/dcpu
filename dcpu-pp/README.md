## DCPU Preprocessor

This app is a DCPU Assembly preprocessor. It parses the given input file(s).

Any references to undefined labels are assumed to be defined
in external files names `<labelname>.dasm`. The 'assembler' looks in a
predefined path to resolve these files. It then simply includes the file
contents into the main source file. This allows us to keep source files
small and manageable.

Everything is minified and stripped of unnecessary bloat and then
spit out as a single chunk of DCPU assembly code. This can be pasted into
any of the numerous emulators and run.

### Usage

    go get github.com/jteeuwen/dcpu/dcpu-pp

### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.
