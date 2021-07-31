include ../env.mk

$(shell go get github.com/googleapis/googleapis)
GOOGLEAPIS_PKG:=$(shell go list -m all | grep github.com/googleapis/googleapis | awk '{print ($$4 != "" ? $$4 : $$1)}')
GOOGLEAPIS_VERSION:=$(shell go list -m all | grep github.com/googleapis/googleapis | awk '{print ($$5 != "" ? $$5 : $$2)}')
GOOGLEAPIS_PATH:=${FIRST_GOPATH}/pkg/mod/${GOOGLEAPIS_PKG}@${GOOGLEAPIS_VERSION}

pwd:
	@pwd

clean:
	rm -f ./pb/strings.pb.go
	rm -f ./pb/strings_grpc.pb.go
	rm -f ./pb/strings.pb.goclay.go
	rm -f main

protoc: .protoc_pb

build: .build

test: pwd clean protoc build
	go test -v ./...
