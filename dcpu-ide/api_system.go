// This file was automatically generated.
// Any changes to it will not be preserved.

package main

import (
	"net/http"
)

// Returns the current configuration data.
func apiConfig(r *http.Request) ([]byte, int) {
	return Pack(config), 200
}
