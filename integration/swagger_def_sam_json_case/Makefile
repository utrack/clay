FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))
GRPC_GATEWAY_PATH?=${FIRST_GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway

pwd:
	@pwd

clean:
	rm -f ./pb/strings/strings.pb.go
	rm -f ./pb/strings/strings.pb.goclay.go
	rm -f ./strings/strings.pb.impl.go
	rm -f main

protoc:
	protoc -I/usr/local/include:${GRPC_GATEWAY_PATH}/third_party/googleapis:. --gofast_out=plugins=grpc:. --goclay_out=impl=true,impl_path=../../strings:. pb/strings/strings.proto

build:
	go build -o main main.go
	vgo build -o main main.go

test: pwd clean protoc build
	go test ./pb/strings
