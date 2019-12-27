module github.com/utrack/clay/integration/grpc_gateway_1_4_1

require (
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-chi/chi v3.3.3+incompatible
	github.com/go-openapi/spec v0.0.0-20180825180323-f1468acb3b29
	github.com/golang/lint v0.0.0-20180702182130-06c8688daad7 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/grpc-gateway v1.12.1
	github.com/pkg/errors v0.8.1
	github.com/utrack/clay/v2 v2.2.5
	golang.org/x/net v0.0.0-20191002035440-2ec189313ef0
	google.golang.org/genproto v0.0.0-20190927181202-20e1ac93f88c
	google.golang.org/grpc v1.24.0
)

replace github.com/utrack/clay/v2 => ../..

replace github.com/grpc-ecosystem/grpc-gateway => github.com/grpc-ecosystem/grpc-gateway v1.4.1

go 1.13
