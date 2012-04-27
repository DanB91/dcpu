## DCPU

This repository contains [DCPU assembly][1] utilities and code.
Mostly commonly used library bits and bobs, along with a preprocessor
which should make the writing of larger programs a little less painful.

[1]: http://dcpu.com

* _dcpu.lang_: This file is a DCPU syntax file for GtkSourceView
  compatible editors. It should be installed in the GtkSourceView
  language-specs directory. For me this is at:
  `/usr/share/gtksourceview-3.0/language-specs/`. 
* _dcpu-pp_: This is a commandline tool that offers some pre-processing
  magic for `.dasm` source code. Refer to its README for more info.
* _lib/_: This directory holds often used assembly code.

### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

