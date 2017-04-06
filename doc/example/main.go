package main

import (
	"errors"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/pressly/chi"
	"github.com/rakyll/statik/fs"

	pb "github.com/utrack/clay/doc/example/pb"
	"github.com/utrack/clay/transport"
	"golang.org/x/net/context"

	// We're using statik-compiled files of Swagger UI
	// for the sake of example.
	_ "github.com/utrack/clay/static/statik"
)

// SumImpl is an implementation of SummatorService.
type SumImpl struct{}

// Sum implements SummatorServer.Sum.
func (s *SumImpl) Sum(ctx context.Context, r *pb.SumRequest) (*pb.SumResponse, error) {
	if r.GetA() == 0 {
		return nil, errors.New("a is zero")
	}

	sum := r.GetA() + r.GetB()
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
	srv := transport.NewServer(12345, transport.OptsHTTPMux(hmux))
	err = srv.Run(impl)
	if err != nil {
		logrus.Fatal(err)
	}
}
