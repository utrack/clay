package server

import (
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type serverSet struct {
	http chi.Router
	grpc *grpc.Server
}

func newServerSet(listeners *listenerSet, opts *serverOpts) *serverSet {
	http := chi.NewMux()
	if len(opts.HTTPMiddlewares) > 0 {
		http.Use(opts.HTTPMiddlewares...)
	}
	http.Mount("/", opts.HTTPMux)

	grpcServer := grpc.NewServer(opts.GRPCOpts...)
	if opts.EnableReflection {
		reflection.Register(grpcServer)
	}

	srv := &serverSet{
		grpc: grpcServer,
		http: http,
	}
	return srv
}
