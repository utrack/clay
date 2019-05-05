package server

import (
	"net/http"
	"os"

	"github.com/utrack/clay/v2/server/middlewares/mwhttp"
	"github.com/utrack/clay/v2/transport"

	"github.com/go-chi/chi"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

// Option is an optional setting applied to the Server.
type Option func(*serverOpts)

type serverOpts struct {
	RPCPort int
	// If HTTPPort is the same then muxing listener is created.
	HTTPPort int
	HTTPMux  transport.Router
	// HTTPServer holds pointer to custom Server instance
	HTTPServer *http.Server
	// HTTPGracefullFunc for gracefull shutdown http server
	HTTPGracefullFunc func(sc chan os.Signal) func() error

	HTTPMiddlewares []func(http.Handler) http.Handler

	GRPCOpts             []grpc.ServerOption
	GRPCUnaryInterceptor grpc.UnaryServerInterceptor
}

func defaultServerOpts(mainPort int) *serverOpts {
	return &serverOpts{
		RPCPort:  mainPort,
		HTTPPort: mainPort,
		HTTPMux:  chi.NewMux(),
	}
}

// WithGRPCOpts sets gRPC server options.
func WithGRPCOpts(opts []grpc.ServerOption) Option {
	return func(o *serverOpts) {
		o.GRPCOpts = append(o.GRPCOpts, opts...)
	}
}

// WithHTTPPort sets HTTP RPC port to listen on.
// Set same port as main to use single port.
func WithHTTPPort(port int) Option {
	return func(o *serverOpts) {
		o.HTTPPort = port
	}
}

// WithHTTPMiddlewares sets up HTTP middlewares to work with.
func WithHTTPMiddlewares(mws ...mwhttp.Middleware) Option {
	mwGeneric := make([]func(http.Handler) http.Handler, 0, len(mws))
	for _, mw := range mws {
		mwGeneric = append(mwGeneric, mw)
	}
	return func(o *serverOpts) {
		o.HTTPMiddlewares = mwGeneric
	}
}

// WithGRPCUnaryMiddlewares sets up unary middlewares for gRPC server.
func WithGRPCUnaryMiddlewares(mws ...grpc.UnaryServerInterceptor) Option {
	mw := grpc_middleware.ChainUnaryServer(mws...)
	return func(o *serverOpts) {
		o.GRPCOpts = append(o.GRPCOpts, grpc.UnaryInterceptor(mw))
		o.GRPCUnaryInterceptor = mw
	}
}

// WithGRPCStreamMiddlewares sets up stream middlewares for gRPC server.
func WithGRPCStreamMiddlewares(mws ...grpc.StreamServerInterceptor) Option {
	return func(o *serverOpts) {
		o.GRPCOpts = append(o.GRPCOpts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(mws...)))
	}
}

// WithHTTPMux sets existing HTTP muxer to use instead of creating new one.
func WithHTTPMux(mux *chi.Mux) Option {
	return func(o *serverOpts) {
		o.HTTPMux = mux
	}
}

// WithHTTPRouterMux func sets HTTPMux Router
func WithHTTPRouterMux(mux transport.Router) Option {
	return func(o *serverOpts) {
		o.HTTPMux = mux
	}
}

// WithHTTPServer sets HTTP Server to use insted of the default
func WithHTTPServer(srv *http.Server) Option {
	if srv == nil {
		panic("sent Server pointer is nil")
	}
	return func(o *serverOpts) {
		o.HTTPServer = srv
	}
}

// WithHTTPGracefull applies Gracefull shutdown func to server
func WithHTTPGracefull(fn func(sc chan os.Signal) func() error) Option {
	if fn == nil {
		panic("sent Gracefull func is nil")
	}
	return func(o *serverOpts) {
		o.HTTPGracefullFunc = fn
	}
}
