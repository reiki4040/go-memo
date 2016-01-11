#!/bin/bash

# build docker image for go build
docker build -t hello-go:latest .

# for binary
OUTPUT_DIR="./bin"
mkdir -p $OUTPUT_DIR

# build
docker run --volume="$(pwd)/$OUTPUT_DIR:/go/src/github.com/reiki4040/go-memo/build_with_docker/bin" hello-go /bin/bash /go/src/github.com/reiki4040/go-memo/build_with_docker/go_build.sh

