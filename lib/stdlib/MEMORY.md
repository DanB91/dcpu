## Memory layout

The memory layout we will be employing in the standard library
is as follows:

	0x0000-0x0FFF: kernel
	0x1000-0x2FFF: OS
	0x3000-0x5FFF: userspace programs 
	0x6000-0x7FFF: memory 
	0x8000-????  : I/O 
	????-0xFFFF  : stack

This is mainly relevant for `malloc` and friends as it gives us some concrete
memory space to work with when allocating. It should be noted that this is
the default behaviour and it can easily be changed if you need it to.

