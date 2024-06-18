#!/bin/bash

INCLUDES="-I=./api/proto/tron -I/usr/lib -I./api/proto/googleapis"
FLAGS="--go_out=./pkg/tronpb --go_opt paths=source_relative --go-grpc_out=./pkg/tronpb --go-grpc_opt=paths=source_relative"
OLD_PACKAGE="github.com/tronprotocol/grpc-gateway"
NEW_PACKAGE="github.com/TheTeaParty/trongrid/pkg/tronpb"
PROTO_DIR="./api/proto/tron"


# Run the Go script to replace the go_package option
go run util/replacer/main.go $PROTO_DIR $OLD_PACKAGE $NEW_PACKAGE

# Run protoc
protoc ${INCLUDES} ${FLAGS} ${PROTO_DIR}/core/*.proto ${PROTO_DIR}/core/contract/*.proto ${PROTO_DIR}/api/*.proto

# Fix the imports
mv pkg/tronpb/core/contract/* pkg/tronpb/core/

echo "Done"