package main

import (
	"net/http"

	"github.com/go-chi/chi"
	pb "github.com/utrack/clay/integration/proto_module_flag/strings"
)

func main() {
	r := chi.NewMux()
	desc := pb.NewStrings().GetDescription()
	desc.RegisterHTTP(r)

	r.Handle("/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(desc.SwaggerDef())
	}))

	http.ListenAndServe(":8080", r)
}
