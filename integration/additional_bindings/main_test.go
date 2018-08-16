package main_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"

	strings_pb "github.com/utrack/clay/integration/additional_bindings/pb"
	strings_srv "github.com/utrack/clay/integration/additional_bindings/strings"
)

func TestToUpper(t *testing.T) {
	ts := testServer()
	defer ts.Close()
	t.Run("GET slice of strings in request and an object in response", func(t *testing.T) {
		httpClient := ts.Client()
		client := strings_pb.NewStringsHTTPClient(httpClient, ts.URL)

		req := &strings_pb.String{
			Str: "foo",
		}
		resp, err := client.ToUpper(context.Background(), req)
		if err != nil {
			t.Fatalf("expected err <nil>, got: %q", err)
		}

		got := resp.GetStr()
		expected := strings.ToUpper(req.GetStr())
		if got != expected {
			t.Fatalf("expected %q, got: %q", expected, got)
		}
	})
}

func testServer() *httptest.Server {
	mux := chi.NewRouter()
	desc := strings_srv.NewStrings().GetDescription()
	desc.RegisterHTTP(mux)
	mux.Handle("/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(desc.SwaggerDef())
	}))
	ts := httptest.NewServer(mux)
	return ts
}
