package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/utrack/clay/integration/impl_standalone_without_swagger/strings"
)

func main() {
	r := chi.NewMux()
	desc := strings.NewStrings().GetDescription()
	desc.RegisterHTTP(r)

	http.ListenAndServe(":8080", r)
}
