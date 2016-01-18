#!/bin/bash

# build docker image for go build
docker build -t golangbuildbase:latest -f Dockerfile.base .

docker build --no-cache -t hello-go:latest -f Dockerfile.build .

# for binary
OUTPUT_DIR="./bin"
mkdir -p $OUTPUT_DIR

# build
docker run --volume="$(pwd)/$OUTPUT_DIR:/go/src/github.com/reiki4040/go-memo/build_with_docker/bin" -w /go/src/github.com/reiki4040/go-memo/build_with_docker hello-go /bin/bash go_build.sh

