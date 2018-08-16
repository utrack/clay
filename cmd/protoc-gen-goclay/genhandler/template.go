package genhandler

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	pbdescriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"github.com/pkg/errors"
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
	*descriptor.File
	Imports          []descriptor.GoPackage
	SwaggerBuffer    []byte
	ApplyMiddlewares bool
	CurrentPath      string
	Method           *descriptor.Method
}

func applyImplTemplate(p param) (string, error) {
	w := bytes.NewBuffer(nil)

	if err := implTemplate.Execute(w, p); err != nil {
		return "", err
	}

	return w.String(), nil
}

func applyDescTemplate(p param) (string, error) {
	// r := &http.Request{}
	// r.URL.Query()
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

func goFieldName(s string) string {
	toks := strings.Split(s, ".")
	for pos := range toks {
		toks[pos] = generator.CamelCase(toks[pos])
	}
	return strings.Join(toks, ".")
}

// addValueTyped returns code, adding the field value to url.Values.
// Depending on the field type, different formatters may be used.
// `range` loop is added if value is repeated.
func addValueTyped(f *descriptor.Field) string {
	isRepeated := false
	if f.GetLabel() == pbdescriptor.FieldDescriptorProto_LABEL_REPEATED {
		isRepeated = true
	}

	goName := goFieldName(f.GetName())

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
		"goFieldName":     goFieldName,
		"byteStr":         func(b []byte) string { return string(b) },
		"escapeBackTicks": func(s string) string { return strings.Replace(s, "`", "` + \"``\" + `", -1) },
		"toGoType":        func(t pbdescriptor.FieldDescriptorProto_Type) string { return primitiveTypeToGo(t) },
		// arrayToPathInterp replaces chi-style path to fmt.Sprint-style path.
		"arrayToPathInterp": func(tpl string) string {
			vv := strings.Split(tpl, "/")
			ret := []string{}
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
		"responseBodyAware": func(binding interface{}) bool {
			_, ok := binding.(interface {
				ResponseBody() *descriptor.Body
			})
			return ok
		},
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
// patterns for {{ $svc.GetName }}
var (
{{ range $m := $svc.Methods }}
{{ range $b := $m.Bindings }}

	pattern_goclay_{{ $svc.GetName }}_{{ $m.GetName }}_{{ $b.Index }} = "{{ $b.PathTmpl.Template }}"

	pattern_goclay_{{ $svc.GetName }}_{{ $m.GetName }}_{{ $b.Index }}_builder = func(in *{{$m.RequestType.GoType $m.Service.File.GoPkg.Path }}) string {
		values := url.Values{}
		{{- if not (hasBody $b) }}
			{{- range $f := $m.RequestType.Fields }}
				{{- if not (inPathParams $f $b) }}
					{{ addValueTyped $f }}
				{{- end }}
			{{- end }}
		{{- end }}

		u := url.URL{
			Path: {{ pkg "fmt" }}Sprintf("{{ arrayToPathInterp $b.PathTmpl.Template }}" {{ range $p := $b.PathParams }}, in.{{ goFieldName $p.String }}{{ end }}),
			RawQuery: values.Encode(),
		}
		return u.String()
	}

	{{ if not (hasAsterisk $b.ExplicitParams) }}
		unmarshaler_goclay_{{ $svc.GetName }}_{{ $m.GetName }}_{{ $b.Index }}_boundParams = {{ NewQueryParamFilter $b }}
	{{ end }}
{{ end }}
{{ end }}
)
{{ end }}
`))

	patternsTemplate = template.Must(template.New("patterns").Funcs(funcMap).Parse(`
{{ define "base" }}
{{ range $svc := .Services }}
// marshalers for {{ $svc.GetName }}
var (
{{ range $m := $svc.Methods }}
{{ range $b := $m.Bindings }}

    unmarshaler_goclay_{{ $svc.GetName }}_{{ $m.GetName }}_{{ $b.Index }} = func(r *{{ pkg "http" }}Request) func(interface{})(error) {
    return func(rif interface{}) error {
        req := rif.(*{{$m.RequestType.GoType $m.Service.File.GoPkg.Path }})

        {{ if not (hasAsterisk $b.ExplicitParams) }}
            if err := {{ pkg "errors" }}Wrap({{ pkg "runtime" }}PopulateQueryParameters(req, r.URL.Query(), unmarshaler_goclay_{{ $svc.GetName }}_{{ $m.GetName }}_{{ $b.Index }}_boundParams),"couldn't populate query parameters"); err != nil {
				return {{ pkg "httpruntime" }}TransformUnmarshalerError(err)
			}
        {{ end }}
        {{- if $b.Body -}}
            {{- template "unmbody" . -}}
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
    inbound,_ := {{ pkg "httpruntime" }}MarshalerForRequest(r)
    if err := {{ pkg "errors" }}Wrap(inbound.Unmarshal(r.Body,&{{.Body.AssignableExpr "req"}}),"couldn't read request JSON"); err != nil {
        return {{ pkg "httptransport" }}NewMarshalerError({{ pkg "httpruntime" }}TransformUnmarshalerError(err))
    }
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
func (i *{{ .Method.Service.GetName }}Implementation) {{ .Method.Name | goTypeName }}(ctx {{ pkg "context" }}Context, req *{{ .Method.RequestType.GoType $.CurrentPath }}) (*{{ .Method.ResponseType.GoType $.CurrentPath }}, error) {
    return nil, {{ pkg "errors" }}New("not implemented")
}
{{ else }}
{{ range $service := .Services }}
type {{ $service.GetName }}Implementation struct {}

func New{{ $service.GetName }}() *{{ $service.GetName }}Implementation {
    return &{{ $service.GetName }}Implementation{}
}
// GetDescription is a simple alias to the ServiceDesc constructor.
// It makes it possible to register the service implementation @ the server.
func (i *{{ $service.GetName }}Implementation) GetDescription() {{ pkg "transport" }}ServiceDesc {
    return {{ pkg "desc" }}New{{ $service.GetName }}ServiceDesc(i)
}
{{ end }}
{{ end }}
`))
)
