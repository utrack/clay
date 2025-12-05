module github.com/utrack/clay/doc/example

go 1.23.0

replace github.com/utrack/clay/v3 => ../../

require (
	github.com/go-chi/chi v3.3.4+incompatible
	github.com/go-openapi/spec v0.0.0-20180415031709-bcff419492ee
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.3
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.1
	github.com/sirupsen/logrus v1.8.3
	github.com/utrack/clay/v3 v3.0.1
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1
	google.golang.org/grpc v1.56.3
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/PuerkitoBio/purell v1.1.0 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/go-openapi/jsonpointer v0.17.2 // indirect
	github.com/go-openapi/jsonreference v0.17.2 // indirect
	github.com/go-openapi/swag v0.17.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/mailru/easyjson v0.0.0-20180823135443-60711f1a8329 // indirect
	github.com/peterbourgon/mergemap v0.0.0-20130613134717-e21c03b7a721 // indirect
	github.com/soheilhy/cmux v0.1.4 // indirect
	golang.org/x/net v0.36.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	gopkg.in/yaml.v2 v2.2.3 // indirect
)
