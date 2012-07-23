// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

func init() {
	Register("/api/newproject", "POST", apiNewProject)
}

// apiNewProject creates a new project.
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

	// Create project object.
	proj := &Project{
		Path:            dir,
		Name:            name,
		AuthorName:      config.AuthorName,
		AuthorCopyright: config.AuthorCopyright,
		Files: []string{
			"README.md",
			"main.dasm",
		},
	}

	// Copy over initial code template files.
	err = createProjectTemplate(proj)
	if err != nil {
		return Error(ErrTemplateFailure, err), http.StatusInternalServerError
	}

	return Pack(proj), 200
}

// createProjectTemplate copies all project templates over to the new
// project directory after running them through Go's template engine.
func createProjectTemplate(proj *Project) (err error) {
	for i := range proj.Files {
		if err = copyFile(
			proj,
			path.Join(proj.Path, proj.Files[i]),
			path.Join("/project/template", proj.Files[i]),
		); err != nil {
			return
		}
	}

	return
}

// copyFile parses a single project temlate file and writes it
// to the destination file.
func copyFile(proj *Project, dst, src string) (err error) {
	fout, err := os.Create(dst)

	if err != nil {
		return
	}

	defer fout.Close()

	file, ok := static[src]
	if !ok {
		return fmt.Errorf("%q not could not be found.", src)
	}

	t, err := template.New("page").Parse(string(file.Data()))

	if err != nil {
		return
	}

	return t.Execute(fout, proj)
}
