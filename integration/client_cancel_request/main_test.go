package main

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	strings_app "github.com/utrack/clay/integration/client_cancel_request/app/strings"
	strings_pb "github.com/utrack/clay/integration/client_cancel_request/pkg/strings"
)

func TestCancelRequest(t *testing.T) {
	ts := testServer()
	errlog := bytes.NewBuffer([]byte{})
	ts.Config.ErrorLog = log.New(errlog, "", 0)
	defer func() {
		ts.Close()
		if errlog.Len() > 0 {
			t.Fatalf("expected no errors, got: %s", errlog.Bytes())
		}
	}()

	httpClient := ts.Client()
	client := strings_pb.NewStringsHTTPClient(httpClient, ts.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	client.ToUpper(ctx, &strings_pb.String{Str: strings.Repeat("s", 10*1024)})
}

func testServer() *httptest.Server {
	mux := http.NewServeMux()
	desc := strings_app.NewStrings().GetDescription()
	desc.RegisterHTTP(mux)
	mux.Handle("/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(desc.SwaggerDef())
	}))

	ts := httptest.NewServer(mux)
	return ts
}
