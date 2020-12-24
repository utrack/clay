module github.com/utrack/clay/integration

replace github.com/utrack/clay/v2 => ../

require (
	github.com/go-chi/chi v3.3.4+incompatible
	github.com/go-openapi/spec v0.0.0-20180415031709-bcff419492ee
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/googleapis/googleapis v0.0.0-20201222224849-8a6f4d9acb16 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.2
	github.com/jmoiron/jsonq v0.0.0-20150511023944-e874b168d07e
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
	github.com/utrack/clay/v2 v2.1.0
	google.golang.org/genproto v0.0.0-20190927181202-20e1ac93f88c
	google.golang.org/grpc v1.27.1
)

go 1.13
