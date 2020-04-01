package genhandler

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"regexp"
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
	guessModule()
	var pkg *ast.Package
	if !g.options.ServiceSubDir {
		pkg, err = astPkg(descriptor.GoPackage{
			Name: file.GoPkg.Name,
			Path: filepath.Join(file.GoPkg.Path, g.options.ImplPath),
		})
		if err != nil {
			return nil, err
		}
	}
	for _, svc := range file.Services {
		if g.options.ServiceSubDir {
			pkg, err = astPkg(descriptor.GoPackage{
				Name: file.GoPkg.Name,
				Path: filepath.Join(file.GoPkg.Path, g.options.ImplPath, internal.KebabCase(svc.GetName())),
			})
			if err != nil {
				return nil, err
			}
		}
		if code, err := g.generateImplService(file, svc, pkg); err == nil {
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
		var output string
		if g.options.ServiceSubDir {
			output = fmt.Sprintf(filepath.Join(file.GoPkg.Path, g.options.ImplPath, internal.KebabCase(svc.GetName()), "%s.go"), implFileName(svc, nil))
		} else {
			output = fmt.Sprintf(filepath.Join(file.GoPkg.Path, g.options.ImplPath, "%s.go"), implFileName(svc, nil))
		}
		implCode, err := g.getServiceImpl(file, svc)

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
		var output string
		if g.options.ServiceSubDir {
			output = fmt.Sprintf(filepath.Join(file.GoPkg.Path, g.options.ImplPath, internal.KebabCase(svc.GetName()), "%s.go"), implFileName(svc, method))
		} else {
			output = fmt.Sprintf(filepath.Join(file.GoPkg.Path, g.options.ImplPath, "%s.go"), implFileName(svc, method))
		}
		output = filepath.Clean(output)
		implCode, err := g.getMethodImpl(file, svc, method)
		if err != nil {
			return nil, err
		}
		formatted, err := format.Source([]byte(implCode))
		if err != nil {
			glog.Errorf("%v: %s", err, annotateString(implCode))
			return nil, err
		}

		glog.V(1).Infof("Will emit %s", output)

		result := []*plugin.CodeGeneratorResponse_File{{
			Name:    proto.String(output),
			Content: proto.String(string(formatted)),
		}}

		if !g.options.WithTests {
			return result, nil
		}

		testCode, err := g.getTestImpl(file, svc, method)
		if err != nil {
			return nil, err
		}

		formatted, err = format.Source([]byte(testCode))
		if err != nil {
			glog.Errorf("%v: %s", err, annotateString(testCode))
			return nil, err
		}

		result = append(result, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(strings.TrimSuffix(output, ".go") + "_test.go"),
			Content: proto.String(string(formatted)),
		})

		return result, nil
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

				// always generate alias for external packages, when types used in req/resp object
				if pkg.Alias == "" {
					pkg.Alias = pkg.Name
					pkgSeen[pkg.Path] = false
				}

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

	p := param{
		File:             f,
		Imports:          imports,
		ApplyMiddlewares: g.options.ApplyDefaultMiddlewares,
		Registry:         g.reg,
	}

	if swagger != nil {
		p.SwaggerBuffer = swagger
	}

	return applyDescTemplate(p)
}

func (g *Generator) getServiceImpl(f *descriptor.File, s *descriptor.Service) (string, error) {
	return applyImplTemplate(g.getImplParam(f, s, nil, []string{"github.com/utrack/clay/v2/transport"}))
}

func (g *Generator) getMethodImpl(f *descriptor.File, s *descriptor.Service, m *descriptor.Method) (string, error) {
	return applyImplTemplate(g.getImplParam(f, s, m, []string{"context", "github.com/pkg/errors"}))
}

func (g *Generator) getTestImpl(f *descriptor.File, s *descriptor.Service, m *descriptor.Method) (string, error) {
	return applyTestTemplate(g.getImplParam(f, s, m, []string{"context", "testing", "github.com/stretchr/testify/require"}))
}

func (g *Generator) getImplParam(f *descriptor.File, s *descriptor.Service, m *descriptor.Method, deps []string) implParam {
	pkgSeen := make(map[string]bool)
	var imports []descriptor.GoPackage
	for _, pkg := range g.imports {
		pkgSeen[pkg.Path] = true
		imports = append(imports, pkg)
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
		descImport := getDescImportPath(f)
		p.ImplGoPkgPath = filepath.Join(descImport, g.options.ImplPath)
		// restore orig f.GoPkg
		defer func() {
			f.GoPkg = fileGoPkg
		}()

		// Generate desc imports only if need
		if m != nil &&
			strings.Index(m.RequestType.File.GoPkg.Path, "/") >= 0 && !strings.HasSuffix(descImport, m.RequestType.File.GoPkg.Path) &&
			strings.Index(m.ResponseType.File.GoPkg.Path, "/") >= 0 && !strings.HasSuffix(descImport, m.ResponseType.File.GoPkg.Path) {
		} else {

			// set relative f.GoPkg for proper determining package for types from desc import
			// f.GoPkg uses in function .Method.RequestType.GoType
			f.GoPkg = g.newGoPackage(descImport, "desc")
			f.GoPkg.Name = fileGoPkg.Name
			pkgSeen[f.GoPkg.Path] = true
			imports = append(imports, f.GoPkg)
		}
	}
	if m != nil {
		checkedAppend := func(pkg descriptor.GoPackage) {
			if pkg.Path == fileGoPkg.Path || pkgSeen[pkg.Path] {
				return
			}
			pkgSeen[pkg.Path] = true

			// always generate alias for external packages, when types used in req/resp object
			if pkg.Alias == "" {
				pkg.Alias = pkg.Name
				pkgSeen[pkg.Path] = false
			}

			imports = append(imports, pkg)
		}
		checkedAppend(m.RequestType.File.GoPkg)
		checkedAppend(m.ResponseType.File.GoPkg)
	}

	p.Imports = imports
	return p
}

func annotateString(str string) string {
	strs := strings.Split(str, "\n")
	for pos := range strs {
		strs[pos] = fmt.Sprintf("%v: %v", pos, strs[pos])
	}
	return strings.Join(strs, "\n")
}

func getDescImportPath(file *descriptor.File) string {
	// wd is current working directory
	wd, err := filepath.Abs(".")
	if err != nil {
		glog.V(-1).Info(err)
	}
	// xwd = wd but after symlink evaluation
	xwd, direrr := filepath.EvalSymlinks(wd)

	// if we know module
	if module != "" {
		return getImportPath(file.GoPkg, wd, "")
	}

	for _, gp := range strings.Split(build.Default.GOPATH, ":") {
		gp = filepath.Clean(gp)
		// xgp = gp but after symlink evaluation
		xgp, gperr := filepath.EvalSymlinks(gp)
		if strings.HasPrefix(wd, gp) {
			return getImportPath(file.GoPkg, wd, gp)
		}
		if direrr == nil && strings.HasPrefix(xwd, gp) {
			return getImportPath(file.GoPkg, xwd, gp)
		}
		if gperr == nil && strings.HasPrefix(wd, xgp) {
			return getImportPath(file.GoPkg, wd, xgp)
		}
		if gperr == nil && direrr == nil && strings.HasPrefix(xwd, xgp) {
			return getImportPath(file.GoPkg, xwd, xgp)
		}
	}
	return ""
}

// getImportPath returns full go import path for specified gopkg
// wd - current working directory
// gopath - current gopath (can be empty if you are not in gopath)
func getImportPath(goPackage descriptor.GoPackage, wd, gopath string) string {
	var wdImportPath, gopkg string
	if goPackage.Path != "." {
		gopkg = goPackage.Path
	}
	if module != "" && moduleDir != "" {
		wdImportPath = filepath.Join(module, strings.TrimPrefix(wd, moduleDir))
	} else {
		wdImportPath = strings.TrimPrefix(wd, filepath.Join(gopath, "src")+string(filepath.Separator))
	}
	if strings.HasPrefix(gopkg, wdImportPath) {
		return gopkg
	} else if gopkg != "" {
		return filepath.Join(wdImportPath, gopkg)
	} else {
		return wdImportPath
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

func astPkg(pkg descriptor.GoPackage) (*ast.Package, error) {
	fileSet := token.NewFileSet()
	astPkgs, err := parser.ParseDir(fileSet, pkg.Path, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && !strings.HasPrefix(name, ".") &&
			!strings.HasSuffix(name, "_test.go") && strings.HasSuffix(name, ".go")
	}, parser.DeclarationErrors)
	if filterError(err) != nil {
		return nil, err
	}
	return astPkgs[pkg.Name], nil
}

func filterError(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*os.PathError); ok {
		return nil
	}
	return err
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

var module string
var moduleDir string
var moduleRegExp = regexp.MustCompile("^module (.*?)(?: //.*)?$")

func guessModule() {
	// dir is current working directory
	dir, err := filepath.Abs(".")
	if err != nil {
		glog.V(-1).Info(err)
	}

	// try to find go.mod
	mod := ""
	root := dir
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			mod = filepath.Join(root, "go.mod")
			break
		}
		if root == "" {
			break
		}
		d := filepath.Dir(root)
		if d == root {
			break
		}
		root = d
	}

	// if go.mod found
	if mod != "" {
		glog.V(1).Infof("Found mod file: %s", mod)
		fd, err := os.Open(mod)
		if err != nil {
			glog.V(-1).Info(err)
		}
		defer fd.Close()
		scanner := bufio.NewScanner(fd)
		for scanner.Scan() {
			line := scanner.Bytes()
			if matches := moduleRegExp.FindSubmatch(line); len(matches) > 1 {
				module = string(matches[1])
				moduleDir = root
				glog.V(1).Infof("Current module: %s", module)
				glog.V(1).Infof("Project directory: %s", moduleDir)
			}
		}
	}
}
