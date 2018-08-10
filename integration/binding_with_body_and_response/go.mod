module github.com/utrack/clay/integration/binding_with_body_and_response

replace github.com/grpc-ecosystem/grpc-gateway => github.com/doroginin/grpc-gateway v1.5.0-alpha4

replace github.com/utrack/grpc-gateway => github.com/doroginin/grpc-gateway v1.5.0-alpha3

replace github.com/utrack/clay/v2 => ../..

require (
	github.com/go-chi/chi v0.0.0-20171222161133-e83ac2304db3
	github.com/go-openapi/spec v0.0.0-20180415031709-bcff419492ee
	github.com/golang/protobuf v1.1.0
	github.com/googleapis/googleapis v0.0.0-20180809173157-2c07771e7948 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.4.1
	github.com/pkg/errors v0.8.0
	github.com/utrack/clay/v2 v2.1.0
	golang.org/x/net v0.0.0-20180629035331-4cb1c02c05b0
	google.golang.org/genproto v0.0.0-20180808183934-383e8b2c3b9e
	google.golang.org/grpc v1.13.0
)
