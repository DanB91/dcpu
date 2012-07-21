# This file is subject to a 1-clause BSD license.
# Its contents can be found in the enclosed LICENSE file.

TOC = ../static.go
DATADIR = data
DEPLOYDIR = deploy
CODEDIR = ../
PREFIX = data_

## These are used to generate error constant and message 
## definitions which are identical in Go and Javascript.
ERR_DATA = api_errors.json
ERR_GOIN = api_errors.gotmpl
ERR_GOOUT = api_errors.go
ERR_JSIN = api_errors.jstmpl
ERR_JSOUT = data/js/api_errors.js

LDFLAGS_BUILD = -ldflags "-X main.AppVersionRev `date -u +%s`"
LDFLAGS_INSTALL = -ldflags "-X main.AppVersionRev `date -u +%s` -s"
CLEANFILES = static.go $(PREFIX)*.go $(DEPLOYDIR) $(ERR_GOOUT)

all: build

build:
	cd codegen && go run *.go -d \
		-i ../$(DATADIR) -o $(CODEDIR) -t $(TOC) -p $(PREFIX) \
		-data ../$(ERR_DATA) \
		-goin ../$(ERR_GOIN) -goout ../$(ERR_GOOUT) \
		-jsin ../$(ERR_JSIN) -jsout ../$(ERR_JSOUT)
	go build $(LDFLAGS_BUILD)
	rm -rf $(CLEANFILES)

install:
	cd codegen && go run *.go -data ../$(ERR_DATA) \
		-goin ../$(ERR_GOIN) -goout ../$(ERR_GOOUT) \
		-jsin ../$(ERR_JSIN) -jsout ../$(ERR_JSOUT)
	./compress_data.sh
	cd codegen && go run *.go -i ../$(DEPLOYDIR) -o $(CODEDIR) -t $(TOC) -p $(PREFIX)
	go install $(LDFLAGS_INSTALL)
	rm -rf $(CLEANFILES) $(ERR_JSOUT)

clean:
	go clean ./...
	rm -rf $(CLEANFILES)

fmt:
	go fmt ./...

