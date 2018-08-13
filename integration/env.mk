DIR:=$(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

GOPATH?=$(HOME)/go
FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
GOBIN:=$(FIRST_GOPATH)/bin

LOCAL_BIN:=$(DIR)/bin
GEN_CLAY_BIN:=$(DIR)/bin/protoc-gen-goclay
export GEN_CLAY_BIN
GEN_GO_BIN:=$(DIR)/bin/protoc-gen-go
export GEN_GO_BIN
GEN_GOFAST_BIN:=$(DIR)/bin/protoc-gen-gofast
export GEN_GOFAST_BIN
GEN_GOGOFAST_BIN:=$(DIR)/bin/protoc-gen-gogofast
export GEN_GOGOFAST_BIN

GRPC_GATEWAY_PKG:=$(shell vgo list -m all | grep github.com/grpc-ecosystem/grpc-gateway | awk '{print ($$4 != "" ? $$4 : $$1)}')
GRPC_GATEWAY_VERSION:=$(shell vgo list -m all | grep github.com/grpc-ecosystem/grpc-gateway | awk '{print ($$5 != "" ? $$5 : $$2)}')
GRPC_GATEWAY_PATH:=${FIRST_GOPATH}/pkg/mod/${GRPC_GATEWAY_PKG}@${GRPC_GATEWAY_VERSION}
export GRPC_GATEWAY_PATH

GREEN:=\033[0;32m
RED:=\033[0;31m
NC=:\033[0m

protoc-build:
	$(info #Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) vgo install github.com/utrack/clay/v2/cmd/protoc-gen-goclay
	GOBIN=$(LOCAL_BIN) vgo install github.com/golang/protobuf/protoc-gen-go
	GOBIN=$(LOCAL_BIN) vgo install github.com/gogo/protobuf/protoc-gen-gofast
	GOBIN=$(LOCAL_BIN) vgo install github.com/gogo/protobuf/protoc-gen-gogofast
