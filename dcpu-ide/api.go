// This file was automatically generated.
// Any changes to it will not be preserved.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var api map[string]ApiHandler

func Register(path string, ah ApiHandler) {
	if api == nil {
		api = make(map[string]ApiHandler)
	}

	api[path] = ah
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

// Errorf creates an error object we can send as an API response.
func Errorf(f string, argv ...interface{}) []byte {
	return Pack(struct {
		Message string
	}{
		fmt.Sprintf(f, argv...),
	})
}
