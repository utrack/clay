package strings

import (
	"testing"

	"bytes"
	"net/http"

	"github.com/stretchr/testify/assert"
)

func TestBindingSubstruct(t *testing.T) {
	so := assert.New(t)

	reqBody := []byte(`{"req":"success"}`)
	req, err := http.NewRequest("PUT", "/strings/123", bytes.NewReader(reqBody))
	so.Nil(err)

	got := String{}

	f := unmarshaler_goclay_Strings_ToUpper_0(req)
	err = f(&got)
	so.Nil(err)
	so.Equal(int32(123), got.Id)
	so.Equal("success", got.Substruct.Req)
}
