FIRST_GOPATH:=$(firstword $(subst :, ,$(GOPATH)))

clean:
	find . -regex "\.\/.*\/.*\.go" -exec rm {} +
	rm -f main

protoc:
	protoc -I/usr/local/include:. --go_out=plugins=grpc:. --goclay_out=impl=true,impl_path=../strings:. pb/strings.proto

build:
	go build -o main main.go && echo

test: clean protoc build