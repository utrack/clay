package httptransport

import (
	"net/http"
	"sync"
)

// CodedResponseWriter saves a response code to be read later.
type CodedResponseWriter struct {
	http.ResponseWriter

	m       sync.Mutex
	code    int
	codeSet bool

	written bool
}

// NewCodedWriter wraps an existing http.ResponseWriter.
func NewCodedWriter(w http.ResponseWriter) *CodedResponseWriter {
	return &CodedResponseWriter{
		ResponseWriter: w,
	}
}

func (w *CodedResponseWriter) Write(b []byte) (int, error) {
	w.written = true

	return w.ResponseWriter.Write(b)
}

func (w *CodedResponseWriter) WriteHeader(statusCode int) {
	w.m.Lock()
	defer w.m.Unlock()

	w.code = statusCode
	w.codeSet = true
	w.written = true
	w.ResponseWriter.WriteHeader(statusCode)
}

// ResponseCode returns a code that was written to the transport.
func (w *CodedResponseWriter) ResponseCode() int {
	w.m.Lock()
	defer w.m.Unlock()

	if !w.codeSet {
		return http.StatusOK
	}
	return w.code
}

// Written returns true if headers were written (or the whole response).
func (w *CodedResponseWriter) Written() bool {
	w.m.Lock()
	defer w.m.Unlock()

	return w.written
}
