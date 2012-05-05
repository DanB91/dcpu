## NE SPC2000

This package implements the SPC2000 (Suspension Chamber 2000)
The SPC2000 is a deep sleep cell based on the ZEF882 time dilation field
generator (available from Polytron Corporation Incorporated).

The implementation requires game mechanics that this emulator
does not supply. Notably access to sensors that can detect certain
physical world properties which determine if this chamber can be
triggered or not. As such, we just supply random data on interrupt
requests.

### Usage

    go get github.com/jteeuwen/dcpu/hw/spc2000

### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

