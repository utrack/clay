package mwgrpc

import (
	"github.com/utrack/clay/v2/server/middlewares/mwgrpc"

	"google.golang.org/grpc"
)

// UnaryPanicHandler handles panics for UnaryHandlers.
func UnaryPanicHandler(logger interface{}) grpc.UnaryServerInterceptor {
	return mwgrpc.UnaryPanicHandler(logger)
}

// StreamPanicHandler handles panics for StreamHandlers.
func StreamPanicHandler(logger interface{}) grpc.StreamServerInterceptor {
	return mwgrpc.StreamPanicHandler(logger)
}
