package strings

import (
	"bytes"
	"testing"

	"encoding/json"

	"github.com/jmoiron/jsonq"
	"github.com/stretchr/testify/assert"
	"github.com/utrack/clay/transport/httpruntime"
)

func TestMarshalUnmarshal(t *testing.T) {
	so := assert.New(t)
	obj := &String{
		SnakeCase:      "no step",
		CamelCase:      "Well Well",
		LowerCamelCase: "well Well",
	}
	m := httpruntime.DefaultMarshaler(nil)

	buf := bytes.NewBuffer(nil)

	err := m.Marshal(buf, obj)
	so.Nil(err)

	got := &String{}

	err = m.Unmarshal(buf, got)
	so.Nil(err)

	so.Equal(*obj, *got)
}

// TestFieldDefs tests that our Swagger names correspond to generated json names.
func TestFieldDefs(t *testing.T) {
	so := assert.New(t)
	obj := &String{
		SnakeCase:      "no step",
		CamelCase:      "Well Well",
		LowerCamelCase: "well Well",
	}
	m := httpruntime.DefaultMarshaler(nil)

	buf := bytes.NewBuffer(nil)

	err := m.Marshal(buf, obj)
	so.Nil(err)

	swgAccess := map[string]interface{}{}
	dec := json.NewDecoder(bytes.NewReader(_swaggerDef_pb_strings_strings_proto))
	err = dec.Decode(&swgAccess)
	so.Nil(err)
	jq := jsonq.NewQuery(swgAccess)

	o, err := jq.Object("definitions", "String", "properties")
	so.Nil(err)

	for k := range o {
		so.Contains(buf.String(), k)
	}
}
