package mwcommon

import (
	"context"
	"fmt"
	"reflect"

	"github.com/utrack/clay/v2/server/log"
)

func GetLogFunc(logger interface{}) func(context.Context, string) {
	if logger, ok := logger.(log.Writer); ok {
		return func(_ context.Context, s string) {
			logger.Log(log.LevelError, s)
		}
	}
	if logger, ok := logger.(log.WriterC); ok {
		return func(ctx context.Context, s string) {
			logger.Logc(ctx, log.LevelError, s)
		}
	}
	panic(fmt.Sprintf("Bad type passed to getLogFunc: %v", reflect.TypeOf(logger)))
}
