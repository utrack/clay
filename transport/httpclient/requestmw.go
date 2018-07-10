package httpclient

import (
	"fmt"
	"net/http"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RequestMiddleware processes HTTP requests and responses vs provided ClientOptions.
type RequestMiddleware struct {
	out []RequestMutator
	in  []ResponseMutator
}

// NewMiddlewareGRPC creates new RequestMiddleware from gRPC call options.
func NewMiddlewareGRPC(opts []grpc.CallOption) (*RequestMiddleware, error) {
	ret := &RequestMiddleware{
		out: []RequestMutator{},
		in:  []ResponseMutator{},
	}
	err := ret.applyClientOptions(opts)
	return ret, err
}

// applyClientOptions applies grpc.Options for calls via HTTP.
func (c *RequestMiddleware) applyClientOptions(opts []grpc.CallOption) error {
	for _, ou := range opts {
		switch o := ou.(type) {
		case grpc.HeaderCallOption:
			c.in = append(c.in, clientRspHeaderCopier(o.HeaderAddr))
		default:
			return fmt.Errorf("Unsupported gRPC-to-HTTP call option: %v", reflect.TypeOf(o).String())
		}
	}
	return nil
}

// ProcessRequest processes outgoing HTTP requests.
func (c *RequestMiddleware) ProcessRequest(r *http.Request) (*http.Request, error) {
	var err error
	for _, m := range DefaultRequestMutators {
		r, err = m(r)
		if err != nil {
			return r, err
		}
	}
	for _, m := range c.out {
		r, err = m(r)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// ProcessResponse processes received HTTP responses.
func (c *RequestMiddleware) ProcessResponse(r *http.Response) (*http.Response, error) {
	var err error
	for _, m := range DefaultResponseMutators {
		r, err = m(r)
		if err != nil {
			return r, err
		}
	}
	for _, m := range c.in {
		r, err = m(r)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

// RequestMutator processes and/or mutates outgoing HTTP requests.
type RequestMutator func(*http.Request) (*http.Request, error)

// ResponseMutator processes and/or mutates HTTP responses.
type ResponseMutator func(*http.Response) (*http.Response, error)

// DefaultRequestMutators are used for every outgoing request.
var DefaultRequestMutators = []RequestMutator{clientReqHeadersFromMD()}

// DefaultResponseMutators are used for every received response.
var DefaultResponseMutators = []ResponseMutator{}

func clientRspHeaderCopier(md *metadata.MD) ResponseMutator {
	return func(rsp *http.Response) (*http.Response, error) {
		h := rsp.Header
		for k := range h {
			md.Append(k, h[k]...)
		}
		return rsp, nil
	}
}

// clientReqHeadersFromMD pushes metadata from OutgoingContext to the
// request headers.
func clientReqHeadersFromMD() RequestMutator {
	return func(req *http.Request) (*http.Request, error) {
		ctxmd, ok := metadata.FromOutgoingContext(req.Context())
		if !ok {
			return req, nil
		}

		for k := range ctxmd {
			vv := ctxmd.Get(k)
			for i := range vv {
				req.Header.Add(k, vv[i])
			}
		}

		return req, nil
	}

}
