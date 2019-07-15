package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/utrack/clay/integration/proto_two_services/strings"
)

func main() {
	r := chi.NewMux()
	desc1 := strings.NewStrings().GetDescription()
	desc1.RegisterHTTP(r)
	desc2 := strings.NewStrings2().GetDescription()
	desc2.RegisterHTTP(r)

	r.Handle("/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(desc1.SwaggerDef())
	}))

	http.ListenAndServe(":8080", r)
}
