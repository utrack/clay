package mwhttp

import (
	"github.com/utrack/clay/v2/server/middlewares/mwhttp"
)

// Recover recovers HTTP server from handlers' panics.
func Recover(logger interface{}) Middleware {
	return mwhttp.Recover(logger)
}
