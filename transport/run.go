package transport

import (
	"bytes"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// Run starts processing requests to the service.
// It blocks indefinitely, run asynchronously to do anything after that.
func (s *Server) Run(svc Service) error {
	desc := svc.GetDescription()

	var err error
	s.listeners, err = getListeners(s.opts)
	if err != nil {
		return errors.Wrap(err, "couldn't create listeners")
	}

	s.srv = getServers(s.listeners, s.opts)
	// Inject static Swagger as root handler
	s.srv.http.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, bytes.NewReader(desc.SwaggerDef()))
	})

	// Register everything
	desc.RegisterHTTP(s.srv.http)
	desc.RegisterGRPC(s.srv.grpc)

	return s.run(s.listeners, s.srv)
}

func (Server) run(l *listenerSet, s *serverSet) error {
	errChan := make(chan error, 5)

	if l.mainListener != nil {
		go func() {
			err := l.mainListener.Serve()
			errChan <- err
		}()
	}
	go func() {
		err := http.Serve(l.HTTP, s.http)
		errChan <- err
	}()
	go func() {
		err := s.grpc.Serve(l.GRPC)
		errChan <- err
	}()

	return <-errChan
}

// Stop stops the server gracefully.
func (s *Server) Stop() {
	// TODO grace HTTP
	s.srv.grpc.GracefulStop()
}
