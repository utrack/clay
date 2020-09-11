include ../env.mk

pwd:
	@pwd

clean:
	rm -f ./pb/strings.pb.go
	rm -f ./pb/strings.pb.goclay.go
	rm -rf ./vendor.pb/
	rm -f main

protoc: protoc-build
	mkdir -p ./vendor.pb/github.com/gogo/protobuf/gogoproto && \
	cp ${GPRC_GOGO_PROTO_PATH}/gogo.proto ./vendor.pb/github.com/gogo/protobuf/gogoproto/gogo.proto && \
	chmod u+w ./vendor.pb/github.com/gogo/protobuf/gogoproto/gogo.proto

	protoc --plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --plugin=protoc-gen-gofast=$(GEN_GOFAST_BIN) \
		-I/usr/local/include:${GRPC_GATEWAY_PATH}/third_party/googleapis:./vendor.pb:. \
		--gofast_out=plugins=grpc:. \
		--goclay_out=impl=true,impl_path=../../strings:. pb/strings/strings.proto

build:
	go build -o main main.go


test: pwd clean protoc build
	go test -v ./strings
