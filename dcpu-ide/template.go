// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"runtime"
	"text/template"
)

var funcs template.FuncMap

func init() {
	funcs = make(template.FuncMap)
	funcs["Version"] = func() string {
		return fmt.Sprintf("%s %d.%d.%s (Go runtime %s).",
			AppName, AppVersionMajor, AppVersionMinor,
			AppVersionRev, runtime.Version())
	}
}

// parseTemplate treats the input data as a template
// and parses it.
func parseTemplate(html []byte) (out []byte, err error) {
	t := template.New("page")
	t.Funcs(funcs)

	t, err = t.Parse(string(html))
	if err != nil {
		return
	}

	addr := config.Address
	if addr[0] == ':' {
		addr = "localhost" + addr
	}

	data := struct {
		SocketAddress string
	}{
		fmt.Sprintf("ws://%s/ws", addr),
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	out = buf.Bytes()
	return
}
