// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"github.com/jteeuwen/dcpu/prof"
	"io"
	"strconv"
	"strings"
)

func Handle(prof *prof.Profile, str []string) (err error) {
	switch strings.ToLower(str[0]) {
	case "q":
		return io.EOF

	case "help":
		usage()

	case "top":
		count := DefaultTopCount
		sort := DefaultTopSort

		if len(str) > 1 {
			n, err := strconv.Atoi(str[1])
			if err == nil && n > 0 {
				count = n
			}
		}

		if len(str) > 2 {
			sort = str[2]
		}

		top(prof, count, sort)
	}

	return
}

func usage() {
	fmt.Println(`List of known commands:
           help : Display this help.
              q : Quit the application.
 top [N [SORT]] : List the top N number of samples for all function calls.
                  N defaults to 10. The optional SORT value denotes the field
                  by which the table should be sorted. Possible values are:

                  count
                    This sorts by number of times each function has been
                    called. This is the default sorting mode.

                  cumulative
                    This sorts by the total cycle cost over the entire
                    program's runtime.
`)
}
