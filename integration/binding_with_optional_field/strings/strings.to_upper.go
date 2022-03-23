// Code generated by protoc-gen-goclay, but your can (must) modify it.
// source: pb/strings.proto

package strings

import (
	"context"
	desc "github.com/utrack/clay/integration/binding_with_optional_field/pb"
	"strings"
)

func (i *StringsImplementation) ToUpper(ctx context.Context, req *desc.String) (rsp *desc.String, err error) {
	rsp = &desc.String{}
	if req.Str == nil {
		rsp.Str = nil
	} else {
		upperStr := strings.ToUpper(*req.Str)
		rsp.Str = &upperStr
	}

	return
}
