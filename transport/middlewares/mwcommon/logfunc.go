package mwcommon

import (
	"context"

	"github.com/utrack/clay/v2/server/middlewares/mwcommon"
)

func GetLogFunc(logger interface{}) func(context.Context, string) {
	return mwcommon.GetLogFunc(logger)
}
