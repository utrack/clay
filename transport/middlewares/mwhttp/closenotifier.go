package mwhttp

import (
	"github.com/utrack/clay/server/middlewares/mwhttp"
)

// CloseNotifier cancels the context if client goes away.
func CloseNotifier() Middleware {
	return mwhttp.CloseNotifier()
}
