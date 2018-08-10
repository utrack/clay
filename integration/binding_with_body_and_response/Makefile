include ../env.mk

$(shell vgo get github.com/googleapis/googleapis)
GOOGLEAPIS_PKG:=$(shell vgo list -m all | grep github.com/googleapis/googleapis | awk '{print ($$4 != "" ? $$4 : $$1)}')
GOOGLEAPIS_VERSION:=$(shell vgo list -m all | grep github.com/googleapis/googleapis | awk '{print ($$5 != "" ? $$5 : $$2)}')
GOOGLEAPIS_PATH:=${FIRST_GOPATH}/pkg/mod/${GOOGLEAPIS_PKG}@${GOOGLEAPIS_VERSION}

pwd:
	@pwd

clean:
	rm -f ./pb/strings.pb.go
	rm -f ./pb/strings.pb.goclay.go
	rm -f main

protoc: protoc-build
	protoc --plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --plugin=protoc-gen-gofast=$(GEN_GOFAST_BIN) -I/usr/local/include:${GOOGLEAPIS_PATH}:. --gofast_out=plugins=grpc:. --goclay_out=impl=true,impl_path=../strings:. pb/strings.proto

build:
	vgo mod vendor
	go build -o main main.go
	vgo build -o main .

test: pwd clean protoc build
	vgo test -v .