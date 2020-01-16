module github.com/utrack/clay/integration

replace github.com/utrack/clay/v2 => ../

require (
	github.com/go-chi/chi v3.3.3+incompatible
	github.com/go-openapi/spec v0.0.0-20180415031709-bcff419492ee
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/googleapis/googleapis v0.0.0-20200115224547-0735b4b09687 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.12.1
	github.com/jmoiron/jsonq v0.0.0-20150511023944-e874b168d07e
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
	github.com/utrack/clay/v2 v2.1.0
	golang.org/x/net v0.0.0-20191002035440-2ec189313ef0
	google.golang.org/genproto v0.0.0-20190927181202-20e1ac93f88c
	google.golang.org/grpc v1.24.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

go 1.13
