module github.com/utrack/clay/doc/example

go 1.14

replace github.com/utrack/clay/v3 => ../../

require (
	github.com/go-chi/chi v3.3.4+incompatible
	github.com/go-openapi/spec v0.0.0-20180415031709-bcff419492ee
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.3
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.1
	github.com/sirupsen/logrus v1.4.2
	github.com/utrack/clay/v3 v3.0.1
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1
	google.golang.org/grpc v1.56.3
	google.golang.org/protobuf v1.30.0
)
