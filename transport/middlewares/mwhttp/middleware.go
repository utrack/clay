package mwhttp

import (
	"github.com/utrack/clay/v2/server/middlewares/mwhttp"
)

// Middleware is the HTTP middleware type.
// It processes the request (potentially mutating it) and
// gives control to the underlying handler.
type Middleware = mwhttp.Middleware
