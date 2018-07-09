FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
LAST_GOPATH:=$(lastword $(subst :, ,$(GOPATH)))
GRPC_GATEWAY_PATH?=${FIRST_GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway
GEN_CLAY_BIN?=$(shell which protoc-gen-goclay)
GEN_GOFAST_BIN?=$(shell which protoc-gen-gofast)

pwd:
	@pwd

clean:
	find . -regex "\.\/.*\/.*\.go" -exec rm {} +
	rm -f main

protoc:
	protoc -I/usr/local/include:. --go_out=. test/test.proto
	protoc --plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --plugin=protoc-gen-gofast=$(GEN_GOFAST_BIN) -I/usr/local/include:${GRPC_GATEWAY_PATH}/third_party/googleapis:${LAST_GOPATH}/src:. --go_out=plugins=grpc:. --goclay_out=impl=true,impl_path=../strings:. pb/strings.proto

build:
	go build -o main main.go
	vgo build -o main main.go

test: pwd clean protoc build
