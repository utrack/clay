package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/jhump/goprotoc/plugins"
	"github.com/pkg/errors"
	"github.com/utrack/clay/v2/cmd/protoc-gen-goclayjh/sdesc"
	"github.com/y0ssar1an/q"
)

func main() {
	output := os.Stdout
	os.Stdout = os.Stderr
	err := plugins.RunPlugin(os.Args[0], doCodeGen, os.Stdin, output)
	if err != nil {
		os.Exit(1)
	}
}

var (
	errNoTargetService = errors.New("no target service defined in the file")
)

func doCodeGen(req *plugins.CodeGenRequest,
	resp *plugins.CodeGenResponse) error {

	// TODO read pkgMap
	goNames := &plugins.GoNames{}
	gen := sdesc.New(goNames)
	for _, f := range req.Files {

		if len(f.GetServices()) == 0 {
			glog.V(0).Infof("%s: %v", f.GetName(), errNoTargetService)
			continue
		}

		name := filepath.Base(f.GetName())
		ext := filepath.Ext(name)
		base := strings.TrimSuffix(name, ext)

		goPkg := goNames.GoPackageForFile(f)

		output := fmt.Sprintf(filepath.Join(goPkg.ImportPath, "%s.pb.goclay.go"), base)
		output = filepath.Clean(output)
		w := resp.OutputFile(output)
		q.Q("processing", f.GetName(), output)

		err := gen.Generate(nil, f, w)
		if err != nil {
			return errors.Wrapf(err, "for %v", f.GetName())
		}
	}
	// ...
	// Process req, generate code to resp
	// ...
	return nil
}
