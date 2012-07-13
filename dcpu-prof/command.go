// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/jteeuwen/dcpu/prof"
	"os"
	"regexp"
	"strings"
)

const DefaultMode = "func"

func Handle(prof *prof.Profile, str []string) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Usage = func() {}
	filemode := fs.Bool("file", false, "")

	switch strings.ToLower(str[0]) {
	case "help":
		usage()

	case "top":
		count := fs.Uint("n", DefaultTopCount, "")
		sort := fs.String("s", DefaultTopSort, "")
		fs.Parse(str[1:])

		top(prof, *filemode, *count, *sort)

	case "list":
		filter := fs.String("f", DefaultListFilter, "")
		fs.Parse(str[1:])

		reg, err := regexp.Compile(*filter)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid filter %q.\n", *filter)
			return
		}

		list(prof, *filemode, reg)
	}
}
