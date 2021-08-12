include ../env.mk

pwd:
	@pwd

clean:
	rm -f ./pb/strings.pb.go
	rm -f ./pb/strings_grpc.pb.go
	rm -f ./pb/strings.pb.goclay.go
	rm -f ./strings/strings.go
	rm -f ./strings/to_lower.go
	rm -f ./strings/to_lower_test.go
	rm -f ./strings/to_upper.go
	rm -f ./strings/to_upper_test.go
	rm -f main

protoc: .protoc_pb

build: .build

test: pwd clean protoc
	go test ./pb
