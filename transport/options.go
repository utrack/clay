package transport

import (
	"github.com/utrack/clay/v2/transport/httptransport"
	"github.com/utrack/clay/v2/transport/swagger"
	"google.golang.org/grpc"
)

// DescOption modifies the ServiceDesc's behaviour.
type DescOption interface {
	Apply(*httptransport.DescOptions)
}

// WithUnaryInterceptor sets up the interceptor for incoming calls.
func WithUnaryInterceptor(i grpc.UnaryServerInterceptor) DescOption {
	return httptransport.OptionUnaryInterceptor{Interceptor: i}
}

// WithSwaggerOptions sets up default Swagger options for the SwaggerDef().
func WithSwaggerOptions(o ...swagger.Option) DescOption {
	return httptransport.OptionSwaggerOpts{Options: o}
}
