package main

import (
	"github.com/go-openapi/spec"
	"github.com/golang/glog"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"github.com/utrack/grpc-gateway/protoc-gen-swagger/genswagger"
)

func genSwaggerDef(req *plugin.CodeGeneratorRequest, pkgMap map[string]string) (map[string]*spec.Swagger, error) {
	reg := descriptor.NewRegistry()
	reg.SetPrefix(*importPrefix)
	reg.SetAllowDeleteBody(*allowDeleteBody)

	for k, v := range pkgMap {
		reg.AddPkgMap(k, v)
	}

	if *grpcAPIConfiguration != "" {
		if err := reg.LoadGrpcAPIServiceFromYAML(*grpcAPIConfiguration); err != nil {
			return nil, err
		}
	}

	gsw := genswagger.New(reg)

	if err := reg.Load(req); err != nil {
		return nil, err
	}

	var targets []*descriptor.File
	for _, target := range req.FileToGenerate {
		f, err := reg.LookupFile(target)
		if err != nil {
			glog.Fatal(err)
		}
		targets = append(targets, f)
	}

	outSwag, err := gsw.Generate(targets)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]*spec.Swagger, len(outSwag))
	for pos := range outSwag {
		s := &spec.Swagger{}
		if err := s.UnmarshalJSON([]byte(outSwag[pos].GetContent())); err != nil {
			return nil, err
		}
		ret[req.FileToGenerate[pos]] = s
	}
	return ret, nil
}
