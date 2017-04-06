package mwhttp

import "net/http"

// Middleware is the HTTP middleware type.
// It processes the request (potentially mutating it) and
// gives control to the underlying handler.
type Middleware func(http.Handler) http.Handler
