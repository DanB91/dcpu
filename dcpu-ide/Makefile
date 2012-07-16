# This file is subject to a 1-clause BSD license.
# Its contents can be found in the enclosed LICENSE file.

TOC = ../static.go
INDIR = ../data
OUTDIR = ../
PREFIX = data_
LDFLAGS = -ldflags "-X main.AppVersionRev `date -u +%s` -s"
CLEANFILES = static.go $(PREFIX)*.go

all: build

build:
	cd codegen && go run *.go -d -i $(INDIR) -o $(OUTDIR) -t $(TOC) -p $(PREFIX)
	go build $(LDFLAGS)

install:
	cd codegen && go run *.go -i $(INDIR) -o $(OUTDIR) -t $(TOC) -p $(PREFIX)
	go install $(LDFLAGS)

clean:
	go clean ./...
	rm -f $(CLEANFILES)

fmt:
	go fmt ./...

