package strings

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSwaggerComments(t *testing.T) {
	so := assert.New(t)

	tags := []string{
		//"RPC_COMMENT",
		"FUNC1_COMMENT",
		"FUNC2_COMMENT",
		"STRUCT_COMMENT",
		"MEM_COMMENT",
	}

	for _, tag := range tags {
		so.True(strings.Contains(string(_swaggerDef_pb_strings_strings_proto), tag), tag)
	}

}
