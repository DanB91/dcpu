// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"github.com/jteeuwen/dcpu/asm"
	"github.com/jteeuwen/dcpu/cpu"
)

// Every post-processor implements this interface.
type PostProcessor interface {
	Process([]cpu.Word, *asm.DebugInfo) error
}

// Common constructor for new post-processors.
type PostProcessorFunc func() PostProcessor

type PostProcessorDef struct {
	proc PostProcessorFunc
	desc string
	use  bool
}

// List of available PostProcessors and a value indicating whether
// we should use them or not.
var postprocessors map[string]*PostProcessorDef

// RegisterPostProcessor registers a new PostProcessor with
// its commandline name, description string.
func RegisterPostProcessor(name, desc string, pf PostProcessorFunc) {
	if postprocessors == nil {
		postprocessors = make(map[string]*PostProcessorDef)
	}

	if _, ok := postprocessors[name]; ok {
		panic("Duplicate post processor: " + name)
	}

	postprocessors[name] = &PostProcessorDef{
		proc: pf,
		desc: desc,
		use:  false,
	}
}

// PostProcess traverses all registered post processors and
// passes the compiled program into them for parsing.
func PostProcess(program []cpu.Word, symbols *asm.DebugInfo) (err error) {
	for _, v := range postprocessors {
		if !v.use {
			continue
		}

		p := v.proc()
		if err = p.Process(program, symbols); err != nil {
			return
		}
	}

	return
}

// CreatePostProcessorFlags adds commandline flags for all registered PostProcessors.
func CreatePostProcessorFlags() {
	for k, v := range postprocessors {
		flag.BoolVar(&v.use, k, v.use, v.desc)
	}
}
