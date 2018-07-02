package strings

import (
	"context"
	"strings"

	desc "github.com/utrack/clay/integration/impl_exists/pb"
	transport "github.com/utrack/clay/v2/transport"
)

type StringsImplementation struct{}

func NewStrings() *StringsImplementation {
	return &StringsImplementation{}
}

func (i *StringsImplementation) ToUpper(ctx context.Context, req *desc.String) (*desc.String, error) {
	return &desc.String{Str: strings.ToUpper(req.Str)}, nil
}

// GetDescription is a simple alias to the ServiceDesc constructor.
// It makes it possible to register the service implementation @ the server.
func (i *StringsImplementation) GetDescription() transport.ServiceDesc {
	return desc.NewStringsServiceDesc(i)
}
