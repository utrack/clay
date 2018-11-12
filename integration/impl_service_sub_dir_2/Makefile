include ../env.mk

pwd:
	@pwd

clean:
	rm -f ./pb/strings/strings.pb.go
	rm -f ./pb/strings/strings.pb.goclay.go
	rm -f ./pb/strings2/strings2.pb.go
	rm -f ./pb/strings2/strings2.pb.goclay.go
	rm -f ./strings/strings.go
	rm -f ./strings/strings.to_lower.go
	rm -f ./strings/strings.to_upper.go
	rm -f ./strings2/strings2.go
	rm -f ./strings2/strings2.to_lower.go
	rm -f ./strings2/strings2.to_upper.go
	rm -f main

protoc: protoc-build
	protoc --plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --plugin=protoc-gen-gofast=$(GEN_GOFAST_BIN) -I/usr/local/include:${GRPC_GATEWAY_PATH}/third_party/googleapis:. --gofast_out=plugins=grpc:. --goclay_out=impl_type_name_tmpl=Server,impl_service_sub_dir=true,impl=true,impl_path=../..:. pb/strings/strings.proto
	protoc --plugin=protoc-gen-goclay=$(GEN_CLAY_BIN) --plugin=protoc-gen-gofast=$(GEN_GOFAST_BIN) -I/usr/local/include:${GRPC_GATEWAY_PATH}/third_party/googleapis:. --gofast_out=plugins=grpc:. --goclay_out=impl_type_name_tmpl=Server,impl_service_sub_dir=true,impl=true,impl_path=../..:. pb/strings2/strings2.proto

build:
	go build -o main main.go
	vgo build -o main .

test: pwd clean protoc build
