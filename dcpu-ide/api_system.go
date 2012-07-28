// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

func init() {
	Register(ApiConfig, apiConfig)
}

// apiConfig returns the current configuration data.
func apiConfig(c *Client, in []byte) {

}
