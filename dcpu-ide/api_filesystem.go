// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import ()

func init() {
	Register(ApiReadFile, apiReadFile)
	Register(ApiDirList, apiDirList)
}

// apiReadFile attempts to read a given file and returns its contents.
func apiReadFile(c *Client, in []byte) {

}

// apiDirList returns the contents of a given directory.
func apiDirList(c *Client, in []byte) {

}
