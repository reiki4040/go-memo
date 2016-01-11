#!/bin/bash
cd /go/src/github.com/reiki4040/go-memo/build_with_docker
glide up
gox -output="bin/hello-go" --osarch="darwin/amd64"
