# This file is subject to a 1-clause BSD license.
# Its contents can be found in the enclosed LICENSE file.

## Misc stuff.
CG = codegen
CGC = ./$(CG)/$(CG)
DATADIR = data
DEPLOYDIR = deploy
TMPLDIR = templates

## Build Stage 1
## Generates Go and javascript code from templates.
CG_DATA = codegen.json
CG_ERR_GOIN = $(TMPLDIR)/api_errors.go.tmpl
CG_ERR_GOOUT = api_errors.go
CG_ERR_JSIN = $(TMPLDIR)/api_errors.js.tmpl
CG_ERR_JSOUT = $(DATADIR)/js/api_errors.js
CG_INDEX_IN = $(TMPLDIR)/index.html.tmpl
CG_INDEX_OUT = $(DATADIR)/index.html

## Build Stage 2
## Merges multiple stylesheets and javascript files into a single
## file for each type. Only applies to install builds.
MERGE_OUT_JS = $(DEPLOYDIR)/app.js
MERGE_OUT_CSS = $(DEPLOYDIR)/app.css
MERGE_PREFIX = ./$(DEPLOYDIR)/

## Build Stage 3
## Generate Go code from static web data for embedding.
PREFIX = data_
TOC = static.go

## Build Stage 4
## Compile the actual ide and clean up our mess.
LDFLAGS_BUILD = -ldflags "-X main.AppVersionRev `date -u +%s`"
LDFLAGS_INSTALL = -ldflags "-X main.AppVersionRev `date -u +%s` -s"
CLEAN_BUILD = static.go $(DEPLOYDIR) $(CG_ERR_GOOUT) $(PREFIX)*.go $(CGC)
CLEAN_INSTALL = $(CLEAN_BUILD) $(CG_ERR_JSOUT) $(CG_INDEX_OUT) \
	$(MERGE_OUT_JS) $(MERGE_OUT_CSS)

all: build

build: 
	cd $(CG) && go build -ldflags="-s"
	$(CGC) -d -data $(CG_DATA) -cgin $(CG_ERR_GOIN) -cgout $(CG_ERR_GOOUT)
	$(CGC) -d -data $(CG_DATA) -cgin $(CG_ERR_JSIN) -cgout $(CG_ERR_JSOUT)
	$(CGC) -d -data $(CG_DATA) -cgin $(CG_INDEX_IN) -cgout $(CG_INDEX_OUT)
	$(CGC) -d -convin $(DATADIR) -convout . -convtoc $(TOC) -convpre $(PREFIX)
	go build $(LDFLAGS_BUILD)
	rm -rf $(CLEAN_BUILD)

install: 
	cd $(CG) && go build -ldflags="-s"
	$(CGC) -data $(CG_DATA) -cgin $(CG_ERR_GOIN) -cgout $(CG_ERR_GOOUT)
	$(CGC) -data $(CG_DATA) -cgin $(CG_ERR_JSIN) -cgout $(CG_ERR_JSOUT)
	$(CGC) -data $(CG_DATA) -cgin $(CG_INDEX_IN) -cgout $(CG_INDEX_OUT)
	./compress_data.sh
	$(CGC) -data $(CG_DATA) -mergetype js -mergepre $(MERGE_PREFIX) -mergeout $(MERGE_OUT_JS)
	$(CGC) -data $(CG_DATA) -mergetype css -mergepre $(MERGE_PREFIX) -mergeout $(MERGE_OUT_CSS)
	$(CGC) -convin $(DEPLOYDIR) -convout . -convtoc $(TOC) -convpre $(PREFIX)
	go install $(LDFLAGS_INSTALL)
	rm -rf $(CLEAN_INSTALL)

clean:
	go clean
	rm -rf $(CLEAN_INSTALL)

fmt:
	go fmt ./...

