## DCPU Assembler 

This is a commandline assembler for DCPU code.
Besides building DCPU binaries, it can perform pre- and post-processing on
the generated AST and binary code. It offers a set of pre- and post-processors
you can select through command line flags to perform various operations on
the generated source- and binary code.

The output modes are as follows:

* **Compiled binary code**, ready to run on an emulator, is written to
  stdout or the target file. It optionally generates a separate
  debug symbol file which contains debug symbols in a JSON data structure. 
* **Source code**; This simply dumps the fully parsed source back to
  stdout or the target file. This is useful if you only want to pack all
  separate sourcefiles into a single one, perform the necessary compile-time
  error checking and perhaps have it pre-processed in some way. The output
  of this mode is still 100% valid DCPU assembly code and can be pasted into
  any of the online emulators.
* **Abstract Syntax Tree**: This covers all the same options as the *source code*
  mode, but instead it writes a human-readable form of the parsed AST.
  This is mostly useful for testing on my part.


### Usage

The assembler expects a single input file which should be the entry point to
your program. From this file, it will automatically resolve any unknown label
references and treat them as being in external files. It will look for such
a label in a `$labelname.dasm` file somewhere in your include paths. If found,
it automatically loads and parses this file into the AST as well.

If the file can not be found, it yields a compile error with appropriate error
context (source file, line and column info where the label is referenced).


### Error reporting

The various stages this tool works in can all yield error reports
related to whatever they are doing. These errors all contain appropriate
context to identify where and why the error occurred.

* **Parser/AST builder**: This yields syntax errors and undefined label
  reference errors.
* **Pre-processor**: This yields errors specific to each selected pre-processor.
* **Assembler**: This one will tell you about invalid/unknown instructions.
* **Post-processor**: This yields errors specific to each selected post-processor.


### Examples

Compile the given source and write to the `foo.bin` file.
Additionally, write debug symbols to `foo.dbg`.

    $ dcpu-asm -d foo.dbg -o foo.bin foo.dasm

Parse the source into an AST and pipe it to another tool for some more
processing:

    $ dcpu-asm -a foo.dasm | other_tool

Parse the given source file; strip it of comments and whitespace and obfuscate
the label names. Then dump the source back to stdout.

    $ dcpu-asm -strip -scramble -s foo.dasm

Pipe source filename into program. Dump its AST:

    $ echo "../foo.dasm" | dcpu-asm -a

Invoke the program with the `-h` flag to see a listing of all options
along with the available pre- and post-processors.


### Debug symbol files

When creating debug symbol files, the assembler outputs a file with JSON
encoded data. It has the following format:

	{
	 "Files": [
	  "../lib/bootstrap/device_detect.dasm"
	 ],
	 "SourceMapping": [
	  {
	   "File": 0,
	   "Line": 27,
	   "Col": 3
	  },
	  {
	   "File": 0,
	   "Line": 28,
	   "Col": 3
	  },
	  {
	   ...
	  }
	 ]
	}

The `Files` array lists all source file paths which were processed.
The order of this list is important as the `File` field in the entries of
`SourceMapping` are indices into this list.

The `SourceMapping` array contains a list of File, Line and Column numbers
where each generated instruction was originally found. There is one entry
for each instruction in the final program. So, in order to get source context
for the currently executing instruction, we can simply index this list using
the current PC (Program Counter) value (provided the program was loaded into
RAM at offset 0. Otherwise it becomes `info := dbg.SourceMapping[PC-offset]`.

The available debug data may be expanded at some point to include more data.

To see this being used, refer to the `dcpu-test` program.


### Dependencies

* github.com/jteeuwen/dcpu/asm
* github.com/jteeuwen/dcpu/cpu
* github.com/jteeuwen/dcpu/parser


### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.
