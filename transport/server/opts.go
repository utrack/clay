package server

import (
	"github.com/utrack/clay/server"
	"github.com/utrack/clay/server/middlewares/mwhttp"
	"github.com/utrack/clay/transport/v2"

	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

// Option is an optional setting applied to the Server.
type Option = server.Option

// WithGRPCOpts sets gRPC server options.
func WithGRPCOpts(opts []grpc.ServerOption) Option {
	return server.WithGRPCOpts(opts)
}

// WithHTTPPort sets HTTP RPC port to listen on.
// Set same port as main to use single port.
func WithHTTPPort(port int) Option {
	return server.WithHTTPPort(port)
}

// WithHTTPMiddlewares sets up HTTP middlewares to work with.
func WithHTTPMiddlewares(mws ...mwhttp.Middleware) Option {
	return server.WithHTTPMiddlewares(mws...)
}

// WithGRPCUnaryMiddlewares sets up unary middlewares for gRPC server.
func WithGRPCUnaryMiddlewares(mws ...grpc.UnaryServerInterceptor) Option {
	return server.WithGRPCUnaryMiddlewares(mws...)
}

// WithGRPCStreamMiddlewares sets up stream middlewares for gRPC server.
func WithGRPCStreamMiddlewares(mws ...grpc.StreamServerInterceptor) Option {
	return server.WithGRPCStreamMiddlewares(mws...)
}

// WithHTTPMux sets existing HTTP muxer to use instead of creating new one.
func WithHTTPMux(mux *chi.Mux) Option {
	return server.WithHTTPMux(mux)
}

func WithHTTPRouterMux(mux transport.Router) Option {
	return server.WithHTTPRouterMux(mux)
}
