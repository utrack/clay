include ../env.mk

pwd:
	@pwd

clean:
	rm -f ./pkg/strings/strings.pb.go
	rm -f ./pkg/strings/strings.pb.goclay.go
	rm -f ./app/strings/strings.go
	rm -f main

protoc: protoc-build
	cd pkg/strings && protoc --plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --plugin=protoc-gen-gofast=$(GEN_GOFAST_BIN) -I../../api:/usr/local/include:${GRPC_GATEWAY_PATH}/third_party/googleapis --gofast_out=plugins=grpc:. --goclay_out=impl=true,impl_path=../../app/strings:. ../../api/strings.proto

build:
	go build -o main main.go


test: pwd clean protoc build
	vgo test -v .
