#!/bin/bash

# test docker command
rs=$(docker ps)
if [ $? -ne 0 ]; then
  echo "your docker machine is running?"
  exit 1
fi

# build docker image for go build
echo "build golang base and install toolchain to docker image..."
docker build -t golangbuildbase:latest -f Dockerfile.base .
if [ $? -eq 1 ]; then
  echo "failed build base image."
  exit 1
fi

echo -e "\nclone target project into docker image..."
docker build -t hello-go:latest --no-cache -f Dockerfile.build .
if [ $? -eq 1 ]; then
  echo "failed build hello-go image."
  exit 1
fi

# for binary
OUTPUT_DIR="./bin"
mkdir -p $OUTPUT_DIR

echo -e "\nstart make binary..."
# build
docker run --volume="$(pwd)/$OUTPUT_DIR:/go/src/github.com/reiki4040/go-memo/build_with_docker/bin" -w /go/src/github.com/reiki4040/go-memo/build_with_docker hello-go /bin/bash go_build.sh
if [ $? -eq 1 ]; then
  echo "failed make binary."
  exit 1
else
  echo "successed make binary."
fi
