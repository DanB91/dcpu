# This file is subject to a 1-clause BSD license.
# Its contents can be found in the enclosed LICENSE file.

## Misc stuff.
DATADIR = data
DEPLOYDIR = deploy
BUILDDIR = build
CODEGENDIR = $(BUILDDIR)/codegen
CGC = ./$(CODEGENDIR)/codegen

## Build Stage 1
## Generates shared Go and javascript code from templates.
CG_DATA = $(BUILDDIR)/codegen.json
CG_MSG_GOIN = $(BUILDDIR)/api_message.go.tmpl
CG_MSG_GOOUT = api_message.go
CG_MSG_JSIN = $(BUILDDIR)/api_message.js.tmpl
CG_MSG_JSOUT = $(DATADIR)/js/api_message.js
CG_INDEX_IN = $(BUILDDIR)/index.html.tmpl
CG_INDEX_OUT = $(DATADIR)/index.html

## Build Stage 2
## Merges multiple stylesheets and javascript files into a single
## file for each type. Only applies to install builds.
##
## This serves to reduce the application load time by reducing the
## amount of roundtrips the browser has to make to fetch all code.
MERGE_OUT_JS = $(DEPLOYDIR)/app.js
MERGE_OUT_CSS = $(DEPLOYDIR)/app.css
MERGE_CONV_PREFIX = ./$(DEPLOYDIR)

## Build Stage 3
## Generate Go code from static web data for embedding.
CONV_PREFIX = data_
CONV_TOC = static.go

## Build Stage 4
## Compile the actual program and clean up our mess.
LDFLAGS_BUILD = -ldflags "-X main.AppVersionRev `date -u +%s`"
LDFLAGS_INSTALL = -ldflags "-X main.AppVersionRev `date -u +%s` -s"
CLEAN_BUILD = static.go $(DEPLOYDIR) $(CG_MSG_GOOUT) $(CONV_PREFIX)*.go $(CGC)
CLEAN_INSTALL = $(CLEAN_BUILD) $(CG_MSG_JSOUT) $(CG_INDEX_OUT)

all: debug

debug: 
	cd $(CODEGENDIR) && go build -ldflags="-s"
	$(CGC) -d -data $(CG_DATA) \
	          -cgin $(CG_MSG_GOIN) \
	          -cgout $(CG_MSG_GOOUT)
	$(CGC) -d -data $(CG_DATA) \
	          -cgin $(CG_MSG_JSIN) \
	          -cgout $(CG_MSG_JSOUT)
	$(CGC) -d -data $(CG_DATA) \
	          -cgin $(CG_INDEX_IN) \
	          -cgout $(CG_INDEX_OUT)
	$(CGC) -d -convin $(DATADIR) \
	          -convout . \
	          -convtoc $(CONV_TOC) \
	          -convpre $(CONV_PREFIX)
	go build $(LDFLAGS_BUILD)
	rm -rf $(CLEAN_BUILD)

install: 
	cd $(CODEGENDIR) && go build -ldflags="-s"
	$(CGC) -data $(CG_DATA) \
	       -cgin $(CG_MSG_GOIN) \
	       -cgout $(CG_MSG_GOOUT)
	$(CGC) -data $(CG_DATA) \
	       -cgin $(CG_MSG_JSIN) \
	       -cgout $(CG_MSG_JSOUT)
	$(CGC) -data $(CG_DATA) \
	       -cgin $(CG_INDEX_IN) \
	       -cgout $(CG_INDEX_OUT)
	./$(BUILDDIR)/compress_data.sh
	$(CGC) -data $(CG_DATA) -mergetype js \
	       -mergepre $(MERGE_CONV_PREFIX) \
	       -mergeout $(MERGE_OUT_JS)
	$(CGC) -data $(CG_DATA) -mergetype css \
	       -mergepre $(MERGE_CONV_PREFIX) \
	       -mergeout $(MERGE_OUT_CSS)
	$(CGC) -convin $(DEPLOYDIR) \
	       -convout . \
	       -convtoc $(CONV_TOC) \
	       -convpre $(CONV_PREFIX)
	go install $(LDFLAGS_INSTALL)
	rm -rf $(CLEAN_INSTALL)

clean:
	go clean
	rm -rf $(CLEAN_INSTALL)

fmt:
	go fmt ./...

