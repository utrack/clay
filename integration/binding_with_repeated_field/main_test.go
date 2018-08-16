package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	strings_pb "github.com/utrack/clay/integration/binding_with_repeated_field/pb"
	strings_srv "github.com/utrack/clay/integration/binding_with_repeated_field/strings"
)

func TestToUpper(t *testing.T) {
	ts := testServer()
	defer ts.Close()
	t.Run("GET slice of strings in request and an object in response", func(t *testing.T) {
		httpClient := ts.Client()
		client := strings_pb.NewStringsHTTPClient(httpClient, ts.URL)

		req := &strings_pb.String{
			Str: []string{"foo", "bar"},
		}
		resp, err := client.ToUpper(context.Background(), req)
		if err != nil {
			t.Fatalf("expected err <nil>, got: %q", err)
		}
		if len(resp.GetStr()) != 2 {
			t.Fatalf("expected response len 2, got: %d", len(resp.GetStr()))
		}
		for i, got := range resp.GetStr() {
			expected := strings.ToUpper(req.GetStr()[i])
			if got != expected {
				t.Fatalf("expected %q, got: %q", expected, got)
			}
		}
	})
}

func testServer() *httptest.Server {
	mux := http.NewServeMux()
	desc := strings_srv.NewStrings().GetDescription()
	desc.RegisterHTTP(mux)
	mux.Handle("/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(desc.SwaggerDef())
	}))
	ts := httptest.NewServer(mux)
	return ts
}
