package mwgrpc

import (
	"fmt"
	"runtime/debug"

	"github.com/utrack/clay/v2/server/middlewares/mwcommon"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// UnaryPanicHandler handles panics for UnaryHandlers.
func UnaryPanicHandler(logger interface{}) grpc.UnaryServerInterceptor {
	logFunc := mwcommon.GetLogFunc(logger)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = grpc.Errorf(codes.Internal, "panic: %v", r)
				logFunc(ctx, fmt.Sprintf("recovered from panic: %v,\n%v ", r, string(debug.Stack())))
			}
		}()
		return handler(ctx, req)
	}

}

// StreamPanicHandler handles panics for StreamHandlers.
func StreamPanicHandler(logger interface{}) grpc.StreamServerInterceptor {
	logFunc := mwcommon.GetLogFunc(logger)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = grpc.Errorf(codes.Internal, "panic: %v", r)
				logFunc(stream.Context(), fmt.Sprintf("recovered from panic: %v, %v ", r, string(debug.Stack())))
			}
		}()

		return handler(srv, stream)
	}
}
