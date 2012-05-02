// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	dp "github.com/jteeuwen/dcpu/parser"
)

// Every post-processor implements this interface.
type Processor interface {
	Process(*dp.AST) error
}

// Common constructor for new processors.
type ProcessorFunc func() Processor

type ProcessorDef struct {
	proc ProcessorFunc
	desc string
	use  bool
}

// List of available processors and a value indicating whether
// we should use them or not.
var processors map[string]*ProcessorDef

// Register registers a new processor with its commandline name,
// description string.
func Register(name, desc string, pf ProcessorFunc) {
	if processors == nil {
		processors = make(map[string]*ProcessorDef)
	}

	if _, ok := processors[name]; ok {
		panic("Duplicate processor: " + name)
	}

	processors[name] = &ProcessorDef{
		proc: pf,
		desc: desc,
		use:  false,
	}
}

// Process traverses all registered processors and
// passes the AST into them for parsing.
func Process(ast *dp.AST) (err error) {
	for _, v := range processors {
		if !v.use {
			continue
		}

		p := v.proc()
		if err = p.Process(ast); err != nil {
			return
		}
	}

	return
}

// CreateProcessorFlags adds commandline flags for all registered processors.
func CreateProcessorFlags() {
	for k, v := range processors {
		flag.BoolVar(&v.use, k, v.use, v.desc)
	}
}
