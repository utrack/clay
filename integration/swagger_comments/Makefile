GOPATH?=$(HOME)/go
FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
GRPC_GATEWAY_PATH?=${FIRST_GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway
GEN_CLAY_BIN?=$(shell which protoc-gen-goclay)

pwd:
	@pwd

clean:
	rm -f ./pb/strings/strings.pb.go
	rm -f ./pb/strings/strings.pb.goclay.go
	rm -f ./strings/strings.pb.impl.go

protoc:
	protoc --plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) -I/usr/local/include:${GRPC_GATEWAY_PATH}/third_party/googleapis:. --gofast_out=plugins=grpc:. --goclay_out=impl=true,impl_path=../../strings:. pb/strings/strings.proto

build:

test: pwd clean protoc build
	vgo test ./pb/strings
