package httpruntime

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

type errResponse struct {
	Error string            `json:"error"`
	Data  map[string]string `json:"data"`
}

// SetError is used to output errors to the client.
// You can override that in the runtime.
var SetError func(context.Context, *http.Request, http.ResponseWriter, error, map[string]string) = DefaultSetError

// DefaultSetError is the default error output.
func DefaultSetError(ctx context.Context, req *http.Request, w http.ResponseWriter, err error, data map[string]string) {
	w.WriteHeader(500)
	enc := json.NewEncoder(w)
	enc.Encode(errResponse{Error: err.Error(), Data: data})
}
