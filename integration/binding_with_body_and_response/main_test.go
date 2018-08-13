package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/utrack/clay/integration/binding_with_body_and_response/strings"
)

func TestToUpper(t *testing.T) {
	ts := testServer()
	defer ts.Close()
	t.Run("POST slice of strings in request and slice of strings in response", func(t *testing.T) {
		rsp, err := ts.Client().Post(ts.URL+"/strings/to_upper", "application/javascript", bytes.NewReader([]byte(`["test","boo"]`)))
		if err != nil {
			t.Fatalf("expected err <nil>, got: %s", err)
		}
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			t.Fatalf("expected err <nil>, got: %s", err)
		}
		if string(body) != `["TEST","BOO"]` {
			t.Fatalf("expected response body `[\"TEST\",\"BOO\"]`, got: %s", body)
		}
	})
	t.Run("GET slice of strings in request and slice of strings in response", func(t *testing.T) {
		rsp, err := ts.Client().Get(ts.URL + "/strings/to_upper/v2?str=test&str=boo")
		if err != nil {
			t.Fatalf("expected err <nil>, got: %s", err)
		}
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			t.Fatalf("expected err <nil>, got: %s", err)
		}
		if string(body) != `["TEST","BOO"]` {
			t.Fatalf("expected response body `[\"TEST\",\"BOO\"]`, got: %s", body)
		}
	})
	t.Run("check response scheme in swagger json definition", func(t *testing.T) {
		rsp, err := ts.Client().Post(ts.URL+"/swagger.json", "application/javascript", nil)
		if err != nil {
			t.Fatalf("expected err <nil>, got: %s", err)
		}
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			t.Fatalf("expected err <nil>, got: %s", err)
		}
		var s = &spec.Swagger{}
		if err = s.UnmarshalJSON(body); err != nil {
			t.Fatalf("expected err <nil> during unmarshal swagger json, got: %s, swagger.json: %s", err, body)
		}
		var schemaType *string
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("panic: %s", r)
				}
			}()
			schemaType = &(s.Paths.Paths["/strings/to_upper"].Post.Responses.StatusCodeResponses[200].Schema.Type[0])
		}()
		if schemaType == nil {
			t.Fatalf("expected schema type for response is array, got: %v", schemaType)
		}
		if *schemaType != "array" {
			t.Fatalf("expected schema type for response is array, got: %v", *schemaType)
		}
	})
}

func testServer() *httptest.Server {
	mux := http.NewServeMux()
	desc := strings.NewStrings().GetDescription()
	desc.RegisterHTTP(mux)
	mux.Handle("/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(desc.SwaggerDef())
	}))
	ts := httptest.NewServer(mux)
	return ts
}
