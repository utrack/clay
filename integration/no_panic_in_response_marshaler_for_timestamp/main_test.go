package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	strings_pb "github.com/utrack/clay/integration/no_panic_in_response_marshaler_for_timestamp/pb"
	strings_srv "github.com/utrack/clay/integration/no_panic_in_response_marshaler_for_timestamp/strings"
)

func TestEcho(t *testing.T) {
	ts := testServer(t)
	defer ts.Close()

	httpClient := ts.Client()
	client := strings_pb.NewStringsHTTPClient(httpClient, ts.URL)
	client.Echo(context.Background(), &strings_pb.EchoReq{})
}

func testServer(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()
	desc := strings_srv.NewStrings().GetDescription()
	desc.RegisterHTTP(mux)
	mux.Handle("/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(desc.SwaggerDef())
	}))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		buf := RW{bytes.NewBuffer([]byte{})}
		defer func() {
			got, _ := ioutil.ReadAll(buf)
			if expected := `{"stdtime":"2018-01-01T01:01:01.000000001Z"}`; expected != strings.TrimSpace(string(got)) {
				t.Errorf("expected response: `%s`, got: `%s`", expected, got)
			}
			p := recover()
			if p != nil {
				t.Fatalf("unexpected panic: `%v`", p)
			}
		}()
		mux.ServeHTTP(buf, req)
	}))

	return ts
}

type RW struct {
	*bytes.Buffer
}

func (w RW) Header() http.Header {
	return map[string][]string{}
}

func (w RW) Write(b []byte) (int, error) {
	return w.Buffer.Write(b)
}

func (w RW) WriteHeader(statusCode int) {

}
