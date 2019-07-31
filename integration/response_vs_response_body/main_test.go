package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/gogo/protobuf/types"
	stringspb "github.com/utrack/clay/integration/response_vs_response_body/pb"
	stringssrv "github.com/utrack/clay/integration/response_vs_response_body/strings"
)

func TestEcho(t *testing.T) {

	ts := testServer()
	defer ts.Close()

	tt := []struct {
		name string
		req  stringspb.Types
	}{
		{
			name: "double",
			req:  stringspb.Types{D: 1.0},
		},
		{
			name: "float",
			req:  stringspb.Types{F: 1.0},
		},
		{
			name: "int32",
			req:  stringspb.Types{I32: 1},
		},
		{
			name: "int64",
			req:  stringspb.Types{I64: 1},
		},
		{
			name: "uint32",
			req:  stringspb.Types{Ui32: 1},
		},
		{
			name: "uint64",
			req:  stringspb.Types{Ui64: 1},
		},
		{
			name: "sint32",
			req:  stringspb.Types{Si32: 1},
		},
		{
			name: "sint64",
			req:  stringspb.Types{Si64: 1},
		},
		{
			name: "fixed32",
			req:  stringspb.Types{Fixed32: 1},
		},
		{
			name: "fixed64",
			req:  stringspb.Types{Fixed64: 1},
		},
		{
			name: "sfixed32",
			req:  stringspb.Types{Sfixed32: 1},
		},
		{
			name: "sfixed64",
			req:  stringspb.Types{Sfixed64: 1},
		},
		{
			name: "bool",
			req:  stringspb.Types{B: true},
		},
		{
			name: "string",
			req:  stringspb.Types{S: "foo"},
		},
		{
			name: "bytes",
			req:  stringspb.Types{Bytes: []byte("bar")},
		},
		{
			name: "enum",
			req:  stringspb.Types{E: stringspb.Enum_FOO},
		},
		{
			name: "time",
			req:  stringspb.Types{Time: types.TimestampNow()},
		},
		{
			name: "duration",
			req:  stringspb.Types{Duration: types.DurationProto(3 * time.Second)},
		},
		// {
		// 	name: "stdtime",
		// 	req:  strings_pb.Types{Stdtime: time.Now().UTC()},
		// },
		// {
		// 	name: "stdduration",
		// 	req:  strings_pb.Types{Stdduration: 3 * time.Second},
		// },
	}

	httpClient := ts.Client()
	client := stringspb.NewStringsHTTPClient(httpClient, ts.URL)

	for _, tc := range tt {
		t.Run(fmt.Sprintf("%s and [%s]", tc.name, tc.name), func(t *testing.T) {
			ts.Lock()
			resp, err := client.Echo(context.Background(), &tc.req)
			echoBody, _ := ioutil.ReadAll(ts.RW)
			ts.Unlock()
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

			ts.Lock()
			resp2, err := client.Echo2(context.Background(), &stringspb.ListTypes{List: []*stringspb.Types{&tc.req}})
			echoBody2, _ := ioutil.ReadAll(ts.RW)
			ts.Unlock()
			if err != nil {
				t.Fatalf("expected err <nil>, got: %q", err)
			}
			if resp2 == nil {
				t.Fatalf("expected non-nil response, got nil")
			}

			if !reflect.DeepEqual(*resp2.List[0], tc.req) {
				t.Fatalf("expected %#v\n"+
					"got: %#v", tc.req, *resp2.List[0])
			}

			if "["+string(echoBody)+"]" != string(echoBody2) {
				t.Fatalf("expected <response from echo2> = `[` + <response from echo> +`]`, got\n"+
					"<response from echo>  = %s\n"+
					"<response from echo2> = %s", echoBody, echoBody2)
			}
		})
	}
}

func testServer() *Server {
	mux := http.NewServeMux()
	desc := stringssrv.NewStrings().GetDescription()
	desc.RegisterHTTP(mux)
	mux.Handle("/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(desc.SwaggerDef())
	}))
	ts := &Server{}
	ts.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ts.RW = NewRW(w)
		mux.ServeHTTP(ts.RW, req)
	}))
	return ts
}

type Server struct {
	sync.Mutex
	*httptest.Server
	RW *RW
}

func NewRW(w http.ResponseWriter) *RW {
	return &RW{
		w,
		bytes.NewBuffer([]byte{}),
	}
}

type RW struct {
	http.ResponseWriter
	*bytes.Buffer
}

func (w RW) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w RW) Write(b []byte) (int, error) {
	w.Buffer.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w RW) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}
