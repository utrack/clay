package transport

import (
	"encoding/json"

	"github.com/peterbourgon/mergemap"
	"github.com/pkg/errors"
)

// swagJoiner glues up several Swagger definitions to one.
// This is one dirty hack...
type swagJoiner struct {
	result map[string]interface{}
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
		return nil
	}
	c.result = mergemap.Merge(c.result, def)
	return nil
}

// SumDefinitions returns a (hopefully) valid Swagger definition combined
// from everything that came up .AddDefinition().
func (c *swagJoiner) SumDefinitions() []byte {
	if c.result == nil {
		c.result = map[string]interface{}{}
	}
	ret, err := json.Marshal(c.result)
	if err != nil {
		panic(err)
	}
	return ret
}
