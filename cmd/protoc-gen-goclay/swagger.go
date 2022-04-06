package main

import (
	"github.com/golang/glog"
	"github.com/utrack/clay/v3/cmd/protoc-gen-goclay/third-party/grpc-gateway/internals/descriptor"
	"github.com/utrack/clay/v3/cmd/protoc-gen-goclay/third-party/grpc-gateway/protoc-gen-openapiv2/internals/genopenapi"
	plugin "google.golang.org/protobuf/types/pluginpb"
)

func genSwaggerDef(reg *descriptor.Registry, req *plugin.CodeGeneratorRequest, swaggerTitle string) (map[string][]byte, error) {
	gsw := genopenapi.New(reg, swaggerTitle)
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
	ret := make(map[string][]byte, len(outSwag))
	for pos := range outSwag {
		ret[req.FileToGenerate[pos]] = []byte(outSwag[pos].GetContent())
	}
	return ret, nil
}
