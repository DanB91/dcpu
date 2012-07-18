// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// An expiration date somewhere in the distant past.
// Used in response handler for API calls. They should not be cached.
var AncientHistory string

// gzipResponseWriter wraps the http.ResponseWriter in a gzip compressor.
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (this gzipResponseWriter) Write(p []byte) (int, error) {
	return this.Writer.Write(p)
}

// Run starts the webserver on the given address.
func Run(address string) (err error) {
	t := time.Unix(0, 0).UTC()

	AncientHistory = t.Format(time.RFC1123)
	AncientHistory = AncientHistory[:len(AncientHistory)-4] + " GMT"

	http.HandleFunc("/", wrappedHandler(handler))
	log.Printf("Listening on %q.\n", address)

	return http.ListenAndServe(address, nil)
}

// handler handles each incoming HTTP request.
func handler(w http.ResponseWriter, r *http.Request) {
	if tracker != nil {
		tracker.Ping()
	}

	log.Printf("%s %s\n", r.RemoteAddr, r.URL)

	if r.URL.Path == "/" {
		r.URL.Path = DefaultPage
	}

	// Are we serving a static file?
	if file, ok := static[r.URL.Path]; ok {
		hdr := w.Header()
		hdr.Set("Content-Type", file.ContentType)

		// These files do not change, so give them a liberal cache policy.
		// One month into the future should be enough.
		t := time.Now().UTC().Add(time.Hour * 24 * 30)
		ts := t.Format(time.RFC1123)
		ts = ts[:len(ts)-4] + " GMT"

		hdr.Set("Cache-Control", "public")
		hdr.Set("Expires", ts)

		w.Write(file.Data())
		return
	}

	// An api call then?
	if handler, ok := api[r.URL.Path]; ok {
		hdr := w.Header()
		hdr.Set("Content-Type", "application/x-javascript")

		// These should expire immediately.
		hdr.Set("Cache-Control", "private")
		hdr.Set("Expires", AncientHistory)

		data, status := handler(r)

		w.WriteHeader(status)
		w.Write(data)
		return
	}

	// Whatever it is we're looking for, it aint here.
	w.WriteHeader(404)
}

// wrappedHandler returns an http.HandlerFunc that sets some default
// response headers and optionally converts our response writer to a gzipped
// response writer.
func wrappedHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hdr := w.Header()
		hdr.Set("Server", fmt.Sprintf("%s/%d.%d",
			AppName, AppVersionMajor, AppVersionMinor))

		// If the client does not support gzip compression,
		// we should just write our response out normally.
		if strings.Index(r.Header.Get("Accept-Encoding"), "gzip") == -1 {
			fn(w, r)
			return
		}

		// They do accept gzip compressed content; make it so.
		hdr.Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)

		fn(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
		gz.Close()
	}
}
