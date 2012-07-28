// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"os"
	"path"
	"text/template"
)

func init() {
	Register(ApiNewProject, apiNewProject)
}

// apiNewProject creates a new project.
func apiNewProject(c *Client, in []byte) {

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
