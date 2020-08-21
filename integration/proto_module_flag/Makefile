include ../env.mk

pwd:
	@pwd

clean:
	rm -f ./pb/strings/strings.pb.go
	rm -f ./pb/strings/strings.pb.goclay.go
	rm -f ./strings/strings.go
	rm -f ./strings/to_lower.go
	rm -f ./strings/to_upper.go
	rm -f main

protoc: protoc-build
	protoc --plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --plugin=protoc-gen-go=$(GEN_GO_BIN) -Iproto:/usr/local/include:${GRPC_GATEWAY_PATH}/third_party/googleapis --go_out=module=github.com/utrack/clay/integration/proto_module_flag,plugins=grpc:. --goclay_out=module=github.com/utrack/clay/integration/proto_module_flag,impl=true,impl_path=../strings:. proto/strings.proto

build:
	go build -o main main.go


test: pwd clean protoc build
