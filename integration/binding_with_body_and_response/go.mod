module github.com/utrack/clay/integration/binding_with_body_and_response

replace github.com/grpc-ecosystem/grpc-gateway => github.com/doroginin/grpc-gateway v1.5.0-alpha

replace github.com/utrack/grpc-gateway => github.com/doroginin/grpc-gateway v1.5.0-alpha-wo-vgo

replace github.com/utrack/clay/v2 => ../..

require (
	github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc // indirect
	github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf // indirect
	github.com/go-chi/chi v0.0.0-20171222161133-e83ac2304db3
	github.com/go-openapi/spec v0.0.0-20180415031709-bcff419492ee
	github.com/gogo/protobuf v1.0.0
	github.com/golang/protobuf v1.1.0
	github.com/grpc-ecosystem/grpc-gateway v1.4.1
	github.com/jmoiron/jsonq v0.0.0-20150511023944-e874b168d07e
	github.com/pkg/errors v0.8.0
	github.com/prometheus/common v0.0.0-20180518154759-7600349dcfe1
	github.com/stretchr/testify v1.2.2
	github.com/utrack/clay/v2 v2.1.0
	golang.org/x/net v0.0.0-20180629035331-4cb1c02c05b0
	google.golang.org/genproto v0.0.0-20180627194029-ff3583edef7d
	google.golang.org/grpc v1.13.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6 // indirect
)
