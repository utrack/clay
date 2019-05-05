package server

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	if s.opts.GracefullFunc != nil {
		fn := s.opts.GracefullFunc(sigint)
		g.Go(fn)
	} else {
		g.Go(func() error {
			<-sigint
			s.Stop()
			return nil
		})
	}
	g.Go(func() error {
		if s.opts.HTTPServer != nil {
			s.opts.HTTPServer.Handler = s.srv.http
			err := s.opts.HTTPServer.Serve(s.listeners.HTTP)
			return err
		}
		err := http.Serve(s.listeners.HTTP, s.srv.http)
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
	var g errgroup.Group
	if s.opts.HTTPServer != nil {
		g.Go(func() error {
			return s.opts.HTTPServer.Shutdown(context.Background())
		})
	}
	g.Go(func() error {
		s.srv.grpc.GracefulStop()
		return nil
	})
	g.Wait()
}

// StopWithTimeout stops the server gracefully.
func (s *Server) StopWithTimeout(timeout time.Duration) {
	var g errgroup.Group
	if s.opts.HTTPServer != nil {
		g.Go(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			return s.opts.HTTPServer.Shutdown(ctx)
		})
	}
	g.Go(func() error {
		s.srv.grpc.GracefulStop()
		return nil
	})
	g.Wait()
}
