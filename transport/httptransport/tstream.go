package httptransport

import (
	"errors"
	"net/http"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TransportStream implements grpc.ServerTransportStream for the HTTP calls.
type TransportStream struct {
	mu sync.Mutex
	w  http.ResponseWriter
}

var _ grpc.ServerTransportStream = &TransportStream{}

// NewTStream creates and returns new TransportStream writing to supplied http.ResponseWriter.
func NewTStream(w http.ResponseWriter) *TransportStream {
	return &TransportStream{
		w: w,
	}
}

// SetHeader implements grpc.ServerTransportStream.
func (ts *TransportStream) SetHeader(md metadata.MD) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	for k := range md {
		vv := md.Get(k)
		for i := range vv {
			ts.w.Header().Add(k, vv[i])
		}
	}
	return nil
}

// Method implements grpc.ServerTransportStream.
func (ts *TransportStream) Method() string {
	panic("Method is not supported for the HTTP TransportStream")
}

// SendHeader implements grpc.ServerTransportStream.
func (ts *TransportStream) SendHeader(md metadata.MD) error {
	for k := range md {
		vv := md.Get(k)
		for i := range vv {
			ts.w.Header().Add(k, vv[i])
		}
	}
	ts.w.WriteHeader(http.StatusOK)
	return nil
}

// SetTrailer implements grpc.ServerTransportStream.
func (ts *TransportStream) SetTrailer(md metadata.MD) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	t, ok := ts.w.(*CodedResponseWriter)
	if !ok {
		return errors.New("SetTrailer is unsupported for this ResponseWriter")
	}

	written := t.Written()

	for k := range md {
		// Prefix the key if headers were written; add name to the Trailer otherwise
		// see http.Header
		if written {
			k = http.TrailerPrefix + k
		} else {
			ts.w.Header().Add("Trailer", k)
		}

		vv := md.Get(k)
		for i := range vv {
			ts.w.Header().Add(k, vv[i])
		}
	}
	return nil
}
