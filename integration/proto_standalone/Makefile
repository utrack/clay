FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
GRPC_GATEWAY_PATH?=${FIRST_GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway

clean:
	find . -regex "\.\/.*\/.*\.go" -exec rm {} +
	rm -f main

protoc:
	protoc -I/usr/local/include:${GRPC_GATEWAY_PATH}/third_party/googleapis:proto --gofast_out=plugins=grpc:pb/strings --goclay_out=impl=true,desc_path=pb/strings,impl_path=strings:. proto/strings.proto

build:
	go build -o main main.go && echo

test: clean protoc build