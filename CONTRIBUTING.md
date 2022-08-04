# Contributing

## Getting Started

Ensure your `$GOPATH` is not blank. It should be something like `export GOPATH=$HOME/go`:

```shell
echo $GOPATH
```

Clone the repository [in your workspace](https://golang.org/doc/code.html#Organization) and move into it:

```shell
mkdir -p $GOPATH/src/github.com/dnsimple && cd $_
git clone git@github.com:dnsimple/strillone.git
cd strillone
```

Install standard Go development tooling:

```shell
cd ~
go get -u golang.org/x/lint/golint
```

Install standard Application development tooling:

```shell
cd ~
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

Go modules are enabled. The file `go.mod` MUST include the `go` directive to determine the language feature used by `go` commands.


## Dependency management

Dependencies are managed using [Go modules](https://github.com/golang/go/wiki/Modules). Learn how to [update the dependencies](https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies).
