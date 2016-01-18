#!/bin/bash

# build docker image for go build
echo "build golang base and toolchain..."
docker build -t golangbuildbase:latest -f Dockerfile.base .

echo -e "\nclone target project..."
docker build -t hello-go:latest --no-cache -f Dockerfile.build .

# for binary
OUTPUT_DIR="./bin"
mkdir -p $OUTPUT_DIR

echo -e "\nmake binary..."
# build
docker run --volume="$(pwd)/$OUTPUT_DIR:/go/src/github.com/reiki4040/go-memo/build_with_docker/bin" -w /go/src/github.com/reiki4040/go-memo/build_with_docker hello-go /bin/bash go_build.sh

