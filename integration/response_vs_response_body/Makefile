include ../env.mk

GOGO_PKG:=$(shell go list -m all | grep github.com/gogo/protobuf | awk '{print ($$4 != "" ? $$4 : $$1)}')
GOGO_VERSION:=$(shell go list -m all | grep github.com/gogo/protobuf | awk '{print ($$5 != "" ? $$5 : $$2)}')
GOGO_PATH:=${FIRST_GOPATH}/pkg/mod/${GOGO_PKG}@${GOGO_VERSION}

pwd:
	@pwd

clean:
	rm -f ./pb/strings.pb.go
	rm -f ./pb/strings_grpc.pb.go
	rm -f ./pb/strings.pb.goclay.go
	rm -f ./strings/strings.pb.impl.go
	rm -f ./strings/strings.go
	rm -f main

protoc: .protoc_pb

build: .build

test: pwd clean protoc build
	go test -v ./...
