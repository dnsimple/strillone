test:
	go test -v ./...

build:
	go build -ldflags "-X main.Version=$$(git rev-parse --short HEAD)" -o bin/strillone

clean:
	rm bin/strillone

start: build
	overmind start

lint:
	golint ./...
