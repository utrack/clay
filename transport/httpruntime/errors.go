package httpruntime

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

type errResponse struct {
	Error string `json:"error"`
}

// SetError is used to output errors to the client.
// You can override that in the runtime.
var SetError func(context.Context, *http.Request, http.ResponseWriter, error) = DefaultSetError

// DefaultSetError is the default error output.
func DefaultSetError(ctx context.Context, req *http.Request, w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)
	enc := json.NewEncoder(w)
	enc.Encode(errResponse{Error: err.Error()})
}

// TransformUnmarshalerError is called for every error reported by unmarshaler.
// It can be used to transform the error returned to the client (embed HTTP code in it,
// mask text, etc.).
var TransformUnmarshalerError = func(err error) error { return err }

// MarshalerError is returned by a marshaler.
type MarshalerError struct {
	Err error
}

func (m MarshalerError) Cause() error {
	return m.Err
}

func (m MarshalerError) Error() string {
	return m.Err.Error()
}

func NewMarshalerError(err error) MarshalerError {
	return MarshalerError{Err: err}
}
