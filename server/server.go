package server

import (
	"bytes"
	"io"
	"net/http"

	"github.com/utrack/clay/v2/transport"
	"golang.org/x/sync/errgroup"

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

	// apply gRPC interceptor
	if d, ok := desc.(transport.ConfigurableServiceDesc); ok {
		d.Apply(transport.WithUnaryInterceptor(s.opts.GRPCUnaryInterceptor))
	}

	// Register everything
	desc.RegisterHTTP(s.srv.http)
	desc.RegisterGRPC(s.srv.grpc)

	return s.run()
}

func (s *Server) run() error {
	var g errgroup.Group
	if s.listeners.mainListener != nil {
		g.Go(func() error {
			err := s.listeners.mainListener.Serve()
			return err
		})
	}
	g.Go(func() error {
		s.srv.httpSrv.Handler = s.srv.http
		err := s.srv.httpSrv.Serve(s.listeners.HTTP)
		return err
	})
	g.Go(func() error {
		err := s.srv.grpc.Serve(s.listeners.GRPC)
		return err
	})

	return g.Wait()
}

// Stop stops the server gracefully.
func (s *Server) Stop() {
	// TODO grace HTTP
	s.srv.grpc.GracefulStop()
}
