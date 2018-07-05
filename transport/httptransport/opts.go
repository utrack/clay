package httptransport

import "google.golang.org/grpc"

// DescOptions provides options for a ServiceDesc compiled code.
type DescOptions struct {
	UnaryInterceptor grpc.UnaryServerInterceptor
}

// OptionUnaryInterceptor sets up the gRPC unary interceptor.
type OptionUnaryInterceptor struct {
	Interceptor grpc.UnaryServerInterceptor
}

// Apply implements transport.DescOption.
func (o OptionUnaryInterceptor) Apply(oo DescOptions) {
	if oo.UnaryInterceptor != nil {
		// Chaining can be done via mwitkow/grpc-middleware for example.
		panic("UnaryInterceptor is already applied, can't apply twice")
	}
	oo.UnaryInterceptor = o.Interceptor
}
