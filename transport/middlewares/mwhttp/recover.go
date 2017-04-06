package mwhttp

import (
	"net/http"
	"runtime/debug"

	"github.com/utrack/clay/transport/httpruntime"

	"github.com/pkg/errors"
)

// Recover recovers HTTP server from handlers' panics.
func Recover() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					stack := debug.Stack()
					httpruntime.SetError(
						r.Context(),
						r, w,
						errors.Errorf("recovered from panic: %v", rec),
						map[string]string{"stack": string(stack)},
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
