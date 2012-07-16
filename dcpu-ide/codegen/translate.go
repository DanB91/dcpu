// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"
)

func translate(fin string) (err error) {
	relpath := strings.Replace(fin, *indir, "", 1)
	fout := strings.Replace(relpath, string(filepath.Separator), "_", -1)

	varname := fout
	varname = strings.ToLower(varname)
	varname = strings.Replace(varname, " ", "_", -1)
	varname = strings.Replace(varname, ".", "_", -1)
	varname = strings.Replace(varname, "-", "_", -1)

	if unicode.IsDigit(rune(varname[0])) {
		varname = "_" + varname
	}

	fout = *prefix + fout + ".go"
	fout = filepath.Join(*outdir, fout)
	files = append(files, File{
		Path: relpath,
		Var:  varname,
		Type: mime.TypeByExtension(path.Ext(fin)),
	})

	fmt.Printf("[*] %s => %s\n", fin, fout)

	fs, err := os.Open(fin)
	if err != nil {
		return
	}

	defer fs.Close()

	fd, err := os.Create(fout)
	if err != nil {
		return
	}

	defer fd.Close()

	if *dev {
		translate_dev(fs, fd, varname, fin)
	} else {
		translate_live(fs, fd, varname)
	}

	return
}

func translate_dev(fs io.Reader, fd io.Writer, varname, infile string) {
	infile = strings.Replace(infile, "../data", "data", 1)

	fmt.Fprintf(fd, `
// This file was automatically generated.
// Any changes to it will not be preserved.

package main

import "io/ioutil"

func %s() []byte {
	data, err := ioutil.ReadFile(%q)

	if err != nil {
		panic("%s: " + err.Error())
	}

	return data
}`, varname, infile, infile)
}

func translate_live(fs io.Reader, fd io.Writer, varname string) {
	fmt.Fprintf(fd, `
// This file was automatically generated.
// Any changes to it will not be preserved.

package main

import (
	"bytes"
	"compress/gzip"
	"io"
)

func %s() []byte {
	gz, err := gzip.NewReader(bytes.NewBuffer([]byte{`, varname)

	gz := gzip.NewWriter(&ByteWriter{Writer: fd})
	io.Copy(gz, fs)
	gz.Close()

	fmt.Fprint(fd, `
	}))

	if err != nil {
		panic("Decompression failed: " + err.Error())
	}

	var b bytes.Buffer
	io.Copy(&b, gz)
	gz.Close()

	return b.Bytes()
}`)
}
