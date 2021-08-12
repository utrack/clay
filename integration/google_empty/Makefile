include ../env.mk

PKGMAP:=M"google/protobuf/any.proto;types"=github.com/gogo/protobuf/types,$\
        M"google/protobuf/api.proto;types"=github.com/gogo/protobuf/types,$\
        M"google/protobuf/descriptor.proto;types"=github.com/gogo/protobuf/types,$\
        M"google/protobuf/duration.proto;types"=github.com/gogo/protobuf/types,$\
        M"google/protobuf/empty.proto;types"=github.com/gogo/protobuf/types,$\
        M"google/protobuf/field_mask.proto;types"=github.com/gogo/protobuf/types,$\
        M"google/protobuf/source_context.proto;types"=github.com/gogo/protobuf/types,$\
        M"google/protobuf/struct.proto;types"=github.com/gogo/protobuf/types,$\
        M"google/protobuf/timestamp.proto;types"=github.com/gogo/protobuf/types,$\
        M"google/protobuf/type.proto;types"=github.com/gogo/protobuf/types,$\
        M"google/protobuf/wrappers.proto;types"=github.com/gogo/protobuf/types

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
	rm -f ./strings/empty_response.go
	rm -f ./strings/empty_response_test.go
	rm -f ./strings/without_desc.go
	rm -f ./strings/without_desc_test.go
	rm -f ./strings/without_desc_and_option.go
	rm -f ./strings/without_desc_and_option_test.go
	rm -f main

protoc: protoc-build
	protoc \
		--plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --goclay_out=. --goclay_opt=$(PKGMAP),impl=true,impl_type_name_tmpl=Implementation,impl_path=../../strings,paths=source_relative \
		--plugin=protoc-gen-go=$(GEN_GO_BIN) --go_out=. --go_opt=$(PKGMAP),paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(GEN_GO_GRPC_BIN) --go-grpc_out=. --go-grpc_opt=$(PKGMAP),paths=source_relative \
		-I/usr/local/include:${THIRD_PARTY_PROTO_PATH}:. \
		pb/strings/strings.proto

build: .build

test: pwd clean protoc build
	go test -v ./...
