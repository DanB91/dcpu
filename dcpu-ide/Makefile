# This file is subject to a 1-clause BSD license.
# Its contents can be found in the enclosed LICENSE file.

# External dependencies on release build:
#
# - htmlcompressor
# - yuicompressor
# - pngcrush
#

TOC = ../static.go
DATADIR = data
DEPLOYDIR = deploy
CODEDIR = ../
PREFIX = data_
LDFLAGS = -ldflags "-X main.AppVersionRev `date -u +%s` -s"
CLEANFILES = static.go $(PREFIX)*.go $(DEPLOYDIR)

all: build

build:
	cd codegen && go run *.go -d -i ../$(DATADIR) -o $(CODEDIR) -t $(TOC) -p $(PREFIX)
	go build $(LDFLAGS)

install:
	./compress_data.sh
	cd codegen && go run *.go -i ../$(DEPLOYDIR) -o $(CODEDIR) -t $(TOC) -p $(PREFIX)
	go install $(LDFLAGS)

clean:
	go clean ./...
	rm -rf $(CLEANFILES)

fmt:
	go fmt ./...

