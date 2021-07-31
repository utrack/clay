include ../env.mk

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