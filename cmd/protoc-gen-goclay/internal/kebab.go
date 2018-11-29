package internal

import (
	"strings"
)

func KebabCase(s string) string {
	return strings.Trim(strings.Replace(SnakeCase(s), "_", "-", -1), "-")
}