package transport

import (
	// We need chi in go.mod, because generated code uses it anyway.
	_ "github.com/go-chi/chi"
)

// IsVersion2 is a static check for lib<->generator version mismatch.
const IsVersion2 = "v2"
