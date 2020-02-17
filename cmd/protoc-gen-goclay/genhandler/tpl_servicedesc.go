package genhandler

import "text/template"

var regTemplate = template.Must(template.New("svc-reg").Funcs(funcMap).Parse(`
{{ define "base" }}
{{ range $svc := .Services }}
// {{ $svc.GetName | goTypeName }}Desc is a descriptor/registrator for the {{ $svc.GetName | goTypeName  }}Server.
type {{ $svc.GetName | goTypeName }}Desc struct {
	svc {{ $svc.GetName | goTypeName }}Server
	opts {{ pkg "httptransport" }}DescOptions
}

// New{{ $svc.GetName | goTypeName }}ServiceDesc creates new registrator for the {{ $svc.GetName | goTypeName }}Server.
// It implements httptransport.ConfigurableServiceDesc as well.
func New{{ $svc.GetName | goTypeName  }}ServiceDesc(svc {{ $svc.GetName | goTypeName }}Server) *{{ $svc.GetName | goTypeName }}Desc {
	return &{{ $svc.GetName | goTypeName  }}Desc{
		svc:svc,
	}
}

// RegisterGRPC implements service registrator interface.
func (d *{{ $svc.GetName | goTypeName }}Desc) RegisterGRPC(s *{{ pkg "grpc" }}Server) {
	Register{{ $svc.GetName | goTypeName }}Server(s,d.svc)
}

// Apply applies passed options.
func (d *{{ $svc.GetName | goTypeName }}Desc) Apply(oo ... {{ pkg "transport" }}DescOption) {
	for _,o := range oo {
		o.Apply(&d.opts)
	}
}

// SwaggerDef returns this file's Swagger definition.
func (d *{{ $svc.GetName | goTypeName }}Desc) SwaggerDef(options ...{{ pkg "swagger" }}Option) (result []byte) {
	{{ if $.SwaggerBuffer }}if len(options) > 0 || len(d.opts.SwaggerDefaultOpts) > 0 {
		var err error
		var s = &{{ pkg "spec" }}Swagger{}
		if err = s.UnmarshalJSON(_swaggerDef_{{ varName $.GetName }}); err != nil {
			panic("Bad swagger definition: " + err.Error())
		}

		for _, o := range d.opts.SwaggerDefaultOpts {
			o(s)
		}
		for _, o := range options {
			o(s)
		}
		if result, err = s.MarshalJSON(); err != nil {
			panic("Failed marshal {{ pkg "spec" }}Swagger definition: " + err.Error())
		}
	} else {
		result = _swaggerDef_{{ varName $.GetName }}
	}
	{{ end -}}
	return result
}

// RegisterHTTP registers this service's HTTP handlers/bindings.
func (d *{{ $svc.GetName | goTypeName }}Desc) RegisterHTTP(mux {{ pkg "transport" }}Router) {
	{{ if $svc | hasBindings -}}
	chiMux, isChi := mux.({{ pkg "chi" }}Router)
	{{ end }}
	{{ range $m := $svc.Methods }}
	{{ range $b := $m.Bindings -}}
	{
		// Handler for {{ $m.GetName }}, binding: {{ $b.HTTPMethod }} {{ $b.PathTmpl.Template }}
		var h http.HandlerFunc
		h = {{ pkg "http" }}HandlerFunc(func(w {{ pkg "http" }}ResponseWriter, r *{{ pkg "http" }}Request) {
			defer r.Body.Close()

			unmFunc := unmarshaler_goclay_{{ $svc.GetName | goTypeName }}_{{ $m.GetName }}_{{ $b.Index }}(r)
			rsp,err := _{{ $svc.GetName | goTypeName }}_{{ $m.GetName | goTypeName }}_Handler(d.svc,r.Context(),unmFunc,d.opts.UnaryInterceptor)

			if err != nil {
				if err,ok := err.({{ pkg "httptransport" }}MarshalerError); ok {
					{{ pkg "httpruntime" }}SetError(r.Context(),r,w,{{ pkg "errors" }}Wrap(err.Err,"couldn't parse request"))
					return
				}
				{{ pkg "httpruntime" }}SetError(r.Context(),r,w,err)
				return
			}

			if ctxErr := r.Context().Err(); ctxErr != nil && ctxErr == context.Canceled {
				w.WriteHeader(499) // Client Closed Request
				return
			}

			_,outbound := {{ pkg "httpruntime" }}MarshalerForRequest(r)
			w.Header().Set("Content-Type", outbound.ContentType())
			{{ if $b | ResponseBody -}}
			xrsp := rsp.(*{{$m.ResponseType.GoType $m.Service.File.GoPkg.Path | goTypeName }})
			err = outbound.Marshal(w, {{ $b.ResponseBody.AssignableExpr "xrsp" }})
			{{ else -}}
			err = outbound.Marshal(w, rsp)
			{{ end -}}
			if err != nil {
				{{ pkg "httpruntime" }}SetError(r.Context(),r,w,{{ pkg "errors" }}Wrap(err,"couldn't write response"))
				return
			}
		})

		{{ if $.ApplyMiddlewares }}
		h = httpmw.DefaultChain(h)
		{{ end }}

		if isChi {
			chiMux.Method("{{ $b.HTTPMethod }}",pattern_goclay_{{ $svc.GetName | goTypeName }}_{{ $m.GetName }}_{{ $b.Index }}, h)
		} else {
			{{ if $b.PathParams -}}
			panic("query URI params supported only for {{ pkg "chi" }}Router")
			{{- else -}}
			mux.Handle(pattern_goclay_{{ $svc.GetName | goTypeName }}_{{ $m.GetName }}_{{ $b.Index }}, {{ pkg "http" }}HandlerFunc(func(w {{ pkg "http" }}ResponseWriter, r *{{ pkg "http" }}Request) {
				if r.Method != "{{ $b.HTTPMethod }}" {
					w.WriteHeader({{ pkg "http" }}StatusMethodNotAllowed)
					return
				}
				h(w, r)
			}))
			{{- end }}
		}
	}
	{{ end }}
	{{ end }}
}
{{ end }}
{{ end }} // base service handler ended
`))
