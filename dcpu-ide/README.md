## DCPU IDE

Some warnings ahead of time:

* This is work in progress and thus subject to lots of breaking changes.
* The release build will only work on unix systems. This should be fixed
  at some point, but I have other priorities now. The debug build should be
  ok on any Go-capable platform, since it involves only Go tools.
* The important parts of this application are tested in both Firefox 13+ and
  Chrome 20+, but Chrome remains my primary target browser. Support for any
  other browsers /may/ come at some point, but don't hold your breath.

This is a webbrowser based development environment for DCPU assembly
projects. It consists of a Go backend and an HTML/Javascript frontend.


### Backend

This is a Go web server which runs on localhost and serves to mediate
commands from the frontend to the dcpu toolchain.

The server listens on `[::1]:7070` by default. This can be changed through the
`-a` commandline flag, altering the `.dcpu-ide` configuration file, or by
setting the `DCPU_IDE_ADDRESS` environment variable. The commandline flags
always takes precedence.


### Frontend

This is an interactive web application that runs in your browser.
It is configured to connect to the backend server running on localhost.
It uses this server to transparently call the DCPU tools in this repository.
 

### Usage

This application does some automated code generation on build.
Which means we require that the Makefile is used to build it.
Not doing so, will cause incorrect builds.

Invoking the IDE is simple:

	$ dcpu-ide

If you want to force it to use a specific browser, you can set the
`BROWSER` environment variable before calling the program:

	$ BROWSER=chromium dcpu-ide

Chromium has a special 'App mode' in which it can run a website.
This removes the standard url bar and other window decorations.
When chromium is supplied as the target browser, this program will
automatically launch in this app mode.


### Debug build

Building in debug mode is simple and has none of the external
tool requirements that are listed in the section for release builds.

    $ make 


### Release build

Building in release mode:

    $ make install

This mode ensures that all the static web app content is first
minified/compressed and then embedded in the server application. This allows
us to move the `dcpu-ide` binary anywhere we want, without having to worry
about dependencies on external files.

During this build, the Makefile invokes `compress_data.sh`.
This script goes through the static web data and compresses/minifies
everything it can (html, js, css, png, etc). For this purpose
it requires some external tools to be present. If these tools do
not exist on your system, build behaviour is undefined.

The external tools are:

* [htmlcompressor](https://aur.archlinux.org/packages.php?ID=48832)
* [yuicompressor](https://aur.archlinux.org/packages.php?ID=22058)
* [pngcrush](https://aur.archlinux.org/packages.php?ID=22877)


### License

DCPU, 0x10c and related materials are Copyright 2012 Mojang.

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

