package server

import (
	"net"

	"google.golang.org/grpc"

	"github.com/go-chi/chi"
	"github.com/soheilhy/cmux"
)

// Server is a transport server.
type Server struct {
	opts      *serverOpts
	listeners *listenerSet
	srv       *serverSet
}

type listenerSet struct {
	mainListener cmux.CMux // nil or CMux. If nil - don't listen
	HTTP         net.Listener
	GRPC         net.Listener
}

type serverSet struct {
	http *chiWrapper
	grpc *grpc.Server
}

func getServers(listeners *listenerSet, opts *serverOpts) *serverSet {
	http := chi.NewMux()
	if len(opts.HTTPMiddlewares) > 0 {
		http.Use(opts.HTTPMiddlewares...)
	}
	http.Mount("/", opts.HTTPMux)

	srv := &serverSet{
		grpc: grpc.NewServer(opts.GRPCOpts...),
		http: &chiWrapper{Router: http},
	}
	return srv
}
