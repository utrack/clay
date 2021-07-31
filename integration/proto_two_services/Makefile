include ../env.mk

pwd:
	@pwd

clean:
	rm -f ./pb/strings/strings.pb.go
	rm -f ./pb/strings/strings_grpc.pb.go
	rm -f ./pb/strings/strings.pb.goclay.go
	rm -f ./strings/strings.go
	rm -f ./strings/to_lower.go
	rm -f ./strings/to_lower_test.go
	rm -f ./strings/to_upper.go
	rm -f ./strings/to_upper_test.go
	rm -f ./strings/strings2.go
	rm -f ./strings/strings2_test.go
	rm -f ./strings/to_lower2.go
	rm -f ./strings/to_lower2_test.go
	rm -f ./strings/to_upper2.go
	rm -f ./strings/to_upper2_test.go
	rm -f main

protoc: .protoc_pb_strings

build: .build

test: pwd clean protoc build
	go test -v ./...
