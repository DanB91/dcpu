// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"encoding/json"
	"github.com/jteeuwen/dcpu/asm"
	"io/ioutil"
)

func writeDebug(d *asm.DebugInfo, file string) (err error) {
	if d == nil || len(file) == 0 {
		return
	}

	data, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		return
	}

	return ioutil.WriteFile(file, data, 0644)
}
