module github.com/utrack/clay/integration

require (
	github.com/go-chi/chi v3.3.4+incompatible
	github.com/go-openapi/spec v0.0.0-20180415031709-bcff419492ee
	github.com/gogo/protobuf v1.3.2
	github.com/google/go-cmp v0.5.9
	github.com/jmoiron/jsonq v0.0.0-20150511023944-e874b168d07e
	github.com/stretchr/testify v1.8.3
	github.com/utrack/clay/v3 v3.0.0
	google.golang.org/grpc v1.56.3
	google.golang.org/protobuf v1.30.0
)

go 1.13

replace github.com/utrack/clay/v3 => ../
