package main_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/utrack/clay/integration/binding_with_body_and_response/strings"
)

func TestToUpper(t *testing.T) {
	ts := testServer()
	defer ts.Close()
	t.Run("slice of strings in request and slice of strings in response", func(t *testing.T) {
		rsp, err := ts.Client().Post(ts.URL+"/strings/to_upper", "", bytes.NewReader([]byte(`["test","boo"]`)))
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
}

func testServer() *httptest.Server {
	mux := http.NewServeMux()
	desc := strings.NewStrings().GetDescription()
	desc.RegisterHTTP(mux)
	ts := httptest.NewServer(mux)
	return ts
}
