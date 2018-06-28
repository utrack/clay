package mwhttp

import (
	"github.com/utrack/clay/v2/server/middlewares/mwhttp"
)

// CloseNotifier cancels the context if client goes away.
func CloseNotifier() Middleware {
	return mwhttp.CloseNotifier()
}
