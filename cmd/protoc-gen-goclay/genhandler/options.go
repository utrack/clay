package genhandler

// protobuf-v2 has an option to store generated pb near protofile (google.golang.org/protobuf@v1.27.1/compiler/protogen/protogen.go:432).
const (
	PathsParamTypeImport         = "import"
	PathsParamTypeSourceRelative = "source_relative"
)

type options struct {
	ImplPath                string
	DescPath                string
	SwaggerDef              map[string][]byte
	SwaggerPath             string
	Impl                    bool
	Force                   bool
	ServiceSubDir           bool
	ApplyDefaultMiddlewares bool
	ImplFileNameTmpl        string
	WithTests               bool
	PathsParamType          string
}

type Option func(*options)

// SwaggerDef sets map of swagger.json per proto file
func SwaggerDef(swaggerDef map[string][]byte) Option {
	return func(o *options) {
		o.SwaggerDef = swaggerDef
	}
}

// SwaggerPath sets path to store separate swagger files
func SwaggerPath(swaggerPath string) Option {
	return func(o *options) {
		o.SwaggerPath = swaggerPath
	}
}

// ApplyDefaultMiddlewares toggles application of httpruntime/httpmw.DefaultChain to
// every generated handler.
func ApplyDefaultMiddlewares(apply bool) Option {
	return func(o *options) {
		o.ApplyDefaultMiddlewares = apply
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

// ServiceSubDir sets impl generation into sub directory
func ServiceSubDir(subDir bool) Option {
	return func(o *options) {
		o.ServiceSubDir = subDir
	}
}

// WithTests sets if it needs for unit tests generation
func WithTests(needTests bool) Option {
	return func(o *options) {
		o.WithTests = needTests
	}
}

// PathsType is for compatibility with google.golang.org/protobuf (aka protobuf-v2).
func PathsType(pathsParamType string) Option {
	return func(o *options) {
		o.PathsParamType = pathsParamType
	}
}
