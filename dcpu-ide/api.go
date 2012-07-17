// This file was automatically generated.
// Any changes to it will not be preserved.

package main

import (
	"encoding/json"
	"net/http"
)

var api map[string]ApiHandler

func init() {
	api = make(map[string]ApiHandler)
	api["/api/config"] = api_config
}

// A handler for api calls.
type ApiHandler func(*http.Request) ([]byte, int)

// Pack turns the given value into a JSON encoded byte slice.
// It panics if something went wrong.
func Pack(v interface{}) []byte {
	data, err := json.Marshal(v)

	if err != nil {
		panic(err)
	}

	return data
}
