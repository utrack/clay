include ../env.mk

pwd:
	@pwd

clean:
	find . -regex "\.\/.*\/.*\.go" -exec rm {} +
	rm -f main

protoc: protoc-build
	protoc --plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --plugin=protoc-gen-gofast=$(GEN_GOFAST_BIN) -I/usr/local/include:${GRPC_GATEWAY_PATH}/third_party/googleapis:. --gofast_out=plugins=grpc:. --goclay_out=impl=true,impl_path=../../strings:. strings.proto

build:
	go build -o main main.go
	vgo build -o main .

test: pwd clean protoc build
