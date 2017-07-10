package transport

import "google.golang.org/grpc"

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

func (d *CompoundServiceDesc) SwaggerDef() []byte {
	j := &swagJoiner{}
	for _, svc := range d.svc {
		j.AddDefinition(svc.SwaggerDef())
	}
	return j.SumDefinitions()
}
