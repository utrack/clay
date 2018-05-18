package server

import (
	"net/http"

	"github.com/go-chi/chi"
)

type chiWrapper struct {
	chi.Router
}

func (c *chiWrapper) MethodFunc(pattern string, method string, h func(http.ResponseWriter, *http.Request)) {
	c.Router.MethodFunc(method, pattern, h)
}
func (c *chiWrapper) Method(pattern string, method string, h http.Handler) {
	c.Router.Method(method, pattern, h)
}

func (c *chiWrapper) Use(middlewares ...func(http.Handler) http.Handler) {
	c.Router.Use(middlewares...)
}
