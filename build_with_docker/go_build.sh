#!/bin/bash
glide up

VERSION=0.1.0
HASH=$(git rev-parse --verify HEAD)
GOVERSION=$(go version)

gox -output="bin/hello-go_{{.OS}}_{{.Arch}}" -os="linux" -os="darwin" -arch="amd64" -ldflags "-X main.version=$VERSION -X main.hash=$HASH -X \"main.goversion=$GOVERSION\""
