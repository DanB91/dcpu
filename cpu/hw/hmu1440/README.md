## HMU1440

This package implements the 1.44 MB 3.5" Harold Media Unit.
This is a floppy that can be plugged into the HMD2043 drive.

It is backed by a file, so we can have its data persisted across
different sessions. It strictly reads/writes data in whole sectors.

Creating a blank backing file is easy enough:

    $ dd if=/dev/zero of=myfile.fdd ibs=1024 count=1440

### Usage

    go get github.com/jteeuwen/dcpu/hw/hmu1440

### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

