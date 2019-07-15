package httptransport

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/utrack/clay/v2/transport/swagger"
	"google.golang.org/grpc"
)

// DescOptions provides options for a ServiceDesc compiled code.
type DescOptions struct {
	UnaryInterceptor   grpc.UnaryServerInterceptor
	SwaggerDefaultOpts []swagger.Option
}

// OptionUnaryInterceptor sets up the gRPC unary interceptor.
type OptionUnaryInterceptor struct {
	Interceptor grpc.UnaryServerInterceptor
}

// Apply implements transport.DescOption.
func (o OptionUnaryInterceptor) Apply(oo *DescOptions) {
	if oo.UnaryInterceptor != nil {
		oo.UnaryInterceptor = grpc_middleware.ChainUnaryServer(
			oo.UnaryInterceptor,
			o.Interceptor,
		)
		return
	}
	oo.UnaryInterceptor = o.Interceptor
}

// OptionSwaggerOpts sets up default options for the SwaggerDef().
type OptionSwaggerOpts struct {
	Options []swagger.Option
}

// Apply implements transport.DescOption.
func (o OptionSwaggerOpts) Apply(oo *DescOptions) {
	oo.SwaggerDefaultOpts = append(oo.SwaggerDefaultOpts, o.Options...)
}
