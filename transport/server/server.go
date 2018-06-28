package server

import (
	"github.com/utrack/clay/v2/server"
)

// Server is a transport server.
type Server = server.Server

// NewServer creates a Server listening on the rpcPort.
// Pass additional Options to mutate its behaviour.
// By default, HTTP JSON handler and gRPC are listening on the same
// port, admin port is p+2 and profile port is p+4.
func NewServer(rpcPort int, opts ...Option) *Server {
	return server.NewServer(rpcPort, opts...)
}
