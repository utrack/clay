// Code generated by protoc-gen-goclay, but you can (must) modify it.
// source: strings.proto

package strings

import (
	"context"

	desc "github.com/utrack/clay/integration/go_package_separate_impl_exists/pkg/strings"
)

func (i *StringsImplementation) ToLower(_ context.Context, req *desc.String) (*desc.String, error) {
	return &desc.String{Str: req.GetStr()}, nil
}
