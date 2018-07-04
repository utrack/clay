package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"

	pb "github.com/utrack/clay/v2/doc/example/pb"
	"github.com/utrack/clay/v2/log"
	"github.com/utrack/clay/v2/transport"
	"github.com/utrack/clay/v2/transport/middlewares/mwgrpc"
	"github.com/utrack/clay/v2/transport/middlewares/mwhttp"
	"github.com/utrack/clay/v2/transport/server"
	"golang.org/x/net/context"

	// We're using statik-compiled files of Swagger UI
	// for the sake of example.
	_ "github.com/utrack/clay/v2/doc/example/static/statik"
)

// SumImpl is an implementation of SummatorService.
type SumImpl struct{}

// Sum implements SummatorServer.Sum.
func (s *SumImpl) Sum(ctx context.Context, r *pb.SumRequest) (*pb.SumResponse, error) {
	if r.GetA() == 0 {
		return nil, errors.New("a is zero")
	}

	if r.GetB().B == 65536 {
		panic(errors.New("we've got a problem"))
	}

	sum := r.GetA() + r.GetB().B
	return &pb.SumResponse{
		Sum: sum,
	}, nil
}

// GetDescription is a simple alias to the ServiceDesc constructor.
// It makes it possible to register the service implementation @ the server.
func (s *SumImpl) GetDescription() transport.ServiceDesc {
	return pb.NewSummatorServiceDesc(s)
}

func main() {
	// Wire up our bundled Swagger UI
	staticFS, err := fs.New()
	if err != nil {
		logrus.Fatal(err)
	}
	hmux := chi.NewRouter()
	hmux.Mount("/", http.FileServer(staticFS))

	impl := &SumImpl{}
	srv := server.NewServer(
		12345,
		// Pass our mux with Swagger UI
		server.WithHTTPMux(hmux),
		// Recover from HTTP panics
		server.WithHTTPMiddlewares(mwhttp.Recover(log.Default), mwhttp.CloseNotifier()),
		// Recover from gRPC panics
		server.WithGRPCUnaryMiddlewares(mwgrpc.UnaryPanicHandler(log.Default)),
	)
	err = srv.Run(impl)
	if err != nil {
		logrus.Fatal(err)
	}
}
