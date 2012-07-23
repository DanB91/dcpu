// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

// Project represents a single DCPU project.
type Project struct {
	Files           []string
	Name            string
	Path            string
	AuthorName      string
	AuthorCopyright string
}
