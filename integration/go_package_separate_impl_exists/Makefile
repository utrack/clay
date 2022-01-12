include ../env.mk

pwd:
	@pwd

clean:
	rm -f ./pkg/strings/strings.pb.go
	rm -f ./pkg/strings/strings_grpc.pb.go
	rm -f ./pkg/strings/strings.pb.goclay.go
	rm -f ./internal/app/strings/to_upper.go
	rm -f ./internal/app/strings/to_upper_test.go

protoc: protoc-build
	protoc \
		--plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --goclay_out=pkg/strings --goclay_opt=impl=true,impl_path=../../internal/app/strings,paths=source_relative \
		--plugin=protoc-gen-go=$(GEN_GO_BIN) --go_out=pkg/strings --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(GEN_GO_GRPC_BIN) --go-grpc_out=pkg/strings --go-grpc_opt=paths=source_relative \
		-I/usr/local/include:${THIRD_PARTY_PROTO_PATH}:. \
		strings.proto

build: .build


test: pwd clean protoc build
	go test ./...
