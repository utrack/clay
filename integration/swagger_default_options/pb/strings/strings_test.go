package strings

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utrack/clay/v2/transport"
	"github.com/utrack/clay/v2/transport/swagger"
)

func TestSwaggerComments(t *testing.T) {
	so := assert.New(t)

	d := NewStringsServiceDesc(nil)
	desc := "some description here"
	d.Apply(transport.WithSwaggerOptions(swagger.WithDescription(desc)))

	sdef := string(d.SwaggerDef())

	so.True(strings.Contains(sdef, desc))
}
