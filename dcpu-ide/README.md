## DCPU IDE

**Note**: This is work in progress.

This is a webbrowser based development environment for DCPU assembly
projects. It consists of a Go backend and an HTML/Javascript frontend.


### Backend

This is a Go web server which runs on localhost and serves to mediate
commands from the frontend to the dcpu toolchain.

The server listens on [::1]:7070 by default. This can be changed through the
`-a` commandline flag, or by setting the `DCPU_IDE_ADDRESS` environment
variable. If both of these are specified, the commandline flag takes
precedence.


### Frontend

This is an interactive web application that runs in your browser.
It is configured to connect to the backend server running on localhost.
It uses this server to transparently call the DCPU tools in this repository.


### Usage

This application does some automated code generation on build.
Which means we require that the Makefile is used to build it.
Not doing so, will cause incorrect builds.

The code generation depends on the external tool `go-bindata`.

    $ make

To run the ide, invoke the `dcpu-ide` program. It starts the server
and automatically launches a browser, pointing it to the correct page.


### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

