package strings

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utrack/clay/v2/transport"
	"google.golang.org/grpc/metadata"
)

// assert type for the headers
type a struct {
	Name   string
	Values []string
}

// TestHTTPHeadersPass_vanillaHTTP tests that default HTTP middleware passes HTTP headers to
// gRPC metadata.
func TestHTTPHeadersPass_vanillaHTTP(t *testing.T) {
	so := assert.New(t)
	impl, ts := getTestSvc()
	defer ts.Close()

	tc := []a{
		a{
			Name:   "X-Something-For-Test",
			Values: []string{"val1", "Val2"},
		},
		a{
			Name:   "User-Agent",
			Values: []string{"Go-http-client/1.1"},
		},
	}

	calledFunc := false

	impl.f = func(ctx context.Context, req *String) (*String, error) {
		calledFunc = true
		md, ok := metadata.FromIncomingContext(ctx)
		so.True(ok)

		for _, c := range tc {
			got := md.Get(strings.ToLower(c.Name))
			so.EqualValues(c.Values, got)
		}

		return nil, nil
	}

	req, err := http.NewRequest("POST", ts.URL+pattern_goclay_Strings_ToLower_0_builder(), bytes.NewReader([]byte("{}")))
	so.Nil(err)

	for _, c := range tc {
		for pos := range c.Values {
			req.Header.Add(c.Name, c.Values[pos])
		}
	}

	_, err = http.DefaultClient.Do(req)
	so.Nil(err)
	so.True(calledFunc)
}

// TestHTTPHeadersPass_genClient tests that default HTTP middleware passes HTTP headers to
// gRPC metadata using generated client.
func TestHTTPHeadersPass_genClient(t *testing.T) {
	so := assert.New(t)
	impl, ts := getTestSvc()
	defer ts.Close()

	tc := []a{
		a{
			Name:   "User-Agent",
			Values: []string{"Go-http-client/1.1"},
		},
		a{
			Name:   "Accept",
			Values: []string{"application/json"},
		},
	}

	calledFunc := false

	impl.f = func(ctx context.Context, req *String) (*String, error) {
		calledFunc = true
		md, ok := metadata.FromIncomingContext(ctx)
		so.True(ok)

		for _, c := range tc {
			got := md.Get(strings.ToLower(c.Name))
			so.EqualValues(c.Values, got)
		}

		return &String{}, nil
	}

	cli := NewStringsHTTPClient(http.DefaultClient, ts.URL)
	_, err := cli.ToLower(context.Background(), &String{})
	so.Nil(err)
	so.True(calledFunc)
}

// TestHTTPHeadersPass_genClient_outgoingContext tests that generated HTTP client
// passes headers from grpc.ToOutgoingContext to the request by default.
func TestHTTPHeadersPass_genClient_outgoingContext(t *testing.T) {
	so := assert.New(t)
	impl, ts := getTestSvc()
	defer ts.Close()

	tc := []a{
		a{
			Name:   "User-Agent",
			Values: []string{"Go-http-client/1.1"},
		},
		a{
			Name:   "Accept",
			Values: []string{"application/json"},
		},
	}
	pt := []a{
		a{
			Name:   "X-Test-Passthrough",
			Values: []string{"v1", "Value2", "3"},
		},
	}

	calledFunc := false

	impl.f = func(ctx context.Context, req *String) (*String, error) {
		calledFunc = true
		md, ok := metadata.FromIncomingContext(ctx)
		so.True(ok)

		for _, c := range tc {
			got := md.Get(strings.ToLower(c.Name))
			so.EqualValues(c.Values, got)
		}

		for _, c := range pt {
			got := md.Get(strings.ToLower(c.Name))
			so.EqualValues(c.Values, got)
		}

		return &String{}, nil
	}

	ctx := context.Background()

	for _, c := range pt {
		for i := range c.Values {
			ctx = metadata.AppendToOutgoingContext(ctx, c.Name, c.Values[i])
		}
	}

	cli := NewStringsHTTPClient(http.DefaultClient, ts.URL)
	_, err := cli.ToLower(ctx, &String{})
	so.Nil(err)
	so.True(calledFunc)
}

func getTestSvc() (*StringsImplementation, *httptest.Server) {
	mux := http.NewServeMux()
	impl := NewStrings()
	d := impl.GetDescription()
	d.RegisterHTTP(mux)

	ts := httptest.NewServer(mux)
	return impl, ts
}

type StringsImplementation struct {
	f func(ctx context.Context, req *String) (*String, error)
}

func NewStrings() *StringsImplementation {
	return &StringsImplementation{}
}

func (i *StringsImplementation) ToLower(ctx context.Context, req *String) (*String, error) {
	return i.f(ctx, req)
}

// GetDescription is a simple alias to the ServiceDesc constructor.
// It makes it possible to register the service implementation @ the server.
func (i *StringsImplementation) GetDescription() transport.ServiceDesc {
	return NewStringsServiceDesc(i)
}
