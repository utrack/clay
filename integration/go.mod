module github.com/utrack/clay/integration

require (
	github.com/go-chi/chi v3.3.4+incompatible
	github.com/go-openapi/spec v0.0.0-20180415031709-bcff419492ee
	github.com/gogo/protobuf v1.3.2
	github.com/google/go-cmp v0.5.9
	github.com/googleapis/googleapis v0.0.0-20240509000043-25a1a57957d9 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.3
	github.com/jmoiron/jsonq v0.0.0-20150511023944-e874b168d07e
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.1
	github.com/utrack/clay/v3 v3.0.0
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
)

go 1.13

replace github.com/utrack/clay/v3 => ../
