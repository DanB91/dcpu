# This file is subject to a 1-clause BSD license.
# Its contents can be found in the enclosed LICENSE file.

# The `-X ...` part sets the main.AppVersionRev variable in all
# executables to the current unix timestamp. This allows us to
# use automatic build number increments on each compile cycle.
LDFLAGS = -ldflags "-X main.AppVersionRev `date -u +%s` -s"

all: build

build:
	cd dcpu-asm && go build $(LDFLAGS)
	cd dcpu-data && go build $(LDFLAGS)
	cd dcpu-fmt && go build $(LDFLAGS)
	cd dcpu-prof && go build $(LDFLAGS)
	cd dcpu-test && go build $(LDFLAGS)
	cd dcpu-ide && make build

install:
	cd dcpu-asm && go install $(LDFLAGS)
	cd dcpu-data && go install $(LDFLAGS)
	cd dcpu-fmt && go install $(LDFLAGS)
	cd dcpu-prof && go install $(LDFLAGS)
	cd dcpu-test && go install $(LDFLAGS)
	cd dcpu-ide && make install

clean:
	go clean ./...
	cd dcpu-ide && make clean

fmt:
	go fmt ./...


