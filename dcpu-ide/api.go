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

	api["/cake"] = func(r *http.Request, ar *ApiResponse) {
		ar.Msg = "Cake rocks."
	}
}

// A handler for api calls.
type ApiHandler func(*http.Request, *ApiResponse)

// A single Api response.
type ApiResponse struct {
	Msg        string `json:",omitempty"`
	HttpStatus int    `json:"-"`
}

// Pack turns the response struct into a JSON encoded byte slice.
// It panics if something went wrong.
func (ar *ApiResponse) Pack() []byte {
	data, err := json.Marshal(ar)

	if err != nil {
		panic(err)
	}

	return data
}
