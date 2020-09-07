package genhandler

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	pbdescriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"github.com/pkg/errors"
	"github.com/utrack/clay/v2/cmd/protoc-gen-goclay/internal"
)

const (
	nullableOption = "65001:0"
)

var (
	errNoTargetService = errors.New("no target service defined in the file")
)

var pkg map[string]string

func getPkg(name string) string {
	if p, ok := pkg[name]; ok && p != "" {
		return p + "."
	}
	return ""
}

type param struct {
	Registry *descriptor.Registry
	*descriptor.File
	Imports          []descriptor.GoPackage
	SwaggerBuffer    []byte
	ApplyMiddlewares bool
}

type implParam struct {
	*descriptor.File
	Imports       []descriptor.GoPackage
	ImplGoPkgPath string
	Method        *descriptor.Method
	Service       *descriptor.Service
}

func applyImplTemplate(p implParam) (string, error) {
	w := bytes.NewBuffer(nil)

	if err := implTemplate.Execute(w, p); err != nil {
		return "", err
	}

	return w.String(), nil
}

func applyTestTemplate(p implParam) (string, error) {
	w := bytes.NewBuffer(nil)

	if err := testTemplate.Execute(w, p); err != nil {
		return "", err
	}

	return w.String(), nil
}

func applyDescTemplate(p param) (string, error) {
	w := bytes.NewBuffer(nil)
	if err := headerTemplate.Execute(w, p); err != nil {
		return "", err
	}

	if err := regTemplate.ExecuteTemplate(w, "base", p); err != nil {
		return "", err
	}

	if err := clientTemplate.Execute(w, p); err != nil {
		return "", err
	}

	if err := marshalersTemplate.Execute(w, p); err != nil {
		return "", err
	}

	if err := patternsTemplate.ExecuteTemplate(w, "base", p); err != nil {
		return "", err
	}

	if p.SwaggerBuffer != nil {
		if err := footerTemplate.Execute(w, p); err != nil {
			return "", err
		}
	}

	return w.String(), nil
}

func goTypeName(s string) string {
	toks := strings.Split(s, ".")
	i := 0
	if len(toks) > 1 {
		i = 1
	}
	for pos := range toks[i:] {
		toks[pos+i] = generator.CamelCase(toks[pos+i])
	}
	return strings.Join(toks, ".")
}

func MustRegisterImplTypeNameTemplate(tmpl string) {
	implTypeNameTmpl = template.Must(template.New("impl-type-name").Parse(tmpl))
}

func MustRegisterImplFileNameTemplate(tmpl string) {
	implFileNameTmpl = template.Must(template.New("impl-file-name").Parse(tmpl))
}

func implTypeName(service *descriptor.Service) string {
	type params struct {
		ServiceName string
	}
	var name bytes.Buffer
	_ = implTypeNameTmpl.Execute(&name, params{ServiceName: goTypeName(service.GetName())})
	return name.String()
}

func implFileName(service *descriptor.Service, method *descriptor.Method) string {
	type params struct {
		ServiceName string
		MethodName  string
	}
	var name bytes.Buffer
	p := params{
		ServiceName: internal.SnakeCase(service.GetName()),
	}
	if method != nil {
		p.MethodName = internal.SnakeCase(goTypeName(method.GetName()))
	}
	_ = implFileNameTmpl.Execute(&name, p)
	return name.String()
}

// addValueTyped returns code, adding the field value to url.Values.
// Depending on the field type, different formatters may be used.
// `range` loop is added if value is repeated.
func addValueTyped(f *descriptor.Field) string {
	isRepeated := false
	if f.GetLabel() == pbdescriptor.FieldDescriptorProto_LABEL_REPEATED {
		isRepeated = true
	}

	goName := goTypeName(f.GetName())

	var valueFormatter string
	var valueVerb string
	switch f.GetType() {
	case pbdescriptor.FieldDescriptorProto_TYPE_BOOL:

		valueVerb = "%t"

	case pbdescriptor.FieldDescriptorProto_TYPE_DOUBLE,
		pbdescriptor.FieldDescriptorProto_TYPE_FLOAT:

		valueVerb = "%f"

	case pbdescriptor.FieldDescriptorProto_TYPE_INT64,
		pbdescriptor.FieldDescriptorProto_TYPE_UINT64,
		pbdescriptor.FieldDescriptorProto_TYPE_SINT64,
		pbdescriptor.FieldDescriptorProto_TYPE_FIXED64,
		pbdescriptor.FieldDescriptorProto_TYPE_SFIXED64,
		pbdescriptor.FieldDescriptorProto_TYPE_INT32,
		pbdescriptor.FieldDescriptorProto_TYPE_UINT32,
		pbdescriptor.FieldDescriptorProto_TYPE_SINT32,
		pbdescriptor.FieldDescriptorProto_TYPE_FIXED32,
		pbdescriptor.FieldDescriptorProto_TYPE_SFIXED32:

		valueVerb = "%d"

	case pbdescriptor.FieldDescriptorProto_TYPE_STRING,
		pbdescriptor.FieldDescriptorProto_TYPE_ENUM:

		valueVerb = "%s"

	case pbdescriptor.FieldDescriptorProto_TYPE_BYTES:

		valueFormatter = `base64.StdEncoding.EncodeToString(%s)`

	default:
		// other types are unsupported in URL Query string
	}

	if valueVerb == "" && valueFormatter == "" {
		// no way to proccess the type value, skipping
		return ""
	}

	// valueTemplater is a closure-helper for getting correct value formatter string
	valueTemplater := func(getter string) string {
		if valueFormatter != "" {
			return fmt.Sprintf(valueFormatter, getter)
		}
		return fmt.Sprintf(`fmt.Sprintf("%s", %s)`, valueVerb, getter)
	}

	if !isRepeated {
		return fmt.Sprintf(`values.Add(%q, %s)`, f.GetName(), valueTemplater("in."+goName))
	}

	format := `for _, v := range in.%s {
	values.Add(%q, %s)
}`
	return fmt.Sprintf(format, goName, f.GetName(), valueTemplater("v"))
}

var (
	varNameReplacer = strings.NewReplacer(
		".", "_",
		"/", "_",
		"-", "_",
		"!", "_",
	)
	funcMap = template.FuncMap{
		"hasAsterisk": func(ss []string) bool {
			for _, s := range ss {
				if s == "*" {
					return true
				}
			}
			return false
		},
		"varName":         func(s string) string { return varNameReplacer.Replace(s) },
		"goTypeName":      goTypeName,
		"implTypeName":    implTypeName,
		"byteStr":         func(b []byte) string { return string(b) },
		"escapeBackTicks": func(s string) string { return strings.Replace(s, "`", "` + \"``\" + `", -1) },
		"toGoType":        func(t pbdescriptor.FieldDescriptorProto_Type) string { return primitiveTypeToGo(t) },
		// arrayToPathInterp replaces chi-style path to fmt.Sprint-style path.
		"arrayToPathInterp": func(tpl string) string {
			vv := strings.Split(tpl, "/")
			var ret []string
			for _, v := range vv {
				if strings.HasPrefix(v, "{") {
					ret = append(ret, "%v")
					continue
				}
				ret = append(ret, v)
			}
			return strings.Join(ret, "/")
		},
		// returns safe package prefix with dot(.) or empty string by imported package name or alias
		"pkg":         getPkg,
		"hasBindings": hasBindings,
		"hasBody": func(b descriptor.Binding) bool {
			if b.Body != nil {
				return true
			}
			return false
		},
		"inPathParams": func(f *descriptor.Field, b descriptor.Binding) bool {
			m := map[string]bool{}
			for _, p := range b.PathParams {
				m[p.Target.GetName()] = true
			}
			return m[f.GetName()]
		},
		"addValueTyped": addValueTyped,
		"NewQueryParamFilter": func(b descriptor.Binding) string {
			var seqs [][]string
			if b.Body != nil {
				seqs = append(seqs, strings.Split(b.Body.FieldPath.String(), "."))
			}
			for _, p := range b.PathParams {
				seqs = append(seqs, strings.Split(p.FieldPath.String(), "."))
			}
			arr := utilities.NewDoubleArray(seqs)
			encodings := make([]string, len(arr.Encoding))
			for str, enc := range arr.Encoding {
				encodings[enc] = fmt.Sprintf("%q: %d", str, enc)
			}
			e := strings.Join(encodings, ", ")
			return fmt.Sprintf("&%sDoubleArray{Encoding: map[string]int{%s}, Base: %#v, Check: %#v}", getPkg("utilities"), e, arr.Base, arr.Check)
		},
		"ResponseBody": func(binding interface{}) *descriptor.Body {
			v := reflect.ValueOf(binding).Elem()
			if f := v.FieldByName("ResponseBody"); f.IsValid() {
				if body, ok := f.Interface().(*descriptor.Body); ok {
					return body
				}
			}
			return nil
		},
		// TODO move reg to map init
		"createBindingBodyTree": func(b *descriptor.Binding, reg *descriptor.Registry, assignExpr, goPkg string) []string {
			var ret []string
			for i := 0; i < len(b.Body.FieldPath); i++ {
				f := b.Body.FieldPath[i]
				if f.Target.GetType() != pbdescriptor.FieldDescriptorProto_TYPE_MESSAGE ||
					f.Target.GetLabel() == pbdescriptor.FieldDescriptorProto_LABEL_REPEATED {
					break
				}
				aExpr := b.Body.FieldPath[:i+1]
				pkg := f.Target.Message.File.GetPackage()
				fMsg, err := reg.LookupMsg(pkg, f.Target.FieldDescriptorProto.GetTypeName())
				if err != nil {
					panic(err)
				}

				t := fMsg.GoType(goPkg)
				isPointerType := true
				options := strings.Split(strings.Trim(f.Target.FieldDescriptorProto.GetOptions().String(), " "), " ")
				for _, o := range options {
					if o == nullableOption {
						isPointerType = false
						break
					}
				}
				if isPointerType {
					t = "&" + t
				}

				ret = append(ret, fmt.Sprintf("%s = %s{}", aExpr.AssignableExpr(assignExpr), t))
			}
			return ret
		},
	}

	headerTemplate = template.Must(template.New("header").Funcs(funcMap).Parse(`
// Code generated by protoc-gen-goclay. DO NOT EDIT.
// source: {{ .GetName }}

/*
Package {{ .GoPkg.Name }} is a self-registering gRPC and JSON+Swagger service definition.

It conforms to the github.com/utrack/clay/v2/transport Service interface.
*/
package {{ .GoPkg.Name }}
import (
	{{ range $i := .Imports }}{{ if $i.Standard }}{{ $i | printf "%s\n" }}{{ end }}{{ end }}

	{{ range $i := .Imports }}{{ if not $i.Standard }}{{ $i | printf "%s\n" }}{{ end }}{{ end }}
)

// Update your shared lib or downgrade generator to v1 if there's an error
var _ = {{ pkg "transport" }}IsVersion2

var _ = {{ pkg "ioutil" }}Discard
var _ {{ pkg "chi" }}Router
var _ {{ pkg "runtime" }}Marshaler
var _ {{ pkg "bytes" }}Buffer
var _ {{ pkg "context" }}Context
var _ {{ pkg "fmt" }}Formatter
var _ {{ pkg "strings" }}Reader
var _ {{ pkg "errors" }}Frame
var _ {{ pkg "httpruntime" }}Marshaler
var _ {{ pkg "http" }}Handler
var _ {{ pkg "url" }}Values
var _ {{ pkg "base64" }}Encoding
var _ {{ pkg "httptransport" }}MarshalerError
var _ {{ pkg "utilities" }}DoubleArray
`))

	footerTemplate = template.Must(template.New("footer").Funcs(funcMap).Parse(`
    var _swaggerDef_{{ varName .GetName }} = []byte(` + "`" + `{{ escapeBackTicks (byteStr .SwaggerBuffer) }}` + `
` + "`)" + `
`))

	marshalersTemplate = template.Must(template.New("patterns").Funcs(funcMap).Parse(`
{{ range $svc := .Services }}
// patterns for {{ $svc.GetName | goTypeName }}
var (
	{{ range $m := $svc.Methods }}
	{{ range $b := $m.Bindings }}

	pattern_goclay_{{ $svc.GetName | goTypeName }}_{{ $m.GetName }}_{{ $b.Index }} = "{{ $b.PathTmpl.Template }}"

	pattern_goclay_{{ $svc.GetName | goTypeName }}_{{ $m.GetName }}_{{ $b.Index }}_builder = func(in *{{ $m.RequestType.GoType $m.Service.File.GoPkg.Path | goTypeName }}) string {
		values := url.Values{}
		{{- if not (hasBody $b) }}
		{{- range $f := $m.RequestType.Fields }}
		{{- if not (inPathParams $f $b) }}
		{{ addValueTyped $f }}
		{{- end }}
		{{- end }}
		{{- end }}

		u := url.URL{
			Path: {{ pkg "fmt" }}Sprintf("{{ arrayToPathInterp $b.PathTmpl.Template }}" {{ range $p := $b.PathParams }}, {{ $p.FieldPath.AssignableExpr "in" }}{{ end }}),
			RawQuery: values.Encode(),
		}
		return u.String()
	}

	{{ if not (hasAsterisk $b.ExplicitParams) }}
	unmarshaler_goclay_{{ $svc.GetName | goTypeName }}_{{ $m.GetName }}_{{ $b.Index }}_boundParams = {{ NewQueryParamFilter $b }}
	{{ end }}
	{{ end }}
	{{ end }}
)
{{ end }}
`))

	patternsTemplate = template.Must(template.New("patterns").Funcs(funcMap).Parse(`
{{ define "base" }}
{{ range $svc := .Services }}
// marshalers for {{ $svc.GetName | goTypeName }}
var (
	{{ range $m := $svc.Methods }}
	{{ range $b := $m.Bindings }}

	unmarshaler_goclay_{{ $svc.GetName | goTypeName }}_{{ $m.GetName }}_{{ $b.Index }} = func(r *{{ pkg "http" }}Request) func(interface{})(error) {
		return func(rif interface{}) error {
			req := rif.(*{{$m.RequestType.GoType $m.Service.File.GoPkg.Path | goTypeName }})

			{{ if not (hasAsterisk $b.ExplicitParams) }}
			if err := {{ pkg "errors" }}Wrap({{ pkg "runtime" }}PopulateQueryParameters(req, r.URL.Query(), unmarshaler_goclay_{{ $svc.GetName | goTypeName }}_{{ $m.GetName }}_{{ $b.Index }}_boundParams),"couldn't populate query parameters"); err != nil {
				return {{ pkg "httpruntime" }}TransformUnmarshalerError(err)
			}
			{{ end }}
			{{- if $b.Body -}}
			{{- range $t := createBindingBodyTree $b $.Registry "req" $.GoPkg.Path }}
			{{ $t }}
			{{- end }}

			inbound,_ := {{ pkg "httpruntime" }}MarshalerForRequest(r)
			if err := {{ pkg "errors" }}Wrap(inbound.Unmarshal(r.Body,&{{.Body.AssignableExpr "req"}}),"couldn't read request JSON"); err != nil {
				return {{ pkg "httptransport" }}NewMarshalerError({{ pkg "httpruntime" }}TransformUnmarshalerError(err))
			}
			{{- end -}}
			{{- if $b.PathParams -}}
			{{- template "unmpath" . -}}
			{{ end }}
			return nil
		}
	}
	{{ end }}
	{{ end }}
)
{{ end }}
{{ end }}
{{ define "unmbody" }}
{{ end }}
{{ define "unmpath" }}
rctx := {{ pkg "chi" }}RouteContext(r.Context())
if rctx == nil {
	panic("Only chi router is supported for GETs atm")
}
for pos,k := range rctx.URLParams.Keys {
	if err := {{ pkg "errors" }}Wrapf({{ pkg "runtime" }}PopulateFieldFromPath(req, k, rctx.URLParams.Values[pos]), "can't read '%v' from path",k); err != nil {
		return {{ pkg "httptransport" }}NewMarshalerError({{ pkg "httpruntime" }}TransformUnmarshalerError(err))
	}
}
{{ end }}
`))

	implTemplate = template.Must(template.New("impl").Funcs(funcMap).Parse(`
// Code generated by protoc-gen-goclay, but your can (must) modify it.
// source: {{ .GetName }}

package  {{ .GoPkg.Name }}

import (
	{{ range $i := .Imports }}{{ if $i.Standard }}{{ $i | printf "%s\n" }}{{ end }}{{ end }}

	{{ range $i := .Imports }}{{ if not $i.Standard }}{{ $i | printf "%s\n" }}{{ end }}{{ end }}
)

{{ if .Method }}
func (i *{{ .Method.Service | implTypeName }}) {{ .Method.Name | goTypeName }}(ctx {{ pkg "context" }}Context, req *{{ .Method.RequestType.GoType $.ImplGoPkgPath | goTypeName }}) (*{{ .Method.ResponseType.GoType $.ImplGoPkgPath | goTypeName }}, error) {
	return nil, {{ pkg "errors" }}New("{{ .Method.Name | goTypeName }} not implemented")
}
{{ else }}
type {{ .Service | implTypeName}} struct {}

// New{{ .Service.GetName | goTypeName }} create new {{ .Service | implTypeName}}
func New{{ .Service.GetName | goTypeName }}() *{{ .Service | implTypeName}} {
	return &{{ .Service | implTypeName}}{}
}
// GetDescription is a simple alias to the ServiceDesc constructor.
// It makes it possible to register the service implementation @ the server.
func (i *{{ .Service | implTypeName}}) GetDescription() {{ pkg "transport" }}ServiceDesc {
	return {{ pkg "desc" }}New{{ .Service.GetName | goTypeName }}ServiceDesc(i)
}
{{ end }}
`))

	testTemplate = template.Must(template.New("test").Funcs(funcMap).Parse(`
// Code generated by protoc-gen-goclay, but your can (must) modify it.
// source: {{ .GetName }}

package  {{ .GoPkg.Name }}

import (
	{{ range $i := .Imports }}{{ if $i.Standard }}{{ $i | printf "%s\n" }}{{ end }}{{ end }}

	{{ range $i := .Imports }}{{ if not $i.Standard }}{{ $i | printf "%s\n" }}{{ end }}{{ end }}
)

func Test{{ .Method.Service | implTypeName }}_{{ .Method.Name | goTypeName }}(t *testing.T) {
	api := New{{ .Service.GetName | goTypeName }}()
	_, err := api.{{ .Method.Name | goTypeName }}({{ pkg "context" }}Background(), &{{ .Method.RequestType.GoType $.ImplGoPkgPath | goTypeName }}{})

	require.NotNil(t, err)
	require.Equal(t, "{{ .Method.Name | goTypeName }} not implemented", err.Error())
}
`))

	implTypeNameTmpl *template.Template

	implFileNameTmpl *template.Template
)
