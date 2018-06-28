package server

import (
	"bytes"
	"io"
	"net/http"

	"github.com/utrack/clay/transport/v2"

	"github.com/pkg/errors"
)

// Server is a transport server.
type Server struct {
	opts      *serverOpts
	listeners *listenerSet
	srv       *serverSet
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

// Run starts processing requests to the service.
// It blocks indefinitely, run asynchronously to do anything after that.
func (s *Server) Run(svc transport.Service) error {
	desc := svc.GetDescription()

	var err error
	s.listeners, err = newListenerSet(s.opts)
	if err != nil {
		return errors.Wrap(err, "couldn't create listeners")
	}

	s.srv = newServerSet(s.listeners, s.opts)
	// Inject static Swagger as root handler
	s.srv.http.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, bytes.NewReader(desc.SwaggerDef()))
	})

	// Register everything
	desc.RegisterHTTP(s.srv.http)
	desc.RegisterGRPC(s.srv.grpc)

	return s.run()
}

func (s *Server) run() error {
	errChan := make(chan error, 5)

	if s.listeners.mainListener != nil {
		go func() {
			err := s.listeners.mainListener.Serve()
			errChan <- err
		}()
	}
	go func() {
		err := http.Serve(s.listeners.HTTP, s.srv.http)
		errChan <- err
	}()
	go func() {
		err := s.srv.grpc.Serve(s.listeners.GRPC)
		errChan <- err
	}()

	return <-errChan
}

// Stop stops the server gracefully.
func (s *Server) Stop() {
	// TODO grace HTTP
	s.srv.grpc.GracefulStop()
}
