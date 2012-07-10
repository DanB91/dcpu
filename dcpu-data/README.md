## DCPU Data

This is a small tool which generates DCPU assembly source from any given
input file.s It constructs a listing of `dat` instructions which encode
all the information. Since DCPU deals with 16-bit words, we pack 2 bytes
into every output word.

### Usage

Turn the given text file into code and write it to `data.dasm`:

    $ dcpu-data -o data.dasm data.txt

Turn the given text file into code and write to stdout to process
it further in another tool:

    $ dcpu-data data.txt | another_tool

Pipe file data into stdin. Label it 'mytext' and append to
another source file:

    $ cat data.txt | dcpu_data -l mytext >> data.dasm

Run the program with the `-h` flag for a listing of all options.

### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.
