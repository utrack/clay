package sdesc

import (
	"io"

	"github.com/jhump/goprotoc/plugins"
	"github.com/jhump/protoreflect/desc"
	"github.com/pkg/errors"
	"github.com/y0ssar1an/q"
)

type Generator struct {
	g  *plugins.GoNames
	pc *PackageCollection
}

func New(g *plugins.GoNames) *Generator {
	return &Generator{
		g:  g,
		pc: newPackageCollection(),
	}
}

func (g *Generator) Generate(
	swaggerDef []byte,
	f *desc.FileDescriptor,
	w io.Writer,
) error {

	filePc := g.getImportList(swaggerDef != nil, f)

	headerTplData := g.getTplDataHeader(f, filePc)
	err := tplHeader(filePc, headerTplData)(w)
	if err != nil {
		return errors.Wrap(err, "when rendering header")
	}

	regTplData := g.getTplDataRegister(f, filePc, false) // TODO ApplyMiddlewares
	err = tplRegister(filePc, regTplData)(w)
	return errors.Wrap(err, "when rendering registration")

}

func (g *Generator) getTplDataHeader(f *desc.FileDescriptor, filePc *PackageCollection) headerTplData {
	hTplData := headerTplData{
		SourceFileName: f.GetFile().GetName(),
		GoPkg:          g.g.GoPackageForFile(f).Name,
	}
	for i := range filePc.pp {
		imp := filePc.pp[i]
		q.Q("imp", imp)
		if imp.IsStdlib() {
			hTplData.StandardImports = append(hTplData.StandardImports, imp)
		} else {
			hTplData.NonStandardImports = append(hTplData.NonStandardImports, imp)
		}
	}

	// TODO sort imports
	return hTplData
}

func (g *Generator) getTplDataRegister(
	f *desc.FileDescriptor,
	filePc *PackageCollection,
	applyMw bool, // TODO conf
) registerTplData {
	ret := registerTplData{
		ApplyMiddlewares: applyMw,
	}

	for _, svc := range f.GetServices() {
		rsvc := tplServiceData{
			Name:   svc.GetName(),
			GoName: g.g.GoTypeForServiceServer(svc).Symbol().Name,
			// TODO hasBindings
		}
		for _, met := range svc.GetMethods() {
			rmet := tplMethodData{
				Name: met.GetName(),
			}
			{
				reqType := g.g.GoTypeForMessage(met.GetInputType()).Symbol()
				if reqType.Package != g.g.GoPackageForFile(f) {
					rmet.RequestGoType = filePc.Get(reqType.Package.ImportPath).Alias + "."
				}
				rmet.RequestGoType += reqType.Name
			}
			{
				rspType := g.g.GoTypeForMessage(met.GetInputType()).Symbol()
				if rspType.Package != g.g.GoPackageForFile(f) {
					rmet.ResponseGoType = filePc.Get(rspType.Package.ImportPath).Alias + "."
				}
				rmet.ResponseGoType += rspType.Name
			}
			rsvc.Methods = append(rsvc.Methods, rmet)
		}
		ret.Services = append(ret.Services, rsvc)
	}
	q.Q(ret)
	return ret
}
