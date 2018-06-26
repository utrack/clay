module github.com/utrack/clay/server

replace github.com/utrack/clay/transport/v2 v2.0.0 => ../transport

require (
	github.com/Sirupsen/logrus v1.0.5
	github.com/go-chi/chi v0.0.0-20171222161133-e83ac2304db3
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/pressly/chi v1.0.0
	github.com/soheilhy/cmux v0.1.4
	github.com/utrack/clay/transport/v2 v2.0.0
)
