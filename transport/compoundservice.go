package transport

import (
	"github.com/utrack/clay/v2/transport/swagger"
	"google.golang.org/grpc"
)

type CompoundServiceDesc struct {
	svc []ServiceDesc
}

func NewCompoundServiceDesc(desc ...ServiceDesc) *CompoundServiceDesc {
	return &CompoundServiceDesc{svc: desc}
}

func (d *CompoundServiceDesc) RegisterGRPC(g *grpc.Server) {
	for _, svc := range d.svc {
		svc.RegisterGRPC(g)
	}
}

func (d *CompoundServiceDesc) RegisterHTTP(r Router) {
	for _, svc := range d.svc {
		svc.RegisterHTTP(r)
	}
}

func (d *CompoundServiceDesc) SwaggerDef(options ...swagger.Option) []byte {
	j := &swagJoiner{}
	for _, svc := range d.svc {
		j.AddDefinition(svc.SwaggerDef(options...))
	}
	return j.SumDefinitions()
}

func (d *CompoundServiceDesc) Apply(oo ...DescOption) {
	for _, ss := range d.svc {
		if s, ok := ss.(ConfigurableServiceDesc); ok {
			s.Apply(oo...)
		}
	}
}
