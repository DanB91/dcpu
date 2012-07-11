## DCPU Profiler

This package contains a profiler which tracks usage information for a single
running DCPU program. Its output can be written to a file after execution,'
which can then be examined for performance bottlenecks.

The output file is written in a binary format. It can be used by the `dcpu-prof`
tool to query and analyze profiling data.

Among other things, it lists cpu cycle costs for each and every instruction
that was executed. This is tied to the original source code through the
use of debug symbols.


### Dependencies

* github.com/jteeuwen/dcpu/cpu

### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.
