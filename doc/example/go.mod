module github.com/utrack/clay/doc/example

go 1.14

replace github.com/utrack/clay/v2 => ../../

require (
	github.com/go-chi/chi v3.3.4+incompatible
	github.com/go-openapi/spec v0.0.0-20180415031709-bcff419492ee
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/pkg/errors v0.8.1
	github.com/rakyll/statik v0.1.1
	github.com/sirupsen/logrus v1.4.2
	github.com/utrack/clay/v2 v2.4.7
	google.golang.org/genproto v0.0.0-20210617175327-b9e0b3197ced
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v2 v2.2.3 // indirect
)
