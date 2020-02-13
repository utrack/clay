package sdesc

import (
	"path"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/y0ssar1an/q"
	"google.golang.org/genproto/googleapis/api/annotations"
)

func (g *Generator) getImportList(hasSwagger bool, f *desc.FileDescriptor) *PackageCollection {
	pkgs := []string{
		"fmt",
		"io/ioutil",
		"strings",
		"bytes",
		"net/http",
		"net/url",
		"encoding/base64",
		"context",

		"github.com/utrack/clay/v2/transport/httpruntime",
		"github.com/utrack/clay/v2/transport/httptransport",
		"github.com/utrack/clay/v2/transport/swagger",
		"github.com/grpc-ecosystem/grpc-gateway/runtime",
		"github.com/grpc-ecosystem/grpc-gateway/utilities",
		"google.golang.org/grpc",
		"github.com/go-chi/chi",
		"github.com/pkg/errors",
		"github.com/utrack/clay/v2/transport",
	}
	if hasSwagger {
		pkgs = append(pkgs, "github.com/go-openapi/spec")
	}
	imp := newPackageCollection()
	for _, pkg := range pkgs {
		q.Q("addpkg", pkg)
		imp.AddOrGetPkg(g.pc.AddOrGet(pkg, ""))
	}

	localPkg := g.g.GoPackageForFile(f).ImportPath

	// import request/response packages for services
	hasBindings := false
	for _, svc := range f.GetServices() {
		for _, m := range svc.GetMethods() {
			q.Q("method")
			if m.GetOptions() == nil || !proto.HasExtension(m.GetOptions(), annotations.E_Http) {
				continue
			}
			hasBindings = true

			{
				reqPkg := g.g.GoPackageForFile(m.GetInputType().GetFile()).ImportPath
				if reqPkg != localPkg {
					imp.AddOrGetPkg(g.pc.AddOrGet(reqPkg, ""))
				}
			}
			{
				rspPkg := g.g.GoPackageForFile(m.GetOutputType().GetFile()).ImportPath
				if rspPkg != localPkg {
					imp.AddOrGetPkg(g.pc.AddOrGet(rspPkg, ""))
				}
			}
		}
	}
	if hasBindings {
		imp.AddOrGetPkg(g.pc.AddOrGet("github.com/utrack/clay/v2/transport/httpclient", ""))
	}
	// TODO if ApplyMiddlewares
	// httpmw := g.pc.AddOrGet("github.com/utrack/clay/v2/transport/httpruntime/httpmw", "")
	return imp
}

type PackageCollection struct {
	pp              map[string]Package
	reservedAliases map[string]string
}

func newPackageCollection() *PackageCollection {
	return &PackageCollection{
		pp:              map[string]Package{},
		reservedAliases: map[string]string{},
	}
}

type Package struct {
	Path  string
	Alias string
	Name  string
}

func (p Package) IsStdlib() bool {
	return !strings.Contains(p.Path, ".")
}

func (p Package) String() string {
	return p.Alias + " " + `"` + p.Path + `"`
}

func (c *PackageCollection) Seen(pkg string) bool {
	_, ret := c.pp[pkg]
	return ret
}

func (c *PackageCollection) Get(pkg string) *Package {
	p, ok := c.pp[pkg]
	if !ok {
		return nil
	}
	return &p
}

func (c *PackageCollection) AddOrGetPkg(p Package) Package {
	return c.AddOrGet(p.Path, p.Alias)
}

func (c *PackageCollection) AddOrGet(pkg string, alias string) Package {
	if ret, ok := c.pp[pkg]; ok {
		return ret
	}
	name := path.Base(pkg)
	aliasRoot := alias
	if alias == "" {
		alias = name
		aliasRoot = name
	}
	for i := 0; ; i++ {
		if _, ok := c.reservedAliases[alias]; !ok {
			break
		}
		alias = aliasRoot + "_" + strconv.Itoa(i)
	}

	c.reservedAliases[alias] = pkg
	ret := Package{
		Path:  pkg,
		Alias: alias,
		Name:  name,
	}
	c.pp[pkg] = ret
	return ret
}
