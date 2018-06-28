package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/grpc-ecosystem/grpc-gateway/codegenerator"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	"github.com/utrack/clay/cmd/protoc-gen-goclay/v2/genhandler"
)

var (
	importPrefix         = flag.String("import_prefix", "", "prefix to be added to go package paths for imported proto files")
	file                 = flag.String("file", "-", "where to load data from")
	allowDeleteBody      = flag.Bool("allow_delete_body", false, "unless set, HTTP DELETE methods may not have a body")
	grpcAPIConfiguration = flag.String("grpc_api_configuration", "", "path to gRPC API Configuration in YAML format")
	withImpl             = flag.Bool("impl", false, "generate simple implementations for proto Services. Implementation will not be generated if it already exists. See also `force` option")
	withSwagger          = flag.Bool("swagger", true, "generate swagger.json")
	applyHTTPMiddlewares = flag.Bool("http_middlewares", true, "apply default HTTP millewares")
	descPath             = flag.String("desc_path", "", "path where the http description is generated")
	implPath             = flag.String("impl_path", "", "path where the implementation is generated (for impl = true)")
	forceImpl            = flag.Bool("force", false, "force regenerate implementation if it already exists (for impl = true)")
)

func main() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	defer glog.Flush()

	reg := descriptor.NewRegistry()

	glog.V(2).Info("Processing code generator request")
	fs := os.Stdin
	if *file != "-" {
		var err error
		fs, err = os.Open(*file)
		if err != nil {
			glog.Fatal(err)
		}
	}

	glog.V(2).Info("Parsing code generator request")
	req, err := codegenerator.ParseRequest(fs)
	if err != nil {
		glog.Fatal(err)
	}

	pkgMap := make(map[string]string)
	if req.Parameter != nil {
		err = parseReqParam(req.GetParameter(), flag.CommandLine, pkgMap)
		if err != nil {
			glog.Fatalf("Error parsing flags: %v", err)
		}
	}

	reg.SetPrefix(*importPrefix)
	for k, v := range pkgMap {
		reg.AddPkgMap(k, v)
	}

	if *grpcAPIConfiguration != "" {
		if err = reg.LoadGrpcAPIServiceFromYAML(*grpcAPIConfiguration); err != nil {
			emitError(err)
			return
		}
	}

	if err = reg.Load(req); err != nil {
		emitError(err)
		return
	}

	opts := []genhandler.Option{
		genhandler.Impl(*withImpl),
		genhandler.ImplPath(*implPath),
		genhandler.DescPath(*descPath),
		genhandler.Force(*forceImpl),
		genhandler.ApplyDefaultMiddlewares(*applyHTTPMiddlewares),
	}

	if *withSwagger {
		swagBuf, err := genSwaggerDef(req, pkgMap)
		if err != nil {
			emitError(err)
			return
		}
		opts = append(opts, genhandler.SwaggerDef(swagBuf))
	}

	g := genhandler.New(reg, opts...)

	var targets []*descriptor.File
	for _, target := range req.FileToGenerate {
		var f *descriptor.File
		f, err = reg.LookupFile(target)
		if err != nil {
			glog.Fatal(err)
		}
		targets = append(targets, f)
	}

	out, err := g.Generate(targets)
	glog.V(2).Info("Processed code generator request")
	if err != nil {
		emitError(err)
		return
	}
	emitFiles(os.Stdout, out)
}

// parseReqParam parses a CodeGeneratorRequest parameter and adds the
// extracted values to the given FlagSet and pkgMap. Returns a non-nil
// error if setting a flag failed.
func parseReqParam(param string, f *flag.FlagSet, pkgMap map[string]string) error {
	if param == "" {
		return nil
	}
	for _, p := range strings.Split(param, ",") {
		spec := strings.SplitN(p, "=", 2)
		if len(spec) == 1 {
			if spec[0] == "allow_delete_body" {
				err := f.Set(spec[0], "true")
				if err != nil {
					return fmt.Errorf("Cannot set flag %s: %v", p, err)
				}
				continue
			}
			err := f.Set(spec[0], "")
			if err != nil {
				return fmt.Errorf("Cannot set flag %s: %v", p, err)
			}
			continue
		}
		name, value := spec[0], spec[1]
		if strings.HasPrefix(name, "M") {
			pkgMap[name[1:]] = value
			continue
		}
		if err := f.Set(name, value); err != nil {
			return fmt.Errorf("Cannot set flag %s: %v", p, err)
		}
	}
	*descPath = strings.Trim(*descPath, "/")
	if *descPath == "." {
		*descPath = ""
	}
	*implPath = strings.Trim(*implPath, "/")
	if *implPath == "." {
		*implPath = ""
	}
	return nil
}

func emitFiles(w io.Writer, out []*plugin.CodeGeneratorResponse_File) {
	emitResp(w, &plugin.CodeGeneratorResponse{File: out})
}

func emitError(err error) {
	emitResp(os.Stdout, &plugin.CodeGeneratorResponse{Error: proto.String(err.Error())})
}

func emitResp(out io.Writer, resp *plugin.CodeGeneratorResponse) {
	buf, err := proto.Marshal(resp)
	if err != nil {
		glog.Fatal(err)
	}
	if _, err := out.Write(buf); err != nil {
		glog.Fatal(err)
	}
}
