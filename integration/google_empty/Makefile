include ../env.mk

PKGMAP:=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,$\
        Mgoogle/protobuf/api.proto=github.com/gogo/protobuf/types,$\
        Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/types,$\
        Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,$\
        Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,$\
        Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,$\
        Mgoogle/protobuf/source_context.proto=github.com/gogo/protobuf/types,$\
        Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,$\
        Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,$\
        Mgoogle/protobuf/type.proto=github.com/gogo/protobuf/types,$\
        Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types

pwd:
	@pwd

clean:
	rm -f ./pb/strings/strings.pb.go
	rm -f ./pb/strings/strings.pb.goclay.go
	rm -f ./strings/strings.go
	rm -f ./strings/to_lower.go
	rm -f ./strings/to_upper.go
	rm -f ./strings/empty_response.go
	rm -f ./strings/without_desc.go
	rm -f ./strings/without_desc_and_option.go
	rm -f main

protoc: protoc-build
	protoc --plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --plugin=protoc-gen-gofast=$(GEN_GOFAST_BIN) -I/usr/local/include:${GRPC_GATEWAY_PATH}/third_party/googleapis:. --gofast_out=$(PKGMAP),plugins=grpc:../../../../.. --goclay_out=$(PKGMAP),impl=true,impl_path=../../strings,impl_type_name_tmpl=Implementation:../../../../.. pb/strings/strings.proto

build:
	go build -o main main.go
	vgo build -o main .

test: pwd clean protoc build
