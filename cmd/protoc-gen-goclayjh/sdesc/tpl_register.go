package sdesc

import (
	"io"
	"text/template"
)

type registerTplData struct {
	Services         []tplServiceData
	ApplyMiddlewares bool
	HasSwagger       bool
}

type tplServiceData struct {
	Name        string
	GoName      string
	HasBindings bool
	Methods     []tplMethodData
}

type tplMethodData struct {
	Name           string
	RequestGoType  string // packageName.Type if foreign type
	ResponseGoType string // packageName.Type if foreign type
	Bindings       []tplBindingData
}

type tplBindingData struct {
	Index                      int
	HTTPMethod                 string
	PathTemplate               string
	ResponseBody               bool
	ResponseBodyAssignableExpr string // Str for .Str
	HasURLPattern              bool   // has /bound/{param}/in/URL
}

func tplRegister(pc *PackageCollection, data registerTplData) func(io.Writer) error {
	tpl := template.Must(template.New("svc-reg").Funcs(template.FuncMap{
		"apkg": func(s string) string { return s }, // TODO proper de-aliasing
	}).Parse(`
{{ range $svc := .Services }}
// {{ $svc.Name }}Desc is a descriptor/registrator for the {{ $svc.GoName }}, which implements RPC {{ $svc.Name }}.
type {{ $svc.Name }}Desc struct {
	svc {{ $svc.GoName }}
	opts {{ apkg "httptransport" }}DescOptions
}

// New{{ $svc.Name }}ServiceDesc creates new registrator for the {{ $svc.GoName }}Server.
// It implements httptransport.ConfigurableServiceDesc as well.
func New{{ $svc.GoName }}Desc(svc {{ $svc.GoName }}) *{{ $svc.Name }}Desc {
	return &{{ $svc.Name }}Desc{
		svc:svc,
	}
}

// RegisterGRPC implements service registrator interface.
func (d *{{ $svc.Name }}Desc) RegisterGRPC(s *{{ apkg "grpc" }}Server) {
	Register{{ $svc.GoName }}(s,d.svc)
}

// Apply applies passed options.
func (d *{{ $svc.Name }}Desc) Apply(oo ... {{ apkg "transport" }}DescOption) {
	for _,o := range oo {
		o.Apply(&d.opts)
	}
}

// SwaggerDef returns Swagger definition for {{ $svc.Name }}.
func (d *{{ $svc.Name }}Desc) SwaggerDef(options ...{{ apkg "swagger" }}Option) (result []byte) {
	{{ if $.HasSwagger -}}
    if len(options) > 0 || len(d.opts.SwaggerDefaultOpts) > 0 {
		var err error
		var s = &{{ apkg "spec" }}Swagger{}
		if err = s.UnmarshalJSON(_swaggerDef_{{ $svc.Name }}); err != nil {
			panic("Bad swagger definition: " + err.Error())
		}

		for _, o := range d.opts.SwaggerDefaultOpts {
			o(s)
		}
		for _, o := range options {
			o(s)
		}
		if result, err = s.MarshalJSON(); err != nil {
			panic("Failed marshal {{ apkg "spec" }}Swagger definition: " + err.Error())
		}
	} else {
		result = _swaggerDef_{{ $svc.Name }}
	}
	{{ end -}}
	return result
}

// RegisterHTTP registers this service's HTTP handlers/bindings.
func (d *{{ $svc.Name }}Desc) RegisterHTTP(mux {{ apkg "transport" }}Router) {
	{{ range $m := $svc.Methods }}
	{{ range $b := $m.Bindings -}}
	// Handler for {{ $m.Name }}, binding: {{ $b.HTTPMethod }} {{ $b.PathTemplate }}
	{

		h := {{ apkg "http" }}HandlerFunc(func(w {{ apkg "http" }}ResponseWriter, r *{{ apkg "http" }}Request) {
			defer r.Body.Close()

			unmFunc := unmarshaler_goclay_{{ $svc.Name }}_{{ $m.Name }}_{{ $b.Index }}(r)
			rsp,err := _{{ $svc.Name }}_{{ $m.Name }}_Handler(d.svc,r.Context(),unmFunc,d.opts.UnaryInterceptor)
			if err != nil {
				if err,ok := err.({{ apkg "httptransport" }}MarshalerError); ok {
					{{ apkg "httpruntime" }}SetError(r.Context(),r,w,{{ apkg "errors" }}Wrap(err.Err,"couldn't parse request"))
					return
				}
				{{ apkg "httpruntime" }}SetError(r.Context(),r,w,err)
				return
			}

			if ctxErr := r.Context().Err(); ctxErr != nil && ctxErr == context.Canceled {
				w.WriteHeader(499) // Client Closed Request
				return
			}

			_,outbound := {{ apkg "httpruntime" }}MarshalerForRequest(r)
			w.Header().Set("Content-Type", outbound.ContentType())
			{{ if $b.ResponseBody -}}
			xrsp := rsp.(*{{$m.ResponseGoType }})
			err = outbound.Marshal(w, xrsp.{{ $b.ResponseBodyAssignableExpr }})
			{{ else -}}
			err = outbound.Marshal(w, rsp)
			{{ end -}}
			if err != nil {
				{{ apkg "httpruntime" }}SetError(r.Context(),r,w,{{ apkg "errors" }}Wrap(err,"couldn't write response"))
				return
			}

		}) {{/* handler func h:= end */}}

		{{ if $.ApplyMiddlewares }}
		h = httpmw.DefaultChain(h)
		{{ end -}}

		{{ if $b.HasURLPattern }}
		chiMux, isChi := mux.({{ apkg "chi" }}Router)
		if !isChi {
			panic("chi router needed to use params in URL path")
		}
		chiMux.Method("{{ $b.HTTPMethod }}",pattern_goclay_{{ $svc.Name }}_{{ $m.Name }}_{{ $b.Index }}, h)
		{{- else -}}
		mux.Handle(pattern_goclay_{{ $svc.Name }}_{{ $m.Name }}_{{ $b.Index }}, {{ apkg "http" }}HandlerFunc(func(w {{ apkg "http" }}ResponseWriter, r *{{ apkg "http" }}Request) {
			if r.Method != "{{ $b.HTTPMethod }}" {
				w.WriteHeader({{ apkg "http" }}StatusMethodNotAllowed)
				return
			}
			h(w, r)
		}))
		{{ end }}
	}
	{{ end }} {{/* bindings $b range end */}}
	{{ end }} {{/* methods $m range end */}}
}
{{ end }} {{/* services $svc range end */}}
`))
	return func(w io.Writer) error {
		return tpl.Execute(w, data)
	}
}
