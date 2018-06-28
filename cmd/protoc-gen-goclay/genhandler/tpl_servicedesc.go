package genhandler

import "text/template"

var regTemplate = template.Must(template.New("svc-reg").Funcs(funcMap).Parse(`
{{ define "base" }}
{{ range $svc := .Services }}
// {{ $svc.GetName }}Desc is a descriptor/registrator for the {{ $svc.GetName }}Server.
type {{ $svc.GetName }}Desc struct {
      svc {{ $svc.GetName }}Server
}

// New{{ $svc.GetName }}ServiceDesc creates new registrator for the {{ $svc.GetName }}Server.
func New{{ $svc.GetName }}ServiceDesc(svc {{ $svc.GetName }}Server) *{{ $svc.GetName }}Desc {
      return &{{ $svc.GetName }}Desc{svc:svc}
}

// RegisterGRPC implements service registrator interface.
func (d *{{ $svc.GetName }}Desc) RegisterGRPC(s *{{ pkg "grpc" }}Server) {
      Register{{ $svc.GetName }}Server(s,d.svc)
}

// SwaggerDef returns this file's Swagger definition.
func (d *{{ $svc.GetName }}Desc) SwaggerDef(options ...{{ pkg "swagger" }}Option) (result []byte) {
    {{ if $.SwaggerBuffer }}if len(options) > 0 {
        var err error
        var s = &{{ pkg "spec" }}Swagger{}
        if err = s.UnmarshalJSON(_swaggerDef_{{ varName $.GetName }}); err != nil {
            panic("Bad swagger definition: " + err.Error())
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
func (d *{{ $svc.GetName }}Desc) RegisterHTTP(mux {{ pkg "transport" }}Router) {
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

        req, err := unmarshaler_goclay_{{ $svc.GetName }}_{{ $m.GetName }}_{{ $b.Index }}(r)
        if err != nil {
            {{ pkg "httpruntime" }}SetError(r.Context(),r,w,{{ pkg "errors" }}Wrap(err,"couldn't parse request"))
            return
        }

        ret,err := d.svc.{{ $m.GetName }}(r.Context(),req)
        if err != nil {
            {{ pkg "httpruntime" }}SetError(r.Context(),r,w,{{ pkg "errors" }}Wrap(err,"returned from handler"))
            return
        }

        _,outbound := {{ pkg "httpruntime" }}MarshalerForRequest(r)
        w.Header().Set("Content-Type", outbound.ContentType())
        err = outbound.Marshal(w, ret)
        if err != nil {
            {{ pkg "httpruntime" }}SetError(r.Context(),r,w,{{ pkg "errors" }}Wrap(err,"couldn't write response"))
            return
        }
    })

{{ if $.ApplyMiddlewares }}
    h = httpmw.DefaultChain(h)
{{ end }}

    if isChi {
        chiMux.Method("{{ $b.HTTPMethod }}",pattern_goclay_{{ $svc.GetName }}_{{ $m.GetName }}_{{ $b.Index }}, h)
    } else {
        {{ if $b.PathParams -}}
            panic("query URI params supported only for {{ pkg "chi" }}Router")
        {{- else -}}
            mux.Handle(pattern_goclay_{{ $svc.GetName }}_{{ $m.GetName }}_{{ $b.Index }}, {{ pkg "http" }}HandlerFunc(func(w {{ pkg "http" }}ResponseWriter, r *{{ pkg "http" }}Request) {
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
