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
	proc  PreProcessorFunc // Processor handler.
	desc  string           // Processor description.
	use   bool             // Whether to perform the processing or not.
	isopt bool             // is this an optimization?
}

// List of available pre-processor and a value indicating whether
// we should use them or not.
var preprocessors map[string]*PreProcessorDef

// Register registers a new pre-processor with its commandline name,
// description string.
func RegisterPreProcessor(name, desc string, pf PreProcessorFunc, isopt bool) {
	if preprocessors == nil {
		preprocessors = make(map[string]*PreProcessorDef)
	}

	if _, ok := preprocessors[name]; ok {
		panic("Duplicate PreProcessor: " + name)
	}

	preprocessors[name] = &PreProcessorDef{
		proc:  pf,
		desc:  desc,
		use:   false,
		isopt: isopt,
	}
}

// PreProcess traverses all registered pre-processor and
// passes the AST into them for parsing.
func PreProcess(ast *dp.AST, force_opt bool) (err error) {
	for _, v := range preprocessors {
		if force_opt && v.isopt {
			v.use = true
		}

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
