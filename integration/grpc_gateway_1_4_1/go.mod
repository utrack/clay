module github.com/utrack/clay/integration/grpc_gateway_1_4_1

require (
	github.com/go-chi/chi v0.0.0-20180827203443-b5294d106738
	github.com/go-openapi/spec v0.0.0-20180825180323-f1468acb3b29
	github.com/golang/protobuf v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.5.0
	github.com/pkg/errors v0.8.0
	github.com/utrack/clay/v2 v2.2.5
	golang.org/x/net v0.0.0-20181005035420-146acd28ed58
	google.golang.org/genproto v0.0.0-20180918203901-c3f76f3b92d1
	google.golang.org/grpc v1.16.0
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

replace github.com/utrack/clay/v2 => ../..

replace github.com/grpc-ecosystem/grpc-gateway => github.com/grpc-ecosystem/grpc-gateway v1.4.1
