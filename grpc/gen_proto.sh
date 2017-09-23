#!/bin/bash

# sample directory structure
# --------------------------------------
# - gen_proto.sh (this file)
# - address/
#     +- address.proto
#     +- address.pb.go (will generate)
# - geocoder_server/ (has file)
# - geocoder_client/ (has file)
# --------------------------------------
#
protoc -I address/ address/address.proto --go_out=plugins=grpc:address

# -I address/ : if not specified, then generate go file to
# address/address/address.pb.go (nest again)
#
# --go_out option explain in below blog entry.
# https://blog.fenrir-inc.com/jp/2016/10/grpc-go.html
# grpc: is generate gRPC code,
# and :address is directory path that store genrated go code
# if :. specified, then address.pb.go is stored to current directory.
#
# so if generate code from .proto that in current directory (in address/ on this sample)
# then run below command
# protoc address.proto --go_out=plugins=grpc:.
