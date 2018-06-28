module github.com/utrack/clay/server

replace github.com/utrack/clay/transport/v2 v2.0.0 => ../transport

require (
	github.com/Sirupsen/logrus v1.0.5
	github.com/go-chi/chi v0.0.0-20171222161133-e83ac2304db3
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/onsi/ginkgo v1.5.0
	github.com/onsi/gomega v1.4.0
	github.com/pressly/chi v1.0.0
	github.com/soheilhy/cmux v0.1.4
	github.com/utrack/clay/transport/v2 v2.0.0
	gopkg.in/airbrake/gobrake.v2 v2.0.9
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2
)
