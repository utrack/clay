package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/utrack/clay/v2/cmd/protoc-gen-goclay/genhandler"
	"github.com/utrack/clay/v2/cmd/protoc-gen-goclay/third-party/grpc-gateway/internals/descriptor"
	"google.golang.org/protobuf/proto"
	plugin "google.golang.org/protobuf/types/pluginpb"
)

var (
	importPrefix         = flag.String("import_prefix", "", "prefix to be added to go package paths for imported proto files")
	file                 = flag.String("file", "-", "where to load data from")
	allowDeleteBody      = flag.Bool("allow_delete_body", false, "unless set, HTTP DELETE methods may not have a body")
	grpcAPIConfiguration = flag.String("grpc_api_configuration", "", "path to gRPC API Configuration in YAML format")
	withImpl             = flag.Bool("impl", false, "generate simple implementations for proto Services. Implementation will not be generated if it already exists. See also `force` option")
	withSwagger          = flag.Bool("swagger", true, "generate swagger.json")
	withSwaggerPath      = flag.String("swagger_path", "", "in addition to swagger in pb.goclay.go, generate separate swagger file at provided path")
	applyHTTPMiddlewares = flag.Bool("http_middlewares", true, "apply default HTTP millewares")
	implPath             = flag.String("impl_path", "", "path where the implementation is generated (for impl = true)")
	forceImpl            = flag.Bool("force", false, "force regenerate implementation if it already exists (for impl = true)")
	serviceSubDir        = flag.Bool("impl_service_sub_dir", false, "generate implementation for each service into sub directory")
	implTypeNameTmpl     = flag.String("impl_type_name_tmpl", "{{ .ServiceName}}Implementation", "template for generating name of implementation structure")
	implFileNameTmpl     = flag.String("impl_file_name_tmpl", "{{ if .MethodName }}{{ .MethodName }}{{ else }}{{ .ServiceName }}{{ end }}", "template for generating implementations filename")
	withTests            = flag.Bool("tests", true, "generate simple unit tests for proto Services")
	pathsParam           = flag.String("paths", "", "if you want to use source_relative instead of import which is default (see google.golang.org/protobuf@v1.27.1/compiler/protogen/protogen.go:177 for more details)")
)

func main() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	defer glog.Flush()

	// for debugging, set it to something like 5
	flag.Lookup("v").Value.Set("0")

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
	req, err := parseRequest(fs)
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

	// Remove this type cast when grpc-gateway will release new version.
	// For now this functions only present in master branch
	var xreg interface{} = reg
	if set, ok := xreg.(interface{ SetAllowRepeatedFieldsInBody(bool) }); ok {
		set.SetAllowRepeatedFieldsInBody(true)
	}
	// Use Field.GetJsonName() for generating swagger definitions
	if set, ok := xreg.(interface{ SetUseJSONNamesForFields(bool) }); ok {
		set.SetUseJSONNamesForFields(true)
	}
	reg.SetAllowDeleteBody(*allowDeleteBody)
	reg.SetPrefix(*importPrefix)
	reg.SetDisableDefaultErrors(true)
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

	genhandler.MustRegisterImplTypeNameTemplate(*implTypeNameTmpl)
	genhandler.MustRegisterImplFileNameTemplate(*implFileNameTmpl)

	opts := []genhandler.Option{
		genhandler.Impl(*withImpl),
		genhandler.ImplPath(*implPath),
		genhandler.Force(*forceImpl),
		genhandler.ServiceSubDir(*serviceSubDir),
		genhandler.ApplyDefaultMiddlewares(*applyHTTPMiddlewares),
		genhandler.WithTests(*withTests),
	}

	if *withSwagger {
		swagBuf, err := genSwaggerDef(reg, req)
		if err != nil {
			emitError(err)
			return
		}
		opts = append(opts, genhandler.SwaggerDef(swagBuf))
	}

	if *pathsParam != "" {
		switch *pathsParam {
		case genhandler.PathsParamTypeImport:
		case genhandler.PathsParamTypeSourceRelative:
			opts = append(opts, genhandler.PathsType(*pathsParam))
		default:
			emitError(fmt.Errorf("unexpected value for paths option: %q, expected one of `import`, `source_relative`", *pathsParam))
		}
	}

	if *withSwaggerPath != "" {
		opts = append(opts, genhandler.SwaggerPath(*withSwaggerPath))
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

// extracted from grpc-gateway pre-1.13.0
func parseRequest(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read code generator request: %v", err)
	}
	req := new(plugin.CodeGeneratorRequest)
	if err = proto.Unmarshal(input, req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal code generator request: %v", err)
	}
	return req, nil
}
