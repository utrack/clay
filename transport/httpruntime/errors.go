package httpruntime

import (
	"context"
	"encoding/json"
	"net/http"
)

type errResponse struct {
	Error string `json:"error"`
}

func SetError(ctx context.Context, req *http.Request, w http.ResponseWriter, err error) {
	// TODO write logs
	w.WriteHeader(500)
	enc := json.NewEncoder(w)
	enc.Encode(errResponse{Error: err.Error()})
}
