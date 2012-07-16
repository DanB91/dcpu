## Codegen

This tool is only built and run when we are building dcpu-ide.
It is designed specifically to create Go code from static web app files,
so they can be included in the build of `dcpu-ide`.

The files it processes are first gzip compressed before turned into Go code.

When the `-d` flag is supplied, it does not embed the file data, but instead
implements a function which reads it directly from the original file. This
is useful for development when data files change often. We do not want to
rebuild and restart the entire ide server every time.


### License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

