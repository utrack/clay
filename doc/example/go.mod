module github.com/utrack/clay/doc/example

go 1.14

replace github.com/utrack/clay/v2 => ../../

require (
	github.com/go-chi/chi v3.3.4+incompatible
	github.com/go-openapi/spec v0.0.0-20180415031709-bcff419492ee
	github.com/golang/lint v0.0.0-20180702182130-06c8688daad7 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/grpc-gateway v1.14.2
	github.com/pkg/errors v0.8.1
	github.com/rakyll/statik v0.1.1
	github.com/sirupsen/logrus v1.4.2
	github.com/utrack/clay/v2 v2.4.7
	golang.org/x/net v0.0.0-20191002035440-2ec189313ef0
	google.golang.org/genproto v0.0.0-20190927181202-20e1ac93f88c
	google.golang.org/grpc v1.27.1
)
