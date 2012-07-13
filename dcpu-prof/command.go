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

func usage() {
	fmt.Println(`List of known commands:

 top [options]
   List the top N number of samples for all function calls where.
   
       -n : The number of results to limit the output to.
       -s : The sort value denotes the field by which the table should be
            sorted. Possible values are:

            count: This sorts by number of times each entry has been called.
             cost: This sorts by the total cycle cost over the entire program's
                   runtime. This is the default sorting mode.
    -file : Display usage stats per file instead of functions.

 list [options]
   This gives an instruction-by-instruction listing of cpu cycle usage for
   all entries that match the given filter.

       -f : This is expected to be a regular expression pattern which will be
            matched against labels or file names. It defaults to 'match everything'.
            Note that for a large codebase, this can generate a large amount of
            output.

            For best results, use the list command in conjunction with 'top' to
            tell you what code needs closer examination.
    -file : Display usage stats per file instead of functions.
`)
}
