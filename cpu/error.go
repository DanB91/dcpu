// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package cpu

// The OverloadError occurs when we exceed the capacity
// of the interrupt queue. The CPU should catch fire here and
// do all manner of strange things to your memory/registers.
//
// For the sake of testing, we simply return an error.
type OverloadError struct{}

func (e OverloadError) Error() string {
	return "Interrupt overload: System on fire!"
}

// TestError occurs when a PANIC instruction fires.
// It has a string message along with some current execution state.
type TestError struct {
	Msg string
	PC  Word
}

func (e TestError) Error() string { return e.Msg }
