package transport

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/utrack/clay/transport/middlewares/mwhttp"

	"google.golang.org/grpc"
)

// Option is an optional setting applied to the Server.
type Option func(*serverOpts)

type serverOpts struct {
	RPCPort int
	// If HTTPPort is the same then muxing listener is created.
	HTTPPort int
	HTTPMux  chi.Router

	HTTPMiddlewares []func(http.Handler) http.Handler

	GRPCOpts []grpc.ServerOption
}

func defaultServerOpts(mainPort int) *serverOpts {
	return &serverOpts{
		RPCPort:  mainPort,
		HTTPPort: mainPort,
		HTTPMux:  chi.NewMux(),
	}
}

// NewServer creates a Server listening on the rpcPort.
// Pass additional Options to mutate its behaviour.
// By default, HTTP JSON handler and gRPC are listening on the same
// port, admin port is p+2 and profile port is p+4.
func NewServer(rpcPort int, opts ...Option) *Server {
	serverOpts := defaultServerOpts(rpcPort)
	for _, opt := range opts {
		opt(serverOpts)
	}
	return &Server{opts: serverOpts}
}

// OptsGRPCOpts sets gRPC server options.
func OptsGRPCOpts(opts []grpc.ServerOption) Option {
	return func(o *serverOpts) {
		o.GRPCOpts = opts
	}
}

// OptsHTTPPort sets HTTP RPC port to listen on.
// Set same port as main to use single port.
func OptsHTTPPort(port int) Option {
	return func(o *serverOpts) {
		o.HTTPPort = port
	}
}

// OptsHTTPMiddlewares sets up HTTP middlewares to work with.
func OptsHTTPMiddlewares(mws ...mwhttp.Middleware) Option {
	mwGeneric := make([]func(http.Handler) http.Handler, 0, len(mws))
	for _, mw := range mws {
		mwGeneric = append(mwGeneric, mw)
	}
	return func(o *serverOpts) {
		o.HTTPMiddlewares = mwGeneric
	}
}

// OptsHTTPMux sets existing HTTP muxer to use instead of creating new one.
func OptsHTTPMux(mux *chi.Mux) Option {
	return func(o *serverOpts) {
		o.HTTPMux = mux
	}
}
