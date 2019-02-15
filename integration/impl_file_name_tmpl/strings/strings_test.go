package strings

import (
	"context"
	"testing"

	desc "github.com/utrack/clay/integration/impl_file_name_tmpl/pb"
)

// TestImplementationExists tests if the implementation is not re-written after generation
func TestImplementationExists(t *testing.T) {
	impl := NewStrings()
	res, err := impl.ToUpper(context.Background(), &desc.String{Str: "foo"})
	if err != nil {
		t.Fatalf("error is expected to be nil, got %q\n Implementation was re-generated!", err)
	}
	if res == nil || res.Str != "FOO" {
		t.Fatalf("wrong result, got %v", res)
	}
}
