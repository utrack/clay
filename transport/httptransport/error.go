package httptransport

// MarshalerError is returned by a marshaler func.
// It is used to decorate errors coming from gRPC-generated _Handler
// to distinguish parser errors from handlers' errors.
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
