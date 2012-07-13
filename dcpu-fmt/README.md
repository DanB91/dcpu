## DCPU Format

This tool formats DCPU source files according to some predefined styling rules.
it is inspired by Go's fmt tool.

It ensures all sources in a codebase look consistently the same, regardless
of who wrote them. This makes code far easier to read through. Specially if
multiple people are working on the same project. It also puts a desisive
stop to the age old style wars which have never served any purpose but to
waste everyone's time.

The style rules employed by this tool have been chosen by me as I feel they
are best. Not everyone will agree with this; what is important is not
what choices it makes, but the fact that it does. Consequently leaving the
programmer free to worry about more important stuff.

The program can be used to format a single file, or a directory tree at once.
Output is sent to sout by default, but with the `-w` flag, can be sent into
the original file; overwriting its contents with the newly formatted code.

The reformatting is achieved by parsing the input sources into an AST and
regenerating the source from there.


### Example

The following input source is messy and needs serious reworking:

	; device_index device_detect( device_id )
	;
	; This finds a specific hardware device index
	; based on the device ID you specify in registers A and B.
	;
	; It returns the device index in A.
	; A will be -1 if the device was not found.

	:device_detect
	set i,  a; line comment
	SET j,  B
	hwn z
	sub Z,  1
	; Some comment on block
	:device_detect_loop
	HWQ Z
	ife a,  i
	IFE b,  j
	set pc, device_detect_ret
	sub z,  1
	ifa z   ,0xffff
	set pc, device_detect_loop
	:device_detect_ret
	set a,  z
	set PC, POP

Passing it through `dcpu-fmt`, yields the following:

	; device_index device_detect( device_id )
	;
	; This finds a specific hardware device index
	; based on the device ID you specify in registers A and B.
	;
	; It returns the device index in A.
	; A will be -1 if the device was not found.
	:device_detect
	   set i, a ; line comment
	   set j, b
	   hwn z
	   sub z, 1

	; Some comment on block
	:device_detect_loop
	   hwq z
	   ife a, i
		  ife b, j
		     set pc, device_detect_ret

	   sub z, 1
	   ifa z, 0xffff
		  set pc, device_detect_loop

	:device_detect_ret
	   set a, z
	   set pc, pop


### Usage

Format the given source file and print to stdout:

    $ dcpu-fmt data.dasm

Format the given source file and overwrite it:

    $ dcpu-fmt -w data.dasm

Pipe source code into stdin. Reformat it and pipe it into `new.dasm`:

    $ cat data.dasm | dcpu_fmt > new.dasm

Reformat all source files in the current directory and
sub-directoriesrecursively:

    $ dcpu-fmt -w .

Run the program with the `-h` flag for a listing of all options.


### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.
