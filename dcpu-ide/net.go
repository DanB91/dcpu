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
	http.HandleFunc("/", wrappedHandler(handler))

	log.Printf("Listening on %q.\n", address)

	return http.ListenAndServe(address, nil)
}

// handler handles each incoming HTTP request.
func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.RemoteAddr, r.URL)

	if r.URL.Path == "/" {
		r.URL.Path = DefaultPage
	}

	// Are we serving a static file?
	if file, ok := static[r.URL.Path]; ok {
		hdr := w.Header()
		hdr.Set("Content-Type", file.ContentType)

		// These files do not change, so give them a liberal cache policy.
		t := time.Now().UTC()
		t = t.Add(time.Duration(t.Second()+3600) * time.Second)
		ts := t.Format(time.RFC1123)

		hdr.Set("Cache-Control", "public")
		hdr.Set("Expires", ts[:len(ts)-4]+" GMT")

		w.Write(file.Data())
		return
	}

	// An api call then?
	if handler, ok := api[r.URL.Path]; ok {
		hdr := w.Header()
		hdr.Set("Content-Type", "application/x-javascript")

		// These should expire immediately.
		t := time.Unix(0, 0).UTC()
		ts := t.Format(time.RFC1123)

		hdr.Set("Cache-Control", "private")
		hdr.Set("Expires", ts[:len(ts)-4]+" GMT")

		var ar ApiResponse
		ar.HttpStatus = 200

		handler(r, &ar)

		w.WriteHeader(ar.HttpStatus)
		w.Write(ar.Pack())
		return
	}

	w.WriteHeader(404)
}

// wrappedHandler returns an http.HandlerFunc that sets some default
// response headers and optionally converts our response writer to a gzipped
// response writer.
func wrappedHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hdr := w.Header()
		hdr.Set("Server", fmt.Sprintf("%s/%d.%d.%s",
			AppName, AppVersionMajor, AppVersionMinor, AppVersionRev))

		if strings.Index(r.Header.Get("Accept-Encoding"), "gzip") == -1 {
			fn(w, r)
			return
		}

		hdr.Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)

		fn(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
		gz.Close()
	}
}
