## Generic Keyboard

This package implements a simple, generic keyboard.

It is not backed by a real keyboard. This may be done using
something like Termbox, SDL, GLFW, etc. To make this work,
the Device.poll() method should be implemented to track input
events, map them to Keys as defined in the keyboard spec and
add them to the buffer.

### Usage

    go get github.com/jteeuwen/dcpu/hw/keyboard

### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

