// This file was automatically generated.
// Any changes to it will not be preserved.

package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	Register("/api/newproject", "POST", apiNewProject)
}

// apiNewProject rcreates a new project.
func apiNewProject(r *http.Request) ([]byte, int) {
	name := r.FormValue("tName")

	// Perform some sanity checks on the new project name/location.
	if len(name) == 0 {
		return Error(ErrMissingName), http.StatusNotAcceptable
	}

	dir := filepath.Join(config.ProjectPath, name)
	stat, err := os.Lstat(dir)

	if err != nil && !os.IsNotExist(err) {
		return Error(ErrUnknown, err.Error()), http.StatusInternalServerError
	}

	if stat != nil {
		return Error(ErrDuplicateProject, name), http.StatusNotAcceptable
	}

	// Create project directory.
	if err = os.MkdirAll(dir, 0700); err != nil {
		return Error(ErrUnknown, err.Error()), http.StatusInternalServerError
	}

	return Pack(struct {
		Path  string
		Name  string
		Files []string
	}{
		dir,
		name,
		nil,
	}), 200
}
