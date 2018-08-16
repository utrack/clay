package strings

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utrack/clay/v2/transport"
	"io/ioutil"
)

// assert type for the headers
type a struct {
	Name   string
	Resp *String
	MarshaledData string
}

// TestHTTPHeadersPass_vanillaHTTP tests that default HTTP middleware passes HTTP headers to
// gRPC metadata.
func TestHTTPHeadersPass_vanillaHTTP(t *testing.T) {
	so := assert.New(t)
	impl, ts := getTestSvc()
	defer ts.Close()

	tc := []a{
		a{
			Name:   "Empty string field",
			Resp: &String{Str: ""},
			MarshaledData: `{"str":""}`,
		},
	}




	req, err := http.NewRequest("POST", ts.URL+pattern_goclay_Strings_ToLower_0_builder(), bytes.NewReader([]byte("{}")))
	so.Nil(err)

	for _, c := range tc {
		impl.f = func(ctx context.Context, req *String) (*String, error) {
			return c.Resp, nil
		}

		resp, err := http.DefaultClient.Do(req)
		so.Nil(err)
		data, err := ioutil.ReadAll(resp.Body)
		so.Equal(c.MarshaledData, string(data))
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
