// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	dp "github.com/jteeuwen/dcpu/parser"
)

// Every pre-processor implements this interface.
type PreProcessor interface {
	Process(*dp.AST) error
}

// Common constructor for new pre-processors.
type PreProcessorFunc func() PreProcessor

type PreProcessorDef struct {
	proc PreProcessorFunc
	desc string
	use  bool
}

// List of available pre-processor and a value indicating whether
// we should use them or not.
var preprocessors map[string]*PreProcessorDef

// Register registers a new pre-processor with its commandline name,
// description string.
func RegisterPreProcessor(name, desc string, pf PreProcessorFunc) {
	if preprocessors == nil {
		preprocessors = make(map[string]*PreProcessorDef)
	}

	if _, ok := preprocessors[name]; ok {
		panic("Duplicate PreProcessor: " + name)
	}

	preprocessors[name] = &PreProcessorDef{
		proc: pf,
		desc: desc,
		use:  false,
	}
}

// PreProcess traverses all registered pre-processor and
// passes the AST into them for parsing.
func PreProcess(ast *dp.AST) (err error) {
	for _, v := range preprocessors {
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

// CreatePreProcessorFlags adds commandline flags for all
// registered pre-processor.
func CreatePreProcessorFlags() {
	for k, v := range preprocessors {
		flag.BoolVar(&v.use, k, v.use, v.desc)
	}
}
