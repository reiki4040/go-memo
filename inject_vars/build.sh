#!/bin/bash

VERSION=0.1.0
HASH=$(git rev-parse --verify HEAD)
BUILDDATE=$(date '+%Y/%m/%d %H:%M:%S %Z')
GOVERSION=$(go version)

go build -ldflags "-X main.version=$VERSION -X main.hash=$HASH -X \"main.builddate=$BUILDDATE\" -X \"main.goversion=$GOVERSION\""
