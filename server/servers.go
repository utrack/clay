package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

type serverSet struct {
	http    chi.Router
	grpc    *grpc.Server
	httpSrv *http.Server
}

func newServerSet(listeners *listenerSet, opts *serverOpts) *serverSet {
	r := chi.NewMux()
	if len(opts.HTTPMiddlewares) > 0 {
		r.Use(opts.HTTPMiddlewares...)
	}
	r.Mount("/", opts.HTTPMux)

	srv := &serverSet{
		http: r,
	}
	if opts.GRPCServer != nil {
		srv.grpc = opts.GRPCServer
	} else {
		srv.grpc = grpc.NewServer(opts.GRPCOpts...)
	}
	if opts.HTTPServer != nil {
		srv.httpSrv = opts.HTTPServer
	} else {
		srv.httpSrv = &http.Server{}
	}
	return srv
}
