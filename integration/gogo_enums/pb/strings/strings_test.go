package strings

import (
	"testing"

	"context"
	"github.com/stretchr/testify/assert"
	"github.com/utrack/clay/v2/transport"
	"net/http"
	"net/http/httptest"
)

func TestMarshalUnmarshal(t *testing.T) {
	so := assert.New(t)
	impl, ts := getTestSvc()
	defer ts.Close()
	impl.f = func(ctx context.Context, req *String) (*String, error) {
		return &String{}, nil
	}

	obj := &String{
		SnakeCase: "foo",
		Strtype:   StringType_STRING_TYPE_BAZ,
	}

	cli := NewStringsHTTPClient(http.DefaultClient, ts.URL)
	_, err := cli.ToLower(context.Background(), obj)
	so.Nil(err)
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
