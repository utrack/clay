include ../env.mk

GOGO_PKG:=$(shell go list -m all | grep github.com/gogo/protobuf | awk '{print ($$4 != "" ? $$4 : $$1)}')
GOGO_VERSION:=$(shell go list -m all | grep github.com/gogo/protobuf | awk '{print ($$5 != "" ? $$5 : $$2)}')
GOGO_PATH:=${FIRST_GOPATH}/pkg/mod/${GOGO_PKG}@${GOGO_VERSION}

pwd:
	@pwd

clean:
	rm -f ./pb/strings.pb.go
	rm -f ./pb/strings.pb.goclay.go
	rm -f ./strings/strings.go
	rm -f main

protoc: protoc-build
	protoc --plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --plugin=protoc-gen-gofast=$(GEN_GOFAST_BIN) -I/usr/local/include:${GRPC_GATEWAY_PATH}/third_party/googleapis:${GOGO_PATH}:. --gofast_out=goproto_registration=true,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,plugins=grpc:. --goclay_out=impl=true,impl_path=../strings:. pb/strings.proto

build:
	go build -o main main.go

test: pwd clean protoc build
	go test -v .
