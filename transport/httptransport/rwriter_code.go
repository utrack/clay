package httptransport

import (
	"net/http"
	"sync"

	"bufio"
	"github.com/pkg/errors"
	"net"
)

// CodedResponseWriter saves a response code to be read later.
type CodedResponseWriter struct {
	http.ResponseWriter

	// https://www.youtube.com/watch?v=HvEhOKkfoSo
	hijacker http.Hijacker

	m       sync.Mutex
	code    int
	codeSet bool

	written bool
}

// NewCodedWriter wraps an existing http.ResponseWriter.
func NewCodedWriter(w http.ResponseWriter) *CodedResponseWriter {
	ret := &CodedResponseWriter{
		ResponseWriter: w,
	}
	hj, ok := w.(http.Hijacker)
	if ok {
		ret.hijacker = hj
	}
	return ret
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

// Hijacker implements http.Hijacker.
//
// Works only if given http.ResponseWriter implements .Hijack().
func (w *CodedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.hijacker == nil {
		return nil, nil, errors.New("ResponseWriter does not support hijacking")
	}
	return w.hijacker.Hijack()
}
