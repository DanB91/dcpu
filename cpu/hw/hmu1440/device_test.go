// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package hmu1440

import (
	"testing"
)

func Test(t *testing.T) {
	dev := New()
	err := dev.Open("test.fdd")

	if err != nil {
		t.Fatal(err)
	}

	defer dev.Close()
}
