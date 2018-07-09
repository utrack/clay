package strings

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/utrack/clay/v2/transport/httpruntime"
)

func TestMarshalUnmarshal(t *testing.T) {
	so := assert.New(t)
	obj := &String{
		SnakeCase: "foo",
		Strtype:   StringType_STRING_TYPE_BAR,
	}
	m := httpruntime.DefaultMarshaler(nil)

	buf := bytes.NewBuffer(nil)

	err := m.Marshal(buf, obj)
	so.Nil(err)

	got := &String{}
	logrus.Info(buf.String())

	err = m.Unmarshal(buf, got)
	so.Nil(err)

	so.Equal(*obj, *got)
}

// TestFieldDefs tests that our Swagger names correspond to generated json names.
func TestFieldDefs(t *testing.T) {
	so := assert.New(t)
	obj := &String{
		SnakeCase: "no step",
		Strtype:   StringType_STRING_TYPE_BAR,
	}
	m := httpruntime.DefaultMarshaler(nil)

	buf := bytes.NewBuffer(nil)

	err := m.Marshal(buf, obj)
	so.Nil(err)
	logrus.Info(buf.String())

	so.Contains(buf.String(), "STRING_TYPE_BAR")
}
