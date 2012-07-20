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
		return Errorf("Missing project name."),
			http.StatusNotAcceptable
	}

	dir := filepath.Join(config.ProjectPath, name)
	stat, err := os.Lstat(dir)

	if err != nil && !os.IsNotExist(err) {
		return Errorf("Unknown error: %v", err),
			http.StatusInternalServerError
	}

	if stat != nil {
		return Errorf("Project %q already exists", name),
			http.StatusNotAcceptable
	}

	// Create project directory.
	if err = os.MkdirAll(dir, 0700); err != nil {
		return Errorf("Unknown error: %v", err),
			http.StatusInternalServerError
	}

	return nil, 200
}
