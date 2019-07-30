module github.com/utrack/clay/integration

replace github.com/utrack/clay/v2 => ../

require (
	github.com/go-chi/chi v3.3.3+incompatible
	github.com/go-openapi/spec v0.0.0-20180415031709-bcff419492ee
	github.com/gogo/protobuf v1.2.0
	github.com/golang/protobuf v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.5.0
	github.com/jmoiron/jsonq v0.0.0-20150511023944-e874b168d07e
	github.com/pkg/errors v0.8.0
	github.com/stretchr/testify v1.2.2
	github.com/utrack/clay/v2 v2.1.0
	golang.org/x/net v0.0.0-20181005035420-146acd28ed58
	google.golang.org/genproto v0.0.0-20180918203901-c3f76f3b92d1
	google.golang.org/grpc v1.16.0
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

replace github.com/gogo/protobuf => github.com/dzendmitry/protobuf v1.3.1
