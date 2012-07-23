// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"encoding/json"
	"net/http"
)

type ApiFunc func(*http.Request) ([]byte, int)

// A handler for api calls.
type ApiHandler struct {
	Func   ApiFunc
	Method string
}

var api map[string]*ApiHandler

func Register(path, method string, ah ApiFunc) {
	if api == nil {
		api = make(map[string]*ApiHandler)
	}

	api[path] = &ApiHandler{ah, method}
}

// Pack turns the given value into a JSON encoded byte slice.
// It panics if something went wrong.
func Pack(v interface{}) []byte {
	data, err := json.Marshal(v)

	if err != nil {
		panic(err)
	}

	return data
}

// Error creates an error object we can send as an API response.
func Error(code int, argv ...interface{}) []byte {
	return Pack(struct {
		Code int
		Args []interface{}
	}{
		code, argv,
	})
}
