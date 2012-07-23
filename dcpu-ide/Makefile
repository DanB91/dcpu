# This file is subject to a 1-clause BSD license.
# Its contents can be found in the enclosed LICENSE file.

## Misc stuff.
CG = codegen
DATADIR = data
DEPLOYDIR = deploy
TMPLDIR = templates

## Build Stage 1
## These are used to generate Go and javascript code.
CG_DATA = codegen.json
CG_ERR_GOIN = $(TMPLDIR)/api_errors.go.tmpl
CG_ERR_GOOUT = api_errors.go
CG_ERR_JSIN = $(TMPLDIR)/api_errors.js.tmpl
CG_ERR_JSOUT = $(DATADIR)/js/api_errors.js
CG_INDEX_IN = $(TMPLDIR)/index.html.tmpl
CG_INDEX_OUT = $(DATADIR)/index.html

## Build Stage 2
## This is used to generate Go code from static web data.
## Release build will embed files in the final executable.
PREFIX = data_
TOC = static.go

## Build Stage 3
LDFLAGS_BUILD = -ldflags "-X main.AppVersionRev `date -u +%s`"
LDFLAGS_INSTALL = -ldflags "-X main.AppVersionRev `date -u +%s` -s"
CLEAN_BUILD = static.go $(DEPLOYDIR) $(CG_ERR_GOOUT) $(PREFIX)*.go
CLEAN_INSTALL = $(CLEAN_BUILD) $(CG_ERR_JSOUT) $(CG_INDEX_OUT) ./$(CG)/$(CG)

all: build

build: 
	cd $(CG) && go build -ldflags="-s"
	./$(CG)/$(CG) -d -cgdata $(CG_DATA) -cgin $(CG_ERR_GOIN) -cgout $(CG_ERR_GOOUT)
	./$(CG)/$(CG) -d -cgdata $(CG_DATA) -cgin $(CG_ERR_JSIN) -cgout $(CG_ERR_JSOUT)
	./$(CG)/$(CG) -d -cgdata $(CG_DATA) -cgin $(CG_INDEX_IN) -cgout $(CG_INDEX_OUT)
	./$(CG)/$(CG) -d -i $(DATADIR) -o . -t $(TOC) -p $(PREFIX)
	go build $(LDFLAGS_BUILD)
	rm -rf $(CLEAN_BUILD)

install: 
	cd $(CG) && go build -ldflags="-s"
	./$(CG)/$(CG) -cgdata $(CG_DATA) -cgin $(CG_ERR_GOIN) -cgout $(CG_ERR_GOOUT)
	./$(CG)/$(CG) -cgdata $(CG_DATA) -cgin $(CG_ERR_JSIN) -cgout $(CG_ERR_JSOUT)
	./$(CG)/$(CG) -cgdata $(CG_DATA) -cgin $(CG_INDEX_IN) -cgout $(CG_INDEX_OUT)
	./compress_data.sh
	./$(CG)/$(CG) -i $(DEPLOYDIR) -o . -t $(TOC) -p $(PREFIX)
	go install $(LDFLAGS_INSTALL)
	rm -rf $(CLEAN_INSTALL)

clean:
	go clean
	rm -rf $(CLEAN_INSTALL)

fmt:
	go fmt ./...

