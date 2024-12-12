#!/bin/sh

mkdir proto/gen
protoc -I proto proto/*.proto --go_out=./proto/gen/ --go_opt=paths=source_relative --go-grpc_out=./proto/gen/ --go-grpc_opt=paths=source_relative
