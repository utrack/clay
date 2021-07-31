include ../env.mk

pwd:
	@pwd

clean:
	rm -f ./pb/strings/strings.pb.go
	rm -f ./pb/strings/strings.pb.goclay.go
	rm -f ./pb/strings/strings_grpc.pb.go
	rm -f ./pb/strings-api/strings_api.pb.go
	rm -f ./pb/strings-api/strings_api.pb.goclay.go
	rm -f ./pb/strings-api/strings_api_grpc.pb.go
	rm -f ./strings/strings.go
	rm -f ./strings/to_lower.go
	rm -f ./strings/to_lower_test.go
	rm -f ./strings/to_upper.go
	rm -f ./strings/to_upper_test.go
	rm -f ./strings-api/strings_api.go
	rm -f ./strings-api/to_lower.go
	rm -f ./strings-api/to_upper.go
	rm -f ./strings-api/to_lower_test.go
	rm -f ./strings-api/to_upper_test.go
	rm -f main

protoc: protoc-build
	protoc \
		--plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --goclay_out=. --goclay_opt=impl_type_name_tmpl=Server,impl_service_sub_dir=true,force=false,impl=true,impl_path=../..,paths=source_relative \
		--plugin=protoc-gen-go=$(GEN_GO_BIN) --go_out=. --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(GEN_GO_GRPC_BIN) --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		-I/usr/local/include:${THIRD_PARTY_PROTO_PATH}:. \
		pb/strings/strings.proto
	protoc \
		--plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --goclay_out=. --goclay_opt=impl_type_name_tmpl=Server,impl_service_sub_dir=true,force=false,impl=true,impl_path=../..,paths=source_relative \
		--plugin=protoc-gen-go=$(GEN_GO_BIN) --go_out=. --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(GEN_GO_GRPC_BIN) --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		-I/usr/local/include:${THIRD_PARTY_PROTO_PATH}:. \
		pb/strings-api/strings_api.proto

build: .build

test: pwd clean protoc build
	go test -v ./...
