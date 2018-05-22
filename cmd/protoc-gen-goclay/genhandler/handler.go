package genhandler

import (
	"fmt"
	"go/format"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	options "google.golang.org/genproto/googleapis/api/annotations"
)

type Generator struct {
	reg         *descriptor.Registry
	baseImports []descriptor.GoPackage
}

// New returns a new generator which generates handler wrappers.
func New(reg *descriptor.Registry) *Generator {
	var imports []descriptor.GoPackage
	for _, pkgpath := range []string{
		"net/http",

		"github.com/utrack/clay/transport",
		"github.com/utrack/clay/transport/httpruntime",

		"github.com/utrack/grpc-gateway/runtime",

		"google.golang.org/grpc",
		"github.com/go-chi/chi",
		"github.com/pkg/errors",
	} {
		pkg := descriptor.GoPackage{
			Path: pkgpath,
			Name: path.Base(pkgpath),
		}
		if err := reg.ReserveGoPackageAlias(pkg.Name, pkg.Path); err != nil {
			for i := 0; ; i++ {
				alias := fmt.Sprintf("%s_%d", pkg.Name, i)
				if err := reg.ReserveGoPackageAlias(alias, pkg.Path); err != nil {
					continue
				}
				pkg.Alias = alias
				break
			}
		}
		imports = append(imports, pkg)
	}
	return &Generator{reg: reg, baseImports: imports}
}

func (g *Generator) Generate(targets []*descriptor.File, fileToSwagger map[string][]byte) ([]*plugin.CodeGeneratorResponse_File, error) {
	var files []*plugin.CodeGeneratorResponse_File
	for _, file := range targets {
		glog.V(1).Infof("Processing %s", file.GetName())
		code, err := g.getTemplate(fileToSwagger[file.GetName()], file)

		if err == errNoTargetService {
			glog.V(1).Infof("%s: %v", file.GetName(), err)
			continue
		}
		if err != nil {
			return nil, err
		}
		formatted, err := format.Source([]byte(code))
		if err != nil {

			glog.Errorf("%v: %s", err, annotateString(code))
			return nil, err
		}
		name := file.GetName()
		ext := filepath.Ext(name)
		base := strings.TrimSuffix(name, ext)
		output := fmt.Sprintf("%s.pb.goclay.go", base)
		files = append(files, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(output),
			Content: proto.String(string(formatted)),
		})
		glog.V(1).Infof("Will emit %s", output)
	}

	return files, nil
}

func (g *Generator) getTemplate(swagBuffer []byte, f *descriptor.File) (string, error) {
	if len(f.Services) == 0 {
		return "", errNoTargetService
	}
	pkgSeen := make(map[string]bool)
	var imports []descriptor.GoPackage
	for _, pkg := range g.baseImports {
		pkgSeen[pkg.Path] = true
		imports = append(imports, pkg)
	}

	for _, svc := range f.Services {
		for _, m := range svc.Methods {
			pkg := m.RequestType.File.GoPkg
			// Add request type package to imports if needed
			if m.Options == nil || !proto.HasExtension(m.Options, options.E_Http) ||
				pkg == f.GoPkg || pkgSeen[pkg.Path] {
				continue
			}
			pkgSeen[pkg.Path] = true
			imports = append(imports, pkg)
		}
	}

	return applyTemplate(param{SwagBuffer: swagBuffer, File: f, Imports: imports})
}

func annotateString(str string) string {
	strs := strings.Split(str, "\n")
	for pos := range strs {
		strs[pos] = fmt.Sprintf("%v: %v", pos, strs[pos])
	}
	return strings.Join(strs, "\n")
}
