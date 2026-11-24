# Strillone on Fly.io.
# The Go version tracks .tool-versions and go.mod. Pinned by digest as well as
# tag, so a rebuild of the same commit gets the same base images. Dependabot
# moves the digests (and Go patch releases) forward.
FROM golang:1.27.2-trixie@sha256:47667934874913a449b96faa7ef4777fa285d9ec29126ee80ed61c45b7cccf0e AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY Makefile ./
COPY cmd ./cmd
COPY internal ./internal
# The deploy passes the commit sha, which the root endpoint does not show but
# the startup log line does.
ARG PACKAGER_VERSION=development
ENV CGO_ENABLED=0
RUN make build PACKAGER_VERSION="${PACKAGER_VERSION}"

FROM gcr.io/distroless/static-debian13:nonroot@sha256:8ce6dd87e9823c5b48dcfb2c8604714ccb351ce59607327bcc4957448da279b5
COPY --from=build /src/bin/strillone /strillone
EXPOSE 4000
CMD ["/strillone"]
