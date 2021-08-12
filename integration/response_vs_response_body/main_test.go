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
	strings_pb "github.com/utrack/clay/integration/response_vs_response_body/pb"
	strings_srv "github.com/utrack/clay/integration/response_vs_response_body/strings"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestEcho(t *testing.T) {
	t.Skip("fix me!")

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
		{
			name: "time",
			req: strings_pb.Types{Time: timestamppb.New(func() time.Time {
				now := types.TimestampNow()
				t, err := types.TimestampFromProto(now)
				if err != nil {
					panic(err)
				}
				return t
			}())},
		},
		{
			name: "duration",
			req: strings_pb.Types{Duration: durationpb.New(func() time.Duration {
				d := types.DurationProto(3 * time.Second)
				dur, err := types.DurationFromProto(d)
				if err != nil {
					panic(err)
				}
				return dur
			}())},
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
	client := strings_pb.NewStringsHTTPClient(httpClient, ts.URL)

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
			resp2, err := client.Echo2(context.Background(), &strings_pb.ListTypes{
				List: []*strings_pb.Types{&tc.req},
			})
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
	desc := strings_srv.NewStrings().GetDescription()
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
