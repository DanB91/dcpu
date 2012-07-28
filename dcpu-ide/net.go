// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"code.google.com/p/go.net/websocket"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
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

func launchServer(address string) (err error) {
	http.Handle("/ws", websocket.Handler(wsHandler))
	http.Handle("/", wrappedHttpHandler(httpHandler))

	log.Printf("Listening on %s", address)
	return http.ListenAndServe(address, nil)
}

func wsHandler(ws *websocket.Conn) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("%s [p] %v", ws.Request().RemoteAddr, x)
		}

		websocket.Message.Send(ws, []byte{ErrUnknown})
	}()

	NewClient(ws).Poll()
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.RemoteAddr, r.URL)

	if r.URL.Path == "/" {
		r.URL.Path = DefaultPage
	}

	// Are we serving a static file?
	file, ok := static[r.URL.Path]
	if !ok {
		w.WriteHeader(404) // Whatever it is we're looking for, it aint here.
		return
	}

	hdr := w.Header()
	hdr.Set("Content-Type", file.ContentType)

	// HTML and JS files go through our template engine first.
	// The rest can just be written as-is.
	if !strings.HasPrefix(file.ContentType, "text/html") &&
		!strings.HasPrefix(file.ContentType, "application/x-javascript") {
		w.Write(file.Data())
		return
	}

	// HTML and JS content might need some extra processing.
	data, err := parseTemplate(file.Data())

	if err != nil {
		log.Printf("Template error for %s: %v", r.URL.Path, err)
		w.WriteHeader(500)
		return
	}

	w.Write(data)
}

// wrappedHttpHandler returns an http.HandlerFunc that sets some default
// response headers and optionally converts our response writer to a gzipped
// response writer.
func wrappedHttpHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Static files do not change, so give them a liberal cache policy.
		// One month into the future should be enough.
		t := time.Now().UTC().Add(time.Hour * 24 * 30)
		ts := t.Format(time.RFC1123)
		ts = ts[:len(ts)-4] + " GMT"

		hdr := w.Header()
		hdr.Set("Server", fmt.Sprintf("%s/%d.%d",
			AppName, AppVersionMajor, AppVersionMinor))
		hdr.Set("Cache-Control", "public")
		hdr.Set("Expires", ts)

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
