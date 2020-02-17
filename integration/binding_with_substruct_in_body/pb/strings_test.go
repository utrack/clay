package strings

import (
	"context"
	"testing"

	"bytes"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestBindingSubstruct(t *testing.T) {
	so := assert.New(t)

	reqBody := []byte(`{"req":"success"}`)
	req, err := http.NewRequest("PUT", "/strings/123", bytes.NewReader(reqBody))

	cctx := chi.NewRouteContext()
	cctx.URLParams.Add("substruct.id", "123")
	req = req.WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, cctx))
	so.Nil(err)

	got := String{}

	f := unmarshaler_goclay_Strings_ToUpper_0(req)
	err = f(&got)
	so.Nil(err)
	so.Equal(int32(123), got.Substruct.Id)
	so.Equal("success", got.Substruct.Reqs1.Req)
}
