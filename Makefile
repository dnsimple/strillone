test:
	go test -v ./...

build:
	go build -o bin/strillone cmd/strillone/*.go

clean:
	rm bin/strillone

start: build
	overmind start

lint:
	golangci-lint run

fmt:
	go install mvdan.cc/gofumpt@latest
	go fmt ./...
	gofumpt -w ./
