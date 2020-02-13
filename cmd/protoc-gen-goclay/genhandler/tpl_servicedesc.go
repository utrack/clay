package genhandler

import "text/template"

var regTemplate = template.Must(template.New("svc-reg").Funcs(funcMap).Parse(`
{{ define "base" }}
{{ range $svc := .Services }}

// RegisterHTTP registers this service's HTTP handlers/bindings.
func (d *{{ $svc.GetName | goTypeName }}Desc) RegisterHTTP(mux {{ pkg "transport" }}Router) {
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
