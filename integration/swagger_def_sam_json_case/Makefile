FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
GRPC_GATEWAY_PATH?=${FIRST_GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway

clean:
	rm ./pb/strings/strings.pb.go
	rm ./pb/strings/strings.pb.goclay.go
	rm ./strings/strings.pb.impl.go
	rm -f main

protoc:
	protoc -I/usr/local/include:${GRPC_GATEWAY_PATH}/third_party/googleapis:. --gofast_out=plugins=grpc:. --goclay_out=impl=true,impl_path=../../strings:. pb/strings/strings.proto

build:
	go build -o main main.go && echo

test: clean protoc build
	go test ./pb/strings
