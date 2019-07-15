package transport

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// swagJoiner glues up several Swagger definitions to one.
// This is one dirty hack...
type swagJoiner struct {
	result map[string]interface{}

	paths map[string]interface{}
	defs  map[string]interface{}
}

// AddDefinition adds another definition to the soup.
func (c *swagJoiner) AddDefinition(buf []byte) error {
	def := map[string]interface{}{}

	err := json.Unmarshal(buf, &def)
	if err != nil {
		return errors.Wrap(err, "couldn't unmarshal JSON def")
	}
	if c.result == nil {
		c.result = def
	}

	paths, _ := def["paths"].(map[string]interface{})
	structs, _ := def["definitions"].(map[string]interface{})
	if c.paths == nil {
		c.paths = paths
		c.defs = structs
		return nil
	}
	for path, sym := range paths {
		c.paths[path] = sym
	}
	for name, s := range structs {
		c.defs[name] = s
	}
	return nil
}

// SumDefinitions returns a (hopefully) valid Swagger definition combined
// from everything that came up .AddDefinition().
func (c *swagJoiner) SumDefinitions() []byte {
	if c.result == nil {
		c.result = map[string]interface{}{}
	}
	c.result["paths"] = c.paths
	c.result["definitions"] = c.defs
	ret, err := json.Marshal(c.result)
	if err != nil {
		panic(err)
	}
	return ret
}
