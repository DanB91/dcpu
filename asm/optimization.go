// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package asm

import dp "github.com/jteeuwen/dcpu/parser"

type OptimizationFunc func(*dp.AST)

// List of registered optimization handlers.
var optimizations []OptimizationFunc

// RegisterOptimization registers a new optimization handler.
func RegisterOptimization(f OptimizationFunc) {
	optimizations = append(optimizations, f)
}

// optimize goes through a list of registered optimization
// handlers and gives them the ast to operate on.
func optimize(ast *dp.AST) {
	for i := range optimizations {
		optimizations[i](ast)
	}
}
