package genhandler

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"github.com/utrack/clay/v2/cmd/protoc-gen-goclay/internal"
	"google.golang.org/genproto/googleapis/api/annotations"
)

type Generator struct {
	options options
	reg     *descriptor.Registry
	imports []descriptor.GoPackage // common imports
}

// New returns a new generator which generates handler wrappers.
func New(reg *descriptor.Registry, opts ...Option) *Generator {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	g := &Generator{
		options: o,
		reg:     reg,
	}
	return g
}

func (g *Generator) newGoPackage(pkgPath string, aalias ...string) descriptor.GoPackage {
	gopkg := descriptor.GoPackage{
		Path: pkgPath,
		Name: path.Base(pkgPath),
	}
	alias := gopkg.Name
	if len(aalias) > 0 {
		alias = aalias[0]
		gopkg.Alias = alias
	}

	reference := alias
	if reference == "" {
		reference = gopkg.Name
	}

	for i := 0; ; i++ {
		if err := g.reg.ReserveGoPackageAlias(alias, gopkg.Path); err == nil {
			break
		}
		alias = fmt.Sprintf("%s_%d", gopkg.Name, i)
		gopkg.Alias = alias
	}

	if pkg == nil {
		pkg = make(map[string]string)
	}
	pkg[reference] = alias

	return gopkg
}

func (g *Generator) generateDesc(file *descriptor.File) (*plugin.CodeGeneratorResponse_File, error) {
	descCode, err := g.getDescTemplate(g.options.SwaggerDef[file.GetName()], file)

	if err != nil {
		return nil, err
	}
	formatted, err := format.Source([]byte(descCode))
	if err != nil {
		glog.Errorf("%v: %s", err, annotateString(descCode))
		return nil, err
	}
	name := filepath.Base(file.GetName())
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)

	output := fmt.Sprintf(filepath.Join(file.GoPkg.Path, "%s.pb.goclay.go"), base)
	output = filepath.Clean(output)

	glog.V(1).Infof("Will emit %s", output)

	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(output),
		Content: proto.String(string(formatted)),
	}, nil
}

func (g *Generator) generateImpl(file *descriptor.File) (files []*plugin.CodeGeneratorResponse_File, err error) {
	astPkg := astPkg(descriptor.GoPackage{
		Name: file.GoPkg.Name,
		Path: filepath.Join(file.GoPkg.Path, g.options.ImplPath),
	})
	for _, svc := range file.Services {
		if code, err := g.generateImplService(file, svc, astPkg); err == nil {
			files = append(files, code...)
		} else {
			return nil, err
		}
	}
	return files, nil
}

func (g *Generator) generateImplService(file *descriptor.File, svc *descriptor.Service, astPkg *ast.Package) ([]*plugin.CodeGeneratorResponse_File, error) {
	var files []*plugin.CodeGeneratorResponse_File

	if exists := astTypeExists(implTypeName(svc), astPkg); !exists || g.options.Force {
		output := fmt.Sprintf(filepath.Join(file.GoPkg.Path, g.options.ImplPath, "%s.pb.impl.go"), internal.SnakeCase(svc.GetName()))
		implCode, err := g.getImplTemplate(file, svc, nil)

		if err != nil {
			return nil, err
		}
		formatted, err := format.Source([]byte(implCode))
		if err != nil {
			glog.Errorf("%v: %s", err, annotateString(implCode))
			return nil, err
		}

		files = append(files, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(output),
			Content: proto.String(string(formatted)),
		})

		glog.V(1).Infof("Will emit %s", output)
	} else {
		glog.V(0).Infof("Implementation of service `%s` will not be emitted: type `%s` already exists in package `%s`", svc.GetName(), implTypeName(svc), file.GoPkg.Name)
	}

	for _, method := range svc.Methods {
		if code, err := g.generateImplServiceMethod(file, svc, method, astPkg); err == nil {
			files = append(files, code...)
		} else {
			return nil, err
		}
	}

	return files, nil
}

func (g *Generator) generateImplServiceMethod(file *descriptor.File, svc *descriptor.Service, method *descriptor.Method, astPkg *ast.Package) ([]*plugin.CodeGeneratorResponse_File, error) {
	methodGoName := goTypeName(method.GetName())
	if exists := astMethodExists(implTypeName(svc), methodGoName, astPkg); !exists || g.options.Force {
		output := fmt.Sprintf(filepath.Join(file.GoPkg.Path, g.options.ImplPath, "%s.%s.pb.impl.go"), internal.SnakeCase(svc.GetName()), internal.SnakeCase(methodGoName))
		output = filepath.Clean(output)
		implCode, err := g.getImplTemplate(file, svc, method)
		if err != nil {
			return nil, err
		}
		formatted, err := format.Source([]byte(implCode))
		if err != nil {
			glog.Errorf("%v: %s", err, annotateString(implCode))
			return nil, err
		}

		glog.V(1).Infof("Will emit %s", output)

		return []*plugin.CodeGeneratorResponse_File{{
			Name:    proto.String(output),
			Content: proto.String(string(formatted)),
		}}, nil
	}
	glog.V(0).Infof("Implementation of method `%s` for service `%s` will not be emitted: method already exists in package: `%s`", methodGoName, svc.GetName(), file.GoPkg.Name)

	return nil, nil
}

func (g *Generator) Generate(targets []*descriptor.File) ([]*plugin.CodeGeneratorResponse_File, error) {
	var files []*plugin.CodeGeneratorResponse_File
	for _, file := range targets {
		glog.V(1).Infof("Processing %s", file.GetName())

		if len(file.Services) == 0 {
			glog.V(0).Infof("%s: %v", file.GetName(), errNoTargetService)
			continue
		}

		if code, err := g.generateDesc(file); err == nil {
			files = append(files, code)
		} else {
			return nil, err
		}

		if g.options.Impl {
			if code, err := g.generateImpl(file); err == nil {
				files = append(files, code...)
			} else {
				return nil, err
			}
		}
	}

	return files, nil
}

func (g *Generator) getDescTemplate(swagger []byte, f *descriptor.File) (string, error) {
	pkgSeen := make(map[string]bool)
	var imports []descriptor.GoPackage
	for _, pkg := range g.imports {
		pkgSeen[pkg.Path] = true
		imports = append(imports, pkg)
	}

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

	if swagger != nil {
		pkgs = append(pkgs, "github.com/go-openapi/spec")
	}

	for _, pkg := range pkgs {
		pkgSeen[pkg] = true
		imports = append(imports, g.newGoPackage(pkg))
	}

	httpmw := g.newGoPackage("github.com/utrack/clay/v2/transport/httpruntime/httpmw")
	httpcli := g.newGoPackage("github.com/utrack/clay/v2/transport/httpclient")
	for _, svc := range f.Services {
		for _, m := range svc.Methods {
			checkedAppend := func(pkg descriptor.GoPackage) {
				// Add request type package to imports if needed
				if m.Options == nil || !proto.HasExtension(m.Options, annotations.E_Http) ||
					pkg == f.GoPkg || pkgSeen[pkg.Path] {
					return
				}
				pkgSeen[pkg.Path] = true
				imports = append(imports, pkg)
			}

			checkedAppend(m.RequestType.File.GoPkg)
			checkedAppend(m.ResponseType.File.GoPkg)
		}

		if hasBindings(svc) && !pkgSeen[httpcli.Path] {
			imports = append(imports, httpcli)
			pkgSeen[httpcli.Path] = true
		}

		if g.options.ApplyDefaultMiddlewares && hasBindings(svc) && !pkgSeen[httpmw.Path] {
			imports = append(imports, httpmw)
			pkgSeen[httpmw.Path] = true
		}
	}

	p := param{File: f, Imports: imports,
		ApplyMiddlewares: g.options.ApplyDefaultMiddlewares,
	}

	if swagger != nil {
		p.SwaggerBuffer = swagger
	}
	return applyDescTemplate(p)
}

func (g *Generator) getImplTemplate(f *descriptor.File, s *descriptor.Service, m *descriptor.Method) (string, error) {
	pkgSeen := make(map[string]bool)
	var imports []descriptor.GoPackage
	for _, pkg := range g.imports {
		pkgSeen[pkg.Path] = true
		imports = append(imports, pkg)
	}
	deps := make([]string, 0)
	if m == nil {
		deps = append(deps, "github.com/utrack/clay/v2/transport")
	} else {
		deps = append(deps, "context", "github.com/pkg/errors")
	}
	for _, pkg := range deps {
		pkgSeen[pkg] = true
		imports = append(imports, g.newGoPackage(pkg))
	}
	p := implParam{
		ImplGoPkgPath: f.GoPkg.Path,
		Service:       s,
		Method:        m,
		File:          f,
	}
	fileGoPkg := f.GoPkg
	if g.options.ImplPath != "" {
		rootImport := getRootImportPath(f)
		p.ImplGoPkgPath = filepath.Join(rootImport, g.options.ImplPath)
		// restore orig f.GoPkg
		defer func() {
			f.GoPkg = fileGoPkg
		}()
		// set relative f.GoPkg for proper determining package for types from desc import
		// f.GoPkg uses in function .Method.RequestType.GoType
		f.GoPkg = g.newGoPackage(rootImport, "desc")
		f.GoPkg.Name = fileGoPkg.Name
		pkgSeen[f.GoPkg.Path] = true
		imports = append(imports, f.GoPkg)
	}
	if m != nil {
		checkedAppend := func(pkg descriptor.GoPackage) {
			if m.Options == nil || !proto.HasExtension(m.Options, annotations.E_Http) ||
				pkg.Path == fileGoPkg.Path || pkgSeen[pkg.Path] {
				return
			}
			pkgSeen[pkg.Path] = true
			imports = append(imports, pkg)
		}
		checkedAppend(m.RequestType.File.GoPkg)
		checkedAppend(m.ResponseType.File.GoPkg)
	}
	p.Imports = imports

	return applyImplTemplate(p)
}

func annotateString(str string) string {
	strs := strings.Split(str, "\n")
	for pos := range strs {
		strs[pos] = fmt.Sprintf("%v: %v", pos, strs[pos])
	}
	return strings.Join(strs, "\n")
}

func getRootImportPath(file *descriptor.File) string {
	goImportPath := ""
	if file.GoPkg.Path != "." {
		goImportPath = file.GoPkg.Path
	}
	// dir is current working directory
	dir, err := filepath.Abs(".")
	if err != nil {
		glog.V(-1).Info(err)
	}
	xdir, direrr := filepath.EvalSymlinks(dir)
	for _, gp := range strings.Split(build.Default.GOPATH, ":") {
		gp = filepath.Clean(gp)
		// xgp = gp but after symlink evaluation
		xgp, gperr := filepath.EvalSymlinks(gp)
		if strings.HasPrefix(dir, gp) {
			return getPackage(dir, gp, goImportPath)
		}
		if direrr == nil && strings.HasPrefix(xdir, gp) {
			return getPackage(xdir, gp, goImportPath)
		}
		if gperr == nil && strings.HasPrefix(dir, xgp) {
			return getPackage(dir, xgp, goImportPath)
		}
		if gperr == nil && direrr == nil && strings.HasPrefix(xdir, xgp) {
			return getPackage(xdir, xgp, goImportPath)
		}
	}
	return ""
}

func getPackage(path, gopath, gopkg string) string {
	currentPath := strings.TrimPrefix(path, filepath.Join(gopath, "src")+string(filepath.Separator))
	if strings.HasPrefix(gopkg, currentPath) {
		return gopkg
	} else if gopkg != "" {
		return filepath.Join(currentPath, gopkg)
	} else {
		return currentPath
	}
}

func hasBindings(service *descriptor.Service) bool {
	for _, m := range service.Methods {
		if len(m.Bindings) > 0 {
			return true
		}
	}
	return false
}

func astPkg(pkg descriptor.GoPackage) *ast.Package {
	fileSet := token.NewFileSet()
	astPkgs, _ := parser.ParseDir(fileSet, pkg.Path, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && !strings.HasPrefix(name, ".") &&
			!strings.HasSuffix(name, "_test.go") && strings.HasSuffix(name, ".go")
	}, parser.DeclarationErrors)
	return astPkgs[pkg.Name]
}

func astTypeExists(typeName string, pkg *ast.Package) bool {
	if pkg == nil {
		return false
	}
	for _, f := range pkg.Files {
		for _, d := range f.Decls {
			if gd, ok := d.(*ast.GenDecl); ok {
				for _, s := range gd.Specs {
					if ts, ok := s.(*ast.TypeSpec); ok && ts.Name != nil && ts.Name.Name == typeName {
						return true
					}
				}
			}
		}
	}
	return false
}

func astMethodExists(typeName, methodName string, pkg *ast.Package) bool {
	if pkg == nil {
		return false
	}
	for _, f := range pkg.Files {
		for _, d := range f.Decls {
			if fd, ok := d.(*ast.FuncDecl); ok && fd.Name != nil && fd.Name.Name == methodName && fd.Recv != nil && len(fd.Recv.List) > 0 {
				if se, ok := fd.Recv.List[0].Type.(*ast.StarExpr); ok {
					if i, ok := se.X.(*ast.Ident); ok && i.Name == typeName {
						return true
					}
				}
			}
		}
	}
	return false
}
