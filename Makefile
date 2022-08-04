test:
	go test -v ./...

build:
	go build -ldflags "-X main.Version=$$(git rev-parse --short HEAD)" -o bin/strillone cmd/strillone/*.go

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
