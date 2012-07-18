// This file was automatically generated.
// Any changes to it will not be preserved.

package main

import (
	"net/http"
)

// apiConfig returns the current configuration data.
func apiConfig(r *http.Request) ([]byte, int) {
	return Pack(config), 200
}

// apiPing serves as a keep-alive pump.
// This is a no-op.
func apiPing(r *http.Request) ([]byte, int) {
	return nil, 200
}
