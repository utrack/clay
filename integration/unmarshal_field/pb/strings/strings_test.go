package strings

import (
	"testing"

	"context"
	"github.com/stretchr/testify/assert"
	"github.com/utrack/clay/v2/transport"
	"google.golang.org/grpc"
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
		Foo:       &FooObj{},
	}

	cli := NewStringsHTTPClient(http.DefaultClient, ts.URL)

	type ft func(context.Context, *String, ...grpc.CallOption) (*String, error)

	ff := []ft{
		cli.ToLower1,
		cli.ToLower2,
		//cli.ToLower3,
		cli.ToLower4,
		cli.ToLower5,
	}

	for pos, f := range ff {
		_, err := f(context.Background(), obj)
		so.Nil(err, "func %v", pos)
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

func (i *StringsImplementation) ToLower1(ctx context.Context, req *String) (*String, error) {
	return i.f(ctx, req)
}
func (i *StringsImplementation) ToLower2(ctx context.Context, req *String) (*String, error) {
	return i.f(ctx, req)
}
func (i *StringsImplementation) ToLower3(ctx context.Context, req *String) (*String, error) {
	return i.f(ctx, req)
}
func (i *StringsImplementation) ToLower4(ctx context.Context, req *String) (*String, error) {
	return i.f(ctx, req)
}
func (i *StringsImplementation) ToLower5(ctx context.Context, req *String) (*String, error) {
	return i.f(ctx, req)
}

// GetDescription is a simple alias to the ServiceDesc constructor.
// It makes it possible to register the service implementation @ the server.
func (i *StringsImplementation) GetDescription() transport.ServiceDesc {
	return NewStringsServiceDesc(i)
}
