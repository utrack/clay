package main

import (
	"net/http"

	"github.com/go-chi/chi"
	desc2 "github.com/utrack/clay/integration/go_package_separate_impl_exists/internal/app/strings"
)

func main() {
	r := chi.NewMux()
	desc := desc2.NewStrings().GetDescription()
	desc.RegisterHTTP(r)

	r.Handle("/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		_, _ = w.Write(desc.SwaggerDef())
	}))

	_ = http.ListenAndServe(":8080", r)
}
