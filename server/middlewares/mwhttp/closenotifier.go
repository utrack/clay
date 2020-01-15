// +build !go1.9

package mwhttp

import "github.com/go-chi/chi/middleware"

// CloseNotifier cancels the context if client goes away.
func CloseNotifier() Middleware {
	return middleware.CloseNotify
}
