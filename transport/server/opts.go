package server

import (
	"net/http"
	"os"

	"github.com/utrack/clay/v2/server"
	"github.com/utrack/clay/v2/server/middlewares/mwhttp"
	"github.com/utrack/clay/v2/transport"

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

// WithHTTPRouterMux func sets HTTPMux Router
func WithHTTPRouterMux(mux transport.Router) Option {
	return server.WithHTTPRouterMux(mux)
}

// WithHTTPServer sets HTTP Server to use insted of the default
func WithHTTPServer(srv *http.Server) Option {
	return server.WithHTTPServer(srv)
}

// WithHTTPGracefull applies Gracefull shutdown func to server
func WithHTTPGracefull(fn func(sc chan os.Signal) func() error) Option {
	return server.WithHTTPGracefull(fn)
}

// WithHost sets default server host
func WithHost(host string) Option {
	return server.WithHost(host)
}
