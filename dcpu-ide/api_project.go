// This file was automatically generated.
// Any changes to it will not be preserved.

package main

import (
	"net/http"
)

func init() {
	Register("/api/newproject", apiNewProject)
}

// apiNewProject rcreates a new project.
func apiNewProject(r *http.Request) ([]byte, int) {
	return nil, 200
}
