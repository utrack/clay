package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"
	sum "github.com/utrack/clay/doc/example/implementation"

	"github.com/utrack/clay/v3/log"
	"github.com/utrack/clay/v3/transport/middlewares/mwgrpc"
	"github.com/utrack/clay/v3/transport/server"

	// We're using statik-compiled files of Swagger UI
	// for the sake of example.
	_ "github.com/utrack/clay/doc/example/static/statik"
)

func main() {
	// Wire up our bundled Swagger UI
	staticFS, err := fs.New()
	if err != nil {
		logrus.Fatal(err)
	}
	hmux := chi.NewRouter()
	hmux.Mount("/", http.FileServer(staticFS))

	impl := sum.NewSummator()
	srv := server.NewServer(
		12345,
		// Pass our mux with Swagger UI
		server.WithHTTPMux(hmux),
		// Recover from both HTTP and gRPC panics and use our own middleware
		server.WithGRPCUnaryMiddlewares(mwgrpc.UnaryPanicHandler(log.Default)),
	)
	err = srv.Run(impl)
	if err != nil {
		logrus.Fatal(err)
	}
}
