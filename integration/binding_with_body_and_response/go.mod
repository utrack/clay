module github.com/utrack/clay/integration/binding_with_body_and_response

replace github.com/grpc-ecosystem/grpc-gateway => github.com/doroginin/grpc-gateway v1.5.0-alpha3

replace github.com/utrack/grpc-gateway => github.com/doroginin/grpc-gateway v1.5.0-alpha3

replace github.com/googleapis/googleapis => github.com/doroginin/googleapis v0.0.0-20180730132820-50a24711b667

replace google.golang.org/genproto => github.com/doroginin/go-genproto v0.0.0-20180730134020-8126e81001e3

replace github.com/utrack/clay/v2 => ../..

require (
	github.com/go-chi/chi v0.0.0-20171222161133-e83ac2304db3
	github.com/go-openapi/spec v0.0.0-20180415031709-bcff419492ee
	github.com/googleapis/googleapis v0.0.0-20180808183217-3b3acb6f44b5 // indirect
	github.com/utrack/clay/v2 v2.1.0
)
