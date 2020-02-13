package strings

import (
	"testing"

	"bytes"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestBindingSubstruct(t *testing.T) {
	so := assert.New(t)

	reqBody := []byte(`{"req":"success"}`)
	req, err := http.NewRequest("PUT", "/strings/123", bytes.NewReader(reqBody))
	so.Nil(err)
	rc := chi.NewRouteContext()
	rc.URLParams.Add("id", "123")
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rc)
	req = req.WithContext(ctx)

	got := String{}

	f := unmarshaler_goclay_Strings_ToUpper_0(req)
	err = f(&got)
	so.Nil(err)
	so.Equal(int32(123), got.Id)
	so.Equal("success", got.Substruct.Req)
}
