// This file was automatically generated.
// Any changes to it will not be preserved.

package main

import "net/http"

// apiConfig returns the current configuration data.
// We cherry-pick only the parts which are relevant to the client.
//
// This can not be done through struct field tags, because we already have
// some in place for the config.Load/Save routines.
func apiConfig(r *http.Request) ([]byte, int) {
	data := struct {
		DcpuPath string
		Address  string
		Timeout  uint
	}{
		config.DcpuPath,
		config.Address,
		config.Timeout,
	}
	return Pack(data), 200
}

// apiPing serves as a keep-alive pump. This is a no-op.
//
// The StateTracker updates its last request time on each and every
// request to the server. But in some cases, there is nothing to fetch.
// This is where these pings are necessary to keep the server alive.
func apiPing(r *http.Request) ([]byte, int) {
	return nil, 200
}
