// This file was automatically generated.
// Any changes to it will not be preserved.

package main

import (
	"net/http"
)

func init() {
	Register("/api/dirlist", apiDirList)
}

// apiDirList returns the contents of a given directory.
func apiDirList(r *http.Request) ([]byte, int) {

	return nil, 200
}
