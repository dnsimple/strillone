ifdef PACKAGER_VERSION
PKG_VERSION := $(PACKAGER_VERSION)
else
PKG_VERSION := development
endif

all: build

.PHONY: build
build:
	go build -ldflags "-X github.com/dnsimple/strillone/internal/config.Version=$(PKG_VERSION)" -o bin/strillone cmd/strillone/*.go

.PHONY: clean
clean:
	rm bin/strillone

.PHONY: test
test:
	go test -v ./...

.PHONY: fmt
fmt:
	gofumpt -l -w .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: start
start: build
	overmind start
