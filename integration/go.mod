module github.com/utrack/clay/integration

replace github.com/utrack/clay/v2 => ../

require (
	github.com/go-chi/chi v0.0.0-20171222161133-e83ac2304db3
	github.com/go-openapi/spec v0.0.0-20180415031709-bcff419492ee
	github.com/gogo/protobuf v1.0.0
	github.com/golang/protobuf v1.1.0
	github.com/grpc-ecosystem/grpc-gateway v1.4.1
	github.com/jmoiron/jsonq v0.0.0-20150511023944-e874b168d07e
	github.com/pkg/errors v0.8.0
	github.com/stretchr/testify v1.2.2
	github.com/utrack/clay/v2 v2.1.0
	github.com/utrack/grpc-gateway v0.0.0-20180624180524-fabb4258c392
	golang.org/x/net v0.0.0-20180629035331-4cb1c02c05b0
	google.golang.org/genproto v0.0.0-20180808183934-383e8b2c3b9e
	google.golang.org/grpc v1.13.0
)
