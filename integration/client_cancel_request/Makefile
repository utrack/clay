include ../env.mk

pwd:
	@pwd

clean:
	rm -f ./pkg/strings/strings.pb.go
	rm -f ./pkg/strings/strings_grpc.pb.go
	rm -f ./pkg/strings/strings.pb.goclay.go
	rm -f ./app/strings/strings.go
	rm -f main

protoc: protoc-build
	cd pkg/strings && \
	protoc \
		--plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --goclay_out=impl=true,impl_path=../../app/strings:. \
		--plugin=protoc-gen-go=$(GEN_GO_BIN) --go_out=. --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(GEN_GO_GRPC_BIN) --go-grpc_out=. --go-grpc_opt paths=source_relative \
		-I../../api:/usr/local/include:${THIRD_PARTY_PROTO_PATH}:. \
		strings.proto

build: .build

test: pwd clean protoc build
	go test -v ./...
