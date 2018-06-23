package genhandler

type options struct {
	ImplPath   string
	DescPath   string
	SwaggerDef map[string][]byte
	Impl       bool
	Force      bool
}

type Option func(*options)

// SwaggerDef sets map of swagger.json per proto file
func SwaggerDef(swaggerDef map[string][]byte) Option {
	return func(o *options) {
		o.SwaggerDef = swaggerDef
	}
}

// Impl sets Impl flag option (if true implementation will be generated)
func Impl(impl bool) Option {
	return func(o *options) {
		o.Impl = impl
	}
}

// ImplPath sets path for implementation file
func ImplPath(path string) Option {
	return func(o *options) {
		o.ImplPath = path
	}
}

// DescPath sets path for description and swagger file
func DescPath(path string) Option {
	return func(o *options) {
		o.DescPath = path
	}
}

// Force sets force mode for generation implementation
func Force(force bool) Option {
	return func(o *options) {
		o.Force = force
	}
}
