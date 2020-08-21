package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	strings_pb "github.com/utrack/clay/integration/binding_with_different_types/pb"
	strings_srv "github.com/utrack/clay/integration/binding_with_different_types/strings"
)

func TestEcho(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	tt := []struct {
		name string
		req  strings_pb.Types
	}{
		{
			name: "double",
			req:  strings_pb.Types{D: 1.0},
		},
		{
			name: "float",
			req:  strings_pb.Types{F: 1.0},
		},
		{
			name: "int32",
			req:  strings_pb.Types{I32: 1},
		},
		{
			name: "int64",
			req:  strings_pb.Types{I64: 1},
		},
		{
			name: "uint32",
			req:  strings_pb.Types{Ui32: 1},
		},
		{
			name: "uint64",
			req:  strings_pb.Types{Ui64: 1},
		},
		{
			name: "sint32",
			req:  strings_pb.Types{Si32: 1},
		},
		{
			name: "sint64",
			req:  strings_pb.Types{Si64: 1},
		},
		{
			name: "fixed32",
			req:  strings_pb.Types{Fixed32: 1},
		},
		{
			name: "fixed64",
			req:  strings_pb.Types{Fixed64: 1},
		},
		{
			name: "sfixed32",
			req:  strings_pb.Types{Sfixed32: 1},
		},
		{
			name: "sfixed64",
			req:  strings_pb.Types{Sfixed64: 1},
		},
		{
			name: "bool",
			req:  strings_pb.Types{B: true},
		},
		{
			name: "string",
			req:  strings_pb.Types{S: "foo"},
		},
		{
			name: "bytes",
			req:  strings_pb.Types{Bytes: []byte("bar")},
		},
		{
			name: "enum",
			req:  strings_pb.Types{E: strings_pb.Enum_FOO},
		},
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("GET echo request for %s", tc.name), func(t *testing.T) {
			httpClient := ts.Client()
			client := strings_pb.NewStringsHTTPClient(httpClient, ts.URL)

			resp, err := client.Echo(context.Background(), &tc.req)
			if err != nil {
				t.Fatalf("expected err <nil>, got: %q", err)
			}
			if resp == nil {
				t.Fatalf("expected non-nil response, got nil")
			}

			if !reflect.DeepEqual(*resp, tc.req) {
				t.Fatalf("expected %#v\n"+
					"got: %#v", tc.req, *resp)
			}
		})
	}
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
