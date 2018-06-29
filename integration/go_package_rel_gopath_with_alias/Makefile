FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
LAST_GOPATH:=$(lastword $(subst :, ,$(GOPATH)))
GRPC_GATEWAY_PATH?=${FIRST_GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway
GEN_CLAY_BIN?=$(shell which protoc-gen-goclay)

pwd:
	@pwd

clean:
	find . -regex "\.\/.*\/.*\.go" -exec rm {} +
	rm -f main

protoc:
	protoc --plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) -I/usr/local/include:${GRPC_GATEWAY_PATH}/third_party/googleapis:. --go_out=plugins=grpc:${LAST_GOPATH}/src --goclay_out=impl=true,impl_path=../strings:${LAST_GOPATH}/src pb/strings.proto

build:
	go build -o main main.go
	vgo build -o main main.go

test: pwd clean protoc build
