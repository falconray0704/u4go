#!/bin/bash

export PATH=/md/gows/bin/:$PATH
protoc -I ../../protos --go_out=plugins=grpc:. ../../protos/uGrpc.proto




