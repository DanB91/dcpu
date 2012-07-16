// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// This tool queries and analyses profiling data from a given input file.
package main

import (
	"flag"
	"fmt"
	"github.com/jteeuwen/dcpu/prof"
	"os"
	"path/filepath"
)

var (
	topcmd   = flag.Bool("top", false, "")
	listcmd  = flag.Bool("list", false, "")
	count    = flag.Uint("n", DefaultTopCount, "")
	sort     = flag.String("s", DefaultTopSort, "")
	filter   = flag.String("f", DefaultListFilter, "")
	filemode = flag.Bool("file", false, "")
)

func main() {
	prof := parseArgs()

	// Run program once for a specific command,
	// or go into interactive mode?
	if *topcmd || *listcmd {
		if *topcmd {
			Handle(prof, []string{
				"top",
				"-s", *sort,
				"-n", fmt.Sprintf("%d", *count),
				fmt.Sprintf("-file=%v", *filemode),
			})
		} else {
			Handle(prof, []string{
				"list",
				"-f", *filter,
				fmt.Sprintf("-file=%v", *filemode),
			})
		}

		return
	}

	// Interactive mode.
	input := pollInput()

	for {
		select {
		case cmd := <-input:
			Handle(prof, cmd)
		}
	}
}

func parseArgs() *prof.Profile {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <file>\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	help := flag.Bool("h", false, "")
	version := flag.Bool("v", false, "")

	flag.Usage = func() {
		fmt.Printf("Usage %s [options] <file>\n\n", os.Args[0])
		usage()
	}

	flag.Parse()

	if *version {
		fmt.Printf("%s\n", Version())
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "No input file.")
		os.Exit(1)
	}

	fd, err := os.Open(filepath.Clean(flag.Arg(0)))

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defer fd.Close()

	p, err := prof.Read(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read: %v\n", err)
		os.Exit(1)
	}

	return p
}

func usage() {
	fmt.Println(`List of known commands:

 -top [-n -s -file]
   List the top N number of samples for all function calls.
   
       -n : The number of results to limit the output to.
       -s : The sort value denotes the field by which the table should be
            sorted. Possible values are:

            count: This sorts by number of times each entry has been called.
             cost: This sorts by the total cycle cost over the entire program's
                   runtime. This is the default sorting mode.
    -file : Display usage stats per file instead of functions.

 -list [-f -file]
   This gives an instruction-by-instruction listing of cpu cycle usage for
   all entries that match the given filter.

       -f : This is expected to be a regular expression pattern which will be
            matched against labels or file names. It defaults to
            'match everything'. Note that for a large codebase, this can
            generate a large amount of output.

            For best results, use the list command in conjunction with 'top' to
            tell you what code needs closer examination.
    -file : Display usage stats per file instead of functions.
`)
}
