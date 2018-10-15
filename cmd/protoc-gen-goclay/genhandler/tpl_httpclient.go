package genhandler

import "text/template"

var clientTemplate = template.Must(template.New("http-client").Funcs(funcMap).Option().Parse(`
{{ range $svc := .Services }}
{{ if $svc | hasBindings -}}
type {{ $svc.GetName | goTypeName }}_httpClient struct {
    c *{{ pkg "http" }}Client
    host string
}

// New{{ $svc.GetName | goTypeName }}HTTPClient creates new HTTP client for {{ $svc.GetName | goTypeName }}Server.
// Pass addr in format "http://host[:port]".
func New{{ $svc.GetName | goTypeName }}HTTPClient(c *{{ pkg "http" }}Client,addr string) *{{ $svc.GetName | goTypeName }}_httpClient {
    if {{ pkg "strings" }}HasSuffix(addr, "/") {
        addr = addr[:len(addr)-1]
    }
    return &{{ $svc.GetName | goTypeName }}_httpClient{c:c,host:addr}
}
{{ end }}

{{ range $m := $svc.Methods }}
{{ if $m.Bindings }}
{{ with $b := index $m.Bindings 0 }}
func (c *{{ $svc.GetName | goTypeName }}_httpClient) {{ $m.GetName | goTypeName }}(ctx {{ pkg "context" }}Context,in *{{ $m.RequestType.GoType $m.Service.File.GoPkg.Path | goTypeName }},opts ...{{ pkg "grpc" }}CallOption) (*{{ $m.ResponseType.GoType $m.Service.File.GoPkg.Path | goTypeName }},error) {
    mw,err := {{ pkg "httpclient" }}NewMiddlewareGRPC(opts)
    if err != nil {
      return nil,err
    }

    path := pattern_goclay_{{ $svc.GetName | goTypeName }}_{{ $m.GetName }}_{{ $b.Index }}_builder(in)

    buf := {{ pkg "bytes" }}NewBuffer(nil)

    m := {{ pkg "httpruntime" }}DefaultMarshaler(nil)
    {{ if $b.Body }}
    if err = m.Marshal(buf, {{.Body.AssignableExpr "in"}}); err != nil {
	return nil, {{ pkg "errors" }}Wrap(err, "can't marshal request")
    }
    {{ end }}


    req, err := {{ pkg "http" }}NewRequest("{{ $b.HTTPMethod }}", c.host+path, buf)
    if err != nil {
        return nil, {{ pkg "errors" }}Wrap(err, "can't initiate HTTP request")
    }
    req = req.WithContext(ctx)

    req.Header.Add("Accept", m.ContentType())

    req,err = mw.ProcessRequest(req)
    if err != nil {
      return nil,err
    }
    rsp, err := c.c.Do(req)
    if err != nil {
        return nil, {{ pkg "errors" }}Wrap(err, "error from client")
    }
    defer rsp.Body.Close()

    rsp,err = mw.ProcessResponse(rsp)
    if err != nil {
      return nil,err
    }

    if rsp.StatusCode >= 400 {
        b,_ := {{ pkg "ioutil" }}ReadAll(rsp.Body)
        return nil,{{ pkg "errors" }}Errorf("%v %v: server returned HTTP %v: '%v'",req.Method,req.URL.String(),rsp.StatusCode,string(b))
    }

    ret := {{$m.ResponseType.GoType $m.Service.File.GoPkg.Path | goTypeName }}{}
    {{ if $b | ResponseBody }}
        err = m.Unmarshal(rsp.Body, &{{ .ResponseBody.AssignableExpr "ret"}})
	{{ else }}
        err = m.Unmarshal(rsp.Body, &ret)
    {{ end }}
    return &ret, {{ pkg "errors" }}Wrap(err, "can't unmarshal response")
}
{{ end }}
{{ end }}
{{ end }}
{{ end }}
`))
