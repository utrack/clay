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
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// assert type for the headers
type a struct {
	Name   string
	Values []string
}

// TestHTTPHeadersResponse tests that server is able to return HTTP headers via
// grpc.SetHeaders.
func TestHTTPHeadersResponse(t *testing.T) {
	so := assert.New(t)
	impl, ts := getTestSvc()
	defer ts.Close()

	tc := []a{
		a{
			Name:   "X-Something-For-Test",
			Values: []string{"val1", "Val2"},
		},
		a{
			Name:   "X-Something-For-Test-2",
			Values: []string{"val3", "Val2"},
		},
	}

	calledFunc := false

	impl.f = func(ctx context.Context, req *String) (*String, error) {
		calledFunc = true

		for _, c := range tc {
			md := metadata.New(map[string]string{})

			for i := range c.Values {
				md.Append(c.Name, c.Values[i])
			}
			grpc.SetHeader(ctx, md)
		}

		return nil, nil
	}

	req, err := http.NewRequest("POST", ts.URL+pattern_goclay_Strings_ToLower_0_builder(&String{}), bytes.NewReader([]byte("{}")))
	so.Nil(err)

	rsp, err := http.DefaultClient.Do(req)
	so.Nil(err)
	so.True(calledFunc)

	for _, c := range tc {
		so.EqualValues(c.Values, rsp.Header[c.Name])
	}
}

// TestHTTPHeadersResponse_genClient tests that genrated client is able to receive
// headers via grpc.Headers.
func TestHTTPHeadersResponse_genClient(t *testing.T) {
	so := assert.New(t)
	impl, ts := getTestSvc()
	defer ts.Close()

	tc := []a{
		a{
			Name:   "X-Something-For-Test",
			Values: []string{"val1", "Val2"},
		},
		a{
			Name:   "X-Something-For-Test-2",
			Values: []string{"val3", "Val2"},
		},
	}

	calledFunc := false

	impl.f = func(ctx context.Context, req *String) (*String, error) {
		calledFunc = true

		md := metadata.New(map[string]string{})
		for _, c := range tc {

			for i := range c.Values {
				md.Append(c.Name, c.Values[i])
			}
		}
		grpc.SetHeader(ctx, md)

		return &String{}, nil
	}

	cli := NewStringsHTTPClient(http.DefaultClient, ts.URL)

	gotHeaders := metadata.MD{}

	_, err := cli.ToLower(
		context.Background(),
		&String{},
		grpc.Header(&gotHeaders),
	)
	so.Nil(err)
	so.True(calledFunc)

	for _, c := range tc {
		so.EqualValues(c.Values, gotHeaders.Get(strings.ToLower(c.Name)))
	}
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
	UnimplementedStringsServer
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
