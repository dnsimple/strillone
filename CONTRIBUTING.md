# Contributing

## Getting Started

Clone the repository and move into it:

```shell
git clone git@github.com:dnsimple/strillone.git
cd strillone
```

Install standard Go development tooling:

```shell
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

Install standard Application development tooling:

```shell
brew install overmind
```

## Compilation

```shell
make build
```

## Testing

To run the test suite:

```shell
make test
```

## Running

```shell
make start
```

## Go standard tooling

[Go Development Tooling Wiki](https://dnsimple.atlassian.net/wiki/spaces/DEV/pages/440139826/Go+Projects)

## Go version management

The current Go version is defined in the `.tool-versions` file. Contributors are expected to use `asdf` to install and manage Go running environments.

## Dependency management

Dependencies are managed using [Go modules](https://github.com/golang/go/wiki/Modules). Learn how to [update the dependencies](https://go.dev/wiki/Modules).
