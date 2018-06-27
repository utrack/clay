package mwcommon

import (
	"context"

	"github.com/utrack/clay/server/middlewares/mwcommon"
)

func GetLogFunc(logger interface{}) func(context.Context, string) {
	return mwcommon.GetLogFunc(logger)
}
