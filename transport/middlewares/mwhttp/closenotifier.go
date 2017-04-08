package mwhttp

import "github.com/pressly/chi/middleware"

// CloseNotifier cancels the context if client goes away.
func CloseNotifier() Middleware {
	return middleware.CloseNotify
}
