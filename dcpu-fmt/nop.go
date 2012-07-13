// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import "io"

type NopCloseReader struct {
	r io.Reader
}

func (n *NopCloseReader) Close() error { return nil }
func (n *NopCloseReader) Read(p []byte) (int, error) {
	return n.r.Read(p)
}

type NopCloseWriter struct {
	w io.Writer
}

func (n *NopCloseWriter) Close() error { return nil }
func (n *NopCloseWriter) Write(p []byte) (int, error) {
	return n.w.Write(p)
}
