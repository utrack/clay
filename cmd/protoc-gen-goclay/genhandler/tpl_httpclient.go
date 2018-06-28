package genhandler

import "text/template"

var clientTemplate = template.Must(template.New("http-client").Funcs(funcMap).Option().Parse(`
{{ range $svc := .Services }}
{{ if $svc | hasBindings -}}
type {{ $svc.GetName }}_httpClient struct {
    c *{{ pkg "http" }}Client
    host string
}

// New{{ $svc.GetName }}HTTPClient creates new HTTP client for {{ $svc.GetName }}Server.
// Pass addr in format "http://host[:port]".
func New{{ $svc.GetName }}HTTPClient(c *{{ pkg "http" }}Client,addr string) {{ $svc.GetName }}Client {
    if {{ pkg "strings" }}HasSuffix(addr, "/") {
        addr = addr[:len(addr)-1]
    }
    return &{{ $svc.GetName }}_httpClient{c:c,host:addr}
}
{{ end }}

{{ range $m := $svc.Methods }}
{{ if $m.Bindings }}
{{ with $b := index $m.Bindings 0 }}
func (c *{{ $svc.GetName }}_httpClient) {{ $m.GetName }}(ctx {{ pkg "context" }}Context,in *{{$m.RequestType.GoType $m.Service.File.GoPkg.Path }},_ ...{{ pkg "grpc" }}CallOption) (*{{$m.ResponseType.GoType $m.Service.File.GoPkg.Path }},error) {
    path := pattern_goclay_{{ $svc.GetName }}_{{ $m.GetName }}_{{ $b.Index }}_builder({{ range $p := $b.PathParams }}in.{{ goTypeName $p.String }},{{ end }})

    buf := {{ pkg "bytes" }}NewBuffer(nil)

    m := {{ pkg "httpruntime" }}DefaultMarshaler(nil)
    {{ if $b.Body }}
    if err := m.Marshal(buf, {{.Body.AssignableExpr "in"}}); err != nil {
	return nil, {{ pkg "errors" }}Wrap(err, "can't marshal request")
    }
    {{ end }}


    req, err := {{ pkg "http" }}NewRequest("{{ $b.HTTPMethod }}", c.host+path, buf)
    if err != nil {
        return nil, {{ pkg "errors" }}Wrap(err, "can't initiate HTTP request")
    }

    req.Header.Add("Accept", m.ContentType())

    rsp, err := c.c.Do(req)
    if err != nil {
        return nil, {{ pkg "errors" }}Wrap(err, "error from client")
    }
    defer rsp.Body.Close()

    if rsp.StatusCode >= 400 {
        b,_ := {{ pkg "ioutil" }}ReadAll(rsp.Body)
        return nil,{{ pkg "errors" }}Errorf("%v %v: server returned HTTP %v: '%v'",req.Method,req.URL.String(),rsp.StatusCode,string(b))
    }

    ret := &{{$m.ResponseType.GoType $m.Service.File.GoPkg.Path }}{}
    err = m.Unmarshal(rsp.Body, ret)
    return ret, {{ pkg "errors" }}Wrap(err, "can't unmarshal response")
}
{{ end }}
{{ end }}
{{ end }}
{{ end }}
`))
