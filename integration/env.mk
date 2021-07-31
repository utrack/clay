DIR:=$(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

GOPATH?=$(HOME)/go
FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
GOBIN:=$(FIRST_GOPATH)/bin

LOCAL_BIN:=$(DIR)/bin
GEN_CLAY_BIN:=$(DIR)/bin/protoc-gen-goclay
export GEN_CLAY_BIN
GEN_GO_BIN:=$(DIR)/bin/protoc-gen-go
export GEN_GO_BIN
GEN_GO_GRPC_BIN:=$(DIR)/bin/protoc-gen-go-grpc
export GEN_GO_GRPC_BIN
GEN_GOFAST_BIN:=$(DIR)/bin/protoc-gen-gofast
export GEN_GOFAST_BIN
GEN_GOGOFAST_BIN:=$(DIR)/bin/protoc-gen-gogofast
export GEN_GOGOFAST_BIN

#GRPC_GATEWAY_PKG:=$(shell go list -m all | grep github.com/grpc-ecosystem/grpc-gateway/v2 | awk '{print ($$4 != "" ? $$4 : $$1)}')
#GRPC_GATEWAY_VERSION:=$(shell go list -m all | grep github.com/grpc-ecosystem/grpc-gateway/v2 | awk '{print ($$5 != "" ? $$5 : $$2)}')
#GRPC_GATEWAY_PATH:=${FIRST_GOPATH}/pkg/mod/${GRPC_GATEWAY_PKG}@${GRPC_GATEWAY_VERSION}
#export GRPC_GATEWAY_PATH
#$(info ${GRPC_GATEWAY_PATH})
export THIRD_PARTY_PROTO_PATH:=$(dir $(abspath $(lastword $(MAKEFILE_LIST))))third_party/proto
$(info ${THIRD_PARTY_PROTO_PATH})

#GRPC_GOGO_PROTO_PKG:=$(shell go list -m all | grep github.com/gogo/protobuf | awk '{print ($$4 != "" ? $$4 : $$1)}')
#GRPC_GOGO_PROTO_VERSION:=$(shell go list -m all | grep github.com/gogo/protobuf | awk '{print ($$5 != "" ? $$5 : $$2)}')
#GPRC_GOGO_PROTO_PATH:=${FIRST_GOPATH}/pkg/mod/${GRPC_GOGO_PROTO_PKG}@${GRPC_GOGO_PROTO_VERSION}/gogoproto
#export GPRC_GOGO_PROTO_PATH

#ifeq ($(RELATIVE_PATH_TO_PROTO),)
#RELATIVE_PATH_TO_PROTO := pb/strings.proto
#endif
#
#ifeq ($(RELATIVE_PATH_IMPL),)
#RELATIVE_PATH_IMPL := ../strings
#endif

GREEN:=\033[0;32m
RED:=\033[0;31m
NC=:\033[0m

protoc-build:
	$(info #Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install -mod=mod github.com/utrack/clay/v2/cmd/protoc-gen-goclay
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/protobuf/cmd/protoc-gen-go
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc

.protoc_pb: protoc-build
	protoc \
		--plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --goclay_out=. --goclay_opt=impl=true,impl_path=../strings,paths=source_relative \
		--plugin=protoc-gen-go=$(GEN_GO_BIN) --go_out=. --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(GEN_GO_GRPC_BIN) --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		-I/usr/local/include:${THIRD_PARTY_PROTO_PATH}:. \
		pb/strings.proto

.protoc_pb_strings: protoc-build
	protoc \
		--plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --goclay_out=. --goclay_opt=impl=true,impl_path=../../strings,paths=source_relative \
		--plugin=protoc-gen-go=$(GEN_GO_BIN) --go_out=. --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(GEN_GO_GRPC_BIN) --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		-I/usr/local/include:${THIRD_PARTY_PROTO_PATH}:. \
		pb/strings/strings.proto

.build:
	go build -mod=mod -o main main.go
